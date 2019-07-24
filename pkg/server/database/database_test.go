package database

import (
	"github.com/dnote/dnote/pkg/assert"
	"testing"
)

func TestValidateConfig(t *testing.T) {
	testCases := []struct {
		input    Config
		expected error
	}{
		{
			input: Config{
				Host:     "mockHost",
				Port:     "mockPort",
				Name:     "mockName",
				User:     "mockUser",
				Password: "mockPassword",
			},
			expected: nil,
		},
		{
			input: Config{
				Host: "mockHost",
				Port: "mockPort",
				Name: "mockName",
				User: "mockUser",
			},
			expected: nil,
		},
		{
			input: Config{
				Port:     "mockPort",
				Name:     "mockName",
				User:     "mockUser",
				Password: "mockPassword",
			},
			expected: ErrConfigMissingHost,
		},
		{
			input: Config{
				Host:     "mockHost",
				Name:     "mockName",
				User:     "mockUser",
				Password: "mockPassword",
			},
			expected: ErrConfigMissingPort,
		},
		{
			input: Config{
				Host:     "mockHost",
				Port:     "mockPort",
				User:     "mockUser",
				Password: "mockPassword",
			},
			expected: ErrConfigMissingName,
		},
		{
			input: Config{
				Host:     "mockHost",
				Port:     "mockPort",
				Name:     "mockName",
				Password: "mockPassword",
			},
			expected: ErrConfigMissingUser,
		},
		{
			input:    Config{},
			expected: ErrConfigMissingHost,
		},
	}

	for _, tc := range testCases {
		result := validateConfig(tc.input)

		assert.Equal(t, result, tc.expected, "result mismatch")
	}
}
