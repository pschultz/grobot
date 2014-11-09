package gobot

import (
	"bytes"
	"fmt"
	"github.com/fgrosse/gobot/log"
	"io"
	"os"
	"os/exec"
	"strings"
)

func init() {
	ShellProvider = &SystemShell{}
}

type Shell interface {
	Execute(cmdLine string) error
}

// ShellProvider is mainly there to enable grobot users to exchange the used shell
// This is especially useful in grobots own tests to mock the shell
var ShellProvider Shell

func Execute(format string, args ...interface{}) error {
	cmdLine := fmt.Sprintf(format, args...)
	log.Shell(cmdLine)
	return ShellProvider.Execute(cmdLine)
}

type SystemShell struct{}

func (s *SystemShell) Execute(cmdLine string) error {
	cmdParts := strings.Split(cmdLine, " ")
	for i, part := range cmdParts {
		cmdParts[i] = strings.Trim(part, `"`)
	}

	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	cmd.Stdout = &ShellWriter{os.Stdout}
	cmd.Stderr = &ShellWriter{os.Stderr}
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

type ShellWriter struct {
	output io.Writer
}

func (w *ShellWriter) Write(p []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	buf.WriteString("> ")
	buf.Write(p)
	n, err = w.output.Write(buf.Bytes())
	if err != nil {
		return n, err
	}
	n = n - 2
	return
}
