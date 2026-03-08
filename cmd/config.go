package cmd

import (
	"fmt"
	"strings"

	"github.com/cod3ddy/mulonda/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Mulonda configuration",
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Print effective config values",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(configFile)
		if err != nil {
			return err
		}

		fmt.Printf("timeout_seconds: %d\n", cfg.TimeoutSeconds)
		fmt.Printf("non_interactive.passthrough: %t\n", cfg.NonInteractive.Passthrough)
		fmt.Printf("files.watchlist_dir: %s\n", cfg.Files.WatchlistDirectory)
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set config value (stub)",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := strings.TrimSpace(args[0])
		value := strings.TrimSpace(args[1])
		fmt.Printf("config set stub: TODO persist %s=%s in %s\n", key, value, configFile)
		return nil
	},
}

func init() {
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
	rootCmd.AddCommand(configCmd)
}
