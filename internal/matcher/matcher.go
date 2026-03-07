package matcher

import (
	"strings"

	"github.com/cod3ddy/mulonda/internal/watchlist"
)

// IsWatchedCommand reports whether command should trigger confirmation.
func IsWatchedCommand(command string, rules []string) bool {
	for _, rule := range rules {
		if strings.TrimSpace(rule) == strings.TrimSpace(command) {
			return true
		}
	}
	return false
}

// MatchRule returns the first watchlist rule that matches command and args.
func MatchRule(command string, args []string, rules []watchlist.Rule) (watchlist.Rule, bool) {
	cmd := strings.TrimSpace(command)
	if cmd == "" {
		return watchlist.Rule{}, false
	}

	for _, rule := range rules {
		if !ruleMatches(cmd, args, rule) {
			continue
		}
		return rule, true
	}

	return watchlist.Rule{}, false
}

func ruleMatches(command string, args []string, rule watchlist.Rule) bool {
	if strings.TrimSpace(rule.Command) != command {
		return false
	}

	if rule.MinArgs > 0 && len(args) < rule.MinArgs {
		return false
	}

	if len(rule.FlagsContain) > 0 && !containsAnyArg(args, rule.FlagsContain) {
		return false
	}

	if len(rule.ArgsMatch) > 0 && !containsAnyArg(args, rule.ArgsMatch) {
		return false
	}

	return true
}

func containsAnyArg(args []string, candidates []string) bool {
	for _, arg := range args {
		a := strings.TrimSpace(arg)
		for _, candidate := range candidates {
			if a == strings.TrimSpace(candidate) {
				return true
			}
		}
	}
	return false
}
