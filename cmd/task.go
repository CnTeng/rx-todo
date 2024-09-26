package cmd

import (
	"github.com/CnTeng/rx-todo/cli"
	"github.com/CnTeng/rx-todo/client"
	"github.com/spf13/cobra"
)

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "list all tasks",
	Run: func(cmd *cobra.Command, args []string) {
		s, err := client.NewStorage("rx-todo/resources.json")
		if err != nil {
			panic(err)
		}

		c := client.NewClient(config.Address, config.Token, s.Path)

		if res, err := c.Sync(); err != nil {
			panic(err)
		} else {
			tasks := cli.TaskSlice(*res.Tasks)
			cmd.Print(tasks.List())
			c.Store(res)
		}
	},
}

func init() {
	rootCmd.AddCommand(taskCmd)
}
