package watchlist

// DefaultRules is the built-in baseline watchlist.
var DefaultRules = []Rule{
	{
		Command:      "rm",
		FlagsContain: []string{"-rf", "-r", "-f"},
		MinArgs:      2,
		Warning:      "Recursive delete",
	},
	{
		Command: "mv",
		Warning: "Move can overwrite files without recovery",
	},
	{
		Command: "cp",
		Warning: "Copy can overwrite existing files",
	},
}
