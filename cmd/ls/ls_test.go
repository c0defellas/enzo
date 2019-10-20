package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"testing"
	"text/template"
)

func setup(t *testing.T) (string, func()) {
	dir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Error(err)
	}

	createTempFiles(t, dir)
	return dir, func() {
		os.RemoveAll(dir)
	}
}

func createNestedDir(baseDir string) error {
	return os.Mkdir(filepath.Join(baseDir, "nestedDirName"), 0444)
}

func checkError(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func createTempFiles(t *testing.T, dir string) {
	files := []string{"file with space", "f1.txt", "f2.pdf", "f3"}
	content := []byte("temporary file's content")

	for _, f := range files {
		fileName := filepath.Join(dir, f)
		if err := ioutil.WriteFile(fileName, content, 0444); err != nil {
			t.Error(err)
		}
	}
}

func replaceUserGroup(t *testing.T, raw string) string {
	var buf bytes.Buffer
	currentUser, err := user.Current()
	checkError(t, err)

	group, err := user.LookupGroupId(currentUser.Gid)
	checkError(t, err)

	data := struct {
		User  string
		Group string
	}{
		User:  currentUser.Username,
		Group: group.Name,
	}

	tpl, err := template.New("list").Parse(raw)
	if err := tpl.Execute(&buf, data); err != nil {
		t.Error(err)
	}

	return buf.String()
}

func TestListDirWithLongList(t *testing.T) {
	tempDir, teardown := setup(t)
	defer teardown()

	var buf bytes.Buffer
	expected := replaceUserGroup(
		t,
		"-r--r--r-- {{.User}} {{.Group}}     24 f1.txt\n"+
			"-r--r--r-- {{.User}} {{.Group}}     24 f2.pdf\n"+
			"-r--r--r-- {{.User}} {{.Group}}     24 f3\n"+
			"-r--r--r-- {{.User}} {{.Group}}     24 'file with space'\n",
	)

	files, _ := ioutil.ReadDir(tempDir)
	ls(files, &buf, printFileList)

	output := string(buf.Bytes())
	if output != expected {
		t.Errorf("got:\n'%v'\nexpected:\n'%v'\n", output, expected)
	}
}

func TestListDirFilesWithNewLineSeparator(t *testing.T) {
	tempDir, teardown := setup(t)
	defer teardown()
	createNestedDir(tempDir)

	var buf bytes.Buffer
	expected := "f1.txt\n" +
		"f2.pdf\n" +
		"f3\n" +
		"'file with space'\n" +
		"nestedDirName\n"

	files, _ := ioutil.ReadDir(tempDir)
	ls(files, &buf, printFileNames)

	output := string(buf.Bytes())
	if output != expected {
		t.Errorf("got:\n'%v'\nexpected:\n'%v'\n", output, expected)
	}
}

func TestListAFileWithLongList(t *testing.T) {
	tempDir, teardown := setup(t)
	defer teardown()
	files, _ := ioutil.ReadDir(tempDir)
	filepath := filepath.Join(tempDir, files[0].Name())

	var buf bytes.Buffer
	runls([]string{filepath}, &buf, printFileList)

	expected := replaceUserGroup(t, "-r--r--r-- {{.User}} {{.Group}}     24 f1.txt\n")

	output := string(buf.Bytes())
	if output != expected {
		t.Errorf("got:\n'%v'\nexpected:\n'%v'\n", output, expected)
	}
}

func TestListFilesWithNewLineSeparator(t *testing.T) {
	tempDir, teardown := setup(t)
	defer teardown()
	files, _ := ioutil.ReadDir(tempDir)

	var buf bytes.Buffer
	expected := "f2.pdf\n" +
		"f3\n"

	files, _ = ioutil.ReadDir(tempDir)
	ls(files[1:3], &buf, printFileNames)

	output := string(buf.Bytes())
	if output != expected {
		t.Errorf("got:\n'%v'\nexpected:\n'%v'\n", output, expected)
	}
}

func TestListFilesAsLongList(t *testing.T) {
	tempDir, teardown := setup(t)
	defer teardown()
	files, _ := ioutil.ReadDir(tempDir)

	var buf bytes.Buffer
	expected := replaceUserGroup(
		t,
		"-r--r--r-- {{.User}} {{.Group}}     24 f1.txt\n"+
			"-r--r--r-- {{.User}} {{.Group}}     24 f2.pdf\n"+
			"-r--r--r-- {{.User}} {{.Group}}     24 f3\n"+
			"-r--r--r-- {{.User}} {{.Group}}     24 'file with space'\n",
	)

	files, _ = ioutil.ReadDir(tempDir)
	ls(files, &buf, printFileList)

	output := string(buf.Bytes())
	if output != expected {
		t.Errorf("got:\n'%v'\nexpected:\n'%v'\n", output, expected)
	}
}

func TestEmptyStringOnEmptyDirectory(t *testing.T) {
	tempDir, teardown := setup(t)
	defer teardown()
	createNestedDir(tempDir)
	files, _ := ioutil.ReadDir(tempDir)
	filepath := filepath.Join(tempDir, files[4].Name())

	var buf bytes.Buffer
	runls([]string{filepath}, &buf, printFileNames)

	expected := replaceUserGroup(t, "")

	output := string(buf.Bytes())
	if output != expected {
		t.Errorf("got:\n'%v'\nexpected:\n'%v'\n", output, expected)
	}
}

func TestHumanizeBytes(t *testing.T) {
	var m = map[int64]string{
		0:             "0",
		1023:          "1023",
		1024:          "1.00K",
		1025:          "1.00K",
		1110000:       "1.06M",
		8073741824:    "7.52G",
		1099511627776: "1.00T",
	}

	for nbytes, expected := range m {
		if humanized := humanizeSize(nbytes); humanized != expected {
			t.Errorf("result [%v] not expected for [%v]", humanized, nbytes)
		}
	}
}
