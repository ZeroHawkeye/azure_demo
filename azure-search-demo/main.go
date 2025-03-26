package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/joho/godotenv"
)

// SearchService 表示Azure搜索服务资源
// 包含了搜索服务的基本信息，如ID、名称、类型、位置等
type SearchService struct {
	ID         string            `json:"id"`             // 搜索服务的唯一标识符
	Name       string            `json:"name"`           // 搜索服务的名称
	Type       string            `json:"type"`           // 资源类型，通常为 "Microsoft.Search/searchServices"
	Location   string            `json:"location"`       // 搜索服务的地理位置
	Properties SearchProperties  `json:"properties"`     // 搜索服务的详细属性
	SKU        SearchSKU         `json:"sku"`            // 搜索服务的SKU（库存单位）信息
	Tags       map[string]string `json:"tags,omitempty"` // 搜索服务的标签信息
}

// SearchProperties 表示搜索服务属性
// 包含了搜索服务的运行状态、配置信息等
type SearchProperties struct {
	Status            string `json:"status"`            // 搜索服务的当前状态
	StatusDetails     string `json:"statusDetails"`     // 状态详细信息
	ProvisioningState string `json:"provisioningState"` // 预配状态
	ReplicaCount      int    `json:"replicaCount"`      // 副本数量
	PartitionCount    int    `json:"partitionCount"`    // 分区数量
	HostingMode       string `json:"hostingMode"`       // 托管模式
}

// SearchSKU 表示搜索服务SKU
// 定义了搜索服务的性能和功能级别
type SearchSKU struct {
	Name string `json:"name"` // SKU名称，如 "free", "basic", "standard" 等
}

// ResourceListResult 表示资源列表结果
// 用于存储从Azure API获取的资源列表响应
type ResourceListResult struct {
	Value []struct {
		ID       string `json:"id"`       // 资源的唯一标识符
		Name     string `json:"name"`     // 资源名称
		Type     string `json:"type"`     // 资源类型
		Location string `json:"location"` // 资源位置
	} `json:"value"` // 资源列表
}

// SearchServiceCreateParams 表示创建搜索服务的参数
// 包含了创建新搜索服务所需的所有配置信息
type SearchServiceCreateParams struct {
	Location   string            `json:"location"`       // 服务部署位置
	Tags       map[string]string `json:"tags,omitempty"` // 服务标签
	Properties SearchProperties  `json:"properties"`     // 服务属性配置
	SKU        SearchSKU         `json:"sku"`            // SKU配置
}

// SearchIndex 表示搜索索引
// 定义了搜索索引的结构和配置
type SearchIndex struct {
	Name        string                 `json:"name"`                      // 索引名称
	Fields      []SearchField          `json:"fields"`                    // 索引字段列表
	Suggesters  []SearchSuggester      `json:"suggesters,omitempty"`      // 建议器配置
	Scorings    []SearchScoringProfile `json:"scoringProfiles,omitempty"` // 评分配置文件
	Analyzers   []SearchAnalyzer       `json:"analyzers,omitempty"`       // 分析器配置
	Corsoptions SearchCorsOptions      `json:"corsOptions,omitempty"`     // CORS配置
}

// SearchField 表示搜索字段
// 定义了索引中单个字段的属性
type SearchField struct {
	Name         string `json:"name"`                   // 字段名称
	Type         string `json:"type"`                   // 字段类型
	Key          bool   `json:"key,omitempty"`          // 是否为键字段
	Searchable   bool   `json:"searchable,omitempty"`   // 是否可搜索
	Filterable   bool   `json:"filterable,omitempty"`   // 是否可过滤
	Sortable     bool   `json:"sortable,omitempty"`     // 是否可排序
	Facetable    bool   `json:"facetable,omitempty"`    // 是否可分面
	Retrievable  bool   `json:"retrievable,omitempty"`  // 是否可检索
	AnalyzerName string `json:"analyzerName,omitempty"` // 分析器名称
}

// SearchSuggester 表示搜索建议器
// 用于实现搜索建议功能
type SearchSuggester struct {
	Name         string   `json:"name"`         // 建议器名称
	SourceFields []string `json:"sourceFields"` // 用于生成建议的源字段
}

// SearchScoringProfile 表示搜索评分配置文件
// 定义了搜索结果的自定义评分规则
type SearchScoringProfile struct {
	Name                string            `json:"name"`                          // 配置文件名称
	FunctionAggregation string            `json:"functionAggregation,omitempty"` // 函数聚合方式
	Functions           []ScoringFunction `json:"functions,omitempty"`           // 评分函数列表
}

// ScoringFunction 表示评分函数
// 定义了单个评分函数的配置
type ScoringFunction struct {
	Type       string                 `json:"type"`                 // 函数类型
	FieldName  string                 `json:"fieldName"`            // 目标字段名
	Boost      float64                `json:"boost"`                // 提升系数
	Parameters map[string]interface{} `json:"parameters,omitempty"` // 函数参数
}

// SearchAnalyzer 表示搜索分析器
// 定义了文本分析的处理方式
type SearchAnalyzer struct {
	Name string `json:"name"`        // 分析器名称
	Type string `json:"@odata.type"` // 分析器类型
}

// SearchCorsOptions 表示CORS选项
// 配置跨域资源共享规则
type SearchCorsOptions struct {
	AllowedOrigins  []string `json:"allowedOrigins"`            // 允许的源
	MaxAgeInSeconds int      `json:"maxAgeInSeconds,omitempty"` // 预检请求的缓存时间
}

// 全局认证和配置
var (
	cred           *azidentity.DefaultAzureCredential // Azure认证凭据
	subscriptionID string                             // Azure订阅ID
	resourceGroup  string                             // 资源组名称
	httpClient     *http.Client                       // HTTP客户端
)

// main 函数是程序的入口点
// 负责初始化Azure认证、加载环境变量并启动交互式菜单
func main() {
	// 尝试从多个位置加载.env文件
	// 支持当前目录和上级目录的.env文件
	envFiles := []string{".env", "../.env"}
	envLoaded := false

	for _, file := range envFiles {
		err := godotenv.Load(file)
		if err == nil {
			log.Printf("成功从 %s 加载环境变量", file)
			envLoaded = true
			break
		}
	}

	if !envLoaded {
		log.Println("警告: 无法加载.env文件，将使用系统环境变量")
	}

	// 从环境变量获取Azure订阅ID
	subscriptionID = os.Getenv("AZURE_SUBSCRIPTION_ID")
	if subscriptionID == "" {
		log.Fatal("请设置环境变量：AZURE_SUBSCRIPTION_ID")
	}

	// 初始化Azure认证凭据
	// 使用默认的Azure认证方式，支持多种认证方法（如环境变量、托管身份等）
	var err error
	cred, err = azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("获取凭据失败: %v", err)
	}

	// 初始化HTTP客户端
	// 设置30秒超时时间，用于所有API请求
	httpClient = &http.Client{Timeout: 30 * time.Second}

	// 获取资源组名称
	// 如果环境变量中未指定，则列出所有资源组供用户选择
	resourceGroup = os.Getenv("AZURE_RESOURCE_GROUP")
	if resourceGroup == "" {
		fmt.Println("未指定资源组名称，将列出所有资源组...")
		getResourceGroups()
		fmt.Print("请输入要使用的资源组名称: ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			resourceGroup = scanner.Text()
		}
		if resourceGroup == "" {
			fmt.Println("未指定资源组，退出程序")
			return
		}
	}

	// 显示交互式菜单
	showMenu()
}

// showMenu 函数显示交互式菜单
// 提供了一系列管理Azure搜索服务的操作选项
func showMenu() {
	for {
		fmt.Println("\n--- Azure 搜索服务管理 ---")
		fmt.Println("1. 列出资源组")
		fmt.Println("2. 列出搜索服务")
		fmt.Println("3. 获取搜索服务详情")
		fmt.Println("4. 创建搜索服务")
		fmt.Println("5. 删除搜索服务")
		fmt.Println("6. 列出搜索索引")
		fmt.Println("7. 创建搜索索引")
		fmt.Println("8. 删除搜索索引")
		fmt.Println("0. 退出程序")
		fmt.Print("请选择操作 [0-8]: ")

		var choice string
		fmt.Scanln(&choice)

		// 根据用户选择执行相应的操作
		switch choice {
		case "0":
			fmt.Println("退出程序")
			return
		case "1":
			getResourceGroups()
		case "2":
			getSearchServices()
		case "3":
			getSearchServiceDetails()
		case "4":
			createSearchService()
		case "5":
			deleteSearchService()
		case "6":
			listSearchIndexes()
		case "7":
			createSearchIndex()
		case "8":
			deleteSearchIndex()
		default:
			fmt.Println("无效的选择，请重试")
		}
	}
}

// getToken 函数获取Azure管理API的访问令牌
// 使用DefaultAzureCredential获取认证令牌，用于后续的API请求
func getToken(ctx context.Context) (string, error) {
	token, err := cred.GetToken(ctx, policy.TokenRequestOptions{
		Scopes: []string{"https://management.azure.com/.default"}, // Azure管理API的作用域
	})
	if err != nil {
		return "", fmt.Errorf("获取令牌失败: %v", err)
	}
	return token.Token, nil
}

// sendRequest 函数是通用的HTTP请求发送函数
// 处理认证、请求发送和响应处理等通用逻辑
func sendRequest(ctx context.Context, method, url string, body io.Reader) ([]byte, error) {
	// 获取访问令牌
	token, err := getToken(ctx)
	if err != nil {
		return nil, err
	}

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	// 设置认证头
	req.Header.Set("Authorization", "Bearer "+token)

	// 如果有请求体，设置Content-Type头
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// 发送请求
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查响应状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(responseBody))
	}

	return responseBody, nil
}

// getResourceGroups 函数列出当前订阅下的所有资源组
// 通过Azure管理API获取资源组列表并显示
func getResourceGroups() {
	ctx := context.Background()
	// 构建获取资源组列表的API URL
	url := fmt.Sprintf("https://management.azure.com/subscriptions/%s/resourcegroups?api-version=2021-04-01", subscriptionID)

	// 发送GET请求获取资源组列表
	body, err := sendRequest(ctx, "GET", url, nil)
	if err != nil {
		log.Fatalf("获取资源组失败: %v", err)
	}

	// 解析响应
	var result ResourceListResult
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatalf("解析响应失败: %v", err)
	}

	// 显示资源组列表
	fmt.Println("\n资源组列表:")
	for _, group := range result.Value {
		fmt.Printf("名称: %s, 位置: %s\n", group.Name, group.Location)
	}
}

// getSearchServices 函数列出指定资源组中的所有搜索服务
// 显示每个搜索服务的详细信息，包括状态、配置等
func getSearchServices() {
	ctx := context.Background()
	// 构建获取搜索服务列表的API URL
	url := fmt.Sprintf("https://management.azure.com/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Search/searchServices?api-version=2020-08-01",
		subscriptionID, resourceGroup)

	// 发送GET请求获取搜索服务列表
	body, err := sendRequest(ctx, "GET", url, nil)
	if err != nil {
		log.Fatalf("获取搜索服务失败: %v", err)
	}

	// 解析响应
	var result struct {
		Value []SearchService `json:"value"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatalf("解析响应失败: %v", err)
	}

	// 显示搜索服务列表
	fmt.Printf("\n资源组 %s 中的搜索服务列表:\n", resourceGroup)
	if len(result.Value) == 0 {
		fmt.Println("未找到搜索服务")
		return
	}

	// 遍历并显示每个搜索服务的详细信息
	for _, service := range result.Value {
		fmt.Printf("名称: %s\n", service.Name)
		fmt.Printf("  ID: %s\n", service.ID)
		fmt.Printf("  位置: %s\n", service.Location)
		fmt.Printf("  SKU: %s\n", service.SKU.Name)
		fmt.Printf("  状态: %s\n", service.Properties.Status)
		fmt.Printf("  预配状态: %s\n", service.Properties.ProvisioningState)
		fmt.Printf("  副本数: %d\n", service.Properties.ReplicaCount)
		fmt.Printf("  分区数: %d\n", service.Properties.PartitionCount)
		fmt.Println("------------------------")
	}
}

// getSearchServiceDetails 函数获取指定搜索服务的详细信息
// 通过服务名称获取并显示该服务的完整配置信息
func getSearchServiceDetails() {
	// 获取用户输入的服务名称
	fmt.Print("请输入搜索服务名称: ")
	var serviceName string
	fmt.Scanln(&serviceName)

	if serviceName == "" {
		fmt.Println("搜索服务名称不能为空")
		return
	}

	ctx := context.Background()
	// 构建获取搜索服务详情的API URL
	url := fmt.Sprintf("https://management.azure.com/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Search/searchServices/%s?api-version=2020-08-01",
		subscriptionID, resourceGroup, serviceName)

	// 发送GET请求获取服务详情
	body, err := sendRequest(ctx, "GET", url, nil)
	if err != nil {
		fmt.Printf("获取搜索服务详情失败: %v\n", err)
		return
	}

	// 解析响应
	var service SearchService
	if err := json.Unmarshal(body, &service); err != nil {
		fmt.Printf("解析响应失败: %v\n", err)
		return
	}

	// 显示服务详情
	fmt.Printf("\n搜索服务详情:\n")
	fmt.Printf("名称: %s\n", service.Name)
	fmt.Printf("ID: %s\n", service.ID)
	fmt.Printf("位置: %s\n", service.Location)
	fmt.Printf("SKU: %s\n", service.SKU.Name)
	fmt.Printf("状态: %s\n", service.Properties.Status)
	fmt.Printf("预配状态: %s\n", service.Properties.ProvisioningState)
	fmt.Printf("副本数: %d\n", service.Properties.ReplicaCount)
	fmt.Printf("分区数: %d\n", service.Properties.PartitionCount)

	// 显示标签信息（如果有）
	if len(service.Tags) > 0 {
		fmt.Println("标签:")
		for k, v := range service.Tags {
			fmt.Printf("  %s: %s\n", k, v)
		}
	}
}

// createSearchService 函数创建新的搜索服务
// 通过用户输入配置创建新的Azure搜索服务
func createSearchService() {
	// 获取用户输入的服务名称
	fmt.Print("请输入搜索服务名称: ")
	var serviceName string
	fmt.Scanln(&serviceName)

	if serviceName == "" {
		fmt.Println("搜索服务名称不能为空")
		return
	}

	// 获取用户输入的位置信息
	fmt.Print("请输入位置 (如: eastus, westus): ")
	var location string
	fmt.Scanln(&location)

	if location == "" {
		fmt.Println("位置不能为空")
		return
	}

	// 获取用户输入的SKU信息
	fmt.Print("请输入SKU (free, basic, standard, standard2, standard3): ")
	var sku string
	fmt.Scanln(&sku)

	if sku == "" {
		sku = "standard" // 默认值
	}

	// 获取用户输入的副本数
	fmt.Print("请输入副本数 (1-12): ")
	var replicaCountStr string
	fmt.Scanln(&replicaCountStr)
	replicaCount, err := strconv.Atoi(replicaCountStr)
	if err != nil || replicaCount < 1 {
		replicaCount = 1 // 默认值
	}

	// 获取用户输入的分区数
	fmt.Print("请输入分区数 (1-12): ")
	var partitionCountStr string
	fmt.Scanln(&partitionCountStr)
	partitionCount, err := strconv.Atoi(partitionCountStr)
	if err != nil || partitionCount < 1 {
		partitionCount = 1 // 默认值
	}

	// 准备创建参数
	createParams := SearchServiceCreateParams{
		Location: location,
		SKU: SearchSKU{
			Name: sku,
		},
		Properties: SearchProperties{
			ReplicaCount:   replicaCount,
			PartitionCount: partitionCount,
			HostingMode:    "default",
		},
		Tags: map[string]string{
			"created-by": "azure-search-demo",
			"created-at": time.Now().Format(time.RFC3339),
		},
	}

	// 序列化请求体
	requestBody, err := json.Marshal(createParams)
	if err != nil {
		fmt.Printf("序列化请求体失败: %v\n", err)
		return
	}

	ctx := context.Background()
	// 构建创建搜索服务的API URL
	url := fmt.Sprintf("https://management.azure.com/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Search/searchServices/%s?api-version=2020-08-01",
		subscriptionID, resourceGroup, serviceName)

	// 发送PUT请求创建搜索服务
	_, err = sendRequest(ctx, "PUT", url, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Printf("创建搜索服务失败: %v\n", err)
		return
	}

	fmt.Printf("搜索服务 %s 创建请求已提交，请等待几分钟后查询服务状态\n", serviceName)
}

// deleteSearchService 函数删除指定的搜索服务
// 包含确认步骤以防止意外删除
func deleteSearchService() {
	// 获取用户输入的服务名称
	fmt.Print("请输入要删除的搜索服务名称: ")
	var serviceName string
	fmt.Scanln(&serviceName)

	if serviceName == "" {
		fmt.Println("搜索服务名称不能为空")
		return
	}

	// 请求用户确认删除操作
	fmt.Printf("警告: 此操作将删除搜索服务 %s 及其所有数据，此操作不可恢复\n", serviceName)
	fmt.Print("请输入 'confirm' 确认删除: ")
	var confirmation string
	fmt.Scanln(&confirmation)

	if confirmation != "confirm" {
		fmt.Println("操作已取消")
		return
	}

	ctx := context.Background()
	// 构建删除搜索服务的API URL
	url := fmt.Sprintf("https://management.azure.com/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Search/searchServices/%s?api-version=2020-08-01",
		subscriptionID, resourceGroup, serviceName)

	// 发送DELETE请求删除搜索服务
	_, err := sendRequest(ctx, "DELETE", url, nil)
	if err != nil {
		fmt.Printf("删除搜索服务失败: %v\n", err)
		return
	}

	fmt.Printf("搜索服务 %s 删除请求已提交\n", serviceName)
}

// listSearchIndexes 函数列出指定搜索服务中的所有索引
// 显示每个索引的名称和字段信息
func listSearchIndexes() {
	// 获取用户输入的服务名称
	fmt.Print("请输入搜索服务名称: ")
	var serviceName string
	fmt.Scanln(&serviceName)

	if serviceName == "" {
		fmt.Println("搜索服务名称不能为空")
		return
	}

	// 获取搜索服务的API密钥
	adminKey, err := getSearchServiceAdminKey(serviceName)
	if err != nil {
		fmt.Printf("获取API密钥失败: %v\n", err)
		return
	}

	// 构建获取索引列表的API URL
	ctx := context.Background()
	url := fmt.Sprintf("https://%s.search.windows.net/indexes?api-version=2020-06-30", serviceName)

	// 创建并发送GET请求
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		fmt.Printf("创建请求失败: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", adminKey)

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Printf("发送请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取响应失败: %v\n", err)
		return
	}

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("请求失败，状态码: %d, 响应: %s\n", resp.StatusCode, string(body))
		return
	}

	// 解析响应
	var result struct {
		Value []struct {
			Name   string `json:"name"` // 索引名称
			Fields []struct {
				Name string `json:"name"` // 字段名称
				Type string `json:"type"` // 字段类型
			} `json:"fields"` // 字段列表
		} `json:"value"` // 索引列表
	}

	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("解析响应失败: %v\n", err)
		return
	}

	// 显示索引列表
	fmt.Printf("\n搜索服务 %s 中的索引列表:\n", serviceName)
	if len(result.Value) == 0 {
		fmt.Println("未找到索引")
		return
	}

	// 遍历并显示每个索引的信息
	for _, index := range result.Value {
		fmt.Printf("索引名称: %s\n", index.Name)
		fmt.Println("  字段:")
		for _, field := range index.Fields {
			fmt.Printf("    %s (%s)\n", field.Name, field.Type)
		}
		fmt.Println("------------------------")
	}
}

// getSearchServiceAdminKey 函数获取搜索服务的管理密钥
// 用于后续的搜索服务API调用认证
func getSearchServiceAdminKey(serviceName string) (string, error) {
	ctx := context.Background()
	// 构建获取管理密钥的API URL
	url := fmt.Sprintf("https://management.azure.com/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Search/searchServices/%s/listAdminKeys?api-version=2020-08-01",
		subscriptionID, resourceGroup, serviceName)

	// 发送POST请求获取管理密钥
	body, err := sendRequest(ctx, "POST", url, nil)
	if err != nil {
		return "", err
	}

	// 解析响应获取主密钥
	var result struct {
		PrimaryKey   string `json:"primaryKey"`   // 主密钥
		SecondaryKey string `json:"secondaryKey"` // 次密钥
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	return result.PrimaryKey, nil
}

// createSearchIndex 函数创建新的搜索索引
// 创建一个包含常用字段的示例索引
func createSearchIndex() {
	// 获取用户输入的服务名称
	fmt.Print("请输入搜索服务名称: ")
	var serviceName string
	fmt.Scanln(&serviceName)

	if serviceName == "" {
		fmt.Println("搜索服务名称不能为空")
		return
	}

	// 获取用户输入的索引名称
	fmt.Print("请输入索引名称: ")
	var indexName string
	fmt.Scanln(&indexName)

	if indexName == "" {
		fmt.Println("索引名称不能为空")
		return
	}

	// 创建一个示例索引定义
	// 包含常用的字段类型和配置
	indexDef := SearchIndex{
		Name: indexName,
		Fields: []SearchField{
			{
				Name:        "id",
				Type:        "Edm.String",
				Key:         true,
				Searchable:  false,
				Filterable:  false,
				Sortable:    false,
				Facetable:   false,
				Retrievable: true,
			},
			{
				Name:        "title",
				Type:        "Edm.String",
				Searchable:  true,
				Filterable:  true,
				Sortable:    true,
				Facetable:   false,
				Retrievable: true,
			},
			{
				Name:        "content",
				Type:        "Edm.String",
				Searchable:  true,
				Filterable:  false,
				Sortable:    false,
				Facetable:   false,
				Retrievable: true,
			},
			{
				Name:        "category",
				Type:        "Edm.String",
				Searchable:  true,
				Filterable:  true,
				Sortable:    true,
				Facetable:   true,
				Retrievable: true,
			},
			{
				Name:        "rating",
				Type:        "Edm.Int32",
				Searchable:  false,
				Filterable:  true,
				Sortable:    true,
				Facetable:   true,
				Retrievable: true,
			},
			{
				Name:        "lastUpdated",
				Type:        "Edm.DateTimeOffset",
				Searchable:  false,
				Filterable:  true,
				Sortable:    true,
				Facetable:   false,
				Retrievable: true,
			},
		},
		Suggesters: []SearchSuggester{
			{
				Name:         "sg",
				SourceFields: []string{"title", "category"},
			},
		},
		Corsoptions: SearchCorsOptions{
			AllowedOrigins:  []string{"*"}, // 允许所有来源访问
			MaxAgeInSeconds: 300,
		},
	}

	// 获取搜索服务管理密钥
	adminKey, err := getSearchServiceAdminKey(serviceName)
	if err != nil {
		fmt.Printf("获取API密钥失败: %v\n", err)
		return
	}

	// 序列化索引定义
	requestBody, err := json.Marshal(indexDef)
	if err != nil {
		fmt.Printf("序列化请求体失败: %v\n", err)
		return
	}

	// 构建创建索引的API URL
	ctx := context.Background()
	url := fmt.Sprintf("https://%s.search.windows.net/indexes/%s?api-version=2020-06-30", serviceName, indexName)

	// 创建并发送PUT请求
	req, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Printf("创建请求失败: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", adminKey)

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Printf("发送请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取响应失败: %v\n", err)
		return
	}

	// 检查响应状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Printf("创建索引失败，状态码: %d, 响应: %s\n", resp.StatusCode, string(body))
		return
	}

	fmt.Printf("索引 %s 创建成功\n", indexName)
}

// deleteSearchIndex 函数删除指定的搜索索引
// 包含确认步骤以防止意外删除
func deleteSearchIndex() {
	// 获取用户输入的服务名称
	fmt.Print("请输入搜索服务名称: ")
	var serviceName string
	fmt.Scanln(&serviceName)

	if serviceName == "" {
		fmt.Println("搜索服务名称不能为空")
		return
	}

	// 获取用户输入的索引名称
	fmt.Print("请输入要删除的索引名称: ")
	var indexName string
	fmt.Scanln(&indexName)

	if indexName == "" {
		fmt.Println("索引名称不能为空")
		return
	}

	// 请求用户确认删除操作
	fmt.Printf("警告: 此操作将删除索引 %s 及其所有数据，此操作不可恢复\n", indexName)
	fmt.Print("请输入 'confirm' 确认删除: ")
	var confirmation string
	fmt.Scanln(&confirmation)

	if confirmation != "confirm" {
		fmt.Println("操作已取消")
		return
	}

	// 获取搜索服务管理密钥
	adminKey, err := getSearchServiceAdminKey(serviceName)
	if err != nil {
		fmt.Printf("获取API密钥失败: %v\n", err)
		return
	}

	// 构建删除索引的API URL
	ctx := context.Background()
	url := fmt.Sprintf("https://%s.search.windows.net/indexes/%s?api-version=2020-06-30", serviceName, indexName)

	// 创建并发送DELETE请求
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		fmt.Printf("创建请求失败: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", adminKey)

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Printf("发送请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("删除索引失败，状态码: %d, 响应: %s\n", resp.StatusCode, string(body))
		return
	}

	fmt.Printf("索引 %s 删除成功\n", indexName)
}
