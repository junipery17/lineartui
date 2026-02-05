package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var issuesCmd = &cobra.Command{
	Use:   "issues",
	Short: "Manage Linear issues",
	Long:  `List, create, and manage issues in Linear.`,
}

var issuesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List issues for a team",
	RunE: func(cmd *cobra.Command, args []string) error {
		teamID, _ := cmd.Flags().GetString("team")
		if teamID == "" {
			teamID = cfg.Linear.TeamID
		}
		if teamID == "" {
			return fmt.Errorf("team ID required. Use --team flag or set linear.team_id in config")
		}
		// TODO: implement issues listing
		fmt.Printf("Listing issues for team %s...\n", teamID)
		return nil
	},
}

var issuesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new issue",
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg.Linear.APIKey == "" {
			return fmt.Errorf("linear API key not configured. Set LINEARTUI_LINEAR_API_KEY or add it to config.yaml")
		}
		title, _ := cmd.Flags().GetString("title")
		if title == "" {
			return fmt.Errorf("title is required")
		}
		teamID, _ := cmd.Flags().GetString("team")
		if teamID == "" {
			teamID = cfg.Linear.TeamID
		}
		if teamID == "" {
			return fmt.Errorf("team ID required. Use --team flag or set linear.team_id in config")
		}
		description, _ := cmd.Flags().GetString("description")
		// TODO: implement issue creation
		fmt.Printf("Creating issue '%s' in team %s...\n", title, teamID)
		if description != "" {
			fmt.Printf("Description: %s\n", description)
		}
		return nil
	},
}

var issuesDeleteCmd = &cobra.Command{
	Use:   "delete [issue-id]",
	Short: "Delete an issue",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg.Linear.APIKey == "" {
			return fmt.Errorf("linear API key not configured. Set LINEARTUI_LINEAR_API_KEY or add it to config.yaml")
		}
		issueID := args[0]
		// TODO: implement issue deletion
		fmt.Printf("Deleting issue %s...\n", issueID)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(issuesCmd)
	issuesCmd.AddCommand(issuesListCmd)
	issuesCmd.AddCommand(issuesCreateCmd)
	issuesCmd.AddCommand(issuesDeleteCmd)

	// Flags for list command
	issuesListCmd.Flags().StringP("team", "t", "", "Team ID to list issues for")

	// Flags for create command
	issuesCreateCmd.Flags().StringP("title", "T", "", "Issue title (required)")
	issuesCreateCmd.Flags().StringP("description", "d", "", "Issue description")
	issuesCreateCmd.Flags().StringP("team", "t", "", "Team ID to create issue in")
	issuesCreateCmd.MarkFlagRequired("title")
}
