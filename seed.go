package seed

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/spf13/viper"
	"github.com/walkerxy/go-seed/parse"
)

var (
	viperConf   = viper.New()
	queryColumn = "select COLUMN_NAME,COLUMN_TYPE,IS_NULLABLE,COLUMN_KEY,COLUMN_DEFAULT from information_schema.COLUMNS where table_name = '%s' and table_schema = '%s';"
)

// Seed Seed
type Seed struct {
	filepath    string
	filename    string
	database    string
	tablePrefix string
	db          *sql.DB
}

// initViper 初始化viper
func initViper() {
	// Init viper path and filename/type
	viperConf.AddConfigPath(".")
	viperConf.SetConfigName("conf")

	if err := viperConf.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatal("No such config file")
		} else {
			log.Fatal("Read config error")
		}
	}
	log.Println("Init viper complete")
}

// NewSeed NewSeed
func NewSeed(filepath string, filename string) *Seed {
	// 初始化viper
	initViper()

	host := viperConf.GetString("database.host")
	port := viperConf.GetString("database.port")
	databaseName := viperConf.GetString("database.database_name")
	username := viperConf.GetString("database.username")
	password := viperConf.GetString("database.password")
	dsn := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + databaseName

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic("Mysql open error : " + err.Error())
	}

	return &Seed{
		filepath:    filepath,
		filename:    filename,
		tablePrefix: viper.GetString("database.table_prefix"),
		database:    databaseName,
		db:          db,
	}
}

// SetTablePrefix 设置表前缀
func (seed *Seed) SetTablePrefix(tablePrefix string) {
	seed.tablePrefix = tablePrefix
}

// SetSeedFile 设置填充文件
func (seed *Seed) SetSeedFile(filepath string, filename string) {
	seed.filepath = filepath
	seed.filename = filename
}

// FetchSqls FetchSqls
func (seed *Seed) FetchSqls() map[string][]string {
	parse := parse.NewYaml(seed.filepath, seed.filename)
	tables := parse.Parse()

	sqls := make(map[string][]string)
	for tableName, tableData := range tables {
		tableName = seed.tablePrefix + "_" + tableName

		sqls["truncate"] = append(sqls["truncate"], "truncate table "+tableName)

		var trans []map[string]string
		for _, table := range tableData.([]interface{}) {
			sub := make(map[string]string)
			for colKey, colData := range table.(map[interface{}]interface{}) {
				sub[OtherToString(colKey)] = OtherToString(colData)
			}
			trans = append(trans, sub)
		}

		log.Println(fmt.Sprintf(queryColumn, tableName, seed.database))
		rows, err := seed.db.Query(fmt.Sprintf(queryColumn, tableName, seed.database))
		if err != nil {
			panic("Exec mysql error : " + err.Error())
		}

		defer rows.Close()

		for rows.Next() {
			var name string
			var colType string
			var colNullable string
			var colKey string
			var colDefault string
			rows.Scan(&name, &colType, &colNullable, &colKey, &colDefault)
			//log.Println(name, colType, colNullable, colKey, colDefault)
			for key := range trans {
				// 主键或者可为空的不用添加默认值
				if colNullable == "YES" || colKey == "PRI" {
					continue
				}
				if trans[key][name] == "" {
					SetDefaultValue(trans[key], name, colType, colDefault)
				}
			}
		}

		// Cellect sqls
		for _, item := range trans {
			sql := "insert into " + tableName + "("
			values := ""
			// 每一行数据
			for colKey, colData := range item {
				sql += colKey + ","
				values += "'" + colData + "',"
			}

			sql = sql[:len(sql)-1]
			values = values[:len(values)-1]

			sql += ") values (" + values + ");"
			sqls["insert"] = append(sqls["insert"], sql)
		}
	}

	return sqls
}

// Fill 填充
func (seed *Seed) Fill() bool {
	sqls := seed.FetchSqls()
	log.Println(sqls)

	tx, _ := seed.db.Begin()
	f := 0
	for _, sql := range sqls["truncate"] {
		log.Println("Exec sql is : " + sql)
		_, err := seed.db.Exec(sql)
		if err != nil {
			f++
			log.Println(err)
		}
	}
	for _, sql := range sqls["insert"] {
		log.Println("Exec sql is : " + sql)
		_, err := seed.db.Exec(sql)
		if err != nil {
			f++
			log.Println(err)
		}
	}

	if f == 0 {
		tx.Commit()
	} else {
		tx.Rollback()
		return false
	}

	return true
}
