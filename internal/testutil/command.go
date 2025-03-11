package testutil

import (
	"os/exec"
)

// Command represents a command to be executed
type Command struct {
	*exec.Cmd
}

// NewCommand creates a new command with the given name and arguments
func NewCommand(name string, args ...string) *Command {
	return &Command{
		Cmd: exec.Command(name, args...),
	}
}

// Run executes the command and returns an error if it fails
func (c *Command) Run() error {
	return c.Cmd.Run()
}

// Output executes the command and returns its output
func (c *Command) Output() ([]byte, error) {
	return c.Cmd.Output()
}

// CombinedOutput executes the command and returns its combined output
func (c *Command) CombinedOutput() ([]byte, error) {
	return c.Cmd.CombinedOutput()
}
