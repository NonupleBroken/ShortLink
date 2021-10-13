package web

import (
	"ShortLink/config"
	"ShortLink/util"
	"fmt"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

var webLogger *zap.SugaredLogger

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latencyTime := time.Now().Sub(start)
		webLogger.Infof("| %3d | %13v | %15s | %-7s %s",
			c.Writer.Status(),
			latencyTime,
			c.ClientIP(),
			c.Request.Method,
			c.Request.URL.Path,
		)
	}
}

// GinRecovery recover掉项目可能出现的panic
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					webLogger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					webLogger.Errorf("[Recovery from panic] error: %s, request: %s, stack:\n%s",
						err,
						string(httpRequest),
						string(debug.Stack()),
					)
				} else {
					webLogger.Errorf("[Recovery from panic] error: %s, request: %s",
						err,
						string(httpRequest),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

func InitWebLogger(logConfig *config.LogConfig) error {
	err := util.CreateDir(logConfig.LogPath)
	if err != nil {
		return err
	}
	logName := fmt.Sprintf("%s/%s.", logConfig.LogPath, logConfig.WebLogName)
	fileWriter, err := getFileWriter(logName)
	if err != nil {
		return err
	}
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(fileWriter, zapcore.AddSync(os.Stdout)), zap.DebugLevel)
	webLogger = zap.New(core).Sugar()
	return nil
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:    "msg",
		TimeKey:       "ts",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.000000"))
		},
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
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