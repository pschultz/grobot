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
	FileInfo(path string) (*File, error)
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte) error
}

var FileSystemProvider FileSystem

type File struct {
	Name             string
	ExistingFile     bool
	IsDir            bool
	ModificationTime time.Time
}

func (t *File) targetExistsMessage() string {
	fileType := "File"
	if t.IsDir {
		fileType = "Folder"
	}
	return fmt.Sprintf("%s [<strong>%s</strong>] does already exist", fileType, t.Name)
}

func (t *File) Typ() string {
	if t.IsDir {
		return "folder"
	}
	return "file"
}

func FileExists(path string) bool {
	targetInfo := FileInfo(path)
	return targetInfo.ExistingFile
}

func DirectoryExists(path string) bool {
	targetInfo := FileInfo(path)
	return targetInfo.ExistingFile && targetInfo.IsDir
}

func FileInfo(path string) *File {
	targetInfo, err := FileSystemProvider.FileInfo(path)
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

func (f *RealFileSystem) FileInfo(path string) (*File, error) {
	fileInfo, err := os.Stat(path)
	if err == nil {
		return &File{ExistingFile: true, IsDir: fileInfo.IsDir(), ModificationTime: fileInfo.ModTime()}, nil
	}
	if os.IsNotExist(err) {
		return &File{ExistingFile: false}, nil
	}
	return &File{ExistingFile: false}, err
}

func (f *RealFileSystem) ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func (f *RealFileSystem) WriteFile(path string, data []byte) error {
	return ioutil.WriteFile(path, data, 0644)
}
