package gomock

import (
	"fmt"
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
	i := 0
	for mockSource, mockConf := range conf.Mocks {
		var mockFile string
		if mockConf.MockFileName != "" {
			mockFile = mockConf.MockFileName
		} else {
			mockFile = path.Base(mockSource)
			mockFile = fileEnding.ReplaceAllString(mockFile, "_mock.go")
		}

		task.mockFiles[i] = fmt.Sprintf("%s/%s", conf.MockFolder, mockFile)
		i++
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
	return true, nil
}
