package dependency

import (
	. "github.com/fgrosse/grobot/tests/mocks"
	. "github.com/fgrosse/grobot/tests/testAPI"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.google.com/p/gomock/gomock"
	"github.com/fgrosse/grobot/modules/dependency"
)

var _ = Describe("Install tasks (new package)", func() {
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

	Context("when lockfile is not yet existent", func() {
		It("should install a new package if given an additional argument", func() {
			AssertFileDoesNotExist(dependency.LockFileName, fileSystem)
			vendorDir := "vendor/src/code.google.com/p/gomock"
			AssertFileDoesNotExist(vendorDir, fileSystem)
			gomock.InOrder(
				shell.EXPECT().Execute("git clone https://code.google.com/p/gomock "+vendorDir, true),
				shell.EXPECT().SetWorkingDirectory(vendorDir),
				shell.EXPECT().Execute("git rev-parse HEAD", true).Return("7891f8646dc62f4e32642ba332bbe7cf0097d8c5", nil),
				shell.EXPECT().SetWorkingDirectory(""),
			)

			expectedLockFileContent := `
			{
				"packages": [
					{
						"name": "code.google.com/p/gomock",
						"source": {
							"type": "git",
							"version": "7891f8646dc62f4e32642ba332bbe7cf0097d8c5"
						}
					}
				]
			}`
			fileSystem.EXPECT().WriteFile(dependency.LockFileName, EqualJsonString(expectedLockFileContent))
			gomock.Any()
			task := dependency.NewInstallTask()
			_, err := task.Invoke("install", "code.google.com/p/gomock")
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("when lockfile does already exist", func() {
		It("should update the existing lockfile", func() {
			existingLockFileContent := `
			{
				"packages": [
					{
						"name": "code.google.com/p/foo",
						"source": {
							"type": "git",
							"version": "8798645651468468464654684684fff"
						}
					},
					{
						"name": "code.google.com/p/bar",
						"source": {
							"type": "git",
							"version": "9848468489a898d4984e984de984984"
						}
					}
				]
			}`
			AssertFileWithContentExists(dependency.LockFileName, existingLockFileContent, AnyTime, fileSystem)
			vendorDir := "vendor/src/code.google.com/p/gomock"
			AssertFileDoesNotExist(vendorDir, fileSystem)
			gomock.InOrder(
				shell.EXPECT().Execute("git clone https://code.google.com/p/gomock "+vendorDir, true),
				shell.EXPECT().SetWorkingDirectory(vendorDir),
				shell.EXPECT().Execute("git rev-parse HEAD", true).Return("7891f8646dc62f4e32642ba332bbe7cf0097d8c5", nil),
				shell.EXPECT().SetWorkingDirectory(""),
			)

			expectedLockFileContent := `
			{
				"packages": [
					{
						"name": "code.google.com/p/foo",
						"source": {
							"type": "git",
							"version": "8798645651468468464654684684fff"
						}
					},
					{
						"name": "code.google.com/p/bar",
						"source": {
							"type": "git",
							"version": "9848468489a898d4984e984de984984"
						}
					},
					{
						"name": "code.google.com/p/gomock",
						"source": {
							"type": "git",
							"version": "7891f8646dc62f4e32642ba332bbe7cf0097d8c5"
						}
					}
				]
			}`
			fileSystem.EXPECT().WriteFile(dependency.LockFileName, EqualJsonString(expectedLockFileContent))

			task := dependency.NewInstallTask()
			_, err := task.Invoke("install", "code.google.com/p/gomock")
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
