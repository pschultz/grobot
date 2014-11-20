package main

import (
	"github.com/fgrosse/grobot"
	"github.com/fgrosse/grobot/log"

	// import modules
	_ "github.com/fgrosse/grobot/modules/dependency"
	_ "github.com/fgrosse/grobot/modules/ginkgo"
	_ "github.com/fgrosse/grobot/modules/gomock"

	"flag"
	"fmt"
	"os"
	"strings"
)

var BotVersion = grobot.NewVersion("0.7")

var debug = flag.Bool("debug", false, "show a lot more debug information on the tasks")
var configFile = flag.String("config", grobot.DefaultConfigFileName, "set the used config file")
var showTasks = flag.Bool("t", false, "Display available tasks with descriptions, then exit.")
var showVersion = flag.Bool("version", false, "Display the current version of bot, then exit.")
var showHelp = flag.Bool("help", false, "Display this help, then exit.")

func main() {
	flag.Parse()
	if *debug == false {
		defer panicHandler()
	} else {
		grobot.EnableDebugMode()
		log.Debug("Running in grobot debug mode")
	}

	parseOutputFlags()
	loadConfigurationFile()
	invokeTask()
}

func parseOutputFlags() {
	if *showTasks {
		grobot.PrintTasks()
		os.Exit(0)
	}

	if *showVersion {
		printVersion()
		os.Exit(0)
	}

	if *showHelp {
		showHelpText()
		os.Exit(0)
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

func printVersion() {
	log.Action("You are running bot version %s", BotVersion)
}

func showHelpText() {
	log.Action("Bot is an automation and build tool for the go programming language.")
	log.Print(`Version: %s`, BotVersion)
	log.Print(`Usage:   <strong>bot</strong> [optonal flags] <task> [optional arguments]`)
	log.Print(`Example: bot -debug install code.google.com/p/gomock`)
	log.Print(``)
	log.Print(`Which tasks are available depends on the used configuration file.`)
	log.Print(``)
	log.Print(`<strong>The following tasks are available with the default configuration file (%s):</strong>`, defaultConfigFile)
	grobot.PrintTasks()
	log.Print(``)
}

func loadConfigurationFile() {
	if *configFile == defaultConfigFile {
		file := grobot.TargetInfo(*configFile)
		if file.ExistingFile == false {
			log.Debug("Default configuration file %S does not exist", defaultConfigFile)
			grobot.LoadBuiltinConfig()
			return
		}
	}

	if err := grobot.LoadConfigFromFile(*configFile, BotVersion); err != nil {
		log.Fatal(err.Error())
	}
}

func invokeTask() {
	taskName := "default"
	args := filterArgs()
	if len(args) == 0 {
		showHelpText()
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
