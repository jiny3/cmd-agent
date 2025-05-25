package client

import (
	"context"

	"github.com/jiny3/cmd-agent/ai/tools"
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
	hookx.Init(&hookx.WithDefault)
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
	systemContent := &genai.Content{
		Role: genai.RoleUser,
		Parts: []*genai.Part{
			{
				// A system prompt to instruct the AI model about the task
				Text: "You are a shell expert, please help me complete the following command, you should output the completed command, no need to include any other explanation. Do not put completed command in a code block. Then execute the command by cmd_executor function call.",
			},
		},
	}
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
		return "", err
	}

	// handle function calls
	for len(result.FunctionCalls()) > 0 {
		_prompt = append(_prompt, result.Candidates[0].Content)
		for _, call := range result.FunctionCalls() {
			// TODO: handle function call
			logrus.WithField("name", call.Name).WithField("args", call.Args).Debug("function call")
			if call.Name == tool[0].FunctionDeclarations[0].Name {
				output, stderr, _err := tools.CmdExecutor(call.Args["command"].(string))
				errResp := ""
				if _err != nil {
					errResp = _err.Error()
				}
				_prompt = append(_prompt, &genai.Content{
					Role: genai.RoleUser,
					Parts: []*genai.Part{
						{
							FunctionResponse: &genai.FunctionResponse{
								ID:   call.ID,
								Name: call.Name,
								Response: map[string]any{
									"output": output,
									"error":  stderr + errResp,
								},
							},
						},
					},
				})
			}
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
