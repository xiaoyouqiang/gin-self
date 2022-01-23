package self_loger

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"path"
	"sync"
	"time"
)

type LogType string

var (
	CliType = LogType("cli")
	HttpType = LogType("http")

	once sync.Once
	loggerInstance *logger
)

type logger struct {
	logRusInstance *logrus.Logger
	traceData *TraceData
	tracePool sync.Pool
}

//GetInstance 获取实例
func GetInstance() *logger {
	once.Do(func() {
		loggerInstance = &logger{
			logRusInstance: logrus.New(),
		}
	})
	return loggerInstance
}


//Info 写入日志到文件
func (l *logger) Info(data logrus.Fields) {
	if l.logRusInstance == nil {
		panic("logger use before SetLogger for type")
	}
	l.logRusInstance.WithFields(data).Info()
}

//SetLogger 设置logger参数
//@param logType
func (l *logger) SetLogger(logType LogType)  {
	if logType == "" {
		log.Fatalln("logger not init for set app_type")
	}

	var (
		logFilePath = "./log" //文件存储路径
		logFileName = fmt.Sprintf("%v", logType) + ".log"
	)
	// 日志文件
	fileName := path.Join(logFilePath, logFileName)
	// 写入文件
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("打开/写入文件失败", err)
	}
	// 实例化
	l.logRusInstance = logrus.New()
	// 日志级别
	l.logRusInstance.SetLevel(logrus.DebugLevel)
	// 设置输出
	l.logRusInstance.Out = file
	// 设置 rotatelogs,实现文件分割
	logWriter, err := rotatelogs.New(
		// 分割后的文件名称
		fileName+".%Y%m%d.log",
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
	l.logRusInstance.AddHook(lfshook.NewHook(writerMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	}))
}