package golint

import (
	"github.com/fgrosse/grobot"
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
	grobot.RegisterTask("lint", NewLintTask())
}
