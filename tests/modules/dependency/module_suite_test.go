package dependency

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/fgrosse/grobot/log"
	"github.com/fgrosse/grobot/tests/testAPI"
	"testing"
)

func TestDependencyModule(t *testing.T) {
	log.SetOutputWriter(new(testAPI.TestLogger))
	RegisterFailHandler(testAPI.LoggedFailHandler)
	RunSpecsWithDefaultAndCustomReporters(t, "Dependency module test suite", []Reporter{&testAPI.LoggedReporter{}})
}
