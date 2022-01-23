package test

import (
	"gin-self/core"
	"gin-self/extend/self_db"
	"gin-self/model/mysql/test1_model"
)

type testStrut struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

func Test(c *core.CliContext)  {
	//query := test_model.NewQueryBuilder()
	//result, _ := query.Select("id", "name", "created_time").WhereOp("id", self_db.Equal, 35).QueryOne()
	//debug.VarDump(result)
	//v,_ := c.GetParam("name")
	//debug.VarDump(v)
	query := test1_model.NewQueryBuilder()
	//query.Select("id", "name", "created_time").WhereOp("id", self_db.Equal, 1).QueryOne()

	 query.Select("id", "name","goods_id").WhereOp("id", self_db.Equal, 114).QueryOne()

	//debug.VarDump(result)
	var t  = testStrut{}
	c.BindParam(&t)
	//debug.VarDump(t)
}