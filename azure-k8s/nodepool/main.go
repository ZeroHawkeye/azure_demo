package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v4"
	"github.com/joho/godotenv"
)

// 定义全局变量，用于存储Azure资源的关键信息
var (
	subscriptionID string // Azure订阅ID
	resourceGroup  string // 资源组名称
	clusterName    string // AKS集群名称
)

// init函数在程序启动时自动执行，用于初始化配置
func init() {
	// 加载.env文件中的环境变量
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	// 从环境变量中获取Azure认证所需的信息
	subscriptionID = os.Getenv("AZURE_SUBSCRIPTION_ID")
	resourceGroup = os.Getenv("AZURE_RESOURCE_GROUP")
	clusterName = "learn-ask-tmp1" // 使用与创建集群时相同的集群名称;上一个视频中使用api创建的集群名称
}

// createNodePool 创建新的节点池
// 参数：
//   - client: Azure AKS节点池客户端
//   - nodepoolName: 要创建的节点池名称
//
// 返回：
//   - error: 如果创建失败则返回错误信息
func createNodePool(client *armcontainerservice.AgentPoolsClient, nodepoolName string) error {
	// 设置节点池的详细配置参数
	parameters := armcontainerservice.AgentPool{
		Properties: &armcontainerservice.ManagedClusterAgentPoolProfileProperties{
			Count:               toPtr[int32](2),                                                 // 节点数量：设置为2个节点
			VMSize:              toPtr("Standard_DS2_v2"),                                        // 虚拟机规格：2核8G内存
			Mode:                toPtr(armcontainerservice.AgentPoolModeUser),                    // 节点池模式：用户节点池（用于运行用户工作负载）
			OrchestratorVersion: toPtr("1.30.10"),                                                // Kubernetes版本：1.30.10
			Type:                toPtr(armcontainerservice.AgentPoolTypeVirtualMachineScaleSets), // 节点池类型：使用虚拟机规模集
		},
	}

	// 开始创建节点池的异步操作
	fmt.Printf("开始创建节点池 %s...\n", nodepoolName)
	poller, err := client.BeginCreateOrUpdate(context.Background(), resourceGroup, clusterName, nodepoolName, parameters, nil)
	if err != nil {
		return err
	}

	// 等待节点池创建完成
	_, err = poller.PollUntilDone(context.Background(), nil)
	if err != nil {
		return err
	}
	fmt.Printf("节点池创建成功！\n")
	return nil
}

// deleteNodePool 删除指定的节点池
// 参数：
//   - client: Azure AKS节点池客户端
//   - nodepoolName: 要删除的节点池名称
//
// 返回：
//   - error: 如果删除失败则返回错误信息
func deleteNodePool(client *armcontainerservice.AgentPoolsClient, nodepoolName string) error {
	// 开始删除节点池的异步操作
	fmt.Printf("开始删除节点池 %s...\n", nodepoolName)
	poller, err := client.BeginDelete(context.Background(), resourceGroup, clusterName, nodepoolName, nil)
	if err != nil {
		return err
	}

	// 等待节点池删除完成
	_, err = poller.PollUntilDone(context.Background(), nil)
	if err != nil {
		return err
	}
	fmt.Printf("节点池删除成功！\n")
	return nil
}

// listNodePools 列出集群中的所有节点池
// 参数：
//   - client: Azure AKS节点池客户端
//
// 返回：
//   - error: 如果列出失败则返回错误信息
func listNodePools(client *armcontainerservice.AgentPoolsClient) error {
	fmt.Println("\n列出所有节点池：")
	// 创建分页器以获取所有节点池
	pager := client.NewListPager(resourceGroup, clusterName, nil)
	for pager.More() {
		// 获取下一页的结果
		nextResult, err := pager.NextPage(context.Background())
		if err != nil {
			return err
		}
		// 遍历并打印每个节点池的详细信息
		for _, v := range nextResult.Value {
			fmt.Printf("节点池名称: %s\n", *v.Name)
			fmt.Printf("节点数量: %d\n", *v.Properties.Count)
			fmt.Printf("虚拟机规格: %s\n", *v.Properties.VMSize)
			fmt.Printf("模式: %s\n", *v.Properties.Mode)
			fmt.Println("---")
		}
	}
	return nil
}

// main 函数是程序的入口点
func main() {
	// 定义命令行参数
	action := flag.String("action", "list", "操作类型：create/delete/list") // 默认操作为list
	nodepoolName := flag.String("name", "userpool1", "节点池名称")          // 默认节点池名称为userpool1
	flag.Parse()

	// 创建Azure认证凭据，使用客户端密钥认证方式
	cred, err := azidentity.NewClientSecretCredential(
		os.Getenv("AZURE_TENANT_ID"),     // Azure租户ID
		os.Getenv("AZURE_CLIENT_ID"),     // Azure客户端ID
		os.Getenv("AZURE_CLIENT_SECRET"), // Azure客户端密钥
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	// 创建AKS节点池客户端
	client, err := armcontainerservice.NewAgentPoolsClient(subscriptionID, cred, nil)
	if err != nil {
		log.Fatal(err)
	}

	// 根据命令行参数执行相应的操作
	switch *action {
	case "create":
		if err := createNodePool(client, *nodepoolName); err != nil {
			log.Fatal(err)
		}
	case "delete":
		if err := deleteNodePool(client, *nodepoolName); err != nil {
			log.Fatal(err)
		}
	case "list":
		if err := listNodePools(client); err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Println("无效的操作类型。请使用 -action 参数指定操作类型：create/delete/list")
		flag.PrintDefaults()
	}
}

// toPtr 是一个泛型辅助函数，用于将任意类型的值转换为指针
// 参数：
//   - v: 任意类型的值
//
// 返回：
//   - *T: 指向该值的指针
func toPtr[T any](v T) *T {
	return &v
}
