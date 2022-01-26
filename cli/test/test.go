package test

import (
	"gin-self/core"
	"gin-self/extend/self_db"
	"gin-self/extend/utils/debug"
	"gin-self/model/mysql/test1_model"
)

type testStrut struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

func Test(c *core.CliContext)  {
	query := test1_model.NewQueryBuilder()
	query.Select("id", "name","goods_id").WhereOp("id", self_db.Equal, 114).QueryOne()

	var t  = testStrut{}
	c.BindParam(&t)
	debug.VarDump(t)
}