package listing

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
)

func Sort(entries []fs.DirEntry) {
	for i := 0; i < len(entries); i++ {
		for j := 0; j < len(entries); j++ {
			// Special cases for '.' and '..'
			if entries[i].Name() == "." && entries[j].Name() != "." {
				entries[i], entries[j] = entries[j], entries[i]
				continue
			}
			if entries[j].Name() == "." {
				continue
			}
			if entries[i].Name() == ".." && entries[j].Name() != ".." && entries[j].Name() != "." {
				entries[i], entries[j] = entries[j], entries[i]
				continue
			}
			if entries[j].Name() == ".." {
				continue
			}

			// General comparison logic
			if strings.ToLower(RemoveSpecialCharacters(entries[i].Name())) < strings.ToLower(RemoveSpecialCharacters(entries[j].Name())) {
				entries[i], entries[j] = entries[j], entries[i]
			}
		}
	}
}

func SortInReverse(entries []fs.DirEntry) {
	n := len(entries)
	for i := 0; i < n/2; i++ {
		entries[i], entries[n-1-i] = entries[n-1-i], entries[i]
	}
}

func TimeSort(entries []fs.DirEntry) {
	for i := 0; i < len(entries); i++ {
		for j := 0; j < len(entries); j++ {
			entriesI, _ := entries[i].Info()
			entriesJ, _ := entries[j].Info()
			if entriesI.ModTime().String() == entriesJ.ModTime().String() {
				if strings.ToLower(RemoveSpecialCharacters(entries[i].Name())) < strings.ToLower(RemoveSpecialCharacters(entries[j].Name())) {
					entries[i], entries[j] = entries[j], entries[i]
				}
				continue
			}

			if entriesI.ModTime().String() > entriesJ.ModTime().String() {
				entries[i], entries[j] = entries[j], entries[i]
			}
		}
	}
}

func TimeSortReverse(entries []fs.DirEntry) {
	for i := 0; i < len(entries); i++ {
		for j := 0; j < len(entries); j++ {
			entriesI, _ := entries[i].Info()
			entriesJ, _ := entries[j].Info()
			if entriesI.ModTime().String() == entriesJ.ModTime().String() {
				if strings.ToLower(RemoveSpecialCharacters(entries[i].Name())) > strings.ToLower(RemoveSpecialCharacters(entries[j].Name())) {
					entries[i], entries[j] = entries[j], entries[i]
				}
				continue
			}
			if entriesI.ModTime().String() < entriesJ.ModTime().String() {
				entries[i], entries[j] = entries[j], entries[i]
			}
		}
	}
}

func PathSort(paths []string) []string {
	var sortedPath []string
	var files []string
	var dirs []string

	for _, v := range paths {
		info, err := os.Lstat(v)
		if err != nil {
			fmt.Println("go run .: cannot access " + v + ": No such file or directory")
			os.Exit(0)
		}

		if info.IsDir() {
			dirs = append(dirs, v)
		} else {
			files = append(files, v)
		}

	}
	sortedPathNames(dirs)
	sortedPathNames(files)
	sortedPath = append(sortedPath, files...)
	sortedPath = append(sortedPath, dirs...)
	return sortedPath
}

func RemoveSpecialCharacters(s string) string {
	output := ""
	for _, v := range s {
		if v >= 'a' && v <= 'z' || v >= 'A' && v <= 'Z' || v >= '0' && v <= '9' {
			output += string(v)
		}
	}
	return output
}

func sortedPathNames(s []string) []string {
	for i := 0; i < len(s); i++ {
		for j := 0; j < len(s); j++ {
			if RemoveSpecialCharacters(s[i]) < RemoveSpecialCharacters(s[j]) {
				s[i], s[j] = s[j], s[i]
			}
		}
	}

	return s
}
