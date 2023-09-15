//go:build !wasm

package plugin

import (
	"net"
)

func rpcDial(c *Client) (net.Conn, error) {
	return net.Dial(c.address.Network(), c.address.String())
}
