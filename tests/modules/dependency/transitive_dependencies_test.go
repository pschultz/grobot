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

var _ = Describe("Install transitive dependencies", func() {
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

	It("should check the installed package for any bot.json files to installed transitive dependencies", func() {
		AssertFileDoesNotExist(dependency.LockFileName, fileSystem)
		vendorDir1 := "vendor/src/code.google.com/p/gomock"
		vendorDir2 := "vendor/src/code.google.com/foo/bar"
		AssertDirectoryDoesNotExist(vendorDir1, fileSystem)
		AssertPackageHasNoDependencies("code.google.com/foo/bar", fileSystem)

		dependencyBotConfig := []byte(`{
			"bot-version": "0.7",
			"dependency": {
				"packages": [
					{
						"name": "code.google.com/foo/bar",
						"type": "git",
						"version": "branch:master"
					}
				]
			}
		}`)

		gomock.InOrder(
			shell.EXPECT().Execute("git clone https://code.google.com/p/gomock "+vendorDir1, true),
			shell.EXPECT().SetWorkingDirectory(vendorDir1),
			shell.EXPECT().Execute("git rev-parse HEAD", true).Return("dependency1_version", nil),
			shell.EXPECT().SetWorkingDirectory(""),

			fileSystem.EXPECT().FileInfo(vendorDir1+"/bot.json").Return(&grobot.File{ExistingFile: true}, nil),
			fileSystem.EXPECT().ReadFile(vendorDir1+"/bot.json").Return(dependencyBotConfig, nil),

			shell.EXPECT().Execute("git clone https://code.google.com/foo/bar "+vendorDir2, true),
			shell.EXPECT().SetWorkingDirectory(vendorDir2),
			shell.EXPECT().Execute("git rev-parse HEAD", true).Return("dependency2_version", nil),
			shell.EXPECT().SetWorkingDirectory(""),
		)

		expectedLockFileContent := `{
			"packages": [
				{
					"name": "code.google.com/p/gomock",
					"source": {
						"type": "git",
						"version": "dependency1_version"
					}
				},
				{
					"name": "code.google.com/foo/bar",
					"source": {
						"type": "git",
						"version": "dependency2_version"
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

	It("should install packages only once", func() {
		AssertFileDoesNotExist(dependency.LockFileName, fileSystem)
		vendorDir1 := "vendor/src/code.google.com/p/gomock"
		vendorDir2 := "vendor/src/code.google.com/foo/bar"
		AssertDirectoryDoesNotExist(vendorDir1, fileSystem)
		AssertDirectoryDoesNotExist(vendorDir2, fileSystem)

		dependencyBotConfig1 := []byte(`{
			"bot-version": "0.7",
			"dependency": {
				"packages": [
					{
						"name": "code.google.com/foo/bar",
						"type": "git",
						"version": "branch:master"
					}
				]
			}
		}`)
		dependencyBotConfig2 := []byte(`{
			"bot-version": "0.7",
			"dependency": {
				"packages": [
					{
						"name": "code.google.com/p/gomock",
						"type": "git",
						"version": "branch:master"
					}
				]
			}
		}`)

		gomock.InOrder(
			shell.EXPECT().Execute("git clone https://code.google.com/p/gomock "+vendorDir1, true),
			shell.EXPECT().SetWorkingDirectory(vendorDir1),
			shell.EXPECT().Execute("git rev-parse HEAD", true).Return("dependency1_version", nil),
			shell.EXPECT().SetWorkingDirectory(""),

			fileSystem.EXPECT().FileInfo(vendorDir1+"/bot.json").Return(&grobot.File{ExistingFile: true}, nil),
			fileSystem.EXPECT().ReadFile(vendorDir1+"/bot.json").Return(dependencyBotConfig1, nil),

			shell.EXPECT().Execute("git clone https://code.google.com/foo/bar "+vendorDir2, true),
			shell.EXPECT().SetWorkingDirectory(vendorDir2),
			shell.EXPECT().Execute("git rev-parse HEAD", true).Return("dependency2_version", nil),
			shell.EXPECT().SetWorkingDirectory(""),

			fileSystem.EXPECT().FileInfo(vendorDir2+"/bot.json").Return(&grobot.File{ExistingFile: true}, nil),
			fileSystem.EXPECT().ReadFile(vendorDir2+"/bot.json").Return(dependencyBotConfig2, nil),
		)

		expectedLockFileContent := `{
			"packages": [
				{
					"name": "code.google.com/p/gomock",
					"source": {
						"type": "git",
						"version": "dependency1_version"
					}
				},
				{
					"name": "code.google.com/foo/bar",
					"source": {
						"type": "git",
						"version": "dependency2_version"
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
