//go:build !wasm

package plugin

import (
	"net"

	"github.com/hashicorp/go-plugin/internal/plugin"
)

// Accept accepts a connection by ID.
//
// This should not be called multiple times with the same ID at one time.
func (b *GRPCBroker) Accept(id uint32) (net.Listener, error) {
	listener, err := serverListener()
	if err != nil {
		return nil, err
	}

	err = b.streamer.Send(&plugin.ConnInfo{
		ServiceId: id,
		Network:   listener.Addr().Network(),
		Address:   listener.Addr().String(),
	})
	if err != nil {
		return nil, err
	}

	return listener, nil
}
