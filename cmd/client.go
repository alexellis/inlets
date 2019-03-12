package cmd

import (
	"log"
	"strings"
	"time"

	"github.com/alexellis/inlets/pkg/client"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	inletsCmd.AddCommand(clientCmd)
	clientCmd.Flags().StringP("remote", "r", "127.0.0.1:8000", "server address i.e. 127.0.0.1:8000")
	clientCmd.Flags().StringP("upstream", "u", "", "upstream server i.e. http://127.0.0.1:3000")
	clientCmd.Flags().StringP("token", "t", "", "token for authentication")
	clientCmd.Flags().DurationP("ping", "p", time.Second*10, "ping internal")
}

type UpstreamParser interface {
	Parse(input string) map[string]string
}

type ArgsUpstreamParser struct {
}

func (a *ArgsUpstreamParser) Parse(input string) map[string]string {
	upstreamMap := buildUpstreamMap(input)

	return upstreamMap
}

func buildUpstreamMap(args string) map[string]string {
	items := make(map[string]string)

	entries := strings.Split(args, ",")
	for _, entry := range entries {
		kvp := strings.Split(entry, "=")
		if len(kvp) == 1 {
			items[""] = strings.TrimSpace(kvp[0])
		} else {
			items[strings.TrimSpace(kvp[0])] = strings.TrimSpace(kvp[1])
		}
	}
	return items
}

// clientCmd represents the client sub command.
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Start the tunnel client.",
	Long: `Start the tunnel client.

Example: inlets client --remote=192.168.0.101:80 --upstream=http://127.0.0.1:3000 
Note: You can pass the --token argument followed by a token value to both the server and client to prevent unauthorized connections to the tunnel.`,
	RunE: runClient,
}

// runClient does the actual work of reading the arguments passed to the client sub command.
func runClient(cmd *cobra.Command, _ []string) error {
	upstream, err := cmd.Flags().GetString("upstream")
	if err != nil {
		return errors.Wrap(err, "failed to get 'upstream' value")
	}

	if len(upstream) == 0 {
		return errors.New("upstream is missing in the client argument")
	}

	argsUpstreamParser := ArgsUpstreamParser{}
	upstreamMap := argsUpstreamParser.Parse(upstream)
	for k, v := range upstreamMap {
		log.Printf("Upstream: %s => %s\n", k, v)
	}

	remote, err := cmd.Flags().GetString("remote")
	if err != nil {
		return errors.Wrap(err, "failed to get 'remote' value.")
	}

	token, err := cmd.Flags().GetString("token")
	if err != nil {
		return errors.Wrap(err, "failed to get 'token' value.")
	}

	pingDuration, err := cmd.Flags().GetDuration("ping")
	if err != nil {
		return errors.Wrap(err, "failed to get 'ping' value.")
	}

	inletsClient := client.Client{
		Remote:           remote,
		UpstreamMap:      upstreamMap,
		Token:            token,
		PingWaitDuration: pingDuration,
	}

	if err := inletsClient.Connect(); err != nil {
		return err
	}

	return nil
}
