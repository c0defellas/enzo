package mbr

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type (
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

func (p partition) Bytes() []byte {
	var data [16]byte

	data[0] = byte(p.status)
	data[1] = p.begin.head
	data[2] = byte((p.begin.cylinder >> 6) | uint16(p.begin.sector&0x3f))
	return data[:]
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
	part.begin.sector = getsector(entry[2])

	if part.begin.sector == 0 {
		return nil, fmt.Errorf("First sector must be >= 1. Found %d",
			part.begin.sector), false
	}

	part.begin.cylinder = getcylinder(entry[2], entry[3])
	part.typ = typ(entry[4])
	part.end.head = uint8(entry[5])
	part.end.sector = getsector(entry[6])
	part.end.cylinder = getcylinder(entry[6], entry[7])

	buf := bytes.NewBuffer(entry[8:12])
	err := binary.Read(buf, binary.LittleEndian, &part.lba)

	if err != nil {
		return nil, err, false
	}

	buf = bytes.NewBuffer(entry[12:16])
	err = binary.Read(buf, binary.LittleEndian, &part.nsectors)

	return part, err, false
}

func NewEmptyPartition() *partition {
	return &partition{}
}

func (p *partition) IsEqual(other *partition) bool {
	if p == other {
		return true
	}

	if p.status != other.status ||
		p.begin.head != other.begin.head ||
		p.begin.sector != other.begin.sector ||
		p.begin.cylinder != other.begin.cylinder ||
		p.typ != other.typ ||
		p.end.head != other.end.head ||
		p.end.sector != other.end.sector ||
		p.end.cylinder != other.end.cylinder ||
		p.lba != other.lba ||
		p.nsectors != other.nsectors {
		return false
	}

	return true
}

func isPartEmpty(buf []byte) bool {
	for i := 0; i < 16; i++ {
		if buf[i] != 0 {
			return false
		}
	}

	return true
}

func getsector(b byte) uint8 {
	return b & 0x3f // sector in bits 5-0
}

func getcylinder(b1, b2 byte) uint16 {
	return uint16(b1&0xc0)<<2 | uint16(b2)
}
