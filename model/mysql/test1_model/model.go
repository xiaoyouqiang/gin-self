
///////////////////////////////
// THE FILE IS AUTO CREATED //
//////////////////////////////

package test1_model

import (
	"gorm.io/gorm"
	"time"
)

// Test1  
type Test1 struct {
 	Id int32 ` gorm:"primaryKey" json:"id" ` //
 	Name string `  json:"name" ` //
 	GoodsId int32 `  json:"goods_id" ` //
 	CreatedTime int32 `  json:"created_time" ` //
 	UpdatedTime int32 `  json:"updated_time" ` //
 	DeletedTime int32 `  json:"deleted_time" ` //
}

func (t Test1) TableName() string {
	 return "test1"
}

func (t Test1) DbName() string {
	 return "test"
}

//GetPrimaryKeyField 返回主键ID是哪个字段
func (t Test1) GetPrimaryKeyField() string {
	return "Id"
}

//GetIsDelField 返回删除状态是哪个字段
func (t Test1) GetIsDelField() string {
	return "deleted"
}

//GetDeleteTimeFiled 返回删除时间是哪个字段
func (t Test1) GetDeleteTimeFiled() string {
	return "DeletedTime"
}

//BeforeCreate 创建记录时自动维护 CreatedTime UpdatedTime 两个字段, 这两个字段名 根据自己的表来设置
func (t Test1) BeforeCreate(tx *gorm.DB) error {
	tx.Statement.SetColumn("CreatedTime", time.Now().Unix())
	tx.Statement.SetColumn("UpdatedTime", time.Now().Unix())

	return nil
}

//BeforeUpdate 更新记录时自动维护 UpdatedTime 字段 这个字段名 根据自己的表来设置
func (t Test1) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("UpdatedTime", time.Now().Unix())

	return nil
}
