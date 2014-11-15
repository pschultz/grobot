package tests

import (
	. "github.com/fgrosse/grobot/tests/mocks"
	. "github.com/fgrosse/grobot/tests/testAPI"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.google.com/p/gomock/gomock"
	"github.com/fgrosse/grobot"
)

var _ = Describe("Tasks", func() {
	var (
		mockCtrl   *gomock.Controller
		shell      *MockShell
		fileSystem *MockFileSystem
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		shell, fileSystem, _ = SetupTestEnvironment(mockCtrl)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("Hooking into other tasks", func() {
		BeforeEach(func() {
			AssertEmptyFileSystem(fileSystem)
		})

		Context("before parent invocation", func() {
			It("should execute all registered hooks in the right order", func() {
				mainTask := AssertRegisteredTask("main", mockCtrl)
				AssertNoDependencies(mainTask)

				subTask1 := AssertRegisteredTask("sub1", mockCtrl)
				AssertNoDependencies(subTask1)
				subTask2 := AssertRegisteredTask("sub2", mockCtrl)
				AssertNoDependencies(subTask2)

				grobot.RegisterTaskHook(grobot.HookBefore, "main", "sub1")
				grobot.RegisterTaskHook(grobot.HookBefore, "main", "sub2")

				gomock.InOrder(
					subTask1.EXPECT().Invoke("sub1"),
					subTask2.EXPECT().Invoke("sub2"),
					mainTask.EXPECT().Invoke("main"),
				)

				_, err := grobot.InvokeTask("main", 0)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("after parent invocation", func() {
			It("should execute the hooked task after the parent task", func() {
				mainTask := AssertRegisteredTask("main", mockCtrl)
				AssertNoDependencies(mainTask)

				subTask := AssertRegisteredTask("sub", mockCtrl)
				AssertNoDependencies(subTask)

				grobot.RegisterTaskHook(grobot.HookAfter, "main", "sub")

				gomock.InOrder(
					mainTask.EXPECT().Invoke("main"),
					subTask.EXPECT().Invoke("sub"),
				)

				_, err := grobot.InvokeTask("main", 0)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should execute all registered hooks in the right order", func() {
				mainTask := AssertRegisteredTask("main", mockCtrl)
				AssertNoDependencies(mainTask)

				subTask1 := AssertRegisteredTask("sub1", mockCtrl)
				AssertNoDependencies(subTask1)
				subTask2 := AssertRegisteredTask("sub2", mockCtrl)
				AssertNoDependencies(subTask2)

				grobot.RegisterTaskHook(grobot.HookAfter, "main", "sub1")
				grobot.RegisterTaskHook(grobot.HookAfter, "main", "sub2")

				gomock.InOrder(
					mainTask.EXPECT().Invoke("main"),
					subTask1.EXPECT().Invoke("sub1"),
					subTask2.EXPECT().Invoke("sub2"),
				)

				_, err := grobot.InvokeTask("main", 0)
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
