package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql" // Needed for connect mysql
	"hcc/flute/config"
	"hcc/flute/logger"
	"strconv"
)

// Db : Pointer of mysql connection
var Db *sql.DB

// Prepare : Connect to mysql and prepare pointer of mysql connection
func Prepare() error {
	var err error
	Db, err = sql.Open("mysql",
		config.Mysql.ID+":"+config.Mysql.Password+"@tcp("+
			config.Mysql.Address+":"+strconv.Itoa(int(config.Mysql.Port))+")/"+
			config.Mysql.Database+"?parseTime=true")
	if err != nil {
		logger.Logger.Println(err)
		return err
	}

	logger.Logger.Println("db is connected")

	err = Db.Ping()
	if err != nil {
		logger.Logger.Println(err.Error())
		return err
	}

	return nil
}
