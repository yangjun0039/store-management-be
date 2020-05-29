package mysql

import (
	"github.com/jmoiron/sqlx"
	"sync"
	//"red-envelope/configer"
	"fmt"
	"github.com/pkg/errors"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
)

type Db struct {
	Name string
	*sqlx.DB
}

//var db *Db
var DbMap map[string]*Db
var mu sync.Mutex
var once sync.Once
var dbNilError = errors.New("db is nil")
var mysqlConfigs []MySQLConfig

func InitMysql(configs []MySQLConfig) {
	if configs == nil || len(configs) == 0 {
		return
	}
	mysqlConfigs = configs
	New()
}

func GetDB(name string) *Db {
	if DbMap == nil {
		openDB()
	}
	return DbMap[name]
}

func New() {
	if DbMap == nil {
		once.Do(openDB)
	} else {
		fmt.Println("db already exists")
	}
}

func openDB() {
	mu.Lock()
	defer mu.Unlock()
	if DbMap != nil {
		return
	}
	DbMap = make(map[string]*Db)
	for _, c := range mysqlConfigs {
		d := NewSqlxDB(c)
		if d == nil {
			panic(dbNilError)
		}
		DbMap[d.Name] = d
	}
}

// 插入
func (self *Db) Insert(sqlstr string, args ...interface{}) (int64, error) {
	if self == nil {
		return 0, dbNilError
	}
	stmtIns, err := self.Prepare(sqlstr)
	if err != nil {
		return 0, err
	}
	defer stmtIns.Close()
	result, err := stmtIns.Exec(args...)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// 插入(有事务处理) 调用时不用Rollback
func (self *Db) InsertTx(tx *sql.Tx, sqlstr string, args ...interface{}) (int64, error) {
	if self == nil {
		return 0, dbNilError
	}
	result, err := tx.Exec(sqlstr, args...)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	n, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return n, nil
}

// 修改和删除(有事务处理) 调用时不用Rollback
func (self *Db) ExecTx(tx *sql.Tx, sqlstr string, args ...interface{}) (int64, error) {
	if self == nil {
		tx.Rollback()
		return 0, dbNilError
	}
	result, err := tx.Exec(sqlstr, args...)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	n, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return n, nil
}

// 修改和删除
func (self *Db) Exec(sqlstr string, args ...interface{}) (int64, error) {
	if self == nil {
		return 0, dbNilError
	}
	stmtIns, err := self.Prepare(sqlstr)
	if err != nil {
		return 0, err
	}
	defer stmtIns.Close()
	result, err := stmtIns.Exec(args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// 查询一行数据
func (self *Db) FetchRowD(sqlstr string, args ...interface{}) (map[string]string, error) {
	if self == nil {
		return nil, dbNilError
	}
	rows, err := self.Query(sqlstr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ret, err := fetchData(rows)
	if err != nil {
		return nil, err
	}
	if len(ret) > 0 {
		return ret[0], nil
	}
	return map[string]string{}, nil
}

// 查询多行数据
func (self *Db) FetchRowsD(sqlstr string, args ...interface{}) ([]map[string]string, error) {
	if self == nil {
		return nil, dbNilError
	}
	rows, err := self.Query(sqlstr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ret, err := fetchData(rows)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func fetchData(rows *sql.Rows) ([]map[string]string, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))

	ret := make([]map[string]string, 0)
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		var value string
		vmap := make(map[string]string, len(scanArgs))
		for i, col := range values {
			if col == nil {
				value = ""
			} else {
				value = string(col)
			}
			vmap[columns[i]] = value
		}
		ret = append(ret, vmap)
	}
	return ret, nil
}
