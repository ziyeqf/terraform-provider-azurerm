package metadata

import (
	"context"
	"net"
	"runtime"
)

func GetDialContext() func(ctx context.Context, network, addr string) (net.Conn, error) {
	if runtime.GOOS == "js" && runtime.GOARCH == "wasm" {
		return nil
	}

	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		d := &net.Dialer{Resolver: &net.Resolver{}}
		return d.DialContext(ctx, network, addr)
	}
}
