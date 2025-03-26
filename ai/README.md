# Azure OpenAI CLI 聊天工具

这是一个简单的命令行聊天工具，使用Azure的认知服务OpenAI接口实现对话功能。

## 前提条件

- Go 编程语言环境
- Azure OpenAI 服务账号和API密钥

## 设置环境变量

在运行程序前，需要设置以下环境变量：

```
AZURE_OPENAI_ENDPOINT - Azure OpenAI服务端点 (例如: https://your-resource-name.openai.azure.com)
AZURE_OPENAI_API_KEY - Azure OpenAI API密钥
AZURE_OPENAI_DEPLOYMENT - Azure OpenAI部署名称 (例如: gpt-35-turbo)
```

在Windows环境中，可以使用以下命令设置环境变量：

```
$env:AZURE_OPENAI_ENDPOINT = "https://your-resource-name.openai.azure.com"
$env:AZURE_OPENAI_API_KEY = "your-api-key"
$env:AZURE_OPENAI_DEPLOYMENT = "your-deployment-name"
```

## 编译和运行

```
cd ai
go build
./ai
```

或者直接运行：

```
cd ai
go run main.go
```

## 使用方法

启动程序后，可以直接在命令行中输入问题，按回车发送。程序会将问题发送到Azure OpenAI服务，并显示回复。

输入 `exit` 或 `quit` 可以退出程序。

## 注意事项

- 程序会保存对话历史，并在每次请求中发送完整的对话历史
- 默认设置了最大令牌数为800
- 请确保您的Azure OpenAI服务已正确配置并部署 