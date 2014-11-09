package log

import (
	"fmt"
	"os"
	"strings"
)

const defaultStyle = "\x1b[0m"
const boldStyle = "\x1b[1m"
const redColor = "\x1b[91m"
const greenColor = "\x1b[32m"
const yellowColor = "\x1b[33m"
const cyanColor = "\x1b[36m"
const grayColor = "\x1b[90m"
const lightGrayColor = "\x1b[37m"

var isDebug = false

func EnableDebug() {
	isDebug = true
}

func Print(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

func Debug(format string, args ...interface{}) {
	if isDebug {
		coloredOutputLn("[DEBUG] "+format, grayColor, args...)
	}
}

func Fatal(format string, args ...interface{}) {
	Error(format, args...)
	os.Exit(1)
}

func Error(format string, args ...interface{}) {
	coloredOutputLn("ERROR: "+format, redColor, args...)
}

func Action(format string, args ...interface{}) {
	coloredOutputLn(format, yellowColor, args...)
}

func Shell(format string, args ...interface{}) {
	coloredOutputLn("$ "+format, lightGrayColor+boldStyle, args...)
}

func coloredOutput(format, color string, args ...interface{}) {
	fmt.Printf(color+format+defaultStyle, args...)
}

func coloredOutputLn(format, color string, args ...interface{}) {
	fmt.Printf(color+format+defaultStyle+"\n", args...)
}

func AskBool(question string, args ...interface{}) bool {
	var input string
	for input != "y" && input != "n" {
		coloredOutput(question, yellowColor, args...)
		_, err := fmt.Scanf("%s", &input)
		if err != nil {
			panic(fmt.Errorf("Could not read input : %s", err.Error()))
		}

		originalInput := input
		input = strings.ToLower(input)
		if input != "y" && input != "n" {
			Error("Invalid input: %s", originalInput)
		}
	}

	return input == "y"
}
