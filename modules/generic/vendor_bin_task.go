package generic

import (
	"github.com/fgrosse/grobot"
	"github.com/fgrosse/grobot/log"
	"strings"
)

type VendorBinTask struct {
	sourcePath string
}

func NewVendorBinTask(sourcePath string) *FileTask {
	return NewFileTask(&VendorBinTask{"vendor/src/" + sourcePath})
}

func (t *VendorBinTask) Dependencies(invokedName string) []string {
	return []string{t.sourcePath}
}

func (t *VendorBinTask) Invoke(name string) (bool, error) {
	sourcePath := stripVendorSource(t.sourcePath)
	log.Action("Compiling %s..", name)
	return true, grobot.Execute(`go build -o "%s" "%s"`, name, sourcePath)
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
	grobot.RegisterTask("vendor/bin/"+binName, NewVendorBinTask(sourceRepo))
	grobot.RegisterTask("vendor/src/"+sourceRepo, NewInstallDependencyTask(sourceRepo))
}
