package cmd

import (
	"time"

	"github.com/CnTeng/rx-todo/cli"
	"github.com/CnTeng/rx-todo/model"
	"github.com/CnTeng/rx-todo/rpc"
	"github.com/spf13/cobra"
)

var (
	id    int64
	name  string
	color string
)

var labelCmd = &cobra.Command{
	Use:   "label",
	Short: "list all labels",
	Run: func(cmd *cobra.Command, args []string) {
		c := rpc.NewClient("unix", "/tmp/rx-todo.sock", 5*time.Second)

		labels := cli.LabelSlice{}
		if err := c.Call("Label.List", nil, &labels); err != nil {
			cobra.CheckErr(err)
		}

		cmd.Print(labels.List(nil))
	},
}

var labelAddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a"},
	Short:   "add new label",
	Run: func(cmd *cobra.Command, args []string) {
		c := rpc.NewClient("unix", "/tmp/rx-todo.sock", 5*time.Second)

		labelCreationRequest := &model.LabelCreationRequest{Name: name, Color: color}
		label := &model.Label{}
		if err := c.Call("Label.Create", labelCreationRequest, label); err != nil {
			cobra.CheckErr(err)
		}
		sm := cli.NewStatusMap([]*model.Label{label}, cli.Add)

		labels := cli.LabelSlice{}
		if err := c.Call("Label.List", nil, &labels); err != nil {
			cobra.CheckErr(err)
		}

		cmd.Print(labels.List(sm))
	},
}

var labelModifyCmd = &cobra.Command{
	Use:     "modify",
	Aliases: []string{"m"},
	Short:   "modify label",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		c := rpc.NewClient("unix", "/tmp/rx-todo.sock", 5*time.Second)

		labelUpdateRequest := &model.LabelUpdateRequestWithID{
			ID: id,
			LabelUpdateRequest: model.LabelUpdateRequest{
				Name:  &name,
				Color: &color,
			},
		}
		if err := model.Validate(labelUpdateRequest); err != nil {
			cobra.CheckErr(err)
		}

		label := &model.Label{}
		if err := c.Call("Label.Update", labelUpdateRequest, &label); err != nil {
			cobra.CheckErr(err)
		}
		sm := cli.NewStatusMap([]*model.Label{label}, cli.Change)

		labels := cli.LabelSlice{}
		if err := c.Call("Label.List", nil, &labels); err != nil {
			cobra.CheckErr(err)
		}

		cmd.Print(labels.List(sm))
	},
}

func init() {
	rootCmd.AddCommand(labelCmd)

	labelCmd.AddCommand(labelAddCmd)
	labelAddCmd.Flags().StringVarP(&name, "name", "n", "", "name of the label")
	labelAddCmd.Flags().StringVarP(&color, "color", "c", "", "color of the label")

	labelCmd.AddCommand(labelModifyCmd)
	labelModifyCmd.Flags().Int64VarP(&id, "id", "i", 0, "ID of the label")
	labelModifyCmd.Flags().StringVarP(&name, "name", "n", "", "name of the label")
	labelModifyCmd.Flags().StringVarP(&color, "color", "c", "", "color of the label")
}
