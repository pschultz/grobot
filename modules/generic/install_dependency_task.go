package generic

import (
	"fmt"
	"github.com/fgrosse/gobot"
	"github.com/fgrosse/gobot/log"
)

type InstallDependencyTask struct {
	dependency string
}

func NewInstallDependencyTask(dependency string) *FolderTask {
	return NewFolderTask(&InstallDependencyTask{dependency})
}

func (t *InstallDependencyTask) Dependencies() []string {
	return []string{}
}

func (t *InstallDependencyTask) Invoke(path string) error {
	log.Print("The dependency %s is not installed in your GOPATH.", path)
	command := "go get " + path
	if log.AskBool("Do you want me to run %s ? [Yn] ", command) == false {
		return fmt.Errorf("User canceled task execution")
	}

	log.Action("Installing %s", path)
	gobot.Shell("echo world")
	return nil
}
