package middlewares

import "os"

func JoinPaths(path, entryName string) string {
	fileinfo, err := os.Lstat(path)
	if err == nil {
		if !fileinfo.IsDir() {
			return path
		}
	}
	// Check if the base path already ends with a separator
	if path[len(path)-1] == os.PathSeparator {
		return path + entryName
	}
	return path + string(os.PathSeparator) + entryName
}
