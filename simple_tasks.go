package grobot

import "github.com/fgrosse/grobot/log"

const BotVersion = "0.6"

func init() {
	RegisterTask("version", &VersionTask{})
}

type VersionTask struct{}

func (t *VersionTask) Description() string {
	return "Show the current version of your bot"
}

func (t *VersionTask) Dependencies(string) []string {
	return []string{}
}

func (t *VersionTask) Invoke(string, ...string) (bool, error) {
	log.Action("Bot version %s", BotVersion)
	return true, nil
}

type CreateDirectoryTask struct{}

func (t *CreateDirectoryTask) Dependencies(string) []string {
	return []string{}
}

func (t *CreateDirectoryTask) Invoke(path string, arguments ...string) (bool, error) {
	log.Action("Creating directory %s", path)
	Execute(`mkdir -p "%s"`, path)
	return true, nil
}
