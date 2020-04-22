package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/rancher/remotedialer"
	"k8s.io/apimachinery/pkg/util/proxy"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Forwarder struct {
	sync.Mutex

	Listen  []Route
	session *remotedialer.Session
}

type Route struct {
	Listen        string
	TargetAddress string
}

func (f *Forwarder) dial(ctx context.Context, network, address string) (net.Conn, error) {
	var (
		s *remotedialer.Session
	)

	f.Lock()
	s = f.session
	f.Unlock()

	if s == nil {
		return nil, fmt.Errorf("no active connection")
	}

	return s.Dial(ctx, network, address)
}

func (f *Forwarder) clientProxy(target string) http.Handler {
	scheme := "http"
	if strings.HasPrefix(target, "https://") {
		scheme = "https"
		target = target[len("https://"):]
	} else if strings.HasPrefix(target, "http://") {
		target = target[len("http://"):]
	}

	host, _, err := net.SplitHostPort(target)
	if err != nil {
		host = target
	}

	transport := &http.Transport{
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
			return f.dial(ctx, network, target)
		},
		TLSClientConfig: &tls.Config{
			// TLS cert will basically never line up right
			InsecureSkipVerify: true,
		},
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := *r.URL
		u.Host = host
		u.Scheme = scheme
		r.Host = host

		httpProxy := proxy.NewUpgradeAwareHandler(&u, transport, false, false, f)
		httpProxy.ServeHTTP(w, r)
	})
}

func (f *Forwarder) Error(w http.ResponseWriter, req *http.Request, err error) {
	remotedialer.DefaultErrorWriter(w, req, http.StatusInternalServerError, err)
}

func (f *Forwarder) Start() error {
	for _, listen := range f.Listen {
		handler := f.clientProxy(listen.TargetAddress)
		go func(listen string) {
			log.Fatal(http.ListenAndServe(listen, handler))
		}(listen.Listen)
	}
	return nil
}

func (f *Forwarder) OnTunnelConnect(ctx context.Context, session *remotedialer.Session) error {
	f.Lock()
	defer f.Unlock()
	f.session = session

	go func() {
		<-ctx.Done()
		f.Lock()
		defer f.Unlock()
		f.session = nil
	}()

	return nil
}

func NewForwarder(listeners []string) (*Forwarder, error) {
	f := &Forwarder{}

	if len(listeners) == 0 {
		return f, nil
	}

	for _, listen := range listeners {
		parts := strings.SplitN(listen, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid target format %s", listen)
		}
		f.Listen = append(f.Listen, Route{
			Listen:        "127.0.0.1:" + parts[0],
			TargetAddress: parts[1],
		})
	}

	return f, nil
}
