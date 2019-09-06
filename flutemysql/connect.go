package cellomysql

import (
	"GraphQL_Cello/celloconfig"
	"GraphQL_Cello/cellologger"
	"database/sql"
	_ "github.com/go-sql-driver/mysql" // Needed for connect mysql
)

// Db : Pointer of mysql connection
var Db *sql.DB

// Prepare : Connect to mysql and prepare pointer of mysql connection
func Prepare() error {
	var err error
	Db, err = sql.Open("mysql", celloconfig.MysqlID+":"+celloconfig.MysqlPassword+"@tcp("+
		celloconfig.MysqlAddress+":"+celloconfig.MysqlPort+")/"+celloconfig.MysqlDatabase+"?parseTime=true")
	if err != nil {
		cellologger.Logger.Println(err)
		return err
	}

	cellologger.Logger.Println("db is connected")

	err = Db.Ping()
	if err != nil {
		cellologger.Logger.Println(err.Error())
		return err
	}

	return nil
}
