package dependency

import (
	. "github.com/fgrosse/grobot/tests/mocks"
	. "github.com/fgrosse/grobot/tests/testAPI"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.google.com/p/gomock/gomock"
	"github.com/fgrosse/grobot/modules/dependency"
)

/*
 * TODO: for this package
 *       - support HTTP (not only HTTPS)
 */

var _ = Describe("Tasks", func() {
	var (
		mockCtrl   *gomock.Controller
		shell      *MockShell
		fileSystem *MockFileSystem
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		shell, fileSystem = SetupTestEnvironment(mockCtrl)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("install task", func() {
		It("should install dependencies from "+dependency.LockFileName, func() {
			lockFileContent := `{
				"packages": [
					{
						"name": "github.com/onsi/ginkgo",
						"source": {
							"type": "git",
							"reference": "7891f8646dc62f4e32642ba332bbe7cf0097d8c5"
						}
					}
				]
			}`
			AssertFileWithContentExists(dependency.LockFileName, lockFileContent, AnyTime, fileSystem)

			gomock.InOrder(
				shell.EXPECT().Execute("git clone https://github.com/onsi/ginkgo vendor/src/github.com/onsi/ginkgo"),
				shell.EXPECT().Execute("cd vendor/src/github.com/onsi/ginkgo && git checkout 7891f8646dc62f4e32642ba332bbe7cf0097d8c5"),
			)

			task := dependency.NewInstallTask()
			_, err := task.Invoke("install")
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
