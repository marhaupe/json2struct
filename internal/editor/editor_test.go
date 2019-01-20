package editor

import (
	"fmt"
	"testing"
)

func TestConsumeEmptyFile(t *testing.T) {
	editor := New()
	content, err := editor.Consume()
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
