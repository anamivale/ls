package formating

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/anamivale/ls/options"
)

// mockDirEntry implements fs.DirEntry for testing
type mockDirEntry struct {
	name  string
	isDir bool
	info  fs.FileInfo
}

func (m mockDirEntry) Name() string               { return m.name }
func (m mockDirEntry) IsDir() bool                { return m.isDir }
func (m mockDirEntry) Type() fs.FileMode          { return m.info.Mode().Type() }
func (m mockDirEntry) Info() (fs.FileInfo, error) { return m.info, nil }

// mockFileInfo implements fs.FileInfo for testing
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

// TestFormatTime tests the formatTime function
func TestFormatTime(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{
			name:     "Recent time",
			time:     now.AddDate(0, -1, 0),
			expected: now.AddDate(0, -1, 0).Format("Jan _2 15:04"),
		},
		{
			name:     "Old time",
			time:     now.AddDate(-1, 0, 0),
			expected: now.AddDate(-1, 0, 0).Format("Jan _2  2006"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatTime(tt.time)
			if result != tt.expected {
				t.Errorf("formatTime() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestGetBlocks tests the GetBlocks function
func TestGetBlocks(t *testing.T) {
	// Create mock stat_t for testing
	mockStat := &syscall.Stat_t{
		Nlink:  1,
		Uid:    1000,
		Gid:    1000,
		Size:   1024,
		Blocks: 2,
	}

	// Create mock entries
	entries := []fs.DirEntry{
		mockDirEntry{
			name: "file1.txt",
			info: mockFileInfo{
				name:    "file1.txt",
				size:    1024,
				mode:    0o644,
				modTime: time.Now(),
				sys:     mockStat,
			},
		},
		mockDirEntry{
			name: "file2.txt",
			info: mockFileInfo{
				name:    "file2.txt",
				size:    2048,
				mode:    0o644,
				modTime: time.Now(),
				sys:     mockStat,
			},
		},
	}

	width := GetBlocks(".", entries)

	// Test block count
	if width.Blocks != 4 { // 2 blocks per file
		t.Errorf("GetBlocks().Blocks = %v, want %v", width.Blocks, 4)
	}

	// Test width calculations
	if width.Permw == 0 || width.Linkw == 0 || width.Sizew == 0 {
		t.Error("GetBlocks() returned zero width for some fields")
	}
}

// TestFormat tests the Format function
func TestFormat(t *testing.T) {
	entries := []fs.DirEntry{
		mockDirEntry{name: "file1.txt"},
		mockDirEntry{name: "file2.txt"},
	}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Format(entries)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	expected := "file1.txt  file2.txt  \n"
	if output != expected {
		t.Errorf("Format() output = %q, want %q", output, expected)
	}
}

// TestLongFormat tests the LongFormat function
func TestLongFormat(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test file
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create a symlink
	testLink := filepath.Join(tmpDir, "test.link")
	if err := os.Symlink(testFile, testLink); err != nil {
		t.Logf("Skipping symlink test: %v", err)
	}

	// Get real directory entries
	dirEntries, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to read directory: %v", err)
	}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	flags := options.Flags{} // Add any necessary flags
	LongFormat(tmpDir, dirEntries, flags)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Check if output contains expected elements
	if !strings.Contains(output, "test.txt") {
		t.Error("LongFormat() output doesn't contain filename")
	}
	if !strings.Contains(output, "-rw-") {
		t.Error("LongFormat() output doesn't contain permissions")
	}
}

// TestWidthAndBlocks tests the WidthAndBlocks struct
func TestWidthAndBlocks(t *testing.T) {
	width := WidthAndBlocks{
		Blocks: 10,
		Permw:  10,
		Userrw: 8,
		Groupw: 8,
		Linkw:  3,
		Sizew:  5,
		Datew:  12,
		Namew:  20,
		Minor:  3,
		Major:  3,
	}

	// Test field values
	tests := []struct {
		name     string
		got      int
		expected int
	}{
		{"Blocks", width.Blocks, 10},
		{"Permw", width.Permw, 10},
		{"Userrw", width.Userrw, 8},
		{"Groupw", width.Groupw, 8},
		{"Linkw", width.Linkw, 3},
		{"Sizew", width.Sizew, 5},
		{"Datew", width.Datew, 12},
		{"Namew", width.Namew, 20},
		{"Minor", width.Minor, 3},
		{"Major", width.Major, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("WidthAndBlocks.%s = %v, want %v", tt.name, tt.got, tt.expected)
			}
		})
	}
}

// TestLongFormatEdgeCases tests edge cases for LongFormat

