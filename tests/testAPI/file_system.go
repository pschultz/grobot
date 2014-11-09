package testAPI

import (
	"github.com/fgrosse/grobot"
	"github.com/fgrosse/grobot/log"
	"github.com/fgrosse/grobot/tests/mocks"
	"time"
)

func AssertFileExists(path string, modTime time.Time, fileSystem *mocks.MockFileSystem) {
	log.Debug(">> Asserting file [%s] exist with modification date %s", path, modTime.Format("15:04:05"))
	targetInfo := grobot.Target{ExistingFile: true, IsDir: false, ModificationTime: modTime}
	fileSystem.EXPECT().TargetInfo(path).Return(&targetInfo, nil).AnyTimes()
}

func AssertFileDoesNotExist(path string, fileSystem *mocks.MockFileSystem) {
	log.Debug(">> Asserting file [%s] does not exist", path)
	targetInfo := grobot.Target{ExistingFile: false}
	fileSystem.EXPECT().TargetInfo(path).Return(&targetInfo, nil).AnyTimes()
}

func AssertNoFiles(fileSystem *mocks.MockFileSystem, args ...string) {
	for _, arg := range args {
		AssertFileDoesNotExist(arg, fileSystem)
	}
}
