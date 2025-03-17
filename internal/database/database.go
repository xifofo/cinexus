package database

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"cinexus/config"
	"cinexus/pkg/logger"
)

// DB 全局数据库连接
var DB *gorm.DB

// Init 初始化数据库连接
func Init() error {
	var err error
	var dialector gorm.Dialector

	switch config.Conf.Database.Type {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.Conf.Database.User,
			config.Conf.Database.Password,
			config.Conf.Database.Host,
			config.Conf.Database.Port,
			config.Conf.Database.Name)
		dialector = mysql.Open(dsn)
	case "sqlite":
		// 确保SQLite数据库目录存在
		dbDir := filepath.Dir(config.Conf.Database.SQLitePath)
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return fmt.Errorf("创建SQLite数据库目录失败: %w", err)
		}
		dialector = sqlite.Open(config.Conf.Database.SQLitePath)
	default:
		return fmt.Errorf("不支持的数据库类型: %s", config.Conf.Database.Type)
	}

	// 配置GORM
	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   config.Conf.Database.TablePrefix,
			SingularTable: true,
		},
		Logger: gormlogger.New(
			&GormLogWriter{},
			gormlogger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  getGormLogLevel(),
				IgnoreRecordNotFoundError: true,
				Colorful:                  false,
			},
		),
	}

	// 连接数据库
	DB, err = gorm.Open(dialector, gormConfig)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// 配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %w", err)
	}

	sqlDB.SetMaxIdleConns(config.Conf.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.Conf.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return nil
}

// Close 关闭数据库连接
func Close() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

// GormLogWriter 实现gorm日志写入器
type GormLogWriter struct{}

// Printf 实现gorm日志写入
func (w *GormLogWriter) Printf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	logger.Info("GORM", zap.String("sql", msg))
}

// 根据配置获取GORM日志级别
func getGormLogLevel() gormlogger.LogLevel {
	switch config.Conf.Log.Level {
	case "debug":
		return gormlogger.Info
	case "info":
		return gormlogger.Info
	case "warn":
		return gormlogger.Warn
	case "error":
		return gormlogger.Error
	default:
		return gormlogger.Info
	}
}
