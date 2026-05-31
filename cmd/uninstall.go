package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove Mulonda shell aliases",
	RunE: func(cmd *cobra.Command, args []string) error {
		shells := []string{"bash", "zsh", "fish"}
		removed := []string{}

		for _, shell := range shells {
			for _, rc := range shellRCFiles(shell) {
				ok, err := removeAliasBlock(rc)
				if err != nil {
					return fmt.Errorf("uninstall from %s: %w", rc, err)
				}
				if ok {
					removed = append(removed, rc)
				}
			}
		}

		if len(removed) == 0 {
			fmt.Println("No Mulonda aliases found.")
			return nil
		}

		fmt.Println("Mulonda aliases removed from:")
		for _, f := range removed {
			fmt.Printf("  %s\n", f)
		}
		fmt.Println("Restart your shell to apply.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}

func removeAliasBlock(rcFile string) (bool, error) {
	content, err := os.ReadFile(rcFile)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	s := string(content)
	if !strings.Contains(s, aliasBlockStart) {
		return false, nil
	}

	start := strings.Index(s, aliasBlockStart)
	end := strings.Index(s, aliasBlockEnd)
	if start == -1 || end == -1 {
		return false, nil
	}

	end += len(aliasBlockEnd)
	if end < len(s) && s[end] == '\n' {
		end++
	}

	// just remove a preceding blank line if the block was appended with one
	if start > 0 && s[start-1] == '\n' {
		start--
	}

	result := s[:start] + s[end:]
	return true, os.WriteFile(rcFile, []byte(result), 0o644)
}
