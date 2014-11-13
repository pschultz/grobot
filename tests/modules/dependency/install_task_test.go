package dependency

import (
	. "github.com/fgrosse/grobot/tests/mocks"
	. "github.com/fgrosse/grobot/tests/testAPI"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.google.com/p/gomock/gomock"
	"fmt"
	"github.com/fgrosse/grobot/modules/dependency"
)

/*
 * TODO: for this package
 *       - support HTTP (not only HTTPS)
 */

var _ = Describe("Install tasks", func() {
	var (
		mockCtrl        *gomock.Controller
		shell           *MockShell
		fileSystem      *MockFileSystem
		lockFileContent string
		vendorDir       = "vendor/src/github.com/onsi/ginkgo"
		cvsRev          = "7891f8646dc62f4e32642ba332bbe7cf0097d8c5"
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		shell, fileSystem = SetupTestEnvironment(mockCtrl)
		lockFileContent = `{
			"packages": [
				{
					"name": "github.com/onsi/ginkgo",
					"source": {
						"type": "git",
						"reference": "` + cvsRev + `"
					}
				}
			]
		}`
		AssertFileWithContentExists(dependency.LockFileName, lockFileContent, AnyTime, fileSystem)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("when repository is not existing", func() {
		BeforeEach(func() {
			AssertFileDoesNotExist(vendorDir, fileSystem)
		})

		It("should install dependencies from "+dependency.LockFileName, func() {
			gomock.InOrder(
				shell.EXPECT().Execute("git clone https://github.com/onsi/ginkgo "+vendorDir, false),
				shell.EXPECT().SetWorkingDirectory(vendorDir),
				shell.EXPECT().Execute("git checkout "+cvsRev+" --quiet", false),
				shell.EXPECT().SetWorkingDirectory(""),
			)

			task := dependency.NewInstallTask()
			_, err := task.Invoke("install")
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("when repository is already existing", func() {
		BeforeEach(func() {
			AssertFileExists(vendorDir, AnyTime, fileSystem)
		})

		It("should not do anything if the checkout version equals the requested version", func() {
			gomock.InOrder(
				shell.EXPECT().SetWorkingDirectory(vendorDir),
				shell.EXPECT().Execute("git rev-parse HEAD", true).Return(cvsRev, nil),
				shell.EXPECT().SetWorkingDirectory(""),
			)
			task := dependency.NewInstallTask()
			_, err := task.Invoke("install")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should report an error if the checkout version does not equal the requested version", func() {
			gomock.InOrder(
				shell.EXPECT().SetWorkingDirectory(vendorDir),
				shell.EXPECT().Execute("git rev-parse HEAD", true).Return("123456", nil),
				shell.EXPECT().SetWorkingDirectory(""),
			)
			task := dependency.NewInstallTask()
			_, err := task.Invoke("install")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(fmt.Sprintf("Repository at %s is not at the required version %s", vendorDir, cvsRev)))
		})
	})
})
