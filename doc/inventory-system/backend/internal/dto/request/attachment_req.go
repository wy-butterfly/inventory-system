package request

// RenameAttachmentReq 重命名附件请求
type RenameAttachmentReq struct {
	DisplayName string `json:"display_name" binding:"required,min=1,max=200"`
}
