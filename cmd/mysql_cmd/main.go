package main

import (
	"database/sql"
	"flag"
	"fmt"
	"gin-self/cmd/mysql_cmd/pkg"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"strings"
)

type tableInfo struct {
	TbName    string         `db:"table_name"`
	TbComment sql.NullString `db:"table_comment"`
}

type tableColumnInfo struct {
	OrdinalPosition uint16         `db:"ORDINAL_POSITION"`
	ColumnName      string         `db:"COLUMN_NAME"`
	ColumnType      string         `db:"COLUMN_TYPE"`
	DataType        string         `db:"DATA_TYPE"`
	ColumnKey       sql.NullString `db:"COLUMN_KEY"`
	IsNullable      string         `db:"IS_NULLABLE"`
	Extra           sql.NullString `db:"EXTRA"`
	ColumnComment   sql.NullString `db:"COLUMN_COMMENT"`
	ColumnDefault   sql.NullString `db:"COLUMN_DEFAULT"`
}

var (
	dbAddr    string
	dbUser    string
	dbPass    string
	dbName    string
	genTables string
)

func init() {
	addr := flag.String("addr", "", "请输入 数据库 地址，如：127.0.0.1:3306\n")
	user := flag.String("user", "", "请输入 数据库 用户名\n")
	pass := flag.String("pass", "", "请输入 数据库 密码\n")
	name := flag.String("name", "", "请输入 数据库 名称\n")
	table := flag.String("tables", "*", "请输入 table 名称，默认“*”，多个 “,”分割\n")

	flag.Parse()

	dbAddr = *addr
	dbUser = *user
	dbPass = *pass
	dbName = strings.ToLower(*name)
	genTables = strings.ToLower(*table)
}

func getDbConn(dbAddr, dbUser, dbPass, dbName string) *gorm.DB{
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=%t&loc=%s",
		dbUser,
		dbPass,
		dbAddr,
		dbName,
		true,
		"Local")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		log.Fatalf("connection failed,Database name: %s\n", dbName)
	}

	return db
}

func main() {
	dbConn := getDbConn(dbAddr, dbUser, dbPass, dbName)
	tables, err := getTables(dbConn, dbName, genTables)
	if err != nil {
		log.Println("query tables of database err", err)
		return
	}

	for _, table := range tables {
		//创建model目录 一个表一个model目录
		filepath := "./model/mysql/" + table.TbName + "_model"
		err := os.Mkdir(filepath, 0766)
		fmt.Printf("%v\n: ", err)
		fmt.Println("create dir : ", filepath)

		//查询当前表字段列
		columnInfo, columnInfoErr := getTableColumn(dbConn, dbName, table.TbName)
		if columnInfoErr != nil {
			continue
		}

		//生成struct model.go 代码
		genModelFile(filepath, table, dbName, columnInfo)

		//生成curd handler.go 代码
		genCurdFile(filepath, table.TbName, columnInfo)
	}
}

func genModelFile(filepath string, table tableInfo, dbName string, columnInfo []tableColumnInfo)  {
	modelFileName := fmt.Sprintf("%s/model.go", filepath)
	modelFile, err := os.OpenFile(modelFileName, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0766)
	if err != nil {
		panic(err)
	}

	type modelStruct struct {
		QueryBuilderName string
		PkgName string
		StructName string
		TableName string
		TableComment sql.NullString
		DbName string
		Quote string
		Fields []tableColumnInfo
	}

	var data = modelStruct {
		table.TbName + "ModelQueryBuilder",
		table.TbName + "_model",
		convertUpperCamelCase(table.TbName),
		table.TbName,
		table.TbComment,
		dbName,
		"`",
		columnInfo,
	}

	err = pkg.ModelTemplate.Execute(modelFile, data)
	if err != nil {
		panic(err)
	}
	fmt.Println("  └── file : ", table.TbName + "_model/model.go")
}

func genCurdFile(filepath, tableName string, columnInfo []tableColumnInfo)  {
	curlFileName := fmt.Sprintf("%s/handler.go", filepath)
	curlFile, err := os.OpenFile(curlFileName, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0766)
	if err != nil {
		panic(err)
	}

	type curdStruct struct {
		QueryBuilderName string
		PkgName string
		StructName string
		PkFieldName string
		PkFieldType string
		Quote string
	}

	pkField,pkFieldType := getPkField(columnInfo)
	var data = curdStruct {
		convertUpperCamelCase(tableName) + "ModelQueryBuilder",
		tableName + "_model",
		convertUpperCamelCase(tableName),
		pkField,
		pkFieldType,
		"`",
	}
	err = pkg.CurdTemplate.Execute(curlFile, data)
	if err != nil {
		panic(err)
	}
	fmt.Println("  └── file : ", tableName + "_model/handler.go")
}

//获得主键字段名称
func getPkField(columnInfo []tableColumnInfo)  (string,string) {
	for _,column := range columnInfo {
		if column.ColumnKey.String == "PRI" {
			return convertUpperCamelCase(column.ColumnName),pkg.GetGoTypeBySqlType(column.DataType)
		}
	}

	return "",""
}

func getTables(db *gorm.DB, dbName string, tableName string) ([]tableInfo, error) {
	var tableInfoMap = map[string]tableInfo{}
	var tableInfoArrayNew []tableInfo

	sql := fmt.Sprintf("SELECT `table_name`,`table_comment` FROM `information_schema`.`tables` WHERE `table_schema`= '%s'", dbName)
	rows, err := db.Raw(sql).Rows()
	if err != nil {
		return []tableInfo{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var info tableInfo
		err = rows.Scan(&info.TbName, &info.TbComment)
		if err != nil {
			fmt.Printf("query tables error,detail: [%v]\n", err.Error())
			continue
		}

		tableInfoMap[info.TbName] = info
	}

	// 指定表的时候，用指定的表 收集 tableInfoArray
	if tableName != "*" {
		inputTables := strings.Split(tableName, ",")

		for _, tableNameItem := range inputTables {
			if  info, OK := tableInfoMap[tableNameItem]; OK {
				tableInfoArrayNew = append(tableInfoArrayNew, info)
			}
		}
	} else {
		for _, info := range tableInfoMap {
			tableInfoArrayNew = append(tableInfoArrayNew, info)
		}
	}

	return tableInfoArrayNew, err
}

func getTableColumn(db *gorm.DB, dbName string, tableName string) ([]tableColumnInfo, error) {
	// 定义承载列信息的切片
	var columns []tableColumnInfo

	sql := fmt.Sprintf("SELECT `ORDINAL_POSITION`,`COLUMN_NAME`,`COLUMN_TYPE`,`DATA_TYPE`,`COLUMN_KEY`,`IS_NULLABLE`,`EXTRA`,`COLUMN_COMMENT`,`COLUMN_DEFAULT` FROM `information_schema`.`columns` WHERE `table_schema`= '%s' AND `table_name`= '%s' ORDER BY `ORDINAL_POSITION` ASC",
		dbName, tableName)

	rows, err := db.Raw(sql).Rows()
	if err != nil {
		fmt.Printf("query table column error, detail: [%v]\n", err.Error())
		return columns, err
	}
	defer rows.Close()

	for rows.Next() {
		var column tableColumnInfo
		err = rows.Scan(
			&column.OrdinalPosition,
			&column.ColumnName,
			&column.ColumnType,
			&column.DataType,
			&column.ColumnKey,
			&column.IsNullable,
			&column.Extra,
			&column.ColumnComment,
			&column.ColumnDefault)
		if err != nil {
			fmt.Printf("query table column error, detail: [%v]\n", err.Error())
			return columns, err
		}
		columns = append(columns, column)
	}

	return columns, err
}

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