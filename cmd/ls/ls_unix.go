// +build linux darwin dragonfly freebsd netbsd openbsd

package main

import (
	"errors"
	"os"
	"os/user"
	"strconv"
	"syscall"
)

func toStatT(fileInfo os.FileInfo) (*syscall.Stat_t, error) {
	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return nil, errors.New("Could not get file stat")
	}
	return stat, nil
}

func lookupUser(fileInfo os.FileInfo) (string, error) {
	statt, err := toStatT(fileInfo)
	if err != nil {
		return "", err
	}

	usr, err := user.LookupId(strconv.FormatUint(uint64(statt.Uid), 10))
	if err != nil {
		return "", err
	}

	return usr.Username, nil
}

func lookupGroup(fileInfo os.FileInfo) (string, error) {
	statt, err := toStatT(fileInfo)
	if err != nil {
		return "", err
	}

	group, err := user.LookupGroupId(strconv.FormatUint(uint64(statt.Gid), 10))
	if err != nil {
		return "", err
	}

	return group.Name, nil
}
