package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// ContentSafetyRequest 表示发送到Azure内容安全API的请求
type ContentSafetyRequest struct {
	Text string `json:"text"`
}

// ContentSafetyResponse 表示从Azure内容安全API返回的响应
type ContentSafetyResponse struct {
	CategoriesAnalysis []struct {
		Category string  `json:"category"`
		Severity float64 `json:"severity"`
	} `json:"categoriesAnalysis"`
	BlocklistsMatch []struct {
		BlocklistName     string `json:"blocklistName"`
		BlocklistItemId   string `json:"blocklistItemId"`
		BlocklistItemText string `json:"blocklistItemText"`
	} `json:"blocklistsMatch"`
}

func main() {
	// 设置Azure Content Safety API的端点和密钥;设置的环境变量在这里读取
	endpoint := os.Getenv("AZURE_CONTENT_SAFETY_ENDPOINT")
	apiKey := os.Getenv("AZURE_CONTENT_SAFETY_KEY")

	if endpoint == "" || apiKey == "" {
		fmt.Println("请设置环境变量: AZURE_CONTENT_SAFETY_ENDPOINT 和 AZURE_CONTENT_SAFETY_KEY")
		fmt.Println("例如:")
		fmt.Println("$env:AZURE_CONTENT_SAFETY_ENDPOINT = 'https://learnclown.cognitiveservices.azure.com/'")
		fmt.Println("$env:AZURE_CONTENT_SAFETY_KEY = 'your-api-key'")
		return
	}

	// 要检查的文本内容
	textToAnalyze := "这是一个测试文本，用于检测内容安全。"

	// 调用内容安全API
	result, err := analyzeText(endpoint, apiKey, textToAnalyze)
	if err != nil {
		fmt.Printf("分析文本时出错: %v\n", err)
		return
	}

	// 打印结果
	fmt.Println("内容安全分析结果:")
	fmt.Println("文本:", textToAnalyze)

	if len(result.CategoriesAnalysis) > 0 {
		fmt.Println("\n类别分析:")
		for _, category := range result.CategoriesAnalysis {
			fmt.Printf("- 类别: %s, 严重程度: %.2f\n", category.Category, category.Severity)
		}
	} else {
		fmt.Println("\n没有检测到任何类别问题")
	}

	if len(result.BlocklistsMatch) > 0 {
		fmt.Println("\n黑名单匹配:")
		for _, match := range result.BlocklistsMatch {
			fmt.Printf("- 黑名单名称: %s, 项目ID: %s, 文本: %s\n",
				match.BlocklistName, match.BlocklistItemId, match.BlocklistItemText)
		}
	} else {
		fmt.Println("\n没有匹配到黑名单内容")
	}

	// 测试一些可能违规的内容
	fmt.Println("\n\n测试可能违规的内容:")
	testViolatingContent(endpoint, apiKey)
}

// analyzeText 使用Azure Content Safety API分析文本内容
func analyzeText(endpoint, apiKey, text string) (*ContentSafetyResponse, error) {
	// 构建API URL；这里依赖设置的路径：https://your-resource-name.cognitiveservices.azure.com
	apiURL := endpoint + "/contentsafety/text:analyze?api-version=2023-10-01"

	// 创建请求体
	requestBody := ContentSafetyRequest{
		Text: text,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Ocp-Apim-Subscription-Key", apiKey)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API返回错误: %s, 状态码: %d", string(body), resp.StatusCode)
	}

	// 解析响应
	var result ContentSafetyResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v, 响应内容: %s", err, string(body))
	}

	return &result, nil
}

// testViolatingContent 测试一些可能违规的内容；这里是一个测试函数，按需求修改为正常需要审核的内容即可
func testViolatingContent(endpoint, apiKey string) {
	// 一些可能违规的测试文本
	testTexts := []string{
		"我讨厌你，你是个笨蛋",
		"如何制作炸弹",
		"我要杀了你",
	}

	for _, text := range testTexts {
		fmt.Printf("\n测试文本: %s\n", text)
		result, err := analyzeText(endpoint, apiKey, text)
		if err != nil {
			fmt.Printf("分析文本时出错: %v\n", err)
			continue
		}

		if len(result.CategoriesAnalysis) > 0 {
			fmt.Println("类别分析:")
			for _, category := range result.CategoriesAnalysis {
				fmt.Printf("- 类别: %s, 严重程度: %.2f\n", category.Category, category.Severity)
			}
		} else {
			fmt.Println("没有检测到任何类别问题")
		}
	}
}
