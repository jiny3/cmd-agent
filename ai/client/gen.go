package client

import (
	"context"

	"github.com/jiny3/cmd-agent/ai/tools"
	"github.com/jiny3/cmd-agent/utils"
	"github.com/jiny3/gopkg/hookx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/genai"
)

var (
	apiKey   string
	modelId  string
	aiClient *genai.Client
)

func init() {
	hookx.Init(&utils.Init)
	apiKey = viper.GetString("api.key")
	if apiKey == "" {
		logrus.Fatal("api key is empty")
	}
	modelId = viper.GetString("api.model")
	if modelId == "" {
		modelId = "gemini-2.0-flash"
		logrus.Warn("model id not found, using default model id: gemini-2.0-flash")
	}
	var err error
	aiClient, err = genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})

	if err != nil {
		logrus.Fatal(err)
	}
}

func GenerateContent(prompt string, tool ...*genai.Tool) (string, error) {
	_prompt := []*genai.Content{systemContent}
	_prompt = append(_prompt, genai.Text(prompt)...)
	result, err := aiClient.Models.GenerateContent(
		context.Background(),
		modelId,
		_prompt,
		&genai.GenerateContentConfig{
			Tools: tool,
		},
	)
	if err != nil {
		logrus.WithError(err).Debug("generate content failed")
		return "", err
	}

	// handle function calls
	for len(result.FunctionCalls()) > 0 {
		_prompt = append(_prompt, result.Candidates[0].Content)
		for _, call := range result.FunctionCalls() {
			logrus.WithField("name", call.Name).WithField("args", call.Args).Debug("function call")

			// get function handler
			handle, exist := tools.GetHandler(call.Name)
			if !exist {
				_prompt = append(_prompt, FormatFunctionCallResponse(call, map[string]any{"error": "function not found"}))
			}

			// execute function
			errResp, outStr := "", ""
			out, _err := handle(call.Args["command"], call.Args["timeout"])
			if _err != nil {
				errResp = _err.Error()
			} else {
				outStr = (out.([]string))[0]
				errResp += (out.([]string))[1]
			}
			logrus.WithField("output", outStr).WithField("error", errResp).Debug("function call response")
			_prompt = append(_prompt, FormatFunctionCallResponse(call, map[string]any{
				"output": outStr,
				"error":  errResp,
			}))
		}
		result, err = aiClient.Models.GenerateContent(
			context.Background(),
			modelId,
			_prompt,
			&genai.GenerateContentConfig{
				Tools: tool,
			},
		)
		if err != nil {
			return "", err
		}
	}

	return result.Text(), nil
}
