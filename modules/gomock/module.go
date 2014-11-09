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

	m.registerTasks()
	return nil
}

func (m *Module) registerTasks() {
	//return RegisterTask(name, NewFolderTask(newTask))
	gobot.RegisterTask("vendor/bin/mockgen", generic.NewVendorBinTask("code.google.com/p/gomock/mockgen"))
	gobot.RegisterTask("vendor/src/code.google.com/p/gomock/mockgen", generic.NewInstallDependencyTask("code.google.com/p/gomock/mockgen"))
	gobot.RegisterTask("mocks", new(MocksTask))
}
