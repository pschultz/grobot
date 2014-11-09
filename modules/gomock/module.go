package gomock

import (
	"encoding/json"
	"fmt"
	"github.com/fgrosse/gobot"
	"github.com/fgrosse/gobot/log"
	"github.com/fgrosse/gobot/modules/generic"
)

func init() {
	gobot.RegisterModule(new(Module))
}

type Module struct {
	conf Configuration
}

func (m *Module) Name() string {
	return "Gomock"
}

func (m *Module) LoadConfiguration(config map[string]*json.RawMessage) error {
	data, keyExists := config["gomock"]
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
	generic.RegisterVendorBin("mockgen", "code.google.com/p/gomock/mockgen")
	gobot.RegisterTask("mocks", NewAllMocksTask(m.conf))
	gobot.RegisterTask(m.conf.MockFolder, generic.NewCreateFolderTask())

	genericMockBuildRule := fmt.Sprintf(`^%s/\w+\.go$`, m.conf.MockFolder)
	gobot.RegisterRule(genericMockBuildRule, NewBuildMockFileTask(m.conf))

}
