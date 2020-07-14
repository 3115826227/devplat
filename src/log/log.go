package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
	"path"
)

var (
	Logger *zap.Logger
	levels = map[string]zapcore.Level{
		"DEBUG": zap.DebugLevel,
		"INFO":  zap.InfoLevel,
		"WARN":  zap.WarnLevel,
		"ERROR": zap.ErrorLevel,
	}
)

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	var logLevel zapcore.Level
	logLevelStr := os.Getenv("LOG_LEVEL")
	var exists bool
	logLevel, exists = levels[logLevelStr]
	if !exists {
		logLevel = zap.DebugLevel
	}

	logPath := os.Getenv("LOG_PATH")
	if logPath != "" {
		logDir := path.Dir(logPath)
		if _, err := os.Stat(logDir); os.IsNotExist(err) {
			log.Fatal("ERROR 日志目录 ", logDir, " 不存在")
		}
		// 打印到文件，自动分裂
		w := zapcore.AddSync(&lumberjack.Logger{
			Filename:   logPath,
			MaxSize:    64, // megabytes
			MaxBackups: 10,
			MaxAge:     28, // days
		})
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			w,
			zap.NewAtomicLevelAt(logLevel),
		)
		Logger = zap.New(core, zap.AddCaller())
	} else {
		// 打印到控制台
		cfg := zap.NewProductionConfig()
		cfg.Level = zap.NewAtomicLevelAt(logLevel)
		cfg.Encoding = "console"
		cfg.EncoderConfig = encoderConfig
		//cfg.DisableStacktrace = true
		var err error
		Logger, err = cfg.Build()
		if err != nil {
			log.Fatal("ERROR ", err)
		}
	}
}
