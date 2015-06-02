package dependency

import (
	"fmt"
	"github.com/fgrosse/grobot"
	"github.com/fgrosse/grobot/log"
)

type UpdateTask struct{}

func NewUpdateTask() *UpdateTask {
	return &UpdateTask{}
}

func (t *UpdateTask) Description() string {
	return "Update a given package to the newest version"
}

func (t *UpdateTask) Dependencies(string) []string {
	return []string{}
}

func (t *UpdateTask) Invoke(invokedName string, args ...string) (updated bool, err error) {
	lockFile, err := loadLockFile()
	if err != nil {
		return false, err
	}

	if len(args) == 0 {
		updated, err = updateAllPackages(lockFile)
	} else {
		packageName := args[0]
		updated, err = updatePackage(packageName, lockFile)
	}
	if err != nil {
		return false, err
	}

	if updated == false {
		return false, nil
	}

	err = writeLockFile(lockFile)
	return err == nil && updated, err
}

func updateAllPackages(lockFile *LockFile) (bool, error) {
	if len(lockFile.Packages) == 0 {
		log.Print("No packages found in lock file %S", LockFileName)
		return false, nil
	}

	if len(lockFile.Packages) >= 1 {
		log.ActionMinor("Updating %d packages...", len(lockFile.Packages))
	}

	somePackageHasBeenUpdated := false
	for _, p := range lockFile.Packages {
		updated, err := updatePackage(p.Name, lockFile)
		if err != nil {
			return somePackageHasBeenUpdated, err
		}
		somePackageHasBeenUpdated = somePackageHasBeenUpdated || updated
	}

	return somePackageHasBeenUpdated, nil
}

func updatePackage(packageName string, lockFile *LockFile) (bool, error) {
	packageInLockFile := lockFile.Package(packageName)
	if packageInLockFile == nil {
		return false, fmt.Errorf("Package %s is not contained in the lockfile %s", packageName, LockFileName)
	}

	log.Action("Updating package %S", packageInLockFile.Name)
	packageName = packageInLockFile.Name
	vendorDir := getInstallDestination(packageName)
	oldVersion := packageInLockFile.Source.Version

	grobot.SetWorkingDirectory(vendorDir)
	grobot.ExecuteSilent("git checkout master --quiet")
	grobot.ExecuteSilent("git pull")
	newVersion := grobot.ExecuteSilent("git rev-parse HEAD")
	grobot.ResetWorkingDirectory()

	if newVersion == oldVersion {
		log.Print("  Package already up to date..")
		return false, nil
	} else {
		log.Print("  Installed new version %S", newVersion)
		packageInLockFile.Source.Version = newVersion
		return true, nil
	}
}
