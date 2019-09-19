package main

import (
	"fmt"
	"os"
	"strings"

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

		headerStrs, _ := cmd.Flags().GetStringSlice("header")
		var headers []client.Header
		for _, v := range headerStrs {
			headers = append(headers, client.Header{Name: strings.Split(v, ":")[0], Value: strings.Split(v, ":")[1]})
		}

		client.Call(client.HttpRequest{Method: "GET", URI: args[0], Headers: headers}, profile)
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
	rootCmd.Flags().StringSliceP("header", "H", nil, "header value")
	rootCmd.Flags().StringP("profile", "P", "", "authentication profile")
	rootCmd.Flags().StringP("body", "b", "", "body")
}
