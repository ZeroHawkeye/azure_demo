# Azure 云服务演示

## 概述

本项目是一系列基于Go语言的示例集合，展示了与各种Azure云服务的集成。每个模块展示了特定Azure服务的实现，为开发者提供了将Azure服务集成到Go应用程序中的实用代码示例。

## 项目结构

仓库组织为以下模块：

- **AI**：Azure OpenAI集成，用于聊天补全和AI功能
- **AzureNsg**：网络安全组(NSG)创建和管理
- **azblob**：Azure Blob存储操作
- **AzureTranslator**：文本翻译服务
- **ContentSafety**：文本和图像内容审核
- **VirtualMachines**：Azure虚拟机部署和管理

## 前提条件

- Go版本 ≥ 1.18
- 有效的Azure订阅
- 服务特定的API密钥和凭证
- 每个模块的环境变量正确配置

## 安装

1. 克隆仓库：
   ```bash
   git clone https://github.com/yourusername/azure_demo.git
   cd azure_demo
   ```

2. 安装依赖：
   ```bash
   go mod tidy
   ```

3. 配置环境变量：
   - 每个模块包含一个`.env.example`文件
   - 将示例复制到新的`.env`文件并添加您的凭证
   - **重要**：永远不要将包含真实凭证的.env文件提交到版本控制

   示例：
   ```bash
   cp AzureNsg/.env.example AzureNsg/.env
   # 使用您的凭证编辑.env文件
   ```

## 服务描述

### AI模块（Azure OpenAI）

演示与Azure OpenAI服务的集成，用于：
- 聊天补全
- 提示工程
- 上下文管理

使用以下环境变量配置：
```
AZURE_OPENAI_ENDPOINT=https://your-resource-name.openai.azure.com/
AZURE_OPENAI_API_KEY=your-api-key
AZURE_OPENAI_DEPLOYMENT_NAME=your-deployment-name
```

### Azure 翻译器

提供使用Azure翻译器服务的文本翻译功能：
- 多语言翻译
- 语言检测
- 文本分析

配置：
```
AZURE_TRANSLATOR_KEY=your-translator-key
AZURE_TRANSLATOR_ENDPOINT=https://api.cognitive.microsofttranslator.com/
AZURE_TRANSLATOR_REGION=your-azure-region
```

### Azure NSG（网络安全组）

演示创建和配置网络安全组：
- 安全规则管理
- 网络接口配置
- 安全策略实施

配置：
```
AZURE_TENANT_ID=your-tenant-id
AZURE_CLIENT_ID=your-client-id
AZURE_CLIENT_SECRET=your-client-secret-placeholder
AZURE_SUBSCRIPTION_ID=your-subscription-id
```

### Azure Blob 存储

展示Azure Blob存储的操作：
- 文件上传和下载
- 容器管理
- 访问控制

配置：
```
AZURE_STORAGE_ACCOUNT=your-storage-account-name
AZURE_STORAGE_ACCESS_KEY=your-access-key-placeholder
```

### Azure 内容安全

提供内容审核的示例：
- 文本内容分析
- 图像内容审核
- 策略实施

### Azure 虚拟机

演示虚拟机管理操作：
- 虚拟机创建和配置
- 资源组管理
- 虚拟机状态管理

配置：
```
AZURE_TENANT_ID=your-tenant-id
AZURE_CLIENT_ID=your-client-id
AZURE_CLIENT_SECRET=your-client-secret-placeholder
AZURE_SUBSCRIPTION_ID=your-subscription-id
```

## 安全注意事项

- **永远不要**将实际的密钥、密钥或凭证提交到版本控制
- 始终使用环境变量或安全保管库存储敏感信息
- 每个模块的示例文件显示带有占位符值的必需变量
- 实际凭证应存储在`.env`文件中并添加到`.gitignore`

## 文档

每个模块在其自己的README文件中包含更详细的文档。有关使用和配置的更多信息，请参阅特定模块。

