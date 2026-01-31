package cmd

import (
	"fmt"
	"os"

	"github.com/junipery17/lineartui/internal/client"
	"github.com/junipery17/lineartui/internal/config"
	"github.com/spf13/cobra"
)

// TODO: use dependency injection instead of global vars?
var (
	cfgFile      string
	cfg          *config.Config
	linearClient client.Client
)

var rootCmd = &cobra.Command{
	Use:   "lineartui",
	Short: "A TUI for Linear",
	Long:  `lineartui is a terminal user interface for interacting with Linear project management.`,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		cfg, err = config.New()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if cfg.Linear.APIKey == "" {
			return fmt.Errorf("Missing API key in configuration chain.")
		}

		linearClient = client.NewClient(cfg.Linear.APIKey, cfg.Linear.APIURL)

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "configfile", "", "config file (default is ./.lcli.yaml or $HOME/.lcli.yaml)")
}
