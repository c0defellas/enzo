package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

const (
	// Classical MBR structure
	CPart1 = 0x1be // 16 bytes each
	CPart2 = 0x1ce
	CPart3 = 0x1de
	CPart4 = 0x1ee

	PEntrySZ = 16

	Magic1Off = 0x1fe
	Magic2Off = 0x1ff

	// MBR magic numbers
	Magic1 = 0x55
	Magic2 = 0xaa
)

func diskinfo(fname string) error {
	file, err := os.Open(fname)

	if err != nil {
		return err
	}

	var mbr [512]byte

	_, err = file.Read(mbr[:])

	if err != nil && err != io.EOF {
		return err
	}

	if mbr[Magic1Off] != Magic1 &&
		mbr[Magic2Off] != Magic2 {
		return fmt.Errorf("no MBR found\n")
	}

	var part *partition
	var empty bool

	if part, err, empty = NewPartition(mbr[CPart1 : CPart1+PEntrySZ]); err != nil {
		return err
	} else if !empty {
		fmt.Printf("%s\n", part)
	}

	if part, err, empty = NewPartition(mbr[CPart2 : CPart2+PEntrySZ]); err != nil {
		return err
	} else if !empty {
		fmt.Printf("%s\n", part)
	}

	if part, err, empty = NewPartition(mbr[CPart3 : CPart3+PEntrySZ]); err != nil {
		return err
	} else if !empty {
		fmt.Printf("%s\n", part)
	}

	if part, err, empty = NewPartition(mbr[CPart4 : CPart4+PEntrySZ]); err != nil {
		return err
	} else if !empty {
		fmt.Printf("%s\n", part)
	}

	return nil
}

func main() {
	flag.Parse()
	diskpaths := flag.Args()

	if len(diskpaths) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s dev1 dev2 ...\n", os.Args[0])
		os.Exit(1)
	}

	for _, disk := range diskpaths {
		err := diskinfo(disk)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get MBR/GPT info: %s", err)
			os.Exit(1)
		}
	}

}
