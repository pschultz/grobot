package tests

import (
	. "github.com/fgrosse/grobot/tests/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.google.com/p/gomock/gomock"
	"github.com/fgrosse/grobot"
)

var _ = Describe("Simple Tasks", func() {
	var (
		mockCtrl *gomock.Controller
		shell    *MockShell
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		shell = NewMockShell(mockCtrl)
		grobot.ShellProvider = shell
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("CreateDirectoryTask", func() {
		It("should have no dependencies", func() {
			task := &grobot.CreateDirectoryTask{}
			dependencies := task.Dependencies("anything")
			Expect(dependencies).To(BeEmpty())
		})

		It("should create the folder if invoked", func() {
			task := &grobot.CreateDirectoryTask{}
			shell.EXPECT().Execute(`mkdir -p "some/folder/bla"`).Return(nil)
			result, err := task.Invoke("some/folder/bla")
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(BeTrue())
		})

		It("should give a convenience method to regoister folders", func() {
			path := "some/folder/bla"
			grobot.RegisterDirectory(path)
			task, err := grobot.GetTask(path)
			Expect(err).NotTo(HaveOccurred())
			Expect(task).To(BeAssignableToTypeOf(&grobot.CreateDirectoryTask{}))
		})
	})
})
