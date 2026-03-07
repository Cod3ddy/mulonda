package watchlist

import (
	"errors"
	"slices"
	"strings"

	"github.com/spf13/viper"
)

// Rule describes a watchlist rule with optional matching metadata.
type Rule struct {
	Command      string   `mapstructure:"command" yaml:"command"`
	FlagsContain []string `mapstructure:"flags_contain" yaml:"flags_contain,omitempty"`
	ArgsMatch    []string `mapstructure:"args_match" yaml:"args_match,omitempty"`
	MinArgs      int      `mapstructure:"min_args" yaml:"min_args,omitempty"`
	Warning      string   `mapstructure:"warning" yaml:"warning,omitempty"`
}

type fileRules struct {
	Rules    []Rule   `mapstructure:"rules"`
	Commands []string `mapstructure:"commands"`
}

// Load returns a merged list of built-in and user-defined rules.
func Load(path string) ([]Rule, error) {
	rules := make([]Rule, 0, len(DefaultRules))
	rules = append(rules, DefaultRules...)

	v := viper.New()
	v.SetConfigType("yaml")
	if path == "" {
		path = "watchlist.yaml"
	}
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if errors.As(err, &notFound) {
			return rules, nil
		}
		return nil, err
	}

	var f fileRules
	if err := v.Unmarshal(&f); err != nil {
		return nil, err
	}

	for _, rule := range f.Rules {
		rule = normalizeRule(rule)
		if rule.Command == "" {
			continue
		}
		rules = upsertRule(rules, rule)
	}

	// Legacy flat list support.
	for _, cmd := range f.Commands {
		rule := parseInputToRule(cmd)
		if rule.Command == "" {
			continue
		}
		rules = upsertRule(rules, rule)
	}

	return rules, nil
}

// Add appends a command rule to the YAML watchlist file if missing.
func Add(path, command string) error {
	if path == "" {
		path = "watchlist.yaml"
	}

	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile(path)

	f := fileRules{}
	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			return err
		}
	} else {
		if err := v.Unmarshal(&f); err != nil {
			return err
		}
	}

	// Migrate legacy commands list into structured rules.
	for _, legacy := range f.Commands {
		rule := parseInputToRule(legacy)
		if rule.Command == "" {
			continue
		}
		f.Rules = upsertRule(f.Rules, rule)
	}

	f.Commands = nil

	newRule := parseInputToRule(command)
	if newRule.Command == "" {
		return nil
	}
	f.Rules = upsertRule(f.Rules, newRule)

	v.Set("rules", f.Rules)
	v.Set("commands", nil)

	if err := v.WriteConfigAs(path); err != nil {
		var alreadyExists viper.ConfigFileAlreadyExistsError
		if errors.As(err, &alreadyExists) {
			return v.WriteConfig()
		}
		return err
	}

	return nil
}

func containsRule(rules []Rule, target Rule) bool {
	for _, r := range rules {
		if equalRule(r, target) {
			return true
		}
	}
	return false
}

func equalRule(a, b Rule) bool {
	return equalIdentity(a, b) && a.Warning == b.Warning
}

func equalIdentity(a, b Rule) bool {
	return a.Command == b.Command &&
		a.MinArgs == b.MinArgs &&
		slices.Equal(a.FlagsContain, b.FlagsContain) &&
		slices.Equal(a.ArgsMatch, b.ArgsMatch)
}

func upsertRule(rules []Rule, incoming Rule) []Rule {
	for i, existing := range rules {
		if equalIdentity(existing, incoming) {
			rules[i] = incoming
			return rules
		}
	}
	return append(rules, incoming)
}

func parseInputToRule(input string) Rule {
	parts := strings.Fields(strings.TrimSpace(input))
	if len(parts) == 0 {
		return Rule{}
	}

	rule := Rule{Command: parts[0]}
	if len(parts) > 1 {
		rule.ArgsMatch = parts[1:]
	}
	return normalizeRule(rule)
}

func normalizeRule(rule Rule) Rule {
	rule.Command = strings.TrimSpace(rule.Command)
	rule.Warning = strings.TrimSpace(rule.Warning)
	return rule
}
