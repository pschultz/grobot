package golint

import (
	. "github.com/fgrosse/grobot/tests/mocks"
	. "github.com/fgrosse/grobot/tests/testAPI"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.google.com/p/gomock/gomock"
	"github.com/fgrosse/grobot"
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

	It("should run golint on each package", func() {
		filesInRoot := []*grobot.File{
			NewDirectory("package1"),
			NewDirectory("package2"),
			NewFile("root_file1.go"),
			NewFile("root_file2.go"),
			NewFile("root_file3.go"),
		}
		filesInPackage1 := []*grobot.File{
			NewDirectory("sub_package1a"),
			NewFile("p1_file1.go"),
			NewFile("p1_file2.go"),
		}
		filesInPackage2 := []*grobot.File{
			NewFile("p2_file1.go"),
			NewFile("p2_file2.go"),
		}
		filesInSubPackage := []*grobot.File{
			NewFile("p1a_filex.go"),
		}

		gomock.InOrder(
			fileSystem.EXPECT().WorkingDir().Return("/test/lint", nil),

			// first lint all files in root
			fileSystem.EXPECT().ListFiles("/test/lint").Return(filesInRoot, nil),
			shell.EXPECT().SetWorkingDirectory(""),
			shell.EXPECT().Execute(`golint "root_file1.go" "root_file2.go" "root_file3.go"`, false),

			// then descent in the first directory and lint all files on that level
			fileSystem.EXPECT().ListFiles("/test/lint/package1").Return(filesInPackage1, nil),
			shell.EXPECT().SetWorkingDirectory("package1"),
			shell.EXPECT().Execute(`golint "p1_file1.go" "p1_file2.go"`, false),

			// further descent into the next level
			fileSystem.EXPECT().ListFiles("/test/lint/package1/sub_package1a").Return(filesInSubPackage, nil),
			shell.EXPECT().SetWorkingDirectory("package1/sub_package1a"),
			shell.EXPECT().Execute(`golint "p1a_filex.go"`, false),

			// work on the next (and last) directory beneath the root
			fileSystem.EXPECT().ListFiles("/test/lint/package2").Return(filesInPackage2, nil),
			shell.EXPECT().SetWorkingDirectory("package2"),
			shell.EXPECT().Execute(`golint "p2_file1.go" "p2_file2.go"`, false),

			// reset the shell pwd
			shell.EXPECT().SetWorkingDirectory(""),
		)

		task := golint.NewLintTask()
		_, err := task.Invoke("lint")
		Expect(err).NotTo(HaveOccurred())
	})

	It("should run golint only on *.go files", func() {
		filesInRoot := []*grobot.File{
			NewFile("file1.go"),
			NewFile("file2.txt"),
		}
		gomock.InOrder(
			fileSystem.EXPECT().WorkingDir().Return("/test/lint", nil),
			fileSystem.EXPECT().ListFiles("/test/lint").Return(filesInRoot, nil),
			shell.EXPECT().SetWorkingDirectory(""),
			shell.EXPECT().Execute(`golint "file1.go"`, false),
			shell.EXPECT().SetWorkingDirectory(""),
		)

		task := golint.NewLintTask()
		_, err := task.Invoke("lint")
		Expect(err).NotTo(HaveOccurred())
	})

	It("should not run golint on hidden files", func() {
		filesInRoot := []*grobot.File{
			NewFile("file1.go"),
			NewFile(".file2.go"),
		}
		gomock.InOrder(
			fileSystem.EXPECT().WorkingDir().Return("/test/lint", nil),
			fileSystem.EXPECT().ListFiles("/test/lint").Return(filesInRoot, nil),
			shell.EXPECT().SetWorkingDirectory(""),
			shell.EXPECT().Execute(`golint "file1.go"`, false),
			shell.EXPECT().SetWorkingDirectory(""),
		)

		task := golint.NewLintTask()
		_, err := task.Invoke("lint")
		Expect(err).NotTo(HaveOccurred())
	})

	It("should not run golint on directories without go files", func() {
		filesInRoot := []*grobot.File{
			NewFile("README.md"),
			NewFile("file1.txt"),
		}
		gomock.InOrder(
			fileSystem.EXPECT().WorkingDir().Return("/test/lint", nil),
			fileSystem.EXPECT().ListFiles("/test/lint").Return(filesInRoot, nil),
			shell.EXPECT().SetWorkingDirectory(""),
		)

		task := golint.NewLintTask()
		_, err := task.Invoke("lint")
		Expect(err).NotTo(HaveOccurred())
	})

	It("should not run golint on the vendors dir", func() {
		filesInRoot := []*grobot.File{
			NewDirectory("vendor"),
		}
		gomock.InOrder(
			fileSystem.EXPECT().WorkingDir().Return("/test/lint", nil),
			fileSystem.EXPECT().ListFiles("/test/lint").Return(filesInRoot, nil),
			shell.EXPECT().SetWorkingDirectory(""),
		)

		task := golint.NewLintTask()
		_, err := task.Invoke("lint")
		Expect(err).NotTo(HaveOccurred())
	})

	// TODO: it should not lint files in test dir

	// TODO: check lint will respect settings from golint.json in each package if it exists

	// TODO: check lint is executed before release (if configured)
})
