package dependency

import (
	. "github.com/fgrosse/grobot/tests/mocks"
	. "github.com/fgrosse/grobot/tests/testAPI"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.google.com/p/gomock/gomock"
	"github.com/fgrosse/grobot/modules/dependency"
)

var _ = Describe("Update tasks", func() {
	var (
		mockCtrl   *gomock.Controller
		shell      *MockShell
		fileSystem *MockFileSystem
		httpClient *MockHttpClient
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		shell, fileSystem, httpClient = SetupTestEnvironment(mockCtrl)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should update a given package to the most recent version on master", func() {
		existingLockFileContent := `{
			"packages": [
				{
					"name": "code.google.com/p/foo",
					"source": {
						"type": "git",
						"version": "8798645651468468464654684684fff"
					}
				}
			]
		}`
		AssertFileWithContentExists(dependency.LockFileName, existingLockFileContent, AnyTime, fileSystem)

		vendorDir := "vendor/src/code.google.com/p/foo"
		AssertFileExists(vendorDir, AnyTime, fileSystem)

		gomock.InOrder(
			shell.EXPECT().SetWorkingDirectory(vendorDir),
			shell.EXPECT().Execute("git checkout master --quiet", true),
			shell.EXPECT().Execute("git pull", true),
			shell.EXPECT().Execute("git rev-parse HEAD", true).Return("7891f8646dc62f4e32642ba332bbe7cf0097d8c5", nil),
			shell.EXPECT().SetWorkingDirectory(""),
		)

		expectedLockFileContent := `{
			"packages": [
				{
					"name": "code.google.com/p/foo",
					"source": {
						"type": "git",
						"version": "7891f8646dc62f4e32642ba332bbe7cf0097d8c5"
					}
				}
			]
		}`
		fileSystem.EXPECT().WriteFile(dependency.LockFileName, EqualJsonString(expectedLockFileContent))

		task := dependency.NewUpdateTask()
		_, err := task.Invoke("update", "code.google.com/p/foo")
		Expect(err).NotTo(HaveOccurred())
	})
})
