package formating

import (
	"fmt"
	"io/fs"
	"os/user"
	"syscall"

	"github.com/anamivale/ls/options"
)

func LongFormat(entries []fs.DirEntry, flags options.Flags) {
	for _, entry := range entries {
		info, _ := entry.Info()

		perm := info.Mode()
		link := getLinks(info)
		size := info.Size()
		date := info.ModTime().Format("Jan _2 15:04")
		name := entry.Name()
		gid := fmt.Sprint(info.Sys().(*syscall.Stat_t).Gid)
		uid := fmt.Sprint(info.Sys().(*syscall.Stat_t).Uid)

		Userr, _ := user.LookupId(uid)
		userr := Userr.Username
		Group, _ := user.LookupGroupId(gid)
		group := Group.Name
		fmt.Println(perm, link, userr, group, size, date, name)

	}
}

func Format(entries []fs.DirEntry) {
	for _, entry := range entries {
		fmt.Printf("%s ", entry.Name())
	}
	fmt.Println()
}

func getLinks(info fs.FileInfo) uint64 {
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		return stat.Nlink
	}
	return 1
}
