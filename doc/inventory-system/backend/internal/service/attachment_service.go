package service

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"inventory-system/internal/model"
	"inventory-system/internal/repository"
	"inventory-system/pkg/errcode"
	"inventory-system/pkg/upload"

	"gorm.io/gorm"
)

type AttachmentService struct {
	attachmentRepo *repository.AttachmentRepo
	inventoryRepo  *repository.InventoryRepo
	logRepo        *repository.LogRepo
}

func NewAttachmentService() *AttachmentService {
	return &AttachmentService{
		attachmentRepo: repository.NewAttachmentRepo(),
		inventoryRepo:  repository.NewInventoryRepo(),
		logRepo:        repository.NewLogRepo(),
	}
}

// Upload 上传附件
func (s *AttachmentService) Upload(inventoryID uint, file *multipart.FileHeader,
	operatorID uint, operatorName string) (*model.Attachment, *errcode.ErrCode) {

	_, err := s.inventoryRepo.FindByID(inventoryID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound.WithMessage("库存记录不存在")
		}
		return nil, errcode.ErrDatabase
	}

	fileType, err := upload.ValidateFile(file)
	if err != nil {
		return nil, errcode.ErrFileFormat.WithMessage(err.Error())
	}

	src, err := file.Open()
	if err != nil {
		return nil, errcode.ErrFileUpload
	}
	mimeType, err := upload.ValidateMIME(src)
	src.Close()
	if err != nil {
		return nil, errcode.ErrFileUpload
	}

	dir, fullPath := upload.GenerateFilePath(inventoryID, file.Filename)

	if err := upload.SaveFile(file, dir, fullPath); err != nil {
		return nil, errcode.ErrFileUpload.WithMessage(err.Error())
	}

	var thumbnailPath string
	if fileType == model.FileTypeImage {
		thumbPath, err := upload.GenerateThumbnail(fullPath)
		if err == nil {
			thumbnailPath = thumbPath
		}
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	att := &model.Attachment{
		InventoryID:   inventoryID,
		FileName:      file.Filename,
		DisplayName:   file.Filename,
		FilePath:      fullPath,
		FileSize:      file.Size,
		FileType:      fileType,
		MimeType:      mimeType,
		Extension:     ext,
		ThumbnailPath: thumbnailPath,
		UploadedBy:    operatorID,
	}

	if err := s.attachmentRepo.Create(att); err != nil {
		return nil, errcode.ErrDatabase
	}

	s.logRepo.CreateOperationLog(&model.OperationLog{
		UserID:      operatorID,
		UserName:    operatorName,
		Action:      model.ActionUpload,
		TargetType:  model.TargetAttachment,
		TargetID:    &att.ID,
		Description: fmt.Sprintf("上传文件: %s (库存ID=%d)", file.Filename, inventoryID),
	})

	return att, nil
}

// GetByID 获取附件信息
func (s *AttachmentService) GetByID(id uint) (*model.Attachment, *errcode.ErrCode) {
	att, err := s.attachmentRepo.FindByID(id)
	if err != nil {
		return nil, errcode.ErrNotFound
	}
	return att, nil
}

// ListByInventoryID 获取库存的附件列表
func (s *AttachmentService) ListByInventoryID(inventoryID uint) ([]model.Attachment, *errcode.ErrCode) {
	attachments, err := s.attachmentRepo.ListByInventoryID(inventoryID)
	if err != nil {
		return nil, errcode.ErrDatabase
	}
	return attachments, nil
}

// Rename 重命名附件
func (s *AttachmentService) Rename(id uint, displayName string,
	operatorID uint, operatorName string) *errcode.ErrCode {

	att, err := s.attachmentRepo.FindByID(id)
	if err != nil {
		return errcode.ErrNotFound
	}

	oldName := att.DisplayName
	if err := s.attachmentRepo.UpdateDisplayName(id, displayName); err != nil {
		return errcode.ErrDatabase
	}

	s.logRepo.CreateOperationLog(&model.OperationLog{
		UserID:      operatorID,
		UserName:    operatorName,
		Action:      model.ActionUpdate,
		TargetType:  model.TargetAttachment,
		TargetID:    &id,
		Description: fmt.Sprintf("重命名附件: %s → %s", oldName, displayName),
	})

	return nil
}

// Delete 删除附件
func (s *AttachmentService) Delete(id uint, operatorID uint, operatorName string) *errcode.ErrCode {
	att, err := s.attachmentRepo.FindByID(id)
	if err != nil {
		return errcode.ErrNotFound
	}

	if err := s.attachmentRepo.SoftDelete(id); err != nil {
		return errcode.ErrDatabase
	}

	s.logRepo.CreateOperationLog(&model.OperationLog{
		UserID:      operatorID,
		UserName:    operatorName,
		Action:      model.ActionDelete,
		TargetType:  model.TargetAttachment,
		TargetID:    &id,
		Description: fmt.Sprintf("删除附件: %s", att.DisplayName),
	})

	return nil
}
