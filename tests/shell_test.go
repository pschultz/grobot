package tests

import (
	. "github.com/fgrosse/grobot/tests/mocks"
	. "github.com/fgrosse/grobot/tests/testAPI"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.google.com/p/gomock/gomock"
	"errors"
	"github.com/fgrosse/grobot"
)

var _ = Describe("Shell", func() {
	var (
		mockCtrl *gomock.Controller
		shell    *MockShell
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		shell, _, _ = SetupTestEnvironment(mockCtrl)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should pass the program and arguments to the shell provider", func() {
		shell.EXPECT().Execute("foo barblup 3", false)
		grobot.Execute("foo %sblup %d", "bar", 3)
	})

	It("should return the trimmed output from the shell provider", func() {
		shell.EXPECT().Execute("foo barblup 3", false).Return("\n\n\ttest output   \n", nil)
		output := grobot.Execute("foo %sblup %d", "bar", 3)
		Expect(output).To(Equal("test output"))
	})

	It("should panic if the shell provider returned an error", func() {
		expectedErr := errors.New("hey there")
		shell.EXPECT().Execute("foo barblup 3", false).Return("", expectedErr)

		defer func() {
			r := recover()
			Expect(r).NotTo(BeNil())
			caughtErr := r.(error)
			Expect(caughtErr).To(Equal(expectedErr))
		}()

		grobot.Execute("foo %sblup %d", "bar", 3)
	})

	It("should support the silent mode", func() {
		shell.EXPECT().Execute("foo barblup 3", true)
		grobot.ExecuteSilent("foo %sblup %d", "bar", 3)
	})

	It("should not trigger the silent mode when debug is enabled", func() {
		grobot.EnableDebugMode()
		shell.EXPECT().Execute("foo barblup 3", false)
		grobot.ExecuteSilent("foo %sblup %d", "bar", 3)
	})
})
