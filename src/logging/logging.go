package logging

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-kit/kit/transport"
	kitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/go-kit/log/term"
	"github.com/lestrrat-go/file-rotatelogs"
	gormlogger "gorm.io/gorm/logger"

	"github.com/IceBearAI/aigc/src/encode"
	"github.com/IceBearAI/aigc/src/services"
)

const (
	TraceId = "traceId"
)

// LogErrorHandler is a transport error handler implementation which logs an error.
type LogErrorHandler struct {
	logger kitlog.Logger
	apiSvc services.Service
}

func (l *LogErrorHandler) Handle(ctx context.Context, err error) {
	var errDefined bool
	for k := range encode.ResponseMessage {
		if strings.Contains(err.Error(), k.Error().Error()) {
			errDefined = true
			break
		}
	}

	defer func() {
		_ = l.logger.Log("traceId", ctx.Value(TraceId), "logErrorHandler", "Handle", "err", err.Error())
	}()

	if !errDefined {
		hostname, _ := os.Hostname()
		content := fmt.Sprintf("\nMessage: 未定义错误! \nError: %s \nHostname: %s",
			err.Error(),
			hostname,
		)
		//_ = l.apiSvc.Alarm().Push(ctx, "未定义错误!", content, "err_notfound", alarm.LevelInfo, 1)
		log.Println(fmt.Sprintf("host: %s, err: %s", hostname, content))
	}
}

func NewLogErrorHandler(logger kitlog.Logger, apiSvc services.Service) transport.ErrorHandler {
	return &LogErrorHandler{
		logger: logger,
		apiSvc: apiSvc,
	}
}

func SetLogging(logger kitlog.Logger, logPath, logFileName, legLevel, serviceName, logDrive string) kitlog.Logger {
	if !strings.EqualFold(logDrive, "") && strings.EqualFold("syslog", strings.ToLower(logDrive)) {
		// TODO: windows 不支持syslog 暂时先去掉
		//return kitlog.With(syslogLogger(legLevel, serviceName), "caller", kitlog.DefaultCaller)
		return term.NewLogger(os.Stdout, kitlog.NewLogfmtLogger, colorFunc())
	}

	if !strings.EqualFold(logPath, "") {
		// default log
		logger = kitlog.WithPrefix(defaultLogger(fmt.Sprintf("%s/%s", logPath, logFileName)), "ts", kitlog.TimestampFormat(func() time.Time {
			return time.Now()
		}, "2006-01-02 15:04:05"))
	} else {
		//logger = kitlog.NewLogfmtLogger(kitlog.StdlibWriter{})
		logger = term.NewLogger(os.Stdout, kitlog.NewLogfmtLogger, colorFunc())
		logger = kitlog.WithPrefix(logger, "ts", kitlog.TimestampFormat(func() time.Time {
			return time.Now()
		}, "2006-01-02 15:04:05"))
	}
	return kitlog.With(level.NewFilter(logger, logLevel(legLevel)), "caller", kitlog.DefaultCaller)
}

func logLevel(logLevel string) (opt level.Option) {
	switch logLevel {
	case "warn":
		opt = level.AllowWarn()
	case "error":
		opt = level.AllowError()
	case "debug":
		opt = level.AllowDebug()
	case "info":
		opt = level.AllowInfo()
	case "all":
		opt = level.AllowAll()
	default:
		opt = level.AllowNone()
	}

	return
}

func defaultLogger(filePath string) kitlog.Logger {
	linkFile, err := filepath.Abs(filePath)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	writer, err := rotatelogs.New(
		linkFile+"-%Y-%m-%d",
		rotatelogs.WithLinkName(linkFile),         // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(time.Hour*24*365),   // 文件最大保存时间
		rotatelogs.WithRotationTime(time.Hour*24), // 日志切割时间间隔
	)

	if err != nil {
		log.Fatal(err)
		return nil
	}

	return kitlog.NewLogfmtLogger(writer)
}

func colorFunc() func(keyvals ...interface{}) term.FgBgColor {
	return func(keyvals ...interface{}) term.FgBgColor {
		for i := 0; i < len(keyvals)-1; i += 2 {
			if keyvals[i] != "level" {
				continue
			}
			val := fmt.Sprintf("%v", keyvals[i+1])
			switch val {
			case "debug":
				return term.FgBgColor{Fg: term.DarkGray}
			case "info":
				return term.FgBgColor{Fg: term.Blue}
			case "warn":
				return term.FgBgColor{Fg: term.Yellow}
			case "error":
				return term.FgBgColor{Fg: term.Red}
			case "crit":
				return term.FgBgColor{Fg: term.Gray, Bg: term.DarkRed}
			default:
				return term.FgBgColor{}
			}
		}
		return term.FgBgColor{}
	}
}

type gormLogger struct {
	logger kitlog.Logger
	level  string
}

func (g *gormLogger) LogMode(l gormlogger.LogLevel) gormlogger.Interface {
	panic("implement me")
}

func (g gormLogger) Info(ctx context.Context, s string, i ...interface{}) {
	if g.level == "info" {
		fmt.Println(s, i)
		_ = level.Info(g.logger).Log(s, i)
		//l.Printf(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

func (g gormLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	panic("implement me")
}

func (g gormLogger) Error(ctx context.Context, s string, i ...interface{}) {
	panic("implement me")
}

func (g gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	panic("implement me")
}

func NewGromLogger(logger kitlog.Logger, lv string) gormlogger.Interface {
	return &gormLogger{logger, lv}
}

type gormWriter struct {
	logger kitlog.Logger
	level  string
}

func (g *gormWriter) Printf(s string, i ...interface{}) {
	logger := level.Debug(g.logger)
	if strings.Contains(s, "[error]") {
		logger = level.Error(g.logger)
	} else if strings.Contains(s, "[warn]") {
		logger = level.Warn(g.logger)
	} else if strings.Contains(s, "[info]") {
		logger = level.Warn(g.logger)
	}
	if len(i)%2 == 0 {
		_ = logger.Log("file", i[0], "took", i[1], "rows", i[2], "sql", i[3])
	} else {
		_ = logger.Log("file", i[0], "slow", i[2], "took", i[2], "rows", i[3], "sql", i[4])
	}
}

func NewGormWriter(logger kitlog.Logger) gormlogger.Writer {
	logger = kitlog.With(logger, "gorm", "log")
	return &gormWriter{logger: logger}
}

func NewGormLogging(logger kitlog.Logger) gormlogger.Interface {
	return gormlogger.New(NewGormWriter(logger), gormlogger.Config{
		SlowThreshold:             500 * time.Millisecond,
		Colorful:                  false,
		IgnoreRecordNotFoundError: false,
		LogLevel:                  gormlogger.Info,
	})
}
