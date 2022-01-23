package log_self

import (
	"log"
	"os"
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var (
	logFilePath = "../log_self" //文件存储路径
	logFileName = "system.log_self"
	Logger      *logrus.Logger
)

func init() {
	// 日志文件
	// 日志文件
	fileName := path.Join(logFilePath, logFileName)
	// 写入文件
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalln("日志文件打开失败")
	}
	// 实例化
	Logger = logrus.New()
	// 日志级别
	Logger.SetLevel(logrus.DebugLevel)
	// 设置输出
	Logger.Out = file
	// 设置 rotatelogs,实现文件分割
	logWriter, err := rotatelogs.New(
		// 分割后的文件名称
		fileName+".%Y%m%d.log_self",
		// 生成软链，指向最新日志文件
		rotatelogs.WithLinkName(fileName),
		// 设置最大保存时间(7天)
		rotatelogs.WithMaxAge(7*24*time.Hour), //以hour为单位的整数
		// 设置日志切割时间间隔(1天)
		rotatelogs.WithRotationTime(1*time.Hour),
	)
	// hook机制的设置
	writerMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}
	//给logrus添加hook
	Logger.AddHook(lfshook.NewHook(writerMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	}))
}
