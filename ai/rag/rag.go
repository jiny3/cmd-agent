package rag

import (
	"fmt"
	"os"
	"path/filepath"

	"google.golang.org/genai"
)

type HistoryProvider interface {
	Provide() (string, error)
}

type fileProvider struct {
	path string
}

func (p fileProvider) Provide() (string, error) {
	data, err := os.ReadFile(p.path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

var providers = make(map[string]HistoryProvider)

func RegisterProvider(name string, prov HistoryProvider) {
	providers[name] = prov
}

func init() {
	home := os.Getenv("HOME")
	// Register providers for common shells
	RegisterProvider("bash", fileProvider{path: filepath.Join(home, ".bash_history")})
	RegisterProvider("zsh", fileProvider{path: filepath.Join(home, ".zsh_history")})
	RegisterProvider("fish", fileProvider{path: filepath.Join(home, ".local", "share", "fish", "fish_history")})
	RegisterProvider("default", fileProvider{path: filepath.Join(".", "history.txt")})
}

// GetContextMessages selects the appropriate Provider based on the current SHELL,
// reads the history, and returns message fragments that can be directly appended to a genai prompt.
func GetContextMessages() []*genai.Content {
	shellEnv := os.Getenv("SHELL")
	shellName := filepath.Base(shellEnv)

	prov, ok := providers[shellName]
	if !ok {
		prov = providers["default"]
	}

	raw, err := prov.Provide()
	if err != nil || raw == "" {
		return nil
	}
	header := fmt.Sprintf("The following is the %s history for your reference:\n", shellName)
	return genai.Text(header + raw)
}
