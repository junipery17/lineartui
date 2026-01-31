package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

var teamsCmd = &cobra.Command{
	Use:   "teams",
	Short: "List Linear teams",
	Long:  `Display all teams you have access to in Linear.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		return linearClient.DisplayTeams(ctx)
	},
}

func init() {
	rootCmd.AddCommand(teamsCmd)
}
