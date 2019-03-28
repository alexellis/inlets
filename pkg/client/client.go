package client

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/alexellis/inlets/pkg/transport"
	"github.com/rancher/remotedialer"
	"github.com/twinj/uuid"
)

// Client for inlets
type Client struct {
	// Remote site for websocket address
	Remote string

	// Map of upstream servers dns.entry=http://ip:port
	UpstreamMap map[string]string

	// Token for authentication
	Token string
}

func allowsAllow(network, address string) bool {
	return true
}

// Connect connect and serve traffic through websocket
func (c *Client) Connect() error {
	headers := http.Header{}
	headers.Set(transport.InletsHeader, uuid.Formatter(uuid.NewV4(), uuid.FormatHex))
	for k, v := range c.UpstreamMap {
		headers.Set(transport.UpstreamHeader, fmt.Sprintf("%s=%s", k, v))
	}
	if c.Token != "" {
		headers.Add("Authorization", "Bearer "+c.Token)
	}

	fmt.Println(headers)
	url := c.Remote
	if !strings.HasPrefix(url, "ws") {
		url = "ws://" + url
	}
	for {
		remotedialer.ClientConnect(url+"/tunnel", headers, nil, allowsAllow, nil)
	}
}
