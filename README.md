# 项目目录说明

这是一个包含多个Azure相关Go语言示例代码的仓库。

## 目录结构

```
e:/go/tmp/
├── go.mod
├── go.sum
├── README.md
├── ai/                    # AI相关示例
│   ├── .env.example
│   ├── main.go
│   └── README.md
├── azblob/                # Azure Blob存储操作
│   ├── learn.txt
│   ├── Readme.md
│   ├── upload_test.go
│   └── upload.go
├── azure-k8s/             # Azure Kubernetes服务操作
│   ├── create/            # 创建集群
│   ├── delete/            # 删除集群
│   ├── deploy/            # 部署应用
│   └── nodepool/          # 节点池管理
├── azure-search-demo/     # Azure搜索服务示例
├── azure-translator/      # Azure翻译服务
├── AzureNsg/              # Azure网络安全组
├── cognitiveServicesContentSafety/  # 内容安全服务
└── VirtualMachines/       # 虚拟机管理
```

## 使用说明

1. 每个子目录包含独立的Go项目
2. 大多数项目需要配置.env文件（参考.env.example）
3. 运行程序前请先执行：
```bash
go mod tidy
```

## 项目维护

- 所有代码使用Go 1.20+编写
- 遵循标准Go项目结构
- 各子目录有独立README说明具体用法