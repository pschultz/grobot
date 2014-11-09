// Automatically generated by MockGen. DO NOT EDIT!
// Source: file_system.go

package mocks

import (
	gomock "code.google.com/p/gomock/gomock"
	. "github.com/fgrosse/grobot"
)

// Mock of FileSystem interface
type MockFileSystem struct {
	ctrl     *gomock.Controller
	recorder *_MockFileSystemRecorder
}

// Recorder for MockFileSystem (not exported)
type _MockFileSystemRecorder struct {
	mock *MockFileSystem
}

func NewMockFileSystem(ctrl *gomock.Controller) *MockFileSystem {
	mock := &MockFileSystem{ctrl: ctrl}
	mock.recorder = &_MockFileSystemRecorder{mock}
	return mock
}

func (_m *MockFileSystem) EXPECT() *_MockFileSystemRecorder {
	return _m.recorder
}

func (_m *MockFileSystem) TargetInfo(path string) (*Target, error) {
	ret := _m.ctrl.Call(_m, "TargetInfo", path)
	ret0, _ := ret[0].(*Target)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockFileSystemRecorder) TargetInfo(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "TargetInfo", arg0)
}
