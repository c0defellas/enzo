// +build windows

package killer

import (
        "os"
)

func Kill(pids []int, safe bool) (map[int]error) {
	var errm = make(map[int]error)
	var err error

	for _, pid := range pids {
		proc := os.Process{Pid: pid}

		err = proc.Kill()
		if err != nil {
			errm[pid] = err
		}
	}
	return errm
}
