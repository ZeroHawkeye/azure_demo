package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	"github.com/joho/godotenv"
)

// main 函数是程序的入口点，负责初始化环境变量、建立Azure连接并创建虚拟网络
func main() {
	// 加载环境变量配置文件
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 从环境变量中获取Azure身份认证所需的凭证信息
	clientID := os.Getenv("AZURE_CLIENT_ID")
	tenantID := os.Getenv("AZURE_TENANT_ID")
	clientSecret := os.Getenv("AZURE_CLIENT_SECRET")
	subscriptionID := os.Getenv("AZURE_SUBSCRIPTION_ID")

	// 打印身份验证信息以便调试
	fmt.Printf("Client ID: %s\n", clientID)
	fmt.Printf("Tenant ID: %s\n", tenantID)
	fmt.Printf("Client Secret: %s\n", clientSecret)

	// 调用connectionAzure函数建立与Azure的连接
	conn, err := connectionAzure(clientID, tenantID, clientSecret)
	if err != nil {
		log.Fatal("Error connecting to Azure")
	}

	// 创建网络客户端工厂实例，用于管理虚拟网络资源
	networkClientFactory, err := armnetwork.NewClientFactory(subscriptionID, conn, nil)
	if err != nil {
		log.Fatal(err)
	}

	// 获取虚拟网络客户端
	virtualNetworksClient := networkClientFactory.NewVirtualNetworksClient()

	// 调用createVirtualNetwork函数创建虚拟网络
	// 参数: 上下文、虚拟网络客户端、资源组名称、虚拟机名称
	vnetName, vnetRes, err := createVirtualNetwork(context.Background(), virtualNetworksClient, "learn", "learn-vm-net")
	if err != nil {
		log.Fatal(err)
	}

	// 打印创建的虚拟网络信息
	fmt.Printf("Virtual Network Name: %s\n", vnetName)
	fmt.Printf("Virtual Network Resource: %v\n", vnetRes)

}

// connectionAzure 函数使用提供的凭证建立与Azure的连接
// 参数:
//   - clientID: Azure客户端ID
//   - tenantID: Azure租户ID
//   - clientSecret: Azure客户端密钥
//
// 返回:
//   - azcore.TokenCredential: 用于Azure API认证的令牌凭证
//   - error: 可能发生的错误
func connectionAzure(clientID, tenantID, clientSecret string) (azcore.TokenCredential, error) {
	// 使用默认凭证，会从环境变量中获取需要的值；如果你的环境变量中没有设置请不要使用这个方法
	// cred, err := azidentity.NewDefaultAzureCredential(nil)

	// 使用客户端凭证，需要从环境变量中获取clientID, tenantID, clientSecret
	cred, err := azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, nil)
	if err != nil {
		return nil, err
	}
	return cred, nil
}

// createVirtualNetwork 函数创建Azure虚拟网络
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//   - virtualNetworksClient: 虚拟网络客户端
//   - resourceGroupName: 资源组名称
//   - vmName: 虚拟机名称，用于构建虚拟网络名称
//
// 返回:
//   - string: 创建的虚拟网络名称
//   - *armnetwork.VirtualNetwork: 创建的虚拟网络资源对象
//   - error: 可能发生的错误
func createVirtualNetwork(ctx context.Context, virtualNetworksClient *armnetwork.VirtualNetworksClient, resourceGroupName string, vmName string) (string, *armnetwork.VirtualNetwork, error) {
	// 定义虚拟网络参数
	parameters := armnetwork.VirtualNetwork{
		// 设置虚拟网络的位置为东亚区域；你可以根据需要选择其他区域
		Location: to.Ptr("eastasia"),
		Properties: &armnetwork.VirtualNetworkPropertiesFormat{
			// 设置地址空间
			AddressSpace: &armnetwork.AddressSpace{
				// 设置地址前缀，例如10.1.0.0/16
				AddressPrefixes: []*string{
					to.Ptr("10.1.0.0/16"),
				},
			},
		},
	}

	// 异步调用API创建或更新虚拟网络
	pollerResponse, err := virtualNetworksClient.BeginCreateOrUpdate(ctx, resourceGroupName, fmt.Sprintf("%s-vnet", vmName), parameters, nil)
	if err != nil {
		return "", nil, err
	}

	// 轮询等待操作完成
	resp, err := pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return "", nil, err
	}

	// 返回创建的虚拟网络名称和资源对象
	return fmt.Sprintf("%s-vnet", vmName), &resp.VirtualNetwork, nil
}
