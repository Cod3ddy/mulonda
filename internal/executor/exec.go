package executor

import (
	"context"
	"os/exec"
)

// Command returns an exec.Cmd bound to ctx — the process is killed when ctx is done.
func Command(ctx context.Context, name string, args ...string) *exec.Cmd {
	return exec.CommandContext(ctx, name, args...)
}
