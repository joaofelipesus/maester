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
			cfg:           config{serverUserName: "admin", serverIP: "192.168.1.11", appPath: "/home/luwin/winterfel"},
			expectedError: false,
			expected:      nil,
		},
		{
			name:          "MissingServerUserName",
			cfg:           config{serverUserName: "", serverIP: "192.168.1.11", appPath: "/home/luwin/winterfel"},
			expectedError: true,
			expected:      errors.New("user tag is required"),
		},
		{
			name:          "MissingIP",
			cfg:           config{serverUserName: "admin", serverIP: "", appPath: "/home/luwin/winterfel"},
			expectedError: true,
			expected:      errors.New("IP address is required"),
		},
		{
			name:          "MissingAppPath",
			cfg:           config{serverUserName: "admin", serverIP: "192.168.10.11", appPath: ""},
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
