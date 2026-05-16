package service

import (
	"fmt"
	"strconv"

	"inventory-system/internal/dto/request"
	"inventory-system/internal/model"
	"inventory-system/internal/repository"
	"inventory-system/pkg/errcode"

	"gorm.io/gorm"
)

type InventoryService struct {
	inventoryRepo  *repository.InventoryRepo
	attachmentRepo *repository.AttachmentRepo
	logRepo        *repository.LogRepo
	userRepo       *repository.UserRepo
}

func NewInventoryService() *InventoryService {
	return &InventoryService{
		inventoryRepo:  repository.NewInventoryRepo(),
		attachmentRepo: repository.NewAttachmentRepo(),
		logRepo:        repository.NewLogRepo(),
		userRepo:       repository.NewUserRepo(),
	}
}

// Create 新增库存记录
func (s *InventoryService) Create(req *request.CreateInventoryReq, operatorID uint, operatorName string) (*model.Inventory, *errcode.ErrCode) {
	if !model.IsValidCategory(req.Category) {
		return nil, errcode.ErrInvalidParam.WithMessage("无效的分类值")
	}

	existing, err := s.inventoryRepo.FindByMaterialCode(req.MaterialCode)
	if err == nil && existing != nil {
		return nil, errcode.ErrDuplicate.WithMessage("物料编码已存在")
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errcode.ErrDatabase
	}

	status := int8(1)
	if req.Status != nil {
		status = *req.Status
	}

	inv := &model.Inventory{
		MaterialCode:  req.MaterialCode,
		MaterialName:  req.MaterialName,
		Category:      req.Category,
		Specification: req.Specification,
		Unit:          req.Unit,
		Quantity:      req.Quantity,
		SafetyStock:   req.SafetyStock,
		Location:      req.Location,
		Status:        status,
		Remark:        req.Remark,
		Version:       1,
		CreatedBy:     operatorID,
	}

	if err := s.inventoryRepo.Create(inv); err != nil {
		return nil, errcode.ErrDatabase
	}

	s.logRepo.CreateOperationLog(&model.OperationLog{
		UserID:      operatorID,
		UserName:    operatorName,
		Action:      model.ActionCreate,
		TargetType:  model.TargetInventory,
		TargetID:    &inv.ID,
		Description: fmt.Sprintf("新增库存: %s (%s)", inv.MaterialName, inv.MaterialCode),
	})

	return inv, nil
}

// GetByID 获取库存详情
func (s *InventoryService) GetByID(id uint) (*model.Inventory, *errcode.ErrCode) {
	inv, err := s.inventoryRepo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrDatabase
	}

	if inv.SafetyStock > 0 && inv.Quantity < inv.SafetyStock {
		inv.Warning = true
	}

	attachments, _ := s.attachmentRepo.ListByInventoryID(id)
	inv.Attachments = attachments

	nameMap, _ := s.userRepo.GetNamesByIDs([]uint{inv.CreatedBy})
	if name, ok := nameMap[inv.CreatedBy]; ok {
		inv.CreatorName = name
	}

	return inv, nil
}

// Update 编辑库存记录
func (s *InventoryService) Update(id uint, req *request.UpdateInventoryReq, operatorID uint, operatorName string) *errcode.ErrCode {
	oldInv, err := s.inventoryRepo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errcode.ErrNotFound
		}
		return errcode.ErrDatabase
	}

	var changeLogs []model.ChangeLog

	addChange := func(field, oldVal, newVal string) {
		if oldVal != newVal {
			label := model.FieldLabelMap[field]
			if label == "" {
				label = field
			}
			changeLogs = append(changeLogs, model.ChangeLog{
				InventoryID: id,
				UserID:      operatorID,
				UserName:    operatorName,
				FieldName:   field,
				FieldLabel:  label,
				OldValue:    oldVal,
				NewValue:    newVal,
			})
		}
	}

	if req.MaterialName != nil {
		addChange("material_name", oldInv.MaterialName, *req.MaterialName)
		oldInv.MaterialName = *req.MaterialName
	}
	if req.Category != nil {
		if !model.IsValidCategory(*req.Category) {
			return errcode.ErrInvalidParam.WithMessage("无效的分类值")
		}
		addChange("category", oldInv.Category, *req.Category)
		oldInv.Category = *req.Category
	}
	if req.Specification != nil {
		addChange("specification", oldInv.Specification, *req.Specification)
		oldInv.Specification = *req.Specification
	}
	if req.Unit != nil {
		addChange("unit", oldInv.Unit, *req.Unit)
		oldInv.Unit = *req.Unit
	}
	if req.Quantity != nil {
		addChange("quantity",
			strconv.FormatFloat(oldInv.Quantity, 'f', 2, 64),
			strconv.FormatFloat(*req.Quantity, 'f', 2, 64))
		oldInv.Quantity = *req.Quantity
	}
	if req.SafetyStock != nil {
		addChange("safety_stock",
			strconv.FormatFloat(oldInv.SafetyStock, 'f', 2, 64),
			strconv.FormatFloat(*req.SafetyStock, 'f', 2, 64))
		oldInv.SafetyStock = *req.SafetyStock
	}
	if req.Location != nil {
		addChange("location", oldInv.Location, *req.Location)
		oldInv.Location = *req.Location
	}
	if req.Status != nil {
		addChange("status",
			strconv.Itoa(int(oldInv.Status)),
			strconv.Itoa(int(*req.Status)))
		oldInv.Status = *req.Status
	}
	if req.Remark != nil {
		addChange("remark", oldInv.Remark, *req.Remark)
		oldInv.Remark = *req.Remark
	}

	oldInv.UpdatedBy = &operatorID

	if err := s.inventoryRepo.Update(oldInv, req.Version); err != nil {
		if err == gorm.ErrRecordNotFound {
			return errcode.ErrOptimisticLock
		}
		return errcode.ErrDatabase
	}

	if len(changeLogs) > 0 {
		s.logRepo.CreateChangeLogs(changeLogs)
	}

	s.logRepo.CreateOperationLog(&model.OperationLog{
		UserID:      operatorID,
		UserName:    operatorName,
		Action:      model.ActionUpdate,
		TargetType:  model.TargetInventory,
		TargetID:    &id,
		Description: fmt.Sprintf("编辑库存: %s，变更%d个字段", oldInv.MaterialCode, len(changeLogs)),
	})

	return nil
}

// Delete 软删除库存记录
func (s *InventoryService) Delete(id uint, operatorID uint, operatorName string) *errcode.ErrCode {
	inv, err := s.inventoryRepo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errcode.ErrNotFound
		}
		return errcode.ErrDatabase
	}

	if err := s.inventoryRepo.SoftDelete(id, operatorID); err != nil {
		return errcode.ErrDatabase
	}

	s.attachmentRepo.SoftDeleteByInventoryID(id)

	s.logRepo.CreateOperationLog(&model.OperationLog{
		UserID:      operatorID,
		UserName:    operatorName,
		Action:      model.ActionDelete,
		TargetType:  model.TargetInventory,
		TargetID:    &id,
		Description: fmt.Sprintf("删除库存: %s (%s)", inv.MaterialName, inv.MaterialCode),
	})

	return nil
}

// BatchDelete 批量软删除
func (s *InventoryService) BatchDelete(ids []uint, operatorID uint, operatorName string) *errcode.ErrCode {
	if err := s.inventoryRepo.BatchSoftDelete(ids, operatorID); err != nil {
		return errcode.ErrDatabase
	}

	for _, id := range ids {
		s.attachmentRepo.SoftDeleteByInventoryID(id)
	}

	s.logRepo.CreateOperationLog(&model.OperationLog{
		UserID:      operatorID,
		UserName:    operatorName,
		Action:      model.ActionDelete,
		TargetType:  model.TargetInventory,
		Description: fmt.Sprintf("批量删除%d条库存记录", len(ids)),
	})

	return nil
}

// Restore 恢复已删除的记录
func (s *InventoryService) Restore(id uint, operatorID uint, operatorName string) *errcode.ErrCode {
	if err := s.inventoryRepo.Restore(id); err != nil {
		return errcode.ErrDatabase
	}

	s.logRepo.CreateOperationLog(&model.OperationLog{
		UserID:      operatorID,
		UserName:    operatorName,
		Action:      model.ActionRestore,
		TargetType:  model.TargetInventory,
		TargetID:    &id,
		Description: fmt.Sprintf("恢复库存记录 ID=%d", id),
	})

	return nil
}

// List 分页查询库存列表
func (s *InventoryService) List(query *request.InventoryQuery) ([]model.Inventory, int64, *errcode.ErrCode) {
	inventories, total, err := s.inventoryRepo.List(
		query.Keyword, query.Category, query.Status,
		query.MinQuantity, query.MaxQuantity, query.Location,
		query.ShowDeleted,
		query.GetOffset(), query.GetPageSize(),
	)
	if err != nil {
		return nil, 0, errcode.ErrDatabase
	}

	var creatorIDs []uint
	for i := range inventories {
		if inventories[i].SafetyStock > 0 && inventories[i].Quantity < inventories[i].SafetyStock {
			inventories[i].Warning = true
		}
		creatorIDs = append(creatorIDs, inventories[i].CreatedBy)
	}

	if len(creatorIDs) > 0 {
		nameMap, _ := s.userRepo.GetNamesByIDs(creatorIDs)
		for i := range inventories {
			if name, ok := nameMap[inventories[i].CreatedBy]; ok {
				inventories[i].CreatorName = name
			}
		}
	}

	return inventories, total, nil
}
