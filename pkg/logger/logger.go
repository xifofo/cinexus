package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"cinexus/config"
)

var logger *zap.Logger

// Init 初始化日志
func Init() error {
	// 确保日志目录存在
	logDir := filepath.Dir(config.Conf.Log.Filename)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	// 设置日志级别
	var level zapcore.Level
	switch config.Conf.Log.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// 配置编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 配置输出
	// 控制台输出
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	consoleOutput := zapcore.Lock(os.Stdout)

	// 文件输出 - 使用lumberjack实现按天分割
	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)

	// 获取当前日期作为日志文件名的一部分
	today := time.Now().Format("2006-01-02")

	// 构建日志文件名
	logFilename := config.Conf.Log.Filename
	ext := filepath.Ext(logFilename)
	logFileNameOnly := logFilename[:len(logFilename)-len(ext)]
	dailyLogFile := logFileNameOnly + "." + today + ext

	// 配置lumberjack
	fileOutput := zapcore.AddSync(&lumberjack.Logger{
		Filename:   dailyLogFile,               // 按日期命名的日志文件
		MaxSize:    config.Conf.Log.MaxSize,    // 每个日志文件的最大大小（MB）
		MaxBackups: config.Conf.Log.MaxBackups, // 保留的旧日志文件最大数量
		MaxAge:     config.Conf.Log.MaxAge,     // 保留的旧日志文件最大天数
		Compress:   config.Conf.Log.Compress,   // 是否压缩
		LocalTime:  true,                       // 使用本地时间
	})

	// 创建一个软链接指向最新的日志文件
	if err := createSymlink(dailyLogFile, logFilename); err != nil {
		// 仅记录错误，不中断程序
		fmt.Println("创建日志软链接失败:", err)
	}

	// 创建核心
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleOutput, level),
		zapcore.NewCore(fileEncoder, fileOutput, level),
	)

	// 创建日志记录器
	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return nil
}

// createSymlink 创建软链接指向最新的日志文件
func createSymlink(source, target string) error {
	// 如果目标文件已存在，先删除
	os.Remove(target)

	// 创建软链接
	return os.Symlink(filepath.Base(source), target)
}

// 自定义时间编码器，格式：2006-01-02 15:04:05
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

// Sync 同步日志
func Sync() {
	if logger != nil {
		_ = logger.Sync()
	}
}

// Debug 输出debug级别日志
func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

// Info 输出info级别日志
func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

// Warn 输出warn级别日志
func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

// Error 输出error级别日志
func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

// Fatal 输出fatal级别日志
func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}

// With 创建一个子日志记录器，附加字段
func With(fields ...zap.Field) *zap.Logger {
	return logger.With(fields...)
}
