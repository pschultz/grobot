package dependency

import "fmt"

type UpdateTask struct{}

func NewUpdateTask() *UpdateTask {
	return &UpdateTask{}
}

func (t *UpdateTask) Description() string {
	return "Update a given package to the newest version"
}

func (t *UpdateTask) Dependencies(string) []string {
	return []string{}
}

func (t *UpdateTask) Invoke(invokedName string, args ...string) (bool, error) {
	return false, fmt.Errorf("SORRY : Not yet implemented")
}
