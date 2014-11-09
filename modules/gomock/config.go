package gomock

type Configuration struct {
	MockFolder  string                `json:"folder"`
	MockPackage string                `json:"package"`
	Mocks       map[string]MockConfig `json:"mocks"`
}

type MockConfig struct {
	Imports string `json:"imports"`
}
