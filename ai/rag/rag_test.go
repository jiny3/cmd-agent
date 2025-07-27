// ai/rag/rag_test.go
package rag

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type dummyProvider struct {
	content string
	err     error
}

func (d dummyProvider) Provide() (string, error) {
	return d.content, d.err
}

func Test_GetContextMessages_Fallback(t *testing.T) {
	dir := t.TempDir()
	histPath := filepath.Join(dir, "history.txt")
	want := "line1\nline2"
	if err := os.WriteFile(histPath, []byte(want), 0o644); err != nil {
		t.Fatalf("failed to write history file: %v", err)
	}

	RegisterProvider("default", fileProvider{path: histPath})

	os.Setenv("SHELL", "/usr/bin/unknownshell")

	msgs := GetContextMessages()
	if msgs == nil {
		t.Fatal("expected non-nil messages for fallback provider")
	}

	header := fmt.Sprintf("The following is the %s history for your reference:\n", "unknownshell")
	found := false
	for _, msg := range msgs {
		var fullText string
		for _, part := range msg.Parts {
			fullText += part.Text
		}
		if strings.Contains(fullText, header+want) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected message containing %q, got: %+v", header+want, msgs)
	}
}

func Test_GetContextMessages_ErrorOrEmpty(t *testing.T) {
	RegisterProvider("errshell", dummyProvider{"", os.ErrNotExist})
	os.Setenv("SHELL", "/bin/errshell")

	if msgs := GetContextMessages(); msgs != nil {
		t.Errorf("expected nil for error provider, got: %+v", msgs)
	}

	RegisterProvider("emptyshell", dummyProvider{"", nil})
	os.Setenv("SHELL", "/usr/local/bin/emptyshell")

	if msgs := GetContextMessages(); msgs != nil {
		t.Errorf("expected nil for empty-content provider, got: %+v", msgs)
	}
}
