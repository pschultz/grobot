package grobot

import "github.com/fgrosse/grobot/log"

type CreateFolderTask struct{}

func (t *CreateFolderTask) Dependencies(string) []string {
	return []string{}
}

func (t *CreateFolderTask) Invoke(path string) (bool, error) {
	log.Action("Creating folder %s", path)
	return true, Execute(`mkdir -p "%s"`, path)
}
