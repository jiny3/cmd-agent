package tools

import (
	"sync"

	"github.com/sirupsen/logrus"
	"google.golang.org/genai"
)

var Registry = sync.Map{}

type functionHandler func(args ...any) (any, error)

func register(name string, handler functionHandler) {
	if _, exists := Registry.Load(name); exists {
		logrus.Warnf("Tool %s is already registered, overwriting it", name)
	}
	Registry.Store(name, handler)
}

func GetHandler(name string) (functionHandler, bool) {
	handler, exists := Registry.Load(name)
	if !exists {
		return nil, false
	}
	return handler.(functionHandler), true
}

// if add new tool, please also add it to ListTools function
func ListTools() []*genai.Tool {
	return []*genai.Tool{
		CmdExecutorTool(),
	}
}
