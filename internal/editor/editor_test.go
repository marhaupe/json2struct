package editor

import (
	"fmt"
	"os"
	"testing"
)

func TestConsumeEmptyFile(t *testing.T) {
	editor := New()
	content, err := editor.Consume()
	editor.Delete()
	if content != "" {
		t.Error("Content is not empty but should be")
	}
	if err != nil {
		t.Error("Error is not nil but should be")
	}
}

func TestConsumeWithErrors(t *testing.T) {
	editor := New()
	editor.err = fmt.Errorf("Testerror")
	content, err := editor.Consume()
	if content != "" {
		t.Error("Content is not empty but should be")
	}
	if err.Error() != fmt.Errorf("Testerror").Error() {
		t.Error("Error has not been set properly")
	}
}

func TestNew(t *testing.T) {
	editor := New()
	if len(editor.cmd.Args) != 2 || editor.cmd.Args[1] != editor.file.Name() {
		t.Errorf("Cmd invoked with wrong args %v", editor.cmd.Args)
	}

	if editor.cmd.Stdout != os.Stdout ||
		editor.cmd.Stdin != os.Stdin ||
		editor.cmd.Stderr != os.Stderr {
		t.Error("Cmd pipes have been set wrong")
	}
	editor.Consume()
}

func Test_getUserEditor(t *testing.T) {
	os.Setenv("EDITOR", "nano")

	editor := getUserEditor()

	if editor != "nano" {
		t.Fatalf("want nano but got %v", editor)
	}

	os.Unsetenv("EDITOR")

	editor = getUserEditor()

	if editor != defaultEditor {
		t.Fatalf("want defaultEditor but got %v", editor)
	}
}
