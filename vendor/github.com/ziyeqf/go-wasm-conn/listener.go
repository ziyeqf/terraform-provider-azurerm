package wasmconn

import (
	"net"

	"github.com/magodo/go-wasmww"
	"github.com/magodo/go-webworkers/types"
)

type Listener struct {
	connectStr      string
	postMessageFunc PostFunc
	eventChan       <-chan types.MessageEventMessage
	cancelFunc      wasmww.WebWorkerCloseFunc
	connChans       []chan WasmMsg
}

func NewListener(connectStr string, postMessageFunc PostFunc, eventChan <-chan types.MessageEventMessage, cancelFunc wasmww.WebWorkerCloseFunc) *Listener {
	return &Listener{
		connectStr,
		postMessageFunc,
		eventChan,
		cancelFunc,
		make([]chan WasmMsg, 0),
	}
}

func (w *Listener) Accept() (net.Conn, error) {
	for event := range w.eventChan {
		if data, err := event.Data(); err == nil {
			if msg, err := ParseRawMsg(data); err == nil {
				switch msg := msg.(type) {
				case *wasmConnRequest:
					if msg.ConnectStr == w.connectStr {
						connEventCh := make(chan WasmMsg, 0)
						w.connChans = append(w.connChans, connEventCh)
						conn := NewWasmConn(msg.ConnId, w.postMessageFunc, connEventCh)
						rawMsg := RawMsgFromWasmMsg(&WasmConnResponse{
							ConnId: msg.ConnId,
						})
						if err := w.postMessageFunc(rawMsg.JsMessage()); err != nil {
							return nil, err
						}
						return conn, nil
					}
				case *WasmConnMessage:
					for _, connEventCh := range w.connChans {
						connEventCh <- msg
					}
				}
				continue
			}
		}
	}
	return nil, nil
}

func (w *Listener) Close() error {
	return w.cancelFunc()
}

func (w *Listener) Addr() net.Addr {
	return NewWasmAddr(w.connectStr)
}
