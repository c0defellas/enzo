package main

import (
	"fmt"
	"os"
	"os/exec"
)

func run(cmdargs []string) {
	cmd := exec.Command(cmdargs[0], cmdargs[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Start()

}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s cmdname\n", os.Args[0])
		os.Exit(1)
	}

	run(os.Args[1:])
}
