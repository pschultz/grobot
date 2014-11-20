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
	WriteFile(path string, data []byte) error
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

func FileExists(path string) bool {
	targetInfo := TargetInfo(path)
	return targetInfo.ExistingFile
}

func DirectoryExists(path string) bool {
	targetInfo := TargetInfo(path)
	return targetInfo.ExistingFile && targetInfo.IsDir
}

func TargetInfo(path string) *Target {
	targetInfo, err := FileSystemProvider.TargetInfo(path)
	if err != nil {
		panic(fmt.Errorf("Could not determine whether or not a file or folder exists : %s", err.Error()))

	}
	if targetInfo == nil {
		panic(fmt.Errorf("Internal error: FileSystemProvider must not return nil"))
	}
	targetInfo.Name = path
	return targetInfo
}

func ReadFile(path string) ([]byte, error) {
	data, err := FileSystemProvider.ReadFile(path)
	if err != nil {
		return data, fmt.Errorf(`Could not read file "%s" : %s`, path, err.Error())
	}
	return data, nil
}

func WriteFile(path string, data []byte) error {
	err := FileSystemProvider.WriteFile(path, data)
	if err != nil {
		return fmt.Errorf(`Error while writing file "%s" : %s`, path, err.Error())
	}
	return nil
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

func (f *RealFileSystem) WriteFile(path string, data []byte) error {
	return ioutil.WriteFile(path, data, 0644)
}
