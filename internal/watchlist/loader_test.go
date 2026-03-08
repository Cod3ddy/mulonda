package watchlist

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddRule_CreatesParentDirAndPersistsStructuredFields(t *testing.T) {
	path := filepath.Join(t.TempDir(), "data", "watchlist.yaml")

	rule := Rule{
		Command:      "chmod",
		FlagsContain: []string{"-R", "-f"},
		ArgsMatch:    []string{"777", "a+rwx"},
		MinArgs:      2,
		Warning:      "Setting permissive mode",
	}

	require.NoError(t, AddRule(path, rule))

	rules, err := Load(path)
	require.NoError(t, err)

	var found *Rule
	for i := range rules {
		if rules[i].Command == "chmod" && rules[i].Warning == "Setting permissive mode" {
			found = &rules[i]
			break
		}
	}

	require.NotNil(t, found)
	assert.Equal(t, []string{"-R", "-f"}, found.FlagsContain)
	assert.Equal(t, []string{"777", "a+rwx"}, found.ArgsMatch)
	assert.Equal(t, 2, found.MinArgs)
}
