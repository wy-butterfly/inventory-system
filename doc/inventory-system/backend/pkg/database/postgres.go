package database

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitPostgres 初始化 PostgreSQL 连接 (Supabase)
func InitPostgres() (*gorm.DB, error) {
	// 从环境变量获取数据库配置，并去除前后空白字符
	host := strings.TrimSpace(getEnv("SUPABASE_HOST", ""))
	port := strings.TrimSpace(getEnv("SUPABASE_PORT", "5432"))
	user := strings.TrimSpace(getEnv("SUPABASE_USER", ""))
	password := strings.TrimSpace(getEnv("SUPABASE_PASSWORD", ""))
	dbname := strings.TrimSpace(getEnv("SUPABASE_DB_NAME", ""))
	sslmode := strings.TrimSpace(getEnv("SUPABASE_SSLMODE", "require"))

	if host == "" || user == "" || password == "" || dbname == "" {
		return nil, fmt.Errorf("missing required database environment variables (SUPABASE_HOST, SUPABASE_USER, SUPABASE_PASSWORD, SUPABASE_DB_NAME)")
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
		return nil, fmt.Errorf("failed to connect to PostgreSQL (host=%s): %w", host, err)
	}

	// 获取底层的 *sql.DB 对象进行连接池配置
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
