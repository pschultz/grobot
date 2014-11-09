package testAPI

import (
	"github.com/fgrosse/grobot/log"
	"github.com/fgrosse/grobot/tests/mocks"
	"time"
)

func AssertFileExists(path string, modTime time.Time, fileSystem *mocks.MockFileSystem) {
	log.Debug(">> Asserting file [%s] exist with modification date %s", path, modTime.Format("15:04:05"))
	fileSystem.EXPECT().ModificationDate(path).Return(true, false, modTime, nil).AnyTimes()
}

func AssertFileDoesNotExist(path string, fileSystem *mocks.MockFileSystem) {
	log.Debug(">> Asserting file [%s] does not exist", path)
	fileSystem.EXPECT().ModificationDate(path).Return(false, false, time.Time{}, nil).AnyTimes()
}

func AssertNoFiles(fileSystem *mocks.MockFileSystem, args ...string) {
	for _, arg := range args {
		AssertFileDoesNotExist(arg, fileSystem)
	}
}
