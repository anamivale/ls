package listing

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/anamivale/ls/formating"
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
	parentDir, _ := os.Lstat(path + "/..")
	curretDir, _ := os.Lstat(path + "/.")
	dotDir := fs.FileInfoToDirEntry(curretDir)
	dotDotDir := fs.FileInfoToDirEntry(parentDir)

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, errors.New("go run .: cannot access " + path + ": No such file or directory")
	}
	var output []fs.DirEntry
	if flags.All {
		output = append(output, dotDir, dotDotDir)
	}

	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".") {
			if !flags.All {
				continue
			}
		}
		output = append(output, entry)
	}
	Sort(output)

	if flags.Time {
		TimeSort(output)
		if flags.Reverse {
			TimeSortReverse(output)
		}
		return output, nil
	}
	if flags.Reverse {
		SortInReverse(output)
	}
	return output, nil
}

func GetDirContentRecursively(path string, flags options.Flags) error {
	entries, err := GetDirContent(path, flags)
	if err != nil {
		return errors.New("go run .: cannot access " + path + ": No such file or directory")
	}
	fmt.Println(path + ":")

	if flags.Long {
		width := formating.GetBlocks(path, entries)
		fmt.Printf("total %d\n", width.Blocks/2)
		formating.LongFormat(path, entries, flags)
	} else {
		formating.Format(entries)
	}
	fmt.Println()
	rError := errors.New("")
	for _, entry := range entries {
		// If the entry is a directory, recursively get its contents
		if entry.IsDir() {
			if entry.Name() == "." || entry.Name() == ".." {
				continue
			}
			subDirPath := middlewares.JoinPaths(path, entry.Name())

			err := GetDirContentRecursively(subDirPath, flags)
			rError = err
		}
	}

	return rError
}
