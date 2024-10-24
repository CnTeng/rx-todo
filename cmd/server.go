package cmd

import (
	"github.com/CnTeng/rx-todo/database"
	"github.com/CnTeng/rx-todo/server"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := database.NewDB("postgresql:///rxtodo?host=/run/postgresql")
		if err != nil {
			panic(err)
		}
		defer db.Close()

		// if err := db.Migrate(); err != nil {
		// 	panic(err)
		// }

		server := server.NewServer(&db)
		_ = server.Start()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
