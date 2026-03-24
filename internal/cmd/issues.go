package cmd

import (
	"context"
	"fmt"
	"strconv"
	"strings"

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
		var teamID string
		teamName, _ := cmd.Flags().GetString("team")
		if teamName == "" {
			teamID = cfg.Linear.TeamID
		} else {
			var err error
			teamID, err = linearClient.FindTeamByName(context.Background(), teamName)
			if err != nil {
				return err
			}
		}
		if teamID == "" {
			return fmt.Errorf("team ID required. Use --team flag or set linear.team_id in config")
		}
		titlesOnly, _ := cmd.Flags().GetBool("titles")
		ctx := context.Background()
		fmt.Printf("Listing issues for team %s...\n", teamID)
		return linearClient.DisplayIssues(ctx, teamID, titlesOnly)
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
		fmt.Printf("Creating issue '%s' in team %s...\n", title, teamID)
		if description != "" {
			fmt.Printf("Description: %s\n", description)
		}
		ctx := context.Background()
		_, err := linearClient.AddIssue(ctx, teamID, title, description)

		return err
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
		ctx := context.Background()
		return linearClient.DeleteIssue(ctx, issueID)
	},
}

var issuesUpdateCmd = &cobra.Command{
	Use:   "update [issue-id]",
	Short: "Modify an existing issue",
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg.Linear.APIKey == "" {
			return fmt.Errorf("linear API key not configured. Set LINEARTUI_LINEAR_API_KEY or add it to config.yaml")
		}
		issueID, _ := cmd.Flags().GetString("issueID")
		title, _ := cmd.Flags().GetString("titleSearch")
		if title != "" {
			var err error
			issueID, err = linearClient.FindIssueByTitle(context.Background(), cfg.Linear.TeamID, title)
			if err != nil {
				return err
			}
		}
		assign, _ := cmd.Flags().GetString("assign")
		if assign != "" {
			fmt.Printf("Updating issue %s...\n", issueID)
			err := linearClient.UpdateAssigneeOnIssue(context.Background(), issueID, assign)
			if err != nil {
				return err
			}
		}
		description, _ := cmd.Flags().GetString("description")
		if description != "" {
			fmt.Printf("Updating description on issue %s...\n", issueID)
			err := linearClient.UpdateDescriptionOnIssue(context.Background(), issueID, description)
			if err != nil {
				return err
			}
		}
		priority, _ := cmd.Flags().GetString("priority")
		if priority != "" {
			priority, parseErr := strconv.ParseInt(priority, 10, 64)
			if parseErr != nil || priority < 0 || priority > 4 {
				return fmt.Errorf("Priority must be an integer from 0 to 4")
			}
			err := linearClient.UpdatePriorityOnIssue(context.Background(), issueID, float64(priority))
			if err != nil {
				return err
			}
		}
		status, _ := cmd.Flags().GetString("status")
		if status != "" {
			status = strings.ToLower(status)
			statID := StatusToID[status]
			err := linearClient.UpdateStatusOnIssue(context.Background(), issueID, statID)
			if err != nil {
				return err
			}
		}
		return nil
	},
}

var issueUpdateLabelCmd = &cobra.Command{
	Use:   "label [issue-id]",
	Short: "Update and edit labels on issue",
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg.Linear.APIKey == "" {
			return fmt.Errorf("linear API key not configured. Set LINEARTUI_LINEAR_API_KEY or add it to config.yaml")
		}
		issueID, _ := cmd.Flags().GetString("issueID")
		title, _ := cmd.Flags().GetString("titleSearch")
		if title != "" {
			var err error
			issueID, err = linearClient.FindIssueByTitle(context.Background(), cfg.Linear.TeamID, title)
			if err != nil {
				return err
			}
		}
		add, _ := cmd.Flags().GetString("add")
		if add != "" {
			err := linearClient.AddLabeltoIssue(context.Background(), issueID, add)
			if err != nil {
				return err
			}
		}
		remove, _ := cmd.Flags().GetString("remove")
		if remove != "" {
			err := linearClient.RemoveLabelFromIssue(context.Background(), issueID, remove)
			if err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(issuesCmd)
	issuesCmd.AddCommand(issuesListCmd)
	issuesCmd.AddCommand(issuesCreateCmd)
	issuesCmd.AddCommand(issuesDeleteCmd)
	issuesCmd.AddCommand(issuesUpdateCmd)
	issuesCmd.AddCommand(issueUpdateLabelCmd)

	// Flags for list command
	issuesListCmd.Flags().StringP("team", "t", "", "Team Name to list issues for")
	issuesListCmd.Flags().BoolP("titles", "T", false, "List only titles of Issues")

	// Flags for create command
	issuesCreateCmd.Flags().StringP("title", "T", "", "Issue title (required)")
	issuesCreateCmd.Flags().StringP("description", "d", "", "Issue description")
	issuesCreateCmd.Flags().StringP("team", "t", "", "Team ID to create issue in")
	issuesCreateCmd.MarkFlagRequired("title")

	//Flags for updating Issue command
	issuesUpdateCmd.Flags().StringP("assign", "a", "", "Edit assigned members")
	issuesUpdateCmd.Flags().StringP("description", "d", "", "Edit description")
	issuesUpdateCmd.Flags().StringP("priority", "p", "", "Set new priority for issue")
	issuesUpdateCmd.Flags().StringP("issueID", "i", "", "ID of issue to update")
	issuesUpdateCmd.Flags().StringP("titleSearch", "t", "", "Select issue by title")
	issuesUpdateCmd.Flags().StringP("status", "s", "", "Update status of issue")
	issuesUpdateCmd.MarkFlagsOneRequired("issueID", "titleSearch")
	issuesUpdateCmd.MarkFlagsMutuallyExclusive("issueID", "titleSearch")

	//Flags for labels
	issueUpdateLabelCmd.Flags().StringP("issueID", "i", "", "ID of issue to edit label")
	issueUpdateLabelCmd.Flags().StringP("titleSearch", "t", "", "Issue by title")
	issueUpdateLabelCmd.Flags().StringP("add", "a", "", "Add a label")
	issueUpdateLabelCmd.Flags().StringP("remove", "r", "", "Remove a label")
	issueUpdateLabelCmd.MarkFlagsOneRequired("issueID", "titleSearch")
	issueUpdateLabelCmd.MarkFlagsMutuallyExclusive("issueID", "titleSearch")

}
