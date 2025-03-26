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

	// 创建网络客户端工厂实例，用于管理虚拟网络资源；这里是初始化的凭证管理client，而不是交换机，vpc等实例
	networkClientFactory, err := armnetwork.NewClientFactory(subscriptionID, conn, nil)
	if err != nil {
		log.Fatal(err)
	}
	securityGroupsClient := networkClientFactory.NewSecurityGroupsClient()

	// 创建网络安全组
	nsg, err := createNetworkSecurityGroup(context.Background(), securityGroupsClient, "learn", "learn-vm-nsg")
	if err != nil {
		fmt.Println("创建安全组失败", err)
		return
	}
	fmt.Printf("安全组: %s\n", *nsg.Name)
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

// createNetworkSecurityGroup 创建Azure网络安全组(NSG)
// 参数:
//   - ctx: 上下文对象,用于控制请求的生命周期
//   - securityGroupsClient: Azure安全组客户端实例
//   - resourceGroupName: 资源组名称
//   - vmName: 虚拟机名称,用于生成NSG名称
//
// 返回:
//   - *armnetwork.SecurityGroup: 创建的安全组对象指针
//   - error: 可能发生的错误
func createNetworkSecurityGroup(ctx context.Context, securityGroupsClient *armnetwork.SecurityGroupsClient, resourceGroupName string, vmName string) (*armnetwork.SecurityGroup, error) {

	// 定义安全组配置参数
	parameters := armnetwork.SecurityGroup{
		// 设置安全组所在的Azure区域
		Location: to.Ptr("eastasia"),
		// 配置安全组属性
		Properties: &armnetwork.SecurityGroupPropertiesFormat{
			// 定义安全规则列表；一个安全组可以包含多个规则
			SecurityRules: []*armnetwork.SecurityRule{
				// 入站规则配置
				{
					Name: to.Ptr("sample_inbound_22"), // 规则名称
					Properties: &armnetwork.SecurityRulePropertiesFormat{
						SourceAddressPrefix:      to.Ptr("0.0.0.0/0"),                             // 源IP地址前缀,0.0.0.0/0表示允许任何IP
						SourcePortRange:          to.Ptr("*"),                                     // 源端口范围,*表示所有端口
						DestinationAddressPrefix: to.Ptr("0.0.0.0/0"),                             // 目标IP地址前缀
						DestinationPortRange:     to.Ptr("*"),                                     // 目标端口范围
						Protocol:                 to.Ptr(armnetwork.SecurityRuleProtocolAsterisk), // 协议类型,*表示所有协议
						Access:                   to.Ptr(armnetwork.SecurityRuleAccessAllow),      // 允许访问
						Priority:                 to.Ptr[int32](100),                              // 规则优先级,数字越小优先级越高
						Description:              to.Ptr("允许任意端口和ip访问"),                           // 规则描述
						Direction:                to.Ptr(armnetwork.SecurityRuleDirectionInbound), // 入站方向
					},
				},
				// 出站规则配置
				{
					Name: to.Ptr("sample_outbound_22"), // 规则名称
					Properties: &armnetwork.SecurityRulePropertiesFormat{
						SourceAddressPrefix:      to.Ptr("0.0.0.0/0"),                              // 源IP地址前缀
						SourcePortRange:          to.Ptr("*"),                                      // 源端口范围
						DestinationAddressPrefix: to.Ptr("0.0.0.0/0"),                              // 目标IP地址前缀
						DestinationPortRange:     to.Ptr("*"),                                      // 目标端口范围
						Protocol:                 to.Ptr(armnetwork.SecurityRuleProtocolAsterisk),  // 协议类型
						Access:                   to.Ptr(armnetwork.SecurityRuleAccessAllow),       // 允许访问
						Priority:                 to.Ptr[int32](100),                               // 规则优先级
						Description:              to.Ptr("允许任意端口和ip出站"),                            // 规则描述
						Direction:                to.Ptr(armnetwork.SecurityRuleDirectionOutbound), // 出站方向
					},
				},
			},
		},
	}

	// 生成安全组名称,格式为"{虚拟机名称}-nsg"
	nsgName := fmt.Sprintf("%s-nsg", vmName)

	// 异步创建或更新安全组
	pollerResponse, err := securityGroupsClient.BeginCreateOrUpdate(ctx, resourceGroupName, nsgName, parameters, nil)
	if err != nil {
		return nil, err
	}

	// 等待异步操作完成
	resp, err := pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}

	// 返回创建的安全组对象
	return &resp.SecurityGroup, nil
}
