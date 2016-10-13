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

var (
	flags    *flag.FlagSet
	flagHelp *bool
)

func init() {
	flags = flag.NewFlagSet("mbr", flag.ContinueOnError)
	flagHelp = flags.Bool("help", false, "Show this help")
}

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

func mbr(args []string) error {
	flags.Parse(args[1:])

	if *flagHelp {
		flags.PrintDefaults()
		return nil
	}

	disks := flags.Args()

	if len(disks) == 0 {
		return fmt.Errorf("Require a device file")
	}

	for _, disk := range disks {
		err := diskinfo(disk)

		if err != nil {
			return err
		}
	}

	return nil
}
