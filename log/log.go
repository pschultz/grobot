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
var debugIndent int
var currentColor = defaultStyle

func EnableDebug() {
	isDebug = true
}

func Print(format string, args ...interface{}) {
	format = strings.Replace(format, "<strong>", boldStyle, -1)
	format = strings.Replace(format, "</strong>", defaultStyle+currentColor, -1)
	fmt.Printf(format+"\n", args...)
}

func Debug(format string, args ...interface{}) {
	if isDebug {
		coloredOutputLn("[DEBUG] "+strings.Repeat(" ", debugIndent)+format, grayColor, args...)
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
	coloredOutputLn("➤ "+format, yellowColor, args...)
}

func Shell(format string, args ...interface{}) {
	coloredOutputLn(boldStyle+"$ "+defaultStyle+format, lightGrayColor, args...)
}

func coloredOutput(format, color string, args ...interface{}) {
	fmt.Printf(color+format+defaultStyle, args...)
}

func coloredOutputLn(format, color string, args ...interface{}) {
	currentColor = color
	Print(color+format+defaultStyle, args...)
	currentColor = defaultStyle
}

func AskBool(question string, args ...interface{}) bool {
	var input string
	for input != "y" && input != "n" {
		coloredOutput("➤ "+question+" [Yn] ", yellowColor, args...)
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

func SetDebugIndent(n int) {
	debugIndent = n
}
