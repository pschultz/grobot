package golint

import (
	"encoding/json"
	"fmt"
	"github.com/fgrosse/grobot"
	"github.com/fgrosse/grobot/log"
	"github.com/fgrosse/grobot/modules/generic"
)

func init() {
	grobot.RegisterModule(new(Module))
}

type Module struct {
	conf *Configuration
}

func (m *Module) Name() string {
	return "Golint"
}

func (m *Module) LoadConfiguration(config *grobot.Configuration) error {
	data, keyExists := config.Get(moduleConfigKey)
	if keyExists == false {
		log.Debug("Using default config")
		m.conf = DefaultLintConfig
	} else {
		m.conf = new(Configuration)
		err := json.Unmarshal(*data, m.conf)
		if err != nil {
			return fmt.Errorf("could not parse configuration key '%s' : %s", moduleConfigKey, err.Error())
		}
	}

	m.registerTasks()
	return nil
}

func (m *Module) registerTasks() {
	generic.RegisterVendorBin("golint", "github.com/golang/lint")
	grobot.RegisterTask("lint", NewLintTask(m.conf))
}
