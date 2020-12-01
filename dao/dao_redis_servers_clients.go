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

type DaoRedisClients struct {
	Data *models.RedisClientsDB
	Mu   sync.Mutex
}

//提取查询数据
func (this *DaoRedisClients) exportdatas(rows *sql.Rows) (datas []*models.RedisClientsDB, err error) {
	for rows.Next() {
		data := &models.RedisClientsDB{}
		err := rows.Scan(&data.ID, &data.Hid, &data.GetTime, &data.ConnectedClients, &data.ClientRecentMaxInputBuffer, &data.ClientRecentMaxOutputBuffer,
			&data.BlockedClients)
		if err != nil {
			manlog.Error(err)
		}
		//manlog.Debug(*data)
		datas = append(datas, data)
	}
	return
}

//提取查询数据
func (this *DaoRedisClients) exportdatas1(rows *sql.Rows) (data *models.RedisClientsDB, err error) {
	for rows.Next() {
		data = &models.RedisClientsDB{}
		err := rows.Scan(&data.ID, &data.Hid, &data.GetTime, &data.ConnectedClients, &data.ClientRecentMaxInputBuffer, &data.ClientRecentMaxOutputBuffer,
			&data.BlockedClients)
		if err != nil {
			manlog.Error(err)
		}
	}
	return
}

func (this *DaoRedisClients) Create() error {
	db := sqlitedb.GetDBConn()
	this.Mu.Lock()
	defer this.Mu.Unlock()
	stmt, err := db.Prepare("INSERT INTO table_redis_servers_clients (hid,get_time,connected_clients,client_recent_max_input_buffer,client_recent_max_output_buffer,clocked_clients)" +
		" values(?,?,?,?,?,?)")
	defer stmt.Close()
	if err != nil {
		manlog.Error(err)
		return err
	}
	res, err := stmt.Exec(this.Data.Hid, this.Data.GetTime, this.Data.ConnectedClients, this.Data.ClientRecentMaxInputBuffer, this.Data.ClientRecentMaxOutputBuffer,
		this.Data.BlockedClients)
	manlog.Debug(&res)
	if err != nil {
		manlog.Error(err)
		return err
	}
	return nil
}

//获取最新的一条数据
func (this *DaoRedisClients) GetNewData(rid string) (data *models.RedisClientsDB, err error) {

	sql := fmt.Sprintf("SELECT * FROM table_redis_servers_clients where hid='%s' Order by get_time desc LIMIT 1;", rid)
	rows, err := sqlitedb.Query(sql)
	defer rows.Close()
	if err != nil {
		manlog.Error(err)
	}
	return this.exportdatas1(rows)
}

//查询当前时间之前的n条数据
func (this *DaoRedisClients) GetLastTimeData(rid string, last_time, n int64) (data []*models.RedisClientsDB, err error) {
	sql := fmt.Sprintf("SELECT * FROM table_redis_servers_clients where hid='%s' and get_time<%d Order by id desc LIMIT %d;", rid, last_time, n)
	rows, err := sqlitedb.Query(sql)
	defer rows.Close()
	if err != nil {
		manlog.Error(err)
	}
	return this.exportdatas(rows)
}

//删除多久之前的数据
func (this *DaoRedisClients) DelTimeData(del_time int64) error {
	this.Mu.Lock()
	defer this.Mu.Unlock()
	sql := "delete from table_redis_servers_clients where get_time<?"
	return sqlitedb.Del(sql, del_time)
}

//清空数据
func (this *DaoRedisClients) Empty() error {
	this.Mu.Lock()
	defer this.Mu.Unlock()
	sql := "delete from table_redis_servers_clients"
	return sqlitedb.Del(sql)
}
