package cmd

import (
	_ "embed"
	"fmt"
	"time"

	"github.com/CnTeng/rx-todo/client"
	"github.com/CnTeng/rx-todo/internal/cli"
	"github.com/CnTeng/rx-todo/internal/model"
	"github.com/CnTeng/rx-todo/internal/rpc"
	"github.com/spf13/cobra"
)

//go:embed example/task_create.toml
var taskCreateExample string

var taskListCmd = &cobra.Command{
	Use:     "list [name] [flags]...",
	Aliases: []string{"l"},
	Short:   "List all tasks",
	Run: func(cmd *cobra.Command, args []string) {
		c := rpc.NewClient(network, socketPath, 5*time.Second)

		tasks := client.TaskSlice{}
		if err := c.Call("task.list", nil, &tasks); err != nil {
			cobra.CheckErr(err)
		}

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
			tasks = tasks.FilterByName(*name)
		}
		if len(tasks) == 0 {
			cobra.CheckErr(fmt.Errorf("no task found"))
		}

		for _, project := range projects {
			t := tasks.FilterByProjectID(project.ID)
			if len(t) == 0 {
				continue
			}

			cli.NewCLI(cli.Nerd).PrintTasks(t.SortByID(), project.Name, client.None)
		}
	},
}

var taskAddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a"},
	Short:   "Add new task",
	Run: func(cmd *cobra.Command, args []string) {
		cl := cli.NewCLI(cli.Nerd)

		var priority model.Priority
		p := getValue(cmd, cmd.Flags().GetInt, "priority")

		if p != nil {
			if *p < 0 || *p > 3 {
				cobra.CheckErr(fmt.Errorf("priority must be between 0 and 3"))
			} else {
				priority = model.Priority(*p)
			}
		} else {
			priority = model.Priority(0)
		}

		request := &model.TaskCreationRequest{
			Name:        getValue(cmd, cmd.Flags().GetString, "name"),
			Description: getValue(cmd, cmd.Flags().GetString, "description"),
			Priority:    &priority,
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
			if err := cl.StartInteractiveMode(request, taskCreateExample); err != nil {
				cobra.CheckErr(err)
			}
		}

		c := rpc.NewClient(network, socketPath, 5*time.Second)

		task := &model.Task{}
		if err := c.Call("task.create", request, task); err != nil {
			cobra.CheckErr(err)
		}

		cl.PrintTasks(client.TaskSlice{task}, "task", client.Add)
	},
}

var taskUpdateCmd = &cobra.Command{
	Use:     "modify",
	Aliases: []string{"m"},
	Short:   "modify task",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var taskMoveCmd = &cobra.Command{
	Use:     "move",
	Aliases: []string{"mv"},
	Short:   "move task",
	Run: func(cmd *cobra.Command, args []string) {
		c := rpc.NewClient(network, socketPath, 5*time.Second)
		t := cli.NewCLI(cli.Nerd)

		request := &client.TaskMoveRequestWithID{}

		parent := getValue(cmd, cmd.Flags().GetString, "parent")

		tasks := client.TaskSlice{}
		if err := c.Call("task.list", nil, &tasks); err != nil {
			cobra.CheckErr(err)
		}

		projects := client.ProjectSlice{}
		if err := c.Call("project.list", nil, &projects); err != nil {
			cobra.CheckErr(err)
		}

		var name *string
		if len(args) == 1 {
			name = &args[0]
		} else {
			name = getValue(cmd, cmd.Flags().GetString, "name")
		}

		if id, err := t.SelectOne(tasks.FilterByName(*name)); err != nil {
			cobra.CheckErr(err)
		} else {
			request.ID = id
		}

		if parent != nil {
			if id, error := t.SelectOne(projects.FilterByName(*parent)); error != nil {
				cobra.CheckErr(error)
			} else {
				request.ProjectID = &id
			}

			if request.ProjectID == nil {
				if id, error := t.SelectOne(tasks.FilterByName(*parent)); error != nil {
					cobra.CheckErr(error)
				} else {
					request.ParentID = &id
				}
			}
		}

		task := &model.Task{}
		if err := c.Call("task.move", request, task); err != nil {
			cobra.CheckErr(err)
		}

		t.PrintTasks(client.TaskSlice{task}, "task", client.Change)
	},
}

var taskCloseCmd = &cobra.Command{
	Use:     "close",
	Aliases: []string{"c"},
	Short:   "close task",
	Run: func(cmd *cobra.Command, args []string) {
		t := cli.NewCLI(cli.Nerd)

		request := &struct {
			ID int64 `json:"id"`
		}{}

		id, err := cmd.Flags().GetInt64("id")
		if err != nil {
			cobra.CheckErr(err)
		}
		request.ID = id

		c := rpc.NewClient(network, socketPath, 5*time.Second)

		tasks := client.TaskSlice{}
		if err := c.Call("task.list", nil, &tasks); err != nil {
			cobra.CheckErr(err)
		}

		var name *string
		if len(args) == 1 {
			name = &args[0]
		} else {
			name = getValue(cmd, cmd.Flags().GetString, "name")
		}

		if id, err := t.SelectOne(tasks.FilterByName(*name)); err != nil {
			cobra.CheckErr(err)
		} else {
			request.ID = id
		}

		task := &model.Task{}
		if err := c.Call("task.close", request, task); err != nil {
			cobra.CheckErr(err)
		}

		t.PrintTasks(client.TaskSlice{task}, "task", client.Change)
	},
}

var taskDeleteCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm"},
	Args:    cobra.ExactArgs(1),
	Short:   "remove task",
	Run: func(cmd *cobra.Command, args []string) {
		cl := cli.NewCLI(cli.Nerd)

		request := &struct {
			ID int64 `json:"id"`
		}{}

		id, err := cmd.Flags().GetInt64("id")
		if err != nil {
			cobra.CheckErr(err)
		}
		request.ID = id

		c := rpc.NewClient(network, socketPath, 5*time.Second)

		tasks := client.TaskSlice{}
		if err := c.Call("task.list", nil, &tasks); err != nil {
			cobra.CheckErr(err)
		}

		var name *string
		if len(args) == 1 {
			name = &args[0]
		} else {
			name = getValue(cmd, cmd.Flags().GetString, "name")
		}

		ids, err := cl.SelectMultiple(tasks.FilterByName(*name))
		if err != nil {
			cobra.CheckErr(err)
		}

		cl.PrintTasks(tasks.FilterByIDs(ids), "task", client.Delete)
		for _, id := range ids {
			request.ID = id
			if err := c.Call("task.delete", request, nil); err != nil {
				cobra.CheckErr(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(taskListCmd)
	rootCmd.AddCommand(taskAddCmd)
	rootCmd.AddCommand(taskUpdateCmd)
	rootCmd.AddCommand(taskMoveCmd)
	rootCmd.AddCommand(taskCloseCmd)
	rootCmd.AddCommand(taskDeleteCmd)

	taskListCmd.Flags().StringP("name", "n", "", "filter task by name")

	taskAddCmd.Flags().StringP("name", "n", "", "Task name")
	taskAddCmd.Flags().StringP("description", "D", "", "Task description")
	taskAddCmd.Flags().StringP("due", "d", "", "Task labels")
	taskAddCmd.Flags().BoolP("recurring", "r", false, "Task recurring")
	taskAddCmd.Flags().String("duration", "", "Task duration") // TODO: implement duration
	taskAddCmd.Flags().IntP("priority", "p", 0, "Task priority")
	taskAddCmd.Flags().Int64("project", 0, "Task project ID")
	taskAddCmd.Flags().Int64("parent", 0, "Task parent ID")
	taskAddCmd.Flags().StringSliceP("labels", "l", []string{}, "Task labels")
	taskAddCmd.Flags().BoolP("edit", "e", false, "Eidt task")

	taskMoveCmd.Flags().Int64P("id", "i", 0, "Task ID")
	taskMoveCmd.Flags().StringP("name", "n", "", "Task name")
	taskMoveCmd.Flags().Int64P("destination", "d", 0, "destination index")
	taskMoveCmd.Flags().StringP("parent", "p", "", "project or parent task name")

	taskCloseCmd.Flags().Int64P("id", "i", 0, "Task ID")
	taskCloseCmd.Flags().StringP("name", "n", "", "Task name")

	taskDeleteCmd.Flags().Int64P("id", "i", 0, "Task ID")
	taskDeleteCmd.Flags().StringP("name", "n", "", "Task name")
}
