package grobot

import (
	"fmt"
	"github.com/fgrosse/grobot/log"
	"regexp"
)

var (
	tasks                = map[string]Task{}
	rules                = map[*regexp.Regexp]Task{}
	resolvedDependencies = map[string]bool{}
)

type Task interface {
	Dependencies(invokedName string) []string
	Invoke(name string) (bool, error)
}

type nullTask struct{}

func (t *nullTask) Dependencies(invokedName string) []string { return []string{} }
func (t *nullTask) Invoke(name string) error                 { return nil }

type Describable interface {
	Description() string
}

// Reset is used to make grobot forget about all registered tasks and rules
// This is probably only useful in tests
func Reset() {
	tasks = map[string]Task{}
	rules = map[*regexp.Regexp]Task{}
	resolvedDependencies = map[string]bool{}
}

func RegisterTask(name string, newTask Task) error {
	if _, keyExists := tasks[name]; keyExists == true {
		return fmt.Errorf("Module error: Task '%s' has already been registered", name)
	}

	log.Debug("Registering [%s] as %T", name, newTask)
	tasks[name] = newTask
	return nil
}

func RegisterRule(ruleRegex string, newTask Task) error {
	rule, err := regexp.Compile(ruleRegex)
	if err != nil {
		return fmt.Errorf("Could not compile rule regex: %s", err.Error())
	}

	log.Debug("Registering rule /%s/ as %T", ruleRegex, newTask)
	rules[rule] = newTask
	return nil
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
	for name, task := range tasks {
		switch t := task.(type) {
		case Describable:
			description := t.Description()
			if description != "" {
				fmt.Printf("%s : %s\n", name, description)
			}
		}
	}
}

func InvokeTask(invokedName string, recursionDepth int) (bool, error) {
	resolvedDependencies[invokedName] = true
	targetIsExistentFile, targetIsDir, _, err := ModificationDate(invokedName)
	if err != nil {
		return false, err
	}

	task, err := GetTask(invokedName)
	if targetIsExistentFile && task == nil {
		log.Debug(targetExistsMessage(invokedName, targetIsDir) + " and no specific task or rule has been defined")
		return false, nil
	}

	if err != nil {
		// file does not exist and can not find task to create it
		return false, err
	}

	debugPrefix := ""
	log.SetDebugIndent(0)
	if recursionDepth > 0 {
		log.SetDebugIndent(3 * (recursionDepth - 1))
		debugPrefix = "┗━ "
	}

	dependencies := task.Dependencies(invokedName)
	if len(dependencies) == 0 {
		log.Debug("%sInvoking task [<strong>%s</strong>] with %T", debugPrefix, invokedName, task)
		log.SetDebugIndent(3 * recursionDepth)
		return task.Invoke(invokedName)
	}

	log.Debug("%sResolving task [<strong>%s</strong>] => %v", debugPrefix, invokedName, dependencies)
	someDependencyExecuted, err := checkDependencies(invokedName, dependencies, recursionDepth)
	if err != nil {
		return false, err
	}

	if targetIsExistentFile && someDependencyExecuted == false {
		log.Debug("No need to build target [<strong>%s</strong>]", invokedName)
		return false, nil
	} else {
		log.Debug("Invoking task [<strong>%s</strong>] => %T", invokedName, task)
		return task.Invoke(invokedName)
	}
}

func checkDependencies(invokedName string, dependencies []string, recursionDepth int) (bool, error) {
	log.SetDebugIndent(3 * recursionDepth)

	someDependencyExecuted := false
	for _, dependency := range dependencies {
		dependencyExecuted, err := checkDependency(dependency, recursionDepth)
		if err != nil {
			return false, err
		}
		if dependencyExecuted {
			log.Debug("Dependency [<strong>%s</strong>] has been executed so [<strong>%s</strong>] needs to be rebuild", dependency, invokedName)
			someDependencyExecuted = true
		}
	}

	log.SetDebugIndent(3 * recursionDepth)
	return someDependencyExecuted, nil
}

func checkDependency(depPath string, recursionDepth int) (bool, error) {
	if _, alreadyInvoked := resolvedDependencies[depPath]; alreadyInvoked == true {
		log.Debug("Skipping dependency [<strong>%s</strong>] (already resolved)", depPath)
		return false, nil
	}

	return InvokeTask(depPath, recursionDepth+1)
	/*
		// only invoke target if dep is not existent file or dep is newer than target
		depFileExists, depIsDir, depModDate, err := ModificationDate(dependency)
		if err != nil {
			return false, err
		}

		if depFileExists {
			nrOfSubDependencies := 0
			depTask, err := GetTask(dependency)
			if err == nil {
				nrOfSubDependencies = len(depTask.Dependencies(dependency))
			}

			depType := "File"
			if depIsDir {
				depType = "Folder"
			}
			debugMessage := fmt.Sprintf("%s [<strong>%s</strong>] does already exist", depType, dependency)
			if nrOfSubDependencies > 0 {
				debugMessage = fmt.Sprintf("%s but has own dependencies to check", debugMessage)
				log.Debug(debugMessage)
				invokeThisTask, err = InvokeTask(dependency, recursionDepth+1)
				if err != nil {
					return false, err
				}
			} else if targetIsExistentFile {
				if depModDate.After(targetModDate) {
					debugMessage = fmt.Sprintf("%s but is newer than [<strong>%s</strong>]", debugMessage, invokedName)
					invokeThisTask = true
				} else {
					debugMessage = fmt.Sprintf("%s and is older than [<strong>%s</strong>]", debugMessage, invokedName)
				}
				log.Debug(debugMessage)
			}
		} else {
			invokeThisTask, err = InvokeTask(dependency, recursionDepth+1)
			if err != nil {
				return false, err
			}
		}*/
	return false, nil
}

func targetExistsMessage(taskName string, isDir bool) string {
	fileType := "File"
	if isDir {
		fileType = "Folder"
	}
	return fmt.Sprintf("%s [<strong>%s</strong>] does already exist", fileType, taskName)
}
