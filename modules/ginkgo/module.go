package ginkgo

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
	return "Ginkgo"
}

func (m *Module) LoadConfiguration(config map[string]*json.RawMessage) error {
	data, keyExists := config["ginkgo"]
	if keyExists == false {
		log.Debug("Did not load Ginkgo module: configuration key ginkgo is not set")
		return nil
	}

	var newConfig Configuration
	err := json.Unmarshal(*data, &newConfig)
	if err != nil {
		return fmt.Errorf("could not parse configuration key 'ginkgo' : %s", err.Error())
	}

	m.conf = newConfig
	log.Debug("Using test folder '%s'", m.conf.TestFolder)

	m.registerTasks()
	return nil
}

func (m *Module) registerTasks() {
	generic.RegisterVendorBin("ginkgo", "github.com/onsi/ginkgo/ginkgo")
	// TODO we probably also want gomega dependency
	// go get github.com/onsi/gomega
	grobot.RegisterTask(grobot.StandardTaskTest, NewTestTask(m.conf))
	if m.conf.TestFolder != "" {
		grobot.RegisterDirectory(m.conf.TestFolder)
	}
}
