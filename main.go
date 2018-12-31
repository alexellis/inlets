package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/twinj/uuid"

	"github.com/gorilla/websocket"
)

type Args struct {
	Port              int
	Server            bool
	Remote            string
	Upstream          string
	GatewayTimeoutRaw string
	GatewayTimeout    time.Duration
}

var client *http.Client

func buildUpstreamMap(args string) map[string]string {
	items := make(map[string]string)

	entries := strings.Split(args, ",")
	for _, entry := range entries {
		kvp := strings.Split(entry, "=")
		if len(kvp) == 1 {
			items[""] = strings.TrimSpace(kvp[0])
		} else {
			items[strings.TrimSpace(kvp[0])] = strings.TrimSpace(kvp[1])
		}
	}
	return items
}
func main() {
	args := Args{}
	flag.IntVar(&args.Port, "port", 8000, "port for server")
	flag.BoolVar(&args.Server, "server", true, "server or client")
	flag.StringVar(&args.Remote, "remote", "127.0.0.1:8000", " server address i.e. 127.0.0.1:8000")
	flag.StringVar(&args.Upstream, "upstream", "", "upstream server i.e. http://127.0.0.1:3000")
	flag.StringVar(&args.GatewayTimeoutRaw, "gateway-timeout", "5s", "timeout for upstream gateway")

	flag.Parse()

	if len(args.Upstream) == 0 {
		log.Printf("give --upstream\n")
		return
	}

	gatewayTimeout, gatewayTimeoutErr := time.ParseDuration(args.GatewayTimeoutRaw)
	if gatewayTimeoutErr != nil {
		fmt.Printf("%s\n", gatewayTimeoutErr)
		return
	}

	args.GatewayTimeout = gatewayTimeout
	log.Printf("Gateway timeout: %f secs\n", gatewayTimeout.Seconds())

	upstreamMap := buildUpstreamMap(args.Upstream)
	for key, val := range upstreamMap {
		log.Printf("Upstream: %s => %s\n", key, val)
	}

	client = http.DefaultClient
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	if args.Server {
		startServer(args)
	} else {
		runClient(args, upstreamMap)
	}
}

func runClient(args Args, upstreamMap map[string]string) {

	u := url.URL{Scheme: "ws", Host: args.Remote, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		panic(err)
	}

	fmt.Println(ws.LocalAddr())

	defer ws.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			messageType, message, err := ws.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			switch messageType {
			case websocket.TextMessage:
				log.Printf("TextMessage: %s\n", message)

				break
			case websocket.BinaryMessage:

				requestID := uuid.Formatter(uuid.NewV4(), uuid.FormatHex)

				// proxyToUpstream
				log.Printf("[%s] recv: %d", requestID, len(message))

				buf := bytes.NewBuffer(message)
				bufReader := bufio.NewReader(buf)
				req, readReqErr := http.ReadRequest(bufReader)
				if readReqErr != nil {
					log.Println(readReqErr)
					return
				}
				log.Printf("[%s] %s", requestID, req.RequestURI)

				body, _ := ioutil.ReadAll(req.Body)

				proxyHost := ""
				if val, ok := upstreamMap[req.Host]; ok {
					proxyHost = val
				} else if val, ok := upstreamMap[""]; ok {
					proxyHost = val
				}

				requestURI := fmt.Sprintf("%s%s", proxyHost, req.URL.String())
				if len(req.URL.RawQuery) > 0 {
					requestURI = requestURI + "?" + req.URL.RawQuery
				}

				log.Printf("[%s] proxy => %s", requestID, requestURI)

				newReq, newReqErr := http.NewRequest(req.Method, requestURI, bytes.NewReader(body))
				if newReqErr != nil {
					log.Printf("[%s] newReqErr: %s", requestID, newReqErr.Error())
					return
				}

				copyHeaders(newReq.Header, &req.Header)

				res, resErr := client.Do(newReq)

				if resErr != nil {
					log.Printf("[%s] Upstream tunnel err: %s", requestID, resErr.Error())

					errRes := http.Response{
						StatusCode: http.StatusBadGateway,
						Body:       ioutil.NopCloser(strings.NewReader(resErr.Error())),
					}

					buf2 := new(bytes.Buffer)
					errRes.Write(buf2)
					if errRes.Body != nil {
						defer errRes.Body.Close()
					}

					ws.WriteMessage(websocket.BinaryMessage, buf2.Bytes())

				} else {
					log.Printf("[%s] tunnel res.Status => %s", requestID, res.Status)

					buf2 := new(bytes.Buffer)

					res.Write(buf2)
					if res.Body != nil {
						defer res.Body.Close()
					}

					log.Printf("[%s] %d bytes", requestID, buf2.Len())

					ws.WriteMessage(websocket.BinaryMessage, buf2.Bytes())
				}
				break
			}

		}
	}()

	<-done
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

		copyHeaders(req.Header, &r.Header)

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

			copyHeaders(w.Header(), &res.Header)
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

func copyHeaders(destination http.Header, source *http.Header) {
	for k, v := range *source {
		vClone := make([]string, len(v))
		copy(vClone, v)
		(destination)[k] = vClone
	}
}

func startServer(args Args) {

	ch := make(chan *http.Response)
	outgoing := make(chan *http.Request)
	http.HandleFunc("/ws", serveWs(ch, outgoing))
	http.HandleFunc("/", proxyHandler(ch, outgoing, args.GatewayTimeout))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", args.Port), nil); err != nil {
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
