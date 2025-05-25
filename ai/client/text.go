package client

import (
	"google.golang.org/genai"
)

var systemContent = &genai.Content{
	Role: genai.RoleUser,
	Parts: []*genai.Part{
		{
			// A system prompt to instruct the AI model about the task
			Text: "You are a shell expert, please help me complete the following command and set the appropriate timeout, you should output the completed command, no need to include any other explanation. Do not put completed command in a code block. Then execute the command by cmd_executor function call. If I denys the execution, the command may have some problems or the command will loop forever, you should improve the command to avoid the problems or loops. If you are not sure about the command, please ask user for more information.",
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
