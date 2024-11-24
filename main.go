package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/anamivale/ls/formating"
	"github.com/anamivale/ls/listing"

	"github.com/anamivale/ls/options"
)

func main() {
	// Initialize path and flags arrays
	var paths []string
	var flagsArray []string

	// Get args and sort them into flags and paths
	args := os.Args
	for _, arg := range args[1:] {
		if arg == "-" {
			paths = append(paths, arg)
			continue
		}
		if arg == "-" {
			paths = append(paths, ".")
		} else if arg == "--l" {
			continue
		} else if arg[0] == '-' {
			if strings.HasPrefix(arg, "--") {
				fmt.Printf("go run .: unrecognized option %s", arg)
				return
			}
			flagsArray = append(flagsArray, arg)
		} else {
			paths = append(paths, arg)
		}
	}

	// If paths not given let use the current dir
	if paths == nil {
		paths = append(paths, ".")
	}
	// convert the flagsArray to a string
	flags := strings.Join(flagsArray, "")
	flags = strings.ReplaceAll(flags, "-", "")

	// Optios function gets the flags and sets them to true or false
	flagsStruct := options.Options(flags)
	paths = listing.PathSort(paths)

	for i, path := range paths {
		info, err := os.Lstat(path)
		if err != nil {
			fmt.Println("go run .: cannot access " + path + ": No such file or directory")
		}
		if len(paths) != 1 {
			if info.IsDir() {
				fmt.Println(path + ":")
			}
		}
		entries, err := listing.GetDirContent(path, flagsStruct)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if flagsStruct.Recursive {
			err := listing.GetDirContentRecursively(path, flagsStruct)
			if err != nil {
				return
			}
			return
		}
		if flagsStruct.Long {
			// check if file

			if info.IsDir() {
				width := formating.GetBlocks(path, entries)
				fmt.Printf("total %d\n", width.Blocks/2)

				formating.LongFormat(path, entries, flagsStruct)

			} else {
				formating.LongFormat(path, entries, flagsStruct)
			}
			// print blocks
		} else {
			formating.Format(entries)
		}
		if i != len(paths)-1 {
			if info.IsDir() {
				fmt.Println()
			}
		}
	}
}
