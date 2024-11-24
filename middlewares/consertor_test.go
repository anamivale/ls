package middlewares

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// Mock FileInfo for testing
type mockFileInfo struct {
	name    string
	size    int64
	mode    fs.FileMode
	modTime time.Time
	isDir   bool
	sys     interface{}
}

func (m mockFileInfo) Name() string       { return m.name }
func (m mockFileInfo) Size() int64        { return m.size }
func (m mockFileInfo) Mode() fs.FileMode  { return m.mode }
func (m mockFileInfo) ModTime() time.Time { return m.modTime }
func (m mockFileInfo) IsDir() bool        { return m.isDir }
func (m mockFileInfo) Sys() interface{}   { return m.sys }

// TestFileInfoDirEntry tests the FileInfoDirEntry struct methods
func TestFileInfoDirEntry(t *testing.T) {
	// Create test cases
	tests := []struct {
		name        string
		fileInfo    os.FileInfo
		wantName    string
		wantIsDir   bool
		wantType    fs.FileMode
		wantInfoErr bool
	}{
		{
			name: "Regular file",
			fileInfo: mockFileInfo{
				name:    "test.txt",
				size:    100,
				mode:    0644,
				modTime: time.Now(),
				isDir:   false,
			},
			wantName:    "test.txt",
			wantIsDir:   false,
			wantType:    0,
			wantInfoErr: false,
		},
		{
			name: "Directory",
			fileInfo: mockFileInfo{
				name:    "testdir",
				size:    4096,
				mode:    fs.ModeDir | 0755,
				modTime: time.Now(),
				isDir:   true,
			},
			wantName:    "testdir",
			wantIsDir:   true,
			wantType:    fs.ModeDir,
			wantInfoErr: false,
		},
		{
			name: "Symlink",
			fileInfo: mockFileInfo{
				name:    "testlink",
				size:    0,
				mode:    fs.ModeSymlink | 0777,
				modTime: time.Now(),
				isDir:   false,
			},
			wantName:    "testlink",
			wantIsDir:   false,
			wantType:    fs.ModeSymlink,
			wantInfoErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dirEntry := FileInfoDirEntry{info: tt.fileInfo}

			// Test Name method
			if got := dirEntry.Name(); got != tt.wantName {
				t.Errorf("FileInfoDirEntry.Name() = %v, want %v", got, tt.wantName)
			}

			// Test IsDir method
			if got := dirEntry.IsDir(); got != tt.wantIsDir {
				t.Errorf("FileInfoDirEntry.IsDir() = %v, want %v", got, tt.wantIsDir)
			}

			// Test Type method
			if got := dirEntry.Type(); got != tt.wantType {
				t.Errorf("FileInfoDirEntry.Type() = %v, want %v", got, tt.wantType)
			}

			// Test Info method
			info, err := dirEntry.Info()
			if (err != nil) != tt.wantInfoErr {
				t.Errorf("FileInfoDirEntry.Info() error = %v, wantErr %v", err, tt.wantInfoErr)
				return
			}
			if info != tt.fileInfo {
				t.Errorf("FileInfoDirEntry.Info() = %v, want %v", info, tt.fileInfo)
			}
		})
	}
}

// TestInfoTDir tests the InfoTDir function
func TestInfoTDir(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	// Create test files and directories
	regularFile := filepath.Join(tmpDir, "regular.txt")
	if err := os.WriteFile(regularFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	testDir := filepath.Join(tmpDir, "testdir")
	if err := os.Mkdir(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create test cases
	tests := []struct {
		name      string
		path      string
		wantErr   bool
		wantIsNil bool
	}{
		{
			name:      "Regular file",
			path:      regularFile,
			wantErr:   false,
			wantIsNil: false,
		},
		{
			name:      "Directory",
			path:      testDir,
			wantErr:   false,
			wantIsNil: true,
		},
		{
			name:      "Non-existent path",
			path:      filepath.Join(tmpDir, "nonexistent"),
			wantErr:   true,
			wantIsNil: true,
		},
		{
			name:      "Empty path",
			path:      "",
			wantErr:   true,
			wantIsNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InfoTDir(tt.path)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("InfoTDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check nil return
			if (got == nil) != tt.wantIsNil {
				t.Errorf("InfoTDir() returned %v, want nil: %v", got, tt.wantIsNil)
				return
			}

			// If we expect a non-nil result, verify it's the correct type
			if !tt.wantIsNil && got != nil {
				if _, ok := got.(FileInfoDirEntry); !ok {
					t.Errorf("InfoTDir() returned wrong type: got %T, want FileInfoDirEntry", got)
				}
			}

			// Verify error message format for non-existent files
			if tt.wantErr {
				expectedErr := "go run .: cannot access " + tt.path + ": No such file or directory"
				if err.Error() != expectedErr {
					t.Errorf("InfoTDir() error message = %v, want %v", err.Error(), expectedErr)
				}
			}
		})
	}
}

// TestInfoTDirSymlink tests the behavior with symbolic links
func TestInfoTDirSymlink(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a regular file and a symlink to it
	regularFile := filepath.Join(tmpDir, "regular.txt")
	if err := os.WriteFile(regularFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	symlink := filepath.Join(tmpDir, "symlink")
	if err := os.Symlink(regularFile, symlink); err != nil {
		t.Skipf("Skipping symlink test: %v", err)
		return
	}

	// Test symlink
	got, err := InfoTDir(symlink)
	if err != nil {
		t.Errorf("InfoTDir() error = %v, want nil", err)
		return
	}

	if got == nil {
		t.Error("InfoTDir() returned nil for symlink, want non-nil")
		return
	}

	// Verify the symlink is handled correctly
	info, err := got.Info()
	if err != nil {
		t.Errorf("got.Info() error = %v, want nil", err)
		return
	}

	if info.Mode()&os.ModeSymlink == 0 {
		t.Error("InfoTDir() didn't preserve symlink mode")
	}
}