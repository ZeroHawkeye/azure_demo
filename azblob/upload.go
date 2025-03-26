package azblob

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

// UploadFileToAzure 上传文件到Azure Blob存储
func UploadFileToAzure(localFilePath, accountName, accountKey, containerName string) error {
	// 目标blob名称
	blobName := filepath.Base(localFilePath)

	// 创建凭证
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return fmt.Errorf("无法创建凭证: %v", err)
	}

	// 创建客户端
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)
	client, err := azblob.NewClientWithSharedKeyCredential(serviceURL, credential, nil)
	if err != nil {
		return fmt.Errorf("无法创建客户端: %v", err)
	}

	// 打开文件
	file, err := os.Open(localFilePath)
	if err != nil {
		return fmt.Errorf("无法打开文件: %v", err)
	}
	defer file.Close()

	// 上传文件
	ctx := context.Background()
	_, err = client.UploadFile(ctx, containerName, blobName, file, &azblob.UploadFileOptions{
		BlockSize:   int64(4 * 1024 * 1024), // 4MB块大小
		Concurrency: 4,                      // 并发数
		// 可以设置Content-Type等元数据；这里的数据就是数据基础信息，可以忽略，不知道的话就不需要写
		// HTTPHeaders: &azblob.BlobHTTPHeaders{
		//     ContentType: "text/plain",
		// },
	})

	if err != nil {
		return fmt.Errorf("上传失败: %v", err)
	}

	fmt.Printf("成功上传文件 %s 到容器 %s 中的 %s\n", localFilePath, containerName, blobName)
	return nil
}
