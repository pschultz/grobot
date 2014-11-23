package dependency

import (
	"encoding/json"
	"fmt"
	"github.com/fgrosse/grobot"
	"github.com/fgrosse/grobot/log"
)

func (t *InstallTask) installNewDependency(packageName string, lockFile *LockFile) (updated bool, err error) {
	updated, err = t.installNewDependencyRecursive(packageName, lockFile)
	if err != nil {
		return false, err
	}

	err = writeLockFile(lockFile)
	if err != nil {
		return false, err
	}

	err = t.addNewDependencyToConfiguration(packageName)
	return err != nil && updated, err
}

func (t *InstallTask) installNewDependencyRecursive(packageName string, lockFile *LockFile) (updated bool, err error) {
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
	t.updateLockFile(packageName, installedVersion, lockFile)
	t.installTransitiveDependencies(packageName, vendorDir, lockFile)
	return true, nil
}

func (t *InstallTask) checkInstallationDirectoryDoesNotExist(packageName string) (string, error) {
	vendorDir := getInstallDestination(packageName)
	targetInfo := grobot.FileInfo(vendorDir)
	if targetInfo.ExistingFile {
		return "", fmt.Errorf("Can not install new package %s : directory %s does already exist", packageName, vendorDir)
	}

	log.Debug("Directory %S does not yet exist", vendorDir)
	return vendorDir, nil
}

func (t *InstallTask) checkoutNewPackage(gitURL, vendorDir string) string {
	grobot.ExecuteSilent("git clone %s %s", gitURL, vendorDir)
	grobot.SetShellWorkingDirectory(vendorDir)
	installedVersion := grobot.ExecuteSilent("git rev-parse HEAD")
	log.Debug("Successfully installed version %S into %S", installedVersion, vendorDir)
	grobot.ResetShellWorkingDirectory()
	return installedVersion
}

func (t *InstallTask) updateLockFile(packageName, installedVersion string, lockFile *LockFile) {
	p := newGitPackage(packageName, installedVersion)
	packageInLockFile := lockFile.Package(packageName)
	if packageInLockFile == nil || packageInLockFile.Source.Version != installedVersion {
		lockFile.Packages = append(lockFile.Packages, p)
	} else {
		log.Action("Package was already contained in %S and has been updated", LockFileName)
	}
}

func (t *InstallTask) installTransitiveDependencies(packageName, vendorDir string, lockFile *LockFile) error {
	log.Debug("Checking %S for transitive dependencies", packageName)
	vendorBotConfigFile := vendorDir + "/" + grobot.ConfigFileName
	if grobot.FileExists(vendorBotConfigFile) == false {
		return nil
	}

	vendorDepConf, err := loadVendorDependencyConfig(vendorBotConfigFile)
	if vendorDepConf == nil || err != nil {
		return err
	}

	if len(vendorDepConf.Packages) == 0 {
		log.Debug("No dependencies configured in %S", vendorBotConfigFile)
		return nil
	}

	if len(vendorDepConf.Packages) == 1 {
		log.Debug("Installing one transitive dependency of %S", packageName)
	} else {
		log.Debug("Installing %d transitive dependencies of %S", len(vendorDepConf.Packages), packageName)
	}

	for _, p := range vendorDepConf.Packages {
		if lockFile.Package(p.Name) != nil {
			log.Debug("Package %S is already in lockfile", p.Name)
			continue
		}

		_, err := t.installNewDependencyRecursive(p.Name, lockFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func loadVendorDependencyConfig(confFilePath string) (*Configuration, error) {
	data, err := grobot.ReadFile(confFilePath)
	if err != nil {
		return nil, fmt.Errorf("Could not read configuration : %s", err.Error())
	}

	vendorConfig := new(grobot.Configuration)
	err = json.Unmarshal(data, vendorConfig)
	if err != nil {
		return nil, fmt.Errorf("Error while unmarshalling configuration file '%s' : %s", confFilePath, err.Error())
	}

	dependencyConfData, keyExists := vendorConfig.Get(moduleConfigKey)
	if keyExists == false {
		log.Debug("No dependency configuration found")
		return nil, nil
	}

	vendorDepConf := new(Configuration)
	err = json.Unmarshal(*dependencyConfData, vendorDepConf)
	if err != nil {
		return nil, fmt.Errorf("could not parse configuration key '%s' from %s : %s", moduleConfigKey, confFilePath, err.Error())
	}

	return vendorDepConf, nil
}

func (t *InstallTask) addNewDependencyToConfiguration(packageName string) error {
	moduleConfig := t.module.conf
	log.Debug("Adding new dependency %S to configuration file %S", packageName, grobot.ConfigFileName)

	if isPackageAlreadyExistentInConfiguration(packageName, moduleConfig) {
		return fmt.Errorf("Package dependency %s is already existing in configuration file %s", packageName, grobot.ConfigFileName)
	}

	moduleConfig.Packages = append(moduleConfig.Packages, &PackageConfigDefinition{
		Name:    packageName,
		Type:    "git",
		Version: grobot.NewVersion("branch:master"),
	})
	moduleConfigBytes, err := json.Marshal(moduleConfig)
	moduleConfigRaw := json.RawMessage(moduleConfigBytes)
	botConfig := moduleConfig.globalConfig
	botConfig.RawModuleConfigs[moduleConfigKey] = &moduleConfigRaw

	data, err := json.MarshalIndent(botConfig, "", "    ")
	if err != nil {
		return fmt.Errorf("Could not write configuration %s : %s", grobot.ConfigFileName, err.Error())
	}

	return grobot.WriteFile(grobot.ConfigFileName, data)
}

func isPackageAlreadyExistentInConfiguration(packageName string, config *Configuration) bool {
	for _, p := range config.Packages {
		if p.Name == packageName {
			return true
		}
	}
	return false
}
