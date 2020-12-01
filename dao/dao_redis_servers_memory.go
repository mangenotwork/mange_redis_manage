package dao

// type RedisMemoryDB struct {
// 	ID                     int64   `gorm:"primary_key;column:id"`
// 	Hid                    string  `gorm:"column:hid"`
// 	GetTime                int64   `gorm:"column:get_time"`
// 	UsedMemory             int64   `gorm:"column:used_memory"`               //由redis分配器分配的内存总量，单位字节
// 	UsedMemoryHuman        string  `gorm:"column:used_memory_human"`         //
// 	UsedMemoryRss          int64   `gorm:"column:used_memory_rss"`           //从操作系统角度，返回redis已分配内存总量
// 	UsedMemoryRssHuman     string  `gorm:"column:used_memory_rss_human"`     //
// 	UsedMemoryPeak         int64   `gorm:"column:used_memory_peak"`          //redis的内存消耗峰值（以字节为单位）
// 	UsedMemoryPeakHuman    string  `gorm:"column:used_memory_peak_human"`    //
// 	UsedMemoryPeakPerc     string  `gorm:"column:used_memory_peak_perc"`     //
// 	UsedMemoryOverhead     int64   `gorm:"column:used_memory_overhead"`      //
// 	UsedMemoryStartup      int64   `gorm:"column:used_memory_startup"`       //
// 	UsedMemoryDataset      int64   `gorm:"column:used_memory_dataset"`       //
// 	UsedMemoryDatasetPerc  string  `gorm:"column:used_memory_dataset_perc"`  //
// 	AllocatorAllocated     int64   `gorm:"column:allocator_allocated"`       //
// 	AllocatorActive        int64   `gorm:"column:allocator_active"`          //
// 	AllocatorResident      int64   `gorm:"column:allocator_resident"`        //
// 	TotalSystemMemory      int64   `gorm:"column:total_system_memory"`       //
// 	TotalSystemMemoryHuman string  `gorm:"column:total_system_memory_human"` //
// 	UsedMemoryLua          int64   `gorm:"column:used_memory_lua"`           //lua引擎所使用的内存大小（单位字节）
// 	UsedMemoryLuaHuman     string  `gorm:"column:used_memory_lua_human"`     //
// 	UsedMemoryScripts      int64   `gorm:"column:used_memory_scripts"`       //
// 	UsedMemoryScriptsHuman string  `gorm:"column:used_memory_scripts_human"` //
// 	NumberOfCachedScripts  int64   `gorm:"column:number_of_cached_scripts"`  //
// 	Maxmemory              int64   `gorm:"column:maxmemory"`                 //
// 	MaxmemoryHuman         string  `gorm:"column:maxmemory_human"`           //
// 	MaxmemoryPolicy        string  `gorm:"column:maxmemory_policy"`          //
// 	AllocatorFragRatio     float64 `gorm:"column:allocator_frag_ratio"`      //
// 	AllocatorFragBytes     int64   `gorm:"column:allocator_frag_bytes"`      //
// 	AllocatorRssRatio      float64 `gorm:"column:allocator_rss_ratio"`       //
// 	AllocatorRssBytes      int64   `gorm:"column:allocator_rss_bytes"`       //
// 	RssOverheadRatio       float64 `gorm:"column:rss_overhead_ratio"`        //
// 	RssOverheadBytes       int64   `gorm:"column:rss_overhead_bytes"`        //
// 	MemFragmentationRatio  float64 `gorm:"column:mem_fragmentation_ratio"`   //used_memory_rss 和 used_memory 之间的比率
// 	MemFragmentationBytes  int64   `gorm:"column:mem_fragmentation_bytes"`   //
// 	MemNotCountedForEvict  int64   `gorm:"column:mem_not_counted_for_evict"` //
// 	MemReplicationBacklog  int64   `gorm:"column:mem_replication_backlog"`   //
// 	MemClientsSlaves       int64   `gorm:"column:mem_clients_slaves"`        //
// 	MemClientsNormal       int64   `gorm:"column:mem_clients_normal"`        //
// 	MemAofBuffer           int64   `gorm:"column:mem_aof_buffer"`            //
// 	MemAllocator           string  `gorm:"column:mem_allocator"`             //编译时指定的redis的内存分配器。越好的分配器内存碎片化率越低，低版本建议升级
// 	ActiveDefragRunning    int64   `gorm:"column:active_defrag_running"`     //
// 	LazyfreePendingObjects int64   `gorm:"column:lazyfree_pending_objects"`  //
// }

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"

	"github.com/mangenotwork/mange_redis_manage/common/manlog"
	"github.com/mangenotwork/mange_redis_manage/common/sqlitedb"
	"github.com/mangenotwork/mange_redis_manage/repository"
	_ "github.com/mattn/go-sqlite3"
)

type DaoRedisMemory struct {
	Data *models.RedisMemoryDB
	Mu   sync.Mutex
}

//提取查询数据
func (this *DaoRedisMemory) exportdatas(rows *sql.Rows) (datas []*models.RedisMemoryDB, err error) {
	for rows.Next() {
		data := &models.RedisMemoryDB{}
		err := rows.Scan(&data.ID, &data.Hid, &data.GetTime, &data.UsedMemory, &data.UsedMemoryHuman, &data.UsedMemoryRss, &data.UsedMemoryRssHuman,
			&data.UsedMemoryPeak, &data.UsedMemoryPeakHuman, &data.UsedMemoryPeakPerc, &data.UsedMemoryOverhead, &data.UsedMemoryStartup,
			&data.UsedMemoryDataset, &data.UsedMemoryDatasetPerc, &data.AllocatorAllocated, &data.AllocatorActive, &data.AllocatorResident,
			&data.TotalSystemMemory, &data.TotalSystemMemoryHuman, &data.UsedMemoryLua, &data.UsedMemoryLuaHuman, &data.UsedMemoryScripts,
			&data.UsedMemoryScriptsHuman, &data.NumberOfCachedScripts, &data.Maxmemory, &data.MaxmemoryHuman, &data.MaxmemoryPolicy, &data.AllocatorFragRatio,
			&data.AllocatorFragBytes, &data.AllocatorRssRatio, &data.AllocatorRssBytes, &data.RssOverheadRatio, &data.RssOverheadBytes,
			&data.MemFragmentationRatio, &data.MemFragmentationBytes, &data.MemNotCountedForEvict, &data.MemReplicationBacklog, &data.MemClientsSlaves,
			&data.MemClientsNormal, &data.MemAofBuffer, &data.MemAllocator, &data.ActiveDefragRunning, &data.LazyfreePendingObjects)
		if err != nil {
			manlog.Error(err)
		}
		//manlog.Debug(*data)
		datas = append(datas, data)
	}
	return
}

//提取查询数据
func (this *DaoRedisMemory) exportdatas1(rows *sql.Rows) (data *models.RedisMemoryDB, err error) {
	for rows.Next() {
		data = &models.RedisMemoryDB{}
		err := rows.Scan(&data.ID, &data.Hid, &data.GetTime, &data.UsedMemory, &data.UsedMemoryHuman, &data.UsedMemoryRss, &data.UsedMemoryRssHuman,
			&data.UsedMemoryPeak, &data.UsedMemoryPeakHuman, &data.UsedMemoryPeakPerc, &data.UsedMemoryOverhead, &data.UsedMemoryStartup,
			&data.UsedMemoryDataset, &data.UsedMemoryDatasetPerc, &data.AllocatorAllocated, &data.AllocatorActive, &data.AllocatorResident,
			&data.TotalSystemMemory, &data.TotalSystemMemoryHuman, &data.UsedMemoryLua, &data.UsedMemoryLuaHuman, &data.UsedMemoryScripts,
			&data.UsedMemoryScriptsHuman, &data.NumberOfCachedScripts, &data.Maxmemory, &data.MaxmemoryHuman, &data.MaxmemoryPolicy, &data.AllocatorFragRatio,
			&data.AllocatorFragBytes, &data.AllocatorRssRatio, &data.AllocatorRssBytes, &data.RssOverheadRatio, &data.RssOverheadBytes,
			&data.MemFragmentationRatio, &data.MemFragmentationBytes, &data.MemNotCountedForEvict, &data.MemReplicationBacklog, &data.MemClientsSlaves,
			&data.MemClientsNormal, &data.MemAofBuffer, &data.MemAllocator, &data.ActiveDefragRunning, &data.LazyfreePendingObjects)
		if err != nil {
			manlog.Error(err)
		}
	}
	return
}

//[]int64 拼接成 1,2,3  的字符串
func (this *DaoRedisMemory) intlistdata2str(data []int64) string {
	Str := ""
	for _, v := range data {
		Str = fmt.Sprintf("%s%d,", Str, v)
	}
	Str = strings.TrimRight(Str, ",")
	return Str
}

func (this *DaoRedisMemory) Create() error {
	db := sqlitedb.GetDBConn()
	this.Mu.Lock()
	defer this.Mu.Unlock()
	stmt, err := db.Prepare("INSERT INTO table_redis_servers_memory (hid,get_time,used_memory,used_memory_human,used_memory_rss,used_memory_rss_human,used_memory_peak," +
		"used_memory_peak_human,used_memory_peak_perc,used_memory_overhead,used_memory_startup,used_memory_dataset,used_memory_dataset_perc,allocator_allocated,allocator_active," +
		"allocator_resident,total_system_memory,total_system_memory_human,used_memory_lua,used_memory_lua_human,used_memory_scripts,used_memory_scripts_human," +
		"number_of_cached_scripts,maxmemory,maxmemory_human,maxmemory_policy,allocator_frag_ratio,allocator_frag_bytes,allocator_rss_ratio,allocator_rss_bytes," +
		"rss_overhead_ratio,rss_overhead_bytes,mem_fragmentation_ratio,mem_fragmentation_bytes,mem_not_counted_for_evict,mem_replication_backlog,mem_clients_slaves," +
		"mem_clients_normal,mem_aof_buffer,mem_allocator,active_defrag_running,lazyfree_pending_objects) " +
		"values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	defer stmt.Close()
	if err != nil {
		manlog.Error(err)
		return err
	}
	res, err := stmt.Exec(this.Data.Hid, this.Data.GetTime, this.Data.UsedMemory, this.Data.UsedMemoryHuman, this.Data.UsedMemoryRss, this.Data.UsedMemoryRssHuman,
		this.Data.UsedMemoryPeak, this.Data.UsedMemoryPeakHuman, this.Data.UsedMemoryPeakPerc, this.Data.UsedMemoryOverhead, this.Data.UsedMemoryStartup,
		this.Data.UsedMemoryDataset, this.Data.UsedMemoryDatasetPerc, this.Data.AllocatorAllocated, this.Data.AllocatorActive, this.Data.AllocatorResident,
		this.Data.TotalSystemMemory, this.Data.TotalSystemMemoryHuman, this.Data.UsedMemoryLua, this.Data.UsedMemoryLuaHuman, this.Data.UsedMemoryScripts,
		this.Data.UsedMemoryScriptsHuman, this.Data.NumberOfCachedScripts, this.Data.Maxmemory, this.Data.MaxmemoryHuman, this.Data.MaxmemoryPolicy, this.Data.AllocatorFragRatio,
		this.Data.AllocatorFragBytes, this.Data.AllocatorRssRatio, this.Data.AllocatorRssBytes, this.Data.RssOverheadRatio, this.Data.RssOverheadBytes,
		this.Data.MemFragmentationRatio, this.Data.MemFragmentationBytes, this.Data.MemNotCountedForEvict, this.Data.MemReplicationBacklog, this.Data.MemClientsSlaves,
		this.Data.MemClientsNormal, this.Data.MemAofBuffer, this.Data.MemAllocator, this.Data.ActiveDefragRunning, this.Data.LazyfreePendingObjects)
	manlog.Debug(&res)
	if err != nil {
		manlog.Error(err)
		return err
	}
	return nil
}

//获取最近存的数据是时间
func (this *DaoRedisMemory) NowDataTime(rid string) (int64, error) {
	sql := fmt.Sprintf("SELECT get_time FROM table_redis_servers_memory where hid='%s' Order by get_time desc Limit 1;", rid)
	manlog.Debug(sql)
	rows, err := sqlitedb.Query(sql)
	defer rows.Close()
	if err != nil {
		manlog.Error(err)
	}
	var timedata int64
	for rows.Next() {
		err := rows.Scan(&timedata)
		if err != nil {
			manlog.Error(err)
			return 0, err
		}
	}
	return timedata, nil
}

//获取最老存的数据是时间
func (this *DaoRedisMemory) OldDataTime(rid string) (int64, error) {
	sql := fmt.Sprintf("SELECT get_time FROM table_redis_servers_memory where hid='%s' Order by get_time asc Limit 1;", rid)
	manlog.Debug(sql)
	rows, err := sqlitedb.Query(sql)
	defer rows.Close()
	if err != nil {
		manlog.Error(err)
	}
	var timedata int64
	for rows.Next() {
		err := rows.Scan(&timedata)
		if err != nil {
			manlog.Error(err)
			return 0, err
		}
	}
	return timedata, nil
}

//获取某时间之后的所有get time
func (this *DaoRedisMemory) GetAllGetTime(rid string, agotime int64) ([]int64, error) {
	sql := fmt.Sprintf("SELECT get_time FROM table_redis_servers_memory where hid='%s' and get_time>%d Order by get_time asc;", rid, agotime)
	manlog.Debug(sql)
	rows, err := sqlitedb.Query(sql)
	defer rows.Close()
	if err != nil {
		manlog.Error(err)
	}
	var timedatas []int64
	var timedata int64
	for rows.Next() {
		err := rows.Scan(&timedata)
		if err != nil {
			manlog.Error(err)
			return nil, err
		}
		timedatas = append(timedatas, timedata)
	}
	return timedatas, nil
}

//获取给图表显示的数据
func (this *DaoRedisMemory) GetShowMemory(rid string, showtime []int64) (datas []*models.RedisMemoryDB, err error) {
	showtimeStr := this.intlistdata2str(showtime)
	sql := fmt.Sprintf("SELECT * FROM table_redis_servers_memory where hid='%s' and get_time in (%s) Order by get_time asc;", rid, showtimeStr)
	manlog.Debug(sql)
	rows, err := sqlitedb.Query(sql)
	defer rows.Close()
	if err != nil {
		manlog.Error(err)
	}
	return this.exportdatas(rows)
}

//查询当前时间之前的n条数据
func (this *DaoRedisMemory) GetLastTimeData(rid string, last_time, n int64) (data []*models.RedisMemoryDB, err error) {
	sql := fmt.Sprintf("SELECT * FROM table_redis_servers_memory where hid='%s' and get_time<%d Order by id desc LIMIT %d;", rid, last_time, n)
	rows, err := sqlitedb.Query(sql)
	defer rows.Close()
	if err != nil {
		manlog.Error(err)
	}
	return this.exportdatas(rows)
}

//获取最新的一条数据
func (this *DaoRedisMemory) GetNewData(rid string) (data *models.RedisMemoryDB, err error) {
	sql := fmt.Sprintf("SELECT * FROM table_redis_servers_memory where hid='%s' Order by get_time desc LIMIT 1;", rid)
	rows, err := sqlitedb.Query(sql)
	defer rows.Close()
	if err != nil {
		manlog.Error(err)
	}
	return this.exportdatas1(rows)
}

//删除多久之前的数据
func (this *DaoRedisMemory) DelTimeData(del_time int64) error {
	this.Mu.Lock()
	defer this.Mu.Unlock()
	sql := "delete from table_redis_servers_memory where get_time<?"
	return sqlitedb.Del(sql, del_time)
}

//清空数据
func (this *DaoRedisMemory) Empty() error {
	this.Mu.Lock()
	defer this.Mu.Unlock()
	sql := "delete from table_redis_servers_memory"
	return sqlitedb.Del(sql)
}
