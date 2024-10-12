package cmd

import (
	_ "embed"
	"time"

	"github.com/CnTeng/rx-todo/cli"
	"github.com/CnTeng/rx-todo/model"
	"github.com/CnTeng/rx-todo/rpc"
	"github.com/spf13/cobra"
)

//go:embed example/task_create.toml
var taskCreateExample string

var taskListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List all tasks",
	Run: func(cmd *cobra.Command, args []string) {
		c := rpc.NewClient(network, socketPath, 5*time.Second)

		tasks := cli.TaskSlice{}
		if err := c.Call("Task.List", nil, &tasks); err != nil {
			cobra.CheckErr(err)
		}

		cli.NewCLI(cli.Nerd).ListTasks(tasks.SortByID(), nil)
	},
}

var taskAddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a"},
	Short:   "Add new task",
	Run: func(cmd *cobra.Command, args []string) {
		request := &model.TaskCreationRequest{
			Content:     getValue(cmd, cmd.Flags().GetString, "content"),
			Description: getValue(cmd, cmd.Flags().GetString, "description"),
			Priority:    getValue(cmd, cmd.Flags().GetInt, "priority"),
			ProjectID:   getValue(cmd, cmd.Flags().GetInt64, "project"),
			ParentID:    getValue(cmd, cmd.Flags().GetInt64, "parent"),
			Labels:      getValue(cmd, cmd.Flags().GetStringSlice, "labels"),
		}

		if due := getValue(cmd, cmd.Flags().GetString, "due"); due != nil {
			date, err := model.ParseDueDate(*due)
			if err != nil {
				cobra.CheckErr(err)
			}

			if request.Due == nil {
				request.Due = &model.Due{Date: date}
			} else {
				request.Due.Date = date
			}
		}

		if recurring := getValue(cmd, cmd.Flags().GetBool, "recurring"); recurring != nil {
			if request.Due == nil {
				request.Due = &model.Due{Recurring: recurring}
			} else {
				request.Due.Recurring = recurring
			}
		}

		if cmd.Flags().Changed("edit") {
			edit, err := cli.NewEditFile(request, taskCreateExample)
			if err != nil {
				cobra.CheckErr(err)
			}

			if err := edit.ParseContent(request); err != nil {
				cobra.CheckErr(err)
			}
		}

		c := rpc.NewClient(network, socketPath, 5*time.Second)

		task := &model.Task{}
		if err := c.Call("Task.Create", request, task); err != nil {
			cobra.CheckErr(err)
		}
		sm := cli.NewStatusMap([]*model.Task{task}, cli.Add)

		tasks := cli.TaskSlice{}
		if err := c.Call("Task.List", nil, &tasks); err != nil {
			cobra.CheckErr(err)
		}

		cli.NewCLI(cli.Nerd).ListTasks(tasks.SortByID(), sm)
	},
}

var taskUpdateCmd = &cobra.Command{
	Use:     "modify",
	Aliases: []string{"m"},
	Short:   "modify task",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var taskDeleteCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm"},
	Short:   "remove task",
	Run: func(cmd *cobra.Command, args []string) {
		request := &model.TaskDeleteRequestWithID{}

		id, err := cmd.Flags().GetInt64("id")
		if err != nil {
			cobra.CheckErr(err)
		}
		request.ID = id

		c := rpc.NewClient("unix", "/tmp/rx-todo.sock", 5*time.Second)

		tasks := cli.TaskSlice{}
		if err := c.Call("Task.List", nil, &tasks); err != nil {
			cobra.CheckErr(err)
		}

		task := &model.Task{}
		if err := c.Call("Task.Delete", request, task); err != nil {
			cobra.CheckErr(err)
		}
		sm := cli.NewStatusMap([]*model.Task{task}, cli.Delete)

		cli.NewCLI(cli.Nerd).ListTasks(tasks.SortByID(), sm)
	},
}

func init() {
	rootCmd.AddCommand(taskListCmd)
	rootCmd.AddCommand(taskAddCmd)
	rootCmd.AddCommand(taskDeleteCmd)

	taskAddCmd.Flags().StringP("content", "c", "", "Task content")
	taskAddCmd.Flags().StringP("description", "D", "", "Task description")
	taskAddCmd.Flags().StringP("due", "d", "", "Task labels")
	taskAddCmd.Flags().BoolP("recurring", "r", false, "Task recurring")
	taskAddCmd.Flags().String("duration", "", "Task duration") // TODO: implement duration
	taskAddCmd.Flags().IntP("priority", "p", 0, "Task priority")
	taskAddCmd.Flags().Int64("project", 0, "Task project ID")
	taskAddCmd.Flags().Int64("parent", 0, "Task parent ID")
	taskAddCmd.Flags().StringSliceP("labels", "l", []string{}, "Task labels")
	taskAddCmd.Flags().BoolP("edit", "e", false, "Eidt task")

	taskDeleteCmd.Flags().Int64P("id", "i", 0, "Task ID")
}
