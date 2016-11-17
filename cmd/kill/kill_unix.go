// +build linux darwin dragonfly freebsd netbsd openbsd

package main

import (
	"syscall"
)

func kill(pids []int, safe bool) map[int]error {
	var errs = make(map[int]error)

	for _, pid := range pids {
		err := syscall.Kill(pid, syscall.SIGTERM)

		if err != nil {
			if safe == true {
				errs[pid] = err
				continue
			}

			err = syscall.Kill(pid, syscall.SIGKILL)

			if err != nil {
				errs[pid] = err
				continue
			}
		}
	}

	return errs
}
