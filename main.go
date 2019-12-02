package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type TableField struct {
	Field      string `json:"field"`
	Type       string `json:"type"`
	Collation  string `json:"collation"`
	Null       string `json:"null"`
	Key        string `json:"key"`
	Default    string `json:"default"`
	Extra      string `json:"extra"`
	Privileges string `json:"privileges"`
	Comment    string `json:"comment"`
}

var (
	Config    *viper.Viper
	DbEngin   *gorm.DB
	tableName string
	help      bool
)

func main() {
	flag.Parse()
	if help {
		flag.Usage()
		return
	}
	initConfig()
	switch Config.GetString("database.type") {
	case "mysql":
		initMysqlDb()
		tableFields, err := getMysqlTableFields(tableName)
		if err != nil {
			log.Fatal(err)
		}
		genStruct(tableName, tableFields)
		return
	case "postgresql":
		initPostgresqlDb()
		tableFields, err := getPostgresqlTableFields(tableName)
		if err != nil {
			log.Fatal(err)
		}
		genStruct(tableName, tableFields)
		return
	}
	log.Fatal("error database.type")
}

func genStruct(tableName string, tableFields []TableField) {
	tableNameHandlers := Config.GetStringSlice("tableNameHandler")
	for _, tableNameHeader := range tableNameHandlers {
		tableName = strings.Replace(tableName, tableNameHeader, "", 1)
	}
	fmt.Println(fmt.Sprintf("type %v struct {", marshal(tableName)))
	for _, tableField := range tableFields {
		fmt.Println(fmt.Sprintf("\t %v \t %v \t `json:\"%v\"` //%v", marshal(tableField.Field), changeType(tableField.Type), gorm.ToColumnName(tableField.Field), tableField.Comment))
	}
	fmt.Println("}")
}

/*
获取表字典
*/
func getMysqlTableFields(tableName string) ([]TableField, error) {
	sqlTableDesc := "show full fields from %v"

	rows, err := DbEngin.Raw(fmt.Sprintf(sqlTableDesc, tableName)).Rows()
	if err != nil {
		return nil, err
	}

	tableFields := []TableField{}
	Field, Type, Collation, Null, Key, Default, Extra, Privileges, Comment := []byte{}, []byte{}, []byte{}, []byte{}, []byte{}, []byte{}, []byte{}, []byte{}, []byte{}
	for rows.Next() {
		rows.Scan(&Field, &Type, &Collation, &Null, &Key, &Default, &Extra, &Privileges, &Comment)
		tableFields = append(tableFields, TableField{
			Field:      string(Field),
			Type:       string(Type),
			Collation:  string(Collation),
			Null:       string(Null),
			Key:        string(Key),
			Default:    string(Default),
			Extra:      string(Extra),
			Privileges: string(Privileges),
			Comment:    string(Comment),
		})
	}
	return tableFields, nil
}

/*
获取表字典
*/
func getPostgresqlTableFields(tableName string) ([]TableField, error) {
	sqlTableDesc := `SELECT A.attname AS field, format_type ( A.atttypid, A.atttypmod ) AS type, col_description ( A.attrelid, A.attnum ) AS comment FROM pg_class AS C, pg_attribute AS A WHERE C.relname = '%v' AND A.attrelid = C.oid AND A.attnum > 0`

	rows, err := DbEngin.Raw(fmt.Sprintf(sqlTableDesc, tableName)).Rows()
	if err != nil {
		return nil, err
	}

	tableFields := []TableField{}
	Field, Type, Comment := []byte{}, []byte{}, []byte{}
	for rows.Next() {
		rows.Scan(&Field, &Type, &Comment)
		tableFields = append(tableFields, TableField{
			Field:   string(Field),
			Type:    string(Type),
			Comment: string(Comment),
		})
	}
	return tableFields, nil
}

/**
加载配置文件
*/
func initConfig() error {
	configFile := "t2g.yaml"
	Config = viper.New()
	Config.SetConfigFile(configFile)
	err := Config.ReadInConfig()
	if err != nil {
		log.Println("InitConfig conf file:%v error", configFile)
		panic(err)
		return err
	}
	return nil
}

/**
连接数据库
*/
func initMysqlDb() {
	mysqlDb, err := gorm.Open("mysql", Config.GetString("database.mysqlConn"))
	if err != nil {
		log.Println("start mysql error:%v", err)
		panic("StartDB Error")
	}
	mysqlDb.DB().SetConnMaxLifetime(time.Minute * 5)
	mysqlDb.LogMode(true)
	mysqlDb.SingularTable(true)

	DbEngin = mysqlDb
}

/**
* 启动数据库
 */
func initPostgresqlDb() {
	postgreDb, err := gorm.Open("postgres", Config.GetString("database.postgresqlConn"))
	if err != nil {
		log.Println("start postgresql error:%v", err)
		panic("StartDB Error")
	}
	postgreDb.DB().SetConnMaxLifetime(time.Minute * 5)
	postgreDb.DB().SetMaxIdleConns(0)
	postgreDb.DB().SetMaxOpenConns(5)
	postgreDb.LogMode(true)
	postgreDb.SingularTable(true)
	postgreDb.LogMode(true)

	DbEngin = postgreDb
}

/*
	转换为大驼峰命名法则
	首字母大写，“_” 忽略后大写
*/
func marshal(name string) string {
	if name == "" {
		return ""
	}

	temp := strings.Split(name, "_")
	var s string
	for _, v := range temp {
		vv := []rune(v)
		if len(vv) > 0 {
			if bool(vv[0] >= 'a' && vv[0] <= 'z') { // 首字母大写
				vv[0] -= 32
			}
			s += string(vv)
		}
	}
	return s
}

func changeType(fieldType string) string {
	index := strings.Index(fieldType, "(")
	if index > 0 {
		fieldType = fieldType[0:index]
	}
	switch fieldType {
	case "int":
		return "int"
	case "bigint":
		return "int64"
	case "varchar":
		return "string"
	case "timestamp":
		return "time.Time"
	case "float":
		return "float64"
	case "text":
		return "string"
	case "integer":
		return "int64"
	case "timestamp without time zone":
		return "time.Time"
	}
	return fieldType
}

func init() {
	flag.BoolVar(&help, "h", false, "this help")
	flag.StringVar(&tableName, "t", "", "表名")
}
