package gomock

import "fmt"

type MocksTask struct{}

func (t *MocksTask) Description() string {
	return "Build all mocks defined in the config file"
}

func (t *MocksTask) Dependencies() []string {
	return []string{
		"vendor/bin/mockgen",
	}
}

func (t *MocksTask) Invoke(name string) error {
	fmt.Println("Hello task world")
	return nil
}
