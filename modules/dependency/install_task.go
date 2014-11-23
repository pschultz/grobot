package dependency

import (
	"fmt"
	"github.com/fgrosse/grobot"
	"github.com/fgrosse/grobot/log"
	"strings"
)

type InstallTask struct {
	module *Module
}

func NewInstallTask(m *Module) *InstallTask {
	return &InstallTask{m}
}

func (t *InstallTask) Description() string {
	return "Install new dependency by providing an additional argument or install all dependencies from " + LockFileName
}

func (t *InstallTask) Dependencies(string) []string {
	return []string{}
}

func (t *InstallTask) Invoke(invokedName string, args ...string) (bool, error) {
	lockFile, err := loadLockFile()
	if err != nil {
		return false, err
	}

	if len(args) == 0 {
		return installDependencies(lockFile)
	}

	packageName := args[0]
	return t.installNewDependency(packageName, lockFile)
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

	vendorDir := getInstallDestination(p.Name)
	log.Debug("Trying to install package %S from %s repo version %s", p.Name, p.Source.Typ, p.Source.Version)
	targetInfo := grobot.FileInfo(vendorDir)
	if targetInfo.ExistingFile {
		log.Debug("Directory %S does already exist", vendorDir)
		return false, checkIfPackageHasRequestedVersion(vendorDir, p)
	} else {
		log.Debug("Directory %S does not yet exist", vendorDir)
		return checkoutPackage(vendorDir, p)
	}
}

func getInstallDestination(packageName string) string {
	return fmt.Sprintf("vendor/src/%s", packageName)
}

func checkIfPackageHasRequestedVersion(vendorDir string, p *PackageDefinition) (err error) {
	log.Debug("Checking repository version...")

	grobot.SetShellWorkingDirectory(vendorDir)
	cvsRef := grobot.ExecuteSilent("git rev-parse HEAD")
	if cvsRef == p.Source.Version {
		log.ActionMinor("Package %S already up to date", p.Name)
		err = nil
	} else {
		err = fmt.Errorf("Package %s : repository at %s is not at the required version %s", p.Name, vendorDir, p.Source.Version)
	}
	grobot.ResetShellWorkingDirectory()
	return err
}

func checkoutPackage(vendorDir string, p *PackageDefinition) (updated bool, err error) {
	log.Action("Installing package %S ...", p.Name)
	gitURL, err := gitURL(p.Name)
	if err != nil {
		return false, err
	}

	grobot.ExecuteSilent("git clone %s %s", gitURL, vendorDir)
	grobot.SetShellWorkingDirectory(vendorDir)
	grobot.ExecuteSilent("git checkout %s --quiet", p.Source.Version)
	grobot.ResetShellWorkingDirectory()
	return true, nil
}

func gitURL(packageName string) (url string, err error) {
	log.Debug("Determining repository URL ...")
	if strings.HasPrefix(packageName, "code.google.com/") {
		url = "https://" + packageName
	} else {
		url, err = repoRootForImportDynamic(packageName)
	}

	if err == nil {
		log.Debug("Repository URL for %S is %S", packageName, url)
	}
	return url, err
}
