package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
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
	if os.Getenv("SUPABASE_HOST") != "" {
		// 尝试使用 PostgreSQL (Supabase)
		pgDB, err := database.InitPostgres()
		if err != nil {
			fmt.Printf("⚠️ Supabase 连接失败: %v\n", err)
			// 仅在 MySQL 已配置时才回退，否则直接退出
			if config.GlobalConfig.Database.Host != "" {
				fmt.Println("⏳ 正在回退到 MySQL...")
				if err := database.Init(); err != nil {
					fmt.Printf("❌ MySQL 数据库连接也失败: %v\n", err)
					os.Exit(1)
				}
			} else {
				fmt.Println("❌ Supabase 连接失败且未配置 MySQL (DATABASE_HOST 为空)，请检查 SUPABASE_* 环境变量")
				os.Exit(1)
			}
		} else {
			database.DB = pgDB
			fmt.Println("✅ PostgreSQL (Supabase) 数据库连接成功")
		}
	} else {
		// 使用 MySQL (本地开发)
		if err := database.Init(); err != nil {
			fmt.Printf("❌ 数据库连接失败: %v\n", err)
			os.Exit(1)
		}
	}
	defer database.Close()

	// 3. 初始化路由
	r := router.Setup(config.GlobalConfig.Server.Mode)

	// 4. 启动HTTP服务器
	port := config.GlobalConfig.Server.Port
	if envPort := os.Getenv("PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			port = p
		}
	}
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
