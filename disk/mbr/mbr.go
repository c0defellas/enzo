package mbr

import (
	"fmt"
	"io"
	"math"
	"os"
)

type (
	mbr [512]byte

	// cylinder-head-sector
	chs struct {
		head, sector uint8
		cylinder     uint16
	}
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

	// bootstrap code
	BCOffEnd = 0x01bd
	BCSize   = BCOffEnd
)

func NewCHS(cylinder uint16, head uint8, sector uint8) chs {
	return chs{
		cylinder: cylinder,
		head:     head,
		sector:   sector,
	}
}

func Info(fname string) error {
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

func Create(devfname, bootfname string) error {
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

func CHS2LBA(c uint16, h uint8, s uint8) uint32 {
	return (uint32(c)*16+uint32(h))*63 + uint32((s - 1))
}

func LBA2C(lba uint32) uint16 {
	return uint16(math.Mod(float64(lba), 16*63))
}

func LBA2H(lba uint32) uint8 {
	lbaf := float64(lba)
	spt := float64(63)
	hpt := float64(16)
	return uint8(math.Mod(math.Mod(lbaf, spt), hpt))
}

func LBA2S(lba uint32) uint8 {
	return uint8(math.Mod(float64(lba), float64(63))) + 1
}
