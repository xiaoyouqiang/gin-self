package demo

import (
	"fmt"
	"gin-self/app/demo/request_validator"
	"gin-self/extend/self_db"
	"gin-self/extend/self_redis"
	"gin-self/extend/utils/debug"
	"gin-self/extend/utils/e"
	"gin-self/extend/utils/helpers"
	"gin-self/extend/utils/request"
	"gin-self/model/mysql/test1_model"
	"gin-self/services/test_service"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context/ctxhttp"
	"net/http"
)

func Index(c *gin.Context) {
	id, _ := test_service.GetInstance(c).Create("xiao")

	helpers.ApiSuccess(c, id)

	//helpers.ApiError(c, e.AUTH_FAIL)
}

func Index1(c *gin.Context) {

	//query := test1_model.NewQueryBuilder()
	//query.
	//	WhereOp("id", self_db.GreaterThan, 0).
	//	Updates(gin.H{"GoodsId":"1"})
	query := test1_model.NewQueryBuilder()
	query.
		WhereOp("id", self_db.GreaterThan, 0).
		SoftDelete()
	//data := []test_model.Test{
	//	{
	//		Name: "xx",
	//	},
	//	{
	//		Name: "yy",
	//	},
	//}
	//
	//query := test_model.NewQueryBuilder()
	//
	//result,_:=query.CreateAll(data)

	//query := test1_model.NewQueryBuilder()
	//
	//result, _ := query.Select("id", "name","goods_id").WhereOp("id", self_db.Equal, 114).QueryOne()
	//
	//helpers.ApiSuccess(c, result)
}

func Index2(c *gin.Context) {
	var (
		req request_validator.Register
	)

	err := request.ParseRequest(c, &req)
	if  err != nil {
		helpers.ApiError(c, e.ParamError,err.Error())
		return
	}

	helpers.ApiSuccess(c, req)
}

func RedisTest(c *gin.Context) {
	self_redis.GetConn("master").Get("test")

	helpers.ApiSuccess(c, gin.H{})
}

func HttpTest(c *gin.Context) {
	_,err := ctxhttp.Get(c.Request.Context(),http.DefaultClient,"https://www.google.com/")
	if err != nil {
		fmt.Printf("%v", err)
	}

	query := test1_model.NewQueryBuilder()

	result, _ := query.Select("id", "name","goods_id").WhereOp("id", self_db.Equal, 114).QueryOne()

	debug.VarDump(result)

	helpers.ApiSuccess(c, gin.H{})
}
