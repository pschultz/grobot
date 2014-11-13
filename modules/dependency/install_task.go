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
	targetInfo := grobot.TargetInfo(LockFileName)

	if targetInfo.ExistingFile == false {
		log.Print("Lock file [<strong>%s</strong>] does not exist", LockFileName)
		return nil, nil
	}

	log.Debug("Reading dependency lock file [<strong>%s</strong>]", LockFileName)
	data, err := grobot.ReadFile(LockFileName)
	if err != nil {
		return nil, fmt.Errorf("Error while reading dependency lock file: %s", err.Error())
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("Error while reading dependency lock file: Empty file")
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

	if len(lockFile.Packages) > 1 {
		log.Action("Installing %d packages", len(lockFile.Packages))
	}

	for _, p := range lockFile.Packages {
		if err := installPackage(p); err != nil {
			return false, err
		}
	}
	return true, nil
}

func installPackage(p *PackageDefinition) error {
	log.Action("Installing package %s", p.Name)
	if p.Source.Typ != "git" {
		return fmt.Errorf("bot install does currently only support git over HTTPS. Please come back later or do a pull request :)")
	}

	vendorDir := fmt.Sprintf("vendor/src/%s", p.Name)
	targetInfo := grobot.TargetInfo(vendorDir)
	if targetInfo.ExistingFile {
		return checkIfPackageHasRequestedVersion(vendorDir, p)
	} else {
		gitURL := fmt.Sprintf("https://%s", p.Name)
		grobot.Execute("git clone %s %s", gitURL, vendorDir)
		grobot.SetWorkingDirectory(vendorDir)
		grobot.Execute("git checkout %s --quiet", p.Source.Reference)
		grobot.ResetWorkingDirectory()
	}

	return nil
}

func checkIfPackageHasRequestedVersion(vendorDir string, p *PackageDefinition) (err error) {
	log.Debug("Directory [<strong>%s</strong>] is already existent", vendorDir)

	grobot.SetWorkingDirectory(vendorDir)
	cvsRef := grobot.ExecuteSilent("git rev-parse HEAD")
	if cvsRef == p.Source.Reference {
		log.Print("  package already up to date (%s)", cvsRef)
		err = nil
	} else {
		err = fmt.Errorf("Repository at %s is not at the required version %s", vendorDir, p.Source.Reference)
	}
	grobot.ResetWorkingDirectory()
	return err
}
