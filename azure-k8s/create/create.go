package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v4"
	"github.com/joho/godotenv"
)

// 定义全局变量，用于存储Azure订阅ID和资源组名称
var (
	subscriptionID string
	resourceGroup  string
)

// init函数在程序启动时自动执行
// 用于加载环境变量配置文件(.env)并初始化必要的变量
func init() {
	// 加载.env文件，这个文件包含了Azure的认证信息
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	// 从环境变量中获取Azure订阅ID和资源组名称
	subscriptionID = os.Getenv("AZURE_SUBSCRIPTION_ID")
	resourceGroup = os.Getenv("AZURE_RESOURCE_GROUP")
}

func main() {
	// 创建Azure认证凭据
	// 使用客户端ID和密钥进行身份验证
	cred, err := azidentity.NewClientSecretCredential(
		os.Getenv("AZURE_TENANT_ID"),     // Azure租户ID
		os.Getenv("AZURE_CLIENT_ID"),     // Azure客户端ID
		os.Getenv("AZURE_CLIENT_SECRET"), // Azure客户端密钥
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	// 创建Azure Kubernetes Service (AKS) 客户端
	// 这个客户端用于与Azure的Kubernetes服务进行交互
	client, err := armcontainerservice.NewManagedClustersClient(subscriptionID, cred, nil)
	if err != nil {
		log.Fatal(err)
	}

	// 设置要创建的AKS集群的基本信息
	clusterName := "learn-ask-tmp1" // 集群名称
	location := "eastus"            // 集群部署位置

	// 配置AKS集群的参数
	parameters := armcontainerservice.ManagedCluster{
		Location: &location,
		Properties: &armcontainerservice.ManagedClusterProperties{
			DNSPrefix: &clusterName, // DNS前缀，用于访问集群
			// 配置节点池信息
			AgentPoolProfiles: []*armcontainerservice.ManagedClusterAgentPoolProfile{
				{
					Name:   toPtr("nodepool1"),                             // 节点池名称
					Count:  toPtr[int32](1),                                // 节点数量
					VMSize: toPtr("Standard_DS2_v2"),                       // 虚拟机规格
					Mode:   toPtr(armcontainerservice.AgentPoolModeSystem), // 节点池模式：系统节点池
				},
			},
			// 配置服务主体信息，用于集群访问Azure资源
			ServicePrincipalProfile: &armcontainerservice.ManagedClusterServicePrincipalProfile{
				ClientID: toPtr(os.Getenv("AZURE_CLIENT_ID")),
				Secret:   toPtr(os.Getenv("AZURE_CLIENT_SECRET")),
			},
		},
	}

	// 开始创建或更新AKS集群
	poller, err := client.BeginCreateOrUpdate(context.Background(), resourceGroup, clusterName, parameters, nil)
	if err != nil {
		log.Fatal(err)
	}

	// 等待集群创建完成
	fmt.Printf("开始创建集群 %s...\n", clusterName)
	_, err = poller.PollUntilDone(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("集群创建成功！\n")

	// 列出所有已创建的AKS集群
	fmt.Println("\n列出所有集群：")
	pager := client.NewListPager(nil)
	for pager.More() {
		nextResult, err := pager.NextPage(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		// 打印每个集群的详细信息
		for _, v := range nextResult.Value {
			fmt.Printf("集群名称: %s\n", *v.Name)
			fmt.Printf("位置: %s\n", *v.Location)
			fmt.Printf("状态: %s\n", *v.Properties.ProvisioningState)
			fmt.Println("---")
		}
	}
}

// 辅助函数：将任意类型的值转换为指针
// 在Go中，很多API需要指针类型的参数，这个函数可以方便地创建指针
func toPtr[T any](v T) *T {
	return &v
}
