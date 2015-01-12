package testAPI

import (
	"code.google.com/p/gomock/gomock"
	"github.com/fgrosse/grobot"
	"github.com/fgrosse/grobot/log"
	"github.com/fgrosse/grobot/tests/mocks"
	"time"
)

var (
	// VeryOld is the current time two years ago
	VeryOld = time.Now().Add(-356 * 24 * 2 * time.Hour)

	// AnyTime is any time in the recent past (when doesn't actually matter too much)
	AnyTime = time.Now().Add(-2 * time.Hour)

	// Now is the current time
	Now = time.Now()
)

func AssertFileExists(path string, modTime time.Time, fileSystem *mocks.MockFileSystem) {
	log.Debug(">> Asserting file [%s] exist with modification date %s", path, modTime.Format("15:04:05"))
	fileInfo := grobot.File{ExistingFile: true, IsDir: false, ModificationTime: modTime}
	fileSystem.EXPECT().FileInfo(path).Return(&fileInfo, nil).AnyTimes()
}

func AssertFileDoesNotExist(path string, fileSystem *mocks.MockFileSystem) {
	log.Debug(">> Asserting file [%s] does not exist", path)
	fileInfo := grobot.File{ExistingFile: false}
	fileSystem.EXPECT().FileInfo(path).Return(&fileInfo, nil).AnyTimes()
}

func AssertDirectoryDoesNotExist(path string, fileSystem *mocks.MockFileSystem) {
	log.Debug(">> Asserting directory [%s] does not exist", path)
	fileInfo := grobot.File{ExistingFile: false}
	fileSystem.EXPECT().FileInfo(path).Return(&fileInfo, nil).AnyTimes()
}

func AssertNoFiles(fileSystem *mocks.MockFileSystem, args ...string) {
	for _, arg := range args {
		AssertFileDoesNotExist(arg, fileSystem)
	}
}

func AssertEmptyFileSystem(fileSystem *mocks.MockFileSystem) {
	fileInfo := grobot.File{ExistingFile: false}
	fileSystem.EXPECT().FileInfo(gomock.Any()).Return(&fileInfo, nil).AnyTimes()
}

func AssertFileWithContentExists(path, content string, modTime time.Time, fileSystem *mocks.MockFileSystem) {
	log.Debug(">> Asserting file [%s] with some content exist with modification date %s", path, modTime.Format("15:04:05"))
	fileInfo := grobot.File{ExistingFile: true, IsDir: false, ModificationTime: modTime}
	fileSystem.EXPECT().FileInfo(path).Return(&fileInfo, nil).AnyTimes()
	fileSystem.EXPECT().ReadFile(path).Return([]byte(content), nil)
}

func AssertPackageHasNoDependencies(packageName string, fileSystem *mocks.MockFileSystem) {
	vendorDir := "vendor/src/" + packageName
	AssertDirectoryDoesNotExist(vendorDir, fileSystem)
	AssertFileDoesNotExist(vendorDir+"/"+grobot.ConfigFileName, fileSystem)
}
