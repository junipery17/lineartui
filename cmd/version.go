package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: use go build time variables to set this
		fmt.Printf("%s v%s\n", "LinearTUI", "0.0.1")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
