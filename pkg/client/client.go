// Copyright (c) Inlets Author(s) 2019. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package client

import (
	"context"
	"fmt"
	"github.com/inlets/inlets/pkg/transport"
	"github.com/rancher/remotedialer"
	"github.com/twinj/uuid"
	"net/http"
	"strings"
)

// Client for inlets
type Client struct {
	// Remote site for websocket address
	Remote string

	// Map of upstream servers dns.entry=http://ip:port
	UpstreamMap map[string]string

	// Token for authentication
	Token string

	// Local listener to forward connections
	Listeners []string
}

func AllowsAllow(network, address string) bool {
	return true
}

// Connect connect and serve traffic through websocket
func (c *Client) Connect() error {
	headers := http.Header{}
	headers.Set(transport.InletsHeader, uuid.Formatter(uuid.NewV4(), uuid.FormatHex))
	for k, v := range c.UpstreamMap {
		headers.Add(transport.UpstreamHeader, fmt.Sprintf("%s=%s", k, v))
	}
	if c.Token != "" {
		headers.Add("Authorization", "Bearer "+c.Token)
	}

	url := c.Remote
	if !strings.HasPrefix(url, "ws") {
		url = "ws://" + url
	}

	onConnect := NewForwarder(c.Listeners)

	for {
		remotedialer.ClientConnect(context.Background(), url+"/tunnel", headers, nil, AllowsAllow, onConnect)
	}
}
