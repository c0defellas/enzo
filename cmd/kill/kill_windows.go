// +build windows

package main

import (
	"os"
)

func kill(pids []int, safe bool) map[int]error {
	var errs = make(map[int]error)

	for _, pid := range pids {
		proc := os.Process{Pid: pid}

		err := proc.Kill()
		if err != nil {
			errs[pid] = err
		}
	}
	return errs
}
