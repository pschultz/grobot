package dependency

import "github.com/fgrosse/grobot"

const moduleConfigKey = "dependency"

type Configuration struct {
	VendorsFolder string                     `json:"folder"`
	Packages      []*PackageConfigDefinition `json:"packages"`
	globalConfig  *grobot.Configuration
}

var defaultConfig = &Configuration{
	VendorsFolder: "vendor",
	Packages:      []*PackageConfigDefinition{},
}

type PackageConfigDefinition struct {
	Name    string          `json:"name"`
	Type    string          `json:"type"`
	Version *grobot.Version `json:"version"`
}
