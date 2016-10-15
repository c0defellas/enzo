package mbr

import "testing"

type testTbl struct {
	chs chs
	lba uint32
}

func TestCHS2LBA(t *testing.T) {
	var tests []testTbl

	for i := 1; i < 64; i++ {
		tests = append(tests, struct {
			chs chs
			lba uint32
		}{
			chs: NewCHS(0, 0, uint8(i)),
			lba: uint32(i - 1),
		})
	}

	tests = append(tests, testTbl{
		chs: NewCHS(0, 1, 1),
		lba: 63,
	})

	tests = append(tests, testTbl{
		chs: NewCHS(0, 15, 1),
		lba: 945,
	})

	tests = append(tests, testTbl{
		chs: NewCHS(15, 15, 63),
		lba: 16127,
	})

	tests = append(tests, testTbl{
		chs: NewCHS(16319, 15, 63),
		lba: 16450559,
	})

	for _, test := range tests {
		if val := CHS2LBA(test.chs.cylinder, test.chs.head, test.chs.sector); val != test.lba {
			t.Errorf("Expected %d but got %d", test.lba, val)
			return
		}
	}

}
