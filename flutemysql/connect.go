package flutemysql

import (
	"GraphQL_Flute/fluteconfig"
	"GraphQL_Flute/flutelogger"
	"database/sql"
	_ "github.com/go-sql-driver/mysql" // Needed for connect mysql
)

// Db : Pointer of mysql connection
var Db *sql.DB

// Prepare : Connect to mysql and prepare pointer of mysql connection
func Prepare() error {
	var err error
	Db, err = sql.Open("mysql", fluteconfig.MysqlID+":"+fluteconfig.MysqlPassword+"@tcp("+
		fluteconfig.MysqlAddress+":"+fluteconfig.MysqlPort+")/"+fluteconfig.MysqlDatabase+"?parseTime=true")
	if err != nil {
		flutelogger.Logger.Println(err)
		return err
	}

	flutelogger.Logger.Println("db is connected")

	err = Db.Ping()
	if err != nil {
		flutelogger.Logger.Println(err.Error())
		return err
	}

	return nil
}
