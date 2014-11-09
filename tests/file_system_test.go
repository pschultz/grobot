package tests

import (
	. "github.com/fgrosse/grobot/tests/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.google.com/p/gomock/gomock"
	"errors"
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

	It("should ask the FileSystemProvider if a file exists and for the modification date", func() {
		fileSystem.EXPECT().TargetInfo(somePath).Return(&grobot.Target{}, nil)
		_, err := grobot.TargetInfo(somePath)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should return the result of the shell provider", func() {
		expectedTargetInfo := &grobot.Target{
			ExistingFile:     true,
			IsDir:            true,
			ModificationTime: time.Now(),
		}
		fileSystem.EXPECT().TargetInfo(somePath).Return(expectedTargetInfo, nil)
		targetInfo, err := grobot.TargetInfo(somePath)
		Expect(err).NotTo(HaveOccurred())
		Expect(targetInfo).To(Equal(expectedTargetInfo))
	})

	It("should return errors from the shell provider", func() {
		expectedErr := errors.New("oh noes!!!")
		fileSystem.EXPECT().TargetInfo(somePath).Return(&grobot.Target{}, expectedErr)
		_, err := grobot.TargetInfo(somePath)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("Could not determine whether or not a file or folder exists : oh noes!!!"))
	})
})
