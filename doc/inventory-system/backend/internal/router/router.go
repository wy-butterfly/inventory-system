package router

import (
	"inventory-system/internal/handler"
	"inventory-system/internal/middleware"
	"inventory-system/internal/model"

	"github.com/gin-gonic/gin"
)

// Setup 初始化路由
func Setup(mode string) *gin.Engine {
	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// 全局中间件
	r.Use(middleware.Recovery())
	r.Use(middleware.RequestLogger())
	r.Use(middleware.CORS())

	// 创建Handler实例
	authHandler := handler.NewAuthHandler()
	inventoryHandler := handler.NewInventoryHandler()
	attachmentHandler := handler.NewAttachmentHandler()
	userHandler := handler.NewUserHandler()
	logHandler := handler.NewLogHandler()
	dashboardHandler := handler.NewDashboardHandler()

	// 健康检查（不需要登录）
	r.GET("/api/health", dashboardHandler.HealthCheck)

	// API v1 路由组
	v1 := r.Group("/api/v1")

	// 认证接口（不需要登录）
	auth := v1.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
	}

	// 以下接口都需要登录
	v1.Use(middleware.JWTAuth())

	// 认证接口（需要登录）
	{
		v1.POST("/auth/logout", authHandler.Logout)
		v1.GET("/auth/me", authHandler.GetMe)
	}

	// 仪表盘
	{
		v1.GET("/dashboard", dashboardHandler.GetDashboard)
	}

	// 库存接口
	inventories := v1.Group("/inventories")
	{
		// 所有角色都可以查询
		inventories.GET("", inventoryHandler.List)
		inventories.GET("/:id", inventoryHandler.GetByID)
		inventories.GET("/export", inventoryHandler.Export)

		// 管理员和操作员可以新增、编辑
		inventories.POST("",
			middleware.RequireRoles(model.RoleAdmin, model.RoleOperator),
			inventoryHandler.Create)
		inventories.PUT("/:id",
			middleware.RequireRoles(model.RoleAdmin, model.RoleOperator),
			inventoryHandler.Update)

		// 删除（仅管理员）
		inventories.DELETE("/:id",
			middleware.RequireAdmin(),
			inventoryHandler.Delete)
		inventories.POST("/batch-delete",
			middleware.RequireAdmin(),
			inventoryHandler.BatchDelete)
		inventories.PUT("/:id/restore",
			middleware.RequireAdmin(),
			inventoryHandler.Restore)

		// 附件接口
		inventories.POST("/:id/attachments",
			middleware.RequireRoles(model.RoleAdmin, model.RoleOperator),
			attachmentHandler.Upload)
		inventories.GET("/:id/attachments", attachmentHandler.List)
	}

	// 附件操作接口
	attachments := v1.Group("/attachments")
	{
		attachments.GET("/:id/download", attachmentHandler.Download)
		attachments.GET("/:id/preview", attachmentHandler.Preview)
		attachments.PUT("/:id",
			middleware.RequireRoles(model.RoleAdmin, model.RoleOperator),
			attachmentHandler.Rename)
		attachments.DELETE("/:id",
			middleware.RequireAdmin(),
			attachmentHandler.Delete)
	}

	// 用户管理接口（仅管理员）
	users := v1.Group("/users", middleware.RequireAdmin())
	{
		users.GET("", userHandler.List)
		users.POST("", userHandler.Create)
		users.PUT("/:id", userHandler.Update)
		users.PUT("/:id/reset-password", userHandler.ResetPassword)
	}

	// 日志接口（仅管理员）
	logs := v1.Group("/logs", middleware.RequireAdmin())
	{
		logs.GET("/operations", logHandler.ListOperationLogs)
		logs.GET("/changes", logHandler.ListChangeLogs)
	}

	return r
}
