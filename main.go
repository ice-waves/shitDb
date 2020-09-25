package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"os/exec"
	"strings"
	"unicode"
)

var TypeMapping = map[string]string{
	"tinyint":           "int8",
	"smallint":          "int16",
	"integer":           "int32",
	"int":               "int32",
	"bigint":            "int64",
	"integer unsigned":  "uint",
	"tinyint unsigned":  "uint8",
	"smallint unsigned": "uint16",
	"bigint unsigned":   "uint64",
	"double precision":  "float32",
	"float":             "float32",
	"bool":              "bool",
	"text":              "string",
	"longtext":          "string",
	"mediumtext":        "string",
	"varchar":           "string",
	"char":              "string",
	"enum":              "string",
	"date":              "time.Time",
	"datetime":          "time.Time",
	"timestamp":         "time.Time",
}

type Field struct {
	fieldName string
	fieldDesc string
	dataType  string
	isNull    string
	length    int
}

type SqlDb struct {
	db *sql.DB
}

var host = flag.String("h", "", "数据库地址")
var database = flag.String("d", "", "数据库")
var table = flag.String("t", "", "数据表")
var modelPath = flag.String("model_path", "", "model路径")
var modelName = flag.String("model_name", "", "model名称")

func main() {
	flag.Parse()
	fmt.Println("gen sql struct!")
	// 创建model
	modelPathSplit := strings.Split(*modelPath, "/")
	packageName := modelPathSplit[len(modelPathSplit)-1]
	fileName := *modelPath + "/" + *modelName + ".go"
	db, err := sql.Open("mysql", *host)
	defer db.Close()
	checkError(err)
	sqlDb := SqlDb{db: db}
	fieldInfo := sqlDb.FieldInfo(*database, *table)
	fmt.Println(fieldInfo)
	createModel(*modelPath, fileName, packageName, fieldInfo)

	// 格式化文件
	cmd := exec.Command("go", "fmt", fileName)
	err = cmd.Run()
	checkError(err)
	fmt.Println("生成model成功")
}

func (sqlDb SqlDb) TableInfo(dbName string) map[string]string {
	sqlStr := `SELECT table_name tableName,TABLE_COMMENT tableDesc
			FROM INFORMATION_SCHEMA.TABLES 
			WHERE UPPER(table_type)='BASE TABLE'
			AND LOWER(table_schema) = ? 
			ORDER BY table_name asc`

	var result = make(map[string]string)
	rows, err := sqlDb.db.Query(sqlStr, dbName)

	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var tableName, tableDesc string
		err = rows.Scan(&tableName, &tableDesc)

		if len(tableDesc) == 0 {
			tableDesc = tableName
		}
		result[tableName] = tableDesc
	}

	return result
}

func (sqlDb SqlDb) FieldInfo(dbName, tableName string) []Field {
	sqlStr := `SELECT COLUMN_NAME fName,column_comment fDesc,DATA_TYPE dataType,
						IS_NULLABLE isNull,IFNULL(CHARACTER_MAXIMUM_LENGTH,0) sLength
			FROM information_schema.columns 
			WHERE table_schema = ? AND table_name = ?`

	var result []Field

	rows, err := sqlDb.db.Query(sqlStr, dbName, tableName)

	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var f Field
		err = rows.Scan(&f.fieldName, &f.fieldDesc, &f.dataType, &f.isNull, &f.length)
		if err != nil {
			panic(err)
		}

		result = append(result, f)
	}
	return result
}

func Ucfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func UnderlineToBigCamel(str, delimiter string) string {
	strSplit := strings.Split(str, delimiter)
	if len(strSplit) == 0 {
		return ""
	}

	for key, value := range strSplit {
		strSplit[key] = Ucfirst(value)
	}

	return strings.Join(strSplit, "")
}

func Lcfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}

}

func createModel(modelPath, modelFile, packageName string, fieldInfo []Field) {
	err := os.MkdirAll(modelPath, os.ModePerm)
	checkError(err)
	fmt.Println("创建路径成功...")

	// 创建model
	file, err := os.OpenFile(modelFile, os.O_RDWR|os.O_CREATE, 0766)
	checkError(err)
	defer file.Close()

	// 使用带缓冲区写入文件
	write := bufio.NewWriter(file)
	_, err = write.WriteString("package" + " " + packageName)
	upModelName := UnderlineToBigCamel(*modelName, "_")
	checkError(err)
	write.WriteString("\n\n\ntype " + upModelName + " struct {\n")

	for _, value := range fieldInfo {
		var columnStr = ""
		if value.fieldName == "id" {
			columnStr = "\t" + strings.ToUpper(value.fieldName) + "\t\t" + TypeMapping[value.dataType] + "\t\t" + "`gorm:column:" + value.fieldName + " json:\"" + value.fieldName + "\"" + "`\n"
		} else {
			columnStr = "\t" + UnderlineToBigCamel(value.fieldName, "_") + "\t\t" + TypeMapping[value.dataType] + "\t\t" + "`gorm:column:" + value.fieldName + " json:\""+ value.fieldName+ "\"" + "`\n"
		}
		write.WriteString(columnStr)
	}
	write.WriteString("\n}\n\n")

	write.WriteString("func ("+ upModelName +") TableName() string {\n\treturn" + " \"" + *modelName + "\"\n}")

	write.Flush()

}
