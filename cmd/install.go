package cmd

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Mulonda shell aliases",
	RunE: func(cmd *cobra.Command, args []string) error {
		install := exec.Command("bash", "scripts/install.sh")
		install.Stdout = os.Stdout
		install.Stderr = os.Stderr
		install.Stdin = os.Stdin
		return install.Run()
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
