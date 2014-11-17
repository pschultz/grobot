package dependency

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
