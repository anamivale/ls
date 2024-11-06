package middlewares

import (
	"errors"
	"io/fs"
	"os"
)

type FileInfoDirEntry struct {
	info os.FileInfo
}

func (f FileInfoDirEntry) Name() string {
	return f.info.Name()
}

func (f FileInfoDirEntry) IsDir() bool {
	return f.info.IsDir()
}

func (f FileInfoDirEntry) Type() fs.FileMode {
	return f.info.Mode().Type()
}

func (f FileInfoDirEntry) Info() (fs.FileInfo, error) {
	return f.info, nil
}

// Example usage
func InfoTDir(path string) (fs.DirEntry, error) {
	chacker, err1 := os.Lstat(path)

	if err1 != nil {
		err := errors.New("go run .: cannot access " + path + ": No such file or directory")
		// Handle the error
		return nil, err
	}

	// If it's not a directory, convert it to a DirEntry
	if !chacker.IsDir() {
		dirEntry := FileInfoDirEntry{info: chacker}
		return dirEntry, nil
		// Now you can use dirEntry as an fs.DirEntry
	}
	return nil, nil
}
