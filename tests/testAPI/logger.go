package testAPI

import (
	"bytes"
	"fmt"
	"github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	"github.com/onsi/ginkgo/types"
)

var buffer *bytes.Buffer

func getBuffer() *bytes.Buffer {
	if buffer == nil {
		buffer = &bytes.Buffer{}
	}

	return buffer
}

func Flush() {
	fmt.Println("\nLog for failed spec:")
	output := getBuffer().String()
	if len(output) == 0 {
		fmt.Println("No messages have been logged..")
	} else {
		fmt.Println(output)
	}
}

type TestLogger struct{}

func (l *TestLogger) Write(p []byte) (n int, err error) {
	return getBuffer().Write(p)
}

func LoggedFailHandler(message string, callerSkip ...int) {
	Flush()
	ginkgo.Fail(message, callerSkip...)
}

type LoggedReporter struct{}

func (r LoggedReporter) SpecWillRun(specSummary *types.SpecSummary) {
	getBuffer().Reset()
}

func (r LoggedReporter) SpecSuiteWillBegin(config config.GinkgoConfigType, summary *types.SuiteSummary) {
}

func (r LoggedReporter) BeforeSuiteDidRun(setupSummary *types.SetupSummary) {}

func (r LoggedReporter) SpecDidComplete(specSummary *types.SpecSummary) {
	if specSummary.Failed() {
		Flush()
	}
}

func (r LoggedReporter) AfterSuiteDidRun(setupSummary *types.SetupSummary) {}

func (r LoggedReporter) SpecSuiteDidEnd(summary *types.SuiteSummary) {}
