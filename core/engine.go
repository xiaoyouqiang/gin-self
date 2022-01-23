package core

import (
	"context"
	"errors"
	"gin-self/extend/self_db"
	"gin-self/extend/self_loger"
	"gin-self/extend/self_redis"
	"github.com/sirupsen/logrus"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type CliContext struct {
	context.Context
	params interface{}
}

type CliEngine struct {
	handlers map[string]func(ctx *CliContext)
	args []string
	startTime time.Time
	ctx CliContext
}

func NewCLiEngine(args []string) *CliEngine {
	startTime := time.Now()
	return &CliEngine{
		handlers: make(map[string]func(ctx *CliContext)),
		args: args,
		startTime: startTime,
		ctx:CliContext {
			Context:context.Background(),
		},
	}
}

//SetHandler 设置 cli 命令路由
func (e *CliEngine) SetHandler(commandPath string, handler func(ctx *CliContext))  {
	if _,exists := e.handlers[commandPath];!exists {
		e.handlers[commandPath] = handler
	} else {
		log.Fatalln("Have a Repeat handler")
	}
}

//ExecBefore 执行 cli 前置
func (e *CliEngine) ExecBefore(trace *self_loger.TraceData)  {
	childCtx := context.WithValue(e.ctx,"trace", trace)

	self_db.WithContext(childCtx)

	self_redis.WithContext(childCtx)
}

//ExecHandler 执行 cli 注册的路由：解析命令行 输入 执行对应的 程序
func (e *CliEngine) ExecHandler()  {
	trace := self_loger.NewTrace()
	defer func() {
		self_loger.ReleaseTrace(trace)
	}()

	e.ExecBefore(trace)

	e.ctx.setParams(e.args)

	var cliGinContext = &e.ctx
	commandPath := e.args[0]
	if _,exists := e.handlers[commandPath];!exists {
		log.Fatalln("command handler not find")
	} else {
		e.handlers[commandPath](cliGinContext)
	}

	e.ExecEnd(trace)
}

//ExecEnd 执行 cli 命令 后置操作
func (e *CliEngine) ExecEnd(trace *self_loger.TraceData)  {
	//执行耗时
	endTime := time.Now()
	execTime := endTime.Sub(e.startTime).Seconds()

	//日志
	traceInfo := logrus.Fields{
		"trace_id":    trace.ValueTraceId(),
		"command":  e.args[0],
		"command_param":   e.args[1:],
		"exec_time":   execTime,
		"trace_info": logrus.Fields {
			"sql":           trace.ValueSqlInfo(),
			"redis":         trace.ValueRedisInfo(),
			"error_stack":    trace.ValueErrorInfo(),
		},
	}

	self_loger.GetInstance().Info(traceInfo)
}

//setParams 保存命令行参数
func (ctx *CliContext) setParams(args []string) {
	var params = map[string]interface{}{}
	argsNum := len(args)
	if argsNum >= 2 {
		for i:=1; i < argsNum; i++ {
			tmpArgs := strings.Split(args[i], "&")
			for _,v := range tmpArgs {
				item := strings.Split(v, "=")
				itemValue := ""
				if len(item) == 2 {
					itemValue = item[1]
				}
				params[item[0]] = itemValue
			}
		}
	}

	ctx.params = params
}

//GetParam 获取命令行参数：命令行程序中可获取参数
func (ctx *CliContext) GetParam(name string) (value interface{}, exists bool) {
	if v,ok := ctx.params.(map[string]interface{})[name];ok {
		exists = true
		value = v
		return
	}

	exists = false
	value = ""

	return
}

//BindParam 绑定命令行产生到 结构体 obj
func (ctx *CliContext) BindParam(obj interface{}) error {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	if t.Kind() != reflect.Ptr {
		return errors.New("参数不是引用类型")
	}

	targetType := t.Elem()
	targetValue := v.Elem()

	if targetType.Kind() != reflect.Struct {
		return errors.New("只接受结构体引用参数")
	}

	for i:=0; i < targetType.NumField();i++ {
		var fieldName string
		field := targetType.Field(i)
		fieldTag := field.Tag
		tagString := fieldTag.Get("json")
		if tagString == "" {
			fieldName = strings.ToLower(field.Name)
		} else {
			fieldName = strings.Split(tagString,",")[0]
		}

		//绑定参数到结构体
		if mapV,ok := ctx.params.(map[string]interface{})[fieldName];ok {
			//如果设置了该字段参数
			fieldV := targetValue.Field(i)
			fType := field.Type.Kind()
			if !fieldV.CanSet() {
				continue
			}

			switch fType {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
				reflect.Int64,reflect.Uint,reflect.Uint8,reflect.Uint16,reflect.Uint32,reflect.Uint64:
					tmp,_ := strconv.ParseInt(mapV.(string),10,64)
					fieldV.SetInt(tmp)
				case reflect.Float32, reflect.Float64:
					tmp,_ := strconv.ParseFloat(mapV.(string),64)
					fieldV.SetFloat(tmp)
				case reflect.String:
					tmp := mapV.(string)
					fieldV.SetString(tmp)
				case reflect.Bool:
					tmp,_ := strconv.ParseBool(mapV.(string))
					fieldV.SetBool(tmp)
				default:

			}
		}
	}

	return nil
}