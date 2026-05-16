package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitPostgres 初始化 PostgreSQL 连接 (Supabase)
func InitPostgres() *gorm.DB {
	// 从环境变量获取数据库配置
	host := getEnv("SUPABASE_HOST", "")
	port := getEnv("SUPABASE_PORT", "5432")
	user := getEnv("SUPABASE_USER", "")
	password := getEnv("SUPABASE_PASSWORD", "")
	dbname := getEnv("SUPABASE_DB_NAME", "")
	sslmode := getEnv("SUPABASE_SSLMODE", "require")

	if host == "" || user == "" || password == "" || dbname == "" {
		log.Fatal("Missing required database environment variables")
	}

	// 构建 DSN
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Shanghai",
		host, port, user, password, dbname, sslmode)

	// GORM 配置
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	}

	// 连接数据库
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 获取底层的 *sql.DB 对象进行连接池配置
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("✅ PostgreSQL (Supabase) connected successfully")
	return db
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
