package oss

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// OSSClient 阿里云 OSS 客户端
type OSSClient struct {
	client     *oss.Client
	bucket     *oss.Bucket
	bucketName string
	domain     string
}

// NewOSSClient 创建 OSS 客户端
func NewOSSClient() (*OSSClient, error) {
	endpoint := os.Getenv("OSS_ENDPOINT")
	accessKeyID := os.Getenv("OSS_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("OSS_ACCESS_KEY_SECRET")
	bucketName := os.Getenv("OSS_BUCKET_NAME")
	bucketDomain := os.Getenv("OSS_BUCKET_DOMAIN")

	if endpoint == "" || accessKeyID == "" || accessKeySecret == "" || bucketName == "" {
		return nil, fmt.Errorf("missing required OSS environment variables")
	}

	// 创建 OSS 客户端
	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("failed to create OSS client: %v", err)
	}

	// 获取 Bucket
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to get bucket: %v", err)
	}

	// 如果没有设置域名，使用默认的 OSS 域名
	if bucketDomain == "" {
		bucketDomain = fmt.Sprintf("https://%s.%s", bucketName, endpoint)
	}

	return &OSSClient{
		client:     client,
		bucket:     bucket,
		bucketName: bucketName,
		domain:     bucketDomain,
	}, nil
}

// UploadFile 上传文件到 OSS
func (o *OSSClient) UploadFile(file io.Reader, objectKey string, contentType string) (string, int64, error) {
	// 设置文件属性
	options := []oss.Option{
		oss.ContentType(contentType),
		oss.ObjectACL(oss.ACLPublicRead), // 设置为公共读
	}

	// 上传文件
	err := o.bucket.PutObject(objectKey, file, options...)
	if err != nil {
		return "", 0, fmt.Errorf("failed to upload file: %v", err)
	}

	// 获取文件信息
	props, err := o.bucket.GetObjectDetailedMeta(objectKey)
	if err != nil {
		return "", 0, fmt.Errorf("failed to get object properties: %v", err)
	}

	// 获取文件大小
	size := props.ContentLength

	// 返回文件 URL
	fileURL := fmt.Sprintf("%s/%s", strings.TrimSuffix(o.domain, "/"), objectKey)

	return fileURL, size, nil
}

// DeleteFile 删除 OSS 文件
func (o *OSSClient) DeleteFile(objectKey string) error {
	err := o.bucket.DeleteObject(objectKey)
	if err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}
	return nil
}

// GetFileURL 获取文件访问 URL
func (o *OSSClient) GetFileURL(objectKey string) string {
	return fmt.Sprintf("%s/%s", strings.TrimSuffix(o.domain, "/"), objectKey)
}

// GeneratePresignedURL 生成预签名 URL (用于私有文件)
func (o *OSSClient) GeneratePresignedURL(objectKey string, expiration time.Duration) (string, error) {
	signedURL, err := o.bucket.SignURL(objectKey, oss.HTTPGet, int64(expiration.Seconds()))
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %v", err)
	}
	return signedURL, nil
}

// GenerateObjectKey 生成对象键 (文件路径)
func GenerateObjectKey(originalName string) string {
	// 获取文件扩展名
	ext := filepath.Ext(originalName)
	
	// 生成基于时间的路径
	now := time.Now()
	path := fmt.Sprintf("uploads/%d/%02d", now.Year(), now.Month())
	
	// 生成唯一文件名
	timestamp := now.Unix()
	filename := fmt.Sprintf("%d%s", timestamp, ext)
	
	return fmt.Sprintf("%s/%s", path, filename)
}

// GetFileExtension 获取文件扩展名
func GetFileExtension(filename string) string {
	return strings.ToLower(filepath.Ext(filename))
}

// IsImageType 判断是否为图片类型
func IsImageType(filename string) bool {
	ext := GetFileExtension(filename)
	imageTypes := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp"}
	for _, imgType := range imageTypes {
		if ext == imgType {
			return true
		}
	}
	return false
}

// GetContentType 根据 MIME 类型获取 Content-Type
func GetContentType(filename string) string {
	ext := GetFileExtension(filename)
	contentTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".bmp":  "image/bmp",
		".webp": "image/webp",
		".pdf":  "application/pdf",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".xls":  "application/vnd.ms-excel",
		".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		".txt":  "text/plain",
		".csv":  "text/csv",
	}
	
	if contentType, exists := contentTypes[ext]; exists {
		return contentType
	}
	return "application/octet-stream"
}
