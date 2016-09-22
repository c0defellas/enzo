package main

import "os"

func main() {
	addline := true
	start := 1
	last := len(os.Args) - 1

	if last > 0 && os.Args[1] == "-n" {
		addline = false
		start = 2
	}

	for i := start; i < len(os.Args); i++ {
		os.Stdout.Write([]byte(os.Args[i]))

		if i < last {
			os.Stdout.Write([]byte{' '})
		}
	}

	if addline {
		os.Stdout.Write([]byte{0x0a})
	}
}
