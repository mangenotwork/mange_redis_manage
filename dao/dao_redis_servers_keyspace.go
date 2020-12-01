package dao

// type RedisKeyspaceDB struct {
// 	ID      int64  `gorm:"primary_key;column:id"`
// 	Hid     string `gorm:"column:hid"`
// 	GetTime int64  `gorm:"column:get_time"`
// 	DBID    int64  `gorm:"column:db_id"`
// 	Keys    int64  `gorm:"column:keys_count"`
// 	Expires int64  `gorm:"column:expires"`
// 	AvgTTL  int64  `gorm:"column:avgttl"`
// }

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/mangenotwork/mange_redis_manage/common/manlog"
	"github.com/mangenotwork/mange_redis_manage/common/sqlitedb"
	"github.com/mangenotwork/mange_redis_manage/repository"
	_ "github.com/mattn/go-sqlite3"
)

type DaoRedisKeyspace struct {
	Data *models.RedisKeyspaceDB
	Mu   sync.Mutex
}

//提取查询数据
func (this *DaoRedisKeyspace) exportdatas(rows *sql.Rows) (datas []*models.RedisKeyspaceDB, err error) {
	for rows.Next() {
		data := &models.RedisKeyspaceDB{}
		err := rows.Scan(&data.ID, &data.Hid, &data.GetTime, &data.DBID, &data.Keys, &data.Expires, &data.AvgTTL)
		if err != nil {
			manlog.Error(err)
		}
		//manlog.Debug(*data)
		datas = append(datas, data)
	}
	return
}

//提取查询数据
func (this *DaoRedisKeyspace) exportdatas1(rows *sql.Rows) (data *models.RedisKeyspaceDB, err error) {
	for rows.Next() {
		data = &models.RedisKeyspaceDB{}
		err := rows.Scan(&data.ID, &data.Hid, &data.GetTime, &data.DBID, &data.Keys, &data.Expires, &data.AvgTTL)
		if err != nil {
			manlog.Error(err)
		}
	}
	return
}

func (this *DaoRedisKeyspace) Create() error {
	db := sqlitedb.GetDBConn()
	this.Mu.Lock()
	defer this.Mu.Unlock()
	stmt, err := db.Prepare("INSERT INTO table_redis_servers_keyspace (hid,get_time,db_id,keys_count,expires,avgttl)" +
		" values(?,?,?,?,?,?)")
	defer stmt.Close()
	if err != nil {
		manlog.Error(err)
		return err
	}
	res, err := stmt.Exec(this.Data.Hid, this.Data.GetTime, this.Data.DBID, this.Data.Keys, this.Data.Expires, this.Data.AvgTTL)
	manlog.Debug(&res)
	if err != nil {
		manlog.Error(err)
		return err
	}
	return nil
}

//获取最新的一条数据
func (this *DaoRedisKeyspace) GetNewData(rid string, maxgietime int64) (data []*models.RedisKeyspaceDB, err error) {

	sql := fmt.Sprintf("SELECT * FROM table_redis_servers_keyspace where hid='%s' and get_time=%d;", rid, maxgietime)
	rows, err := sqlitedb.Query(sql)
	defer rows.Close()
	if err != nil {
		manlog.Error(err)
	}
	return this.exportdatas(rows)
}

//删除多久之前的数据
func (this *DaoRedisKeyspace) DelTimeData(del_time int64) error {
	this.Mu.Lock()
	defer this.Mu.Unlock()
	sql := "delete from table_redis_servers_keyspace where get_time<?"
	return sqlitedb.Del(sql, del_time)
}

//清空数据
func (this *DaoRedisKeyspace) Empty() error {
	this.Mu.Lock()
	defer this.Mu.Unlock()
	sql := "delete from table_redis_servers_keyspace"
	return sqlitedb.Del(sql)
}
