package cmd

import (
	"bytes"
	"testing"

	"github.com/cod3ddy/mulonda/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRootCommand_NoArgs_ShowsHelp(t *testing.T) {
	t.Cleanup(func() {
		rootCmd.SetOut(nil)
		rootCmd.SetErr(nil)
		rootCmd.SetArgs(nil)
	})

	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{})

	err := rootCmd.Execute()
	require.NoError(t, err)

	helpText := out.String()
	assert.Contains(t, helpText, "Usage:")
	assert.Contains(t, helpText, "Available Commands:")
}

func TestParseGlobalFlags(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		wantConfig    string
		wantWatchlist string
		wantPass      []string
		wantErr       string
	}{
		{
			name:    "config without value errors",
			args:    []string{"--config"},
			wantErr: "missing value for --config",
		},
		{
			name:    "watchlist without value errors",
			args:    []string{"--watchlist"},
			wantErr: "missing value for --watchlist",
		},
		{
			name:          "config equals form parsed",
			args:          []string{"--config=custom.yml"},
			wantConfig:    "custom.yml",
			wantWatchlist: config.DefaultWatchlistPath,
		},
		{
			name:          "watchlist equals form parsed",
			args:          []string{"--watchlist=custom-watch.yml"},
			wantConfig:    config.DefaultConfigPath,
			wantWatchlist: "custom-watch.yml",
		},
		{
			name:    "unknown flag errors",
			args:    []string{"--bad"},
			wantErr: `unknown flag "--bad"`,
		},
		{
			name:          "delimiter passes through remaining args",
			args:          []string{"--", "echo", "hello"},
			wantConfig:    config.DefaultConfigPath,
			wantWatchlist: config.DefaultWatchlistPath,
			wantPass:      []string{"echo", "hello"},
		},
		{
			name:          "non management command becomes passthrough",
			args:          []string{"echo", "hello"},
			wantConfig:    config.DefaultConfigPath,
			wantWatchlist: config.DefaultWatchlistPath,
			wantPass:      []string{"echo", "hello"},
		},
		{
			name:          "empty passthrough after only flags",
			args:          []string{"--config=foo.yml", "--watchlist=bar.yml"},
			wantConfig:    "foo.yml",
			wantWatchlist: "bar.yml",
		},
		{
			name:          "first non flag token selected after valid globals",
			args:          []string{"--config", "cfg.yml", "--watchlist", "wl.yml", "ls", "-la"},
			wantConfig:    "cfg.yml",
			wantWatchlist: "wl.yml",
			wantPass:      []string{"ls", "-la"},
		},
		{
			name:          "management command name passes through from parser",
			args:          []string{"list"},
			wantConfig:    config.DefaultConfigPath,
			wantWatchlist: config.DefaultWatchlistPath,
			wantPass:      []string{"list"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotConfig, gotWatchlist, gotPass, err := parseGlobalFlags(tt.args)

			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantConfig, gotConfig)
			assert.Equal(t, tt.wantWatchlist, gotWatchlist)
			assert.Equal(t, tt.wantPass, gotPass)
		})
	}
}

func TestIsManagementCommand(t *testing.T) {
	tests := []struct {
		name string
		cmd  string
		want bool
	}{
		{name: "add", cmd: "add", want: true},
		{name: "remove", cmd: "remove", want: true},
		{name: "list", cmd: "list", want: true},
		{name: "install", cmd: "install", want: true},
		{name: "uninstall", cmd: "uninstall", want: true},
		{name: "config", cmd: "config", want: true},
		{name: "completion", cmd: "completion", want: true},
		{name: "help", cmd: "help", want: true},
		{name: "non management command", cmd: "echo", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isManagementCommand(tt.cmd)
			assert.Equal(t, tt.want, got, "isManagementCommand(%q)", tt.cmd)
		})
	}
}
