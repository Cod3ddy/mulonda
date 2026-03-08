package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/cod3ddy/mulonda/internal/config"
	"github.com/cod3ddy/mulonda/internal/executor"
	"github.com/cod3ddy/mulonda/internal/matcher"
	"github.com/cod3ddy/mulonda/internal/prompter"
	"github.com/cod3ddy/mulonda/internal/watchlist"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mulonda",
	Short: "Guard dangerous shell commands with a confirmation prompt",
	Long: `Mulonda is a lightweight CLI that proxies selected shell commands
and asks for confirmation before potentially destructive operations.`,
}

var configFile string
var watchlistFile string

var (
	loadConfigFn     = config.Load
	loadWatchlistFn  = watchlist.Load
	matchRuleFn      = matcher.MatchRule
	confirmPromptFn  = prompter.Confirm
	commandFactoryFn = executor.Command
	isInteractiveFn  = isInteractiveSession
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	configPath, watchlistPath, passthroughArgs, parseErr := parseGlobalFlags(os.Args[1:])
	if parseErr == nil && len(passthroughArgs) > 0 && !isManagementCommand(passthroughArgs[0]) {
		if err := executeProxy(configPath, watchlistPath, passthroughArgs); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "mulonda.yaml", "Path to YAML config file")
	rootCmd.PersistentFlags().StringVar(&watchlistFile, "watchlist", "watchlist.yaml", "Path to YAML watchlist file")
}

func executeProxy(configPath, watchlistPath string, args []string) error {
	command := args[0]
	commandArgs := args[1:]

	cfg, err := loadConfigFn(configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	rules, err := loadWatchlistFn(watchlistPath)
	if err != nil {
		return fmt.Errorf("load watchlist: %w", err)
	}

	if rule, matched := matchRuleFn(command, commandArgs, rules); matched {
		warning := strings.TrimSpace(rule.Warning)
		if warning == "" {
			warning = "Watched command requires confirmation"
		}

		if isInteractiveFn() {
			fullCommand := command
			if len(commandArgs) > 0 {
				fullCommand = fmt.Sprintf("%s %s", command, strings.Join(commandArgs, " "))
			}
			message := fmt.Sprintf("mulonda: %s\nwarning: %s\nProceed?", fullCommand, warning)
			if !confirmPromptFn(strings.TrimSpace(message)) {
				return fmt.Errorf("aborted: command not executed")
			}
		} else if !cfg.NonInteractive.Passthrough {
			return fmt.Errorf("blocked in non-interactive mode: %s", command)
		}
	}

	cmd := commandFactoryFn(command, commandArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("execute command: %w", err)
	}

	return nil
}

func parseGlobalFlags(args []string) (configPath, watchPath string, passthrough []string, err error) {
	configPath = "mulonda.yaml"
	watchPath = "watchlist.yaml"

	for i := 0; i < len(args); i++ {
		token := args[i]

		if token == "--" {
			return configPath, watchPath, args[i+1:], nil
		}

		if token == "--config" {
			if i+1 >= len(args) {
				return "", "", nil, fmt.Errorf("missing value for --config")
			}
			configPath = args[i+1]
			i++
			continue
		}

		if strings.HasPrefix(token, "--config=") {
			configPath = strings.TrimPrefix(token, "--config=")
			continue
		}

		if token == "--watchlist" {
			if i+1 >= len(args) {
				return "", "", nil, fmt.Errorf("missing value for --watchlist")
			}
			watchPath = args[i+1]
			i++
			continue
		}

		if strings.HasPrefix(token, "--watchlist=") {
			watchPath = strings.TrimPrefix(token, "--watchlist=")
			continue
		}

		if strings.HasPrefix(token, "-") {
			return "", "", nil, fmt.Errorf("unknown flag %q", token)
		}

		return configPath, watchPath, args[i:], nil
	}

	return configPath, watchPath, nil, nil
}

func isManagementCommand(name string) bool {
	switch name {
	case "add", "list", "install", "uninstall", "config", "completion", "help":
		return true
	default:
		return false
	}
}

func isInteractiveSession() bool {
	stdinInfo, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	stdoutInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}

	return (stdinInfo.Mode()&os.ModeCharDevice) != 0 && (stdoutInfo.Mode()&os.ModeCharDevice) != 0
}
