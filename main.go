package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/alexellis/inlets/pkg/client"
	"github.com/alexellis/inlets/pkg/server"

	"github.com/urfave/cli"
)

var (
	Version = "dev"
)

func main() {
	// Moving to cobra
	app := cli.NewApp()
	app.Version = Version
	app.Name = "inlets"
	app.HelpName = "inlets"
	app.Usage = `Expose your local endpoints to the Internet

Start the tunnel server on a machine with a publicly-accessible IPv4 IP address such as a VPS
	inlets server --port80

Start the tunnel client
	inlets client --remote 192.168.0.101:80 --upstream http://127.0.0.1:3000
`

	app.Commands = []cli.Command{
		{
			Name:  "server",
			Usage: "start the tunnel server on a machine with a publicly-accessible IPv4 IP address such as a VPS",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "port",
					Value: 8000,
					Usage: "port for server",
				},
				cli.StringFlag{
					Name:  "token",
					Value: "",
					Usage: "token for authentication",
				},
				cli.StringFlag{
					Name:  "gateway-timeout",
					Value: "5s",
					Usage: "timeout for upstream gateway",
				},
				cli.BoolFlag{
					Name:  "print-token",
					Usage: "prints the token in server mode",
				},
			},
			Action: func(ctx *cli.Context) error {
				port := ctx.Int("port")
				token := ctx.String("token")
				timeout := ctx.String("gateway-timeout")
				printToken := ctx.Bool("print-token")

				if len(token) > 0 && printToken {
					log.Printf("Server token: %s", token)
				}

				gatewayTimeout, gatewayTimeoutErr := time.ParseDuration(timeout)
				if gatewayTimeoutErr != nil {
					fmt.Printf("%s\n", gatewayTimeoutErr)
					return nil
				}

				log.Printf("Gateway timeout: %f secs\n", gatewayTimeout.Seconds())

				s := server.Server{
					Port:           port,
					GatewayTimeout: gatewayTimeout,
					Token:          token,
				}
				s.Serve()
				return nil
			},
		},

		{
			Name:  "client",
			Usage: "start the tunnel client",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "remote",
					Value: "127.0.0.1:8000",
					Usage: "server address i.e. 127.0.0.1:8000",
				},
				cli.StringFlag{
					Name:  "upstream",
					Value: "",
					Usage: "upstream server i.e. http://127.0.0.1:3000, http://127.0.0.1:3001",
				},
				cli.StringFlag{
					Name:  "token",
					Value: "",
					Usage: "token for authentication",
				},
			},
			Action: func(ctx *cli.Context) error {
				remote := ctx.String("remote")
				upstream := ctx.String("upstream")
				token := ctx.String("token")

				if len(upstream) == 0 {
					log.Printf("give --upstream\n")
					return nil
				}

				argsUpstreamParser := ArgsUpstreamParser{}
				upstreamMap := argsUpstreamParser.Parse(upstream)

				for key, val := range upstreamMap {
					log.Printf("Upstream: %s => %s\n", key, val)
				}

				c := client.Client{
					Remote:      remote,
					UpstreamMap: upstreamMap,
					Token:       token,
				}

				return c.Connect()
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
