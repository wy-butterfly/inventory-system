package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"inventory-system/config"
	"inventory-system/internal/router"
	"inventory-system/pkg/database"
)

func main() {
	// 1. 加载配置文件
	fmt.Println("⏳ 正在加载配置...")
	if err := config.Load("."); err != nil {
		fmt.Printf("❌ 加载配置失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ 配置加载成功")

	// 2. 初始化数据库连接
	fmt.Println("⏳ 正在连接数据库...")
	if err := database.Init(); err != nil {
		fmt.Printf("❌ 数据库连接失败: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	// 3. 初始化路由
	r := router.Setup(config.GlobalConfig.Server.Mode)

	// 4. 启动HTTP服务器
	port := config.GlobalConfig.Server.Port
	fmt.Printf("🚀 服务器启动成功，监听端口: %d\n", port)
	fmt.Printf("   API地址: http://localhost:%d/api/v1\n", port)
	fmt.Printf("   健康检查: http://localhost:%d/api/health\n", port)
	fmt.Println("   按 Ctrl+C 停止服务器")

	go func() {
		if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
			fmt.Printf("❌ 服务器启动失败: %v\n", err)
			os.Exit(1)
		}
	}()

	// 5. 优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("\n🛑 服务器正在关闭...")
}
