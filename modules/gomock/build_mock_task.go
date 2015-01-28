package gomock

import (
	"fmt"
	"github.com/fgrosse/grobot"
	"github.com/fgrosse/grobot/log"
	"path"
	"regexp"
)

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

func (t *BuildMockFileTask) Invoke(targetName string, args ...string) (bool, error) {
	mockSource, err := t.getMockSourcePath(targetName)
	if err != nil {
		return false, err
	}

	mockConfig, isSet := t.conf.Mocks[mockSource]
	if isSet == false {
		return false, fmt.Errorf("Logic error in BuildMockFileTask: could not find mock configuration for mock source %s (target [%s])", mockSource, targetName)
	}

	log.Action("Generating mock: %s from %s", targetName, mockSource)
	mockGenBinary := "vendor/bin/mockgen"
	command := fmt.Sprintf(`%s -source "%s" -destination "%s"`, mockGenBinary, mockSource, targetName)
	if t.conf.MockPackage != "" {
		command = fmt.Sprintf("%s -package %s", command, t.conf.MockPackage)
	}

	if mockConfig.Imports != "" {
		command = fmt.Sprintf(`%s -imports "%s"`, command, mockConfig.Imports)
	}

	if len(mockConfig.AuxFiles) > 0 {
		command = fmt.Sprintf(`%s -aux_files=%s"`, command, mockConfig.AuxFilesString())
	}

	grobot.Execute(command)
	return true, nil
}

func (t *BuildMockFileTask) getMockSourcePath(invokedName string) (string, error) {
	fileEnding, _ := regexp.Compile(`_mock\.go$`)

	baseName := path.Base(invokedName)
	baseNameWithoutMockExtension := fileEnding.ReplaceAllString(baseName, ".go")

	for sourceName, mockConf := range t.conf.Mocks {
		if mockConf.MockFileName == baseName {
			return sourceName, nil
		}
		if path.Base(sourceName) == baseNameWithoutMockExtension {
			return sourceName, nil
		}
	}

	return "", fmt.Errorf("Could not determine source file of mock %s", invokedName)
}
