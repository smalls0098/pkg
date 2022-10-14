package log

import (
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Prod       bool
	Filename   string
	MaxSize    int32
	MaxBackups int32
	MaxAge     int32
	Compress   bool
}

func NewZapLogger(config Config) *zap.Logger {
	encoder := zapcore.EncoderConfig{
		TimeKey:        "t",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	level := zapcore.DebugLevel
	ws := make([]zapcore.WriteSyncer, 0)
	ws = append(ws, zapcore.AddSync(os.Stdout)) // 打印到控制台
	if config.Prod {
		level = zapcore.InfoLevel
		ljLogger := &lumberjack.Logger{
			Filename:   config.Filename,        //指定日志存储位置
			MaxSize:    int(config.MaxSize),    //日志的最大大小（M）
			MaxBackups: int(config.MaxBackups), //日志的最大保存数量
			MaxAge:     int(config.MaxAge),     //日志文件存储最大天数
			Compress:   config.Compress,        //是否执行压缩
			LocalTime:  true,                   //使用本地时间命名 默认为utc时间命名
		}
		ws = append(ws, zapcore.AddSync(ljLogger))
	}

	ore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoder),    // 编码器配置
		zapcore.NewMultiWriteSyncer(ws...), // 多个打印器
		zap.NewAtomicLevelAt(level),        // 日志级别
	)
	zapLogger := zap.New(
		ore,
		zap.AddStacktrace(zap.NewAtomicLevelAt(zapcore.ErrorLevel)),
		zap.AddCaller(),
		zap.AddCallerSkip(2),
		zap.Development(),
	)
	return zapLogger
}
