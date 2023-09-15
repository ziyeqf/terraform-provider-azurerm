package wasmconn

import (
	"bytes"
	"context"
	"io"
	"net"
	"time"
)

type WasmConn struct {
	readBuf         bytes.Buffer
	readChan        chan []byte
	msgChan         <-chan WasmMsg
	postMessageFunc PostFunc

	connId string
	done   bool

	readDDL, writeDDL time.Time
}

func NewWasmConn(connId string, postMessageFunc PostFunc, msgChan <-chan WasmMsg) *WasmConn {
	conn := &WasmConn{
		connId:          connId,
		readChan:        make(chan []byte, 4096),
		msgChan:         msgChan,
		postMessageFunc: postMessageFunc,
		done:            false,
	}

	go func(msgChan <-chan WasmMsg) {
		for msg := range msgChan {
			switch msg := msg.(type) {
			case *wasmConnClose:
				if msg.ConnId == conn.connId {
					conn.Close()
					return
				}
			case *WasmConnMessage:

				if msg.ConnId == conn.connId {
					conn.readChan <- msg.Bytes
					continue
				}
				if conn.done {
					return
				}
			}
		}
	}(msgChan)

	return conn
}

func (conn *WasmConn) Read(p []byte) (n int, err error) {
	if conn.done {
		return 0, io.EOF
	}

	var ctx context.Context
	var cancel context.CancelFunc
	if conn.readDDL.IsZero() {
		ctx, cancel = context.WithCancel(context.Background())
	} else {
		ctx, cancel = context.WithTimeout(context.Background(), time.Until(conn.readDDL))
	}
	defer cancel()

	type res struct {
		cp  int
		err error
	}
	cpCh := make(chan res, 1)

	go func() {
		if conn.readBuf.Len() != 0 {
			cp, err := io.ReadAtLeast(&conn.readBuf, p, min(len(p), conn.readBuf.Len()))
			cpCh <- res{cp, err}
			return
		}

		b := <-conn.readChan

		c := copy(p, b)
		if c < len(b) {
			if _, err := conn.readBuf.Write(b[c:]); err != nil {
				cpCh <- res{c, nil}
				return
			}
		}

		cpCh <- res{c, nil}
	}()

	select {
	case <-ctx.Done():
		if ctxErr := ctx.Err(); ctxErr != nil {
			return 0, ctxErr
		}
	case cp := <-cpCh:
		return cp.cp, cp.err
	}
	return 0, io.EOF
}

func (conn *WasmConn) Write(p []byte) (n int, err error) {
	if conn.done {
		return 0, io.ErrClosedPipe
	}

	var ctx context.Context
	var cancel context.CancelFunc
	if conn.readDDL.IsZero() {
		ctx, cancel = context.WithCancel(context.Background())
	} else {
		ctx, cancel = context.WithTimeout(context.Background(), time.Until(conn.readDDL))
	}
	defer cancel()

	type res struct {
		cp  int
		err error
	}
	postCh := make(chan res, 1)
	go func() {
		rawMsg := RawMsgFromWasmMsg(&WasmConnMessage{
			ConnId: conn.connId,
			Bytes:  p,
		})

		if err := conn.postMessageFunc(rawMsg.JsMessage()); err != nil {
			postCh <- res{0, err}
			return
		}
		postCh <- res{len(p), nil}
	}()

	select {
	case <-ctx.Done():
		if ctxErr := ctx.Err(); ctxErr != nil {
			return 0, ctxErr
		}
	case p := <-postCh:
		return p.cp, p.err
	}
	return 0, io.EOF
}

func (conn *WasmConn) Close() error {
	rawMsg := RawMsgFromWasmMsg(&wasmConnClose{
		ConnId: conn.connId,
	})
	conn.done = true
	if err := conn.postMessageFunc(rawMsg.JsMessage()); err != nil {
		return err
	}
	return nil
}

func (conn *WasmConn) LocalAddr() net.Addr {
	return NewWasmAddr(conn.connId)
}

func (conn *WasmConn) RemoteAddr() net.Addr {
	return NewWasmAddr(conn.connId)
}

func (conn *WasmConn) SetDeadline(t time.Time) error {
	if err := conn.SetReadDeadline(t); err != nil {
		return err
	}
	if err := conn.SetWriteDeadline(t); err != nil {
		return err
	}
	return nil
}

func (conn *WasmConn) SetReadDeadline(t time.Time) error {
	conn.readDDL = t
	return nil
}

func (conn *WasmConn) SetWriteDeadline(t time.Time) error {
	conn.writeDDL = t
	return nil
}
