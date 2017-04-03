package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

const (
	_          = iota
	KB float64 = 1 << (10 * iota)
	MB
	GB
)

type formatter func(os.FileInfo) (string, error)

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "error: %s\n", err)
	os.Exit(1)
}

func humanizeSize(size int64) string {
	switch sz := float64(size); {
	case (sz >= GB):
		return fmt.Sprintf("%.2fG", sz/GB)
	case (sz >= MB):
		return fmt.Sprintf("%.2fM", sz/MB)
	case (sz >= KB):
		return fmt.Sprintf("%.2fK", sz/KB)
	}

	return fmt.Sprintf("%d", size)
}

func formatFileName(name string) string {
	if strings.Contains(name, " ") {
		return "'" + name + "'"
	}
	return name
}

func printFileList(fileInfo os.FileInfo) (string, error) {
	userName, err := lookupUser(fileInfo)
	if err != nil {
		return "", err
	}
	groupName, err := lookupGroup(fileInfo)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(
		"%s %s %s %6s %s\n",
		fileInfo.Mode(),
		userName,
		groupName,
		humanizeSize(fileInfo.Size()),
		formatFileName(fileInfo.Name()),
	), nil
}

func printFileNames(fileInfo os.FileInfo) (string, error) {
	return fmt.Sprintf("%s\n", formatFileName(fileInfo.Name())), nil
}

func ls(files []os.FileInfo, writer io.Writer, fn formatter) error {
	for _, f := range files {
		txt, err := fn(f)
		if err != nil {
			return err
		}
		fmt.Fprint(writer, txt)
	}
	return nil
}

func main() {
	path := "."

	l := flag.Bool("l", false, "use a long listing format")
	flag.Parse()

	if len(flag.Args()) > 0 {
		path = flag.Args()[0]
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		fatal(err)
	}

	f := printFileNames

	if *l {
		f = printFileList
	}

	err = ls(files, os.Stdout, f)
	if err != nil {
		fatal(err)
	}
}
