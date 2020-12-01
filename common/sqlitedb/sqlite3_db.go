package sqlitedb

import (
	"database/sql"
	_ "fmt"

	"github.com/mangenotwork/mange_redis_manage/common/manlog"
	_ "github.com/mattn/go-sqlite3"
)

var sdb *sql.DB
var err error

//初始化连接
func init() {
	sdb, err = sql.Open("sqlite3", "./db/mange_redis_manage.db")
	if err != nil {
		manlog.Panic("连接本地数据文件错误, err = ", err)
	}
}

func Exec(sql string) {
	manlog.Debug("\n [执行sql] -> ", sql)
	sdb.Exec(sql)
}

func Query(sql string) (*sql.Rows, error) {
	manlog.Debug("\n [执行sql] -> ", sql)
	rows, err := sdb.Query(sql)
	// defer rows.Close()
	return rows, err
}

//删除命令的封装
func Del(sql string, arg ...interface{}) error {
	manlog.Debug("\n [执行sql] -> ", sql, arg)
	db := GetDBConn()
	stmt, err := db.Prepare(sql)
	defer stmt.Close()
	if err != nil {
		manlog.Error(err)
		return err
	}
	res, err := stmt.Exec(arg...)
	manlog.Debug(&res)
	if err != nil {
		manlog.Error(err)
		return err
	}
	return nil
}

func GetDBConn() *sql.DB {
	return sdb
}

func Prepare(sql string, args ...interface{}) error {
	manlog.Debug("\n [执行sql] -> ", sql)
	manlog.Debug("\n [接收sql参数] -> ", args)
	stmt, err := sdb.Prepare(sql)
	if err != nil {
		manlog.Error("sql err = ", err)
		return err
	}
	_, err = stmt.Exec(args...)
	if err != nil {
		manlog.Error("args err = ", err)
		return err
	}
	return nil
}
