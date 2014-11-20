package generic

import (
	"github.com/fgrosse/grobot"
	"github.com/fgrosse/grobot/log"
	"strings"
)

type VendorBinTask struct {
	sourcePath string
	binName    string
}

func NewVendorBinTask(sourcePath, binName string) *VendorBinTask {
	return &VendorBinTask{"vendor/src/" + sourcePath, binName}
}

func (t *VendorBinTask) Dependencies(invokedName string) []string {
	return []string{t.sourcePath}
}

func (t *VendorBinTask) Invoke(name string, args ...string) (bool, error) {
	sourcePath := stripVendorSource(t.sourcePath)
	log.Action("Compiling %s..", name)
	grobot.Execute(`go build -o "%s" "%s/%s"`, name, sourcePath, t.binName)
	return true, nil
}

func stripVendorSource(path string) string {
	if strings.HasPrefix(path, "vendor/src/") {
		path = path[len("vendor/src/"):]
	}
	return path
}

// RegisterVendorBin registers a binary file which is compiled from
// some vendor source.
// Example:
//   RegisterVendorBin("mockgen", "code.google.com/p/gomock/mockgen")
func RegisterVendorBin(binName, sourceRepo string) {
	grobot.RegisterTask("vendor/bin/"+binName, NewVendorBinTask(sourceRepo, binName))
	grobot.RegisterTask("vendor/src/"+sourceRepo, NewInstallDependencyTask(sourceRepo))
}
