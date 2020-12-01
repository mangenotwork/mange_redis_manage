package dao

// type RedisReplicationDB struct {
// 	ID                         int64  `gorm:"primary_key;column:id"`
// 	Hid                        string `gorm:"column:hid"`
// 	GetTime                    int64  `gorm:"column:get_time"`
// 	Role                       string `gorm:"column:role_vlue"`                          //
// 	ConnectedSlaves            string `gorm:"column:connected_slaves"`              //
// 	MasterReplid               string `gorm:"column:master_replid"`                 //
// 	MasterReplid2              string `gorm:"column:master_replid2"`                //
// 	MasterReplOffset           string `gorm:"column:master_repl_offset"`            //主从同步偏移量（通过主从对比判断主从同步是否一致）
// 	SecondReplOffset           string `gorm:"column:second_repl_offset"`            //
// 	ReplBacklogActive          string `gorm:"column:repl_backlog_active"`           //
// 	ReplBacklogSize            int64  `gorm:"column:repl_backlog_size"`             //
// 	ReplBacklogFirstByteOffset int64  `gorm:"column:repl_backlog_first_byteoffset"` //
// 	ReplBacklogHistlen         int64  `gorm:"column:repl_backlog_histlen"`          //
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

type DaoRedisReplication struct {
	Data *models.RedisReplicationDB
	Mu   sync.Mutex
}

//提取查询数据
func (this *DaoRedisReplication) exportdatas(rows *sql.Rows) (datas []*models.RedisReplicationDB, err error) {
	for rows.Next() {
		data := &models.RedisReplicationDB{}
		err := rows.Scan(&data.ID, &data.Hid, &data.GetTime, &data.Role, &data.ConnectedSlaves, &data.MasterReplid, &data.MasterReplid2,
			&data.MasterReplOffset, &data.SecondReplOffset, &data.ReplBacklogActive, &data.ReplBacklogSize, &data.ReplBacklogFirstByteOffset,
			&data.ReplBacklogHistlen)
		if err != nil {
			manlog.Error(err)
		}
		//manlog.Debug(*data)
		datas = append(datas, data)
	}
	return
}

//提取查询数据
func (this *DaoRedisReplication) exportdatas1(rows *sql.Rows) (data *models.RedisReplicationDB, err error) {
	for rows.Next() {
		data = &models.RedisReplicationDB{}
		err := rows.Scan(&data.ID, &data.Hid, &data.GetTime, &data.Role, &data.ConnectedSlaves, &data.MasterReplid, &data.MasterReplid2,
			&data.MasterReplOffset, &data.SecondReplOffset, &data.ReplBacklogActive, &data.ReplBacklogSize, &data.ReplBacklogFirstByteOffset,
			&data.ReplBacklogHistlen)
		if err != nil {
			manlog.Error(err)
		}
	}
	return
}

func (this *DaoRedisReplication) Create() error {
	db := sqlitedb.GetDBConn()
	this.Mu.Lock()
	defer this.Mu.Unlock()
	stmt, err := db.Prepare("INSERT INTO table_redis_servers_replication (hid,get_time,role_value,connected_slaves,master_replid,master_replid2,master_repl_offset,second_repl_offset," +
		"repl_backlog_active,repl_backlog_size,repl_backlog_first_byteoffset,repl_backlog_histlen)" +
		" values(?,?,?,?,?,?,?,?,?,?,?,?)")
	defer stmt.Close()
	if err != nil {
		manlog.Error(err)
		return err
	}
	res, err := stmt.Exec(this.Data.Hid, this.Data.GetTime, this.Data.Role, this.Data.ConnectedSlaves, this.Data.MasterReplid, this.Data.MasterReplid2,
		this.Data.MasterReplOffset, this.Data.SecondReplOffset, this.Data.ReplBacklogActive, this.Data.ReplBacklogSize, this.Data.ReplBacklogFirstByteOffset,
		this.Data.ReplBacklogHistlen)
	manlog.Debug(&res)
	if err != nil {
		manlog.Error(err)
		return err
	}
	return nil
}

//查询当前时间之前的n条数据
func (this *DaoRedisReplication) GetLastTimeData(rid string, last_time, n int64) (data []*models.RedisReplicationDB, err error) {
	sql := fmt.Sprintf("SELECT * FROM table_redis_servers_replication where hid='%s' and get_time<%d Order by id desc LIMIT %d;", rid, last_time, n)
	rows, err := sqlitedb.Query(sql)
	defer rows.Close()
	if err != nil {
		manlog.Error(err)
	}
	return this.exportdatas(rows)
}

//删除多久之前的数据
func (this *DaoRedisReplication) DelTimeData(del_time int64) error {
	this.Mu.Lock()
	defer this.Mu.Unlock()
	sql := "delete from table_redis_servers_replication where get_time<?"
	return sqlitedb.Del(sql, del_time)
}

//清空数据
func (this *DaoRedisReplication) Empty() error {
	this.Mu.Lock()
	defer this.Mu.Unlock()
	sql := "delete from table_redis_servers_replication"
	return sqlitedb.Del(sql)
}
