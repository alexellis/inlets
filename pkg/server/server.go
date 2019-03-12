package server

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/alexellis/inlets/pkg/transport"
	"github.com/alexellis/inlets/pkg/types"
	"github.com/gorilla/websocket"
	"github.com/twinj/uuid"
)

// Server for the exit-node of inlets
type Server struct {
	GatewayTimeout time.Duration
	Port           int
	Token          string
}

// Serve traffic
func (s *Server) Serve() {
	outgoing := make(chan *http.Request)

	bus := types.NewBus()

	http.HandleFunc("/", proxyHandler(outgoing, bus, s.GatewayTimeout))
	http.HandleFunc("/tunnel", serveWs(outgoing, bus, s.Token))

	collectInterval := time.Second * 10
	go garbageCollectBus(bus, collectInterval, s.GatewayTimeout*2)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.Port), nil); err != nil {
		log.Fatal(err)
	}
}

func garbageCollectBus(bus *types.Bus, interval time.Duration, expiry time.Duration) {
	ticker := time.NewTicker(interval)
	select {
	case <-ticker.C:
		list := bus.SubscriptionList()
		for _, item := range list {
			if bus.Expired(item, expiry) {
				bus.Unsubscribe(item)
			}
		}
		break
	}
}

func proxyHandler(outgoing chan *http.Request, bus *types.Bus, gatewayTimeout time.Duration) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		inletsID := uuid.Formatter(uuid.NewV4(), uuid.FormatHex)

		sub := bus.Subscribe(inletsID)

		defer func() {
			bus.Unsubscribe(inletsID)
		}()

		log.Printf("[%s] proxy %s %s %s", inletsID, r.Host, r.Method, r.URL.String())
		r.Header.Set(transport.InletsHeader, inletsID)

		if r.Body != nil {
			defer r.Body.Close()
		}

		body, _ := ioutil.ReadAll(r.Body)

		qs := ""
		if len(r.URL.RawQuery) > 0 {
			qs = "?" + r.URL.RawQuery
		}

		req, _ := http.NewRequest(r.Method, fmt.Sprintf("http://%s%s%s", r.Host, r.URL.Path, qs),
			bytes.NewReader(body))

		transport.CopyHeaders(req.Header, &r.Header)

		wg := sync.WaitGroup{}
		wg.Add(2)

		go func() {
			log.Printf("[%s] waiting for response", inletsID)

			select {
			case res := <-sub.Data:

				if res != nil && res.Body != nil {
					innerBody, _ := ioutil.ReadAll(res.Body)

					transport.CopyHeaders(w.Header(), &res.Header)
					w.WriteHeader(res.StatusCode)
					w.Write(innerBody)
					log.Printf("[%s] wrote %d bytes", inletsID, len(innerBody))
				} else {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Request or body from request was nil, client may have disconnected"))
					log.Printf("Request or body from request was nil, client may have disconnected")
				}

				wg.Done()
				break
			case <-time.After(gatewayTimeout):
				log.Printf("[%s] timeout after %f secs\n", inletsID, gatewayTimeout.Seconds())

				w.WriteHeader(http.StatusGatewayTimeout)
				wg.Done()
				break
			}
		}()

		go func() {
			outgoing <- req
			wg.Done()
		}()

		wg.Wait()
	}
}

func serveWs(outgoing chan *http.Request, bus *types.Bus, token string) func(w http.ResponseWriter, r *http.Request) {

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		err := authorized(token, r)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
		}

		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			if _, ok := err.(websocket.HandshakeError); !ok {
				log.Println(err)
			}
			return
		}

		log.Printf("Connecting websocket on: %s", ws.RemoteAddr())

		connectionDone := make(chan struct{})

		go func() {
			defer close(connectionDone)
			for {
				msgType, message, err := ws.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					return
				}

				if msgType == websocket.TextMessage {
					log.Println("TextMessage: ", message)
				} else if msgType == websocket.BinaryMessage {

					reader := bytes.NewReader(message)
					scanner := bufio.NewReader(reader)
					res, _ := http.ReadResponse(scanner, nil)

					if id := res.Header.Get(transport.InletsHeader); len(id) > 0 {
						bus.Send(id, res)
					}
				}
			}
		}()

		go func() {
			defer close(connectionDone)
			for {
				log.Printf("wait for request")
				outboundRequest := <-outgoing
				log.Printf("[%s] request written to websocket", outboundRequest.Header.Get(transport.InletsHeader))

				buf := new(bytes.Buffer)

				outboundRequest.Write(buf)

				ws.WriteMessage(websocket.BinaryMessage, buf.Bytes())
			}

		}()

		<-connectionDone
	}
}

func authorized(token string, r *http.Request) error {

	auth := r.Header.Get("Authorization")
	valid := false
	if len(token) == 0 {
		valid = true
	} else {
		prefix := "Bearer "
		if strings.HasPrefix(auth, prefix); len(auth) > len(prefix) && auth[len(prefix):] == token {
			valid = true
		}
	}

	if !valid {
		return fmt.Errorf("send token in header Authorization: Bearer <token>")
	}

	return nil
}
