package grobot

import (
	"bytes"
	"fmt"
	"github.com/fgrosse/grobot/log"
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
	SetWorkingDirectory(workingDirectory string)
}

// ShellProvider is mainly there to enable bot users to exchange the used shell
// This is especially useful in grobots own tests to mock the shell
var ShellProvider Shell

// SetWorkingDirectory changes the current working directory for each subsequent call to Execute
// Note that this is a permanent change so you probably need to use ResetWorkingDirectory when you are done
func SetWorkingDirectory(format string, args ...interface{}) {
	workingDir := fmt.Sprintf(format, args...)
	ShellProvider.SetWorkingDirectory(workingDir)
}

// ResetWorkingDirectory sets the current working directory of the shell to its initial state
// This is equivalent to the call SetWorkingDirectory("") but much nicer to read
func ResetWorkingDirectory() {
	SetWorkingDirectory("")
}

// Execute a command on the shell
// The format and arguments are directly piped into fmt.Sprintf
//
// NOTE: Execute will panic if an error occurred.
// This way you can write multiple Executes right after one another without
// having to handle each and every error separately
// Defer some panic handler if you need to recover from the error.
// Grobots main function will handle the panic if nobody else does it
func Execute(format string, args ...interface{}) {
	cmdLine := fmt.Sprintf(format, args...)
	log.Shell(cmdLine)
	err := ShellProvider.Execute(cmdLine)
	if err != nil {
		panic(err)
	}
}

type SystemShell struct {
	Dir string
}

func (s *SystemShell) SetWorkingDirectory(workingDirectory string) {
	s.Dir = workingDirectory
}

func (s *SystemShell) Execute(cmdLine string) error {
	cmdParts := strings.Split(cmdLine, " ")
	for i, part := range cmdParts {
		cmdParts[i] = strings.Trim(part, `"`)
	}

	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = s.Dir
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
