package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/alexellis/inlets/pkg/router"
	"github.com/alexellis/inlets/pkg/transport"
	"github.com/rancher/remotedialer"
	"github.com/twinj/uuid"
	"k8s.io/apimachinery/pkg/util/proxy"
)

// Server for the exit-node of inlets
type Server struct {
	Port   int
	Token  string
	router router.Router
	server *remotedialer.Server

	DisableWrapTransport bool
}

// Serve traffic
func (s *Server) Serve() {
	s.server = remotedialer.New(s.authorized, remotedialer.DefaultErrorWriter)
	s.router.Server = s.server

	http.HandleFunc("/", s.proxy)
	http.HandleFunc("/tunnel", s.tunnel)

	log.Printf("Listening on :%d\n", s.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.Port), nil); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) tunnel(w http.ResponseWriter, r *http.Request) {
	s.server.ServeHTTP(w, r)
	s.router.Remove(r)
}

func (s *Server) proxy(w http.ResponseWriter, r *http.Request) {
	route := s.router.Lookup(r)
	if route == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	inletsID := uuid.Formatter(uuid.NewV4(), uuid.FormatHex)
	log.Printf("[%s] proxy %s %s %s", inletsID, r.Host, r.Method, r.URL.String())
	r.Header.Set(transport.InletsHeader, inletsID)

	u := *r.URL
	u.Host = r.Host
	u.Scheme = route.Scheme

	httpProxy := proxy.NewUpgradeAwareHandler(&u, route.Transport, !s.DisableWrapTransport, false, s)
	httpProxy.ServeHTTP(w, r)
}

func (s Server) Error(w http.ResponseWriter, req *http.Request, err error) {
	remotedialer.DefaultErrorWriter(w, req, http.StatusInternalServerError, err)
}

func (s *Server) dialerFor(id, host string) remotedialer.Dialer {
	return func(network, address string) (net.Conn, error) {
		return s.server.Dial(id, time.Minute, network, host)
	}
}

func (s *Server) tokenValid(req *http.Request) bool {
	auth := req.Header.Get("Authorization")
	return len(s.Token) == 0 || auth == "Bearer "+s.Token
}

func (s *Server) authorized(req *http.Request) (id string, ok bool, err error) {
	defer func() {
		if id == "" {
			// empty id is also an auth failure
			ok = false
		}
		if !ok || err != nil {
			// don't let non-authed request clear routes
			req.Header.Del(transport.InletsHeader)
		}
	}()

	if !s.tokenValid(req) {
		return "", false, nil
	}

	return s.router.Add(req), true, nil
}
