package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/alexellis/inlets/pkg/transport"
	"github.com/gorilla/websocket"
	"github.com/twinj/uuid"
)

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

				transport.CopyHeaders(newReq.Header, &req.Header)

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
