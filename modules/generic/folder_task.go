package generic

import (
	"fmt"
	"github.com/fgrosse/gobot"
	"github.com/fgrosse/gobot/log"
	"os"
)

type FolderTask struct {
	createTask   gobot.Task
	dependencies []string
}

func NewFolderTask(createTask gobot.Task) *FolderTask {
	return &FolderTask{createTask, []string{}}
}

func (t *FolderTask) Dependencies() []string {
	return t.dependencies
}

func (t *FolderTask) Invoke(path string) error {
	exists, err := pathExists(path)
	if err != nil {
		return err
	}

	if exists {
		log.Debug("Folder %s does already exist", path)
		return nil
	}

	return t.createTask.Invoke(path)
}

// pathExists returns whether the given file or directory exists or not
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, fmt.Errorf("Could not determine whether or not a file or folder exists : %s", err.Error())
}
