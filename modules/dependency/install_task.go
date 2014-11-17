package dependency

import (
	"encoding/json"
	"fmt"
	"github.com/fgrosse/grobot"
	"github.com/fgrosse/grobot/log"
	"strings"
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
		log.Print("Lock file %S does not exist", LockFileName)
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
		log.Print("No packages to install found in %S", LockFileName)
		return false, nil
	}

	if len(lockFile.Packages) > 1 {
		log.Action("Processing %d dependencies", len(lockFile.Packages))
	}

	nrOfUpdates := 0
	for _, p := range lockFile.Packages {
		updated, err := installPackage(p)
		if err != nil {
			fmt.Print("  ")
			log.Error(err.Error())
		}
		if updated == true {
			nrOfUpdates++
		}
	}

	switch {
	case nrOfUpdates == 0:
		log.Action("No new packages have been installed")
	case nrOfUpdates == 1:
		log.Action("One package has been installed")
	default:
		log.Action("<strong>%d</strong> packages have been installed", nrOfUpdates)
	}

	return nrOfUpdates > 0, nil
}

func installPackage(p *PackageDefinition) (bool, error) {
	if p.Source.Typ != "git" {
		return false, fmt.Errorf("bot install does currently only support git over HTTPS. Please come back later or do a pull request :)")
	}

	vendorDir := fmt.Sprintf("vendor/src/%s", p.Name)
	log.Debug("Trying to install package %S from %s repo version %s", p.Name, p.Source.Typ, p.Source.Reference)
	targetInfo := grobot.TargetInfo(vendorDir)
	if targetInfo.ExistingFile {
		log.Debug("Directory %S does already exist", vendorDir)
		return false, checkIfPackageHasRequestedVersion(vendorDir, p)
	} else {
		log.Debug("Directory %S does not yet exist", vendorDir)
		return checkoutPackage(vendorDir, p)
	}
}

func checkIfPackageHasRequestedVersion(vendorDir string, p *PackageDefinition) (err error) {
	log.Debug("Checking repository version...")

	grobot.SetWorkingDirectory(vendorDir)
	cvsRef := grobot.ExecuteSilent("git rev-parse HEAD")
	if cvsRef == p.Source.Reference {
		log.ActionMinor("Package %S already up to date", p.Name)
		err = nil
	} else {
		err = fmt.Errorf("Package %s : repository at %s is not at the required version %s", p.Name, vendorDir, p.Source.Reference)
	}
	grobot.ResetWorkingDirectory()
	return err
}

func checkoutPackage(vendorDir string, p *PackageDefinition) (updated bool, err error) {
	log.Action("Installing package %S ...", p.Name)
	log.Debug("Determining repository URL ...")
	gitURL, err := gitURL(p.Name)
	if err != nil {
		return false, err
	}

	log.Debug("Repository URL for %S is %S", p.Name, gitURL)
	grobot.ExecuteSilent("git clone %s %s", gitURL, vendorDir)
	grobot.SetWorkingDirectory(vendorDir)
	grobot.ExecuteSilent("git checkout %s --quiet", p.Source.Reference)
	grobot.ResetWorkingDirectory()
	return true, nil
}

func gitURL(packageName string) (string, error) {
	if strings.HasPrefix(packageName, "code.google.com/") {
		return "https://" + packageName, nil
	}

	return repoRootForImportDynamic(packageName)
}
