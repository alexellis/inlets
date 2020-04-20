package client

import (
	"context"
	"fmt"
	"inet.af/tcpproxy"
	"github.com/rancher/remotedialer"
	"net"
	"strings"
	"sync"
)

type Forwarder struct {
	sync.Mutex

	Listen []string
	session *remotedialer.Session
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

func (f *Forwarder) Start() error {
	var p tcpproxy.Proxy
	for _, listen := range f.Listen {
		parts := strings.SplitN(listen, ":", 2)
		port, target := parts[0], parts[1]
		d := &tcpproxy.DialProxy{
			Addr:                 target,
			DialContext:          f.dial,
		}
		p.AddRoute(":" + port, d)
	}
	return p.Run()
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

func NewForwarder(listeners []string) func(ctx context.Context, session *remotedialer.Session) error {
	if len(listeners) == 0 {
		return nil
	}

	f := &Forwarder{
		Listen: listeners,
	}
	return f.OnTunnelConnect
}
