package grobot

import (
	"encoding/json"
	"fmt"
	"github.com/fgrosse/grobot/log"
)

var isDebug = false

var defaultConfig = Configuration{
	Version: NewVersion("none"),
}

func EnableDebugMode() {
	isDebug = true
	log.EnableDebug()
}

func IsDebugMode() bool {
	return isDebug
}

type Configuration struct {
	Version          *Version `json:"bot-version"`
	RawModuleConfigs map[string]*json.RawMessage
}

func (c *Configuration) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &c.RawModuleConfigs)
	if version, versionIsDefined := c.RawModuleConfigs["bot-version"]; versionIsDefined {
		delete(c.RawModuleConfigs, "bot-version")
		c.Version = new(Version)
		err = json.Unmarshal(*version, c.Version)
	} else {
		c.Version = NewVersion("none")
	}

	return err
}

func (c *Configuration) Get(field string) (raw *json.RawMessage, exists bool) {
	raw, exists = c.RawModuleConfigs[field]
	return
}

func LoadConfigFromFile(confFilePath string, currentVersion *Version) error {
	log.Debug("Loading configuration from file '%s'", confFilePath)
	data, err := ReadFile(confFilePath)
	if err != nil {
		return fmt.Errorf("Could not read configuration : %s", err.Error())
	}

	var config Configuration
	err = json.Unmarshal(data, &config)
	if err != nil {
		return fmt.Errorf("Error while unmarshalling configuration file '%s' : %s", confFilePath, err.Error())
	}

	if config.Version.GreaterThen(currentVersion) {
		return fmt.Errorf(`Error while read configuration file %s : The minimum required bot version is "%s" but you are running bot version "%s"`, confFilePath, config.Version.String(), currentVersion.String())
	}

	loadModules(&config)
	return nil
}

func LoadBuiltinConfig() {
	log.Debug("Loading modules from builtin configuration")
	loadModules(&defaultConfig)
}

func loadModules(config *Configuration) {
	for _, module := range modules {
		log.Debug("Loading configuration for module [<strong>%s</strong>]", module.Name())
		err := module.LoadConfiguration(config)
		if err != nil {
			log.Error("Error whileloading module %s : %s", module.Name(), err.Error())
		}
		log.Debug("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	}
}
