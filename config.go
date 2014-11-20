package grobot

import (
	"encoding/json"
	"fmt"
	"github.com/fgrosse/grobot/log"
)

const ConfigFileName = "bot.json"

var (
	isDebug       = false
	defaultConfig = Configuration{
		Version:          NewVersion("none"),
		RawModuleConfigs: map[string]*json.RawMessage{},
	}
)

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
	fileName         string
}

func (c *Configuration) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &c.RawModuleConfigs)
	if version, versionIsDefined := c.RawModuleConfigs["bot-version"]; versionIsDefined {
		delete(c.RawModuleConfigs, "bot-version")
		c.Version = new(Version)
		err = json.Unmarshal(*version, c.Version)
	} else {
		c.Version = NoVersion
	}

	return err
}

func (c *Configuration) MarshalJSON() ([]byte, error) {
	data := c.RawModuleConfigs
	if data == nil {
		data = map[string]*json.RawMessage{}
	}

	if c.Version != NoVersion {
		v, err := json.Marshal(c.Version)
		if err != nil {
			return []byte{}, err
		}
		rawVersion := json.RawMessage(v)
		data["bot-version"] = &rawVersion
	}
	return json.Marshal(data)
}

func (c *Configuration) Get(field string) (raw *json.RawMessage, exists bool) {
	raw, exists = c.RawModuleConfigs[field]
	return
}

// LoadConfigFromFile loads a bot configuration from the given file and initializes all modules
// with this configuration.
func LoadConfigFromFile(confFilePath string, currentVersion *Version) (*Configuration, error) {
	log.Debug("Loading configuration from file '%s'", confFilePath)
	data, err := ReadFile(confFilePath)
	if err != nil {
		return nil, fmt.Errorf("Could not read configuration : %s", err.Error())
	}

	config := new(Configuration)
	config.fileName = confFilePath
	err = json.Unmarshal(data, config)
	if err != nil {
		return nil, fmt.Errorf("Error while unmarshalling configuration file '%s' : %s", confFilePath, err.Error())
	}

	if config.Version.GreaterThen(currentVersion) {
		return nil, fmt.Errorf(`Error while read configuration file %s : The minimum required bot version is "%s" but you are running bot version "%s"`, confFilePath, config.Version.String(), currentVersion.String())
	}

	loadModules(config)
	return config, nil
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
