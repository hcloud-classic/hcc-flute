package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql" // Needed for connect mysql
	"hcc/flute/lib/config"
	"hcc/flute/lib/logger"
	"strconv"
)

// Db : Pointer of mysql connection
var Db *sql.DB

func prepare() error {
	var err error
	Db, err = sql.Open("mysql",
		config.Mysql.ID+":"+config.Mysql.Password+"@tcp("+
			config.Mysql.Address+":"+strconv.Itoa(int(config.Mysql.Port))+")/"+
			config.Mysql.Database+"?parseTime=true")
	if err != nil {
		logger.Logger.Println(err)
		return err
	}

	err = Db.Ping()
	if err != nil {
		logger.Logger.Println(err)
		return err
	}

	logger.Logger.Println("db is connected")

	return nil
}
