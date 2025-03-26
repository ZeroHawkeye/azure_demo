- 使用azure时，存在一个比阿里云等国产服务商不存在的一个概念，账户是可以理解为数据库的一个库，实际你去安装mysql等数据库时，安装后可以创建多个库，azure的存储账户也是这个概念，每个账户下存在多个不同服务的账户。


- blob(oss)下存在一个容器，这个容器可以理解为alicloud的backup(桶)

- 这里使用go去演示，其他语言可以参考官网进行下载各自的sdk：
github.com/Azure/azure-sdk-for-go/sdk/storage/azblob


- 这里可以看到已经成功上传文件到oss了