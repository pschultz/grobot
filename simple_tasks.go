package grobot

import "github.com/fgrosse/grobot/log"

type CreateDirectoryTask struct{}

func (t *CreateDirectoryTask) Dependencies(string) []string {
	return []string{}
}

func (t *CreateDirectoryTask) Invoke(path string) (bool, error) {
	log.Action("Creating directory %s", path)
	Execute(`mkdir -p "%s"`, path)
	return true, nil
}
