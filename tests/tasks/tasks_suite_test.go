package tasks

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGobotTasks(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gobot tasks test suite")
}
