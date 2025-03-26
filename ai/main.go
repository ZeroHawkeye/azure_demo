package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// Azure OpenAI配置
type Config struct {
	Endpoint   string
	ApiKey     string
	Deployment string
	ApiVersion string
}

// 聊天消息结构
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// 请求体结构
type RequestBody struct {
	Messages            []Message `json:"messages"`
	MaxCompletionTokens int       `json:"max_completion_tokens"`
	Temperature         float64   `json:"temperature"`
}

// 响应结构
type ResponseChoice struct {
	Message Message `json:"message"`
}

type Response struct {
	Choices []ResponseChoice `json:"choices"`
}

func main() {
	// 从环境变量获取配置
	azureOpenAIEndpoint := os.Getenv("AZURE_OPENAI_ENDPOINT")
	azureOpenAIKey := os.Getenv("AZURE_OPENAI_API_KEY")
	deploymentName := os.Getenv("AZURE_OPENAI_DEPLOYMENT")

	// 检查必要的环境变量
	if azureOpenAIEndpoint == "" || azureOpenAIKey == "" || deploymentName == "" {
		fmt.Println("请设置以下环境变量:")
		fmt.Println("AZURE_OPENAI_ENDPOINT - Azure OpenAI服务端点")
		fmt.Println("AZURE_OPENAI_API_KEY - Azure OpenAI API密钥")
		fmt.Println("AZURE_OPENAI_DEPLOYMENT - Azure OpenAI部署名称")
		return
	}

	// 使用基于API密钥的身份验证来初始化OpenAI客户端
	cred := azcore.NewKeyCredential(azureOpenAIKey)
	client, err := azopenai.NewClientWithKeyCredential(azureOpenAIEndpoint, cred, nil)
	if err != nil {
		log.Fatalf("初始化客户端错误: %s", err)
	}

	// 存储对话历史
	var messages []azopenai.ChatRequestMessageClassification

	// 添加系统消息
	systemMessage := azopenai.ChatRequestSystemMessage{
		Content: azopenai.NewChatRequestSystemMessageContent("你是一个有用的AI助手，可以回答用户的问题。"),
	}
	messages = append(messages, &systemMessage)

	fmt.Println("欢迎使用Azure OpenAI聊天工具！")
	fmt.Println("输入'exit'或'quit'退出程序")
	fmt.Println("------------------------------")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("用户: ")
		if !scanner.Scan() {
			break
		}

		userInput := scanner.Text()
		if userInput == "exit" || userInput == "quit" {
			break
		}

		// 添加用户消息到历史
		userMessage := azopenai.ChatRequestUserMessage{
			Content: azopenai.NewChatRequestUserMessageContent(userInput),
		}
		messages = append(messages, &userMessage)

		// 设置最大令牌数；这里可以自由调整，根据模型不同限制是不同的
		maxTokens := int32(800)

		// 发出聊天完成请求
		resp, err := client.GetChatCompletions(
			context.TODO(),
			azopenai.ChatCompletionsOptions{
				Messages:            messages,
				DeploymentName:      &deploymentName,
				MaxCompletionTokens: &maxTokens,
			},
			nil,
		)
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			continue
		}

		// 处理响应
		if len(resp.Choices) > 0 && resp.Choices[0].Message != nil && resp.Choices[0].Message.Content != nil {
			assistantContent := *resp.Choices[0].Message.Content
			fmt.Printf("AI: %s\n", assistantContent)

			// 添加助手回复到历史
			assistantMessage := azopenai.ChatRequestAssistantMessage{
				Content: azopenai.NewChatRequestAssistantMessageContent(assistantContent),
			}
			messages = append(messages, &assistantMessage)
		} else {
			fmt.Println("AI: 抱歉，我无法生成回复。")
		}
	}
}
