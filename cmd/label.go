package cmd

import (
	"time"

	"github.com/CnTeng/rx-todo/cli"
	"github.com/CnTeng/rx-todo/client"
	"github.com/CnTeng/rx-todo/model"
	"github.com/CnTeng/rx-todo/rpc"
	"github.com/spf13/cobra"
)

var labelCmd = &cobra.Command{
	Use:   "label",
	Short: "list all labels",
	Run: func(cmd *cobra.Command, args []string) {
		c := rpc.NewClient(network, socketPath, 5*time.Second)

		labels := client.LabelSlice{}
		if err := c.Call("Label.List", nil, &labels); err != nil {
			cobra.CheckErr(err)
		}

		cli.NewCLI(cli.Nerd).PrintLabels(labels, "Labels", client.None)
	},
}

var labelAddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a"},
	Short:   "add new label",
	Run: func(cmd *cobra.Command, args []string) {
		request := &model.LabelCreationRequest{
			Name:  getValue(cmd, cmd.Flags().GetString, "name"),
			Color: getValue(cmd, cmd.Flags().GetString, "color"),
		}

		c := rpc.NewClient(network, socketPath, 5*time.Second)

		label := &model.Label{}
		if err := c.Call("Label.Create", request, label); err != nil {
			cobra.CheckErr(err)
		}

		cli.NewCLI(cli.Nerd).PrintLabels(client.LabelSlice{label}, "Labels", client.Add)
	},
}

var labelModifyCmd = &cobra.Command{
	Use:     "modify",
	Aliases: []string{"m"},
	Short:   "modify label",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		c := rpc.NewClient(network, socketPath, 5*time.Second)

		request := &model.LabelUpdateRequestWithID{
			ID: getValue(cmd, cmd.Flags().GetInt64, "id"),
			LabelUpdateRequest: &model.LabelUpdateRequest{
				Name:  getValue(cmd, cmd.Flags().GetString, "name"),
				Color: getValue(cmd, cmd.Flags().GetString, "color"),
			},
		}

		label := &model.Label{}
		if err := c.Call("Label.Update", request, &label); err != nil {
			cobra.CheckErr(err)
		}

		cli.NewCLI(cli.Nerd).PrintLabels(client.LabelSlice{label}, "Labels", client.Change)
	},
}

func init() {
	rootCmd.AddCommand(labelCmd)

	labelCmd.AddCommand(labelAddCmd)
	labelAddCmd.Flags().StringP("name", "n", "", "name of the label")
	labelAddCmd.Flags().StringP("color", "c", "", "color of the label")

	labelCmd.AddCommand(labelModifyCmd)
	labelModifyCmd.Flags().Int64P("id", "i", 0, "ID of the label")
	labelModifyCmd.Flags().StringP("name", "n", "", "name of the label")
	labelModifyCmd.Flags().StringP("color", "c", "", "color of the label")
}
