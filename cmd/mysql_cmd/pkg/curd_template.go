package pkg

import "text/template"

func newCurdParse(templateStr string) *template.Template {
	tpl, err := template.New("curd_template").Parse(templateStr)
	if err != nil {
		panic(err)
	}
	return tpl
}

var CurdTemplate = newCurdParse(`
///////////////////////////////
// THE FILE IS AUTO CREATED //
//////////////////////////////

package {{.PkgName}}

import (
	"fmt"
	"sync"
	"time"

	"gin-self/extend/self_db"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var modelOnce sync.Once
var modelInstance *{{.StructName}}

type {{.QueryBuilderName}} struct {
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

func NewModel() *{{.StructName}} {
	modelOnce.Do(func() {
		modelInstance = new({{.StructName}})
	})
	return modelInstance
}

func NewQueryBuilder() *{{.QueryBuilderName}} {
	return &{{.QueryBuilderName}}{
		db: self_db.GetDbConn(NewModel().DbName()),
	}
}

func (qb *{{.QueryBuilderName}}) getDbConn() *gorm.DB {
	return qb.db
}

func (qb *{{.QueryBuilderName}}) Create(model {{.StructName}}) (id int32, err error) {
	if err = qb.db.Create(&model).Error; err != nil {
		return 0, errors.Wrap(err, "create err")
	}
	return model.{{.PkFieldName}}, nil
}


func (qb *{{.QueryBuilderName}}) CreateAll(model []{{.StructName}}) (result []{{.StructName}}, err error) {
	if err = qb.db.Create(&model).Error; err != nil {
		return model, errors.Wrap(err, "create err")
	}
	return model, nil
}

func (qb *{{.QueryBuilderName}}) buildQuery() *gorm.DB {
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

func (qb *{{.QueryBuilderName}}) Updates(m map[string]interface{}) (err error) {
	qb.db = qb.db.Model(&{{.StructName}}{})

	for _, where := range qb.where {
		qb.db.Where(where.prefix, where.value)
	}

	if err = qb.db.Updates(m).Error; err != nil {
		return errors.Wrap(err, "updates err")
	}
	return nil
}

func (qb *{{.QueryBuilderName}}) Delete() (err error) {
	for _, where := range qb.where {
		qb.db = qb.db.Where(where.prefix, where.value)
	}

	if err = qb.db.Delete(&{{.StructName}}{}).Error; err != nil {
		return errors.Wrap(err, "delete err")
	}
	return nil
}

//SoftDelete 软删除
func (qb *{{.QueryBuilderName}}) SoftDelete() (err error) {
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

func (qb *{{.QueryBuilderName}}) Count() (int64, error) {
	var c int64
	res := qb.buildQuery().Model(&{{.StructName}}{}).Count(&c)
	if res.Error != nil && res.Error == gorm.ErrRecordNotFound {
		c = 0
	}
	return c, res.Error
}

func (qb *{{.QueryBuilderName}}) First() ({{.StructName}}, error) {
	var ret {{.StructName}}

	res := qb.buildQuery().First(&ret)

	return ret, res.Error
}

func (qb *{{.QueryBuilderName}}) QueryOne() ({{.StructName}}, error) {
	var ret {{.StructName}}

	err := qb.buildQuery().Take(&ret).Error

	return ret, err
}

func (qb *{{.QueryBuilderName}}) QueryAll() ([]{{.StructName}}, error) {
	var ret []{{.StructName}}
	err := qb.buildQuery().Find(&ret).Error
	return ret, err
}

func (qb *{{.QueryBuilderName}}) Limit(limit int) *{{.QueryBuilderName}} {
	qb.limit = limit
	return qb
}

func (qb *{{.QueryBuilderName}}) Offset(offset int) *{{.QueryBuilderName}} {
	qb.offset = offset
	return qb
}

func (qb *{{.QueryBuilderName}}) Select(query interface{}, args ...interface{}) *{{.QueryBuilderName}} {
	qb.fields.query = query
	qb.fields.args = args

	return qb
}

func (qb *{{.QueryBuilderName}}) WhereOp(field string, op self_db.SqlOp, value interface{}) *{{.QueryBuilderName}} {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", field, op),
		value,
	})
	return qb
}

func (qb *{{.QueryBuilderName}}) WhereIn(field string, value interface{}) *{{.QueryBuilderName}} {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", field, "IN"),
		value,
	})
	return qb
}

func (qb *{{.QueryBuilderName}}) WhereNotIn(field string, value interface{}) *{{.QueryBuilderName}} {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", field, "NOT IN"),
		value,
	})
	return qb
}

func (qb *{{.QueryBuilderName}}) Order(field string, asc bool) *{{.QueryBuilderName}} {
	order := "DESC"
	if asc {
		order = "ASC"
	}

	qb.order = append(qb.order, field+" "+order)
	return qb
}

`)
