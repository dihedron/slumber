package configuration

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/dihedron/slumberd/pointer"
	"github.com/dihedron/slumberd/timex"
)

func TestConfigurationCheck(t *testing.T) {
	// Create a temporary file for valid package test
	tmpFile, err := os.CreateTemp("", "packages-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Create a temporary file for unreadable package test
	unreadableFile, err := os.CreateTemp("", "unreadable-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(unreadableFile.Name())
	unreadableFile.Chmod(0000)
	unreadableFile.Close()

	tests := []struct {
		name          string
		config        Configuration
		value         string
		expectedError string
	}{
		{
			name: "Valid file",
			config: Configuration{
				Packages:  pointer.To(tmpFile.Name()),
				Debounce:  pointer.To(timex.Duration(500 * time.Millisecond)),
				Timeout:   pointer.To(timex.Duration(10 * time.Minute)),
				Frequency: pointer.To(timex.Duration(5 * time.Minute)),
			},
			expectedError: "",
		},
		{
			name: "Non-existent file",
			config: Configuration{
				Packages:  pointer.To("/path/to/non/existent/file"),
				Debounce:  pointer.To(timex.Duration(500 * time.Millisecond)),
				Timeout:   pointer.To(timex.Duration(10 * time.Minute)),
				Frequency: pointer.To(timex.Duration(5 * time.Minute)),
			},
			expectedError: "packages file /path/to/non/existent/file does not exist",
		},
		{
			name: "Unreadable file",
			config: Configuration{
				Packages:  pointer.To(unreadableFile.Name()),
				Debounce:  pointer.To(timex.Duration(500 * time.Millisecond)),
				Timeout:   pointer.To(timex.Duration(10 * time.Minute)),
				Frequency: pointer.To(timex.Duration(5 * time.Minute)),
			},
			expectedError: "packages file " + unreadableFile.Name() + " is not readable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// we need to create a config file that points to the packages
			// file we want to test

			configFile, err := os.CreateTemp("", "config-*.yaml")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(configFile.Name())

			configContent := "packages: " + *tt.config.Packages + "\n"
			if _, err := configFile.WriteString(configContent); err != nil {
				t.Fatal(err)
			}
			configFile.Close()

			c := &Configuration{}
			err = c.UnmarshalFlag(configFile.Name())

			if tt.expectedError == "" {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError && !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("expected error containing %q, got %v", tt.expectedError, err)
				}
			}
		})
	}
}
