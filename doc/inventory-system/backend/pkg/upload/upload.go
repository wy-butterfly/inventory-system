package upload

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"inventory-system/config"
	"inventory-system/internal/model"

	"github.com/google/uuid"
	"github.com/nfnt/resize"
)

// ValidateFile 校验上传的文件（格式+大小）
func ValidateFile(file *multipart.FileHeader) (string, error) {
	cfg := config.GlobalConfig.Upload

	if file.Size > cfg.MaxSize {
		return "", fmt.Errorf("文件大小超过限制（最大%dMB）", cfg.MaxSize/1024/1024)
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))

	for _, allowedExt := range cfg.AllowedImageExt {
		if ext == allowedExt {
			return model.FileTypeImage, nil
		}
	}

	for _, allowedExt := range cfg.AllowedFileExt {
		if ext == allowedExt {
			return model.FileTypeDocument, nil
		}
	}

	return "", fmt.Errorf("不支持的文件格式: %s", ext)
}

// ValidateMIME 校验文件的MIME类型
func ValidateMIME(file multipart.File) (string, error) {
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return "", err
	}
	file.Seek(0, io.SeekStart)
	mimeType := http.DetectContentType(buffer)
	return mimeType, nil
}

// GenerateFilePath 生成文件存储路径
func GenerateFilePath(inventoryID uint, originalName string) (string, string) {
	cfg := config.GlobalConfig.Upload
	ext := strings.ToLower(filepath.Ext(originalName))
	fileName := uuid.New().String() + ext
	yearMonth := time.Now().Format("200601")
	dir := filepath.Join(cfg.StoragePath, yearMonth, fmt.Sprintf("%d", inventoryID))
	fullPath := filepath.Join(dir, fileName)
	return dir, fullPath
}

// SaveFile 保存上传的文件到磁盘
func SaveFile(file *multipart.FileHeader, dir, fullPath string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("打开上传文件失败: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}

// GenerateThumbnail 生成图片缩略图
func GenerateThumbnail(originalPath string) (string, error) {
	cfg := config.GlobalConfig.Upload

	file, err := os.Open(originalPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	thumbnail := resize.Resize(uint(cfg.ThumbnailWidth), 0, img, resize.Lanczos3)

	dir := filepath.Dir(originalPath)
	base := filepath.Base(originalPath)
	thumbPath := filepath.Join(dir, "thumb_"+base)

	thumbFile, err := os.Create(thumbPath)
	if err != nil {
		return "", err
	}
	defer thumbFile.Close()

	switch format {
	case "jpeg":
		err = jpeg.Encode(thumbFile, thumbnail, &jpeg.Options{Quality: 80})
	case "png":
		err = png.Encode(thumbFile, thumbnail)
	default:
		err = jpeg.Encode(thumbFile, thumbnail, &jpeg.Options{Quality: 80})
	}

	if err != nil {
		return "", err
	}

	return thumbPath, nil
}
