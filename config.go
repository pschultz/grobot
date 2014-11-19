package grobot

import (
	"encoding/json"
	"fmt"
	"github.com/fgrosse/grobot/log"
	"io/ioutil"
)

var isDebug = false

func EnableDebugMode() {
	isDebug = true
	log.EnableDebug()
}

func IsDebugMode() bool {
	return isDebug
}

type Configuration struct {
	Version          Version `json:"version"`
	RawModuleConfigs map[string]*json.RawMessage
}

func (c *Configuration) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &c.RawModuleConfigs)
	if version, versionIsDefined := c.RawModuleConfigs["version"]; versionIsDefined {
		delete(c.RawModuleConfigs, "version")
		err = json.Unmarshal(*version, &c.Version)
	}

	return err
}

func (c *Configuration) Get(field string) (raw *json.RawMessage, exists bool) {
	raw, exists = c.RawModuleConfigs[field]
	return
}

func LoadConfigFromFile(confFilePath string) error {
	log.Debug("Loading configuration from file '%s'", confFilePath)
	data, err := ioutil.ReadFile(confFilePath)
	if err != nil {
		return fmt.Errorf("Could not read configuration : %s", err.Error())
	}

	var config Configuration
	err = json.Unmarshal(data, &config)
	if err != nil {
		return fmt.Errorf("Error while unmarshalling configuration file '%s' : %s", confFilePath, err.Error())
	}

	for _, module := range modules {
		log.Debug("Loading configuration for module [<strong>%s</strong>]", module.Name())
		err = module.LoadConfiguration(&config)
		if err != nil {
			log.Error("Error whileloading module %s : %s", module.Name(), err.Error())
		}
		log.Debug("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	}

	return nil
}
