package database

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
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
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Shanghai connect_timeout=10",
		host, port, user, password, dbname, sslmode)

	// 解析 pgx 配置
	connConfig, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PostgreSQL DSN: %w", err)
	}

	// 强制优先使用 IPv4，避免 Render 等平台因 IPv6 不可达导致连接失败
	connConfig.LookupFunc = func(ctx context.Context, h string) ([]string, error) {
		addrs, err := net.DefaultResolver.LookupIPAddr(ctx, h)
		if err != nil {
			return nil, err
		}
		ipv4 := make([]string, 0, len(addrs))
		ipv6 := make([]string, 0, len(addrs))
		for _, a := range addrs {
			if a.IP.To4() != nil {
				ipv4 = append(ipv4, a.IP.String())
			} else {
				ipv6 = append(ipv6, a.IP.String())
			}
		}
		if len(ipv4) > 0 {
			return ipv4, nil
		}
		return ipv6, nil
	}

	// GORM 配置
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	}

	// Transaction pooler (port 6543) 不支持 prepared statements，需使用 simple protocol
	connConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	// 使用自定义 pgx 配置打开连接
	sqlDB := stdlib.OpenDB(*connConfig)
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn:                 sqlDB,
		PreferSimpleProtocol: true,
	}), gormConfig)
	if err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("failed to connect to PostgreSQL (host=%s): %w", host, err)
	}

	// 设置连接池参数并验证连通性
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := sqlDB.Ping(); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("failed to connect to PostgreSQL (host=%s): %w", host, err)
	}

	return db, nil
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
