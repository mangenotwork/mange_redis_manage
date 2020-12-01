package dao

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/mangenotwork/mange_redis_manage/common/manlog"
	"github.com/mangenotwork/mange_redis_manage/common/sqlitedb"
	"github.com/mangenotwork/mange_redis_manage/repository"
	_ "github.com/mattn/go-sqlite3"
)

type DaoRedisCluster struct {
	Data *models.RedisClusterDB
	Mu   sync.Mutex
}

//提取查询数据
func (this *DaoRedisCluster) exportdatas(rows *sql.Rows) (datas []*models.RedisClusterDB, err error) {
	for rows.Next() {
		data := &models.RedisClusterDB{}
		err := rows.Scan(&data.ID, &data.Hid, &data.GetTime, &data.ClusterEnabled)
		if err != nil {
			manlog.Error(err)
		}
		//manlog.Debug(*data)
		datas = append(datas, data)
	}
	return
}

//提取查询数据
func (this *DaoRedisCluster) exportdatas1(rows *sql.Rows) (data *models.RedisClusterDB, err error) {
	for rows.Next() {
		data = &models.RedisClusterDB{}
		err := rows.Scan(&data.ID, &data.Hid, &data.GetTime, &data.ClusterEnabled)
		if err != nil {
			manlog.Error(err)
		}
	}
	return
}

func (this *DaoRedisCluster) Create() error {
	db := sqlitedb.GetDBConn()
	this.Mu.Lock()
	defer this.Mu.Unlock()
	stmt, err := db.Prepare("INSERT INTO table_redis_servers_cluster (hid,get_time,cluster_enabled)" +
		" values(?,?,?)")
	defer stmt.Close()
	if err != nil {
		manlog.Error(err)
		return err
	}
	res, err := stmt.Exec(this.Data.Hid, this.Data.GetTime, this.Data.ClusterEnabled)
	manlog.Debug(&res)
	if err != nil {
		manlog.Error(err)
		return err
	}
	return nil
}

//查询当前时间之前的n条数据
func (this *DaoRedisCluster) GetLastTimeData(rid string, last_time, n int64) (data []*models.RedisClusterDB, err error) {
	sql := fmt.Sprintf("SELECT * FROM table_redis_servers_cluster where hid='%s' and get_time<%d Order by id desc LIMIT %d;", rid, last_time, n)
	rows, err := sqlitedb.Query(sql)
	defer rows.Close()
	if err != nil {
		manlog.Error(err)
	}
	return this.exportdatas(rows)
}

//删除多久之前的数据
func (this *DaoRedisCluster) DelTimeData(del_time int64) error {
	this.Mu.Lock()
	defer this.Mu.Unlock()
	sql := "delete from table_redis_servers_cluster where get_time<?"
	return sqlitedb.Del(sql, del_time)
}

//清空数据
func (this *DaoRedisCluster) Empty() error {
	this.Mu.Lock()
	defer this.Mu.Unlock()
	sql := "delete from table_redis_servers_cluster"
	return sqlitedb.Del(sql)
}
