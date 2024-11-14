package formating

import (
	"fmt"
	"io/fs"
	"os/user"
	"strconv"
	"syscall"

	"github.com/anamivale/ls/options"
)

func LongFormat(entries []fs.DirEntry, flags options.Flags) {
	blocks,permw,userrw, groupw, linkw, sizew, datwe, namew := GetBlocks(entries)
	
	fmt.Printf("total %d\n", blocks/2)
	for _, entry := range entries {
		info, _ := entry.Info()

		perm := info.Mode()
		link := info.Sys().(*syscall.Stat_t).Nlink
		size := info.Size()
		date := info.ModTime().Format("Jan _2 15:04")
		name := entry.Name()
		gid := fmt.Sprint(info.Sys().(*syscall.Stat_t).Gid)
		uid := fmt.Sprint(info.Sys().(*syscall.Stat_t).Uid)

		Userr, _ := user.LookupId(uid)
		userr := Userr.Username
		Group, _ := user.LookupGroupId(gid)
		group := Group.Name
		format := fmt.Sprintf("%%%ds %%-%dd %%-%ds %%%ds %%%dd %%%ds %%-%ds\n",
		permw, linkw, userrw, groupw, sizew, datwe, namew)

	// Print the formatted line
	fmt.Printf(format, perm, link, userr, group, size, date, name)

	}
}

func Format(entries []fs.DirEntry) {
	for _, entry := range entries {
		fmt.Printf("%s  ", entry.Name())
	}
	fmt.Println()
}
func GetBlocks(entries []fs.DirEntry) (int, int, int, int, int, int, int, int) {
	blocks := 0
	info, _ := entries[0].Info()
	perm := len(info.Mode().String())

	link := len(strconv.Itoa(int(info.Sys().(*syscall.Stat_t).Nlink)))
	size := len(strconv.Itoa(int(info.Size())))
	date := len(info.ModTime().Format("Jan _2 15:04"))
	name := len(info.Name())
	Userr, _ := user.LookupId(fmt.Sprint(info.Sys().(*syscall.Stat_t).Uid))
	userr := len(Userr.Username)
	Group, _ := user.LookupGroupId(fmt.Sprint(info.Sys().(*syscall.Stat_t).Gid))
	group := len(Group.Name)

	for _, entry := range entries {
		info1, _ := entry.Info()

		link1 := len(strconv.Itoa(int(info1.Sys().(*syscall.Stat_t).Nlink)))
		size1 := len(strconv.Itoa(int(info1.Size())))
		date1 := len(info1.ModTime().Format("Jan _2 15:04"))
		name1 := len(info1.Name())
		Userr1, _ := user.LookupId(fmt.Sprint(info1.Sys().(*syscall.Stat_t).Uid))
		userr1 := len(Userr1.Username)
		Group1, _ := user.LookupGroupId(fmt.Sprint(info1.Sys().(*syscall.Stat_t).Gid))
		group1 := len(Group1.Name)

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
	// fmt.Println(perm, link, userr, group, size, date, name)

	return blocks,perm, link, userr, group, size, date, name
}
