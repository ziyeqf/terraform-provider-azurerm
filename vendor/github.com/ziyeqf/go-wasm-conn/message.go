package wasmconn

import (
	"syscall/js"

	"github.com/hack-pad/safejs"
)

type wasmConnRequest struct {
	ConnectStr string
	ConnId     string
}

func (m *wasmConnRequest) Encode() map[string]any {
	return map[string]any{
		"ConnectStr": m.ConnectStr,
		"ConnId":     m.ConnId,
	}
}

func (m *wasmConnRequest) Decode(e safejs.Value) error {
	connectStr, err := e.Get("ConnectStr")
	if err != nil {
		return err
	}
	m.ConnectStr, err = connectStr.String()
	if err != nil {
		return err
	}
	ConnId, err := e.Get("ConnId")
	if err != nil {
		return err
	}
	m.ConnId, err = ConnId.String()
	if err != nil {
		return err
	}
	return nil
}

func (m *wasmConnRequest) MsgType() string {
	return "connect"
}

type WasmConnResponse struct {
	ConnId string
}

func (m *WasmConnResponse) Encode() map[string]any {
	r := make(map[string]any)
	r["ConnId"] = m.ConnId
	return r

}

func (m *WasmConnResponse) Decode(e safejs.Value) error {
	ConnId, err := e.Get("ConnId")
	if err != nil {
		return err
	}
	m.ConnId, err = ConnId.String()
	if err != nil {
		return err
	}
	return nil
}

func (m *WasmConnResponse) MsgType() string {
	return "ack"
}

type wasmConnClose struct {
	ConnId string
}

func (m *wasmConnClose) Encode() map[string]any {
	return map[string]any{
		"ConnId": m.ConnId,
	}
}

func (m *wasmConnClose) Decode(e safejs.Value) error {
	ConnId, err := e.Get("ConnId")
	if err != nil {
		return err
	}
	m.ConnId, err = ConnId.String()
	if err != nil {
		return err
	}
	return nil
}

func (m *wasmConnClose) MsgType() string {
	return "close"
}

type WasmConnMessage struct {
	ConnId string
	Bytes  []byte
}

func (m *WasmConnMessage) Encode() map[string]any {
	jsBytesArray := js.Global().Get("Uint8Array").New(len(m.Bytes))
	js.CopyBytesToJS(jsBytesArray, m.Bytes)
	return map[string]any{
		"ConnId": m.ConnId,
		"Bytes":  jsBytesArray,
	}
}

func (m *WasmConnMessage) Decode(e safejs.Value) error {
	ConnId, err := e.Get("ConnId")
	if err != nil {
		return err
	}
	m.ConnId, err = ConnId.String()
	if err != nil {
		return err
	}
	bytes, err := e.Get("Bytes")
	if err != nil {
		return err
	}
	msgLen, err := bytes.Length()
	if err != nil {
		return err
	}
	dst := make([]byte, msgLen)
	_, err = safejs.CopyBytesToGo(dst, bytes)
	if err != nil {
		return err
	}
	m.Bytes = dst
	return nil
}

func (m *WasmConnMessage) MsgType() string {
	return "msg"
}
