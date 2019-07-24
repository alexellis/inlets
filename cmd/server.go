package cmd

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/alexellis/inlets/pkg/server"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	inletsCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntP("port", "p", 8000, "port for server")
	serverCmd.Flags().StringP("token", "t", "", "token for authentication")
	serverCmd.Flags().Bool("print-token", true, "prints the token in server mode")
	serverCmd.Flags().StringP("token-from", "f", "", "read the authentication token from a file")
	serverCmd.Flags().Bool("disable-transport-wrapping", false, "disable wrapping the transport that removes CORS headers for example")
}

// serverCmd represents the server sub command.
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the tunnel server on a machine with a publicly-accessible IPv4 IP address such as a VPS.",
	Long: `Start the tunnel server on a machine with a publicly-accessible IPv4 IP address such as a VPS.

Example: inlets server -p 80 
Note: You can pass the --token argument followed by a token value to both the server and client to prevent unauthorized connections to the tunnel.`,
	RunE: runServer,
}

// runServer does the actual work of reading the arguments passed to the server sub command.
func runServer(cmd *cobra.Command, _ []string) error {

	tokenFile, err := cmd.Flags().GetString("token-from")
	if err != nil {
		return errors.Wrap(err, "failed to get 'token-from' value.")
	}

	var token string
	if len(tokenFile) > 0 {
		fileData, err := ioutil.ReadFile(tokenFile)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("unable to load file: %s", tokenFile))
		}
		token = string(fileData)
	} else {
		tokenVal, err := cmd.Flags().GetString("token")
		if err != nil {
			return errors.Wrap(err, "failed to get 'token' value.")
		}
		token = tokenVal
	}

	printToken, err := cmd.Flags().GetBool("print-token")
	if err != nil {
		return errors.Wrap(err, "failed to get 'print-token' value.")
	}

	if len(token) > 0 && printToken {
		log.Printf("Server token: %q", token)
	}

	port, err := cmd.Flags().GetInt("port")
	if err != nil {
		return errors.Wrap(err, "failed to get the 'port' value.")
	}

	disableWrapTransport, err := cmd.Flags().GetBool("disable-transport-wrapping")
	if err != nil {
		return errors.Wrap(err, "failed to get the 'disable-transport-wrapping' value.")
	}

	inletsServer := server.Server{
		Port:  port,
		Token: token,

		DisableWrapTransport: disableWrapTransport,
	}

	inletsServer.Serve()
	return nil
}
