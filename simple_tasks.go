package grobot

type CreateFolderTask struct{}

func (t *CreateFolderTask) Dependencies(string) []string {
	return []string{}
}

func (t *CreateFolderTask) Invoke(path string) (bool, error) {
	return true, Execute(`mkdir -p "%s"`, path)
}
