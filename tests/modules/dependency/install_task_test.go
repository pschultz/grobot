package dependency

import (
	. "github.com/fgrosse/grobot/tests/mocks"
	. "github.com/fgrosse/grobot/tests/testAPI"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.google.com/p/gomock/gomock"
	"github.com/fgrosse/grobot/modules/dependency"
	"net/http"
)

var _ = Describe("Install tasks", func() {
	var (
		mockCtrl        *gomock.Controller
		shell           *MockShell
		fileSystem      *MockFileSystem
		httpClient      *MockHttpClient
		lockFileContent string
		vendorDir       = "vendor/src/foo.bar/fgrosse/test"
		cvsRev          = "7891f8646dc62f4e32642ba332bbe7cf0097d8c5"
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		shell, fileSystem, httpClient = SetupTestEnvironment(mockCtrl)
		lockFileContent = `{
			"packages": [
				{
					"name": "foo.bar/fgrosse/test",
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

		It("should use the go get logic to determine the repository url and checkout the correct version", func() {
			responseBody := `
				<!DOCTYPE html>
				<html>
					<head>
						<meta name='go-import' content='foo.bar/fgrosse/test git https://repository.foo.bar/fgrosse/test.git'>
					</head>
					<body/>
				</html>
			`
			req1, _ := http.NewRequest("GET", "https://foo.bar/fgrosse/test?go-get=1", nil)
			httpClient.EXPECT().Send(req1).Return(NewHttpResponse(responseBody), nil)

			gomock.InOrder(
				shell.EXPECT().Execute("git clone https://repository.foo.bar/fgrosse/test.git "+vendorDir, true),
				shell.EXPECT().SetWorkingDirectory(vendorDir),
				shell.EXPECT().Execute("git checkout "+cvsRev+" --quiet", true),
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
	})
})
