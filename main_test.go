package main

import (
	"bytes"
	"errors"
	"reflect"
	"testing"
)

func TestValidateRequiredConfigs(t *testing.T) {
	testCases := []struct {
		name          string
		cfg           Config
		expectedError bool
		expected      error
	}{
		{
			name:          "NoMissingParams",
			cfg:           Config{serverUserName: "admin", serverIP: "192.168.1.11", appPath: "/home/luwin/winterfel"},
			expectedError: false,
			expected:      nil,
		},
		{
			name:          "MissingServerUserName",
			cfg:           Config{serverUserName: "", serverIP: "192.168.1.11", appPath: "/home/luwin/winterfel"},
			expectedError: true,
			expected:      errors.New("user is required"),
		},
		{
			name:          "MissingIP",
			cfg:           Config{serverUserName: "admin", serverIP: "", appPath: "/home/luwin/winterfel"},
			expectedError: true,
			expected:      errors.New("IP address is required"),
		},
		{
			name:          "MissingAppPath",
			cfg:           Config{serverUserName: "admin", serverIP: "192.168.10.11", appPath: ""},
			expectedError: true,
			expected:      errors.New("App path is required"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := validateRequiredConfigs(testCase.cfg)

			if testCase.expectedError && result.Error() != testCase.expected.Error() {
				t.Errorf("Expected %q, got %q instead\n", testCase.expected, result)
			}
		})
	}
}

type fakeCommand struct {
	output []byte
	err    error
}

func (command fakeCommand) CombinedOutput() ([]byte, error) {
	return command.output, command.err
}

// records every command built by the factory, so the tests can assert which
// commands run and with which arguments.
type recordedCall struct {
	name string
	args []string
}

func TestRunSuccess(t *testing.T) {
	cfg := Config{
		serverUserName: "jon",
		serverIP:       "192.168.10.11",
		appPath:        "/user/castle-black",
		downloadLogs:   true,
	}

	expectedLogs := []byte("app | winter is coming\n")
	// store every call done to createCommand, which mocks the execution of commands,
	// so we are using the injection of an interface to "mock" the execution of commands
	var calls []recordedCall

	createCommand := func(name string, args ...string) command {
		calls = append(calls, recordedCall{name: name, args: args})

		if name == "ssh" {
			return fakeCommand{output: expectedLogs}
		}

		return fakeCommand{}
	}

	var mockWriter bytes.Buffer

	run(cfg, createCommand, &mockWriter)

	expectedCalls := []recordedCall{
		{name: "ping", args: []string{"-c", "1", "192.168.10.11"}},
		{name: "nc", args: []string{"-zv", "192.168.10.11", "22"}},
		{name: "ssh", args: []string{
			"jon@192.168.10.11",
			"cd", "/user/castle-black", "&&",
			"docker", "compose",
			"-f", "docker-compose.yml",
			"-f", "docker-compose.prod.yml",
			"logs", "--since", "2h", "app",
		}},
	}

	if !reflect.DeepEqual(calls, expectedCalls) {
		t.Fatalf("expected commands %v, got %v", expectedCalls, calls)
	}

	if mockWriter.String() != string(expectedLogs) {
		t.Errorf("expected %q, got %q", expectedLogs, mockWriter.String())
	}
}
