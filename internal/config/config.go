package config

import (
	"errors"
	"os"

	"github.com/spf13/viper"
)

// Config stores user preferences for Mulonda behavior.
type Config struct {
	TimeoutSeconds int                  `mapstructure:"timeout_seconds"`
	NonInteractive NonInteractiveConfig `mapstructure:"non_interactive"`
	Files          *FileConfig          `mapstructure:"files"`
}

type NonInteractiveConfig struct {
	Passthrough bool `mapstructure:"passthrough"`
}

// Default returns sane defaults for first run.
func Default() Config {
	cfg := Config{TimeoutSeconds: 30}
	cfg.NonInteractive.Passthrough = true
	cfg.Files = &FileConfig{WatchlistDirectory: DefaultWatchlistPath}
	return cfg
}

// Load reads YAML config using viper. Missing files return defaults.
func Load(path string) (Config, error) {
	cfg := Default()

	v := viper.New()
	v.SetConfigType("yaml")
	v.SetDefault("timeout_seconds", cfg.TimeoutSeconds)
	v.SetDefault("non_interactive.passthrough", cfg.NonInteractive.Passthrough)

	if path != "" {
		v.SetConfigFile(path)
	} else {
		v.SetConfigName("mulonda")
		v.AddConfigPath(".")
	}

	if err := v.ReadInConfig(); err != nil {
		if !isConfigMissing(err) {
			return cfg, err
		}
	} else {
		if err := v.Unmarshal(&cfg); err != nil {
			return cfg, err
		}
	}

	if cfg.Files == nil {
		cfg.Files = &FileConfig{WatchlistDirectory: DefaultWatchlistPath}
	}

	return cfg, nil
}

func isConfigMissing(err error) bool {
	var notFound viper.ConfigFileNotFoundError
	return errors.As(err, &notFound) || errors.Is(err, os.ErrNotExist)
}
