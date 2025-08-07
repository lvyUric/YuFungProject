package logger

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *logrus.Logger

// LogConfig 日志配置
type LogConfig struct {
	Level      string `yaml:"level"`       // 日志级别
	Format     string `yaml:"format"`      // 日志格式: json/text
	Output     string `yaml:"output"`      // 输出: stdout/file/both
	FilePath   string `yaml:"file_path"`   // 文件路径
	MaxSize    int    `yaml:"max_size"`    // 单个文件最大大小(MB)
	MaxAge     int    `yaml:"max_age"`     // 文件保留天数
	MaxBackups int    `yaml:"max_backups"` // 最大备份文件数
	Compress   bool   `yaml:"compress"`    // 是否压缩旧文件
}

// InitLogger 初始化日志系统
func InitLogger(config LogConfig) error {
	Log = logrus.New()

	// 设置日志级别
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	Log.SetLevel(level)

	// 设置日志格式
	if config.Format == "json" {
		Log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		Log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     true,
		})
	}

	// 配置输出
	var writers []io.Writer

	// 控制台输出
	if config.Output == "stdout" || config.Output == "both" {
		writers = append(writers, os.Stdout)
	}

	// 文件输出
	if config.Output == "file" || config.Output == "both" {
		// 确保日志目录存在
		logDir := filepath.Dir(config.FilePath)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return err
		}

		// 配置日志文件轮转
		fileWriter := &lumberjack.Logger{
			Filename:   config.FilePath,
			MaxSize:    config.MaxSize,    // 单个文件最大大小(MB)
			MaxAge:     config.MaxAge,     // 文件保留天数
			MaxBackups: config.MaxBackups, // 最大备份文件数
			Compress:   config.Compress,   // 是否压缩
			LocalTime:  true,              // 使用本地时间
		}
		writers = append(writers, fileWriter)
	}

	// 设置多重输出
	if len(writers) > 0 {
		Log.SetOutput(io.MultiWriter(writers...))
	}

	Log.Info("日志系统初始化完成")
	return nil
}

// WithFields 创建带字段的日志条目
func WithFields(fields logrus.Fields) *logrus.Entry {
	return Log.WithFields(fields)
}

// WithField 创建带单个字段的日志条目
func WithField(key string, value interface{}) *logrus.Entry {
	return Log.WithField(key, value)
}

// Info 信息日志
func Info(args ...interface{}) {
	Log.Info(args...)
}

// Infof 格式化信息日志
func Infof(format string, args ...interface{}) {
	Log.Infof(format, args...)
}

// Debug 调试日志
func Debug(args ...interface{}) {
	Log.Debug(args...)
}

// Debugf 格式化调试日志
func Debugf(format string, args ...interface{}) {
	Log.Debugf(format, args...)
}

// Warn 警告日志
func Warn(args ...interface{}) {
	Log.Warn(args...)
}

// Warnf 格式化警告日志
func Warnf(format string, args ...interface{}) {
	Log.Warnf(format, args...)
}

// Error 错误日志
func Error(args ...interface{}) {
	Log.Error(args...)
}

// Errorf 格式化错误日志
func Errorf(format string, args ...interface{}) {
	Log.Errorf(format, args...)
}

// Fatal 致命错误日志
func Fatal(args ...interface{}) {
	Log.Fatal(args...)
}

// Fatalf 格式化致命错误日志
func Fatalf(format string, args ...interface{}) {
	Log.Fatalf(format, args...)
}

// Panic panic日志
func Panic(args ...interface{}) {
	Log.Panic(args...)
}

// Panicf 格式化panic日志
func Panicf(format string, args ...interface{}) {
	Log.Panicf(format, args...)
}

// DBLog 数据库操作日志
func DBLog(operation, collection string, filter interface{}, duration time.Duration) {
	WithFields(logrus.Fields{
		"type":       "database",
		"operation":  operation,
		"collection": collection,
		"filter":     filter,
		"duration":   duration.String(),
	}).Info("数据库操作")
}

// APILog API请求日志
func APILog(method, path, clientIP, userAgent string, statusCode int, duration time.Duration) {
	WithFields(logrus.Fields{
		"type":        "api",
		"method":      method,
		"path":        path,
		"client_ip":   clientIP,
		"user_agent":  userAgent,
		"status_code": statusCode,
		"duration":    duration.String(),
	}).Info("API请求")
}

// AuthLog 认证相关日志
func AuthLog(action, username, clientIP string, success bool, message string) {
	level := Log.Info
	if !success {
		level = Log.Warn
	}

	level(WithFields(logrus.Fields{
		"type":      "auth",
		"action":    action,
		"username":  username,
		"client_ip": clientIP,
		"success":   success,
		"message":   message,
	}))
}

// BusinessLog 业务操作日志
func BusinessLog(module, action, userID, details string) {
	WithFields(logrus.Fields{
		"type":    "business",
		"module":  module,
		"action":  action,
		"user_id": userID,
		"details": details,
	}).Info("业务操作")
}
