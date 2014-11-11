package dependency

type LockFile struct {
	Packages []*PackageDefinition `json:"packages"`
}

type PackageDefinition struct {
	Name   string               `json:"name"`
	Source *SourceConfiguration `json:"source"`
}

type SourceConfiguration struct {
	Typ       string `json:"type"`
	Reference string `json:"reference"`
}
