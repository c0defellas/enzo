// Based on Plan9 cat

package main

import (
	"errors"
	"io"
	"os"
)

func fatal(err error) {
	os.Stderr.Write([]byte(err.Error()))
	os.Stderr.Write([]byte{0x0a})
	os.Exit(1)
}

func cat(in io.Reader, out io.Writer, name string) error {
	var (
		buf [8192]byte
	)

	for {
		n, err := in.Read(buf[:])

		if n <= 0 {
			break
		}

		if err != nil {
			return errors.New("error reading " + name + ": " + err.Error())
		}

		b := buf[:n]

		if n, err = out.Write(b); n != len(b) {
			return errors.New("write error copying " + name + ": " + err.Error())
		}
	}

	return nil
}

func runcat(files []string, out io.Writer) error {
	for _, fname := range files {
		f, err := os.Open(fname)

		if err != nil {
			return err
		}
		err = cat(f, out, fname)
		f.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		err := cat(os.Stdin, os.Stdout, "<stdin>")

		if err != nil {
			fatal(err)
		}

		return
	}

	if err := runcat(os.Args[1:], os.Stdout); err != nil {
		fatal(err)
	}
}
