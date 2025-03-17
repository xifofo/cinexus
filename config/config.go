package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	RunMode      string `mapstructure:"run_mode"`
	Port         string `mapstructure:"port"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type         string `mapstructure:"type"`
	Host         string `mapstructure:"host"`
	Port         string `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	Name         string `mapstructure:"name"`
	TablePrefix  string `mapstructure:"table_prefix"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	SQLitePath   string `mapstructure:"sqlite_path"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	Issuer     string `mapstructure:"issuer"`
	ExpireTime int    `mapstructure:"expire_time"` // 过期时间（小时）
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`       // 日志级别
	Filename   string `mapstructure:"filename"`    // 日志文件名
	MaxSize    int    `mapstructure:"max_size"`    // 每个日志文件的最大大小（MB）
	MaxBackups int    `mapstructure:"max_backups"` // 保留的旧日志文件最大数量
	MaxAge     int    `mapstructure:"max_age"`     // 保留的旧日志文件最大天数
	Compress   bool   `mapstructure:"compress"`    // 是否压缩
}

// Conf 全局配置变量
var Conf = &Config{}

// Init 初始化配置
func Init() error {
	workDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// 设置配置文件路径
	viper.AddConfigPath(filepath.Join(workDir, "config"))
	viper.SetConfigName("config")
	viper.SetConfigType("toml")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析到结构体
	if err := viper.Unmarshal(Conf); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 监控配置文件变化并热加载
	viper.WatchConfig()

	return nil
}
