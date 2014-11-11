package grobot

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func init() {
	FileSystemProvider = &RealFileSystem{}
}

type FileSystem interface {
	TargetInfo(path string) (*Target, error)
	ReadFile(path string) ([]byte, error)
}

var FileSystemProvider FileSystem

type Target struct {
	Name             string
	ExistingFile     bool
	IsDir            bool
	ModificationTime time.Time
}

func (t *Target) targetExistsMessage() string {
	fileType := "File"
	if t.IsDir {
		fileType = "Folder"
	}
	return fmt.Sprintf("%s [<strong>%s</strong>] does already exist", fileType, t.Name)
}

func (t *Target) Typ() string {
	if t.IsDir {
		return "folder"
	}
	return "file"
}

func TargetInfo(path string) (*Target, error) {
	targetInfo, err := FileSystemProvider.TargetInfo(path)
	if err != nil {
		err = fmt.Errorf("Could not determine whether or not a file or folder exists : %s", err.Error())
	}
	if targetInfo == nil {
		return nil, fmt.Errorf("Internal error: FileSystemProvider must not return nil")
	}
	targetInfo.Name = path
	return targetInfo, err
}

func ReadFile(path string) ([]byte, error) {
	data, err := FileSystemProvider.ReadFile(path)
	if err != nil {
		return data, fmt.Errorf(`Could not read file "%s" : %s`, path, err.Error())
	}
	return data, nil
}

type RealFileSystem struct{}

func (f *RealFileSystem) TargetInfo(path string) (*Target, error) {
	fileInfo, err := os.Stat(path)
	if err == nil {
		return &Target{ExistingFile: true, IsDir: fileInfo.IsDir(), ModificationTime: fileInfo.ModTime()}, nil
	}
	if os.IsNotExist(err) {
		return &Target{ExistingFile: false}, nil
	}
	return &Target{ExistingFile: false}, err
}

func (f *RealFileSystem) ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}
