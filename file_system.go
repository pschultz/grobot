package grobot

import (
	"fmt"
	"os"
	"time"
)

func init() {
	FileSystemProvider = &RealFileSystem{}
}

type FileSystem interface {
	ModificationDate(path string) (bool, bool, time.Time, error)
}

var FileSystemProvider FileSystem

func ModificationDate(path string) (bool, bool, time.Time, error) {
	exists, isDir, modTime, err := FileSystemProvider.ModificationDate(path)
	if err != nil {
		err = fmt.Errorf("Could not determine whether or not a file or folder exists : %s", err.Error())
	}
	return exists, isDir, modTime, err
}

type RealFileSystem struct{}

func (f *RealFileSystem) ModificationDate(path string) (bool, bool, time.Time, error) {
	fileInfo, err := os.Stat(path)
	if err == nil {
		return true, fileInfo.IsDir(), fileInfo.ModTime(), nil
	}
	if os.IsNotExist(err) {
		return false, false, time.Time{}, nil
	}
	return false, false, time.Time{}, err
}
