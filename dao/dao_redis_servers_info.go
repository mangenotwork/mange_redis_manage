package dao

// type RedisServerInfosDB struct {
// 	ID              int64  `gorm:"primary_key;column:id"`
// 	Hid             string `gorm:"column:hid"`
// 	GetTime         int64  `gorm:"column:get_time"`
// 	RedisVersion    string `gorm:"column:redis_version"`     //Redis 服务器版本
// 	RedisGitSha1    string `gorm:"column:redis_git_sha1"`    //Git SHA1
// 	RedisGitDirty   string `gorm:"column:redis_git_dirty"`   //Git dirty flag
// 	RedisBuildId    string `gorm:"column:redis_build_id"`    //Redis build id
// 	RedisMode       string `gorm:"column:redis_mode"`        //运行模式，单机或集群
// 	Os              string `gorm:"column:os"`                //Redis 服务器的宿主操作系统
// 	ArchBits        string `gorm:"column:arch_bits"`         //架构64位
// 	MultiplexingApi string `gorm:"column:multiplexing_api"`  //redis所使用的事件处理模型
// 	AtomicvarApi    string `gorm:"column:atomicvar_api"`     //
// 	GccVersion      string `gorm:"column:gcc_version"`       //编译redis时gcc版本
// 	ProcessId       string `gorm:"column:process_id"`        //redis服务器进程的pid
// 	RunId           string `gorm:"column:run_id"`            //redis服务器的随机标识符（sentinel和集群）
// 	TcpPort         int64  `gorm:"column:tcp_port"`          //
// 	UptimeInSeconds int64  `gorm:"column:uptime_in_seconds"` //redis服务器启动总时间，单位秒
// 	UptimeInDays    int64  `gorm:"column:uptime_in_days"`    //redis服务器启动总时间，单位天
// 	Hz              int64  `gorm:"column:hz"`                //redis内部调度频率（关闭timeout客户端，删除过期key）
// 	ConfiguredHz    int64  `gorm:"column:configured_hz"`     //
// 	Lru_clock       int64  `gorm:"column:lru_clock"`         //自增时间，用于LRU管理
// 	Executable      string `gorm:"column:executable"`        //
// 	ConfigFile      string `gorm:"column:config_file"`       //配置文件路径
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

type DaoRedisServerInfos struct {
	Data *models.RedisServerInfosDB
	Mu   sync.Mutex
}

//提取查询数据
func (this *DaoRedisServerInfos) exportdatas(rows *sql.Rows) (datas []*models.RedisServerInfosDB, err error) {
	for rows.Next() {
		data := &models.RedisServerInfosDB{}
		err := rows.Scan(&data.ID, &data.Hid, &data.GetTime, &data.RedisVersion, &data.RedisGitSha1, &data.RedisGitDirty, &data.RedisBuildId,
			&data.RedisMode, &data.Os, &data.ArchBits, &data.MultiplexingApi, &data.AtomicvarApi, &data.GccVersion, &data.ProcessId,
			&data.RunId, &data.TcpPort, &data.UptimeInSeconds, &data.UptimeInDays, &data.Hz, &data.ConfiguredHz, &data.Lru_clock,
			&data.Executable, &data.ConfigFile)
		if err != nil {
			manlog.Error(err)
		}
		//manlog.Debug(*data)
		datas = append(datas, data)
	}
	return
}

//提取查询数据
func (this *DaoRedisServerInfos) exportdatas1(rows *sql.Rows) (data *models.RedisServerInfosDB, err error) {
	for rows.Next() {
		data = &models.RedisServerInfosDB{}
		err := rows.Scan(&data.ID, &data.Hid, &data.GetTime, &data.RedisVersion, &data.RedisGitSha1, &data.RedisGitDirty, &data.RedisBuildId,
			&data.RedisMode, &data.Os, &data.ArchBits, &data.MultiplexingApi, &data.AtomicvarApi, &data.GccVersion, &data.ProcessId,
			&data.RunId, &data.TcpPort, &data.UptimeInSeconds, &data.UptimeInDays, &data.Hz, &data.ConfiguredHz, &data.Lru_clock,
			&data.Executable, &data.ConfigFile)
		if err != nil {
			manlog.Error(err)
		}
	}
	return
}

func (this *DaoRedisServerInfos) Create() error {
	db := sqlitedb.GetDBConn()
	this.Mu.Lock()
	defer this.Mu.Unlock()
	stmt, err := db.Prepare("INSERT INTO table_redis_servers_infos (hid,get_time,redis_version,redis_git_sha1,redis_git_dirty,redis_build_id,redis_mode," +
		"os,arch_bits,multiplexing_api,atomicvar_api,gcc_version,process_id,run_id,tcp_port,uptime_in_seconds,uptime_in_days,hz,configured_hz,lru_clock," +
		"executable,config_file) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	defer stmt.Close()
	if err != nil {
		manlog.Error(err)
		return err
	}
	res, err := stmt.Exec(this.Data.Hid, this.Data.GetTime, this.Data.RedisVersion, this.Data.RedisGitSha1, this.Data.RedisGitDirty, this.Data.RedisBuildId,
		this.Data.RedisMode, this.Data.Os, this.Data.ArchBits, this.Data.MultiplexingApi, this.Data.AtomicvarApi, this.Data.GccVersion, this.Data.ProcessId,
		this.Data.RunId, this.Data.TcpPort, this.Data.UptimeInSeconds, this.Data.UptimeInDays, this.Data.Hz, this.Data.ConfiguredHz, this.Data.Lru_clock,
		this.Data.Executable, this.Data.ConfigFile)
	manlog.Debug(&res)
	if err != nil {
		manlog.Error(err)
		return err
	}
	return nil
}

//获取最新的一条数据
func (this *DaoRedisServerInfos) GetNewData(rid string) (data *models.RedisServerInfosDB, err error) {

	sql := fmt.Sprintf("SELECT * FROM table_redis_servers_infos where hid='%s' Order by get_time desc LIMIT 1;", rid)
	rows, err := sqlitedb.Query(sql)
	defer rows.Close()
	if err != nil {
		manlog.Error(err)
	}
	return this.exportdatas1(rows)
}

//查询当前时间之前的n条数据
func (this *DaoRedisServerInfos) GetLastTimeData(rid string, last_time, n int64) (data []*models.RedisServerInfosDB, err error) {
	sql := fmt.Sprintf("SELECT * FROM table_redis_servers_infos where hid='%s' and get_time<%d Order by id desc LIMIT %d;", rid, last_time, n)
	rows, err := sqlitedb.Query(sql)
	defer rows.Close()
	if err != nil {
		manlog.Error(err)
	}
	return this.exportdatas(rows)
}

//删除多久之前的数据
func (this *DaoRedisServerInfos) DelTimeData(del_time int64) error {
	this.Mu.Lock()
	defer this.Mu.Unlock()
	sql := "delete from table_redis_servers_infos where get_time<?"
	return sqlitedb.Del(sql, del_time)
}

//清空数据
func (this *DaoRedisServerInfos) Empty() error {
	this.Mu.Lock()
	defer this.Mu.Unlock()
	sql := "delete from table_redis_servers_infos"
	return sqlitedb.Del(sql)
}
