package gomock

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
	conf Configuration
}

func (m *Module) Name() string {
	return "Gomock"
}

func (m *Module) LoadConfiguration(config *grobot.Configuration) error {
	data, keyExists := config.Get("gomock")
	if keyExists == false {
		log.Debug("Did not load Gomock module: configuration key gomock is not set")
		return nil
	}

	var newConfig Configuration
	err := json.Unmarshal(*data, &newConfig)
	if err != nil {
		return fmt.Errorf("could not parse configuration key 'gomock' : %s", err.Error())
	}

	m.conf = newConfig
	log.Debug("Using mock folder '%s'", m.conf.MockFolder)

	m.registerTasks()
	return nil
}

func (m *Module) registerTasks() {
	grobot.RegisterTask("mocks", NewAllMocksTask(m.conf))
	grobot.RegisterTaskHook(grobot.HookBefore, grobot.StandardTaskTest, "mocks")

	generic.RegisterVendorBin("mockgen", "code.google.com/p/gomock")
	grobot.RegisterDirectory(m.conf.MockFolder)

	genericMockBuildRule := fmt.Sprintf(`^%s/\w+_mock\.go$`, m.conf.MockFolder)
	grobot.RegisterRule(genericMockBuildRule, NewBuildMockFileTask(m.conf))

}
