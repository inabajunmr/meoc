package main

import (
	"fmt"
	"os"

	"github.com/inabajunmr/meoc/client"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "API Client for authorized endpoint by OAuth 2.0",
	Run: func(cmd *cobra.Command, args []string) {

		profile, err := cmd.Flags().GetString("profile")
		if err != nil {
			os.Exit(1)
		}

		client.Call(client.HttpRequest{Method: "GET", URI: args[0]}, profile)
	},
}

// Execute cmd
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("url", "u", "", "endpoint")
	rootCmd.Flags().StringP("profile", "p", "", "authentication profile")
}
