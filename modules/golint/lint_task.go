package golint

import (
	"bytes"
	"github.com/fgrosse/grobot"
	"github.com/fgrosse/grobot/log"
	"github.com/fgrosse/grobot/modules/dependency"
	"strings"
)

type LintTask struct {
	rootDir      string
	ignoredFiles []string
}

func NewLintTask() *LintTask {
	return &LintTask{}
}

func (t *LintTask) Dependencies(invokedName string) []string {
	return []string{"vendor/bin/golint"}
}

func (t *LintTask) Description() string {
	return "Run golint on all source files"
}

func (t *LintTask) Invoke(targetName string, args ...string) (bool, error) {
	t.setup()
	t.lintDirectory("")
	grobot.ResetShellWorkingDirectory()
	return false, nil
}

func (t *LintTask) setup() {
	t.rootDir = grobot.WorkingDir()
	module := grobot.GetModule("Depenency").(*dependency.Module)
	if module == nil {
		return
	}

	t.ignoredFiles = append(t.ignoredFiles, module.VendorDir())
}

func (t *LintTask) lintDirectory(relativePath string) {
	absolutePath := t.rootDir
	if relativePath != "" {
		absolutePath = absolutePath + "/" + relativePath
	}

	log.Debug("Linting files in %S", absolutePath)
	filesAndDirs := grobot.ListFiles(absolutePath)
	golintCmd := bytes.NewBuffer([]byte(`golint`))
	dirs := []*grobot.File{}
	atLeastOneGoFileFound := false
	for _, f := range filesAndDirs {
		if f.IsDir {
			if t.isIgnored(relativePath, f.Name) {
				log.Debug("Ignoring directory %S", f.Name)
				continue
			}
			dirs = append(dirs, f)
			continue
		}

		if strings.HasPrefix(f.Name, ".") {
			log.Debug("Ignoring file %S (hidden file)", f.Name)
			continue
		}
		if strings.HasSuffix(f.Name, ".go") == false {
			log.Debug("Ignoring file %S (no *.go extension)", f.Name)
			continue
		}

		golintCmd.WriteString(` "`)
		golintCmd.WriteString(f.Name)
		golintCmd.WriteString(`"`)
		atLeastOneGoFileFound = true
	}

	if atLeastOneGoFileFound {
		grobot.SetShellWorkingDirectory(relativePath)
		grobot.Execute(golintCmd.String())
	}

	if relativePath != "" {
		relativePath = relativePath + "/"
	}

	for _, d := range dirs {
		t.lintDirectory(relativePath + d.Name)
	}
}

func (t *LintTask) isIgnored(relativePath, fileName string) bool {
	name := relativePath + fileName
	for _, ignored := range t.ignoredFiles {
		if name == ignored {
			return true
		}
	}
	return false
}
