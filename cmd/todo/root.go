package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type Config struct {
	Address string `json:"address"`
	Token   string `json:"token"`
}

var cfgFile string

var config Config

const (
	network    = "unix"
	socketPath = "@rx-todo.sock"
)

var rootCmd = &cobra.Command{
	Use:     "todo",
	Aliases: []string{"rx-todo"},
	Short:   "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is config.json)")
}

func getValue[T any](cmd *cobra.Command, getter func(string) (T, error), name string) *T {
	if !cmd.Flags().Changed(name) {
		return nil
	}

	v, err := getter(name)
	if err != nil {
		cobra.CheckErr(err)
	}

	return &v
}

func initConfig() {
	file, err := os.ReadFile(cfgFile)
	if err != nil {
		fmt.Printf("Error opening config file: %v\n", err)
		os.Exit(1)
	}

	if err := json.Unmarshal(file, &config); err != nil {
		fmt.Printf("Error decoding config file: %v\n", err)
		os.Exit(1)
	}
}
