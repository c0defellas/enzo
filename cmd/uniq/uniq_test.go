package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
)

var input = []byte(`hello
world

hello

世界
世界
世
1
3
4
日本語
4
1
`)

func newLine(s string, ln []int) *Line {
	str := s + "\n"
	return &Line{
		text: &str,
		nums: ln,
	}
}

var ttScanLines = []struct {
	opts     options
	expected []*Line
}{
	{
		options{printEmptyLines: false},
		[]*Line{
			newLine("hello", []int{1, 4}),
			newLine("world", []int{2}),
			newLine("世界", []int{6, 7}),
			newLine("世", []int{8}),
			newLine("1", []int{9, 14}),
			newLine("3", []int{10}),
			newLine("4", []int{11, 13}),
			newLine("日本語", []int{12}),
		},
	},
	{
		options{printEmptyLines: true},
		[]*Line{
			newLine("hello", []int{1, 4}),
			newLine("world", []int{2}),
			newLine("", []int{3, 5}),
			newLine("世界", []int{6, 7}),
			newLine("世", []int{8}),
			newLine("1", []int{9, 14}),
			newLine("3", []int{10}),
			newLine("4", []int{11, 13}),
			newLine("日本語", []int{12}),
		},
	},
}

var ttUniq = []struct {
	opts     options
	lines    []*Line
	expected string
}{
	{
		// uniq
		options{},
		ttScanLines[0].expected,
		`world
世
3
日本語
`,
	},
	{
		// uniq -num
		options{printLineNumber: true},
		ttScanLines[0].expected,
		`2: world
8: 世
10: 3
12: 日本語
`,
	},
	{
		// uniq -empty
		options{printEmptyLines: true},
		ttScanLines[1].expected,
		`world

世
3
日本語
`,
	},
	{
		// uniq -empty -num
		options{
			printEmptyLines: true,
			printLineNumber: true,
		},
		ttScanLines[1].expected,
		"2: world\n" +
			"3,5: \n" +
			`8: 世
10: 3
12: 日本語
`,
	},
	{
		// uniq -dup
		options{printDuplicates: true},
		ttScanLines[0].expected,
		`hello
世界
1
4
`,
	},
	{
		// uniq -dup -num
		options{
			printDuplicates: true,
			printLineNumber: true,
		},
		ttScanLines[0].expected,
		`1,4: hello
6,7: 世界
9,14: 1
11,13: 4
`,
	},
	{
		// uniq -dup -empty
		options{
			printDuplicates: true,
			printEmptyLines: true,
		},
		ttScanLines[1].expected,
		`hello

世界
1
4
`,
	},
	{
		// uniq -dup -empty -num
		options{
			printDuplicates: true,
			printEmptyLines: true,
			printLineNumber: true,
		},
		ttScanLines[1].expected,
		"1,4: hello\n" +
			"3,5: \n" +
			`6,7: 世界
9,14: 1
11,13: 4
`,
	},

	{
		// uniq -every
		options{printEveryOnce: true},
		ttScanLines[0].expected,
		`hello
world
世界
世
1
3
4
日本語
`,
	},
	{
		// uniq -every -num
		options{
			printEveryOnce:  true,
			printLineNumber: true,
		},
		ttScanLines[0].expected,
		`1,4: hello
2: world
6,7: 世界
8: 世
9,14: 1
10: 3
11,13: 4
12: 日本語
`,
	},
	{
		// uniq -every -empty
		options{
			printEveryOnce:  true,
			printEmptyLines: true,
		},
		ttScanLines[1].expected,
		`hello
world

世界
世
1
3
4
日本語
`,
	},
	{
		// uniq -every -empty -num
		options{
			printEveryOnce:  true,
			printEmptyLines: true,
			printLineNumber: true,
		},
		ttScanLines[1].expected,
		"1,4: hello\n" +
			"2: world\n" +
			"3,5: \n" +
			`6,7: 世界
8: 世
9,14: 1
10: 3
11,13: 4
12: 日本語
`,
	},
}

func cmpLines(a, b []*Line) error {
	if len(a) != len(b) {
		return fmt.Errorf("%T sizes diverge: %v != %v", a, len(a), len(b))
	}
	for i, _ := range a {
		linea := a[i]
		lineb := b[i]
		if *linea.text != *lineb.text {
			return fmt.Errorf("line %d:\n%v !=\n%v",
				i+1, *linea.text, *lineb.text)
		}

		if lena, lenb := len(linea.nums), len(lineb.nums); lena != lenb {
			return fmt.Errorf("line counts diverge: %v != %v", lena, lenb)
		}
		for n, _ := range linea.nums {
			lna := linea.nums[n]
			lnb := lineb.nums[n]
			if lna != lnb {
				return fmt.Errorf("line number index: %d -"+
					"numbers diverge: %v != %v", n, lna, lnb)
			}
		}
	}
	return nil
}

func TestScanLines(t *testing.T) {
	for i, tScanLines := range ttScanLines {
		bytesReader := bytes.NewReader(input)
		scannedLines, err := scanLines(bytesReader, tScanLines.opts)
		if err != nil {
			t.Fatalf("test index %d: %v", i, err)
		}
		if scannedLines == nil {
			t.Fatalf("test index %d: %v %[2]T %[2]v", i, scannedLines)
		}

		err = cmpLines(scannedLines, tScanLines.expected)
		if err != nil {
			t.Fatalf("test index %d: %v", i, err)
		}
	}
}

func getUniqOutput(lines []*Line, opts options) (string, error) {
	var buf bytes.Buffer
	var err error

	rEnd, wEnd, err := os.Pipe()
	if err != nil {
		return "", err
	}

	savedStdout := os.Stdout
	defer func() { os.Stdout = savedStdout }()
	os.Stdout = wEnd
	uniq(lines, opts)
	wEnd.Close()

	_, err = io.Copy(&buf, rEnd)
	rEnd.Close()
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func TestUniq(t *testing.T) {
	for i, tUniq := range ttUniq {
		output, err := getUniqOutput(tUniq.lines, tUniq.opts)
		if err != nil {
			t.Fatal(err)
		}
		if output != tUniq.expected {
			t.Fatalf("test index: %d\nexpected:\n%v\noutput:\n%v",
				i, tUniq.expected, output)
		}
	}
}
