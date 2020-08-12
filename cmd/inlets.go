// Copyright (c) Inlets Author(s) 2019. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

import (
	"fmt"
	"os"

	"github.com/morikuni/aec"
	"github.com/spf13/cobra"
)

var (
	Version   string
	GitCommit string
)

const WelcomeMessage = "Welcome to inlets.dev! Find out more at https://github.com/inlets/inlets"

func init() {
	inletsCmd.AddCommand(versionCmd)
}

// inletsCmd represents the base command when called without any sub commands.
var inletsCmd = &cobra.Command{
	Use:   "inlets",
	Short: "Expose your local endpoints to the Internet.",
	Long: `
Inlets combines a reverse proxy and websocket tunnels to expose your internal 
and development endpoints on another network, or to the public Internet via 
an exit-node.

An exit-node may be a 5-10 USD VPS or any other computer with an IPv4 IP address. 
You can also use inlets to bridge connect between private networks, only 
exposing the services on loopback or a private IP address.

It is strongly recommended to put a reverse proxy with TLS/SSL enabled such as 
Nginx or Caddy in front of your inlets server to enable an encrypted tunnel.

You can also use inlets PRO (https://github.com/inlets/inlets-pro) to enable 
automatic TLS, TCP support, and private-mode for your tunnel.

See: https://github.com/inlets/inlets for more information.`,
	Run: runInlets,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the clients version information.",
	Run:   parseBaseCommand,
}

func getVersion() string {
	if len(Version) != 0 {
		return Version
	}
	return "dev"
}

func parseBaseCommand(_ *cobra.Command, _ []string) {
	printLogo()

	fmt.Println("Version:", getVersion())
	fmt.Println("Git Commit:", GitCommit)
	os.Exit(0)
}

// Execute adds all child commands to the root command(InletsCmd) and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the InletsCmd.
func Execute(version, gitCommit string) error {

	// Get Version and GitCommit values from main.go.
	Version = version
	GitCommit = gitCommit

	if err := inletsCmd.Execute(); err != nil {
		return err
	}
	return nil
}

func runInlets(cmd *cobra.Command, args []string) {
	printLogo()
	cmd.Help()
}

func printLogo() {
	inletsLogo := aec.WhiteF.Apply(inletsFigletStr)
	fmt.Println(inletsLogo)
}

const inletsFigletStr = ` _       _      _            _
(_)_ __ | | ___| |_ ___   __| | _____   __
| | '_ \| |/ _ \ __/ __| / _` + "`" + ` |/ _ \ \ / /
| | | | | |  __/ |_\__ \| (_| |  __/\ V /
|_|_| |_|_|\___|\__|___(_)__,_|\___| \_/
`
