// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build wasm

package plugin

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/magodo/chanio"
	"github.com/magodo/go-wasmww"
	"github.com/ziyeqf/go-wasm-conn"

	"google.golang.org/grpc"
)

// CoreProtocolVersion is the ProtocolVersion of the plugin system itself.
// We will increment this whenever we change any protocol behavior. This
// will invalidate any prior plugins but will at least allow us to iterate
// on the core in a safe way. We will do our best to do this very
// infrequently.
const CoreProtocolVersion = 1

// HandshakeConfig is the configuration used by client and servers to
// handshake before starting a plugin connection. This is embedded by
// both ServeConfig and ClientConfig.
//
// In practice, the plugin host creates a HandshakeConfig that is exported
// and plugins then can easily consume it.
type HandshakeConfig struct {
	// ProtocolVersion is the version that clients must match on to
	// agree they can communicate. This should match the ProtocolVersion
	// set on ClientConfig when using a plugin.
	// This field is not required if VersionedPlugins are being used in the
	// Client or Server configurations.
	ProtocolVersion uint

	// MagicCookieKey and value are used as a very basic verification
	// that a plugin is intended to be launched. This is not a security
	// measure, just a UX feature. If the magic cookie doesn't match,
	// we show human-friendly output.
	MagicCookieKey   string
	MagicCookieValue string
}

// PluginSet is a set of plugins provided to be registered in the plugin
// server.
type PluginSet map[string]Plugin

// ServeConfig configures what sorts of plugins are served.
type ServeConfig struct {
	// HandshakeConfig is the configuration that must match clients.
	HandshakeConfig

	// TLSProvider is a function that returns a configured tls.Config.
	TLSProvider func() (*tls.Config, error)

	// Plugins are the plugins that are served.
	// The implied version of this PluginSet is the Handshake.ProtocolVersion.
	Plugins PluginSet

	// VersionedPlugins is a map of PluginSets for specific protocol versions.
	// These can be used to negotiate a compatible version between client and
	// server. If this is set, Handshake.ProtocolVersion is not required.
	VersionedPlugins map[int]PluginSet

	// GRPCServer should be non-nil to enable serving the plugins over
	// gRPC. This is a function to create the server when needed with the
	// given server options. The server options populated by go-plugin will
	// be for TLS if set. You may modify the input slice.
	//
	// Note that the grpc.Server will automatically be registered with
	// the gRPC health checking service. This is not optional since go-plugin
	// relies on this to implement Ping().
	GRPCServer func([]grpc.ServerOption) *grpc.Server

	// Logger is used to pass a logger into the server. If none is provided the
	// server will create a default logger.
	Logger hclog.Logger

	// Test, if non-nil, will put plugin serving into "test mode". This is
	// meant to be used as part of `go test` within a plugin's codebase to
	// launch the plugin in-process and output a ReattachConfig.
	//
	// This changes the behavior of the server in a number of ways to
	// accomodate the expectation of running in-process:
	//
	//   * The handshake cookie is not validated.
	//   * Stdout/stderr will receive plugin reads and writes
	//   * Connection information will not be sent to stdout
	//
	Test *ServeTestConfig

	WASMConnectStr string
}

// ServeTestConfig configures plugin serving for test mode. See ServeConfig.Test.
type ServeTestConfig struct {
	// Context, if set, will force the plugin serving to end when cancelled.
	// This is only a test configuration because the non-test configuration
	// expects to take over the process and therefore end on an interrupt or
	// kill signal. For tests, we need to kill the plugin serving routinely
	// and this provides a way to do so.
	//
	// If you want to wait for the plugin process to close before moving on,
	// you can wait on CloseCh.
	Context context.Context

	// If this channel is non-nil, we will send the ReattachConfig via
	// this channel. This can be encoded (via JSON recommended) to the
	// plugin client to attach to this plugin.
	ReattachConfigCh chan<- *ReattachConfig

	// CloseCh, if non-nil, will be closed when serving exits. This can be
	// used along with Context to determine when the server is fully shut down.
	// If this is not set, you can still use Context on its own, but note there
	// may be a period of time between canceling the context and the plugin
	// server being shut down.
	CloseCh chan<- struct{}

	// SyncStdio, if true, will enable the client side "SyncStdout/Stderr"
	// functionality to work. This defaults to false because the implementation
	// of making this work within test environments is particularly messy
	// and SyncStdio functionality is fairly rare, so we default to the simple
	// scenario.
	SyncStdio bool
}

// protocolVersion determines the protocol version and plugin set to be used by
// the server. In the event that there is no suitable version, the last version
// in the config is returned leaving the client to report the incompatibility.
func protocolVersion(opts *ServeConfig) (int, Protocol, PluginSet) {
	protoVersion := int(opts.ProtocolVersion)
	pluginSet := opts.Plugins
	protoType := ProtocolNetRPC
	// Check if the client sent a list of acceptable versions
	var clientVersions []int
	if vs := os.Getenv("PLUGIN_PROTOCOL_VERSIONS"); vs != "" {
		for _, s := range strings.Split(vs, ",") {
			v, err := strconv.Atoi(s)
			if err != nil {
				fmt.Fprintf(os.Stderr, "server sent invalid plugin version %q", s)
				continue
			}
			clientVersions = append(clientVersions, v)
		}
	}

	// We want to iterate in reverse order, to ensure we match the newest
	// compatible plugin version.
	sort.Sort(sort.Reverse(sort.IntSlice(clientVersions)))

	// set the old un-versioned fields as if they were versioned plugins
	if opts.VersionedPlugins == nil {
		opts.VersionedPlugins = make(map[int]PluginSet)
	}

	if pluginSet != nil {
		opts.VersionedPlugins[protoVersion] = pluginSet
	}

	// Sort the version to make sure we match the latest first
	var versions []int
	for v := range opts.VersionedPlugins {
		versions = append(versions, v)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(versions)))

	// See if we have multiple versions of Plugins to choose from
	for _, version := range versions {
		// Record each version, since we guarantee that this returns valid
		// values even if they are not a protocol match.
		protoVersion = version
		pluginSet = opts.VersionedPlugins[version]

		// If we have a configured gRPC server we should select a protocol
		if opts.GRPCServer != nil {
			// All plugins in a set must use the same transport, so check the first
			// for the protocol type
			for _, p := range pluginSet {
				switch p.(type) {
				case GRPCPlugin:
					protoType = ProtocolGRPC
				default:
					protoType = ProtocolNetRPC
				}
				break
			}
		}

		for _, clientVersion := range clientVersions {
			if clientVersion == protoVersion {
				return protoVersion, protoType, pluginSet
			}
		}
	}

	// Return the lowest version as the fallback.
	// Since we iterated over all the versions in reverse order above, these
	// values are from the lowest version number plugins (which may be from
	// a combination of the Handshake.ProtocolVersion and ServeConfig.Plugins
	// fields). This allows serving the oldest version of our plugins to a
	// legacy client that did not send a PLUGIN_PROTOCOL_VERSIONS list.
	return protoVersion, protoType, pluginSet
}

// Serve serves the plugins given by ServeConfig.
//
// Serve doesn't return until the plugin is done being executed. Any
// fixable errors will be output to os.Stderr and the process will
// exit with a status code of 1. Serve will panic for unexpected
// conditions where a user's fix is unknown.
//
// This is the method that plugins should call in their main() functions.
func Serve(opts *ServeConfig) {
	exitCode := -1
	// We use this to trigger an `os.Exit` so that we can execute our other
	// deferred functions. In test mode, we just output the err to stderr
	// and return.
	defer func() {
		if opts.Test == nil && exitCode >= 0 {
			os.Exit(exitCode)
		}

		if opts.Test != nil && opts.Test.CloseCh != nil {
			close(opts.Test.CloseCh)
		}
	}()

	if opts.Test == nil {
		// Validate the handshake config
		if opts.MagicCookieKey == "" || opts.MagicCookieValue == "" {
			fmt.Fprintf(os.Stderr,
				"Misconfigured ServeConfig given to serve this plugin: no magic cookie\n"+
					"key or value was set. Please notify the plugin author and report\n"+
					"this as a bug.\n")
			exitCode = 1
			return
		}

		// First check the cookie
		if os.Getenv(opts.MagicCookieKey) != opts.MagicCookieValue {
			fmt.Fprintf(os.Stderr,
				"This binary is a plugin. These are not meant to be executed directly.\n"+
					"Please execute the program that consumes these plugins, which will\n"+
					"load any plugins automatically\n")
			exitCode = 1
			return
		}
	}

	if opts.WASMConnectStr == "" {
		fmt.Fprintf(os.Stderr, "WASMConnectStr must be set.")
		exitCode = 1
		return
	}

	protoVersion, protoType, pluginSet := protocolVersion(opts)

	logger := opts.Logger
	if logger == nil {
		// internal logger to os.Stderr
		logger = hclog.New(&hclog.LoggerOptions{
			Level:      hclog.Trace,
			Output:     os.Stderr,
			JSONFormat: true,
		})
	}

	jsSelf, err := wasmww.NewSelfConn()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting self conn: %s\n", err)
	}

	// Register a listener so we can accept a connection
	listener, err := serverListener(opts.WASMConnectStr, jsSelf)
	if err != nil {
		logger.Error("plugin init error", "error", err)
		return
	}

	// Close the listener on return. We wrap this in a func() on purpose
	// because the "listener" reference may change to TLS.
	defer func() {
		listener.Close()
	}()

	// Create the channel to tell us when we're done
	doneCh := make(chan struct{})

	// Create our new stdout, stderr files. These will override our built-in
	// stdout/stderr so that it works across the stream boundary.
	stderr_r, stderr_w, err := chanio.Pipe()
	stdout_r, stdout_w, err := chanio.Pipe()

	// Build the server type
	var server ServerProtocol
	switch protoType {
	case ProtocolNetRPC:
		// Create the RPC server to dispense
		server = &RPCServer{
			Plugins: pluginSet,
			Stdout:  stderr_r,
			Stderr:  stdout_r,
			DoneCh:  doneCh,
		}
	case ProtocolGRPC:
		server = &GRPCServer{
			Plugins: pluginSet,
			Server:  opts.GRPCServer,
			//TLS:     tlsConfig,
			Stdout: stdout_r,
			Stderr: stderr_r,
			DoneCh: doneCh,
			logger: logger,
		}
	default:
		jsSelf.ResetWriteSync()
		panic("unknown server protocol: " + protoType)
	}

	// Initialize the servers
	if err := server.Init(); err != nil {
		logger.Error("protocol init", "error", err)
		return
	}

	logger.Debug("plugin address", "network", listener.Addr().Network(), "address", listener.Addr().String())

	fmt.Printf("%d|%d|%s|%s|%s\n",
		CoreProtocolVersion,
		protoVersion,
		listener.Addr().Network(),
		listener.Addr().String(),
		protoType,
	)
	os.Stdout.Sync()

	wasmww.SetWriteSync(
		[]wasmww.MsgWriter{
			wasmww.NewMsgWriterToIoWriter(stdout_w),
		},
		[]wasmww.MsgWriter{
			wasmww.NewMsgWriterToIoWriter(stderr_w),
		},
	)
	go server.Serve(listener)

	ctx := context.Background()
	if opts.Test != nil && opts.Test.Context != nil {
		ctx = opts.Test.Context
	}
	select {
	case <-ctx.Done():
		// Cancellation. We can stop the server by closing the listener.
		// This isn't graceful at all but this is currently only used by
		// tests and its our only way to stop.
		listener.Close()

		// If this is a grpc server, then we also ask the server itself to
		// end which will kill all connections. There isn't an easy way to do
		// this for net/rpc currently but net/rpc is more and more unused.
		//if s, ok := server.(*GRPCServer); ok {
		//	s.Stop()
		//}

		// Wait for the server itself to shut down
		<-doneCh

	case <-doneCh:
		// Note that given the documentation of Serve we should probably be
		// setting exitCode = 0 and using os.Exit here. That's how it used to
		// work before extracting this library. However, for years we've done
		// this so we'll keep this functionality.
	}
}

func serverListener(connectStr string, jsSelf *wasmww.SelfConn) (net.Listener, error) {
	eventCh, err := jsSelf.SetupConn()
	if err != nil {
		return nil, fmt.Errorf("Error setting up wasm conn: %s\n", err)
	}

	return wasmconn.NewListener(connectStr, jsSelf.PostMessage, eventCh, jsSelf.Close), nil
}
