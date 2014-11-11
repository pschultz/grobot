package dependency

import (
	"encoding/json"
	"fmt"
	"github.com/fgrosse/grobot"
	"github.com/fgrosse/grobot/log"
)

const LockFileName = "bot.lock.json"

func init() {
	grobot.RegisterModule(new(Module))
}

type Module struct {
	conf Configuration
}

func (m *Module) Name() string {
	return "Depenency"
}

func (m *Module) LoadConfiguration(config map[string]*json.RawMessage) error {
	data, keyExists := config["dependency"]
	if keyExists == false {
		log.Debug("Using default config")
		m.conf = defaultConfig
	} else {
		var newConfig Configuration
		err := json.Unmarshal(*data, &newConfig)
		if err != nil {
			return fmt.Errorf("could not parse configuration key 'dependency' : %s", err.Error())
		}

		m.conf = newConfig
	}

	log.Debug("Using vendors folder '%s'", m.conf.VendorsFolder)
	m.registerTasks()
	return nil
}

func (m *Module) registerTasks() {
	grobot.RegisterTask("install", NewInstallTask())
}
