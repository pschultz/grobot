package generic

import (
	"fmt"
	"github.com/fgrosse/grobot"
	"github.com/fgrosse/grobot/log"
)

type InstallDependencyTask struct {
	dependency string
}

func NewInstallDependencyTask(dependency string) *FolderTask {
	return NewFolderTask(&InstallDependencyTask{dependency})
}

func (t *InstallDependencyTask) Dependencies(invokedName string) []string {
	return []string{}
}

func (t *InstallDependencyTask) Invoke(path string) (bool, error) {
	path = stripVendorSource(path)
	log.Print("The dependency %s is not installed in your GOPATH.", path)
	command := "go get " + path
	if log.AskBool("Do you want me to run %s ?", command) == false {
		return false, fmt.Errorf("User canceled task execution")
	}

	log.Action("Installing %s", path)
	return true, grobot.Execute(command)
}
