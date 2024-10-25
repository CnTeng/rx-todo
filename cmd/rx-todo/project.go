package cmd

import (
	"fmt"
	"time"

	"github.com/CnTeng/rx-todo/client"
	"github.com/CnTeng/rx-todo/internal/cli"
	"github.com/CnTeng/rx-todo/internal/rpc"
	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "list all labels",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var projectListCmd = &cobra.Command{
	Use:     "list [name] [flags]...",
	Aliases: []string{"l"},
	Short:   "List all projects",
	Run: func(cmd *cobra.Command, args []string) {
		c := rpc.NewClient(network, socketPath, 5*time.Second)

		projects := client.ProjectSlice{}
		if err := c.Call("project.list", nil, &projects); err != nil {
			cobra.CheckErr(err)
		}

		// Filter by name
		var name *string
		if len(args) == 1 {
			name = &args[0]
		} else {
			name = getValue(cmd, cmd.Flags().GetString, "name")
		}
		if name != nil {
			projects = projects.FilterByName(*name)
		}
		if len(projects) == 0 {
			cobra.CheckErr(fmt.Errorf("no project found"))
		}

		cli.NewCLI(cli.Nerd).PrintProjects(projects, "projects", client.None)
	},
}

func init() {
	rootCmd.AddCommand(projectCmd)
	projectCmd.AddCommand(projectListCmd)

	projectListCmd.Flags().StringP("name", "n", "", "filter projects by name")
}
