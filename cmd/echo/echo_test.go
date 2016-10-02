package main

import (
	"bytes"
	"testing"
)

type testEcho struct {
	args     []string
	expected string
	newline  bool
}

func testecho(args []string, expected string, newline bool, t *testing.T) {
	var out bytes.Buffer

	echo(&out, args, newline)

	if string(out.Bytes()) != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, string(out.Bytes()))
		return
	}
}

func TestEcho(t *testing.T) {
	testTbl := []testEcho{
		{
			[]string{}, "", false, // echo -n
		},
		{
			[]string{"hello"}, "hello", false, // echo -n hello
		},
		{
			[]string{"hello", "world"}, "hello world", false, // echo -n hello world
		},
		{
			[]string{"hello", "world"}, `hello world
`, true, // echo hello world
		},
		{
			[]string{`hello
world`}, `hello
world`, false},
	}

	for _, test := range testTbl {
		testecho(test.args, test.expected, test.newline, t)
	}
}

type testArgs struct {
	args    []string
	newline bool
}

func testParseArgs(args []string, expected testArgs, t *testing.T) {
	parsed, newline := parsearg(args)

	if len(parsed) != len(expected.args) {
		t.Errorf("Expect no args")
		return
	}

	if newline != expected.newline {
		t.Errorf("Expected %v newline", expected.newline)
		return
	}

	for i := 0; i < len(parsed); i++ {
		if parsed[i] != expected.args[i] {
			t.Errorf("Expected '%s' but got '%s'", parsed[i], expected.args[i])
			return
		}
	}
}

func TestParseArg(t *testing.T) {
	testTbl := []struct {
		args     []string
		expected testArgs
	}{
		{
			[]string{},
			testArgs{[]string{}, true},
		},
		{
			[]string{"echo"},
			testArgs{[]string{}, true},
		},
		{
			[]string{"echo", "-n"},
			testArgs{[]string{}, false},
		},
		{
			[]string{"echo", "-n", "hello"},
			testArgs{[]string{"hello"}, false},
		},
	}

	for _, test := range testTbl {
		testParseArgs(test.args, test.expected, t)
	}
}
