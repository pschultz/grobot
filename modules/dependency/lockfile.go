package dependency

import (
	"github.com/fgrosse/grobot"
	"github.com/fgrosse/grobot/log"

	"encoding/json"
	"fmt"
	"strings"
)

type LockFile struct {
	Packages []*PackageDefinition `json:"packages"`
}

type PackageDefinition struct {
	Name   string               `json:"name"`
	Source *SourceConfiguration `json:"source"`
}

type SourceConfiguration struct {
	Typ     string `json:"type"`
	Version string `json:"version"`
}

func newGitPackage(name, version string) *PackageDefinition {
	return &PackageDefinition{
		Name: name,
		Source: &SourceConfiguration{
			Typ:     "git",
			Version: version,
		},
	}
}

func loadLockFile() (*LockFile, error) {
	targetInfo := grobot.TargetInfo(LockFileName)

	if targetInfo.ExistingFile == false {
		log.Print("Lock file %S does not yet exist and will not be loaded", LockFileName)
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

func writeLockFile(lockFile *LockFile) error {
	log.Debug("Writing lock file %S", LockFileName)
	data, err := json.MarshalIndent(lockFile, "", "    ")
	if err != nil {
		return fmt.Errorf("Could not write lockfile %s : %s", LockFileName, err.Error())
	}

	return grobot.WriteFile(LockFileName, data)
}

func (l *LockFile) Package(packageName string) *PackageDefinition {
	matches := []*PackageDefinition{}
	for _, p := range l.Packages {
		if p.Name == packageName {
			return p
		}

		if strings.Contains(p.Name, packageName) {
			matches = append(matches, p)
		}
	}

	if len(matches) == 1 {
		log.Debug("Autocompleted package name %S to %S", packageName, matches[0].Name)
		return matches[0]
	}

	return nil
}
