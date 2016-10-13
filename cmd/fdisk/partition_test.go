package main

import "testing"

func TestPartitionParse(t *testing.T) {
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
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
	})

	if err == nil {
		t.Error("Must fail, bootable partition but none valid entries")
		return
	}

	if empty {
		t.Error("is not empty")
		return
	}
}
