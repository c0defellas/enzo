package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type (
	status []string
)

const procDir = "/proc"

const (
	STpid      uint8 = 0
	STcomm           = 1
	STstat           = 2
	STppid           = 3
	STpgrp           = 4
	STsid            = 5
	STnice           = 18
	STnthreads       = 19
	STvmsize         = 22
	STrss            = 23
	STproc           = 38
	STLAST           = 52
)

var (
	states map[rune]string
)

func init() {
	states = map[rune]string{
		'R': "running",
		'S': "sleeping",
		'D': "sleep/disk",
		'Z': "zombie",
		'T': "stopped",
		't': "tracing stop",
		'X': "dead",
		'x': "dead",
		'K': "wakekill",
		'W': "waking",
		'P': "Parked",
	}
}

func newStatus(statstr string) (status, error) {
	stat := strings.Split(statstr, " ")

	if len(stat) != STLAST {
		return nil, fmt.Errorf("malformed stat line: (length %d) %s",
			len(stat), statstr)
	}

	return status(stat), nil
}

func (s status) pid() string      { return s[STpid] }
func (s status) comm() string     { return s[STcomm] }
func (s status) state() string    { return s[STstat] }
func (s status) ppid() string     { return s[STppid] }
func (s status) nice() string     { return s[STnice] }
func (s status) nthreads() string { return s[STnthreads] }
func (s status) vmsize() string   { return s[STvmsize] }
func (s status) rss() string      { return s[STrss] }
func (s status) lastproc() string { return s[STproc] }

func (s status) String() string {
	return strings.Join(s, " ")
}

func getpids() ([]os.FileInfo, error) {
	var pids []os.FileInfo

	files, err := ioutil.ReadDir(procDir)

	if err != nil {
		return nil, err
	}

	for _, f := range files {
		fname := f.Name()

		if !f.IsDir() {
			continue
		}

		if len(fname) > 0 && (fname[0] < '0' || fname[0] > '9') {
			continue
		}

		pids = append(pids, f)
	}

	return pids, nil
}

func readall(r io.Reader) ([]byte, error) {
	data := make([]byte, 0, 8192)
	buf := make([]byte, 8192)

	for {
		n, err := r.Read(buf)

		if err != nil && err != io.EOF {
			return nil, err
		}

		if n <= 0 {
			break
		}

		data = append(data, buf[:n]...)
	}

	return data, nil
}

func cmdline(pid os.FileInfo) (string, error) {
	file, err := os.Open(procDir + "/" + pid.Name() + "/cmdline")

	if err != nil {
		return "", err
	}

	defer file.Close()

	cmdbytes, err := readall(file)

	if err != nil {
		return "", err
	}

	cmdline := string(bytes.Replace(cmdbytes, []byte{0}, []byte{' '}, -1))
	return cmdline, nil
}

func getstat(pid os.FileInfo) (status, error) {
	file, err := os.Open(procDir + "/" + pid.Name() + "/stat")

	if err != nil {
		return nil, err
	}

	defer file.Close()

	statbytes, err := readall(file)

	if err != nil {
		return nil, err
	}

	stat, err := newStatus(string(statbytes))

	if err != nil {
		return nil, err
	}

	return stat, nil
}

func showpid(pid os.FileInfo) error {
	user, err := pidowner(pid)

	if err != nil {
		return err
	}

	cmdline, err := cmdline(pid)

	if err != nil {
		return err
	}

	status, err := getstat(pid)

	if err != nil {
		return err
	}

	fmt.Printf("%s\t%s\t%.20s\n", status.pid(), user, cmdline)
	return nil
}

func ps() error {
	pids, err := getpids()

	if err != nil {
		return err
	}

	for _, pid := range pids {
		if err = showpid(pid); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	if err := ps(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
}
