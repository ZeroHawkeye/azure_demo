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

// 全局变量定义
var (
	subscriptionID string // Azure 订阅 ID
	resourceGroup  string // Azure 资源组名称
)

// init 函数在程序启动时执行，用于初始化必要的配置
func init() {
	// 加载 .env 文件中的环境变量
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	// 从环境变量中获取 Azure 订阅 ID 和资源组名称
	subscriptionID = os.Getenv("AZURE_SUBSCRIPTION_ID")
	resourceGroup = os.Getenv("AZURE_RESOURCE_GROUP")
}

func main() {
	// 创建 Azure 认证凭据
	// 使用客户端密钥认证方式，需要提供租户 ID、客户端 ID 和客户端密钥
	cred, err := azidentity.NewClientSecretCredential(
		os.Getenv("AZURE_TENANT_ID"),     // Azure AD 租户 ID
		os.Getenv("AZURE_CLIENT_ID"),     // 应用程序（服务主体）的客户端 ID
		os.Getenv("AZURE_CLIENT_SECRET"), // 应用程序的客户端密钥
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	// 创建 AKS 管理集群的客户端
	// 使用订阅 ID 和认证凭据初始化客户端
	client, err := armcontainerservice.NewManagedClustersClient(subscriptionID, cred, nil)
	if err != nil {
		log.Fatal(err)
	}

	// 定义要删除的 AKS 集群名称
	clusterName := "learn-ask-tmp1"
	fmt.Printf("开始删除集群 %s...\n", clusterName)

	// 开始删除集群操作
	// BeginDelete 返回一个轮询器，用于跟踪删除操作的进度
	deletePoller, err := client.BeginDelete(context.Background(), resourceGroup, clusterName, nil)
	if err != nil {
		log.Fatal(err)
	}

	// 等待删除操作完成
	// PollUntilDone 会持续轮询直到操作完成或发生错误
	// 频率是在没有重试标头的情况下在轮询间隔之间等待的时间。允许的最小值是一秒钟。
	// 通过零以接受默认值（30s）。
	// 注意这里的超时限制是函数内部的轮训，不是这个方法返回的时间，例如刚刚等待了超过 30s。
	_, err = deletePoller.PollUntilDone(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("集群删除成功！\n")
}
