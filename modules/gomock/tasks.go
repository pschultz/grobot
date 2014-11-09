package gomock

import (
	"fmt"
	"github.com/fgrosse/gobot"
	"github.com/fgrosse/gobot/log"
	"path"
	"regexp"
)

type AllMocksTask struct {
	mockFiles []string
}

func NewAllMocksTask(conf Configuration) *AllMocksTask {
	task := AllMocksTask{}
	task.mockFiles = make([]string, len(conf.Mocks))
	fileEnding, _ := regexp.Compile(`\.go$`)
	for i, mockSource := range conf.Mocks {
		mockFile := path.Base(mockSource)
		fileEnding.ReplaceAllString(mockFile, "_mock.go")
		task.mockFiles[i] = fmt.Sprintf("%s/%s", conf.MockFolder, mockFile)
	}
	return &task
}

func (t *AllMocksTask) Description() string {
	return "Build all mocks defined in the config file"
}

func (t *AllMocksTask) Dependencies(name string) []string {
	return t.mockFiles
}

func (t *AllMocksTask) Invoke(invokedName string) (bool, error) {
	return false, nil
}

type BuildMockFileTask struct {
	conf Configuration
}

func NewBuildMockFileTask(conf Configuration) *BuildMockFileTask {
	return &BuildMockFileTask{conf}
}

func (t *BuildMockFileTask) Dependencies(invokedName string) []string {
	source, err := t.getMockSourcePath(invokedName)
	if err != nil {
		log.Debug("Issue while resolving dependency for [%s] : %s", invokedName, err.Error())
		return []string{}
	}

	return []string{
		t.conf.MockFolder,
		"vendor/bin/mockgen",
		source,
	}
}

func (t *BuildMockFileTask) Invoke(invokedName string) (bool, error) {
	mockSource, err := t.getMockSourcePath(invokedName)
	if err != nil {
		return false, err
	}

	log.Action("Generating mock: %s from %s", invokedName, mockSource)
	mockGenBinary := "vendor/bin/mockgen"
	command := fmt.Sprintf(`%s -source "%s" -destination "%s"`, mockGenBinary, mockSource, invokedName)
	if t.conf.MockPackage != "" {
		command = fmt.Sprintf("%s -package %s", command, t.conf.MockPackage)
	}

	return true, gobot.Execute(command)
}

func (t *BuildMockFileTask) getMockSourcePath(invokedName string) (string, error) {
	fileEnding, _ := regexp.Compile(`_mock\.go$`)
	baseName := path.Base(invokedName)
	baseName = fileEnding.ReplaceAllString(baseName, ".go")

	for _, sourceName := range t.conf.Mocks {
		if path.Base(sourceName) == baseName {
			return sourceName, nil
		}
	}

	return "", fmt.Errorf("Could not determine source file of mock %s", invokedName)
}
