# Azure 搜索服务管理

此示例展示如何使用 Go 语言调用 Azure API 来管理 Microsoft.Search/searchServices 服务。

## 功能

该示例提供以下功能：

1. 查询所有资源组
2. 查询指定资源组中的 Azure 搜索服务资源
3. 获取搜索服务详细信息
4. 创建新的搜索服务
5. 删除搜索服务
6. 列出搜索索引
7. 创建搜索索引
8. 删除搜索索引

## 使用要求

使用前需要满足以下条件：

1. 安装 Go 1.18 或更高版本
2. 配置 Azure 凭据
3. 拥有 Azure 订阅和足够的权限

## 环境变量配置

程序需要设置以下环境变量：

```
# 必填
AZURE_SUBSCRIPTION_ID=your_subscription_id  # Azure 订阅 ID,用于标识你的 Azure 订阅,可以在 Azure 门户的"订阅"页面中找到

# 如果不设置此变量，程序会先列出所有资源组，然后要求用户输入
AZURE_RESOURCE_GROUP=your_resource_group_name

# Azure 认证信息（以下三个变量需要至少设置一组）
AZURE_TENANT_ID=your_tenant_id    # Azure AD 租户 ID,可以在 Azure 门户 -> Azure Active Directory -> 概述页面中找到
AZURE_CLIENT_ID=your_client_id    # 应用程序(客户端)ID,在 Azure AD 中注册应用时获取
AZURE_CLIENT_SECRET=your_client_secret  # 客户端密码,在 Azure AD 应用程序中生成的密钥

# 或者使用托管身份
# AZURE_CLIENT_ID=your_managed_identity_client_id
```

## 运行方式

在 Windows + Nushell 环境下可以按以下方式运行：

```nushell
# 设置环境变量
$env.AZURE_SUBSCRIPTION_ID = "your-subscription-id"
$env.AZURE_RESOURCE_GROUP = "your-resource-group"
$env.AZURE_TENANT_ID = "your-tenant-id"
$env.AZURE_CLIENT_ID = "your-client-id"
$env.AZURE_CLIENT_SECRET = "your-client-secret"

# 运行程序
cd azure-search-demo
go run .
```

## 交互式菜单

程序启动后会显示交互式菜单，提供以下选项：

```
--- Azure 搜索服务管理 ---
1. 列出资源组
2. 列出搜索服务
3. 获取搜索服务详情
4. 创建搜索服务
5. 删除搜索服务
6. 列出搜索索引
7. 创建搜索索引
8. 删除搜索索引
0. 退出程序
请选择操作 [0-8]:
```

## 输出示例

成功执行后，程序会输出类似以下内容：

```
资源组 your-resource-group 中的搜索服务列表:
名称: my-search-service
  ID: /subscriptions/xxxx/resourceGroups/your-resource-group/providers/Microsoft.Search/searchServices/my-search-service
  位置: eastus
  SKU: standard
  状态: running
  预配状态: succeeded
  副本数: 1
  分区数: 1
------------------------
```

## 注意事项

1. 请确保已安装必要的 Go 依赖包
2. 需要导入 `github.com/Azure/azure-sdk-for-go/sdk/azcore/policy` 和 `github.com/Azure/azure-sdk-for-go/sdk/azidentity` 包
3. 删除搜索服务或索引操作需要确认，以防止误操作
4. 推荐使用 Azure CLI 登录，这样可以自动获取凭据
5. 创建搜索服务和索引可能需要几分钟时间才能完成 

### 结语

这里可以看到全部api调用成功了，需要注意的是应用创建后需要进行应用授权角色。

这两个入口都可以进行角色授权，如果订阅中进行了授权则资源组可以继承授权

后续我会增加注释发布的github demo中