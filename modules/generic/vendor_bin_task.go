package generic

import (
	"github.com/fgrosse/gobot"
	"github.com/fgrosse/gobot/log"
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
	return true, gobot.Shell(`go build -o "%s" "%s"`, name, sourcePath)
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
	gobot.RegisterTask("vendor/bin/"+binName, NewVendorBinTask(sourceRepo))
	gobot.RegisterTask("vendor/src/"+sourceRepo, NewInstallDependencyTask(sourceRepo))
}
