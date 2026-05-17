package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config 整个应用的配置根结构
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Upload   UploadConfig   `mapstructure:"upload"`
	Log      LogConfig      `mapstructure:"log"`
}

// ServerConfig 服务器相关配置
type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

// DatabaseConfig 数据库连接配置
type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Name         string `mapstructure:"name"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

// JWTConfig JWT认证配置
type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

// UploadConfig 文件上传配置
type UploadConfig struct {
	MaxSize         int64    `mapstructure:"max_size"`
	AllowedImageExt []string `mapstructure:"allowed_image_ext"`
	AllowedFileExt  []string `mapstructure:"allowed_file_ext"`
	StoragePath     string   `mapstructure:"storage_path"`
	ThumbnailWidth  int      `mapstructure:"thumbnail_width"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level  string `mapstructure:"level"`
	Output string `mapstructure:"output"`
	Path   string `mapstructure:"path"`
}

// GlobalConfig 全局配置实例
var GlobalConfig *Config

// Load 读取并解析配置文件
func Load(path string) error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	// 设置默认值
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("jwt.expire_hours", 24)
	viper.SetDefault("upload.max_size", 20971520)
	viper.SetDefault("upload.storage_path", "./uploads")
	viper.SetDefault("upload.thumbnail_width", 200)
	viper.SetDefault("database.max_open_conns", 100)
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.output", "console")

	// 启用环境变量自动映射，如 server.port -> SERVER_PORT
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 尝试读取配置文件，不存在时仅使用默认值和环境变量
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("读取配置文件失败: %w", err)
		}
		fmt.Println("⚠️ 未找到配置文件，将使用默认值和环境变量")
	}

	GlobalConfig = &Config{}
	if err := viper.Unmarshal(GlobalConfig); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	return nil
}

// GetDSN 生成数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
	)
}
