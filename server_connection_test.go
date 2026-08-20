package main

import (
	"errors"
	"testing"
)

func TestPingServerSuccess(t *testing.T) {
	cfg := config{serverIP: "192.168.1.11"}

	createCommand := func(name string, args ...string) command {
		return fakeCommand{}
	}

	err := PingServer(cfg, createCommand)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestPingServerError(t *testing.T) {
	cfg := config{serverIP: "192.168.1.11"}
	expectedErr := errors.New("ping failed")

	createCommand := func(name string, args ...string) command {
		return fakeCommand{err: expectedErr}
	}

	err := PingServer(cfg, createCommand)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
}

func TestCheckSSHAvailableSuccess(t *testing.T) {
	cfg := config{serverIP: "192.168.1.11"}

	createdCommand := func(name string, args ...string) command {
		return fakeCommand{}
	}

	err := CheckSSHAvailable(cfg, createdCommand)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestCheckSSHAvailableError(t *testing.T) {
	cfg := config{serverIP: "192.168.1.11"}
	expectedError := errors.New("ssh fail")

	createCommand := func(name string, args ...string) command {
		return fakeCommand{err: expectedError}
	}

	err := CheckSSHAvailable(cfg, createCommand)

	if !errors.Is(err, expectedError) {
		t.Fatalf("expected %v, got %v", expectedError, err)
	}
}
