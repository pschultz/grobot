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

func (t *UpdateTask) Invoke(invokedName string, args ...string) (bool, error) {
	if len(args) == 0 {
		return false, fmt.Errorf("No package name given. Please tell me the full name of the package to update")
	}

	lockFile, err := loadLockFile()
	if err != nil {
		return false, err
	}

	packageName := args[0]
	return updatePackage(packageName, lockFile)
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
		log.Print("  Package allready up to date..")
		return false, nil
	} else {
		log.Print("  Installed new version %S", newVersion)
		packageInLockFile.Source.Version = newVersion
		err := writeLockFile(lockFile)
		return err == nil, err
	}
}
