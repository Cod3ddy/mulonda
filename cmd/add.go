package cmd

import (
	"fmt"
	"strings"

	"github.com/cod3ddy/mulonda/internal/watchlist"
	"github.com/spf13/cobra"
)

var (
	addWarning      string
	addFlagsContain []string
	addArgsMatch    []string
	addMinArgs      int
)

var addCmd = &cobra.Command{
	Use:   "add <command> [args_match...]",
	Short: "Add or update a structured watchlist rule",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		command := strings.TrimSpace(args[0])
		if command == "" {
			return fmt.Errorf("command cannot be empty")
		}

		if addMinArgs < 0 {
			return fmt.Errorf("--min-args cannot be less than 0, try adding something  dude , >1 atleast.")
		}

		rule := watchlist.Rule{
			Command:      command,
			Warning:      strings.TrimSpace(addWarning),
			FlagsContain: cleaned(addFlagsContain),
			MinArgs:      addMinArgs,
		}

		// This is somewhat backward-compatible, for instance when now you run
		// `mulonda add chmod 777` => args_match: ["777"].
		if matches := cleaned(addArgsMatch); len(matches) > 0 {
			rule.ArgsMatch = matches
		} else if len(args) > 1 {
			rule.ArgsMatch = cleaned(args[1:])
		}

		if err := watchlist.AddRule(watchlistFile, rule); err != nil {
			return err
		}

		fmt.Printf("added to watchlist: %s\n", rule.Command)
		return nil
	},
}

func init() {
	addCmd.Flags().StringVar(&addWarning, "warning", "", "Warning message shown before execution")
	addCmd.Flags().StringSliceVar(&addFlagsContain, "flags-contain", nil, "Match if any of these flags are present")
	addCmd.Flags().StringSliceVar(&addArgsMatch, "args-match", nil, "Match if any of these args are present")
	addCmd.Flags().IntVar(&addMinArgs, "min-args", 0, "Minimum number of args required for the rule to match")
	rootCmd.AddCommand(addCmd)
}

func cleaned(values []string) []string {
	out := make([]string, 0, len(values))
	for _, v := range values {
		s := strings.TrimSpace(v)
		if s == "" {
			continue
		}
		out = append(out, s)
	}
	return out
}
