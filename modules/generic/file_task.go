package generic

import (
	"github.com/fgrosse/gobot"
	"github.com/fgrosse/gobot/log"
)

type FileTask struct {
	creationTask gobot.Task
}

type FolderTask struct {
	FileTask
}

func NewFileTask(createTask gobot.Task) *FileTask {
	return &FileTask{createTask}
}

func NewFolderTask(createTask gobot.Task) *FolderTask {
	return &FolderTask{FileTask{createTask}}
}

func NewCreateFolderTask() *FolderTask {
	return NewFolderTask(&CreateFolderTask{})
}

func (t *FileTask) Dependencies(invokedName string) []string {
	return t.creationTask.Dependencies(invokedName)
}

func (t *FileTask) Invoke(path string) (bool, error) {
	exists, _, _, err := gobot.ModificationDate(path)
	if err != nil {
		return false, err
	}

	if exists {
		log.Debug("Nothing to do: file or folder '%s' does already exist", path)
		return false, nil
	}

	return t.creationTask.Invoke(path)
}

type CreateFolderTask struct{}

func (t *CreateFolderTask) Dependencies(invokedName string) []string {
	return []string{}
}

func (t *CreateFolderTask) Invoke(path string) (bool, error) {
	return true, gobot.Execute(`mkdir -p "%s"`, path)
}
