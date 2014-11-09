package generic

import (
	"github.com/fgrosse/grobot"
	"github.com/fgrosse/grobot/log"
)

type FileTask struct {
	creationTask grobot.Task
}

type FolderTask struct {
	FileTask
}

func NewFileTask(createTask grobot.Task) *FileTask {
	return &FileTask{createTask}
}

func NewFolderTask(createTask grobot.Task) *FolderTask {
	return &FolderTask{FileTask{createTask}}
}

func NewCreateFolderTask() *FolderTask {
	return NewFolderTask(&CreateFolderTask{})
}

func (t *FileTask) Dependencies(invokedName string) []string {
	return t.creationTask.Dependencies(invokedName)
}

func (t *FileTask) Invoke(path string) (bool, error) {
	targetInfo, err := grobot.TargetInfo(path)
	if err != nil {
		return false, err
	}

	if targetInfo.ExistingFile {
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
	return true, grobot.Execute(`mkdir -p "%s"`, path)
}
