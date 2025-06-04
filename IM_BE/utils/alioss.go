package utils

import (
	"bytes"
	"context"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/spf13/viper"
	"strings"
)

func UpdateFile(fileBytes []byte, objectName string, contentType string) (string, error) {
	endpoint := viper.GetString("alioss.endpoint")
	region := viper.GetString("alioss.region")
	bucketName := viper.GetString("alioss.bucket_name")
	accessKeyID := viper.GetString("alioss.access_key_id")
	accessKeySecret := viper.GetString("alioss.access_key_secret")

	// 使用配置中的 AccessKey 初始化凭证
	provider := credentials.NewStaticCredentialsProvider(
		accessKeyID,
		accessKeySecret,
	)

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(provider).WithRegion(region)

	// 创建OSS客户端
	client := oss.NewClient(cfg)

	// 创建上传对象的请求
	request := &oss.PutObjectRequest{
		Bucket:      oss.Ptr(bucketName),        // 存储空间名称
		Key:         oss.Ptr(objectName),        // 对象名称
		Body:        bytes.NewReader(fileBytes), // 要上传的内容
		ContentType: oss.Ptr(contentType),
	}

	if _, err := client.PutObject(context.Background(), request); err != nil {
		return "", err
	}

	// TODO: 为什么使用 builder
	/*
		性能更高：
		+ 每次拼接都会创建一个新字符串，涉及多次内存分配；
		strings.Builder 只分配一次内存，内部维护一个可扩展的缓冲区。

		内存效率更好：
		多次拼接时不反复复制旧数据，减少 GC 压力。

		代码更清晰：
		明确表示“这是在构造字符串”，比多行 + 拼接更语义化。
	*/
	sb := strings.Builder{}
	// 拼接 URL
	sb.WriteString("https://")
	sb.WriteString(bucketName)
	sb.WriteString(".")
	sb.WriteString(endpoint)
	sb.WriteString("/")
	sb.WriteString(objectName)

	return sb.String(), nil
}
