package service

import (
	"fmt"
	"strconv"

	"inventory-system/internal/dto/request"
	"inventory-system/internal/model"
	"inventory-system/internal/repository"
	"inventory-system/pkg/errcode"

	"github.com/xuri/excelize/v2"
)

type ExportService struct {
	inventoryRepo *repository.InventoryRepo
	logRepo       *repository.LogRepo
}

func NewExportService() *ExportService {
	return &ExportService{
		inventoryRepo: repository.NewInventoryRepo(),
		logRepo:       repository.NewLogRepo(),
	}
}

// ExportInventories 导出库存数据为Excel
func (s *ExportService) ExportInventories(query *request.InventoryQuery,
	operatorID uint, operatorName string) (*excelize.File, *errcode.ErrCode) {

	inventories, err := s.inventoryRepo.ListAll(
		query.Keyword, query.Category, query.Status,
		query.MinQuantity, query.MaxQuantity,
	)
	if err != nil {
		return nil, errcode.ErrDatabase
	}

	f := excelize.NewFile()
	sheetName := "库存数据"
	f.SetSheetName("Sheet1", sheetName)

	headers := []string{"物料编码", "物料名称", "分类", "规格型号", "单位",
		"库存数量", "安全库存", "仓库位置", "状态", "备注", "创建时间", "更新时间"}
	for i, h := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, h)
	}

	categoryMap := map[string]string{
		"raw_material": "原料", "finished_product": "成品",
		"spare_part": "备件", "other": "其他",
	}
	for i, inv := range inventories {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), inv.MaterialCode)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), inv.MaterialName)
		catName := categoryMap[inv.Category]
		if catName == "" {
			catName = inv.Category
		}
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), catName)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), inv.Specification)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), inv.Unit)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), strconv.FormatFloat(inv.Quantity, 'f', 2, 64))
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), strconv.FormatFloat(inv.SafetyStock, 'f', 2, 64))
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), inv.Location)
		statusText := "启用"
		if inv.Status == 0 {
			statusText = "停用"
		}
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), statusText)
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), inv.Remark)
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", row), inv.CreatedAt.Format("2006-01-02 15:04:05"))
		f.SetCellValue(sheetName, fmt.Sprintf("L%d", row), inv.UpdatedAt.Format("2006-01-02 15:04:05"))
	}

	s.logRepo.CreateOperationLog(&model.OperationLog{
		UserID:      operatorID,
		UserName:    operatorName,
		Action:      model.ActionExport,
		TargetType:  model.TargetInventory,
		Description: fmt.Sprintf("导出库存数据 %d 条", len(inventories)),
	})

	return f, nil
}
