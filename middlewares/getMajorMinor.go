package middlewares

import (
	"fmt"
	"os"
	"syscall"
)

func MajorMinor(devicePath string) (int, int) {

	// Get file info
	fileInfo, err := os.Stat(devicePath)
	if err != nil {
		fmt.Println("error")
		return 0, 0

	}

	// Get system-specific file attributes
	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		fmt.Println("error")
		return 0, 0
	}

	// Extract major and minor numbers
	major := (stat.Rdev >> 8) & 0xff // Major number
	minor := stat.Rdev & 0xff        // Minor number

	return int(major), int(minor)
}
