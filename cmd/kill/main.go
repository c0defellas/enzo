package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

func usage() {
	fmt.Println("Usage:")
	fmt.Println("kill [-safe] pids")
	flag.PrintDefaults()
	os.Exit(1)
}

func sliceatoi(strNumbers []string) ([]int, error) {
	numbers := make([]int, 0, len(strNumbers))

	for _, str := range strNumbers {
		i, err := strconv.Atoi(str)

		if err != nil {
			return []int{}, err
		}

		numbers = append(numbers, i)
	}

	return numbers, nil
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
	errs := kill(pids, safe)

	// some went wrong
	if len(errs) > 0 {
		for pid, err := range errs {
			fmt.Printf("%s: [%d] - %s\n", os.Args[0], pid, err.Error())
		}

		os.Exit(1)
	}
}
