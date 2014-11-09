package gobot

import (
	"fmt"
	"github.com/fgrosse/gobot/log"
	"os/exec"
	"strings"
)

func Shell(format string, args ...interface{}) error {
	cmdLine := fmt.Sprintf(format, args...)
	log.Shell(cmdLine)
	cmdParts := strings.SplitN(cmdLine, " ", 2)

	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	return cmd.Run()
}
