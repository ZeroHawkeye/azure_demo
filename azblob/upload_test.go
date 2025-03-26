package azblob

import (
	"fmt"
	"log"
	"testing"
)

func TestUploadFileToAzure(t *testing.T) {
	// Azure存储账号信息
	// 这些凭证可以在Azure门户网站中找到:
	// 1. 登录Azure门户(https://portal.azure.com)
	// 2. 导航到你的存储账号
	// 3. 在"访问密钥"部分可以找到存储账号名称和密钥
	// 4. 或者在"共享访问签名"部分创建SAS令牌
	accountName := ""   // 存储账号(可以理解为数据库)名称
	accountKey := ""    // 密钥
	containerName := "" // 容器(可以理解为数据库中的表)名称

	// 本地文件路径
	localFilePath := "learn.txt"

	// 调用azblob包中的上传函数
	err := UploadFileToAzure(localFilePath, accountName, accountKey, containerName)
	if err != nil {
		log.Fatalf("上传失败: %v", err)
	}

	fmt.Println("文件上传完成")
}
