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

		url, err := cmd.Flags().GetString("url")
		if err != nil {
			os.Exit(1)
		}

		client.Call(client.HttpRequest{"GET", url})
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
}
