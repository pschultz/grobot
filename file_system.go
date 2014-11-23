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
	ListFiles(path string) ([]*File, error)
	WorkingDir() (string, error)
}

var FileSystemProvider FileSystem

type File struct {
	Name             string
	ExistingFile     bool
	IsDir            bool
	ModificationTime time.Time
}

func (t *File) targetExistsMessage() string {
	return fmt.Sprintf("%s [<strong>%s</strong>] does already exist", t.Typ(), t.Name)
}

func (t *File) Typ() string {
	if t.IsDir {
		return "directory"
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
		err = fmt.Errorf(`Could not read file "%s" : %s`, path, err.Error())
	}
	return data, err
}

func WriteFile(path string, data []byte) error {
	err := FileSystemProvider.WriteFile(path, data)
	if err != nil {
		return fmt.Errorf(`Error while writing file "%s" : %s`, path, err.Error())
	}
	return nil
}

func ListFiles(path string) []*File {
	files, err := FileSystemProvider.ListFiles(path)
	if err != nil {
		panic(fmt.Errorf(`Could not list file of directory "%s" : %s`, path, err.Error()))
	}
	return files
}

func WorkingDir() string {
	pwd, err := FileSystemProvider.WorkingDir()
	if err != nil {
		panic(fmt.Errorf(`Could not get current working directory : %s`, err.Error()))
	}
	return pwd
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

func (f *RealFileSystem) ListFiles(path string) ([]*File, error) {
	directoryEntries, err := ioutil.ReadDir(path)
	if err != nil {
		return []*File{}, err
	}

	var files []*File
	for _, fileInfo := range directoryEntries {
		file := fileFromOsFileInfo(fileInfo)
		files = append(files, file)
	}

	return files, nil
}

func fileFromOsFileInfo(fileInfo os.FileInfo) *File {
	return &File{Name: fileInfo.Name(), ExistingFile: true, IsDir: fileInfo.IsDir(), ModificationTime: fileInfo.ModTime()}
}

func (f *RealFileSystem) WorkingDir() (string, error) {
	return os.Getwd()
}
