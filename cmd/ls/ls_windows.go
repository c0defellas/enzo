// +build windows
package main

import (
	"os"
)

func lookupUser(fileinfo os.FileInfo) (string, error) {
	return "unknown", nil
}

func lookupGroup(fileinfo os.FileInfo) (string, error) {
	return "unknown", nil
}