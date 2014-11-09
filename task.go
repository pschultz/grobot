package gobot

import (
	"fmt"
	"github.com/fgrosse/gobot/log"
)

var tasks = map[string]Task{}

type Task interface {
	Dependencies() []string
	Invoke(name string) error
}

type Describable interface {
	Description() string
}

func RegisterTask(name string, newTask Task) error {
	if _, keyExists := tasks[name]; keyExists == true {
		return fmt.Errorf("Module error: Task '%s' has already been registered", name)
	}

	log.Debug("Registering task %s", name)
	tasks[name] = newTask
	return nil
}

func GetTask(name string) (Task, error) {
	if task, taskExists := tasks[name]; taskExists == true {
		return task, nil
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

func InvokeTask(taskName string) error {
	task, err := GetTask(taskName)
	if err != nil {
		return err
	}

	var invokedDependencies map[string]bool
	dependencies := task.Dependencies()
	if len(dependencies) > 1 {
		log.Debug("There are %d dependencies: %v", len(dependencies), dependencies)
	} else if len(dependencies) == 1 {
		log.Debug("There is 1 dependency: %s", dependencies[0])
	}

	for _, dependency := range dependencies {
		if _, alreadyInvoked := invokedDependencies[dependency]; alreadyInvoked == true {
			log.Debug("Skipping dependency [%s] (already invoked)", dependency)
			continue
		}
		err := InvokeTask(dependency)
		if err != nil {
			return err
		}
	}

	log.Debug("Invoking task [%s]", taskName)
	return task.Invoke(taskName)
}
