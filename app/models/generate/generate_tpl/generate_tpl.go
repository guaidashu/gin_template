package generate_tpl

import (
	"bufio"
	"fmt"
	"gin_template/app/config"
	"gin_template/app/libs"
	"gin_template/app/models"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"unicode"

	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

func Generate(tableNames ...string) {
	var (
		tables        []Table
		tableNamesStr string
	)

	for _, name := range tableNames {
		if tableNamesStr != "" {
			tableNamesStr += ","
		}
		tableNamesStr += name
	}

	// 初始化配置文件
	config.InitConf()
	// 初始化数据库
	if err := models.InitDB(); err != nil {
		fmt.Println(err)
	} else {
		tables = getTables(tableNamesStr)
		for _, table := range tables {
			fields := getFields(table.Name)
			GenerateTpl(table, fields)
		}
	}

	if len(tables) == 0 {
		GenerateTpl(Table{
			Name:    tableNamesStr,
			Comment: "",
		}, nil)
	}
}

func GenerateTpl(table Table, fields []Field) {
	var (
		path      string
		content   string // fields 字段拼接
		readFile  *os.File
		writeFile *os.File
		err       error
		isExist   bool
	)

	tableName := table.Name
	// 检测表是否存在
	if models.GDB != nil {
		isExist = models.GDB.Migrator().HasTable(tableName)
	}
	if isExist {
		path = "./generate_tpl/tpl_from_db.tpl"
	} else {
		path = "./generate_tpl/tpl.tpl"
	}

	// 首先加载模板文件
	readFile, err = os.Open(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		_ = readFile.Close()
	}()

	// 打开 要一个临时文件，用于写入
	// 这里要判断用哪个模板文件
	newPath := tableName + ".go"
	writeFile, err = os.OpenFile(newPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return
	}
	defer func() {
		_ = writeFile.Close()
	}()

	reader := bufio.NewReader(readFile)
	tableUpper := libs.Case2Camel(tableName)

	if isExist {
		// 生成字段
		content = generator.CamelCase(table.Name) + "Model struct {\r\n"
		for _, field := range fields {
			fieldName := generator.CamelCase(field.Field)
			fieldGorm := getFieldGorm(field)
			fieldJson := getFieldJson(field)
			fieldType := getFiledType(field)
			fieldComment := getFieldComment(field)

			if fieldComment != "" {
				content += "\t\t" + fieldComment + "\r\n"
			}
			content += "\t\t" + fieldName + " " + fieldType + " `" + fieldGorm + " " + fieldJson + "`" + "\r\n"
		}
		content += "\t}"
	}

	// 读数据
	data, readErr := ioutil.ReadAll(reader)
	if readErr != nil {
		fmt.Println(readErr)
		return
	}
	s := string(data)

	s = strings.Replace(s, "TemplateId", tableUpper, -1)
	s = strings.Replace(s, "template_name", tableName, -1)
	s = strings.Replace(s, "TemplateModel", tableUpper+"Model", -1)
	s = strings.Replace(s, "templateModel", libs.LcFirst(tableUpper)+"Model", -1)
	s = strings.Replace(s, "_templateOnce", "_"+libs.LcFirst(tableUpper)+"Once", -1)
	s = strings.Replace(s, "template_id", tableName, -1)
	s = strings.Replace(s, "templateId", libs.LcFirst(tableUpper), -1)
	s = strings.Replace(s, "templateCacheKey", libs.LcFirst(tableUpper)+"CacheKey", -1)
	s = strings.Replace(s, "template#cache#", libs.LcFirst(tableUpper)+"#cache#", -1)
	s = strings.Replace(s, "TemplateList", tableUpper+"List", -1)
	s = strings.Replace(s, "${TemplateStruct}", content, -1)
	// 写入
	if _, err = writeFile.Write([]byte(s)); err != nil {
		fmt.Println(err)
	}
}

// 获取表信息
func getTables(tableNames string) []Table {
	var tables []Table
	if tableNames == "" {
		models.GDB.Raw("SELECT TABLE_NAME as Name,TABLE_COMMENT as Comment FROM information_schema.TABLES WHERE table_schema='" + config.Config.Mysql.Database + "';").Find(&tables)
	} else {
		models.GDB.Raw("SELECT TABLE_NAME as Name,TABLE_COMMENT as Comment FROM information_schema.TABLES WHERE TABLE_NAME IN ('" + tableNames + "') AND table_schema='" + config.Config.Mysql.Database + "';").Find(&tables)
	}
	return tables
}

// 获取所有字段信息
func getFields(tableName string) []Field {
	var fields []Field
	models.GDB.Raw("show FULL COLUMNS from `" + tableName + "`;").Find(&fields)
	return fields
}

// 获取字段类型
func getFiledType(field Field) string {
	typeArr := strings.Split(field.Type, "(")
	reg2 := regexp.MustCompile(" unsigned") // 正则匹配 unsigned
	isUnsigned := reg2.MatchString(field.Type)

	switch typeArr[0] {
	case "int", "integer", "mediumint", "year":
		if isUnsigned {
			return "int64"
		}
		return "int64"
	case "tinyint", "bit":
		if isUnsigned {
			return "int64"
		}
		return "int64"
	case "smallint":
		if isUnsigned {
			return "int64"
		}
		return "int64"
	case "bigint", "bigint unsigned":
		return "int64"
	case "decimal", "double", "float", "real", "numeric":
		return "float32"
	case "timestamp", "datetime", "time":
		return "time.Time"
	default:
		return "string"
	}
}

// 获取字段gorm描述
func getFieldGorm(field Field) string {
	if field.Field == "update_time" || field.Field == "create_time" {
		return `gorm:"column:` + field.Field + `;default:null"`
	}
	if field.Field == "created" {
		return `gorm:"column:` + field.Field + ";autoCreateTime" + `"`
	}
	if field.Field == "edited" {
		return `gorm:"column:` + field.Field + ";autoUpdateTime" + `"`
	}
	// 判断字段是否是主键
	if field.Key == "PRI" {
		return `gorm:"column:` + field.Field + ";primaryKey;autoIncrement" + `"`
	}
	if field.Type == "text" {
		return `gorm:"column:` + field.Field + `;` + "type:text" + `"`
	}
	return `gorm:"column:` + field.Field + `"`
}

// 获取字段json描述
func getFieldJson(field Field) string {
	return `json:"` + field.Field + `"`
}

// 下划线写法转为驼峰写法
func CaseCamel(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = strings.Title(name)
	name = Lcfirst(name)
	return strings.Replace(name, " ", "", -1)
}

func Lcfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

// 获取字段说明
func getFieldComment(field Field) string {
	if len(field.Comment) > 0 {
		return "// " + field.Comment
	}
	return ""
}

// 检查文件是否存在
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
