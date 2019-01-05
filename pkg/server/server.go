package server

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/alexellis/inlets/pkg/transport"
	"github.com/gorilla/websocket"
)

// Server for the exit-node of inlets
type Server struct {
	GatewayTimeout time.Duration
	Port           int
}

// Serve traffic
func (s *Server) Serve() {
	ch := make(chan *http.Response)
	outgoing := make(chan *http.Request)

	http.HandleFunc("/ws", serveWs(ch, outgoing))
	http.HandleFunc("/", proxyHandler(ch, outgoing, s.GatewayTimeout))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.Port), nil); err != nil {
		log.Fatal(err)
	}
}

func serveWs(msg chan *http.Response, outgoing chan *http.Request) func(w http.ResponseWriter, r *http.Request) {

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			if _, ok := err.(websocket.HandshakeError); !ok {
				log.Println(err)
			}
			return
		}

		fmt.Println(ws.RemoteAddr())

		done := make(chan struct{})

		go func() {
			defer close(done)
			for {
				msgType, message, err := ws.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					return
				}

				if msgType == websocket.TextMessage {
					log.Println("TextMessage: ", message)
				} else if msgType == websocket.BinaryMessage {
					// log.Printf("Server recv: %s", message)

					reader := bytes.NewReader(message)
					scanner := bufio.NewReader(reader)
					res, _ := http.ReadResponse(scanner, nil)

					msg <- res
				}
			}
		}()

		go func() {
			defer close(done)
			for {
				fmt.Println("wait for outboundRequest")
				outboundRequest := <-outgoing

				buf := new(bytes.Buffer)

				outboundRequest.Write(buf)

				ws.WriteMessage(websocket.BinaryMessage, buf.Bytes())
			}

		}()

		<-done
	}
}

func proxyHandler(msg chan *http.Response, outgoing chan *http.Request, gatewayTimeout time.Duration) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Reverse proxy", r.Host, r.Method, r.URL.String())

		if r.Body != nil {
			defer r.Body.Close()
		}

		body, _ := ioutil.ReadAll(r.Body)
		// fmt.Println("RequestURI/Host", r.RequestURI, r.Host)
		qs := ""
		if len(r.URL.RawQuery) > 0 {
			qs = "?" + r.URL.RawQuery
		}

		req, _ := http.NewRequest(r.Method, fmt.Sprintf("http://%s%s%s", r.Host, r.URL.Path, qs),
			bytes.NewReader(body))

		transport.CopyHeaders(req.Header, &r.Header)

		outgoing <- req

		log.Println("waiting for response")

		cancel := make(chan bool)

		timeout := time.AfterFunc(gatewayTimeout, func() {
			cancel <- true
		})

		select {
		case res := <-msg:
			timeout.Stop()

			log.Println("writing response from tunnel", res.ContentLength)

			innerBody, _ := ioutil.ReadAll(res.Body)

			transport.CopyHeaders(w.Header(), &res.Header)
			w.WriteHeader(res.StatusCode)
			w.Write(innerBody)
			break
		case <-cancel:
			log.Printf("Cancelled due to timeout after %f secs\n", gatewayTimeout.Seconds())

			w.WriteHeader(http.StatusGatewayTimeout)
			break
		}

	}
}
