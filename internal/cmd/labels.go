package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var labelsCmd = &cobra.Command{
	Use:   "labels",
	Short: "Manage Linear labels",
	Long:  `List, create, and delete labels in Linear.`,
}

var labelsListCmd = &cobra.Command{
	Use:   "list labels",
	Short: "List existing labels",
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg.Linear.APIKey == "" {
			return fmt.Errorf("linear API key not configured. Set LINEARTUI_LINEAR_API_KEY or add it to config.yaml")
		}
		issueID, _ := cmd.Flags().GetString("issueID")
		title, _ := cmd.Flags().GetString("title")
		if title != "" {
			var err error
			issueID, err = linearClient.FindIssueByTitle(context.Background(), cfg.Linear.TeamID, title)
			if err != nil {
				return err
			}
			fmt.Printf("Listing labels applied to issue: %s\n", title)
		}
		if issueID != "" {
			if title == "" {
				fmt.Printf("listing labels applied to issue: %s\n", issueID)
			}
			return linearClient.ListLabels(context.Background(), issueID)
		}
		fmt.Print("Listing all existing labels\n")
		return linearClient.ListLabels(context.Background(), "")
	},
}

// apply a given label to multiple issues?

func init() {
	rootCmd.AddCommand(labelsCmd)
	labelsCmd.AddCommand(labelsListCmd)

	labelsListCmd.Flags().StringP("issueID", "i", "", "issueID to list labels for")
	labelsListCmd.Flags().StringP("title", "t", "", "title of issue to list labels for")
	labelsListCmd.MarkFlagsMutuallyExclusive("issueID", "title")
}
