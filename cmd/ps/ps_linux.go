// +build linux
package main

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"syscall"
)

func pidowner(pid os.FileInfo) (string, error) {
	var (
		stat *syscall.Stat_t
		ok   bool
	)

	if stat, ok = pid.Sys().(*syscall.Stat_t); !ok {
		return "", fmt.Errorf("Failed to get file owner")
	}

	user, err := user.LookupId(strconv.Itoa(int(stat.Uid)))

	if err != nil {
		return "", err
	}

	return user.Username, nil

}
