package executor

import "os/exec"

// Command proxies execution to the underlying system binary.
func Command(name string, args ...string) *exec.Cmd {
	return exec.Command(name, args...)
}
