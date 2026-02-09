package services

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"my-blog-backend/internal/config"
	"my-blog-backend/internal/pkg/logger"

	"github.com/google/uuid"
)

type UploadService interface {
	UploadImage(file *multipart.FileHeader) (string, error)
	DeleteImage(filename string) error
}

type uploadService struct {
	config *config.UploadConfig
}

func NewUploadService(cfg *config.UploadConfig) UploadService {
	return &uploadService{
		config: cfg,
	}
}

// UploadImage 上传图片文件
func (s *uploadService) UploadImage(file *multipart.FileHeader) (string, error) {
	// 验证文件扩展名
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !s.isAllowedExt(ext) {
		return "", fmt.Errorf("不支持的文件类型: %s", ext)
	}

	// 验证文件大小
	if file.Size > s.config.Local.MaxSize {
		return "", fmt.Errorf("文件大小超过限制: 最大 %d MB", s.config.Local.MaxSize/(1024*1024))
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		logger.Error("打开文件失败", logger.Err("error", err))
		return "", fmt.Errorf("打开文件失败: %w", err)
	}
	defer src.Close()

	// 读取文件内容以生成MD5
	fileContent, err := io.ReadAll(src)
	if err != nil {
		logger.Error("读取文件内容失败", logger.Err("error", err))
		return "", fmt.Errorf("读取文件内容失败: %w", err)
	}

	// 生成MD5哈希
	hash := md5.Sum(fileContent)
	md5Hash := hex.EncodeToString(hash[:])

	// 生成唯一文件名: MD5前8位_UUID + 原始扩展名
	uniqueFilename := fmt.Sprintf("%s_%s%s", md5Hash[:8], uuid.New().String()[:8], ext)

	// 根据配置选择上传方式
	var relativePath string
	if s.config.Type == "local" {
		relativePath, err = s.uploadToLocal(uniqueFilename, fileContent)
	} else {
		relativePath, err = s.uploadToRemote(uniqueFilename, fileContent)
	}

	if err != nil {
		return "", err
	}

	// 拼接完整URL
	fullURL := s.buildFullURL(relativePath)

	logger.Info("图片上传成功",
		logger.String("filename", file.Filename),
		logger.Int64("size", file.Size),
		logger.String("md5", md5Hash),
		logger.String("url", fullURL),
	)

	return fullURL, nil
}

// uploadToLocal 上传到本地存储
func (s *uploadService) uploadToLocal(filename string, content []byte) (string, error) {
	// 获取上传目录
	uploadPath := s.config.Local.UploadPath
	if uploadPath == "" {
		uploadPath = "./uploads"
	}

	// 创建目录
	uploadDir := filepath.Join(uploadPath, "images")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		logger.Error("创建上传目录失败", logger.Err("error", err))
		return "", fmt.Errorf("创建上传目录失败: %w", err)
	}

	// 拼接完整路径
	filePath := filepath.Join(uploadDir, filename)

	// 写入文件
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		logger.Error("保存文件失败", logger.Err("error", err))
		return "", fmt.Errorf("保存文件失败: %w", err)
	}

	// 返回相对路径（使用正斜杠，Windows上filepath.Join会使用反斜杠）
	return "/uploads/images/" + filename, nil
}

// uploadToRemote 上传到远程存储（OSS/COS等）
func (s *uploadService) uploadToRemote(filename string, content []byte) (string, error) {
	// 这里可以根据不同的provider实现不同的上传逻辑
	// 目前仅作为示例，实际使用时需要集成对应的SDK

	switch s.config.Remote.Provider {
	case "oss":
		return s.uploadToOSS(filename, content)
	case "cos":
		return s.uploadToCOS(filename, content)
	case "qiniu":
		return s.uploadToQiniu(filename, content)
	default:
		return "", fmt.Errorf("不支持的远程存储服务提供商: %s", s.config.Remote.Provider)
	}
}

// uploadToOSS 上传到阿里云OSS（示例，需要安装OSS SDK）
func (s *uploadService) uploadToOSS(_ string, _ []byte) (string, error) {
	// TODO: 实际实现需要安装阿里云OSS SDK
	// import "github.com/aliyun/aliyun-oss-go-sdk/oss"

	logger.Error("阿里云OSS上传功能尚未实现，请配置本地上传")
	return "", fmt.Errorf("阿里云OSS上传功能尚未实现，请配置本地上传")

	/*
	// 示例代码
	client, err := oss.New(s.config.Remote.Endpoint, s.config.Remote.AccessKey, s.config.Remote.SecretKey)
	if err != nil {
		return "", err
	}

	bucket, err := client.Bucket(s.config.Remote.Bucket)
	if err != nil {
		return "", err
	}

	objectKey := "images/" + filename
	err = bucket.PutObject(objectKey, bytes.NewReader(content))
	if err != nil {
		return "", err
	}

	return "/" + objectKey, nil
	*/
}

// uploadToCOS 上传到腾讯云COS（示例，需要安装COS SDK）
func (s *uploadService) uploadToCOS(_ string, _ []byte) (string, error) {
	logger.Error("腾讯云COS上传功能尚未实现，请配置本地上传")
	return "", fmt.Errorf("腾讯云COS上传功能尚未实现，请配置本地上传")
}

// uploadToQiniu 上传到七牛云（示例，需要安装七牛SDK）
func (s *uploadService) uploadToQiniu(_ string, _ []byte) (string, error) {
	logger.Error("七牛云上传功能尚未实现，请配置本地上传")
	return "", fmt.Errorf("七牛云上传功能尚未实现，请配置本地上传")
}

// DeleteImage 删除图片
func (s *uploadService) DeleteImage(filename string) error {
	// 从完整URL中提取相对路径
	relativePath := s.extractRelativePath(filename)
	if relativePath == "" {
		return fmt.Errorf("无效的文件路径: %s", filename)
	}

	// 根据配置选择删除方式
	if s.config.Type == "local" {
		return s.deleteFromLocal(relativePath)
	} else {
		return s.deleteFromRemote(relativePath)
	}
}

// deleteFromLocal 从本地删除
func (s *uploadService) deleteFromLocal(relativePath string) error {
	// 移除前缀 /uploads/
	filePath := strings.TrimPrefix(relativePath, "/uploads/")

	uploadPath := s.config.Local.UploadPath
	if uploadPath == "" {
		uploadPath = "./uploads"
	}

	fullPath := filepath.Join(uploadPath, filePath)

	if err := os.Remove(fullPath); err != nil {
		logger.Error("删除文件失败", logger.Err("error", err), logger.String("path", fullPath))
		return fmt.Errorf("删除文件失败: %w", err)
	}

	logger.Info("图片删除成功", logger.String("path", fullPath))
	return nil
}

// deleteFromRemote 从远程存储删除
func (s *uploadService) deleteFromRemote(_ string) error {
	// TODO: 实现远程删除逻辑
	logger.Error("远程存储删除功能尚未实现")
	return fmt.Errorf("远程存储删除功能尚未实现")
}

// isAllowedExt 检查文件扩展名是否允许
func (s *uploadService) isAllowedExt(ext string) bool {
	if len(s.config.Local.AllowedExts) == 0 {
		return true
	}

	for _, allowedExt := range s.config.Local.AllowedExts {
		if strings.EqualFold(ext, allowedExt) {
			return true
		}
	}
	return false
}

// buildFullURL 拼接完整URL
func (s *uploadService) buildFullURL(relativePath string) string {
	if relativePath == "" {
		return ""
	}

	// 如果已经包含协议头，直接返回
	if strings.HasPrefix(relativePath, "http://") || strings.HasPrefix(relativePath, "https://") {
		return relativePath
	}

	// 使用配置的URL前缀
	urlPrefix := s.config.URLPrefix
	if urlPrefix == "" {
		urlPrefix = "http://localhost:8081"
	}

	// 确保urlPrefix不以/结尾
	urlPrefix = strings.TrimSuffix(urlPrefix, "/")

	// 确保relativePath以/开头
	if !strings.HasPrefix(relativePath, "/") {
		relativePath = "/" + relativePath
	}

	return urlPrefix + relativePath
}

// extractRelativePath 从完整URL中提取相对路径
func (s *uploadService) extractRelativePath(url string) string {
	if url == "" {
		return ""
	}

	// 如果是相对路径，直接返回
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return url
	}

	// 如果是本地URL，提取/uploads/之后的部分
	if strings.Contains(url, "/uploads/") {
		parts := strings.Split(url, "/uploads/")
		if len(parts) > 1 {
			return "/uploads/" + parts[1]
		}
	}

	// 如果是CDN URL，提取路径部分
	if s.config.Remote.CDN != "" && strings.HasPrefix(url, s.config.Remote.CDN) {
		return strings.TrimPrefix(url, s.config.Remote.CDN)
	}

	return url
}
