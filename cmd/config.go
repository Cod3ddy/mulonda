package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cod3ddy/mulonda/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
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
	Short: "Set a config value (e.g. timeout_seconds=60, non_interactive.passthrough=false)",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := strings.TrimSpace(args[0])
		value := strings.TrimSpace(args[1])

		path := configFile
		if path == "" {
			path = config.DefaultConfigPath
		}

		if err := setConfigValue(path, key, value); err != nil {
			return err
		}
		fmt.Printf("set %s = %s\n", key, value)
		return nil
	},
}

func init() {
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
	rootCmd.AddCommand(configCmd)
}

func setConfigValue(path, key, value string) error {
	data := map[string]any{}
	if content, err := os.ReadFile(path); err == nil {
		_ = yaml.Unmarshal(content, &data)
	}

	setNestedKey(data, key, coerceValue(value))

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	out, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	return os.WriteFile(path, out, 0o644)
}

func setNestedKey(m map[string]interface{}, key string, value interface{}) {
	parts := strings.SplitN(key, ".", 2)
	if len(parts) == 1 {
		m[key] = value
		return
	}
	sub, ok := m[parts[0]].(map[string]any)
	if !ok {
		sub = map[string]any{}
	}
	setNestedKey(sub, parts[1], value)
	m[parts[0]] = sub
}

func coerceValue(s string) interface{} {
	switch strings.ToLower(s) {
	case "true":
		return true
	case "false":
		return false
	}
	if n, err := strconv.Atoi(s); err == nil {
		return n
	}
	return s
}
