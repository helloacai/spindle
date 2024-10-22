package main

import (
	"github.com/spf13/cobra"

	"github.com/helloacai/spindle/pkg/keepalive"
	"github.com/helloacai/spindle/pkg/server"
	"github.com/helloacai/spindle/pkg/substreams"
)

func init() {
	substreams.AddFlags(rootCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "spindle",
	Short: "ACS thread manager",
	Long:  `ACS thread manager`,
	Run:   root,
}

func root(cmd *cobra.Command, args []string) {
	server.Start()
	keepalive.Start()

	substreams.Listen(cmd)
}
