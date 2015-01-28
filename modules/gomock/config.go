package gomock

import "strings"

type Configuration struct {
	MockFolder  string                `json:"folder"`
	MockPackage string                `json:"package"`
	Mocks       map[string]MockConfig `json:"mocks"`
}

type MockConfig struct {
	Imports      string            `json:"imports"`
	MockFileName string            `json:"mock_file_name"`
	AuxFiles     map[string]string `json:"aux_files"`
}

func (c MockConfig) AuxFilesString() string {
	definitions := []string{}
	for pkg, source := range c.AuxFiles {
		definitions = append(definitions, pkg+"="+source)
	}

	return strings.Join(definitions, ",")
}
