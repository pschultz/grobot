package dependency

import (
	. "github.com/fgrosse/grobot/tests/mocks"
	. "github.com/fgrosse/grobot/tests/testAPI"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.google.com/p/gomock/gomock"
	"github.com/fgrosse/grobot"
	"github.com/fgrosse/grobot/modules/dependency"
)

var _ = Describe("Install tasks (new package)", func() {
	var (
		mockCtrl   *gomock.Controller
		shell      *MockShell
		fileSystem *MockFileSystem
		httpClient *MockHttpClient
		module     *dependency.Module
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		shell, fileSystem, httpClient = SetupTestEnvironment(mockCtrl)
		AssertFileWithContentExists(grobot.ConfigFileName, `{
			"bot-version": "0.7",
			"foo": {
				"bar": "blup"
			},
			"dependency": {
				"folder": "vendor"
			}
		}`, AnyTime, fileSystem)

		_, err := grobot.LoadConfigFromFile(grobot.ConfigFileName, grobot.NewVersion("0.7"))
		Expect(err).NotTo(HaveOccurred())
		module = grobot.GetModule("Depenency").(*dependency.Module)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("when lockfile is not yet existent", func() {
		It("should install a new package if given an additional argument", func() {
			AssertFileDoesNotExist(dependency.LockFileName, fileSystem)
			vendorDir := "vendor/src/code.google.com/p/gomock"
			AssertDirectoryDoesNotExist(vendorDir, fileSystem)
			AssertFileDoesNotExist(vendorDir+"/bot.json", fileSystem)
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
			fileSystem.EXPECT().WriteFile(grobot.ConfigFileName, gomock.Any()) // this is tested separately

			task := dependency.NewInstallTask(module)
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
			AssertPackageHasNoDependencies("code.google.com/p/gomock", fileSystem)

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
			fileSystem.EXPECT().WriteFile(grobot.ConfigFileName, gomock.Any()) // this is tested separately

			task := dependency.NewInstallTask(module)
			_, err := task.Invoke("install", "code.google.com/p/gomock")
			Expect(err).NotTo(HaveOccurred())
		})
	})

	It("should add the new dependencies to the dependency module in the config file", func() {
		AssertFileDoesNotExist(dependency.LockFileName, fileSystem)
		AssertPackageHasNoDependencies("code.google.com/p/gomock", fileSystem)
		shell.EXPECT().Execute(gomock.Any(), gomock.Any()).AnyTimes() // not tested here
		shell.EXPECT().SetWorkingDirectory(gomock.Any()).AnyTimes()   // not tested here

		expectedConfigFileContent := `
			{
				"bot-version": "0.7",
				"foo": {
					"bar": "blup"
				},
				"dependency": {
					"folder": "vendor",
					"packages": [
						{
							"name": "code.google.com/p/gomock",
							"type": "git",
							"version": "branch:master"
						}
					]
				}
			}`
		fileSystem.EXPECT().WriteFile(dependency.LockFileName, gomock.Any()) // not tested here
		fileSystem.EXPECT().WriteFile(grobot.ConfigFileName, EqualJsonString(expectedConfigFileContent))

		task := dependency.NewInstallTask(module)
		_, err := task.Invoke("install", "code.google.com/p/gomock")
		Expect(err).NotTo(HaveOccurred())
	})
})
