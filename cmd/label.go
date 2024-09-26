package cmd

import (
	"github.com/CnTeng/rx-todo/cli"
	"github.com/CnTeng/rx-todo/client"
	"github.com/CnTeng/rx-todo/internal/model"
	"github.com/spf13/cobra"
)

var labelCmd = &cobra.Command{
	Use:   "label",
	Short: "list all labels",
	Run: func(cmd *cobra.Command, args []string) {
		s, err := client.NewStorage("rx-todo/resources.json")
		if err != nil {
			cli.PrintErr(err)
		}

		c := client.NewClient(config.Address, config.Token, s.Path)

		if res, err := c.Sync(); err != nil {
		} else {
			labels := cli.LabelSlice(*res.Labels)
			cmd.Print(labels.List(nil))
			c.Store(res)
		}
	},
}

var labelAddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a"},
	Short:   "add new label",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		color, _ := cmd.Flags().GetString("color")

		s, err := client.NewStorage("rx-todo/resources.json")
		if err != nil {
			cli.PrintErr(err)
			return
		}

		c := client.NewClient(config.Address, config.Token, s.Path)

		label, err := c.CreateLabel(name, color)
		if err != nil {
			cli.PrintErr(err)
			return
		}

		sm := cli.NewStatusMap([]*model.Resource{&label.Resource}, cli.Add)
		res, _ := c.Patch([]*model.Label{label})

		labels := cli.LabelSlice(*res.Labels)
		cmd.Print(labels.List(sm))
	},
}

func init() {
	rootCmd.AddCommand(labelCmd)

	labelCmd.AddCommand(labelAddCmd)
	labelAddCmd.Flags().StringP("name", "n", "", "name of the label")
	labelAddCmd.Flags().StringP("color", "c", "", "color of the label")
}
