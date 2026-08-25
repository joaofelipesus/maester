package internal

import "os/exec"

// interface extracted to enable cover function pingServer and other functions that uses
// direct command calls.
type Command interface {
	CombinedOutput() ([]byte, error)
}

type CommandFactory func(name string, args ...string) Command

func RealCommand(name string, args ...string) Command {
	return exec.Command(name, args...)
}
