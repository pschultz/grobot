package grobot

import (
	"bytes"
	"fmt"
	"github.com/fgrosse/grobot/log"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func init() {
	ShellProvider = &SystemShell{}
}

type Shell interface {
	Execute(cmdLine string, silent bool) (string, error)
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
func Execute(format string, args ...interface{}) string {
	cmdLine := fmt.Sprintf(format, args...)
	log.Shell(cmdLine)
	output, err := ShellProvider.Execute(cmdLine, false)
	if err != nil {
		log.Error("Output from shell:")
		log.Print(output)
		panic(err)
	}
	return strings.TrimSpace(output)
}

// ExecuteSilent does exactly what Execute does but only printing the
// shell command and output if debug mode is enabled
func ExecuteSilent(format string, args ...interface{}) string {
	if isDebug {
		return Execute(format, args...)
	}

	cmdLine := fmt.Sprintf(format, args...)
	output, err := ShellProvider.Execute(cmdLine, true)
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(output)
}

type SystemShell struct {
	Dir string
}

func (s *SystemShell) SetWorkingDirectory(workingDirectory string) {
	s.Dir = workingDirectory
}

func (s *SystemShell) Execute(cmdLine string, silent bool) (string, error) {
	cmdParts := strings.Split(cmdLine, " ")
	for i, part := range cmdParts {
		cmdParts[i] = strings.Trim(part, `"`)
	}

	var stdOut, stdErr *shellWriter
	if silent {
		stdOut = newShellWriter(ioutil.Discard)
		stdErr = newShellWriter(ioutil.Discard)
	} else {
		stdOut = newShellWriter(os.Stdout)
		stdErr = newShellWriter(os.Stderr)
	}

	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	cmd.Stdout = stdOut
	cmd.Stderr = stdErr
	cmd.Stdin = os.Stdin
	cmd.Dir = s.Dir

	err := cmd.Run()
	return stdOut.Output(), err
}

type shellWriter struct {
	output io.Writer
	buf    *bytes.Buffer
}

func newShellWriter(output io.Writer) *shellWriter {
	return &shellWriter{output, &bytes.Buffer{}}
}

func (w *shellWriter) Write(p []byte) (int, error) {
	n, err := w.buf.Write(p)
	if err != nil {
		return n, err
	}

	m, err := w.output.Write(p)
	if err != nil {
		return m, err
	}

	if n != m {
		return 0, fmt.Errorf("Error while writing shell output to internal buffer: wrote %d bytes to buffer and %d bytes to the output", n, m)
	}

	return n, err
}

func (w *shellWriter) Output() string {
	return w.buf.String()
}
