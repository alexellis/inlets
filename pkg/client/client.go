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
	"time"

	"github.com/alexellis/inlets/pkg/transport"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
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

	// PingWaitDuration duration to wait between pings
	PingWaitDuration time.Duration
}

// Connect connect and serve traffic through websocket
func (c *Client) Connect() error {

	httpClient = http.DefaultClient
	httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	remote := c.Remote
	if !strings.HasPrefix(remote, "ws") {
		remote = "ws://" + remote
	}

	remoteURL, urlErr := url.Parse(remote)
	if urlErr != nil {
		return errors.Wrap(urlErr, "bad remote URL")
	}

	u := url.URL{Scheme: remoteURL.Scheme, Host: remoteURL.Host, Path: "/tunnel"}

	log.Printf("connecting to %s with ping=%s", u.String(), c.PingWaitDuration.String())

	wsc, _, err := websocket.DefaultDialer.Dial(u.String(), http.Header{
		"Authorization": []string{"Bearer " + c.Token},
	})

	ws := transport.NewWebsocketConn(wsc, c.PingWaitDuration)

	if err != nil {
		return err
	}

	log.Printf("Connected to websocket: %s", ws.LocalAddr())

	defer wsc.Close()

	// Send pings
	tickerDone := make(chan bool)

	go func() {
		log.Printf("Writing pings")

		ticker := time.NewTicker((c.PingWaitDuration * 9) / 10) // send on a period which is around 9/10ths of original value
		for {
			select {
			case <-ticker.C:
				if err := ws.Ping(); err != nil {
					close(tickerDone)
				}
				break
			case <-tickerDone:
				log.Printf("tickerDone, no more pings will be sent from client\n")
				return
			}
		}
	}()

	// Work with websocket
	done := make(chan struct{})

	go func() {
		defer close(done)

		for {
			messageType, message, err := ws.ReadMessage()
			fmt.Printf("Read a message from websocket.\n")
			if err != nil {
				fmt.Printf("Read error: %s.\n", err)
				return
			}

			switch messageType {
			case websocket.TextMessage:
				log.Printf("TextMessage: %s\n", message)

				break
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
