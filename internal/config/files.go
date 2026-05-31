package config

import (
	"os"
	"path/filepath"
)

var (
	DefaultConfigPath    = computeDefaultConfigPath()
	DefaultWatchlistPath = computeDefaultWatchlistPath()
)

// FileConfig holds file-path preferences read from config.
type FileConfig struct {
	WatchlistDirectory string `mapstructure:"watchlist_dir"`
}

func mulondaConfigDir() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		home, _ := os.UserHomeDir()
		dir = filepath.Join(home, ".config")
	}
	return filepath.Join(dir, "mulonda")
}

func computeDefaultConfigPath() string {
	return filepath.Join(mulondaConfigDir(), "config.yaml")
}

func computeDefaultWatchlistPath() string {
	return filepath.Join(mulondaConfigDir(), "watchlist.yaml")
}
