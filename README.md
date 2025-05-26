# cmd-agent

`cmd-agent` 是一个使用自然语言执行命令的工具。它利用 AI 理解命令需求并与 shell 交互。

目前仅支持使用 Gemini API 驱动

## Description

自然语言处理: 使用 AI 模型理解您的指令。

命令执行: 能够根据您的提示执行相应的系统命令。

可扩展工具: 支持通过 [`ai/tools`](ai/tools/) 目录添加自定义工具。

## Build

```sh
git clone https://github.com/jiny3/cmd-agent.git
cd cmd-agent

go build -o bin/x .
# 将 bin/x 移至环境变量中的文件夹下
```

## Usage

> 首次运行将会创建 `~/.config/xAI/config.yaml` , 需填写相应配置信息（如下）
>
> ```yaml
> # Optional
> log:
>   path: default.log
>   level: info
> 
> api:
>   # Required
>   key: <your-api-key>
>   # Optional
>   # api.model will be gemini-2.0-flash if not set
>   model: gemini-2.0-flash
> ```

```sh
x "<your prompt>"
```

## Architecture

```sh
.
├── ai/                   # AI 相关逻辑
│   ├── client/           # AI 服务客户端
│   └── tools/            # AI 可使用的工具 (例如: ai/tools/cmdexecutor.go)
├── bin/                  
├── cmd/                  # Cobra 命令定义
├── utils/                # 通用工具函数
├── go.mod                
├── go.sum                
├── main.go               # 程序主入口点
└── README.md             
```

## Advanced

本工具支持通过在 tools 目录下添加新的 Go 文件来扩展其功能。每个工具通常会定义一个或多个 genai.FunctionDeclaration，并提供相应的处理函数。

例如，ai/tools/cmdexecutor.go 实现了一个名为 cmd_executor 的工具，允许 AI 执行任意 shell 命令。

要添加新工具：

1. 在 tools 目录下创建一个新的 .go 文件。该文件中定义您的函数声明和处理逻辑。
2. 在 `ai/tools/registry.go` 中注册，以便 AI 客户端可以发现并使用它。

## TODO

- [ ] funtion tool 通过 mcp 进一步解耦
- [ ] 支持使用自部署模型
- [ ] 支持 windows powershell