package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGobot(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gobot main test suite")
}
