package main

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func readAll(t *testing.T, reader io.Reader) string {
	content, err := ioutil.ReadAll(reader)

	if err != nil {
		t.Fatal(err)
	}

	return string(content)
}

func writeAll(t *testing.T, w io.Writer, data string) {
	_, err := w.Write([]byte(data))

	if err != nil {
		t.Fatal(err)
	}
}

// Only tests that cat copies the content of in into out
func testCatInOut(testList []string, t *testing.T) {
	for _, expected := range testList {
		t.Run(expected, func(t *testing.T) {
			var in, out bytes.Buffer

			writeAll(t, &in, expected)

			err := cat(&in, &out, "test")

			if err != nil {
				t.Error(err)
				return
			}

			got := readAll(t, &out)

			if got != expected {
				t.Errorf("Expected '%s' but got '%s'", expected, got)
				return
			}
		})
	}
}

func TestCatInOut(t *testing.T) {
	testList := []string{
		"",
		"test",
		"AAAAAAAAAAAAAAAAAAAAAAAAAA",
		`1
2
3
4
5
6
7
8
9
10
`,
	}

	testCatInOut(testList, t)
}

func writeOnTempfile(t *testing.T, contents string) string {
	f, err := ioutil.TempFile("/tmp", "enzo-test-cat")

	if err != nil {
		t.Fatal(err)
	}

	writeAll(t, f, contents)
	f.Close()
	return f.Name()
}

func testCatFile(expected string, t *testing.T) {
	filename := writeOnTempfile(t, expected)

	var out bytes.Buffer

	err := runcat([]string{filename}, &out)

	if err != nil {
		t.Fatal(err)
	}

	got := readAll(t, &out)

	if string(got) != expected {
		t.Fatalf("Got %q but expected was %q", string(got), expected)
	}
}

func TestHandlesMultipleFiles(t *testing.T) {
	data1 := "data1"
	data2 := "data2"
	filename1 := writeOnTempfile(t, data1)
	filename2 := writeOnTempfile(t, data2)

	var out bytes.Buffer

	err := runcat([]string{filename1, filename2}, &out)

	if err != nil {
		t.Fatal(err)
	}

	got := readAll(t, &out)
	expected := data1 + data2
	if string(got) != expected {
		t.Fatalf("Got %q but expected was %q", string(got), expected)
	}
}

func TestHandleOneFile(t *testing.T) {
	testTbl := []string{
		"",
		"1 line",
		`multi
		line`,
		`
`,
	}

	for _, test := range testTbl {
		testCatFile(test, t)
	}
}

func TestHandleFileNotFound(t *testing.T) {
	err := runcat([]string{"/<path-do-not-exists>"}, os.Stdout)

	if err == nil {
		t.Errorf("Must fail")
		return
	}
}

type fakeIO struct {
	Err error
	N   int
}

func (f fakeIO) Read(p []byte) (n int, err error) {
	return f.N, f.Err
}

func (f fakeIO) Write(p []byte) (n int, err error) {
	return f.N, f.Err
}

func TestHandleReadError(t *testing.T) {
	in := fakeIO{N: 1, Err: errors.New("injectedReadError")}
	out := fakeIO{}
	err := cat(in, out, "readError")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestHandleWriteError(t *testing.T) {
	in := fakeIO{N: 2}
	out := fakeIO{N: 1, Err: errors.New("injectedWriteError")}
	err := cat(in, out, "writeError")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}
