
///////////////////////////////
// THE FILE IS AUTO CREATED //
//////////////////////////////

package test1_model

import (
	"fmt"
	"sync"
	"time"

	"gin-self/extend/self_db"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var modelOnce sync.Once
var modelInstance *Test1

type Test1ModelQueryBuilder struct {
	order []string
	where []struct {
		prefix string
		value  interface{}
	}
	fields struct {
		query interface{}
		args  []interface{}
	}
	limit  int
	offset int
	db     *gorm.DB
}

func NewModel() *Test1 {
	modelOnce.Do(func() {
		modelInstance = new(Test1)
	})
	return modelInstance
}

func NewQueryBuilder() *Test1ModelQueryBuilder {
	return &Test1ModelQueryBuilder{
		db: self_db.GetDbConn(NewModel().DbName()),
	}
}

func (qb *Test1ModelQueryBuilder) getDbConn() *gorm.DB {
	return qb.db
}

func (qb *Test1ModelQueryBuilder) Create(model Test1) (id int32, err error) {
	if err = qb.db.Create(&model).Error; err != nil {
		return 0, errors.Wrap(err, "create err")
	}
	return model.Id, nil
}


func (qb *Test1ModelQueryBuilder) CreateAll(model []Test1) (result []Test1, err error) {
	if err = qb.db.Create(&model).Error; err != nil {
		return model, errors.Wrap(err, "create err")
	}
	return model, nil
}

func (qb *Test1ModelQueryBuilder) buildQuery() *gorm.DB {
	ret := qb.db
	for _, where := range qb.where {
		ret = ret.Where(where.prefix, where.value)
	}
	for _, order := range qb.order {
		ret = ret.Order(order)
	}
	if qb.fields.query != nil {
		ret = ret.Select(qb.fields.query, qb.fields.args...)
	}

	ret = ret.Limit(qb.limit).Offset(qb.offset)
	return ret
}

func (qb *Test1ModelQueryBuilder) Updates(m map[string]interface{}) (err error) {
	qb.db = qb.db.Model(&Test1{})

	for _, where := range qb.where {
		qb.db.Where(where.prefix, where.value)
	}

	if err = qb.db.Updates(m).Error; err != nil {
		return errors.Wrap(err, "updates err")
	}
	return nil
}

func (qb *Test1ModelQueryBuilder) Delete() (err error) {
	for _, where := range qb.where {
		qb.db = qb.db.Where(where.prefix, where.value)
	}

	if err = qb.db.Delete(&Test1{}).Error; err != nil {
		return errors.Wrap(err, "delete err")
	}
	return nil
}

//SoftDelete 软删除
func (qb *Test1ModelQueryBuilder) SoftDelete() (err error) {
	updateDate := self_db.M{}
	//模型设置了软删除字段
	if NewModel().GetIsDelField() != "" && NewModel().GetDeleteTimeFiled() != "" {
		updateDate[NewModel().GetIsDelField()] = 1
		updateDate[NewModel().GetDeleteTimeFiled()] = time.Now().Unix()
	}
	if len(updateDate) == 0 {
		return errors.Wrap(err, "not set soft delete fields")
	}

	return qb.Updates(updateDate)
}

func (qb *Test1ModelQueryBuilder) Count() (int64, error) {
	var c int64
	res := qb.buildQuery().Model(&Test1{}).Count(&c)
	if res.Error != nil && res.Error == gorm.ErrRecordNotFound {
		c = 0
	}
	return c, res.Error
}

func (qb *Test1ModelQueryBuilder) First() (Test1, error) {
	var ret Test1

	res := qb.buildQuery().First(&ret)

	return ret, res.Error
}

func (qb *Test1ModelQueryBuilder) QueryOne() (Test1, error) {
	var ret Test1

	err := qb.buildQuery().Take(&ret).Error

	return ret, err
}

func (qb *Test1ModelQueryBuilder) QueryAll() ([]Test1, error) {
	var ret []Test1
	err := qb.buildQuery().Find(&ret).Error
	return ret, err
}

func (qb *Test1ModelQueryBuilder) Limit(limit int) *Test1ModelQueryBuilder {
	qb.limit = limit
	return qb
}

func (qb *Test1ModelQueryBuilder) Offset(offset int) *Test1ModelQueryBuilder {
	qb.offset = offset
	return qb
}

func (qb *Test1ModelQueryBuilder) Select(query interface{}, args ...interface{}) *Test1ModelQueryBuilder {
	qb.fields.query = query
	qb.fields.args = args

	return qb
}

func (qb *Test1ModelQueryBuilder) WhereOp(field string, op self_db.SqlOp, value interface{}) *Test1ModelQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", field, op),
		value,
	})
	return qb
}

func (qb *Test1ModelQueryBuilder) WhereIn(field string, value interface{}) *Test1ModelQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", field, "IN"),
		value,
	})
	return qb
}

func (qb *Test1ModelQueryBuilder) WhereNotIn(field string, value interface{}) *Test1ModelQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", field, "NOT IN"),
		value,
	})
	return qb
}

func (qb *Test1ModelQueryBuilder) Order(field string, asc bool) *Test1ModelQueryBuilder {
	order := "DESC"
	if asc {
		order = "ASC"
	}

	qb.order = append(qb.order, field+" "+order)
	return qb
}

