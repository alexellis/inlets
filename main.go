package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/alexellis/inlets/pkg/client"
	"github.com/alexellis/inlets/pkg/server"
)

// Args parsed from the command-line
type Args struct {
	Port              int
	Server            bool
	Remote            string
	Upstream          string
	GatewayTimeoutRaw string
	GatewayTimeout    time.Duration
}

func main() {
	args := Args{}
	flag.IntVar(&args.Port, "port", 8000, "port for server")
	flag.BoolVar(&args.Server, "server", true, "server or client")
	flag.StringVar(&args.Remote, "remote", "127.0.0.1:8000", " server address i.e. 127.0.0.1:8000")
	flag.StringVar(&args.Upstream, "upstream", "", "upstream server i.e. http://127.0.0.1:3000")
	flag.StringVar(&args.GatewayTimeoutRaw, "gateway-timeout", "5s", "timeout for upstream gateway")

	flag.Parse()

	argsUpstreamParser := ArgsUpstreamParser{}

	upstreamMap := map[string]string{}

	if args.Server == false {

		if len(args.Upstream) == 0 {
			log.Printf("give --upstream\n")
			return
		}
		upstreamMap = argsUpstreamParser.Parse(args.Upstream)
		for key, val := range upstreamMap {
			log.Printf("Upstream: %s => %s\n", key, val)
		}
	}

	if args.Server {
		gatewayTimeout, gatewayTimeoutErr := time.ParseDuration(args.GatewayTimeoutRaw)
		if gatewayTimeoutErr != nil {
			fmt.Printf("%s\n", gatewayTimeoutErr)
			return
		}

		args.GatewayTimeout = gatewayTimeout
		log.Printf("Gateway timeout: %f secs\n", gatewayTimeout.Seconds())
	}

	if args.Server {
		server := server.Server{
			Port:           args.Port,
			GatewayTimeout: args.GatewayTimeout,
		}
		server.Serve()
	} else {
		client := client.Client{
			Remote:      args.Remote,
			UpstreamMap: upstreamMap,
		}

		err := client.Connect()

		if err != nil {
			panic(err)
		}
	}
}
