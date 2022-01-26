package pkg

import (
	"strings"
	"text/template"
)

//下划线字符串 转换成 大驼峰 格式
func convertUpperCamelCase(s string) string {
	var upperStr string
	charsSplitSlice := strings.Split(s, "_")
	for _, chars := range charsSplitSlice {
		words := []rune(chars)
		for i := 0; i < len(words); i++ {
			if i == 0 {
				//如果是 a - z 之间
				if words[i] >= 97 && words[i] <= 122 {
					words[i] -= 32 // -32 转成大写字母
				}
				upperStr += string(words[i])
			} else {
				upperStr += string(words[i])
			}
		}
	}
	return upperStr
}

func GetGoTypeBySqlType(sqlType string) string {
	var mysqlTypeGoTypeMap = map[string]string{
		"tinyint":    "int32",
		"smallint":   "int32",
		"mediumint":  "int32",
		"int":        "int32",
		"integer":    "int64",
		"bigint":     "int64",
		"float":      "float64",
		"double":     "float64",
		"decimal":    "float64",
		"date":       "string",
		"time":       "string",
		"year":       "string",
		"enum":       "string",
		"datetime":   "time.Time",
		"timestamp":  "time.Time",
		"char":       "string",
		"varchar":    "string",
		"tinyblob":   "string",
		"tinytext":   "string",
		"blob":       "string",
		"text":       "string",
		"mediumblob": "string",
		"mediumtext": "string",
		"longblob":   "string",
		"longtext":   "string",
	}
	return mysqlTypeGoTypeMap[sqlType]
}

func newModelParse(templateStr string) *template.Template {
	tpl, err := template.New("model_template").
		Funcs(template.FuncMap{"convertUpperCamelCase": convertUpperCamelCase,"GetGoTypeBySqlType":GetGoTypeBySqlType}).
		Parse(templateStr)
	if err != nil {
		panic(err)
	}
	return tpl
}

var ModelTemplate = newModelParse(`
///////////////////////////////
// THE FILE IS AUTO CREATED //
//////////////////////////////
{{$quote := .Quote}}
package {{.PkgName}}

import (
	"gorm.io/gorm"
	"time"
)

// {{.StructName}} {{.TableComment.String}} 
type {{.StructName}} struct {
{{- range .Fields -}}
    {{- $tag := "" -}}
   	{{- if eq .ColumnKey.String "PRI" -}}
		{{- $tag = "gorm:\"primaryKey\"" | printf "%s" -}}
	{{end}}
 	{{if eq .DataType "time.Time"}} 
		{{- $tag =  " time" | printf "%s" -}}
   	{{end}}
	{{- $tag = " json:\"" | printf "%s%s" $tag -}}
	{{- $tag =  .ColumnName | printf "%s%s" $tag -}}
	{{- $tag =  "\"" | printf "%s%s"  $tag -}}
   	{{- .ColumnName | convertUpperCamelCase}} {{.DataType | GetGoTypeBySqlType}} {{$quote}} {{$tag}} {{$quote}} //{{.ColumnComment.String -}}
{{- end}}
}

func (t {{.StructName}}) TableName() string {
	 return "{{.TableName}}"
}

func (t {{.StructName}}) DbName() string {
	 return "{{.DbName}}"
}

//GetPrimaryKeyField 返回主键ID是哪个字段
func (t {{.StructName}}) GetPrimaryKeyField() string {
	return "Id"
}

//GetIsDelField 返回删除状态是哪个字段
func (t {{.StructName}}) GetIsDelField() string {
	return "deleted"
}

//GetDeleteTimeFiled 返回删除时间是哪个字段
func (t {{.StructName}}) GetDeleteTimeFiled() string {
	return "DeletedTime"
}

//BeforeCreate 创建记录时自动维护 CreatedTime UpdatedTime 两个字段, 这两个字段名 根据自己的表来设置
func (t {{.StructName}}) BeforeCreate(tx *gorm.DB) error {
	tx.Statement.SetColumn("CreatedTime", time.Now().Unix())
	tx.Statement.SetColumn("UpdatedTime", time.Now().Unix())

	return nil
}

//BeforeUpdate 更新记录时自动维护 UpdatedTime 字段 这个字段名 根据自己的表来设置
func (t {{.StructName}}) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("UpdatedTime", time.Now().Unix())

	return nil
}
`)
