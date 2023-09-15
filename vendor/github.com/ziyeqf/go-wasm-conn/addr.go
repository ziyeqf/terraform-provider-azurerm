package wasmconn

type WasmAddr struct {
	connectStr string
}

func NewWasmAddr(connectStr string) *WasmAddr {
	return &WasmAddr{
		connectStr: connectStr,
	}
}

func (w WasmAddr) Network() string {
	return "wasm"
}

func (w WasmAddr) String() string {
	return w.connectStr
}
