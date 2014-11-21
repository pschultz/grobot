package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/fgrosse/grobot/log"
	"github.com/fgrosse/grobot/tests/testAPI"
	"testing"
)

func TestGrobot(t *testing.T) {
	log.SetOutputWriter(new(testAPI.TestLogger))
	RegisterFailHandler(testAPI.LoggedFailHandler)
	RunSpecsWithDefaultAndCustomReporters(t, "Grobot main test suite", []Reporter{&testAPI.LoggedReporter{}})
}
