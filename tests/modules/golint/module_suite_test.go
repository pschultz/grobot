package golint

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/fgrosse/grobot/log"
	"github.com/fgrosse/grobot/tests/testAPI"
	"testing"
)

func TestGolintModule(t *testing.T) {
	log.SetOutputWriter(new(testAPI.TestLogger))
	RegisterFailHandler(testAPI.LoggedFailHandler)
	RunSpecsWithDefaultAndCustomReporters(t, "golint module    ", []Reporter{&testAPI.LoggedReporter{}})
}
