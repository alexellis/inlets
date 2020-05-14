// Copyright (c) Inlets Author(s) 2019. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/inlets/inlets/pkg/server"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// serverCmd represents the server sub command.
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: `Start the tunnel server.`,
	Long: `Start the tunnel server on a machine with a publicly-accessible IPv4 IP 
address such as a VPS.

Example: inlets server -p 80 
Example: inlets server --port 80 --control-port 8080

Note: You can pass the --token argument followed by a token value to both the 
server and client to prevent unauthorized connections to the tunnel.`,
	RunE: runServer,
}

func init() {

	serverCmd.Flags().IntP("port", "p", 8000, "port for server and for tunnel")
	serverCmd.Flags().StringP("token", "t", "", "token for authentication")
	serverCmd.Flags().Bool("print-token", true, "prints the token in server mode")
	serverCmd.Flags().StringP("token-from", "f", "", "read the authentication token from a file")
	serverCmd.Flags().Bool("disable-transport-wrapping", false, "disable wrapping the transport that removes CORS headers for example")
	serverCmd.Flags().IntP("control-port", "c", 8080, "control port for tunnel")
	serverCmd.Flags().StringP("bind-addr", "b", "", "address the server should be bound to")

	inletsCmd.AddCommand(serverCmd)
}

// runServer does the actual work of reading the arguments passed to the server sub command.
func runServer(cmd *cobra.Command, _ []string) error {

	log.Printf("%s", WelcomeMessage)
	log.Printf("Starting server - version %s", getVersion())

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

		// new-lines will be stripped, this is not configurable and is to
		// make the code foolproof for beginners
		token = strings.TrimRight(string(fileData), "\n")
	} else {
		tokenVal, err := cmd.Flags().GetString("token")
		if err != nil {
			return errors.Wrap(err, "failed to get 'token' value.")
		}
		token = tokenVal
	}

	if tokenEnv, ok := os.LookupEnv("TOKEN"); ok && len(tokenEnv) > 0 {
		fmt.Printf("Token read from environment variable.\n")
		token = tokenEnv
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

	controlPort := port
	if cmd.Flags().Changed("control-port") {
		val, err := cmd.Flags().GetInt("control-port")
		if err != nil {
			return errors.Wrap(err, "failed to get the 'control-port' value.")
		}
		controlPort = val
	}

	if portVal, exists := os.LookupEnv("PORT"); exists && len(portVal) > 0 {
		port, _ = strconv.Atoi(portVal)
		controlPort = port
	}

	disableWrapTransport, err := cmd.Flags().GetBool("disable-transport-wrapping")
	if err != nil {
		return errors.Wrap(err, "failed to get the 'disable-transport-wrapping' value.")
	}

	bindAddr, err := cmd.Flags().GetString("bind-addr")
	if err != nil {
		return errors.Wrap(err, "failed to get the 'bind-addr' value.")
	}

	inletsServer := server.Server{
		Port:        port,
		ControlPort: controlPort,
		BindAddr:    bindAddr,
		Token:       token,

		DisableWrapTransport: disableWrapTransport,
	}

	inletsServer.Serve()
	return nil
}
