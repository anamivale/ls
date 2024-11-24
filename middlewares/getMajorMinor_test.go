package middlewares

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMajorMinor(t *testing.T) {
	// Define test cases
	tests := []struct {
		name          string
		devicePath    string
		expectedMajor int
		expectedMinor int
		shouldError   bool
	}{
		{
			name:          "Test /dev/null",
			devicePath:    "/dev/null",
			expectedMajor: 1,
			expectedMinor: 3,
			shouldError:   false,
		},
		{
			name:          "Test /dev/zero",
			devicePath:    "/dev/zero",
			expectedMajor: 1,
			expectedMinor: 5,
			shouldError:   false,
		},
		{
			name:          "Test non-existent device",
			devicePath:    "/dev/nonexistent",
			expectedMajor: 0,
			expectedMinor: 0,
			shouldError:   true,
		},
		{
			name:          "Test regular file",
			devicePath:    "testfile.txt",
			expectedMajor: 0,
			expectedMinor: 0,
			shouldError:   true,
		},
	}

	// Create a temporary regular file for testing
	tmpfile := filepath.Join(t.TempDir(), "testfile.txt")
	if err := os.WriteFile(tmpfile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// If we're testing the regular file case, use our temp file
			testPath := tt.devicePath
			if tt.devicePath == "testfile.txt" {
				testPath = tmpfile
			}

			major, minor := MajorMinor(testPath)

			// For device files that should exist, verify the numbers
			if !tt.shouldError && (major != tt.expectedMajor || minor != tt.expectedMinor) {
				t.Errorf("MajorMinor(%s) = (%d, %d), want (%d, %d)",
					tt.devicePath, major, minor, tt.expectedMajor, tt.expectedMinor)
			}

			// For cases that should error, verify we get (0, 0)
			if tt.shouldError && (major != 0 || minor != 0) {
				t.Errorf("MajorMinor(%s) = (%d, %d), want (0, 0) for error case",
					tt.devicePath, major, minor)
			}
		})
	}
}

// TestMajorMinorPermissions tests the behavior when the user doesn't have
// sufficient permissions to access the device file
func TestMajorMinorPermissions(t *testing.T) {
	// Skip if running as root since root can access everything
	if os.Geteuid() == 0 {
		t.Skip("Skipping permission test when running as root")
	}

	// Create a device file with restricted permissions
	tmpDir := t.TempDir()
	restrictedFile := filepath.Join(tmpDir, "restricted")
	if err := os.WriteFile(restrictedFile, []byte{}, 0000); err != nil {
		t.Fatalf("Failed to create restricted file: %v", err)
	}

	major, minor := MajorMinor(restrictedFile)
	if major != 0 || minor != 0 {
		t.Errorf("MajorMinor(%s) = (%d, %d), want (0, 0) for permission denied",
			restrictedFile, major, minor)
	}
}

// TestMajorMinorEdgeCases tests edge cases like empty paths and special characters
func TestMajorMinorEdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		devicePath string
	}{
		{"Empty path", ""},
		{"Space only path", " "},
		{"Path with special characters", "!@#$%^&*()"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			major, minor := MajorMinor(tt.devicePath)
			if major != 0 || minor != 0 {
				t.Errorf("MajorMinor(%s) = (%d, %d), want (0, 0) for invalid path",
					tt.devicePath, major, minor)
			}
		})
	}
}