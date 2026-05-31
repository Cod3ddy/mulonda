package cmd

import (
	"fmt"
	"strings"

	"github.com/cod3ddy/mulonda/internal/watchlist"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove <command>",
	Short: "Remove a rule from the watchlist",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		command := strings.TrimSpace(args[0])
		if command == "" {
			return fmt.Errorf("command cannot be empty")
		}

		removed, err := watchlist.RemoveRule(watchlistFile, command)
		if err != nil {
			return err
		}

		if removed {
			fmt.Printf("removed from watchlist: %s\n", command)
		} else {
			fmt.Printf("not found in watchlist: %s\n", command)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
