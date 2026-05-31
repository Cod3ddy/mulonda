package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cod3ddy/mulonda/internal/watchlist"
	"github.com/spf13/cobra"
)

const (
	aliasBlockStart = "# >>> mulonda >>>"
	aliasBlockEnd   = "# <<< mulonda <<<"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Mulonda shell aliases",
	RunE: func(cmd *cobra.Command, args []string) error {
		shell := detectShell()
		cmds := defaultWatchedCommands()
		block := buildAliasBlock(shell, cmds)

		rcFiles := shellRCFiles(shell)
		for _, rc := range rcFiles {
			if err := injectAliasBlock(rc, block); err != nil {
				return fmt.Errorf("inject into %s: %w", rc, err)
			}
		}

		fmt.Println("Mulonda aliases installed.")
		fmt.Println("Restart your shell or run:")
		for _, f := range rcFiles {
			fmt.Printf("  source %s\n", f)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

func detectShell() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return "bash"
	}
	return filepath.Base(shell)
}

// shellRCFiles returns the rc file paths for the given shell.
func shellRCFiles(shell string) []string {
	home, _ := os.UserHomeDir()
	switch shell {
	case "zsh":
		return []string{filepath.Join(home, ".zshrc")}
	case "fish":
		cfgDir, _ := os.UserConfigDir()

		return []string{filepath.Join(cfgDir, "fish", "conf.d", "mulonda.fish")}
	default:
		return []string{filepath.Join(home, ".bashrc")}
	}
}

func defaultWatchedCommands() []string {
	cmds := make([]string, 0, len(watchlist.DefaultRules))
	for _, r := range watchlist.DefaultRules {
		cmds = append(cmds, r.Command)
	}
	return cmds
}

func buildAliasBlock(shell string, cmds []string) string {
	var sb strings.Builder
	sb.WriteString(aliasBlockStart + "\n")
	for _, c := range cmds {
		if shell == "fish" {
			fmt.Fprintf(&sb, "alias %s \"mulonda %s\"\n", c, c)
		} else {
			fmt.Fprintf(&sb, "alias %s=\"mulonda %s\"\n", c, c)
		}
	}
	sb.WriteString(aliasBlockEnd + "\n")
	return sb.String()
}

func injectAliasBlock(rcFile, block string) error {
	content, err := os.ReadFile(rcFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	s := string(content)

	if strings.Contains(s, aliasBlockStart) {
		return os.WriteFile(rcFile, []byte(replaceAliasBlock(s, block)), 0o644)
	}

	if err := os.MkdirAll(filepath.Dir(rcFile), 0o755); err != nil {
		return err
	}

	f, err := os.OpenFile(rcFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	prefix := "\n"
	if len(s) == 0 {
		prefix = ""
	}
	_, err = fmt.Fprintf(f, "%s%s\n", prefix, block)
	return err
}

func replaceAliasBlock(content, newBlock string) string {
	before, _, ok := strings.Cut(content, aliasBlockStart)
	end := strings.Index(content, aliasBlockEnd)
	if !ok || end == -1 {
		return content
	}
	end += len(aliasBlockEnd)
	if end < len(content) && content[end] == '\n' {
		end++
	}
	return before + newBlock + "\n" + content[end:]
}
