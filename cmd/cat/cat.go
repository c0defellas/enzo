// Based on Plan9 cat

package main

import (
	"log"
	"os"
)

func cat(f *os.File, name string) {
	var (
		buf [8192]byte
	)

	for {
		n, err := f.Read(buf[:])

		if err != nil {
			log.Fatal("error reading %s: %s", name, err)
		}

		if n <= 0 {
			break
		}

		b := buf[:n]

		if n, err = os.Stdout.Write(b); n != len(b) {
			log.Fatal("write error copying %s: %s", name, err)
		}
	}
}

func main() {
	if len(os.Args) == 1 {
		cat(os.Stdin, "<stdin>")
	} else {
		for i := 1; i < len(os.Args); i++ {
			f, err := os.Open(os.Args[i])

			if err != nil {
				log.Fatal("can't open %s", os.Args[i])
			}

			cat(f, os.Args[i])
			f.Close()
		}
	}
}
