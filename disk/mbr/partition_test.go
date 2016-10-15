package mbr

import "testing"

func testParse(entry [16]byte, expected *partition, t *testing.T) {
	part, err, empty := NewPartition(entry[:])

	if err != nil {
		t.Error(err)
		return
	}

	if empty {
		t.Errorf("Partition is empty")
		return
	}

	if !part.IsEqual(expected) {
		t.Errorf("Partitions differ: %v != %v", part, expected)
		return
	}
}

func TestPartitionFailParse(t *testing.T) {
	// empty partition entry
	_, err, _ := NewPartition([]byte{})

	if err == nil {
		t.Errorf("Partition entry must have at least 16 bytes")
		return
	}

	_, err, _ = NewPartition([]byte{
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff,
	}) // 15 bytes

	if err == nil {
		t.Errorf("Partition entry must have at least 16 bytes")
		return
	}

	_, err, empty := NewPartition([]byte{
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	})

	if err != nil {
		t.Error(err)
		return
	}

	if !empty {
		t.Errorf("Partition entry is empty")
		return
	}

	_, err, empty = NewPartition([]byte{
		0x80, // active partition
		0,    // wrong, sector must be >= 1
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	})

	if err == nil {
		t.Errorf("Must fail, first sector must be >= 1")
		return
	}

	if empty {
		t.Error("is not empty")
		return
	}
}

func TestPartitionOK(t *testing.T) {
	for _, test := range []struct {
		entry    [16]byte
		expected partition
	}{
		{
			entry: [16]byte{
				0, 0, 1, 0,
				0x83,
				255, 0x3f | 0xc0, 0xff, 0, 0,
				0, 0},
			expected: partition{
				status: 0,
				begin: chs{
					head:     0,
					sector:   1,
					cylinder: 0,
				},
				typ: 0x83,
				end: chs{
					head:     255,
					sector:   63,
					cylinder: 1023,
				},
				lba:      0,
				nsectors: 0,
			},
		},
	} {
		testParse(test.entry, &test.expected, t)
	}
}
