package main

import (
	"io"
	"os"
)

func echo(out io.Writer, args []string, newline bool) {
	last := len(args) - 1

	for i := 0; i < len(args); i++ {
		out.Write([]byte(args[i]))

		if i < last {
			out.Write([]byte{' '})
		}
	}

	if newline {
		out.Write([]byte{0x0a})
	}
}

func parsearg(args []string) ([]string, bool) {
	newline := true

	if len(args) == 0 {
		return args, true
	}

	start := 1

	if len(args) > 1 && args[1] == "-n" {
		newline = false
		start = 2
	}

	return args[start:], newline
}

func main() {
	args, newline := parsearg(os.Args)
	echo(os.Stdout, args, newline)
}
