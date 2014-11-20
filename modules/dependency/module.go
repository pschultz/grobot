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
	conf *Configuration
}

func NewModule() *Module {
	return &Module{conf: defaultConfig}
}

func (m *Module) Name() string {
	return "Depenency"
}

func (m *Module) LoadConfiguration(config *grobot.Configuration) error {
	data, keyExists := config.Get(moduleConfigKey)
	if keyExists == false {
		log.Debug("Using default config")
		m.conf = defaultConfig
	} else {
		m.conf = new(Configuration)
		err := json.Unmarshal(*data, m.conf)
		if err != nil {
			return fmt.Errorf("could not parse configuration key '%s' : %s", moduleConfigKey, err.Error())
		}
	}

	log.Debug("Using vendors folder '%s'", m.conf.VendorsFolder)
	m.registerTasks()
	return nil
}

func (m *Module) registerTasks() {
	grobot.RegisterTask("install", NewInstallTask(m))
	grobot.RegisterTask("update", NewUpdateTask())
}
