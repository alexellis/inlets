// Copyright (c) Inlets Author(s) 2019. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package router

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/inlets/inlets/pkg/transport"
	"github.com/rancher/remotedialer"
)

type target struct {
	id     string
	domain string
	target string
}

type transportKey struct {
	id   string
	host string
}

type transportValue struct {
	tranport *http.Transport
	scheme   string
}

type Router struct {
	sync.RWMutex
	transportLock sync.RWMutex

	Server     *remotedialer.Server
	domains    map[string][]target
	clients    map[string][]target
	transports map[transportKey]transportValue
}

type Route struct {
	ID        string
	Scheme    string
	Transport *http.Transport
}

func (r *Router) Lookup(req *http.Request) *Route {
	r.RLock()
	defer r.RUnlock()

	targets := r.domains[req.Host]
	if len(targets) == 0 {
		targets = r.domains[""]
	}
	if len(targets) == 0 {
		return nil
	}

	id, host := targets[0].id, targets[0].target
	scheme, transport := r.getTransport(id, host)
	return &Route{
		ID:        id,
		Scheme:    scheme,
		Transport: transport,
	}
}

func (r *Router) getTransport(id, host string) (string, *http.Transport) {
	key := transportKey{id: id, host: host}

	r.transportLock.RLock()
	val, ok := r.transports[key]
	if ok {
		r.transportLock.RUnlock()
		return val.scheme, val.tranport
	}
	r.transportLock.RUnlock()

	r.transportLock.Lock()
	defer r.transportLock.Unlock()

	targetHost := host
	scheme := "http"
	if strings.HasPrefix(host, "https://") {
		scheme = "https"
		targetHost = host[len("https://"):]
	} else if strings.HasPrefix(host, "http://") {
		targetHost = host[len("http://"):]
	}

	transport := &http.Transport{
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
			return r.Server.Dialer(id)(ctx, network, targetHost)
		},
		TLSClientConfig: &tls.Config{
			// TLS cert will basically never line up right
			InsecureSkipVerify: true,
		},
	}

	if r.transports == nil {
		r.transports = map[transportKey]transportValue{}
	}

	r.transports[transportKey{id: id, host: host}] = transportValue{scheme: scheme, tranport: transport}
	return scheme, transport
}

func (r *Router) Add(req *http.Request) string {
	var targets []target

	id := req.Header.Get(transport.InletsHeader)
	upstreams := req.Header[http.CanonicalHeaderKey(transport.UpstreamHeader)]

	for _, upstream := range upstreams {
		parts := strings.SplitN(upstream, "=", 2)
		if len(parts) != 2 {
			continue
		}
		targets = append(targets, target{
			id:     id,
			domain: parts[0],
			target: parts[1],
		})
	}

	if id == "" || len(targets) == 0 {
		return ""
	}

	r.Lock()
	defer r.Unlock()

	if r.domains == nil {
		r.domains = map[string][]target{}
		r.clients = map[string][]target{}
	}

	for _, target := range targets {
		r.domains[target.domain] = append(r.domains[target.domain], target)
	}

	r.clients[id] = targets
	return id
}

func (r *Router) Remove(req *http.Request) {
	r.Lock()
	defer r.Unlock()

	id := req.Header.Get(transport.InletsHeader)
	targets := r.clients[id]
	delete(r.clients, id)

	for _, idTarget := range targets {
		var newTargets []target
		domainTargets := r.domains[idTarget.domain]

		for _, domainTarget := range domainTargets {
			if domainTarget.id != id {
				newTargets = append(newTargets, domainTarget)
			}
		}

		if len(newTargets) == 0 {
			delete(r.domains, idTarget.domain)
		} else {
			r.domains[idTarget.domain] = newTargets
		}
	}
}
