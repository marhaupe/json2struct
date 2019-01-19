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
	file *os.File
}

// New initializes an Editor instance ready to be used
func New() *Editor {
	e := &Editor{}
	file, err := os.OpenFile("json2struct.temp", os.O_RDWR|os.O_CREATE, 0666)
	e.file = file
	e.err = err
	return e
}

// Display spawns a new Vim process and pipes stdin, stdout and stderr to it
func (e *Editor) Display() {
	cmd := exec.Command("vim", e.file.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// Consume consumes the content in the created file and closes and deletes it
func (e *Editor) Consume() (string, error) {
	filename := e.file.Name()
	defer os.Remove(filename)
	defer e.file.Close()
	if e.err != nil {
		return "", e.err
	}
	b, err := ioutil.ReadAll(e.file)
	return string(b), err
}
