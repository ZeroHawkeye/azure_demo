# Azure 翻译服务调用示例

这是一个简单的Go语言示例，演示如何调用Azure翻译服务API进行文本翻译。

## 前提条件

1. 拥有Azure账号
2. 在Azure门户中创建翻译服务资源
3. 获取翻译服务的订阅密钥和区域信息

## 环境变量设置

在运行程序前，需要设置以下环境变量：

```
# Windows CMD;这里按需使用自己使用的命令行进行设置。
set AZURE_TRANSLATOR_KEY=你的订阅密钥
set AZURE_TRANSLATOR_REGION=你的区域(如eastasia)

# Windows PowerShell
$env:AZURE_TRANSLATOR_KEY="你的订阅密钥"
$env:AZURE_TRANSLATOR_REGION="你的区域"

# Nushell
$env.AZURE_TRANSLATOR_KEY = "你的订阅密钥"
$env.AZURE_TRANSLATOR_REGION = "你的区域"
```

## 运行程序

```
go run translator.go
```

## 代码说明

代码主要实现了以下功能：

1. 从环境变量获取Azure翻译服务的密钥和区域
2. 构建HTTP请求，调用Azure翻译服务API
3. 解析API返回的JSON响应
4. 输出翻译结果

## 自定义翻译

如需修改翻译内容或目标语言，可以编辑代码中的以下变量：

- `textToTranslate`: 要翻译的文本
- `targetLanguage`: 目标语言代码（如"zh-CN"表示简体中文，"en"表示英文）

## 常见语言代码

- 简体中文: zh-CN
- 繁体中文: zh-TW
- 英语: en
- 日语: ja
- 韩语: ko
- 法语: fr
- 德语: de
- 西班牙语: es

更多语言代码请参考[Azure文档](https://docs.microsoft.com/zh-cn/azure/cognitive-services/translator/language-support) ;网络不好就不看了。

## 接口调用地址

> https://api.cognitive.microsofttranslator.com