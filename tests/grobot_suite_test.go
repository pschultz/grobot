package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGrobot(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Grobot main test suite")
}
