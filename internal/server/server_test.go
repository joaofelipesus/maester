package server

import (
	"errors"
	"maester/internal"
	"testing"
)

type fakeCommand struct {
	output []byte
	err    error
}

func (command fakeCommand) CombinedOutput() ([]byte, error) {
	return command.output, command.err
}

func TestPingServerSuccess(t *testing.T) {
	cfg := internal.Config{ServerIP: "192.168.1.11"}

	createCommand := func(name string, args ...string) internal.Command {
		return fakeCommand{}
	}

	err := PingServer(cfg, createCommand)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestPingServerError(t *testing.T) {
	cfg := internal.Config{ServerIP: "192.168.1.11"}
	expectedErr := errors.New("ping failed")

	createCommand := func(name string, args ...string) internal.Command {
		return fakeCommand{err: expectedErr}
	}

	err := PingServer(cfg, createCommand)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
}

func TestCheckSSHAvailableSuccess(t *testing.T) {
	cfg := internal.Config{ServerIP: "192.168.1.11"}

	createdCommand := func(name string, args ...string) internal.Command {
		return fakeCommand{}
	}

	err := CheckSSHAvailable(cfg, createdCommand)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestCheckSSHAvailableError(t *testing.T) {
	cfg := internal.Config{ServerIP: "192.168.1.11"}
	expectedError := errors.New("ssh fail")

	createCommand := func(name string, args ...string) internal.Command {
		return fakeCommand{err: expectedError}
	}

	err := CheckSSHAvailable(cfg, createCommand)

	if !errors.Is(err, expectedError) {
		t.Fatalf("expected %v, got %v", expectedError, err)
	}
}
