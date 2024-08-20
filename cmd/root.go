package cmd

import (
	"errors"
	"fmt"
	"github-user-activity/activity"
	"github.com/spf13/cobra"
)

var username string
var perPage int
var page int

func init() {
	rootCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "username for github user")
	rootCmd.PersistentFlags().IntVarP(&perPage, "perPage", "", 30, "number of items to list")
	rootCmd.PersistentFlags().IntVarP(&page, "page", "", 1, "number of items to list")
}

var rootCmd = &cobra.Command{
	Use:   "github-user-activity",
	Short: "Display GitHub User Activity",
	RunE: func(cmd *cobra.Command, args []string) error {
		if username == "" {
			return errors.New("username is required")
		}
		activity.ListAllActivities(username, perPage, page)
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
