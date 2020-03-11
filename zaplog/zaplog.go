package zaplog

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

// Info Info
func Info(msg string, args ...interface{}) {
	logger.Info(msg)
}

// Panic Panic
func Panic(msg string, args ...interface{}) {
	logger.Panic(msg)
}

// Error Error
func Error(msg string, args ...interface{}) {
	logger.Error(msg)
}

// Warn Warn
func Warn(msg string, args ...interface{}) {
	logger.Warn(msg)
}

// Debug Debug
func Debug(msg string, args ...interface{}) {
	logger.Debug(msg)
}

// InitZapLog InitZapLog
func InitZapLog() {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,    // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // 全路径编码器
	}
	// 设置日志级别
	atom := zap.NewAtomicLevelAt(zap.DebugLevel)

	config := zap.Config{
		Level:         atom,          // 日志级别
		Development:   true,          // 开发模式，堆栈跟踪
		Encoding:      "console",     // 输出格式 console 或 json
		EncoderConfig: encoderConfig, // 编码器配置
		//InitialFields:    map[string]interface{}{"serviceName": "serviceName"}, // 初始化字段，如：添加一个服务器名称
		OutputPaths:      []string{"stdout", "./logs/info.log"}, // 输出到指定文件 stdout（标准输出，正常颜色） stderr（错误输出，红色）
		ErrorOutputPaths: []string{"stderr"},
	}

	// 构建日志
	var err error
	logger, err = config.Build()
	if err != nil {
		panic(fmt.Sprintf("log 初始化失败: %v", err))
	}
	logger.Info("log 初始化成功")
	defer logger.Sync()
}
