package webmgmt

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
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

// WriteNewFile will attempt to write a file with the filename and path, a Reader and the FileMode of the file to be created.
// If an error is encountered an error will be returned.
func WriteNewFile(fpath string, in io.Reader) error {
	err := os.MkdirAll(filepath.Dir(fpath), 0775)
	if err != nil {
		return fmt.Errorf("%s: making directory for file: %v", fpath, err)
	}

	out, err := os.Create(fpath)
	if err != nil {
		return fmt.Errorf("%s: creating new file: %v", fpath, err)
	}
	defer out.Close()

	//fmt.Printf("WriteNewFile: %v\n", fm)
	//
	//err = out.Chmod(fm)
	//if err != nil && runtime.GOOS != "windows" {
	//	return fmt.Errorf("%s: changing file mode: %v", fpath, err)
	//}

	_, err = io.Copy(out, in)
	if err != nil {
		return fmt.Errorf("%s: writing file: %v", fpath, err)
	}
	return nil
}
