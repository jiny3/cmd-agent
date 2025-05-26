package client

import (
	"google.golang.org/genai"
)

var systemContents = []*genai.Content{
	{
		Role: genai.RoleUser,
		Parts: []*genai.Part{
			{
				// A system prompt to instruct the AI model about the task
				Text: taskPrompt,
			},
		},
	},
}

func FormatFunctionCallResponse(call *genai.FunctionCall, resp map[string]any) *genai.Content {
	return &genai.Content{
		Role: genai.RoleUser,
		Parts: []*genai.Part{
			{
				FunctionResponse: &genai.FunctionResponse{
					ID:       call.ID,
					Name:     call.Name,
					Response: resp,
				},
			},
		},
	}
}
