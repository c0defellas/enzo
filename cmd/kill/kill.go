package main

import (
	"os"
	"strconv"
	"fmt"
	"flag"

	killer "github.com/tiago4orion/enzo/cmd/kill/killer"
)

// save own name for easier reference
var me = os.Args[0]

func usage() {
	fmt.Println("Usage:")
	fmt.Println("kill [-safe] pids")
	flag.PrintDefaults()
	os.Exit(22); // EINVAL
}

func sliceatoi(ss []string) ([]int, error) {
	is := make([]int, 0, len(ss))

	for _, s := range ss {
		i, err := strconv.Atoi(s)
		if err != nil {
			return []int{}, err
		}
		is = append(is, i)
	}
	return is, nil
}

func parseargs() ([]int, bool) {
	var safe bool

	flag.BoolVar(&safe, "safe", false, "doesn't use SIGKILL when SIGTERM fail (unix systems)")
	flag.Parse()

	pids, err := sliceatoi(flag.Args())
	if err != nil || len(pids) == 0 {
		usage()
	}
	return pids, safe
}

func main() {
	pids, safe := parseargs()
	errm := killer.Kill(pids, safe)

	// some went wrong
	if len(errm) > 0 {
		for p, e := range errm {
			fmt.Printf("%s: [%d] - %s\n", me, p, e.Error())
		}
		os.Exit(1)
	}
}
