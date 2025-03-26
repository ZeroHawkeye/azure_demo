package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// TranslationRequest 表示翻译请求的结构
type TranslationRequest struct {
	Text string `json:"Text"`
}

// TranslationResponse 表示翻译响应的结构
type TranslationResponse []struct {
	Translations []struct {
		Text string `json:"text"`
		To   string `json:"to"`
	} `json:"translations"`
}

func main() {
	// 从环境变量获取Azure翻译服务的密钥和区域
	// 注意：在实际使用前，需要在Azure门户中创建翻译服务资源并获取这些值;这里我使用的环境变量进行获取，按需修改即可。
	subscriptionKey := os.Getenv("AZURE_TRANSLATOR_KEY")
	location := os.Getenv("AZURE_TRANSLATOR_REGION")

	// 检查必要的环境变量是否设置
	if subscriptionKey == "" || location == "" {
		log.Fatal("请设置AZURE_TRANSLATOR_KEY和AZURE_TRANSLATOR_REGION环境变量")
	}

	// 要翻译的文本
	textToTranslate := "Hello, world!"

	// 目标语言；目标语言需要设置。
	targetLanguage := "ko"

	// 调用翻译函数
	translatedText, err := translateText(textToTranslate, targetLanguage, subscriptionKey, location)
	if err != nil {
		log.Fatalf("翻译失败: %v", err)
	}

	fmt.Printf("原文: %s\n", textToTranslate)
	fmt.Printf("译文: %s\n", translatedText)
}

// translateText 使用Azure翻译服务将文本从一种语言翻译为另一种语言
// 参数:
//   - text: 要翻译的文本
//   - targetLang: 目标语言代码（如"zh-CN"表示简体中文）
//   - subscriptionKey: Azure翻译服务的订阅密钥
//   - location: Azure翻译服务的区域（如"eastasia"）
//
// 返回:
//   - 翻译后的文本
//   - 可能的错误
func translateText(text, targetLang, subscriptionKey, location string) (string, error) {
	// Azure翻译服务的API端点；这里的调用接口地址需要记录
	endpoint := "https://api.cognitive.microsofttranslator.com"

	// 构建完整的API URL;目标语言在路径中拼接
	uri := endpoint + "/translate?api-version=3.0&to=" + targetLang

	// 准备请求体
	body := []TranslationRequest{
		{
			Text: text,
		},
	}

	// 将请求体转换为JSON
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("JSON编码失败: %v", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 添加必要的请求头；这些header头是基础必须的。
	req.Header.Add("Ocp-Apim-Subscription-Key", subscriptionKey)
	req.Header.Add("Ocp-Apim-Subscription-Region", location)
	req.Header.Add("Content-Type", "application/json")

	// 发送HTTP请求；后续是简单的创建一个http请求。不一定要使用net/http库，按需使用即可。
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API返回非成功状态码: %d", resp.StatusCode)
	}

	// 读取响应体
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	// 解析JSON响应
	var translationResp TranslationResponse
	if err := json.Unmarshal(respBytes, &translationResp); err != nil {
		return "", fmt.Errorf("解析JSON响应失败: %v", err)
	}

	// 检查响应是否包含翻译结果
	if len(translationResp) == 0 || len(translationResp[0].Translations) == 0 {
		return "", fmt.Errorf("翻译响应为空")
	}

	// 返回翻译后的文本
	return translationResp[0].Translations[0].Text, nil
}
