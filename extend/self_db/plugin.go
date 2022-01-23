package self_db

import (
	"time"

	"gin-self/extend/self_loger"

	"gorm.io/gorm"
	"gorm.io/gorm/utils"
)

const (
	callBackBeforeName = "core:before"
	callBackAfterName  = "core:after"
	startTime          = "_start_time"
)

type LoggerPlugin struct{}

func (op *LoggerPlugin) Name() string {
	return "loggerPlugin"
}

func (op *LoggerPlugin) Initialize(db *gorm.DB) (err error) {
	// 开始前
	_ = db.Callback().Create().Before("gorm:before_create").Register(callBackBeforeName, before)
	_ = db.Callback().Query().Before("gorm:query").Register(callBackBeforeName, before)
	_ = db.Callback().Delete().Before("gorm:before_delete").Register(callBackBeforeName, before)
	_ = db.Callback().Update().Before("gorm:setup_reflect_value").Register(callBackBeforeName, before)
	_ = db.Callback().Row().Before("gorm:row").Register(callBackBeforeName, before)
	_ = db.Callback().Raw().Before("gorm:raw").Register(callBackBeforeName, before)

	// 结束后
	_ = db.Callback().Create().After("gorm:after_create").Register(callBackAfterName, after)
	_ = db.Callback().Query().After("gorm:after_query").Register(callBackAfterName, after)
	_ = db.Callback().Delete().After("gorm:after_delete").Register(callBackAfterName, after)
	_ = db.Callback().Update().After("gorm:after_update").Register(callBackAfterName, after)
	_ = db.Callback().Row().After("gorm:row").Register(callBackAfterName, after)
	_ = db.Callback().Raw().After("gorm:raw").Register(callBackAfterName, after)
	return
}

var _ gorm.Plugin = &LoggerPlugin{}

func before(db *gorm.DB) {
	db.InstanceSet(startTime, time.Now())
	return
}

func after(db *gorm.DB) {
	_ts, isExist := db.InstanceGet(startTime)
	if !isExist {
		return
	}

	ctx := db.Statement.Context
	if ctx.Value("trace") == nil {
		return
	}

	ts, ok := _ts.(time.Time)
	if !ok {
		return
	}

	errorStr := ""
	if db.Error != nil && db.Error != gorm.ErrRecordNotFound {
		errorStr = db.Error.Error()
	}

	sql := db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)

	ctx.Value("trace").(*self_loger.TraceData).AddSqlLog(
		time.Now().Format("2006/01/02 15:04:05"),
		utils.FileWithLineNum(),
		sql,
		errorStr,
		db.Statement.RowsAffected,
		time.Since(ts).Seconds(),
	)

	return
}
