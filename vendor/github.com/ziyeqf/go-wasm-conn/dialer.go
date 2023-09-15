package wasmconn

import (
	"net"

	"github.com/google/uuid"
	"github.com/magodo/go-wasmww"
	"github.com/magodo/go-webworkers/types"
)

type Dialer struct {
	connectStr string
	workerConn *wasmww.WasmWebWorkerConn
}

func NewWasmDialer(connectStr string, workerConn *wasmww.WasmWebWorkerConn) *Dialer {
	return &Dialer{
		connectStr,
		workerConn,
	}
}

func (d *Dialer) Dial() (net.Conn, error) {
	connId := uuid.New().String()

	connReceived := make(chan interface{}, 0)
	// it needs to listen before sending the request
	go func() {
		for event := range d.workerConn.EventChannel() {
			if data, err := event.Data(); err == nil {
				if msg, err := ParseRawMsg(data); err == nil {
					if resp, ok := msg.(*WasmConnResponse); ok {
						if resp.ConnId == connId {
							resp := RawMsgFromWasmMsg(&WasmConnResponse{
								ConnId: connId,
							})
							if err := d.workerConn.PostMessage(resp.JsMessage()); err != nil {
								panic(err)
							}
							connReceived <- struct{}{}
							return
						}
					}
				}
			}
		}
	}()

	connectMsg := RawMsgFromWasmMsg(&wasmConnRequest{
		ConnectStr: d.connectStr,
		ConnId:     connId,
	})

	if err := d.workerConn.PostMessage(connectMsg.JsMessage()); err != nil {
		panic(err)
	}

	<-connReceived

	return NewWasmConn(connId, d.workerConn.PostMessage, startMsgChanProxy(d.workerConn.EventChannel())), nil
}

func startMsgChanProxy(eventChan <-chan types.MessageEventMessage) <-chan WasmMsg {
	msgCh := make(chan WasmMsg, 0)
	go func() {
		for event := range eventChan {
			if data, err := event.Data(); err == nil {
				if msg, err := ParseRawMsg(data); err == nil {
					switch msg := msg.(type) {
					case *WasmConnMessage:
						msgCh <- msg
					}
				}
			}
		}
	}()
	return msgCh
}
