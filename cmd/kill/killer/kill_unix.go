// +build linux darwin dragonfly freebsd netbsd openbsd

package killer

import (
	"syscall"
)

func Kill(pids []int, safe bool) (map[int]error) {
	var errm = make(map[int]error)
	var err error

	for _, pid := range pids {
		err = syscall.Kill(pid, syscall.SIGTERM)
		if err != nil {
			if safe == true {
				goto errloop
			}
			err = syscall.Kill(pid, syscall.SIGKILL)
			if err != nil {
				goto errloop
			}
		}
		continue
errloop:
		errm[pid] = err
	}
	return errm
}
