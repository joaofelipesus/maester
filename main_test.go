package main

import (
	"errors"
	"testing"
)

func TestValidateRequiredConfigs(t *testing.T) {
	testCases := []struct {
		name          string
		cfg           config
		expectedError bool
		expected      error
	}{
		{
			name:          "NoMissingParams",
			cfg:           config{serverUserName: "admin", serverIP: "192.168.1.11"},
			expectedError: false,
			expected:      nil,
		},
		{
			name:          "MissingServerUserName",
			cfg:           config{serverUserName: "", serverIP: "192.168.1.11"},
			expectedError: true,
			expected:      errors.New("user tag is required"),
		},
		{
			name:          "MissingIP",
			cfg:           config{serverUserName: "admin", serverIP: ""},
			expectedError: true,
			expected:      errors.New("IP address ir required"),
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

// TODO: move test to a specific module when move the function from main package
type fakeCommand struct {
	output []byte
	err    error
}

func (command fakeCommand) CombinedOutput() ([]byte, error) {
	return command.output, command.err
}

func TestPingServerSuccess(t *testing.T) {
	cfg := config{serverIP: "192.168.1.11"}

	createCommand := func(name string, args ...string) command {
		return fakeCommand{}
	}

	err := pingServer(cfg, createCommand)

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

	err := pingServer(cfg, createCommand)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
}
