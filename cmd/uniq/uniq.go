package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

type options struct {
	printDuplicates bool
	printEmptyLines bool
	printEveryOnce  bool
	printLineNumber bool
}

type Line struct {
	text *string
	nums []int
}

func usage() {
	fmt.Println(`Usage:
uniq [[-dup | -every] -empty | -num]`)
}

func parseArgs(args []string) (options, error) {
	var opts options

	for _, opt := range args[1:] {
		switch opt {
		case "-dup":
			opts.printDuplicates = true
		case "-empty":
			opts.printEmptyLines = true
		case "-every":
			opts.printEveryOnce = true
		case "-num":
			opts.printLineNumber = true
		default:
			usage()
			return options{}, errors.New("Wrong option")
		}
	}
	if opts.printEveryOnce && opts.printDuplicates {
		usage()
		return options{}, errors.New("Choose -dup OR -every")
	}

	return opts, nil
}

func shouldAddLine(linep *Line, opts options, emptyAdded bool) (bool, bool) {
	if *linep.text == "\n" {
		if opts.printEmptyLines && !emptyAdded {
			return true, true
		}
	} else if len(linep.nums) == 1 {
		return true, emptyAdded
	}

	return false, emptyAdded
}

func scanLines(input io.Reader, opts options) ([]*Line, error) {
	reader := bufio.NewReader(input)
	linesPtrMap := make(map[string]*Line)
	var linesOrdered []*Line

	lineNum := 0
	add := false
	emptyAdded := false
	for {
		lineNum++
		lineb, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		lineStr := string(lineb)
		linep := linesPtrMap[lineStr]
		if linep == nil {
			linep = &Line{text: &lineStr}
			linesPtrMap[lineStr] = linep
		}
		linep.nums = append(linep.nums, lineNum)
		add, emptyAdded = shouldAddLine(linep, opts, emptyAdded)
		if add {
			linesOrdered = append(linesOrdered, linep)
		}
	}

	return linesOrdered, nil
}

func shouldPrint(linep *Line, opts options) bool {
	lineCount := len(linep.nums)

	if *linep.text == "\n" {
		return true
	}
	if opts.printDuplicates {
		if lineCount > 1 {
			return true
		}
	} else if opts.printEveryOnce || lineCount == 1 {
		return true
	}

	return false
}

func printLineNumbers(linep *Line) {
	for i, ln := range linep.nums {
		fmt.Print(ln)
		if i < len(linep.nums)-1 {
			fmt.Print(",")
		}
	}
	fmt.Print(": ")
}

func uniq(linesp []*Line, opts options) {
	for _, linep := range linesp {
		if !shouldPrint(linep, opts) {
			continue
		}
		if opts.printLineNumber {
			printLineNumbers(linep)
		}
		fmt.Print(*linep.text)
	}
}

func main() {
	opts, err := parseArgs(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	linesp, err := scanLines(os.Stdin, opts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	uniq(linesp, opts)
}
