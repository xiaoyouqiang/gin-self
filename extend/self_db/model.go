package self_db

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"gin-self/extend/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var dbPool = make(map[string]*gorm.DB) //数据库连接对象数组 可存储多个 数据库连接实例

//GetDbConn 根据数据库名获取 数据库连接
//@param opts Option 设置参数函数
func GetDbConn(dbName string) *gorm.DB {
	if _, ok := dbPool[dbName]; !ok {
		panic("db not init,maybe no db config")
	}

	db := dbPool[dbName]

	//end

	return db
}

//WithContext gin 中间件中 追加上下文 到 gorm中 进行日志 记录
func WithContext(c context.Context) {
	for k, _ := range dbPool {
		dbPool[k] = dbPool[k].WithContext(c)
	}
}

func Open() {
	needInitDatabases := strings.TrimSpace(config.Get("app", "need_init_database").String())
	if needInitDatabases == "" {
		log.Fatalln("not init_databases config for app")
	}
	databaseArray := strings.Split(needInitDatabases,",")
	for _,dbNameItem := range databaseArray {
		dbName := strings.TrimSpace(dbNameItem)
		if dbName == "" {
			continue
		}
		if _, ok := dbPool[dbName]; ok {
			//该库已经初始化
			continue
		}
		connect(dbName)
	}
}

func connect(dbName string)  {
	dbConfigSection := "database-" + dbName
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Get(dbConfigSection, "user").String(),
		config.Get(dbConfigSection, "password").String(),
		config.Get(dbConfigSection, "host").String(),
		config.Get(dbConfigSection, "db_name").String(),
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名，启用该选项后，`User` 表将是`user`
		},
		SkipDefaultTransaction: true, //默认不开始事务 需要时 手动启动事务 以下有启动事务方法
	})

	db = db.Set("db_name", dbName)

	if err != nil {
		log.Fatalln(err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalln(err.Error())
	}

	// 设置连接池 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	sqlDB.SetMaxOpenConns(config.Get("database", "max_open_conn").MustInt(100))

	// 设置最大连接数 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	sqlDB.SetMaxIdleConns(config.Get("database", "max_idle_conn").MustInt(100))

	// 设置最大连接超时
	sqlDB.SetConnMaxLifetime(time.Duration(config.Get("database", "time_out").MustInt(5)) * time.Second)

	// 使用插件 统一追踪日志
	db.Use(&LoggerPlugin{})

	dbPool[dbName] = db
}
