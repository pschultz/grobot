package gobot

import (
	"encoding/json"
	"fmt"
	"github.com/fgrosse/gobot/log"
	"io/ioutil"
)

func LoadConfigFromFile(confFilePath string) error {
	log.Debug("Loading configuration from file '%s'", confFilePath)
	data, err := ioutil.ReadFile(confFilePath)
	if err != nil {
		return fmt.Errorf("Could not read configuration : %s", err.Error())
	}

	var config map[string]*json.RawMessage
	err = json.Unmarshal(data, &config)
	if err != nil {
		return fmt.Errorf("Error while unmarshalling configuration file '%s' : %s", confFilePath, err.Error())
	}

	for _, module := range modules {
		err = module.LoadConfiguration(config)
		if err != nil {
			log.Error("Error whileloading module %s : %s", module.Name(), err.Error())
		}
	}
	return nil
}
