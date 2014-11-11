package ginkgo

import (
	"fmt"
	"github.com/fgrosse/grobot"
)

type TestTask struct {
	conf Configuration
}

func NewTestTask(conf Configuration) *TestTask {
	return &TestTask{conf}
}

func (t *TestTask) Description() string {
	return "Run all ginkgo tests"
}

func (t *TestTask) Dependencies(name string) []string {
	return []string{"vendor/bin/ginkgo"}
}

func (t *TestTask) Invoke(invokedName string) (bool, error) {
	command := "ginkgo -r"
	if t.conf.TestFolder != "" {
		command = fmt.Sprintf(`%s "%s"`, command, t.conf.TestFolder)
	}

	err := grobot.Execute(command)
	return true, err
}
