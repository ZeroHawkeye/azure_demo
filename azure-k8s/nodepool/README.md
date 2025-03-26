# Azure AKS 节点池管理工具

这是一个用于管理 Azure Kubernetes Service (AKS) 节点池的命令行工具。使用该工具可以方便地创建、删除和列出 AKS 集群中的节点池。

代码中有详细注释，可在github自行查看: https://github.com/ZeroHawkeye/azure_demo

## 环境要求

- Go 1.16 或更高版本
- Azure 订阅
- Azure AKS 集群

## 配置

在使用此工具之前，需要设置以下环境变量。你可以创建一个 `.env` 文件并填入以下信息：

```env
# 这里的凭证信息是必须的
AZURE_SUBSCRIPTION_ID=
AZURE_TENANT_ID=
AZURE_CLIENT_ID=
AZURE_CLIENT_SECRET=
AZURE_RESOURCE_GROUP=default
```

## 安装依赖

```bash
go mod tidy
```

## 使用方法

该工具支持以下命令：

### 列出节点池

```bash
go run main.go -action list
```

### 创建节点池

```bash
go run main.go -action create -name <节点池名称>
```

### 删除节点池

```bash
go run main.go -action delete -name <节点池名称>
```

## 命令行参数

- `-action`: 操作类型，可选值：create/delete/list（默认为 list）
- `-name`: 节点池名称（默认为 userpool1）

## 默认配置

创建节点池时使用以下默认配置：
- 节点数量：2
- 虚拟机规格：Standard_DS2_v2
- 节点池模式：用户节点池
- Kubernetes 版本：1.30.10
- 节点池类型：VirtualMachineScaleSets

## 注意事项

1. 确保已正确配置 Azure 认证信息
2. 确保有足够的权限管理 AKS 集群
3. 删除节点池前请确保没有工作负载在该节点池上运行 