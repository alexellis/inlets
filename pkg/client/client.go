package client

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
)

var httpClient *http.Client

// Client for inlets
type Client struct {
	// Remote site for websocket address
	Remote string

	// Map of upstream servers dns.entry=http://ip:port
	UpstreamMap map[string]string

	// Token for authentication
	Token string
}

// Connect connect and serve traffic through websocket
func (c *Client) Connect() error {

	httpClient = http.DefaultClient
	httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	u := url.URL{Scheme: "ws", Host: c.Remote, Path: "/tunnel"}
	log.Printf("connecting to %s", u.String())

	ws, _, err := websocket.DefaultDialer.Dial(u.String(), http.Header{
		"Authorization": []string{"Bearer " + c.Token},
	})

	if err != nil {
		return err
	}

	log.Printf("Connected to websocket: %s", ws.LocalAddr())

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

			case websocket.BinaryMessage:
				// proxyToUpstream

				buf := bytes.NewBuffer(message)
				bufReader := bufio.NewReader(buf)
				req, readReqErr := http.ReadRequest(bufReader)
				if readReqErr != nil {
					log.Println(readReqErr)
					return
				}

				inletsID := req.Header.Get(transport.InletsHeader)
				// log.Printf("[%s] recv: %d", requestID, len(message))

				log.Printf("[%s] %s", inletsID, req.RequestURI)

				body, _ := ioutil.ReadAll(req.Body)

				proxyHost := ""
				if val, ok := c.UpstreamMap[req.Host]; ok {
					proxyHost = val
				} else if val, ok := c.UpstreamMap[""]; ok {
					proxyHost = val
				}

				requestURI := fmt.Sprintf("%s%s", proxyHost, req.URL.String())
				if len(req.URL.RawQuery) > 0 {
					requestURI = requestURI + "?" + req.URL.RawQuery
				}

				log.Printf("[%s] proxy => %s", inletsID, requestURI)

				newReq, newReqErr := http.NewRequest(req.Method, requestURI, bytes.NewReader(body))
				if newReqErr != nil {
					log.Printf("[%s] newReqErr: %s", inletsID, newReqErr.Error())
					return
				}

				transport.CopyHeaders(newReq.Header, &req.Header)

				res, resErr := httpClient.Do(newReq)

				if resErr != nil {
					log.Printf("[%s] Upstream tunnel err: %s", inletsID, resErr.Error())

					errRes := http.Response{
						StatusCode: http.StatusBadGateway,
						Body:       ioutil.NopCloser(strings.NewReader(resErr.Error())),
						Header:     http.Header{},
					}

					errRes.Header.Set(transport.InletsHeader, inletsID)
					buf2 := new(bytes.Buffer)
					errRes.Write(buf2)
					if errRes.Body != nil {
						errRes.Body.Close()
					}

					ws.WriteMessage(websocket.BinaryMessage, buf2.Bytes())

				} else {
					log.Printf("[%s] tunnel res.Status => %s", inletsID, res.Status)

					buf2 := new(bytes.Buffer)
					res.Header.Set(transport.InletsHeader, inletsID)

					res.Write(buf2)
					if res.Body != nil {
						res.Body.Close()
					}

					log.Printf("[%s] %d bytes", inletsID, buf2.Len())

					ws.WriteMessage(websocket.BinaryMessage, buf2.Bytes())
				}
			}

		}
	}()

	<-done

	return nil
}
