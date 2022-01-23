package self_loger

import (
	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
	"sync"
	"sync/atomic"
)
var (
	Count uint32

	TracePool = &sync.Pool {
		New: func() interface{} {
			atomic.AddUint32(&Count,1)
			return new(TraceData)
		},
	}
)

type TraceData struct {
	traceId			   	string
	sqlList            	[]*sqlLog
	redisList          	[]*redisLog
	errorStack		 	[]string
}

type sqlLog struct {
	timestamp    string   `json "t"` // 时间，格式：2006-01-02 15:04:05
	stack        string   // 文件地址和行号
	sql          string   // SQL 语句
	rowsAffected int64    // 影响行数
	exeSeconds   float64   // 执行时长(单位秒)
	errorMsg 	 string
}

type redisLog struct {
	timestamp   string 	    // 时间，格式：2006-01-02 15:04:05
	operation   string   // 操作，SET/GET 等
	key         string   // Key
	value       string  // Value
	ttl         float64 // 超时时长(单位分)
	exeSeconds float64   // 执行时间(单位秒)
}

func NewTrace() *TraceData {
	traceData := TracePool.Get().(*TraceData)
	traceData.traceId = GenUUID()
	return traceData
}

func ReleaseTrace(traceData *TraceData) {
	if traceData == nil {
		println("tracedata is nil")
	} else {
		traceData.traceId = ""
		traceData.sqlList = nil
		traceData.redisList = nil
		TracePool.Put(traceData)
	}
}

func GetTraceByCtx(c *gin.Context) *TraceData {
	return c.Request.Context().Value("trace").(*TraceData)
}

func (t *TraceData) ValueTraceId() string {
	return t.traceId
}

func (t *TraceData) AddSqlLog(timestamp,stack,sql,errorMsg string, rowsAffected int64, exeSeconds float64)  {
	log := new(sqlLog)
	log.timestamp = timestamp
	log.stack = stack
	log.sql = sql
	log.rowsAffected = rowsAffected
	log.exeSeconds = exeSeconds
	log.errorMsg = errorMsg

	t.sqlList = append(t.sqlList, log)
}

func (t *TraceData) GetSqlLog() []*sqlLog{
	return t.sqlList
}

func (t *TraceData) AddRedisLog(timestamp,operation,key,value string, ttl, exeSeconds float64)  {

	log := new(redisLog)
	log.timestamp = timestamp
	log.operation = operation
	log.key = key
	log.value = value
	log.ttl = ttl
	log.exeSeconds = exeSeconds

	t.redisList = append(t.redisList, log)
}

func (t *TraceData) GetRedisLog() []*redisLog{
	return t.redisList
}

//AddErrorStackLog 异常日志
func (t *TraceData) AddErrorStackLog(stackLog []string) {
	t.errorStack  = stackLog
}
func (t *TraceData) GetErrorStackLog() []string{
	return t.errorStack
}


func (t *TraceData) ValueSqlInfo() []logrus.Fields {
	sqlLog := t.GetSqlLog()
	var logFieldsList []logrus.Fields

	for _,log := range sqlLog {
		logFieldsList = append(logFieldsList, logrus.Fields {
			"timestamp" : log.timestamp,
			"stack" : log.stack,
			"sql" : log.sql,
			"rowsAffected" : log.rowsAffected,
			"exeSeconds" : log.exeSeconds,
			"errorMsg" : log.errorMsg,
		})
	}

	return logFieldsList
}

func (t *TraceData) ValueRedisInfo() []logrus.Fields {
	redisLog := t.GetRedisLog()
	var logFieldsList []logrus.Fields

	for _,log := range redisLog {
		logFieldsList = append(logFieldsList, logrus.Fields {
			"timestamp" : log.timestamp,
			"operation" : log.operation,
			"key" : log.key,
			"value" : log.value,
			"ttl" : log.ttl,
			"exeSeconds" : log.exeSeconds,
		})
	}

	return logFieldsList
}

func (t *TraceData) ValueErrorInfo() []string {

	return t.GetErrorStackLog()
}

func GenUUID() string {
	id := ksuid.New()
	return id.String()
}