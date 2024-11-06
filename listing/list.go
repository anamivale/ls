package listing

import (
	"errors"
	"io/fs"
	"os"
	"strings"

	"github.com/anamivale/ls/middlewares"
	"github.com/anamivale/ls/options"
)

func GetDirContent(path string, flags options.Flags) ([]fs.DirEntry, error) {
	chacker, readErr := middlewares.InfoTDir(path)
	if readErr != nil {
		return nil, errors.New("go run .: cannot access " + path + ": No such file or directory")
	}
	if chacker != nil {
		return []fs.DirEntry{chacker}, nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, errors.New("go run .: cannot access " + path + ": No such file or directory")
	}
	var output []fs.DirEntry

	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".") {
			if !flags.All {
				continue
			}
		}
		output = append(output, entry)
	}
	return output, nil
}
