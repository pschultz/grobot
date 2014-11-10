package testAPI

import (
	"code.google.com/p/gomock/gomock"
	"github.com/fgrosse/grobot"
	"github.com/fgrosse/grobot/log"
	. "github.com/fgrosse/grobot/tests/mocks"
)

func SetupTestEnvironment(mockCtrl *gomock.Controller) (shell *MockShell, fileSystem *MockFileSystem) {
	grobot.Reset()
	log.EnableDebug()
	shell = NewMockShell(mockCtrl)
	fileSystem = NewMockFileSystem(mockCtrl)
	grobot.ShellProvider = shell
	grobot.FileSystemProvider = fileSystem
	return shell, fileSystem
}

// AssertRegisteredTask registers a new mock task to the given name
func AssertRegisteredTask(name string, mockCtrl *gomock.Controller) *MockTask {
	task := NewMockTask(mockCtrl)
	grobot.RegisterTask(name, task)
	return task
}

// AssertTaskIsInvoked registers a new mock task to the given name and expects that it
// is invoked with the name any times.
// The mock task will return (true, nil) on any invocation
// @see AssertRegisteredTask
func AssertTaskIsInvoked(name string, mockCtrl *gomock.Controller) *MockTask {
	task := AssertRegisteredTask(name, mockCtrl)
	task.EXPECT().Invoke(name).Return(true, nil).AnyTimes()
	return task
}

func AssertLeafDependency(name string, mockCtrl *gomock.Controller) *MockTask {
	dep := AssertTaskIsInvoked(name, mockCtrl)
	AssertNoDependencies(dep)
	return dep
}

func AssertDependencies(task *MockTask, args ...string) {
	task.EXPECT().Dependencies(gomock.Any()).Return(args).AnyTimes()
}

func AssertNoDependencies(task *MockTask) {
	task.EXPECT().Dependencies(gomock.Any()).Return([]string{}).AnyTimes()
}
