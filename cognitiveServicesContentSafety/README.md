# Azure 内容安全 API 示例

这是一个使用 Azure 认知服务内容安全 API 的 Go 语言示例程序，用于检测文本内容中可能存在的有害内容。

## 功能

- 分析文本内容的安全性
- 检测多种类别的有害内容
- 支持黑名单匹配
- 包含违规内容测试示例

## 前提条件

- Go 编程语言环境
- Azure 认知服务账户（内容安全服务）
- 有效的 API 密钥和终端节点

1. 设置环境变量：

在 Windows + Powershell 环境中：

```
$env:AZURE_CONTENT_SAFETY_ENDPOINT = 'https://your-resource-name.cognitiveservices.azure.com/'
$env:AZURE_CONTENT_SAFETY_KEY = 'your-api-key'
```

## 使用方法

运行主程序：

```
go run cognitiveServicesContentSafety.go
```

## 示例输出

程序将分析默认文本和一些可能违规的测试文本，输出类似以下内容：

```
内容安全分析结果:
文本: 这是一个测试文本，用于检测内容安全。

没有检测到任何类别问题

没有匹配到黑名单内容


测试可能违规的内容:

测试文本: 我讨厌你，你是个笨蛋
类别分析:
- 类别: Hate, 严重程度: 0.50

测试文本: 如何制作炸弹
类别分析:
- 类别: Violence, 严重程度: 0.70

测试文本: 我要杀了你
类别分析:
- 类别: Violence, 严重程度: 0.90
- 类别: SelfHarm, 严重程度: 0.60
```

## API 响应说明

API 返回的主要信息包括：

- `categoriesAnalysis`: 内容分类分析结果，包含类别和严重程度
- `blocklistsMatch`: 黑名单匹配结果，包含匹配到的黑名单名称、项目ID和文本

## 注意事项

- 使用前请确保已正确设置环境变量
- API 调用可能会产生费用，请查阅 Azure 定价页面
- 本示例仅用于演示目的，生产环境中请添加适当的错误处理和安全措施 