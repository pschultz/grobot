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

func RegisterFolder(path string) error {
	task := CreateFolderTask{}
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
	target, err := TargetInfo(invokedName)
	if err != nil {
		return false, err
	}

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

	dependencies := task.Dependencies(invokedName)
	if len(dependencies) == 0 {
		log.Debug("%sInvoking task [<strong>%s</strong>] with %T", debugPrefix, invokedName, task)
		log.SetDebugIndent(3 * recursionDepth)
		return task.Invoke(invokedName)
	}

	log.Debug("%sResolving task [<strong>%s</strong>] => %v", debugPrefix, invokedName, dependencies)
	someDependencyUpdatedOrNewer, err := checkDependencies(target, dependencies, recursionDepth)
	if err != nil {
		return false, err
	}

	if target.ExistingFile && someDependencyUpdatedOrNewer == false {
		log.Debug("No need to build target [<strong>%s</strong>]", invokedName)
		return false, nil
	} else {
		message := fmt.Sprintf("Invoking task [<strong>%s</strong>] => %T", invokedName, task)
		if someDependencyUpdatedOrNewer {
			message = message + " (dependencies updated or newer)"
		}
		log.Debug(message)
		return task.Invoke(invokedName)
	}
}

func checkDependencies(target *Target, dependencies []string, recursionDepth int) (bool, error) {
	log.SetDebugIndent(3 * recursionDepth)

	someDependencyUpdatedOrNewer := false
	for _, dependency := range dependencies {
		depInfo, err := TargetInfo(dependency)
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
