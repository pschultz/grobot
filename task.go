package grobot

import (
	"fmt"
	"github.com/fgrosse/grobot/log"
	"regexp"
	"strings"
)

const (
	StandardTaskDefault = "default"
	StandardTaskTest    = "test"
)

var (
	tasks                = map[string]Task{}
	rules                = map[*regexp.Regexp]Task{}
	resolvedDependencies = map[string]bool{}
)

type Task interface {
	Dependencies(invokedName string) []string
	Invoke(name string, arguments ...string) (bool, error)
}

type nullTask struct{}

func (t *nullTask) Dependencies(string) []string   { return []string{} }
func (t *nullTask) Invoke(string, ...string) error { return nil }

type Describable interface {
	Description() string
}

// Reset is used to make bot forget about all registered tasks and rules
// This is probably only useful in tests
func Reset() {
	tasks = map[string]Task{}
	rules = map[*regexp.Regexp]Task{}
	resolvedDependencies = map[string]bool{}
	hooks = map[string][]*TaskHook{}
	isDebug = false
}

func RegisterTask(name string, newTask Task) error {
	if _, keyExists := tasks[name]; keyExists == true {
		return fmt.Errorf("Module error: Task '%s' has already been registered", name)
	}

	log.Debug("Registering task [<strong>%s</strong>] as %T", name, newTask)
	tasks[name] = newTask
	return nil
}

func RegisterRule(ruleRegex string, newTask Task) error {
	rule, err := regexp.Compile(ruleRegex)
	if err != nil {
		return fmt.Errorf("Could not compile rule regex: %s", err.Error())
	}

	log.Debug("Registering rule [<strong>/%s/</strong>] as %T", ruleRegex, newTask)
	rules[rule] = newTask
	return nil
}

func RegisterDirectory(path string) error {
	task := CreateDirectoryTask{}
	return RegisterTask(path, &task)
}

func GetTask(name string) (Task, error) {
	if task, taskExists := tasks[name]; taskExists == true {
		return task, nil
	}

	for rule, task := range rules {
		if rule.MatchString(name) {
			return task, nil
		}
	}

	return nil, fmt.Errorf("Don't know how to build task '%s'", name)
}

func PrintTasks() {
	taskDescriptions := map[string]string{}
	longestTaskName := 0
	for name, task := range tasks {
		switch t := task.(type) {
		case Describable:
			description := t.Description()
			if description == "" {
				continue
			}
			taskDescriptions[name] = description
			if len(name) > longestTaskName {
				longestTaskName = len(name)
			}

		}
	}

	for name, description := range taskDescriptions {
		whiteSpace := strings.Repeat(" ", longestTaskName-len(name))
		log.Print("<em>%s</em> %s: %s", name, whiteSpace, description)
	}
}

func InvokeTask(invokedName string, recursionDepth int, args ...string) (bool, error) {
	checkHooks(HookBefore, invokedName, recursionDepth)

	resolvedDependencies[invokedName] = true
	target := FileInfo(invokedName)

	debugPrefix := ""
	log.SetDebugIndent(0)
	if recursionDepth > 0 {
		log.SetDebugIndent(3 * (recursionDepth - 1))
		debugPrefix = "┗━ "
	}

	task, err := GetTask(invokedName)
	if target.ExistingFile && task == nil {
		log.Debug("%s"+target.targetExistsMessage()+" and no specific task or rule has been defined", debugPrefix)
		return false, nil
	}

	if err != nil {
		// file does not exist and we can not find a task to create it
		return false, err
	}

	someDependencyUpdatedOrNewer := false
	dependencies := task.Dependencies(invokedName)
	if len(dependencies) == 0 {
		log.SetDebugIndent(3 * recursionDepth)
	} else {
		log.Debug("%sResolving task [<strong>%s</strong>] => %v", debugPrefix, invokedName, dependencies)
		someDependencyUpdatedOrNewer, err = checkDependencies(target, dependencies, recursionDepth)
		if err != nil {
			return false, err
		}
	}

	targetWasUpdated := false
	if target.ExistingFile && someDependencyUpdatedOrNewer == false {
		log.Debug("No need to build target [<strong>%s</strong>]", invokedName)
	} else {
		argumentsMessage := ""
		if len(args) > 0 {
			argumentsMessage = fmt.Sprintf("with args %v ", args)
		}
		message := fmt.Sprintf("%sInvoking task [<strong>%s</strong>] %son %T", debugPrefix, invokedName, argumentsMessage, task)
		if someDependencyUpdatedOrNewer {
			message = message + " (dependencies updated or newer)"
		}
		log.Debug(message)
		targetWasUpdated, err = task.Invoke(invokedName, args...)
		if err != nil {
			return false, err
		}
	}

	hooksUpdated, err := checkHooks(HookAfter, invokedName, recursionDepth)
	if err != nil {
		return false, err
	}
	return targetWasUpdated || hooksUpdated, nil
}

func checkDependencies(target *Target, dependencies []string, recursionDepth int) (bool, error) {
	log.SetDebugIndent(3 * recursionDepth)

	someDependencyUpdatedOrNewer := false
	for _, dependency := range dependencies {
		depInfo := FileInfo(dependency)
		if target.ExistingFile && depInfo.ExistingFile && depInfo.ModificationTime.After(target.ModificationTime) {
			log.Debug("Dependency [<strong>%s</strong>] is newer than [<strong>%s</strong>] so that needs to be rebuild", dependency, target.Name)
			someDependencyUpdatedOrNewer = true
		}

		dependencyUpdated, err := checkDependency(depInfo, recursionDepth)
		if err != nil {
			return false, err
		}
		if dependencyUpdated {
			log.Debug("Dependency [<strong>%s</strong>] has been updated so [<strong>%s</strong>] needs to be rebuild", dependency, target.Name)
			someDependencyUpdatedOrNewer = true
		}
	}

	log.SetDebugIndent(3 * recursionDepth)
	return someDependencyUpdatedOrNewer, nil
}

func checkDependency(dependency *Target, recursionDepth int) (bool, error) {
	if _, alreadyInvoked := resolvedDependencies[dependency.Name]; alreadyInvoked == true {
		log.Debug("Skipping dependency [<strong>%s</strong>] (already resolved)", dependency.Name)
		return false, nil
	}

	if dependency.ExistingFile {
		if depTask, _ := GetTask(dependency.Name); depTask != nil && len(depTask.Dependencies(dependency.Name)) == 0 {
			log.Debug("Skipping dependency [<strong>%s</strong>] (%s exists and has no sub dependencies)", dependency.Name, dependency.Typ())
			return false, nil
		}
	}
	return InvokeTask(dependency.Name, recursionDepth+1)
}
