package main

import (
	"fmt"
	"os"

	"github.com/tiago4orion/enzo/cmd/fdisk/cmds"
)

var (
	perr = func(format string, args ...interface{}) (int, error) {
		return fmt.Fprintf(os.Stderr, format, args...)
	}
)

func usage() {
	perr("fdisk <subcommand> [options] device/file\n")
	perr("Subcommands:\n")
	perr("\tmbr\n")
	perr("\n")
	perr("Use: fdisk <subcommand> -h for more info\n")
}

func main() {
	var err error

	if len(os.Args) <= 1 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "mbr":
		err = cmds.MBR(os.Args[1:])
	default:
		perr("Invalid subcommand: %s\n", os.Args[1])
		os.Exit(1)
	}

	if err != nil {
		perr("error: %s\n", err)
		os.Exit(1)
	}
}
