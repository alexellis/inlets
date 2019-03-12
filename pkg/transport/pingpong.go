package transport

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type WebsocketConn struct {
	sync.RWMutex
	conn             *websocket.Conn
	PingWaitDuration time.Duration
}

func NewWebsocketConn(conn *websocket.Conn, pingWaitDuration time.Duration) *WebsocketConn {
	w := WebsocketConn{
		conn:             conn,
		PingWaitDuration: pingWaitDuration,
	}
	w.setupDeadline()

	return &w
}

func (w *WebsocketConn) WriteMessage(messageType int, data []byte) error {
	w.Lock()
	w.conn.SetWriteDeadline(time.Now().Add(w.PingWaitDuration))
	w.Unlock()
	return w.conn.WriteMessage(messageType, data)
}

func (w *WebsocketConn) setupDeadline() {
	log.Printf("Ping duration: %fs", w.PingWaitDuration.Seconds())
	w.conn.SetReadDeadline(time.Now().Add(w.PingWaitDuration))
	w.conn.SetPongHandler(func(string) error {
		log.Printf("PongHandler. Extend deadline.")

		newDeadline := time.Now().Add(w.PingWaitDuration)
		return w.conn.SetReadDeadline(newDeadline)
	})

	w.conn.SetPingHandler(func(string) error {
		log.Printf("PingHandler. Send pong")
		w.Lock()
		w.conn.WriteControl(websocket.PongMessage, []byte(""), time.Now().Add(time.Second))
		w.Unlock()
		return w.conn.SetReadDeadline(time.Now().Add(w.PingWaitDuration))
	})

}

func (w *WebsocketConn) LocalAddr() net.Addr {
	return w.conn.LocalAddr()
}

func (w *WebsocketConn) ReadMessage() (messageType int, p []byte, err error) {
	return w.conn.ReadMessage()
}

func (w *WebsocketConn) Ping() error {
	w.Lock()
	defer w.Unlock()
	if err := w.conn.WriteControl(websocket.PingMessage, []byte(""), time.Now().Add(time.Second)); err != nil {
		log.Println("Error sendind ping")
		return err
	}
	log.Println("Ping sent")
	return nil
}

// func (w *WebsocketConn) NextReader() (messageType int, r io.Reader, err error) {
// 	return w.conn.NextReader()
// }
