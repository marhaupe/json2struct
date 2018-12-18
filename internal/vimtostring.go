package internal

import (
	"io/ioutil"
	"os"
	"os/exec"
)

func VimToString(filename string) (string, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer file.Close()
	defer os.Remove(filename)

	cmd := exec.Command("vim", file.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	b, _ := ioutil.ReadAll(file)
	return string(b), nil
}
