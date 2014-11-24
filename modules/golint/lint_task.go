package golint

import (
	"bytes"
	"fmt"
	"github.com/fgrosse/grobot"
	"github.com/fgrosse/grobot/log"
	"github.com/fgrosse/grobot/modules/dependency"
	"strings"
)

type LintTask struct {
	conf         *Configuration
	rootDir      string
	ignoredFiles []string
}

func NewLintTask(conf *Configuration) *LintTask {
	return &LintTask{conf: conf}
}

func (t *LintTask) Dependencies(string) []string {
	return []string{"vendor/bin/golint"}
}

func (t *LintTask) Description() string {
	return "Run golint on all source files"
}

func (t *LintTask) Invoke(string, ...string) (bool, error) {
	t.setup()
	nrOfIssues := t.lintDirectory("")
	grobot.ResetShellWorkingDirectory()

	switch nrOfIssues {
	case 0:
		return false, nil
	case 1:
		return false, fmt.Errorf("Golint detected one issue")
	default:
		return false, fmt.Errorf("Golint detected %d issues", nrOfIssues)
	}
}

func (t *LintTask) setup() {
	t.rootDir = grobot.WorkingDir()
	module := grobot.GetModule("Depenency").(*dependency.Module)
	if module == nil {
		return
	}

	t.ignoredFiles = append(t.ignoredFiles, module.VendorDir())
}

func (t *LintTask) lintDirectory(relativePath string) (nrOfIssues int) {
	absolutePath := t.rootDir
	if relativePath != "" {
		absolutePath = absolutePath + "/" + relativePath
	}

	log.Debug("Searching for go files in %S", relativePath)
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
		if relativePath != "" {
			log.Action("Linting files in %S", relativePath)
			grobot.SetShellWorkingDirectory(relativePath)
		}
		lintingResult := grobot.Execute(golintCmd.String())
		nrOfIssues = nrOfIssues + t.getNrOfIssuesFrom(lintingResult)
	}

	if relativePath != "" {
		relativePath = relativePath + "/"
	}

	for _, d := range dirs {
		nrOfIssues = nrOfIssues + t.lintDirectory(relativePath+d.Name)
	}

	return
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

func (t *LintTask) getNrOfIssuesFrom(lintingResult string) int {
	lintingResult = strings.TrimSpace(lintingResult)
	if lintingResult == "" {
		return 0
	}

	issues := strings.Split(lintingResult, "\n")
	reportedIssues := []string{}
	for _, issue := range issues {
		if strings.HasSuffix(issue, "foooo") && t.conf.WarnCommentOrBeUnexported == false {
			// TODO continue here
			log.Debug("Ignoring issue %S", issue)
			continue
		}
		reportedIssues = append(reportedIssues, issue)
	}

	return len(reportedIssues)
}
