package dao

// type RedisPersistenceDB struct {
// 	ID                       int64  `gorm:"primary_key;column:id"`
// 	Hid                      string `gorm:"column:hid"`
// 	GetTime                  int64  `gorm:"column:get_time"`
// 	Loading                  int64  `gorm:"column:loading"`                      //服务器是否正在载入持久化文件
// 	RdbChangesSinceLastSave  int64  `gorm:"column:rdb_changes_since_last_save"`  //有多少个已经写入的命令还未被持久化
// 	RdbBgsaveInProgress      int64  `gorm:"column:rdb_bgsave_in_progress"`       //服务器是否正在创建rdb文件
// 	RdbLastSaveTime          int64  `gorm:"column:rdb_last_save_time"`           //已经有多长时间没有进行持久化了
// 	RdbLastBgsaveStatus      string `gorm:"column:rdb_last_bgsave_status"`       //最后一次的rdb持久化是否成功
// 	RdbLastBgsaveTimeSec     int64  `gorm:"column:rdb_last_bgsave_time_sec"`     //最后一次生成rdb文件耗时秒数
// 	RdbCurrentBgsaveTimeSec  int64  `gorm:"column:rdb_current_bgsave_time_sec"`  //如果服务器正在创建rdb文件，那么当前这个记录就是创建操作耗时秒数
// 	RdbLastCowSize           int64  `gorm:"column:rdb_last_cow_size"`            //
// 	AofEnabled               int64  `gorm:"column:aof_enabled"`                  //是否开启了aof
// 	AofRewriteInProgress     int64  `gorm:"column:aof_rewrite_in_progress"`      //标识aof的rewrite操作是否进行中
// 	AofRewriteScheduled      int64  `gorm:"column:aof_rewrite_scheduled"`        //
// 	AofLastRewriteTimeSec    int64  `gorm:"column:aof_last_rewrite_time_sec"`    //
// 	AofCurrentRewriteTimeSec int64  `gorm:"column:aof_current_rewrite_time_sec"` //
// 	AofLastBgrewriteStatus   string `gorm:"column:aof_last_bgrewrite_status"`    //上次bgrewriteaof操作的状态
// 	AofLastWriteStatus       string `gorm:"column:aof_last_write_status"`        //上一次aof写入状态
// 	AofLastCowSize           int64  `gorm:"column:aof_last_cow_size"`            //
// 	AofCurrentSize           int64  `gorm:"column:aof_current_size"`             //
// 	AofBaseSize              int64  `gorm:"column:aof_base_size"`                //
// 	AofPendingRewrite        int64  `gorm:"column:aof_pending_rewrite"`          //
// 	AofBufferLength          int64  `gorm:"column:aof_buffer_length"`            //
// 	AofRewriteBufferLength   int64  `gorm:"column:aof_rewrite_buffer_length"`    //
// 	AofPendingBioFsync       int64  `gorm:"column:aof_pending_bio_fsync"`        //
// 	AofDelayedFsync          int64  `gorm:"column:aof_delayed_fsync"`            //
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

type DaoRedisPersistence struct {
	Data *models.RedisPersistenceDB
	Mu   sync.Mutex
}

//提取查询数据
func (this *DaoRedisPersistence) exportdatas(rows *sql.Rows) (datas []*models.RedisPersistenceDB, err error) {
	for rows.Next() {
		data := &models.RedisPersistenceDB{}
		err := rows.Scan(&data.ID, &data.Hid, &data.GetTime, &data.Loading, &data.RdbChangesSinceLastSave, &data.RdbBgsaveInProgress, &data.RdbLastSaveTime,
			&data.RdbLastBgsaveStatus, &data.RdbLastBgsaveTimeSec, &data.RdbCurrentBgsaveTimeSec, &data.RdbLastCowSize,
			&data.AofEnabled, &data.AofRewriteInProgress, &data.AofRewriteScheduled, &data.AofLastRewriteTimeSec, &data.AofCurrentRewriteTimeSec,
			&data.AofLastBgrewriteStatus, &data.AofLastWriteStatus, &data.AofLastCowSize, &data.AofCurrentSize, &data.AofBaseSize,
			&data.AofPendingRewrite, &data.AofBufferLength, &data.AofRewriteBufferLength, &data.AofPendingBioFsync, &data.AofDelayedFsync)
		if err != nil {
			manlog.Error(err)
		}
		//manlog.Debug(*data)
		datas = append(datas, data)
	}
	return
}

//提取查询数据
func (this *DaoRedisPersistence) exportdatas1(rows *sql.Rows) (data *models.RedisPersistenceDB, err error) {
	for rows.Next() {
		data = &models.RedisPersistenceDB{}
		err := rows.Scan(&data.ID, &data.Hid, &data.GetTime, &data.Loading, &data.RdbChangesSinceLastSave, &data.RdbBgsaveInProgress, &data.RdbLastSaveTime,
			&data.RdbLastBgsaveStatus, &data.RdbLastBgsaveTimeSec, &data.RdbCurrentBgsaveTimeSec, &data.RdbLastCowSize,
			&data.AofEnabled, &data.AofRewriteInProgress, &data.AofRewriteScheduled, &data.AofLastRewriteTimeSec, &data.AofCurrentRewriteTimeSec,
			&data.AofLastBgrewriteStatus, &data.AofLastWriteStatus, &data.AofLastCowSize, &data.AofCurrentSize, &data.AofBaseSize,
			&data.AofPendingRewrite, &data.AofBufferLength, &data.AofRewriteBufferLength, &data.AofPendingBioFsync, &data.AofDelayedFsync)
		if err != nil {
			manlog.Error(err)
		}
	}
	return
}

func (this *DaoRedisPersistence) Create() error {
	db := sqlitedb.GetDBConn()
	this.Mu.Lock()
	defer this.Mu.Unlock()
	stmt, err := db.Prepare("INSERT INTO table_redis_servers_persistence (hid,get_time,loading,rdb_changes_since_last_save,rdb_bgsave_in_progress,rdb_last_save_time," +
		"rdb_last_bgsave_status,rdb_last_bgsave_time_sec,rdb_current_bgsave_time_sec,rdb_last_cow_size,aof_enabled,aof_rewrite_in_progress,aof_rewrite_scheduled," +
		"aof_last_rewrite_time_sec,aof_current_rewrite_time_sec,aof_last_bgrewrite_status,aof_last_write_status,aof_last_cow_size,aof_current_size,aof_base_size," +
		"aof_pending_rewrite,aof_buffer_length,aof_rewrite_buffer_length,aof_pending_bio_fsync,aof_delayed_fsync)" +
		" values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	defer stmt.Close()
	if err != nil {
		manlog.Error(err)
		return err
	}
	res, err := stmt.Exec(this.Data.Hid, this.Data.GetTime, this.Data.Loading, this.Data.RdbChangesSinceLastSave, this.Data.RdbBgsaveInProgress,
		this.Data.RdbLastSaveTime, this.Data.RdbLastBgsaveStatus, this.Data.RdbLastBgsaveTimeSec, this.Data.RdbCurrentBgsaveTimeSec, this.Data.RdbLastCowSize,
		this.Data.AofEnabled, this.Data.AofRewriteInProgress, this.Data.AofRewriteScheduled, this.Data.AofLastRewriteTimeSec, this.Data.AofCurrentRewriteTimeSec,
		this.Data.AofLastBgrewriteStatus, this.Data.AofLastWriteStatus, this.Data.AofLastCowSize, this.Data.AofCurrentSize, this.Data.AofBaseSize,
		this.Data.AofPendingRewrite, this.Data.AofBufferLength, this.Data.AofRewriteBufferLength, this.Data.AofPendingBioFsync, this.Data.AofDelayedFsync)
	manlog.Debug(&res)
	if err != nil {
		manlog.Error(err)
		return err
	}
	return nil
}

//查询当前时间之前的n条数据
func (this *DaoRedisPersistence) GetLastTimeData(rid string, last_time, n int64) (data []*models.RedisPersistenceDB, err error) {
	sql := fmt.Sprintf("SELECT * FROM table_redis_servers_persistence where hid='%s' and get_time<%d Order by id desc LIMIT %d;", rid, last_time, n)
	rows, err := sqlitedb.Query(sql)
	defer rows.Close()
	if err != nil {
		manlog.Error(err)
	}
	return this.exportdatas(rows)
}

//删除多久之前的数据
func (this *DaoRedisPersistence) DelTimeData(del_time int64) error {
	this.Mu.Lock()
	defer this.Mu.Unlock()
	sql := "delete from table_redis_servers_persistence where get_time<?"
	return sqlitedb.Del(sql, del_time)
}

//清空数据
func (this *DaoRedisPersistence) Empty() error {
	this.Mu.Lock()
	defer this.Mu.Unlock()
	sql := "delete from table_redis_servers_persistence"
	return sqlitedb.Del(sql)
}
