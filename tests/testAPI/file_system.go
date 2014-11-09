package testAPI

import (
	"github.com/fgrosse/grobot/tests/mocks"
	"time"
)

func AssertFileExists(path string, modTime time.Time, fileSystem *mocks.MockFileSystem) {
	fileSystem.EXPECT().ModificationDate(path).Return(true, false, modTime, nil)
}

func AssertFileDoesNotExist(path string, fileSystem *mocks.MockFileSystem) {
	fileSystem.EXPECT().ModificationDate(path).Return(false, false, time.Time{}, nil)
}

func AssertNoFiles(fileSystem *mocks.MockFileSystem, args ...string) {
	for _, arg := range args {
		AssertFileDoesNotExist(arg, fileSystem)
	}
}
