package main

import (
	"flag"
	"fmt"
	"github.com/fgrosse/gobot"
	"github.com/fgrosse/gobot/log"

	_ "github.com/fgrosse/gobot/modules/gomock"
	"os"
)

const defaultConfigFile = "gobot.json"

var debug = flag.Bool("debug", false, "only useful in gobot development")
var configFile = flag.String("config", defaultConfigFile, "set the used config file")
var showTasks = flag.Bool("t", false, "Display available tasks with descriptions, then exit.")

func main() {
	flag.Parse()
	if *debug == false {
		defer panicHandler()
	} else {
		log.EnableDebug()
		log.Debug("Running in gobot debug mode")
	}

	if err := gobot.LoadConfigFromFile(*configFile); err != nil {
		log.Fatal(err.Error())
	}

	if *showTasks {
		gobot.PrintTasks()
		os.Exit(0)
	}

	taskName := "mocks" // TODO change this to `default`
	if err := gobot.InvokeTask(taskName); err != nil {
		log.Fatal(err.Error())
	}
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
