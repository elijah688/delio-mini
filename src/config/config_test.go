package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {

	// set up test cases
	testCases := []struct {
		auth, name  string
		shouldError bool
	}{
		{
			name:        "auth is set",
			auth:        "token",
			shouldError: false,
		},

		{
			name:        "auth is not set",
			auth:        "",
			shouldError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tc.shouldError {
						t.Errorf("config.New() returned %s, expected to succeed", r)
					}
				}
			}()
			os.Setenv(FH_TOKEN, tc.auth)
			cfg := New()

			if !tc.shouldError {
				assert.Equal(t, cfg.DefaultHeader, map[string]string{"X-Finnhub-Token": tc.auth})
			}
		})
	}
}
