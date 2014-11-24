package dependency

import (
	. "github.com/fgrosse/grobot/tests/mocks"
	. "github.com/fgrosse/grobot/tests/testAPI"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.google.com/p/gomock/gomock"
	"github.com/fgrosse/grobot/modules/dependency"
)

var _ = Describe("update task", func() {
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

	It("should update all packages if no argument is given", func() {
		existingLockFileContent := `{
			"packages": [
				{
					"name": "code.google.com/p/foo",
					"source": {
						"type": "git",
						"version": "foo_version"
					}
				},
				{
					"name": "code.google.com/p/bar",
					"source": {
						"type": "git",
						"version": "bar_version"
					}
				}
			]
		}`
		AssertFileWithContentExists(dependency.LockFileName, existingLockFileContent, AnyTime, fileSystem)

		vendorDir1 := "vendor/src/code.google.com/p/foo"
		AssertPackageHasNoDependencies("code.google.com/p/foo", fileSystem)
		vendorDir2 := "vendor/src/code.google.com/p/bar"
		AssertPackageHasNoDependencies("code.google.com/p/bar", fileSystem)

		gomock.InOrder(
			shell.EXPECT().SetWorkingDirectory(vendorDir1),
			shell.EXPECT().Execute("git checkout master --quiet", true),
			shell.EXPECT().Execute("git pull", true),
			shell.EXPECT().Execute("git rev-parse HEAD", true).Return("foo_version_new", nil),
			shell.EXPECT().SetWorkingDirectory(""),

			shell.EXPECT().SetWorkingDirectory(vendorDir2),
			shell.EXPECT().Execute("git checkout master --quiet", true),
			shell.EXPECT().Execute("git pull", true),
			shell.EXPECT().Execute("git rev-parse HEAD", true).Return("bar_version_new", nil),
			shell.EXPECT().SetWorkingDirectory(""),
		)

		expectedLockFileContent := `{
			"packages": [
				{
					"name": "code.google.com/p/foo",
					"source": {
						"type": "git",
						"version": "foo_version_new"
					}
				},
				{
					"name": "code.google.com/p/bar",
					"source": {
						"type": "git",
						"version": "bar_version_new"
					}
				}
			]
		}`
		fileSystem.EXPECT().WriteFile(dependency.LockFileName, EqualJsonString(expectedLockFileContent))

		task := dependency.NewUpdateTask()
		_, err := task.Invoke("update")
		Expect(err).NotTo(HaveOccurred())
	})

	Context("Autocompletion", func() {
		It("should autocomplete package name if only one installed package matches the name", func() {
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
				shell.EXPECT().Execute("git rev-parse HEAD", true).Return("8798645651468468464654684684fff", nil),
				shell.EXPECT().SetWorkingDirectory(""),
			)

			task := dependency.NewUpdateTask()
			_, err := task.Invoke("update", "foo")
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
