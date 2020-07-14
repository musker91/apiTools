package logger

import (
	"apiTools/libs/config"
	"apiTools/utils"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

var (
	Echo *logrus.Logger
)

func InitLogger() (err error) {
	Echo = logrus.New()

	// 记录文件名和行号
	//Echo.SetReportCaller(true)

	// 设置日志输入格式
	switch config.GetString("web::logOutType") {
	case "json":
		Echo.SetFormatter(&logrus.JSONFormatter{})
	default:
		Echo.SetFormatter(&logrus.TextFormatter{})
	}

	// 设置开发模式下的日志输出
	if config.GetString("web::appMode") == "development" {
		// 设置日志输出级别
		Echo.SetLevel(logrus.DebugLevel)
		// 设置日志输出的路径
		Echo.Out = os.Stdout
		return
	}

	// 设置日志输出级别
	switch config.GetString("web::logLevel") {
	case "info":
		Echo.SetLevel(logrus.InfoLevel)
	case "warn":
		Echo.SetLevel(logrus.WarnLevel)
	case "debug":
		Echo.SetLevel(logrus.DebugLevel)
	case "error":
		Echo.SetLevel(logrus.ErrorLevel)
	case "panic":
		Echo.SetLevel(logrus.PanicLevel)
	case "fatal":
		Echo.SetLevel(logrus.FatalLevel)
	default:
		Echo.SetLevel(logrus.DebugLevel)
	}

	// 设置日志输出位置
	switch config.GetString("web::logOutPath") {
	case "file":
		// 日志打印到指定的目录
		logFileName := path.Join(utils.GetRootPath(), "logs", "apitools.log")
		//禁止logrus的输出
		logOut, err := os.OpenFile(logFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend)
		if err != nil {
			break
		}
		// 设置日志输出的路径
		Echo.Out = logOut
		// 创建日志输出对象
		var logSaveDay = time.Duration(config.GetInt("web::logSaveDay"))
		var logSplitTime = time.Duration(config.GetInt("web::logSplitTime"))
		logWriter, err := rotatelogs.New(
			logFileName+".%Y-%m-%d-%H-%M.log",                   // 日志切割名称
			rotatelogs.WithLinkName(logFileName),                // 生成软链，指向最新日志文件
			rotatelogs.WithMaxAge(logSaveDay*24*time.Hour),      // 文件最大保存时间
			rotatelogs.WithRotationTime(logSplitTime*time.Hour), // 日志切割时间间隔
		)
		// 为不同级别设置不同的输出目的
		writeMap := lfshook.WriterMap{
			logrus.InfoLevel:  logWriter,
			logrus.FatalLevel: logWriter,
			logrus.DebugLevel: logWriter,
			logrus.WarnLevel:  logWriter,
			logrus.ErrorLevel: logWriter,
			logrus.PanicLevel: logWriter,
		}

		// 创建logrus的本地文件系统钩子
		lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{})
		Echo.AddHook(lfHook)
	default:
		Echo.Out = os.Stdout
	}
	return
}
