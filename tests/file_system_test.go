package tests

import (
	. "github.com/fgrosse/grobot/tests/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.google.com/p/gomock/gomock"
	"errors"
	"fmt"
	"github.com/fgrosse/grobot"
	"time"
)

var _ = Describe("FileSystem", func() {
	var (
		mockCtrl   *gomock.Controller
		fileSystem *MockFileSystem
		somePath   = "foo/bar/baz.go"
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		fileSystem = NewMockFileSystem(mockCtrl)
		grobot.FileSystemProvider = fileSystem
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("FileInfo", func() {
		It("should ask the FileSystemProvider if a file exists and for the modification date", func() {
			fileSystem.EXPECT().FileInfo(somePath).Return(&grobot.File{}, nil)
			grobot.FileInfo(somePath)
		})

		It("should return the result from the file system provider", func() {
			expectedFileInfo := &grobot.File{
				ExistingFile:     true,
				IsDir:            true,
				ModificationTime: time.Now(),
			}
			fileSystem.EXPECT().FileInfo(somePath).Return(expectedFileInfo, nil)
			targetInfo := grobot.FileInfo(somePath)
			Expect(targetInfo).To(Equal(expectedFileInfo))
		})

		It("should panic if the FileSystemProvider returns any errors", func() {
			expectedErr := errors.New("oh noes!!!")
			fileSystem.EXPECT().FileInfo(somePath).Return(&grobot.File{}, expectedErr)

			defer func() {
				r := recover()
				Expect(r).NotTo(BeNil())
				caughtErr := r.(error)
				Expect(caughtErr.Error()).To(Equal("Could not determine whether or not a file or folder exists : oh noes!!!"))
			}()

			grobot.FileInfo(somePath)
		})
	})

	Describe("ReadFile", func() {
		It("should use the FileSystemProvider to read files", func() {
			expectedContent := []byte("123456")
			fileSystem.EXPECT().ReadFile(somePath).Return(expectedContent, nil)
			returnedContent, err := grobot.ReadFile(somePath)
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedContent).To(Equal(expectedContent))
		})

		It("should return errors from the shell provider", func() {
			expectedErr := errors.New("oh noes!!!")
			fileSystem.EXPECT().ReadFile(somePath).Return([]byte("123456"), expectedErr)
			_, err := grobot.ReadFile(somePath)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(fmt.Sprintf(`Could not read file "%s" : oh noes!!!`, somePath)))
		})
	})
})
