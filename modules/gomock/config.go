package gomock

type Configuration struct {
	MockFolder  string   `json:"folder"`
	MockPackage string   `json:"package"`
	Mocks       []string `json:"mocks"`
}
