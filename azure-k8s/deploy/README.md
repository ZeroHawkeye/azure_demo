# Azure AKS 应用部署 Demo

这个demo展示了如何使用Go语言通过Azure SDK和Kubernetes client-go来管理Azure Kubernetes Service (AKS)中的应用部署。

## 功能特性

- 创建Deployment和Service
- 删除Deployment和Service
- 列出所有Deployment
- 支持自定义应用名称和副本数量
- 使用LoadBalancer类型的Service暴露应用

## 前置条件

1. 已安装Go 1.21或更高版本
2. 已创建Azure AKS集群
3. 已配置Azure服务主体（Service Principal）并获取以下信息：
   - 订阅ID (AZURE_SUBSCRIPTION_ID)
   - 租户ID (AZURE_TENANT_ID)
   - 客户端ID (AZURE_CLIENT_ID)
   - 客户端密钥 (AZURE_CLIENT_SECRET)
   - 资源组名称 (AZURE_RESOURCE_GROUP)

## 环境配置

1. 在项目根目录创建`.env`文件，包含以下内容：
```env
AZURE_SUBSCRIPTION_ID=你的订阅ID
AZURE_TENANT_ID=你的租户ID
AZURE_CLIENT_ID=你的客户端ID
AZURE_CLIENT_SECRET=你的客户端密钥
AZURE_RESOURCE_GROUP=你的资源组名称
```

2. 安装依赖：
```bash
go mod tidy
```

## 使用方法

### 创建应用

```bash
go run main.go -action create -name my-app -replicas 2
```

参数说明：
- `-action create`: 创建应用
- `-name`: 应用名称（默认为"demo-app"）
- `-replicas`: 副本数量（默认为2）

### 删除应用

```bash
go run main.go -action delete -name my-app
```

参数说明：
- `-action delete`: 删除应用
- `-name`: 要删除的应用名称

### 列出所有应用

```bash
go run main.go -action list
```

参数说明：
- `-action list`: 列出所有应用

## 默认配置

- 默认使用`default`命名空间
- 默认部署nginx:latest镜像
- 默认暴露80端口
- 默认使用LoadBalancer类型的Service

## 注意事项

1. 确保Azure服务主体有足够的权限访问AKS集群
2. 确保AKS集群已经正常运行
3. 默认使用`default`命名空间，如需修改请更改代码中的`namespace`变量
4. 创建应用后，需要等待LoadBalancer分配外部IP地址才能访问

## 代码结构

- `main.go`: 主程序文件
- `go.mod`: Go模块依赖文件
- `.env`: 环境变量配置文件（需要自行创建）

## 依赖包

主要依赖：
- github.com/Azure/azure-sdk-for-go/sdk/azidentity
- github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v4
- github.com/joho/godotenv
- k8s.io/api
- k8s.io/apimachinery
- k8s.io/client-go 