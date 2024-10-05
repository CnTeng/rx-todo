package cmd

import (
	"github.com/CnTeng/rx-todo/client"
	"github.com/CnTeng/rx-todo/daemon"
	"github.com/spf13/cobra"
)

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "start daemon",
	Run: func(cmd *cobra.Command, args []string) {
		socketPath := "/tmp/rx-todo.sock"

		c, err := client.NewClient(config.Address, config.Token, "rx-todo/resources.json")
		if err != nil {
			cobra.CheckErr(err)
		}
		daemon := daemon.NewDaemon(socketPath, c)

		if err := c.Sync(); err != nil {
			cobra.CheckErr(err)
		}

		if err := daemon.Serve(); err != nil {
			cobra.CheckErr(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(daemonCmd)
}
