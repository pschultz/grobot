package tests

import (
	. "github.com/fgrosse/gobot/tests/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.google.com/p/gomock/gomock"
	"errors"
	"github.com/fgrosse/gobot"
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
		gobot.FileSystemProvider = fileSystem
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should ask the FileSystemProvider if a file exists and for the modification date", func() {
		fileSystem.EXPECT().ModificationDate(somePath)
		gobot.ModificationDate(somePath)
	})

	It("should return the result of the shell provider", func() {
		expectedResult := true
		expectedIsDir := true
		expectedModTime := time.Now()
		fileSystem.EXPECT().ModificationDate(somePath).Return(expectedResult, expectedIsDir, expectedModTime, nil)
		exists, isDir, modTime, err := gobot.ModificationDate(somePath)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(Equal(expectedResult))
		Expect(isDir).To(Equal(expectedIsDir))
		Expect(modTime).To(Equal(expectedModTime))
	})

	It("should return errors from the shell provider", func() {
		expectedErr := errors.New("oh noes!!!")
		fileSystem.EXPECT().ModificationDate(somePath).Return(false, false, time.Time{}, expectedErr)
		_, _, _, err := gobot.ModificationDate(somePath)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("Could not determine whether or not a file or folder exists : oh noes!!!"))
	})
})
