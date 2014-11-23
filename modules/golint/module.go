package golint

import (
	"github.com/fgrosse/grobot"
	"github.com/fgrosse/grobot/modules/generic"
)

func init() {
	grobot.RegisterModule(new(Module))
}

type Module struct{}

func (m *Module) Name() string {
	return "Golint"
}

func (m *Module) LoadConfiguration(config *grobot.Configuration) error {
	// TODO load lint config
	m.registerTasks()
	return nil
}

func (m *Module) registerTasks() {
	generic.RegisterVendorBin("golint", "github.com/golang/lint")
	grobot.RegisterTask("lint", NewLintTask())
}
