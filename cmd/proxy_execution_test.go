package cmd

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"testing"

	"github.com/cod3ddy/mulonda/internal/config"
	"github.com/cod3ddy/mulonda/internal/watchlist"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func resetProxyHooks(t *testing.T) {
	origLoadConfig := loadConfigFn
	origLoadWatchlist := loadWatchlistFn
	origMatch := matchRuleFn
	origConfirm := confirmPromptFn
	origCommandFactory := commandFactoryFn
	origInteractive := isInteractiveFn

	t.Cleanup(func() {
		loadConfigFn = origLoadConfig
		loadWatchlistFn = origLoadWatchlist
		matchRuleFn = origMatch
		confirmPromptFn = origConfirm
		commandFactoryFn = origCommandFactory
		isInteractiveFn = origInteractive
	})
}

func TestExecuteProxy_LoadConfigError(t *testing.T) {
	resetProxyHooks(t)

	loadConfigFn = func(path string) (config.Config, error) {
		return config.Config{}, errors.New("config boom")
	}

	err := executeProxy("ignored", "ignored", []string{"echo", "hi"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "load config: config boom")
}

func TestExecuteProxy_LoadWatchlistError(t *testing.T) {
	resetProxyHooks(t)

	loadConfigFn = func(path string) (config.Config, error) {
		return config.Default(), nil
	}
	loadWatchlistFn = func(path string) ([]watchlist.Rule, error) {
		return nil, errors.New("watchlist boom")
	}

	err := executeProxy("ignored", "ignored", []string{"echo", "hi"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "load watchlist: watchlist boom")
}

func TestExecuteProxy_MatchedInteractive_ConfirmYesExecutes(t *testing.T) {
	resetProxyHooks(t)

	loadConfigFn = func(path string) (config.Config, error) {
		return config.Default(), nil
	}
	loadWatchlistFn = func(path string) ([]watchlist.Rule, error) {
		return []watchlist.Rule{{Command: "danger", Warning: "be careful"}}, nil
	}
	matchRuleFn = func(command string, args []string, rules []watchlist.Rule) (watchlist.Rule, bool) {
		return watchlist.Rule{Command: "danger", Warning: "be careful"}, true
	}
	isInteractiveFn = func() bool { return true }
	confirmPromptFn = func(message string) bool { return true }
	commandFactoryFn = func(name string, args ...string) *exec.Cmd {
		return exec.Command("bash", "-lc", "true")
	}

	err := executeProxy("ignored", "ignored", []string{"danger", "file.txt"})
	require.NoError(t, err)
}

func TestExecuteProxy_MatchedInteractive_ConfirmNoAborts(t *testing.T) {
	resetProxyHooks(t)

	loadConfigFn = func(path string) (config.Config, error) {
		return config.Default(), nil
	}
	loadWatchlistFn = func(path string) ([]watchlist.Rule, error) {
		return []watchlist.Rule{{Command: "danger", Warning: "be careful"}}, nil
	}
	matchRuleFn = func(command string, args []string, rules []watchlist.Rule) (watchlist.Rule, bool) {
		return watchlist.Rule{Command: "danger", Warning: "be careful"}, true
	}
	isInteractiveFn = func() bool { return true }
	confirmPromptFn = func(message string) bool { return false }

	commandCalled := false
	commandFactoryFn = func(name string, args ...string) *exec.Cmd {
		commandCalled = true
		return exec.Command("bash", "-lc", "true")
	}

	err := executeProxy("ignored", "ignored", []string{"danger", "file.txt"})
	require.EqualError(t, err, "aborted: command not executed")
	assert.False(t, commandCalled)
}

func TestExecuteProxy_MatchedNonInteractive_PassthroughTrueExecutes(t *testing.T) {
	resetProxyHooks(t)

	loadConfigFn = func(path string) (config.Config, error) {
		cfg := config.Default()
		cfg.NonInteractive.Passthrough = true
		return cfg, nil
	}
	loadWatchlistFn = func(path string) ([]watchlist.Rule, error) {
		return []watchlist.Rule{{Command: "danger"}}, nil
	}
	matchRuleFn = func(command string, args []string, rules []watchlist.Rule) (watchlist.Rule, bool) {
		return watchlist.Rule{Command: "danger"}, true
	}
	isInteractiveFn = func() bool { return false }
	commandFactoryFn = func(name string, args ...string) *exec.Cmd {
		return exec.Command("bash", "-lc", "true")
	}

	err := executeProxy("ignored", "ignored", []string{"danger"})
	require.NoError(t, err)
}

func TestExecuteProxy_MatchedNonInteractive_PassthroughFalseBlocks(t *testing.T) {
	resetProxyHooks(t)

	loadConfigFn = func(path string) (config.Config, error) {
		cfg := config.Default()
		cfg.NonInteractive.Passthrough = false
		return cfg, nil
	}
	loadWatchlistFn = func(path string) ([]watchlist.Rule, error) {
		return []watchlist.Rule{{Command: "danger"}}, nil
	}
	matchRuleFn = func(command string, args []string, rules []watchlist.Rule) (watchlist.Rule, bool) {
		return watchlist.Rule{Command: "danger"}, true
	}
	isInteractiveFn = func() bool { return false }

	err := executeProxy("ignored", "ignored", []string{"danger"})
	require.EqualError(t, err, "blocked in non-interactive mode: danger")
}

func TestExecuteProxy_UnmatchedCommandExecutesWithoutPrompt(t *testing.T) {
	resetProxyHooks(t)

	loadConfigFn = func(path string) (config.Config, error) {
		return config.Default(), nil
	}
	loadWatchlistFn = func(path string) ([]watchlist.Rule, error) {
		return []watchlist.Rule{{Command: "danger"}}, nil
	}
	matchRuleFn = func(command string, args []string, rules []watchlist.Rule) (watchlist.Rule, bool) {
		return watchlist.Rule{}, false
	}

	promptCalled := false
	confirmPromptFn = func(message string) bool {
		promptCalled = true
		return true
	}

	commandFactoryFn = func(name string, args ...string) *exec.Cmd {
		return exec.Command("bash", "-lc", "true")
	}

	err := executeProxy("ignored", "ignored", []string{"echo", "hi"})
	require.NoError(t, err)
	assert.False(t, promptCalled)
}

func TestExecuteProxy_BlankWarningUsesFallbackAndFormatsZeroArgsPrompt(t *testing.T) {
	resetProxyHooks(t)

	loadConfigFn = func(path string) (config.Config, error) {
		return config.Default(), nil
	}
	loadWatchlistFn = func(path string) ([]watchlist.Rule, error) {
		return []watchlist.Rule{{Command: "danger", Warning: "   "}}, nil
	}
	matchRuleFn = func(command string, args []string, rules []watchlist.Rule) (watchlist.Rule, bool) {
		return watchlist.Rule{Command: "danger", Warning: "   "}, true
	}
	isInteractiveFn = func() bool { return true }

	var gotPrompt string
	confirmPromptFn = func(message string) bool {
		gotPrompt = message
		return true
	}

	commandFactoryFn = func(name string, args ...string) *exec.Cmd {
		return exec.Command("bash", "-lc", "true")
	}

	err := executeProxy("ignored", "ignored", []string{"danger"})
	require.NoError(t, err)
	assert.Contains(t, gotPrompt, "warning: Watched command requires confirmation")
	assert.Contains(t, gotPrompt, "mulonda: danger\nwarning:")
	assert.NotContains(t, gotPrompt, "danger \nwarning:")
}

func TestExecuteProxy_CommandFailureIsWrapped(t *testing.T) {
	resetProxyHooks(t)

	loadConfigFn = func(path string) (config.Config, error) {
		return config.Default(), nil
	}
	loadWatchlistFn = func(path string) ([]watchlist.Rule, error) {
		return []watchlist.Rule{}, nil
	}
	matchRuleFn = func(command string, args []string, rules []watchlist.Rule) (watchlist.Rule, bool) {
		return watchlist.Rule{}, false
	}
	commandFactoryFn = func(name string, args ...string) *exec.Cmd {
		return exec.Command("bash", "-lc", "exit 7")
	}

	err := executeProxy("ignored", "ignored", []string{"echo", "hi"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "execute command:")
}

func TestExecuteProxy_ForwardsStdIOToCommand(t *testing.T) {
	resetProxyHooks(t)

	loadConfigFn = func(path string) (config.Config, error) {
		return config.Default(), nil
	}
	loadWatchlistFn = func(path string) ([]watchlist.Rule, error) {
		return []watchlist.Rule{}, nil
	}
	matchRuleFn = func(command string, args []string, rules []watchlist.Rule) (watchlist.Rule, bool) {
		return watchlist.Rule{}, false
	}
	commandFactoryFn = func(name string, args ...string) *exec.Cmd {
		return exec.Command("bash", "-lc", "cat")
	}

	origStdin := os.Stdin
	origStdout := os.Stdout
	origStderr := os.Stderr
	t.Cleanup(func() {
		os.Stdin = origStdin
		os.Stdout = origStdout
		os.Stderr = origStderr
	})

	inFile, err := os.CreateTemp("", "mulonda-stdin-*")
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = os.Remove(inFile.Name())
	})

	_, err = inFile.WriteString("proxy-stdio\n")
	require.NoError(t, err)
	_, err = inFile.Seek(0, io.SeekStart)
	require.NoError(t, err)
	os.Stdin = inFile

	outRead, outWrite, err := os.Pipe()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = outRead.Close()
		_ = outWrite.Close()
	})
	os.Stdout = outWrite

	err = executeProxy("ignored", "ignored", []string{"cat"})
	require.NoError(t, err)

	require.NoError(t, outWrite.Close())
	outBytes, err := io.ReadAll(outRead)
	require.NoError(t, err)

	assert.Equal(t, "proxy-stdio\n", string(outBytes))
}
