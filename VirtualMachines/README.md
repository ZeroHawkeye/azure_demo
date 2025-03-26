# Azure 创建虚拟网络

## DEMO简介
该项目提供了一个简单的工具，用于在Azure云平台上创建和管理虚拟网络资源。通过Azure SDK for Go，实现了虚拟网络的自动化创建过程。

## 功能特点
- Azure身份认证与连接
- 虚拟网络（VNet）创建

## 技术栈
- Go语言
- Azure SDK for Go
- 环境变量配置（.env）

## 快速开始

### 前提条件
- 安装Go环境（1.16+）
- Azure账户及订阅
- 配置Azure服务主体（Service Principal）

### 环境配置
1. 复制`.env.example`文件并重命名为`.env`
2. 在`.env`文件中填入您的Azure凭证:
```
AZURE_CLIENT_ID=你的客户端ID
AZURE_TENANT_ID=你的租户ID
AZURE_CLIENT_SECRET=你的客户端密钥
AZURE_SUBSCRIPTION_ID=你的订阅ID
```

### 安装依赖
```
go mod tidy
```

### 运行
```
go run VirtualMachines.go
```

## 项目结构
- `VirtualMachines.go` - 主程序文件
- `.env` - 环境变量配置文件
- `azure_flow.svg` - Azure资源流程图

## 流程图
项目包含一个SVG格式的Azure资源流程图（azure_flow.svg），展示了虚拟网络创建的工作流程。 