package webmgmt

import (
	"fmt"
	"os"
)

// FileExists will check the name passed in and return a bool if the file exists or not
func FileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

// Mkdir will attempt to create a directory and set the FileMode of the directory
func Mkdir(dirPath string, dirMode os.FileMode) error {
	err := os.MkdirAll(dirPath, dirMode)
	if err != nil {
		return fmt.Errorf("%s: making directory: %v", dirPath, err)
	}
	return nil
}
