package formating

import (
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"strconv"
	"syscall"

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

	for _, entry := range entries {
		info, _ := entry.Info()
		path1 := middlewares.JoinPaths(path, entry.Name())

		perm := middlewares.ModeToString(info.Mode())
		link := info.Sys().(*syscall.Stat_t).Nlink
		size := strconv.Itoa(int(info.Size()))

		date := info.ModTime().Format("Jan _2 15:04")
		name := entry.Name()
		gid := fmt.Sprint(info.Sys().(*syscall.Stat_t).Gid)
		uid := fmt.Sprint(info.Sys().(*syscall.Stat_t).Uid)

		Userr, _ := user.LookupId(uid)
		userr := Userr.Username
		Group, _ := user.LookupGroupId(gid)
		group := Group.Name
		if group == "video" {
			perm = perm + "+"
		}

		if info.Mode()&os.ModeCharDevice != 0 || info.Mode()&os.ModeDevice != 0 {
			major, minor := middlewares.MajorMinor(path1)
			if info.Mode()&os.ModeSymlink != 0 {

				target, err := os.Readlink(path1)
				name = name + " -> " + target
				if err == nil {
					format := fmt.Sprintf("%%%ds %%%dd %%-%ds %%%-ds %%%dd, %%%dd %%%ds %%-%ds\n",
						width.Permw, width.Linkw, width.Userrw, width.Groupw, width.Major, width.Minor, width.Datew, width.Namew)

					fmt.Printf(format, perm, link, userr, group, major, minor, date, name)
				}
			} else {
				format := fmt.Sprintf("%%%ds %%%dd %%-%ds %%%-ds %%%dd, %%%dd %%%ds %%-%ds\n",
					width.Permw, width.Linkw, width.Userrw, width.Groupw, width.Major, width.Minor, width.Datew, width.Namew)

				fmt.Printf(format, perm, link, userr, group, major, minor, date, name)
			}
		} else {
			if info.Mode()&os.ModeSymlink != 0 {

				target, err := os.Readlink(path1)
				name = name + " -> " + target
				if err == nil {
					format := fmt.Sprintf("%%%ds %%%dd %%-%ds %%%-ds %%%ds %%%ds %%-%ds\n",
						width.Permw, width.Linkw, width.Userrw, width.Groupw, width.Sizew, width.Datew, width.Namew)

					// Print the formatted line
					fmt.Printf(format, perm, link, userr, group, size, date, name)
				}
			} else {
				format := fmt.Sprintf("%%%ds %%%dd %%-%ds %%%-ds %%%ds %%%ds %%-%ds\n",
					width.Permw, width.Linkw, width.Userrw, width.Groupw, width.Sizew, width.Datew, width.Namew)

				// Print the formatted line
				fmt.Printf(format, perm, link, userr, group, size, date, name)
			}
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
			size1 = len(strconv.Itoa(major)+",") + len(strconv.Itoa(minor)) + 2
		}
		perm1 := len(middlewares.ModeToString(info.Mode()))

		date1 := len(info1.ModTime().Format("Jan _2 15:04"))
		name1 := len(info1.Name())
		Userr1, _ := user.LookupId(fmt.Sprint(info1.Sys().(*syscall.Stat_t).Uid))
		userr1 := len(Userr1.Username)
		Group1, _ := user.LookupGroupId(fmt.Sprint(info1.Sys().(*syscall.Stat_t).Gid))
		if Group1.Name == "video" {
			perm1 = perm1 + len("+")
		}
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
