package testAPI

import (
	"code.google.com/p/gomock/gomock"
	"github.com/fgrosse/grobot"
	"github.com/fgrosse/grobot/tests/mocks"
)

func AssertTask(name string, mockCtrl *gomock.Controller) *mocks.MockTask {
	task := mocks.NewMockTask(mockCtrl)
	grobot.RegisterTask(name, task)
	task.EXPECT().Invoke(name).Return(true, nil).AnyTimes()
	return task
}

func AssertLeafDependency(name string, mockCtrl *gomock.Controller) *mocks.MockTask {
	dep := AssertTask(name, mockCtrl)
	AssertNoDependencies(dep)
	return dep
}

func AssertDependencies(task *mocks.MockTask, args ...string) {
	task.EXPECT().Dependencies(gomock.Any()).Return(args).AnyTimes()
}

func AssertNoDependencies(task *mocks.MockTask) {
	task.EXPECT().Dependencies(gomock.Any()).Return([]string{}).AnyTimes()
}
