package logger

import (
	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net"
	"net/http/httputil"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
	"time"
	"used-car-deal-gobackend/base/config"
)

// Logger 封装zap 让项目中使用logger与zap松耦合
type Logger interface {
	// Debug uses fmt.Sprint to construct and log a message.
	Debug(args ...interface{})
	// Info uses fmt.Sprint to construct and log a message.
	Info(args ...interface{})
	// Warn uses fmt.Sprint to construct and log a message.
	Warn(args ...interface{})
	// Error uses fmt.Sprint to construct and log a message.
	Error(args ...interface{})
	// DPanic uses fmt.Sprint to construct and log a message. In development, the
	// logger then panics. (See DPanicLevel for details.)
	DPanic(args ...interface{})
	// Panic uses fmt.Sprint to construct and log a message, then panics.
	Panic(args ...interface{})
	// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
	Fatal(args ...interface{})
	// Debugf uses fmt.Sprintf to log a templated message.
	Debugf(template string, args ...interface{})
	// Infof uses fmt.Sprintf to log a templated message.
	Infof(template string, args ...interface{})

	// Warnf uses fmt.Sprintf to log a templated message.
	Warnf(template string, args ...interface{})
	// Errorf uses fmt.Sprintf to log a templated message.
	Errorf(template string, args ...interface{})

	// DPanicf uses fmt.Sprintf to log a templated message. In development, the
	// logger then panics. (See DPanicLevel for details.)
	DPanicf(template string, args ...interface{})

	// Panicf uses fmt.Sprintf to log a templated message, then panics.
	Panicf(template string, args ...interface{})

	// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
	Fatalf(template string, args ...interface{})

	// Debugw logs a message with some additional context. The variadic key-value
	// pairs are treated as they are in With.
	//
	// When debug-level logging is disabled, this is much faster than
	//  s.With(keysAndValues).Debug(msg)
	Debugw(msg string, keysAndValues ...interface{})

	// Infow logs a message with some additional context. The variadic key-value
	// pairs are treated as they are in With.
	Infow(msg string, keysAndValues ...interface{})

	// Warnw logs a message with some additional context. The variadic key-value
	// pairs are treated as they are in With.
	Warnw(msg string, keysAndValues ...interface{})

	// Errorw logs a message with some additional context. The variadic key-value
	// pairs are treated as they are in With.
	Errorw(msg string, keysAndValues ...interface{})

	// DPanicw logs a message with some additional context. In development, the
	// logger then panics. (See DPanicLevel for details.) The variadic key-value
	// pairs are treated as they are in With.
	DPanicw(msg string, keysAndValues ...interface{})
	// Panicw logs a message with some additional context, then panics. The
	// variadic key-value pairs are treated as they are in With.
	Panicw(msg string, keysAndValues ...interface{})

	// Fatalw logs a message with some additional context, then calls os.Exit. The
	// variadic key-value pairs are treated as they are in With.
	Fatalw(msg string, keysAndValues ...interface{})
}

var logger = struct {
	*zap.Logger
}{}
var sugaredLogger = struct {
	*zap.SugaredLogger
}{}
var mutex sync.Mutex

func GetInstance() Logger {
	return &sugaredLogger
}

// InitLogger 初始化Logger
func InitLogger() (err error) {
	var lvl = new(zapcore.Level)
	err = lvl.UnmarshalText([]byte(config.LogConfig.Level))
	if err != nil {
		return
	}

	var allCore []zapcore.Core

	if config.LogConfig.ConsoleOnly {
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		allCore = append(allCore, zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), lvl))
	} else {
		writeSyncer := getLogWriter()
		encoder := getEncoder()
		allCore = append(allCore, zapcore.NewCore(encoder, writeSyncer, lvl))
	}
	core := zapcore.NewTee(allCore...)
	logger.Logger = zap.New(core, zap.AddCaller())
	sugaredLogger.SugaredLogger = logger.Sugar()
	zap.ReplaceGlobals(logger.Logger) // 替换zap包中全局的logger实例，后续在其他包中只需使用zap.L()调用即可
	return
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	if config.LogConfig.ConsoleOnly {
		return zapcore.NewConsoleEncoder(encoderConfig)
	}
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {

	lumberJackLogger := &lumberjack.Logger{
		Filename:   config.LogConfig.Filename,   //日志文件位置
		MaxSize:    config.LogConfig.MaxSize,    //进行切割之前，日志文件最大值（单位：MB),默认100MB
		MaxBackups: config.LogConfig.MaxBackups, //保留旧文件的最大个数
		MaxAge:     config.LogConfig.MaxAge,     //保留旧文伯最大天数
		Compress:   false,                       //是否压缩/归档旧文件
	}
	return zapcore.AddSync(lumberJackLogger)
}

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		logger.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			err := recover()
			if err != nil {
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
					if config.LogConfig.ConsoleOnly {
						sugaredLogger.Errorf("%v\nerror:%v\nrequest:%v", c.Request.URL.Path,
							err, string(httpRequest))
					} else {
						logger.Error(c.Request.URL.Path,
							zap.Any("error", err),
							zap.String("request", string(httpRequest)),
						)
					}

					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					if config.LogConfig.ConsoleOnly {
						runtime.Caller(1)
						sugaredLogger.Errorf("[Recovery from panic] err:%v \n",
							err)
					} else {
						logger.Error("[Recovery from panic]",
							zap.Any("error", err),
							zap.String("request", string(httpRequest)),
							zap.String("stack", string(debug.Stack())),
						)
					}

				} else {
					if config.LogConfig.ConsoleOnly {
						sugaredLogger.Errorf("[Recovery from panic] error:%v \nrequest:%v",
							err, string(httpRequest))
					} else {
						logger.Error("[Recovery from panic]",
							zap.Any("error", err),
							zap.String("request", string(httpRequest)),
						)
					}

				}
				//c.AbortWithStatus(http.StatusInternalServerError)
				c.Abort()
			}
		}()
		c.Next()
	}
}
