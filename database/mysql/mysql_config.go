package mysql

import (
	"time"

	_ "github.com/go-sql-driver/mysql" // import your used driver
	//"winkim/baselib/logger"
	"store-management-be/baselib/logger"
	"github.com/jmoiron/sqlx"
)

type MySQLConfig struct {
	Name   string // for trace
	DSN    string // data source name
	Active int    // pool
	Idle   int    // pool
}

func NewSqlxDB(c MySQLConfig) (db *Db) {
	d, err := sqlx.Connect("mysql", c.DSN)
	if err != nil {
		logger.LogSugar.Errorf("Connect db error: %v", err)
		return
	}
	db = &Db{c.Name, d}
	db.SetMaxOpenConns(c.Active)
	db.SetMaxIdleConns(c.Idle)
	db.SetConnMaxLifetime(30 * time.Second)
	return
}
