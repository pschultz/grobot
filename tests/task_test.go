package tests

import (
	. "github.com/fgrosse/gobot/tests/mocks"
	. "github.com/fgrosse/gobot/tests/testAPI"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.google.com/p/gomock/gomock"
	"github.com/fgrosse/gobot"
	"github.com/fgrosse/gobot/log"
	"time"
)

var _ = Describe("Tasks", func() {
	var (
		mockCtrl   *gomock.Controller
		shell      *MockShell
		fileSystem *MockFileSystem
	)

	BeforeEach(func() {
		gobot.Reset()
		log.EnableDebug()
		mockCtrl = gomock.NewController(GinkgoT())
		shell = NewMockShell(mockCtrl)
		fileSystem = NewMockFileSystem(mockCtrl)

		gobot.ShellProvider = shell
		gobot.FileSystemProvider = fileSystem
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("GetTask", func() {
		Context("with no registered tasks", func() {
			It("should return an error", func() {
				task, err := gobot.GetTask("foo/bar")
				Expect(err).To(HaveOccurred())
				Expect(task).To(BeNil())
				Expect(err.Error()).To(Equal("Don't know how to build task 'foo/bar'"))
			})
		})

		Context("with registered tasks", func() {
			var (
				task1 = NewMockTask(mockCtrl)
				task2 = NewMockTask(mockCtrl)
			)
			BeforeEach(func() {
				gobot.RegisterTask("foo/bar", task1)
				gobot.RegisterTask("test", task2)
			})

			It("should return the registered task", func() {
				returnedTask, err := gobot.GetTask("foo/bar")
				Expect(err).NotTo(HaveOccurred())
				Expect(returnedTask).To(Equal(task1))

				returnedTask, err = gobot.GetTask("test")
				Expect(err).NotTo(HaveOccurred())
				Expect(returnedTask).To(Equal(task2))
			})

			It("should still return an error if an unkown task is requested", func() {
				returnedTask, err := gobot.GetTask("blub")
				Expect(err).To(HaveOccurred())
				Expect(returnedTask).To(BeNil())
				Expect(err.Error()).To(Equal("Don't know how to build task 'blub'"))
			})
		})

		Context("with registered rules", func() {
			var (
				task1 = NewMockTask(mockCtrl)
				task2 = NewMockTask(mockCtrl)
			)
			BeforeEach(func() {
				gobot.RegisterRule(`^foo/\w+\.go$`, task1)
				gobot.RegisterTask("foo/bar", task2)
			})

			It("should return the registered task", func() {
				returnedTask, err := gobot.GetTask("foo/bar.go")
				Expect(err).NotTo(HaveOccurred())
				Expect(returnedTask).To(Equal(task1))
			})
		})
	})

	Describe("InvokeTask", func() {
		It("should just invoke tasks without any dependencies", func() {
			AssertFileDoesNotExist("main", fileSystem)
			AssertLeafDependency("main", mockCtrl)
			_, err := gobot.InvokeTask("main", 0)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should invoke all dependencies of a task", func() {
			AssertNoFiles(fileSystem, "main", "dep1", "dep2")

			task := AssertTask("main", mockCtrl)
			AssertLeafDependency("dep1", mockCtrl)
			AssertLeafDependency("dep2", mockCtrl)

			AssertDependencies(task, "dep1", "dep2")
			_, err := gobot.InvokeTask("main", 0)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should invoke all dependencies of the dependencies", func() {
			AssertNoFiles(fileSystem, "main", "dep1", "dep2", "dep1/a", "dep1/b", "dep2/c", "dep2/d")

			task := AssertTask("main", mockCtrl)
			dep1 := AssertTask("dep1", mockCtrl)
			dep2 := AssertTask("dep2", mockCtrl)
			AssertLeafDependency("dep1/a", mockCtrl)
			AssertLeafDependency("dep1/b", mockCtrl)
			AssertLeafDependency("dep2/c", mockCtrl)
			AssertLeafDependency("dep2/d", mockCtrl)

			AssertDependencies(task, "dep1", "dep2")
			AssertDependencies(dep1, "dep1/a", "dep1/b")
			AssertDependencies(dep2, "dep2/c", "dep2/d")

			_, err := gobot.InvokeTask("main", 0)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should invoke each target only once", func() {
			AssertNoFiles(fileSystem, "main", "foo", "dep1", "dep2")

			task := AssertTask("main", mockCtrl)
			dep1 := AssertTask("dep1", mockCtrl)
			dep2 := AssertTask("dep2", mockCtrl)
			foo := NewMockTask(mockCtrl)
			gobot.RegisterTask("foo", foo)
			foo.EXPECT().Invoke("foo").Return(true, nil).Times(1)
			AssertNoDependencies(foo)

			AssertDependencies(task, "dep1", "dep2")
			AssertDependencies(dep1, "foo")
			AssertDependencies(dep2, "foo")

			_, err := gobot.InvokeTask("main", 0)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("target is not an existent file", func() {
			path := "foo/bar.go"
			BeforeEach(func() {
				AssertFileDoesNotExist(path, fileSystem)
			})

			It("should return an error if task has not been registered", func() {
				_, err := gobot.InvokeTask(path, 0)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Don't know how to build task 'foo/bar.go'"))
			})

			It("should invoke the task even though none of the dependencies returned (true, nil) on invoke", func() {
				AssertFileDoesNotExist("dep1", fileSystem)
				AssertFileDoesNotExist("dep2", fileSystem)

				task := NewMockTask(mockCtrl)
				AssertDependencies(task, "dep1", "dep2")
				gobot.RegisterTask(path, task)

				dep1 := NewMockTask(mockCtrl)
				dep2 := NewMockTask(mockCtrl)
				gobot.RegisterTask("dep1", dep1)
				gobot.RegisterTask("dep2", dep2)
				AssertNoDependencies(dep1)
				AssertNoDependencies(dep2)
				dep1.EXPECT().Invoke("dep1").Return(false, nil)
				dep2.EXPECT().Invoke("dep2").Return(false, nil)

				task.EXPECT().Invoke(path)
				_, err := gobot.InvokeTask(path, 0)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("target is an existent file", func() {
			path := "foo/bar.go"
			BeforeEach(func() {
				AssertFileExists(path, time.Now(), fileSystem)
			})

			It("should not return an error if task has not been registered", func() {
				result, err := gobot.InvokeTask(path, 0)
				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(BeFalse())
			})

			Context("with dependencies", func() {
				var (
					dep1 *MockTask
					dep2 *MockTask
				)
				BeforeEach(func() {
					AssertFileDoesNotExist("dep1", fileSystem)
					AssertFileDoesNotExist("dep2", fileSystem)

					dep1 = NewMockTask(mockCtrl)
					gobot.RegisterTask("dep1", dep1)
					AssertNoDependencies(dep1)

					dep2 = NewMockTask(mockCtrl)
					gobot.RegisterTask("dep2", dep2)
					AssertNoDependencies(dep2)
				})

				It("should invoke the task if any of the dependencies returns (true, nil) on invoke", func() {
					task := NewMockTask(mockCtrl)
					AssertDependencies(task, "dep1", "dep2")
					gobot.RegisterTask(path, task)

					dep1.EXPECT().Invoke("dep1").Return(false, nil)
					dep2.EXPECT().Invoke("dep2").Return(true, nil)
					task.EXPECT().Invoke(path).Return(true, nil)

					_, err := gobot.InvokeTask(path, 0)
					Expect(err).NotTo(HaveOccurred())
				})

				It("should not invoke the task if none of the dependencies returned (true, nil) on invoke", func() {
					task := NewMockTask(mockCtrl)
					AssertDependencies(task, "dep1", "dep2")
					gobot.RegisterTask(path, task)

					dep1.EXPECT().Invoke("dep1").Return(false, nil)
					dep2.EXPECT().Invoke("dep2").Return(false, nil)

					_, err := gobot.InvokeTask(path, 0)
					Expect(err).NotTo(HaveOccurred())
				})
			})
		})
	})
})
