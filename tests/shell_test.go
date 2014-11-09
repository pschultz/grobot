package tests

import (
	. "github.com/fgrosse/gobot/tests/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.google.com/p/gomock/gomock"
	"errors"
	"github.com/fgrosse/gobot"
)

var _ = Describe("Shell", func() {
	var (
		mockCtrl *gomock.Controller
		shell    *MockShell
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		shell = NewMockShell(mockCtrl)
		gobot.ShellProvider = shell
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should pass the program and arguments to the shell provider", func() {
		shell.EXPECT().Execute("foo barblup 3")
		gobot.Execute("foo %sblup %d", "bar", 3)
	})

	It("should return the result of the shell provider", func() {
		expectedErr := errors.New("hey there")
		shell.EXPECT().Execute("foo barblup 3").Return(expectedErr)
		err := gobot.Execute("foo %sblup %d", "bar", 3)
		Expect(err).To(Equal(expectedErr))
	})
})
