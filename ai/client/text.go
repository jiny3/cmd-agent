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
				Text: "You are a shell expert, please help me with the following command, you should generate the correct command and then execute the command with the cmd_executor function call. If I refuse to execute or get an error in the execution response, the command may be faulty or will loop forever and you should improve the command to avoid them. If you are unsure about a command, ask the user for more information.Translated with www.DeepL.com/Translator (free version)",
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
