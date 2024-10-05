package cmd

import (
	"time"

	"github.com/CnTeng/rx-todo/cli"
	"github.com/CnTeng/rx-todo/rpc"
	"github.com/spf13/cobra"
)

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "list all tasks",
	Run: func(cmd *cobra.Command, args []string) {
		c := rpc.NewClient("unix", "/tmp/rx-todo.sock", 5*time.Second)

		Tasks := cli.TaskSlice{}
		if err := c.Call("Task.List", nil, &Tasks); err != nil {
			cobra.CheckErr(err)
		}

		cmd.Print(Tasks.List(nil))
	},
}

func init() {
	rootCmd.AddCommand(taskCmd)
}
