package listing

import (
	"io/fs"
	"strings"
)

func Sort(entries []fs.DirEntry) {
	for i := 0; i < len(entries); i++ {
		for j := 0; j < len(entries); j++ {
			if strings.ToLower(strings.TrimPrefix(entries[i].Name(), ".")) < strings.ToLower(strings.TrimPrefix(entries[j].Name(), ".")) {
				entries[i], entries[j] = entries[j], entries[i]
			}
		}
	}
}

func SortInReverse(entries []fs.DirEntry) {
	for i := 0; i < len(entries); i++ {
		for j := 0; j < len(entries); j++ {
			if strings.ToLower(strings.TrimPrefix(entries[i].Name(), ".")) > strings.ToLower(strings.TrimPrefix(entries[j].Name(), ".")) {
				entries[i], entries[j] = entries[j], entries[i]
			}
		}
	}
}

func TimeSort(entries []fs.DirEntry) {
	for i := 0; i < len(entries); i++ {
		for j := 0; j < len(entries); j++ {
			entriesI, _ := entries[i].Info()
			entriesJ, _ := entries[j].Info()
			if entriesI.ModTime().String() == entriesJ.ModTime().String() {
				if strings.ToLower(strings.TrimPrefix(entries[i].Name(), ".")) < strings.ToLower(strings.TrimPrefix(entries[j].Name(), ".")) {
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
				if strings.ToLower(strings.TrimPrefix(entries[i].Name(), ".")) > strings.ToLower(strings.TrimPrefix(entries[j].Name(), ".")) {
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
