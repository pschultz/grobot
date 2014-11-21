package golint

type LintTask struct{}

func NewLintTask() *LintTask {
	return &LintTask{}
}

func (t *LintTask) Dependencies(invokedName string) []string {
	return []string{}
}

func (t *LintTask) Invoke(targetName string, args ...string) (bool, error) {
	return false, nil
}
