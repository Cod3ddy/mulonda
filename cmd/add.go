package cmd

import (
	"fmt"
	"strings"

	"github.com/cod3ddy/mulonda/internal/watchlist"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <command>",
	Short: "Add a command pattern to the watchlist",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		rule := strings.TrimSpace(strings.Join(args, " "))
		if rule == "" {
			return fmt.Errorf("command cannot be empty")
		}

		if err := watchlist.Add(watchlistFile, rule); err != nil {
			return err
		}

		fmt.Printf("added to watchlist: %s\n", rule)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
