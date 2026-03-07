package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove Mulonda shell aliases",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("uninstall stub: TODO remove aliases from shell rc files")
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
