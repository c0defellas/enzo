package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

type mbr [512]byte

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

	// bootstrap code
	BCOffEnd = 0x01bd
	BCSize   = BCOffEnd
)

var (
	flags        *flag.FlagSet
	flagHelp     *bool
	flagCreate   *bool
	flagBootcode *string
)

func init() {
	flags = flag.NewFlagSet("mbr", flag.ContinueOnError)
	flagHelp = flags.Bool("help", false, "Show this help")
	flagCreate = flags.Bool("create", false, "Create new MBR")
	flagBootcode = flags.String("bootcode", "", "Bootsector binary code")
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

func NewMBR() *mbr {
	mbr := mbr{}
	mbr[Magic1Off] = Magic1
	mbr[Magic2Off] = Magic2
	return &mbr
}

func (m *mbr) SetBootcode(bcode []byte) error {
	if len(bcode) > BCSize {
		return fmt.Errorf("bootcode must have less than %d bytes", BCSize)
	}

	if copied := copy(m[0:BCOffEnd], bcode[:]); copied != len(bcode) {
		return fmt.Errorf("Failed to copy bootcode to mbr")
	}

	return nil
}

func mbrCreate(devfname, bootfname string) error {
	mbr := NewMBR()

	devfile, err := os.OpenFile(devfname, os.O_RDWR, 0)

	if err != nil {
		return err
	}

	if bootfname != "" {
		bfile, err := os.Open(bootfname)

		if err != nil {
			return err
		}

		var bootcode [BCSize + 1]byte

		n, err := bfile.Read(bootcode[:])

		if err == io.EOF || n > BCSize {
			return fmt.Errorf("bootcode must have less than %d bytes. Got %d", BCSize, n)
		}

		err = mbr.SetBootcode(bootcode[0:n])

		if err != nil {
			return err
		}
	}

	_, err = devfile.Write(mbr[:])

	return err
}

func runmbr(args []string) error {
	flags.Parse(args[1:])

	if *flagHelp {
		flags.PrintDefaults()
		return nil
	}

	disks := flags.Args()

	if len(disks) != 1 {
		return fmt.Errorf("Require one device file")
	}

	if *flagCreate {
		return mbrCreate(disks[0], *flagBootcode)
	}

	for _, disk := range disks {
		err := diskinfo(disk)

		if err != nil {
			return err
		}
	}

	return nil
}
