package dao

// type RedisStatsDB struct {
// 	ID                         int64   `gorm:"primary_key;column:id"`
// 	Hid                        string  `gorm:"column:hid"`
// 	GetTime                    int64   `gorm:"column:get_time"`
// 	TotalConnectionsReceived   int64   `gorm:"column:total_connections_received"`     //新创建的链接个数，如果过多，会影响性能
// 	TotalCommandsProcessed     int64   `gorm:"column:total_commands_processed"`       //redis处理的命令数
// 	InstantaneousOpsPerSec     int64   `gorm:"column:instantaneous_ops_per_sec"`      //redis当前的qps，redis内部较实时的每秒执行命令数
// 	TotalNetInputBytes         int64   `gorm:"column:total_net_input_bytes"`          //redis网络入口流量字节数
// 	TotalNetOutputBytes        int64   `gorm:"column:total_net_output_bytes"`         //redis网络出口流量字节数
// 	InstantaneousInputKbps     float64 `gorm:"column:instantaneous_input_kbps"`       //redis网络入口kps
// 	InstantaneousOutputKbps    float64 `gorm:"column:instantaneous_output_kbps"`      //redis网络出口kps
// 	RejectedConnections        int64   `gorm:"column:rejected_connections"`           //拒绝的连接个数，redis连接个数已经达到maxclients限制。
// 	SyncFull                   int64   `gorm:"column:sync_full"`                      //主从完全同步成功次数
// 	SyncPartialOk              int64   `gorm:"column:sync_partial_ok"`                //主从部分同步成功次数
// 	SyncPartialErr             int64   `gorm:"column:sync_partial_err"`               //主从部分同步失败次数
// 	ExpiredKeys                int64   `gorm:"column:expired_keys"`                   //运行以来过期的key的数量
// 	ExpiredStalePerc           string  `gorm:"column:expired_stale_perc"`             //
// 	ExpiredTimeCapReachedCount int64   `gorm:"column:expired_time_cap_reached_count"` //
// 	EvictedKeys                int64   `gorm:"column:evicted_keys"`                   //运行以来剔除（超过maxmemory）的key的数量s
// 	KeyspaceHits               int64   `gorm:"column:keyspace_hits"`                  //命中次数
// 	KeyspaceMisses             int64   `gorm:"column:keyspace_misses"`                //没命中次数
// 	PubsubChannels             int64   `gorm:"column:pubsub_channels"`                //当前使用中的频道数量
// 	PubsubPatterns             int64   `gorm:"column:pubsub_patterns"`                //当前使用的模式数量
// 	LatestForkUsec             int64   `gorm:"column:latest_fork_usec"`               //
// 	MigrateCachedSockets       int64   `gorm:"column:migrate_cached_sockets"`         //
// 	SlaveExpiresTrackedKeys    int64   `gorm:"column:slave_expires_tracked_keys"`     //
// 	ActiveDefragHits           int64   `gorm:"column:active_defrag_hits"`             //
// 	ActiveDefragMisses         int64   `gorm:"column:active_defrag_misses"`           //
// 	ActiveDefragKeyHits        int64   `gorm:"column:active_defrag_key_hits"`         //
// 	ActiveDefragKeyMisses      int64   `gorm:"column:active_defrag_key_misses"`       //
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

type DaoRedisStats struct {
	Data *models.RedisStatsDB
	Mu   sync.Mutex
}

//提取查询数据
func (this *DaoRedisStats) exportdatas(rows *sql.Rows) (datas []*models.RedisStatsDB, err error) {
	for rows.Next() {
		data := &models.RedisStatsDB{}
		err := rows.Scan(&data.ID, &data.Hid, &data.GetTime, &data.TotalConnectionsReceived, &data.TotalCommandsProcessed, &data.InstantaneousOpsPerSec,
			&data.TotalNetInputBytes, &data.TotalNetOutputBytes, &data.InstantaneousInputKbps, &data.InstantaneousOutputKbps, &data.RejectedConnections,
			&data.SyncFull, &data.SyncPartialOk, &data.SyncPartialErr, &data.ExpiredKeys, &data.ExpiredStalePerc, &data.ExpiredTimeCapReachedCount,
			&data.EvictedKeys, &data.KeyspaceHits, &data.KeyspaceMisses, &data.PubsubChannels, &data.PubsubPatterns, &data.LatestForkUsec,
			&data.MigrateCachedSockets, &data.SlaveExpiresTrackedKeys, &data.ActiveDefragHits, &data.ActiveDefragMisses, &data.ActiveDefragKeyHits,
			&data.ActiveDefragKeyMisses)
		if err != nil {
			manlog.Error(err)
		}
		//manlog.Debug(*data)
		datas = append(datas, data)
	}
	return
}

//提取查询数据
func (this *DaoRedisStats) exportdatas1(rows *sql.Rows) (data *models.RedisStatsDB, err error) {
	for rows.Next() {
		data = &models.RedisStatsDB{}
		err := rows.Scan(&data.ID, &data.Hid, &data.GetTime, &data.TotalConnectionsReceived, &data.TotalCommandsProcessed, &data.InstantaneousOpsPerSec,
			&data.TotalNetInputBytes, &data.TotalNetOutputBytes, &data.InstantaneousInputKbps, &data.InstantaneousOutputKbps, &data.RejectedConnections,
			&data.SyncFull, &data.SyncPartialOk, &data.SyncPartialErr, &data.ExpiredKeys, &data.ExpiredStalePerc, &data.ExpiredTimeCapReachedCount,
			&data.EvictedKeys, &data.KeyspaceHits, &data.KeyspaceMisses, &data.PubsubChannels, &data.PubsubPatterns, &data.LatestForkUsec,
			&data.MigrateCachedSockets, &data.SlaveExpiresTrackedKeys, &data.ActiveDefragHits, &data.ActiveDefragMisses, &data.ActiveDefragKeyHits,
			&data.ActiveDefragKeyMisses)
		if err != nil {
			manlog.Error(err)
		}
	}
	return
}

func (this *DaoRedisStats) Create() error {
	db := sqlitedb.GetDBConn()
	this.Mu.Lock()
	defer this.Mu.Unlock()
	stmt, err := db.Prepare("INSERT INTO table_redis_servers_stats (hid,get_time,total_connections_received,total_commands_processed,instantaneous_ops_per_sec," +
		"total_net_input_bytes,total_net_output_bytes,instantaneous_input_kbps,instantaneous_output_kbps,rejected_connections,sync_full,sync_partial_ok,sync_partial_err," +
		"expired_keys,expired_stale_perc,expired_time_cap_reached_count,evicted_keys,keyspace_hits,keyspace_misses,pubsub_channels,pubsub_patterns,latest_fork_usec," +
		"migrate_cached_sockets,slave_expires_tracked_keys,active_defrag_hits,active_defrag_misses,active_defrag_key_hits,active_defrag_key_misses)" +
		" values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	defer stmt.Close()
	if err != nil {
		manlog.Error(err)
		return err
	}
	res, err := stmt.Exec(this.Data.Hid, this.Data.GetTime, this.Data.TotalConnectionsReceived, this.Data.TotalCommandsProcessed, this.Data.InstantaneousOpsPerSec,
		this.Data.TotalNetInputBytes, this.Data.TotalNetOutputBytes, this.Data.InstantaneousInputKbps, this.Data.InstantaneousOutputKbps, this.Data.RejectedConnections,
		this.Data.SyncFull, this.Data.SyncPartialOk, this.Data.SyncPartialErr, this.Data.ExpiredKeys, this.Data.ExpiredStalePerc, this.Data.ExpiredTimeCapReachedCount,
		this.Data.EvictedKeys, this.Data.KeyspaceHits, this.Data.KeyspaceMisses, this.Data.PubsubChannels, this.Data.PubsubPatterns, this.Data.LatestForkUsec,
		this.Data.MigrateCachedSockets, this.Data.SlaveExpiresTrackedKeys, this.Data.ActiveDefragHits, this.Data.ActiveDefragMisses, this.Data.ActiveDefragKeyHits,
		this.Data.ActiveDefragKeyMisses)
	manlog.Debug(&res)
	if err != nil {
		manlog.Error(err)
		return err
	}
	return nil
}

//获取最新的一条数据
func (this *DaoRedisStats) GetNewData(rid string) (data *models.RedisStatsDB, err error) {

	sql := fmt.Sprintf("SELECT * FROM table_redis_servers_stats where hid='%s' Order by get_time desc LIMIT 1;", rid)
	rows, err := sqlitedb.Query(sql)
	defer rows.Close()
	if err != nil {
		manlog.Error(err)
	}
	return this.exportdatas1(rows)
}

//查询当前时间之前的n条数据
func (this *DaoRedisStats) GetLastTimeData(rid string, last_time, n int64) (data []*models.RedisStatsDB, err error) {
	sql := fmt.Sprintf("SELECT * FROM table_redis_servers_stats where hid='%s' and get_time<%d Order by id desc LIMIT %d;", rid, last_time, n)
	rows, err := sqlitedb.Query(sql)
	defer rows.Close()
	if err != nil {
		manlog.Error(err)
	}
	return this.exportdatas(rows)
}

//删除多久之前的数据
func (this *DaoRedisStats) DelTimeData(del_time int64) error {
	this.Mu.Lock()
	defer this.Mu.Unlock()
	sql := "delete from table_redis_servers_stats where get_time<?"
	return sqlitedb.Del(sql, del_time)
}

//清空数据
func (this *DaoRedisStats) Empty() error {
	this.Mu.Lock()
	defer this.Mu.Unlock()
	sql := "delete from table_redis_servers_stats"
	return sqlitedb.Del(sql)
}
