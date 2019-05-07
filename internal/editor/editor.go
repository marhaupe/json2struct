// Package editor contains fields to create, edit, save and then delete files
// with currently only vim
package editor

import (
	"io/ioutil"
	"os"
	"os/exec"
)

// Editor contains fields to display Vim, write contents to a file
// and consume that file and it's contents
type Editor struct {
	err  error
	cmd  *exec.Cmd
	file *os.File
}

// New initializes an Editor instance ready to be used
func New() *Editor {
	e := &Editor{}
	e.setupFile()
	e.setupCmd()
	return e
}

func (e *Editor) setupFile() {
	e.file, e.err = os.OpenFile("json2struct.temp", os.O_RDWR|os.O_CREATE, 0666)
}

func (e *Editor) setupCmd() {
	e.cmd = exec.Command("vim", e.file.Name())
	e.cmd.Stdin = os.Stdin
	e.cmd.Stdout = os.Stdout
	e.cmd.Stderr = os.Stderr
}

// Display spawns a new Vim process and pipes stdin, stdout and stderr to it
func (e *Editor) Display() {
	if e.cmd.Process != nil {
		e.setupCmd()
	}
	e.cmd.Run()
}

// Consume consumes the content in the created file and closes and deletes it
func (e *Editor) Consume() (string, error) {
	defer e.Delete()
	return e.Read()
}

func (e *Editor) Read() (string, error) {
	if e.err != nil {
		return "", e.err
	}
	b, err := ioutil.ReadAll(e.file)
	return string(b), err
}

func (e *Editor) Delete() error {
	filename := e.file.Name()
	return os.Remove(filename)
}
