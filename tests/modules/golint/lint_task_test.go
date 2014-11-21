package golint

import (
	. "github.com/fgrosse/grobot/tests/mocks"
	. "github.com/fgrosse/grobot/tests/testAPI"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.google.com/p/gomock/gomock"
	"github.com/fgrosse/grobot/modules/golint"
)

var _ = Describe("Codestyle tasks", func() {
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

	It("should run golint", func() {
		// TODO: should lint all files in this package and all subpackages except vendor dir and test dir

		shell.EXPECT().Execute("golint", true)
		task := golint.NewLintTask()
		_, err := task.Invoke("lint")
		Expect(err).NotTo(HaveOccurred())
	})

	// TODO: check lint will respect settings from golint.json in each package if it exists

	// TODO: check lint is executed before release (if configured)
})
