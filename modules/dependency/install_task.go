package dependency

import (
	"encoding/json"
	"fmt"
	"github.com/fgrosse/grobot"
	"github.com/fgrosse/grobot/log"
)

type InstallTask struct{}

func NewInstallTask() *InstallTask {
	return &InstallTask{}
}

func (t *InstallTask) Description() string {
	return "Install all dependencies that are documented in " + LockFileName
}

func (t *InstallTask) Dependencies(string) []string {
	return []string{}
}

func (t *InstallTask) Invoke(string) (bool, error) {
	lockFile, err := loadLockFile()
	if err != nil {
		return false, err
	}

	return installDependencies(lockFile)
}

func loadLockFile() (*LockFile, error) {
	targetInfo, err := grobot.TargetInfo(LockFileName)
	if err != nil {
		return nil, fmt.Errorf("Error while reading meta data of dependency lock file: %s", err.Error())
	}

	if targetInfo.ExistingFile == false {
		log.Print("Lock file [<strong>%s</strong>] does not exist", LockFileName)
		return nil, nil
	}

	log.Debug("Reading dependency lock file [<strong>%s</strong>]", LockFileName)
	data, err := grobot.ReadFile(LockFileName)
	if err != nil {
		return nil, fmt.Errorf("Error while reading dependency lock file: %s", err.Error())
	}

	var lockFile LockFile
	err = json.Unmarshal(data, &lockFile)
	if err != nil {
		return nil, fmt.Errorf("Could not decode JSON dependency lock file: %s", err.Error())
	}

	return &lockFile, nil
}

func installDependencies(lockFile *LockFile) (bool, error) {
	if len(lockFile.Packages) == 0 {
		log.Print("No packages to install found in [<strong>%s</strong>]", LockFileName)
		return false, nil
	}

	log.Action("Installing %d %s", len(lockFile.Packages), log.Pluralize("package", len(lockFile.Packages)))
	for _, p := range lockFile.Packages {
		if err := installPackage(p); err != nil {
			return false, err
		}
	}
	return true, nil
}

func installPackage(p *PackageDefinition) error {
	if p.Source.Typ != "git" {
		return fmt.Errorf("bot install does currently only support git over HTTPS. Please come back later or do a pull request :)")
	}

	gitURL := fmt.Sprintf("https://%s", p.Name)
	grobot.Execute("git clone %s vendor/src/%s", gitURL, p.Name)
	grobot.Execute("cd vendor/src/%s && git checkout %s", p.Name, p.Source.Reference)
	return nil
}
