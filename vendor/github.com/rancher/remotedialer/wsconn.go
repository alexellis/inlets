package remotedialer

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type wsConn struct {
	sync.Mutex
	conn *websocket.Conn
}

func newWSConn(conn *websocket.Conn) *wsConn {
	w := &wsConn{
		conn: conn,
	}
	w.setupDeadline()
	return w
}

func (w *wsConn) WriteMessage(messageType int, deadline time.Time, data []byte) error {
	if deadline.IsZero() {
		w.Lock()
		defer w.Unlock()
		return w.conn.WriteMessage(messageType, data)
	}

	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		w.Lock()
		defer w.Unlock()
		done <- w.conn.WriteMessage(messageType, data)
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("i/o timeout")
	case err := <-done:
		return err
	}
}

func (w *wsConn) NextReader() (int, io.Reader, error) {
	return w.conn.NextReader()
}

func (w *wsConn) setupDeadline() {
	w.conn.SetReadDeadline(time.Now().Add(PingWaitDuration))
	w.conn.SetPingHandler(func(string) error {
		w.Lock()
		err := w.conn.WriteControl(websocket.PongMessage, []byte(""), time.Now().Add(PingWaitDuration))
		w.Unlock()
		if err != nil {
			return err
		}
		if err := w.conn.SetReadDeadline(time.Now().Add(PingWaitDuration)); err != nil {
			return err
		}
		return w.conn.SetWriteDeadline(time.Now().Add(PingWaitDuration))
	})
	w.conn.SetPongHandler(func(string) error {
		if err := w.conn.SetReadDeadline(time.Now().Add(PingWaitDuration)); err != nil {
			return err
		}
		return w.conn.SetWriteDeadline(time.Now().Add(PingWaitDuration))
	})

}
