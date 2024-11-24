package formating

import (
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"strconv"
	"syscall"
	"time"

	"github.com/anamivale/ls/middlewares"
	"github.com/anamivale/ls/options"
)

type WidthAndBlocks struct {
	Blocks int
	Permw  int
	Userrw int
	Groupw int
	Linkw  int
	Sizew  int
	Datew  int
	Namew  int
	Minor  int
	Major  int
}

func LongFormat(path string, entries []fs.DirEntry, flags options.Flags) {
	width := GetBlocks(path, entries)

	// Print total blocks if more than one entry

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		path1 := middlewares.JoinPaths(path, entry.Name())

		// Get file metadata
		perm := middlewares.ModeToString(info.Mode())
		link := info.Sys().(*syscall.Stat_t).Nlink
		size := strconv.FormatInt(info.Size(), 10)
		date := formatTime(info.ModTime())
		name := entry.Name()

		// Get user and group info
		uid := fmt.Sprint(info.Sys().(*syscall.Stat_t).Uid)
		gid := fmt.Sprint(info.Sys().(*syscall.Stat_t).Gid)

		userr := uid // Fallback to UID if username lookup fails
		group := gid // Fallback to GID if group lookup fails

		if Userr, err := user.LookupId(uid); err == nil {
			userr = Userr.Username
		}
		if Group, err := user.LookupGroupId(gid); err == nil {
			group = Group.Name
		}


		// Handle different file types
		switch {
		case info.Mode()&os.ModeCharDevice != 0 || info.Mode()&os.ModeDevice != 0:
			// Device files: show major,minor instead of size
			major, minor := middlewares.MajorMinor(path1)
			fmt.Printf("%s %*d %-*s %-*s %*d, %*d %s %s\n",
				perm, width.Linkw, link, width.Userrw, userr, width.Groupw, group,
				width.Major, major, width.Minor, minor, date, name)

		case info.Mode()&os.ModeSymlink != 0:
			// Symbolic links: show target
			target, err := os.Readlink(path1)
			if err == nil {
				fmt.Printf("%s %*d %-*s %-*s %*s %s %s -> %s\n",
					perm, width.Linkw, link, width.Userrw, userr, width.Groupw, group,
					width.Sizew, size, date, name, target)
			}

		default:
			// Regular files and directories
			fmt.Printf("%s %*d %-*s %-*s %*s %s %s\n",
				perm, width.Linkw, link, width.Userrw, userr, width.Groupw, group,
				width.Sizew, size, date, name)
		}
	}
}

func Format(entries []fs.DirEntry) {
	for _, entry := range entries {
		fmt.Printf("%s  ", entry.Name())
	}
	fmt.Println()
}

func GetBlocks(path string, entries []fs.DirEntry) WidthAndBlocks {
	var width WidthAndBlocks
	if len(entries) < 1 {
		return width
	}
	blocks := 0
	info, _ := entries[0].Info()
	perm := len(middlewares.ModeToString(info.Mode()))

	link := len(strconv.Itoa(int(info.Sys().(*syscall.Stat_t).Nlink)))
	size := len(strconv.Itoa(int(info.Size())))
	date := len(info.ModTime().Format("Jan _2 15:04"))
	name := len(info.Name())
	Userr, _ := user.LookupId(fmt.Sprint(info.Sys().(*syscall.Stat_t).Uid))
	userr := len(Userr.Username)
	Group, _ := user.LookupGroupId(fmt.Sprint(info.Sys().(*syscall.Stat_t).Gid))
	group := len(Group.Name)
	Major := 0
	Minor := 0

	for _, entry := range entries {
		info1, _ := entry.Info()

		link1 := len(strconv.Itoa(int(info1.Sys().(*syscall.Stat_t).Nlink)))
		size1 := len(strconv.Itoa(int(info1.Size())))
		if info.Mode()&os.ModeCharDevice != 0 || info.Mode()&os.ModeDevice != 0 {
			path := middlewares.JoinPaths(path, entry.Name())
			major, minor := middlewares.MajorMinor(path)
			majors := len(strconv.Itoa(major))
			minors := len(strconv.Itoa(minor))
			if minors > Minor {
				Minor = minors
			}
			if majors > Major {
				Major = majors
			}
			size1 = len(strconv.Itoa(major)+",") + len(strconv.Itoa(minor)) + 1
		}
		perm1 := len(middlewares.ModeToString(info.Mode()))

		date1 := len(info1.ModTime().Format("Jan _2 15:04"))
		name1 := len(info1.Name())
		Userr1, _ := user.LookupId(fmt.Sprint(info1.Sys().(*syscall.Stat_t).Uid))
		userr1 := len(Userr1.Username)
		Group1, _ := user.LookupGroupId(fmt.Sprint(info1.Sys().(*syscall.Stat_t).Gid))

		group1 := len(Group1.Name)
		if perm1 > perm {
			perm = perm1
		}
		if link < link1 {
			link = link1
		}
		if size < size1 {
			size = size1
		}
		if date < date1 {
			date = date1
		}
		if name < name1 {
			name = name1
		}

		if userr < userr1 {
			userr = userr1
		}
		if group < group1 {
			group = group1
		}

		blocks += int(info1.Sys().(*syscall.Stat_t).Blocks)

	}
	width.Blocks = blocks
	width.Permw = perm
	width.Linkw = link
	width.Userrw = userr
	width.Groupw = group
	width.Sizew = size
	width.Datew = date
	width.Namew = name
	width.Major = Major
	width.Minor = Minor
	return width
	// fmt.Println(perm, link, userr, group, size, date, name)
}

func formatTime(t time.Time) string {
	now := time.Now()
	sixMonthsAgo := now.AddDate(0, -6, 0)

	if t.After(sixMonthsAgo) {
		return t.Format("Jan _2 15:04")
	} else {
		return t.Format("Jan _2  2006")
	}
}
