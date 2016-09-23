package main

import "os"

func echo(args []string, newline bool) {
	last := len(args) - 1

	for i := 0; i < len(args); i++ {
		os.Stdout.Write([]byte(args[i]))

		if i < last {
			os.Stdout.Write([]byte{' '})
		}
	}

	if newline {
		os.Stdout.Write([]byte{0x0a})
	}
}

func main() {
	newline := true
	start := 1
	last := len(os.Args) - 1

	if last > 0 && os.Args[1] == "-n" {
		newline = false
		start = 2
	}

	echo(os.Args[start:], newline)
}
