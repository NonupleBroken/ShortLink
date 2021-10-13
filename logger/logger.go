package logger

import (
	"ShortLink/config"
	"ShortLink/util"
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var L *zap.Logger
var S *zap.SugaredLogger

func InitLogger(logConfig *config.LogConfig) error {
	err := util.CreateDir(logConfig.LogPath)
	if err != nil {
		return err
	}
	logName := fmt.Sprintf("%s/%s.", logConfig.LogPath, logConfig.ServerLogName)
	fileWriter, err := getFileWriter(logName)
	if err != nil {
		return err
	}
	encoder := getEncoder()
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, fileWriter, zap.InfoLevel), // info以上写入文件
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zap.DebugLevel), // 全部打输出到stdout
	)
	L = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	S = L.Sugar()
	return nil
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:    "msg",
		LevelKey:      "level",
		TimeKey:       "ts",
		CallerKey:     "file",
		StacktraceKey: "trace",
		EncodeLevel:   zapcore.CapitalLevelEncoder, // level级别大写
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.000000"))
		},
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	return encoder
}

func getFileWriter(logName string) (zapcore.WriteSyncer, error) {
	writer, err := rotatelogs.New(
		logName + "%Y%m%d",
		rotatelogs.WithMaxAge(time.Duration(3600 * 24 * 30) * time.Second),
		rotatelogs.WithRotationTime(time.Duration(3600 * 24) * time.Second),
	)
	return zapcore.AddSync(writer), err
}
