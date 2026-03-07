package cmd

import (
	"fmt"

	"github.com/cod3ddy/mulonda/internal/watchlist"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type listOutput struct {
	Rules []watchlist.Rule `yaml:"rules"`
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show active watchlist commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		rules, err := watchlist.Load(watchlistFile)
		if err != nil {
			return err
		}

		out := listOutput{Rules: rules}
		payload, err := yaml.Marshal(out)
		if err != nil {
			return fmt.Errorf("marshal watchlist: %w", err)
		}

		_, err = cmd.OutOrStdout().Write(payload)
		if err != nil {
			return fmt.Errorf("write output: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
