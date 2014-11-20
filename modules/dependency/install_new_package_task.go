package dependency

import (
	"encoding/json"
	"fmt"
	"github.com/fgrosse/grobot"
	"github.com/fgrosse/grobot/log"
)

func (t *InstallTask) installNewDependency(packageName string, lockFile *LockFile) (updated bool, err error) {
	log.Action("Installing new package %S", packageName)
	vendorDir, err := t.checkInstallationDirectoryDoesNotExist(packageName)
	if err != nil {
		return false, err
	}

	gitURL, err := gitURL(packageName)
	if err != nil {
		return false, err
	}

	installedVersion := t.checkoutNewPackage(gitURL, vendorDir)
	updated, err = t.updateLockFile(packageName, installedVersion, lockFile)
	if err != nil {
		return false, err
	}

	err = t.addNewDependencyToConfiguration(packageName)
	return err != nil, err
}

func (t *InstallTask) checkInstallationDirectoryDoesNotExist(packageName string) (string, error) {
	vendorDir := getInstallDestination(packageName)
	targetInfo := grobot.TargetInfo(vendorDir)
	if targetInfo.ExistingFile {
		return "", fmt.Errorf("Can not install new package %s : directory %s does already exist", packageName, vendorDir)
	}

	log.Debug("Directory %S does not yet exist", vendorDir)
	return vendorDir, nil
}

func (t *InstallTask) checkoutNewPackage(gitURL, vendorDir string) string {
	grobot.ExecuteSilent("git clone %s %s", gitURL, vendorDir)
	grobot.SetWorkingDirectory(vendorDir)
	installedVersion := grobot.ExecuteSilent("git rev-parse HEAD")
	log.Debug("Successfully installed version %S into %S", installedVersion, vendorDir)
	grobot.ResetWorkingDirectory()
	return installedVersion
}

func (t *InstallTask) updateLockFile(packageName, installedVersion string, lockFile *LockFile) (updated bool, err error) {
	p := newGitPackage(packageName, installedVersion)
	if lockFile == nil {
		log.Debug("Creating new lockfile %S", LockFileName)
		lockFile = &LockFile{[]*PackageDefinition{}}
	} else {
		log.Debug("Updating existing lockfile %S", LockFileName)
	}

	packageInLockFile := lockFile.Package(packageName)
	if packageInLockFile == nil || packageInLockFile.Source.Version != installedVersion {
		lockFile.Packages = append(lockFile.Packages, p)
	} else {
		log.Action("Package was already contained in %S and has been updated", LockFileName)
	}

	// TODO only write lock file if there was an update
	err = writeLockFile(lockFile)
	return err == nil, err
}

func (t *InstallTask) addNewDependencyToConfiguration(packageName string) error {
	botConfig := grobot.CurrentConfig()
	configFileName := botConfig.FileName()
	moduleConfig := t.module.conf
	log.Debug("Adding new dependency %S to configuration file %S", packageName, configFileName)

	if isPackageAlreadyExistentInConfiguration(packageName, moduleConfig) {
		return fmt.Errorf("Package dependency %s is already existing in configuration file %s", packageName, configFileName)
	}

	moduleConfig.Packages = append(moduleConfig.Packages, &PackageConfigDefinition{
		Name:    packageName,
		Type:    "git",
		Version: grobot.NewVersion("branch:master"),
	})
	moduleConfigBytes, err := json.Marshal(moduleConfig)
	moduleConfigRaw := json.RawMessage(moduleConfigBytes)
	botConfig.RawModuleConfigs[moduleConfigKey] = &moduleConfigRaw

	data, err := json.MarshalIndent(botConfig, "", "    ")
	if err != nil {
		return fmt.Errorf("Could not write configuration %s : %s", configFileName, err.Error())
	}

	return grobot.WriteFile(configFileName, data)
}

func isPackageAlreadyExistentInConfiguration(packageName string, config *Configuration) bool {
	for _, p := range config.Packages {
		if p.Name == packageName {
			return true
		}
	}
	return false
}
