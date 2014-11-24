package testAPI

import (
	"code.google.com/p/gomock/gomock"
	"github.com/fgrosse/grobot/tests/mocks"
)

func ExpectWorkingDirectoryIsReset(shell *mocks.MockShell) *gomock.Call {
	return shell.EXPECT().SetWorkingDirectory("")
}
