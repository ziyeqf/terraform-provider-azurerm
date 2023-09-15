//go:build wasm

package plugin

import (
	"net"

	"github.com/ziyeqf/go-wasm-conn"
)

func rpcDial(c *Client) (net.Conn, error) {
	return wasmconn.NewWasmDialer(c.address.String(), c.workerConn).Dial()
}
