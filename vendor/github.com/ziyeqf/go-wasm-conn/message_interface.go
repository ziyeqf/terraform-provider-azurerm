package wasmconn

import (
	"fmt"

	"github.com/hack-pad/safejs"
)

type RawMsg struct {
	msgType string
	extra   map[string]any
}

func (m *RawMsg) JsMessage() (safejs.Value, []safejs.Value) {
	jsMsg := make(map[string]any)
	jsMsg["msgType"] = m.msgType
	jsMsg["extra"] = m.extra
	t, err := safejs.ValueOf(jsMsg)
	if err != nil {
		panic(err)
	}
	return t, nil
}

func ParseRawMsg(jsMsg safejs.Value) (WasmMsg, error) {
	t, err := jsMsg.Get("msgType")
	if err != nil {
		return nil, err
	}
	e, err := jsMsg.Get("extra")
	if err != nil {
		return nil, err
	}
	tStr, err := t.String()
	if err != nil {
		return nil, fmt.Errorf("failed to parse msgType: %v", t.Type())
	}

	var m WasmMsg
	switch tStr {
	case "connect":
		m = &wasmConnRequest{}
		if err := m.Decode(e); err != nil {
			return nil, err
		}
	case "ack":
		m = &WasmConnResponse{}
		if err := m.Decode(e); err != nil {
			return nil, err
		}
	case "close":
		m = &wasmConnClose{}
		if err := m.Decode(e); err != nil {
			return nil, err
		}
	case "msg":
		m = &WasmConnMessage{}
		if err := m.Decode(e); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown msgType: %q", tStr)
	}

	return m, nil
}

func RawMsgFromWasmMsg(msg WasmMsg) RawMsg {
	return RawMsg{
		msgType: msg.MsgType(),
		extra:   msg.Encode(),
	}
}
