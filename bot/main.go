package main

import (
	"flag"
	"fmt"
	"github.com/fgrosse/grobot"
	"github.com/fgrosse/grobot/log"

	// import modules
	_ "github.com/fgrosse/grobot/modules/dependency"
	_ "github.com/fgrosse/grobot/modules/ginkgo"
	_ "github.com/fgrosse/grobot/modules/gomock"

	"os"
	"strings"
)

const defaultConfigFile = "grobot.json"

var debug = flag.Bool("debug", false, "show a lot more debug information on the tasks")
var configFile = flag.String("config", defaultConfigFile, "set the used config file")
var showTasks = flag.Bool("t", false, "Display available tasks with descriptions, then exit.")

func main() {
	flag.Parse()
	if *debug == false {
		defer panicHandler()
	} else {
		grobot.EnableDebugMode()
		log.Debug("Running in grobot debug mode")
	}

	if err := grobot.LoadConfigFromFile(*configFile); err != nil {
		log.Fatal(err.Error())
	}

	if *showTasks {
		grobot.PrintTasks()
		os.Exit(0)
	}

	taskName := "default"
	args := filterArgs()
	if len(args) == 0 {
		log.Action("No target given. Please give me a target as first parameter.")
		log.Print("Usage: bot <task> [parameters]\n")
		log.Print("The following tasks are available:")
		grobot.PrintTasks()
		log.Print("")
		os.Exit(1)
	}

	taskName = args[0]
	somethingWasDone, err := grobot.InvokeTask(taskName, 0, args[1:]...)
	if err != nil {
		log.Fatal(err.Error())
	}

	if somethingWasDone == false {
		log.Debug("Task [<strong>%s</strong>] is up to date", taskName)
	} else {
		log.Debug("Task [<strong>%s</strong>] has been updated", taskName)
	}
}

func filterArgs() []string {
	args := []string{}
	for _, a := range os.Args[1:] {
		if strings.HasPrefix(a, "-") {
			continue
		}
		args = append(args, a)
	}

	log.Debug("Args: %v", args)
	return args
}

func panicHandler() {
	if r := recover(); r != nil {
		var err error
		switch caughtErr := r.(type) {
		case error:
			err = fmt.Errorf("Caught unexpected panic: %s", caughtErr.Error())
		case string:
			err = fmt.Errorf("Caught unexpected panic: %s", caughtErr)
		default:
			err = fmt.Errorf("Unknown unexpected panic occurred")
		}
		log.Fatal(err.Error())
	}
}
