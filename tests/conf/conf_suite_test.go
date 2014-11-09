package conf

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGobotConf(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gobot conf test suite")
}
