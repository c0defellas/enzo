package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type (
	// cylinder-head-sector
	chs struct {
		head, sector uint8
		cylinder     uint16
	}

	status byte
	typ    byte

	partition struct {
		number   uint8
		status   status
		begin    chs
		typ      typ
		end      chs
		lba      int32
		nsectors int32
	}
)

func (chs chs) String() string {
	return fmt.Sprintf("%d/%d/%d", chs.cylinder, chs.head, chs.sector)
}

func (t typ) String() string {
	switch t {
	case 0x83:
		return fmt.Sprintf("%x (Linux)", int(t))
	}
	return "unknown"

}

func (st status) String() string {
	err := ""
	if st&0x80 != 0 && st&0x80 != 0x80 {
		err += " (wrong)"
	}

	switch st >> 7 {
	case 1:
		return "active" + err
	}

	return "inactive" + err
}

func (p partition) String() string {
	format := `Partition #%d
Status: %s
FS type: %s
First C/H/S: %s
Last C/H/S: %s
LBA: %d
Number of sectors: %d
`
	return fmt.Sprintf(format, p.number, p.status, p.typ, p.begin, p.end, p.lba, p.nsectors)
}

func NewPartition(entry []byte) (*partition, error, bool) {
	part := &partition{}

	if len(entry) != 16 {
		return nil, fmt.Errorf("Invalid partition entry: %v", entry), false
	}

	if isPartEmpty(entry) {
		return nil, nil, true
	}

	part.status = status(entry[0])
	part.begin.head = uint8(entry[1])
	part.begin.sector = uint8(entry[2])
	part.begin.cylinder = uint16(entry[3])
	part.typ = typ(entry[4])
	part.end.head = uint8(entry[5])
	part.end.sector = uint8(entry[6] & 0x3f) // sector in bits 5-0
	part.end.cylinder = uint16(entry[6]&0xc0)<<2 | uint16(entry[7])

	buf := bytes.NewBuffer(entry[8:12])
	err := binary.Read(buf, binary.LittleEndian, &part.lba)

	if err != nil {
		return nil, err, false
	}

	buf = bytes.NewBuffer(entry[12:16])
	err = binary.Read(buf, binary.LittleEndian, &part.nsectors)

	return part, err, false
}

func isPartEmpty(buf []byte) bool {
	for i := 0; i < 16; i++ {
		if buf[i] != 0 {
			return false
		}
	}

	return true
}
