//
//	包含redis服务器信息
//
package redis

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/mangenotwork/mange_redis_manage/common"
	"github.com/mangenotwork/mange_redis_manage/common/manlog"
	"github.com/mangenotwork/mange_redis_manage/repository"
)

//获取redis 服务器的详细信息
func GetRedisServersInfos(c redis.Conn, rcid string) (data *RedisServersInfo) {
	manlog.Debug("[Execute redis command]: ", "INFO")
	res, err := redis.String(c.Do("INFO"))
	if err != nil {
		manlog.Error(err)
		return
	}
	//manlog.Debug(res)

	data = &RedisServersInfo{}

	data.NowTime = time.Now().Unix()

	wg := new(sync.WaitGroup)

	wg.Add(1)
	go func() {
		defer wg.Done()
		data.Server = GetRedisServersInfo(res, rcid, data.NowTime)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		data.Clients = GetRedisClients(res, rcid, data.NowTime)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		data.Memory = GetRedisMemory(res, rcid, data.NowTime)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		data.Persistence = GetRedisPersistence(res, rcid, data.NowTime)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		data.Stats = GetRedisStats(res, rcid, data.NowTime)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		data.Replication = GetRedisReplication(res, rcid, data.NowTime)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		data.CPU = GetRedisCPU(res, rcid, data.NowTime)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		data.Cluster = GetRedisCluster(res, rcid, data.NowTime)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		data.Keyspace = GetRedisKeyspace(res, rcid, data.NowTime)
	}()

	wg.Wait()

	//manlog.Debug(*data)

	return
}

//获取redis 服务的时间
func GetRedisServersTime(c redis.Conn) {
	manlog.Debug("[Execute redis command]: ", "TIME")
	res, err := redis.String(c.Do("TIME"))
	if err != nil {
		manlog.Error(err)
		return
	}
	manlog.Debug(res)
}

//获取redis 服务的所有配置信息
func GetRedisServersAllConfig(c redis.Conn) (res []string, err error) {
	manlog.Debug("[Execute redis command]: ", "CONFIG GET *")
	res, err = redis.Strings(c.Do("CONFIG", "GET", "*"))
	if err != nil {
		manlog.Error(err)
		return
	}
	manlog.Debug(res)
	return
}

//获取redis 服务的指定配置信息
func GetRedisServersConfig(c redis.Conn, parameter string) (res []string, err error) {
	manlog.Debug("[Execute redis command]: ", "CONFIG GET ", "parameter")
	res, err = redis.Strings(c.Do("CONFIG", "GET", parameter))
	if err != nil {
		manlog.Error(err)
		return
	}
	manlog.Debug(res)
	return
}

//ping
func Ping(c redis.Conn) bool {
	manlog.Debug("[Execute redis command]: ", "PING")
	res, err := redis.String(c.Do("PING"))
	if err != nil {
		manlog.Error(err)
		return false
	}
	manlog.Debug(res)
	if res == "PONG" {
		return true
	}
	return false
}

//查看 show log : SLOWLOG GET
/*
1) 1) (integer) 12                      # 唯一性(unique)的日志标识符
   2) (integer) 1324097834              # 被记录命令的执行时间点，以 UNIX 时间戳格式表示
   3) (integer) 16                      # 查询执行时间，以微秒为单位
   4) 1) "CONFIG"                       # 执行的命令，以数组的形式排列
      2) "GET"                          # 这里完整的命令是 CONFIG GET slowlog-log-slower-than
      3) "slowlog-log-slower-than"
   5)   "192.168.0.101:61819"		#客户端
   6)   ""
*/
func GetSlowlog(c redis.Conn, rcid string) (datas []*models.RedisSlowLog) {
	manlog.Debug("[Execute redis command]: ", "SLOWLOG GET")
	res, err := redis.Values(c.Do("SLOWLOG", "GET"))
	if err != nil {
		manlog.Error(err)
	}
	manlog.Debug(res)
	datas = make([]*models.RedisSlowLog, 0)
	git_time := time.Now().Unix()
	for _, v := range res {
		fmt.Printf("%T", v)
		v_list := v.([]interface{})
		if len(v_list) == 6 {
			id := v_list[0].(int64)
			manlog.Debug("唯一性ID ： ", id)

			t := v_list[1].(int64)
			manlog.Debug("被记录命令 ： ", t)

			run_time := v_list[2].(int64)
			manlog.Debug("执行时间 ： ", run_time)

			cmdlist := make([]string, 0)
			for _, m := range v_list[3].([]interface{}) {
				cmdlist = append(cmdlist, common.Uint82Str(m.([]uint8)))
			}
			cmd := strings.Join(cmdlist, " ")
			manlog.Debug("命令 ： ", cmd)

			client := common.Uint82Str(v_list[4].([]uint8))
			fmt.Printf("%T", client)
			manlog.Debug("客户端 ： ", client)

			datas = append(datas, &models.RedisSlowLog{
				Hid:      rcid,
				GetTime:  git_time,
				OnlyId:   id,
				Time:     t,
				Duration: run_time,
				Cmd:      cmd,
				Client:   client,
			})
		}
	}
	return
}

//查看当前日志数量  ： SLOWLOG LEN

//清空日志： SLOWLOG RESET

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
//
//停止所有客户端
//执行 SHUTDOWN SAVE 会强制让数据库执行保存操作，即使没有设定(configure)保存点
//执行 SHUTDOWN NOSAVE 会阻止数据库执行保存操作，即使已经设定有一个或多个保存点(你可以将这一用法看作是强制停止服务器的一个假想的 ABORT 命令)
//
//

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
//
//SAVE 命令执行一个同步保存操作，将当前 Redis 实例的所有数据快照(snapshot)以 RDB 文件的形式保存到硬盘。
//
//

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
//
//MONITOR
//实时打印出 Redis 服务器接收到的命令，调试用。
//
//

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
//
//LASTSAVE
//返回最近一次 Redis 成功将数据保存到磁盘上的时间，以 UNIX 时间戳格式表示。
//
//

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
//
//FLUSHDB
//清空当前数据库中的所有 key。
//
//

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
//
//FLUSHALL
//清空整个 Redis 服务器的数据(删除所有数据库的所有 key )。
//
//

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
//
//DEBUG SEGFAULT
//执行一个不合法的内存访问从而让 Redis 崩溃，仅在开发时用于 BUG 模拟。
//
//

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
//
//DBSIZE
//返回当前db数据库的 key 的数量。
//
//

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
//
//CONFIG SET 命令可以动态地调整 Redis 服务器的配置(configuration)而无须重启。
//CONFIG SET 可以修改的配置参数可以使用命令 CONFIG GET * 来列出，所有被 CONFIG SET 修改的配置参数都会立即生效。

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
//
//CONFIG REWRITE
//CONFIG REWRITE 命令对启动 Redis 服务器时所指定的 redis.conf 文件进行改写： 因为 CONFIG SET 命令可以对服务器的当前配置进行修改，
// 而修改后的配置可能和 redis.conf 文件中所描述的配置不一样， CONFIG REWRITE 的作用就是通过尽可能少的修改，
//将服务器当前所使用的配置记录到 redis.conf 文件中。
// 127.0.0.1:6379> CONFIG GET appendonly           # appendonly 处于关闭状态
// 1) "appendonly"
// 2) "no"
// 127.0.0.1:6379> CONFIG SET appendonly yes       # 打开 appendonly
// OK
// 127.0.0.1:6379> CONFIG GET appendonly
// 1) "appendonly"
// 2) "yes"
// 127.0.0.1:6379> CONFIG REWRITE                  # 将 appendonly 的修改写入到 redis.conf 中
// OK
//
//

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
//
//CONFIG RESETSTAT
// 重置 INFO 命令中的某些统计数据，包括：
// Keyspace hits (键空间命中次数)
// Keyspace misses (键空间不命中次数)
// Number of commands processed (执行命令的次数)
// Number of connections received (连接服务器的次数)
// Number of expired keys (过期key的数量)
// Number of rejected connections (被拒绝的连接数量)
// Latest fork(2) time(最后执行 fork(2) 的时间)
// The aof_delayed_fsync counter(aof_delayed_fsync 计数器的值)
//
//

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
//
// 在后台异步(Asynchronously)保存当前数据库的数据到磁盘。
// BGSAVE 命令执行之后立即返回 OK ，然后 Redis fork 出一个新子进程，原来的 Redis 进程(父进程)继续处理客户端请求，而子进程则负责将数据保存到磁盘，然后退出。
// 客户端可以通过 LASTSAVE 命令查看相关信息，判断 BGSAVE 命令是否执行成功。
//
//

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
//
type RedisServersInfo struct {
	NowTime     int64                      `json:"now_time"`         //当前时间
	Server      *models.RedisServerInfosDB `json:"server_data"`      //一般 Redis 服务器信息
	Clients     *models.RedisClientsDB     `json:"clients_data"`     //已连接客户端信息
	Memory      *models.RedisMemoryDB      `json:"memory_data"`      //内存信息
	Persistence *models.RedisPersistenceDB `json:"persistence_data"` //RDB 和 AOF 的相关信息
	Stats       *models.RedisStatsDB       `json:"stats_data"`       //一般统计信息
	Replication *models.RedisReplicationDB `json:"replication_data"` //主/从复制信息
	CPU         *models.RedisCPUDB         `json:"cpu_data"`         //CPU 计算量统计信息
	Cluster     *models.RedisClusterDB     `json:"cluster_data"`     //Redis 集群信息
	Keyspace    []*models.RedisKeyspaceDB  `json:"keyspace_data"`    //数据库相关的统计信息
	//Commandstats                   // Redis 命令统计信息
}

func GetRedisServersInfo(strs, rcid string, get_time int64) (data *models.RedisServerInfosDB) {
	redis_version := GetStrValue(strs, `redis_version:(.*?)\r\n`)
	redis_git_sha1 := GetStrValue(strs, `redis_git_sha1:(.*?)\r\n`)
	redis_git_dirty := GetStrValue(strs, `redis_git_dirty:(.*?)\r\n`)
	redis_build_id := GetStrValue(strs, `redis_build_id:(.*?)\r\n`)
	redis_mode := GetStrValue(strs, `redis_mode:(.*?)\r\n`)
	os := GetStrValue(strs, `os:(.*?)\r\n`)
	arch_bits := GetStrValue(strs, `arch_bits:(.*?)\r\n`)
	multiplexing_api := GetStrValue(strs, `multiplexing_api:(.*?)\r\n`)
	atomicvar_api := GetStrValue(strs, `atomicvar_api:(.*?)\r\n`)
	gcc_version := GetStrValue(strs, `gcc_version:(.*?)\r\n`)
	process_id := GetStrValue(strs, `process_id:(.*?)\r\n`)
	run_id := GetStrValue(strs, `run_id:(.*?)\r\n`)
	tcp_port := GetIntValue(strs, `tcp_port:(.*?)\r\n`)
	uptime_in_seconds := GetIntValue(strs, `uptime_in_seconds:(.*?)\r\n`)
	uptime_in_days := GetIntValue(strs, `uptime_in_days:(.*?)\r\n`)
	hz := GetIntValue(strs, `hz:(.*?)\r\n`)
	configured_hz := GetIntValue(strs, `configured_hz:(.*?)\r\n`)
	lru_clock := GetIntValue(strs, `lru_clock:(.*?)\r\n`)
	executable := GetStrValue(strs, `executable:(.*?)\r\n`)
	config_file := GetStrValue(strs, `config_file:(.*?)\r\n`)

	data = &models.RedisServerInfosDB{
		Hid:             rcid,
		GetTime:         get_time,
		RedisVersion:    redis_version,
		RedisGitSha1:    redis_git_sha1,
		RedisGitDirty:   redis_git_dirty,
		RedisBuildId:    redis_build_id,
		RedisMode:       redis_mode,
		Os:              os,
		ArchBits:        arch_bits,
		MultiplexingApi: multiplexing_api,
		AtomicvarApi:    atomicvar_api,
		GccVersion:      gcc_version,
		ProcessId:       process_id,
		RunId:           run_id,
		TcpPort:         tcp_port,
		UptimeInSeconds: uptime_in_seconds,
		UptimeInDays:    uptime_in_days,
		Hz:              hz,
		ConfiguredHz:    configured_hz,
		Lru_clock:       lru_clock,
		Executable:      executable,
		ConfigFile:      config_file,
	}
	manlog.Debug(data)
	manlog.Debug(*data)
	//err = nil
	return
}

func GetRedisClients(strs, rcid string, get_time int64) (data *models.RedisClientsDB) {
	connected_clients := GetIntValue(strs, `connected_clients:(.*?)\r\n`)
	client_recent_max_input_buffer := GetIntValue(strs, `client_recent_max_input_buffer:(.*?)\r\n`)
	client_recent_max_output_buffer := GetIntValue(strs, `client_recent_max_output_buffer:(.*?)\r\n`)
	blocked_clients := GetIntValue(strs, `blocked_clients:(.*?)\r\n`)

	data = &models.RedisClientsDB{
		Hid:                         rcid,
		GetTime:                     get_time,
		ConnectedClients:            connected_clients,
		ClientRecentMaxInputBuffer:  client_recent_max_input_buffer,
		ClientRecentMaxOutputBuffer: client_recent_max_output_buffer,
		BlockedClients:              blocked_clients,
	}

	manlog.Debug(data)
	manlog.Debug(*data)
	//err = nil
	return
}

func GetRedisMemory(strs, rcid string, get_time int64) (data *models.RedisMemoryDB) {
	used_memory := GetIntValue(strs, `used_memory:(.*?)\r\n`)
	used_memory_human := GetStrValue(strs, `used_memory_human:(.*?)\r\n`)
	used_memory_rss := GetIntValue(strs, `used_memory_rss:(.*?)\r\n`)
	used_memory_rss_human := GetStrValue(strs, `used_memory_rss_human:(.*?)\r\n`)
	used_memory_peak := GetIntValue(strs, `used_memory_peak:(.*?)\r\n`)
	used_memory_peak_human := GetStrValue(strs, `used_memory_peak_human:(.*?)\r\n`)
	used_memory_peak_perc := GetStrValue(strs, `used_memory_peak_perc:(.*?)\r\n`)
	used_memory_overhead := GetIntValue(strs, `used_memory_overhead:(.*?)\r\n`)
	used_memory_startup := GetIntValue(strs, `used_memory_startup:(.*?)\r\n`)
	used_memory_dataset := GetIntValue(strs, `used_memory_dataset:(.*?)\r\n`)
	used_memory_dataset_perc := GetStrValue(strs, `used_memory_dataset_perc:(.*?)\r\n`)
	allocator_allocated := GetIntValue(strs, `allocator_allocated:(.*?)\r\n`)
	allocator_active := GetIntValue(strs, `allocator_active:(.*?)\r\n`)
	allocator_resident := GetIntValue(strs, `allocator_resident:(.*?)\r\n`)
	total_system_memory := GetIntValue(strs, `total_system_memory:(.*?)\r\n`)
	total_system_memory_human := GetStrValue(strs, `total_system_memory_human:(.*?)\r\n`)
	used_memory_lua := GetIntValue(strs, `used_memory_lua:(.*?)\r\n`)
	used_memory_lua_human := GetStrValue(strs, `used_memory_lua_human:(.*?)\r\n`)
	used_memory_scripts := GetIntValue(strs, `used_memory_scripts:(.*?)\r\n`)
	used_memory_scripts_human := GetStrValue(strs, `used_memory_scripts_human:(.*?)\r\n`)
	number_of_cached_scripts := GetIntValue(strs, `number_of_cached_scripts:(.*?)\r\n`)
	maxmemory := GetIntValue(strs, `maxmemory:(.*?)\r\n`)
	maxmemory_human := GetStrValue(strs, `maxmemory_human:(.*?)\r\n`)
	maxmemory_policy := GetStrValue(strs, `maxmemory_policy:(.*?)\r\n`)
	allocator_frag_ratio := GetFloatValue(strs, `allocator_frag_ratio:(.*?)\r\n`)
	allocator_frag_bytes := GetIntValue(strs, `allocator_frag_bytes:(.*?)\r\n`)
	allocator_rss_ratio := GetFloatValue(strs, `allocator_rss_ratio:(.*?)\r\n`)
	allocator_rss_bytes := GetIntValue(strs, `allocator_rss_bytes:(.*?)\r\n`)
	rss_overhead_ratio := GetFloatValue(strs, `rss_overhead_ratio:(.*?)\r\n`)
	rss_overhead_bytes := GetIntValue(strs, `rss_overhead_bytes:(.*?)\r\n`)
	mem_fragmentation_ratio := GetFloatValue(strs, `mem_fragmentation_ratio:(.*?)\r\n`)
	mem_fragmentation_bytes := GetIntValue(strs, `mem_fragmentation_bytes:(.*?)\r\n`)
	mem_not_counted_for_evict := GetIntValue(strs, `mem_not_counted_for_evict:(.*?)\r\n`)
	mem_replication_backlog := GetIntValue(strs, `mem_replication_backlog:(.*?)\r\n`)
	mem_clients_slaves := GetIntValue(strs, `mem_clients_slaves:(.*?)\r\n`)
	mem_clients_normal := GetIntValue(strs, `mem_clients_normal:(.*?)\r\n`)
	mem_aof_buffer := GetIntValue(strs, `mem_aof_buffer:(.*?)\r\n`)
	mem_allocator := GetStrValue(strs, `mem_allocator:(.*?)\r\n`)
	active_defrag_running := GetIntValue(strs, `active_defrag_running:(.*?)\r\n`)
	lazyfree_pending_objects := GetIntValue(strs, `lazyfree_pending_objects:(.*?)\r\n`)

	data = &models.RedisMemoryDB{
		Hid:                    rcid,
		GetTime:                get_time,
		UsedMemory:             used_memory,
		UsedMemoryHuman:        used_memory_human,
		UsedMemoryRss:          used_memory_rss,
		UsedMemoryRssHuman:     used_memory_rss_human,
		UsedMemoryPeak:         used_memory_peak,
		UsedMemoryPeakHuman:    used_memory_peak_human,
		UsedMemoryPeakPerc:     used_memory_peak_perc,
		UsedMemoryOverhead:     used_memory_overhead,
		UsedMemoryStartup:      used_memory_startup,
		UsedMemoryDataset:      used_memory_dataset,
		UsedMemoryDatasetPerc:  used_memory_dataset_perc,
		AllocatorAllocated:     allocator_allocated,
		AllocatorActive:        allocator_active,
		AllocatorResident:      allocator_resident,
		TotalSystemMemory:      total_system_memory,
		TotalSystemMemoryHuman: total_system_memory_human,
		UsedMemoryLua:          used_memory_lua,
		UsedMemoryLuaHuman:     used_memory_lua_human,
		UsedMemoryScripts:      used_memory_scripts,
		UsedMemoryScriptsHuman: used_memory_scripts_human,
		NumberOfCachedScripts:  number_of_cached_scripts,
		Maxmemory:              maxmemory,
		MaxmemoryHuman:         maxmemory_human,
		MaxmemoryPolicy:        maxmemory_policy,
		AllocatorFragRatio:     allocator_frag_ratio,
		AllocatorFragBytes:     allocator_frag_bytes,
		AllocatorRssRatio:      allocator_rss_ratio,
		AllocatorRssBytes:      allocator_rss_bytes,
		RssOverheadRatio:       rss_overhead_ratio,
		RssOverheadBytes:       rss_overhead_bytes,
		MemFragmentationRatio:  mem_fragmentation_ratio,
		MemFragmentationBytes:  mem_fragmentation_bytes,
		MemNotCountedForEvict:  mem_not_counted_for_evict,
		MemReplicationBacklog:  mem_replication_backlog,
		MemClientsSlaves:       mem_clients_slaves,
		MemClientsNormal:       mem_clients_normal,
		MemAofBuffer:           mem_aof_buffer,
		MemAllocator:           mem_allocator,
		ActiveDefragRunning:    active_defrag_running,
		LazyfreePendingObjects: lazyfree_pending_objects,
	}

	manlog.Debug(data)
	manlog.Debug(*data)
	//err = nil
	return
}

func GetRedisPersistence(strs, rcid string, get_time int64) (data *models.RedisPersistenceDB) {
	loading := GetIntValue(strs, `loading:(.*?)\r\n`)
	rdb_changes_since_last_save := GetIntValue(strs, `rdb_changes_since_last_save:(.*?)\r\n`)
	rdb_bgsave_in_progress := GetIntValue(strs, `rdb_bgsave_in_progress:(.*?)\r\n`)
	rdb_last_save_time := GetIntValue(strs, `rdb_last_save_time:(.*?)\r\n`)
	rdb_last_bgsave_status := GetStrValue(strs, `rdb_last_bgsave_status:(.*?)\r\n`)
	rdb_last_bgsave_time_sec := GetIntValue(strs, `rdb_last_bgsave_time_sec:(.*?)\r\n`)
	rdb_current_bgsave_time_sec := GetIntValue(strs, `rdb_current_bgsave_time_sec:(.*?)\r\n`)
	rdb_last_cow_size := GetIntValue(strs, `rdb_last_cow_size:(.*?)\r\n`)
	aof_enabled := GetIntValue(strs, `aof_enabled:(.*?)\r\n`)
	aof_rewrite_in_progress := GetIntValue(strs, `aof_rewrite_in_progress:(.*?)\r\n`)
	aof_rewrite_scheduled := GetIntValue(strs, `aof_rewrite_scheduled:(.*?)\r\n`)
	aof_last_rewrite_time_sec := GetIntValue(strs, `aof_last_rewrite_time_sec:(.*?)\r\n`)
	aof_current_rewrite_time_sec := GetIntValue(strs, `aof_current_rewrite_time_sec:(.*?)\r\n`)
	aof_last_bgrewrite_status := GetStrValue(strs, `aof_last_bgrewrite_status:(.*?)\r\n`)
	aof_last_write_status := GetStrValue(strs, `aof_last_write_status:(.*?)\r\n`)
	aof_last_cow_size := GetIntValue(strs, `aof_last_cow_size:(.*?)\r\n`)
	aof_current_size := GetIntValue(strs, `aof_current_size:(.*?)\r\n`)
	aof_base_size := GetIntValue(strs, `aof_base_size:(.*?)\r\n`)
	aof_pending_rewrite := GetIntValue(strs, `aof_pending_rewrite:(.*?)\r\n`)
	aof_buffer_length := GetIntValue(strs, `aof_buffer_length:(.*?)\r\n`)
	aof_rewrite_buffer_length := GetIntValue(strs, `aof_rewrite_buffer_length:(.*?)\r\n`)
	aof_pending_bio_fsync := GetIntValue(strs, `aof_pending_bio_fsync:(.*?)\r\n`)
	aof_delayed_fsync := GetIntValue(strs, `aof_delayed_fsync:(.*?)\r\n`)

	data = &models.RedisPersistenceDB{
		Hid:                      rcid,
		GetTime:                  get_time,
		Loading:                  loading,
		RdbChangesSinceLastSave:  rdb_changes_since_last_save,
		RdbBgsaveInProgress:      rdb_bgsave_in_progress,
		RdbLastSaveTime:          rdb_last_save_time,
		RdbLastBgsaveStatus:      rdb_last_bgsave_status,
		RdbLastBgsaveTimeSec:     rdb_last_bgsave_time_sec,
		RdbCurrentBgsaveTimeSec:  rdb_current_bgsave_time_sec,
		RdbLastCowSize:           rdb_last_cow_size,
		AofEnabled:               aof_enabled,
		AofRewriteInProgress:     aof_rewrite_in_progress,
		AofRewriteScheduled:      aof_rewrite_scheduled,
		AofLastRewriteTimeSec:    aof_last_rewrite_time_sec,
		AofCurrentRewriteTimeSec: aof_current_rewrite_time_sec,
		AofLastBgrewriteStatus:   aof_last_bgrewrite_status,
		AofLastWriteStatus:       aof_last_write_status,
		AofLastCowSize:           aof_last_cow_size,
		AofCurrentSize:           aof_current_size,
		AofBaseSize:              aof_base_size,
		AofPendingRewrite:        aof_pending_rewrite,
		AofBufferLength:          aof_buffer_length,
		AofRewriteBufferLength:   aof_rewrite_buffer_length,
		AofPendingBioFsync:       aof_pending_bio_fsync,
		AofDelayedFsync:          aof_delayed_fsync,
	}

	manlog.Debug(data)
	manlog.Debug(*data)
	//err = nil
	return
}

func GetRedisStats(strs, rcid string, get_time int64) (data *models.RedisStatsDB) {
	total_connections_received := GetIntValue(strs, `total_connections_received:(.*?)\r\n`)
	total_commands_processed := GetIntValue(strs, `total_commands_processed:(.*?)\r\n`)
	instantaneous_ops_per_sec := GetIntValue(strs, `instantaneous_ops_per_sec:(.*?)\r\n`)
	total_net_input_bytes := GetIntValue(strs, `total_net_input_bytes:(.*?)\r\n`)
	total_net_output_bytes := GetIntValue(strs, `total_net_output_bytes:(.*?)\r\n`)
	instantaneous_input_kbps := GetFloatValue(strs, `instantaneous_input_kbps:(.*?)\r\n`)
	instantaneous_output_kbps := GetFloatValue(strs, `instantaneous_output_kbps:(.*?)\r\n`)
	rejected_connections := GetIntValue(strs, `rejected_connections:(.*?)\r\n`)
	sync_full := GetIntValue(strs, `sync_full:(.*?)\r\n`)
	sync_partial_ok := GetIntValue(strs, `sync_partial_ok:(.*?)\r\n`)
	sync_partial_err := GetIntValue(strs, `sync_partial_err:(.*?)\r\n`)
	expired_keys := GetIntValue(strs, `expired_keys:(.*?)\r\n`)
	expired_stale_perc := GetStrValue(strs, `expired_stale_perc:(.*?)\r\n`)
	expired_time_cap_reached_count := GetIntValue(strs, `expired_time_cap_reached_count:(.*?)\r\n`)
	evicted_keys := GetIntValue(strs, `evicted_keys:(.*?)\r\n`)
	keyspace_hits := GetIntValue(strs, `keyspace_hits:(.*?)\r\n`)
	keyspace_misses := GetIntValue(strs, `keyspace_misses:(.*?)\r\n`)
	pubsub_channels := GetIntValue(strs, `pubsub_channels:(.*?)\r\n`)
	pubsub_patterns := GetIntValue(strs, `pubsub_patterns:(.*?)\r\n`)
	latest_fork_usec := GetIntValue(strs, `latest_fork_usec:(.*?)\r\n`)
	migrate_cached_sockets := GetIntValue(strs, `migrate_cached_sockets:(.*?)\r\n`)
	slave_expires_tracked_keys := GetIntValue(strs, `slave_expires_tracked_keys:(.*?)\r\n`)
	active_defrag_hits := GetIntValue(strs, `active_defrag_hits:(.*?)\r\n`)
	active_defrag_misses := GetIntValue(strs, `active_defrag_misses:(.*?)\r\n`)
	active_defrag_key_hits := GetIntValue(strs, `active_defrag_key_hits:(.*?)\r\n`)
	active_defrag_key_misses := GetIntValue(strs, `active_defrag_key_misses:(.*?)\r\n`)

	data = &models.RedisStatsDB{
		Hid:                        rcid,
		GetTime:                    get_time,
		TotalConnectionsReceived:   total_connections_received,
		TotalCommandsProcessed:     total_commands_processed,
		InstantaneousOpsPerSec:     instantaneous_ops_per_sec,
		TotalNetInputBytes:         total_net_input_bytes,
		TotalNetOutputBytes:        total_net_output_bytes,
		InstantaneousInputKbps:     instantaneous_input_kbps,
		InstantaneousOutputKbps:    instantaneous_output_kbps,
		RejectedConnections:        rejected_connections,
		SyncFull:                   sync_full,
		SyncPartialOk:              sync_partial_ok,
		SyncPartialErr:             sync_partial_err,
		ExpiredKeys:                expired_keys,
		ExpiredStalePerc:           expired_stale_perc,
		ExpiredTimeCapReachedCount: expired_time_cap_reached_count,
		EvictedKeys:                evicted_keys,
		KeyspaceHits:               keyspace_hits,
		KeyspaceMisses:             keyspace_misses,
		PubsubChannels:             pubsub_channels,
		PubsubPatterns:             pubsub_patterns,
		LatestForkUsec:             latest_fork_usec,
		MigrateCachedSockets:       migrate_cached_sockets,
		SlaveExpiresTrackedKeys:    slave_expires_tracked_keys,
		ActiveDefragHits:           active_defrag_hits,
		ActiveDefragMisses:         active_defrag_misses,
		ActiveDefragKeyHits:        active_defrag_key_hits,
		ActiveDefragKeyMisses:      active_defrag_key_misses,
	}

	manlog.Debug(data)
	manlog.Debug(*data)
	//err = nil
	return
}

func GetRedisReplication(strs, rcid string, get_time int64) (data *models.RedisReplicationDB) {

	role := GetStrValue(strs, `role:(.*?)\r\n`)
	connected_slaves := GetStrValue(strs, `connected_slaves:(.*?)\r\n`)
	master_replid := GetStrValue(strs, `master_replid:(.*?)\r\n`)
	master_replid2 := GetStrValue(strs, `master_replid2:(.*?)\r\n`)
	master_repl_offset := GetStrValue(strs, `master_repl_offset:(.*?)\r\n`)
	second_repl_offset := GetStrValue(strs, `second_repl_offset:(.*?)\r\n`)
	repl_backlog_active := GetStrValue(strs, `repl_backlog_active:(.*?)\r\n`)
	repl_backlog_size := GetIntValue(strs, `repl_backlog_size:(.*?)\r\n`)
	repl_backlog_first_byte_offset := GetIntValue(strs, `repl_backlog_first_byte_offset:(.*?)\r\n`)
	repl_backlog_histlen := GetIntValue(strs, `repl_backlog_histlen:(.*?)\r\n`)

	data = &models.RedisReplicationDB{
		Hid:                        rcid,
		GetTime:                    get_time,
		Role:                       role,
		ConnectedSlaves:            connected_slaves,
		MasterReplid:               master_replid,
		MasterReplid2:              master_replid2,
		MasterReplOffset:           master_repl_offset,
		SecondReplOffset:           second_repl_offset,
		ReplBacklogActive:          repl_backlog_active,
		ReplBacklogSize:            repl_backlog_size,
		ReplBacklogFirstByteOffset: repl_backlog_first_byte_offset,
		ReplBacklogHistlen:         repl_backlog_histlen,
	}

	manlog.Debug(data)
	manlog.Debug(*data)
	//err = nil
	return
}

func GetRedisCPU(strs, rcid string, get_time int64) (data *models.RedisCPUDB) {

	used_cpu_sys := GetFloatValue(strs, `used_cpu_sys:(.*?)\r\n`)
	used_cpu_user := GetFloatValue(strs, `used_cpu_user:(.*?)\r\n`)
	used_cpu_sys_children := GetFloatValue(strs, `used_cpu_sys_children:(.*?)\r\n`)
	used_cpu_user_children := GetFloatValue(strs, `used_cpu_user_children:(.*?)\r\n`)

	data = &models.RedisCPUDB{
		Hid:                 rcid,
		GetTime:             get_time,
		UsedCpuSys:          used_cpu_sys,
		UsedCpuUser:         used_cpu_user,
		UsedCpuSysChildren:  used_cpu_sys_children,
		UsedCpuUserChildren: used_cpu_user_children,
	}

	manlog.Debug(data)
	manlog.Debug(*data)
	//err = nil
	return
}

func GetRedisCluster(strs, rcid string, get_time int64) (data *models.RedisClusterDB) {

	cluster_enabled := GetStrValue(strs, `cluster_enabled:(.*?)\r\n`)

	data = &models.RedisClusterDB{
		Hid:            rcid,
		GetTime:        get_time,
		ClusterEnabled: cluster_enabled,
	}

	manlog.Debug(data)
	manlog.Debug(*data)
	//err = nil
	return
}

func GetRedisKeyspace(strs, rcid string, get_time int64) (datas []*models.RedisKeyspaceDB) {

	spit := strings.Index(strs, "# Keyspace")
	value := string([]rune(strs)[spit:])
	manlog.Debug(value)
	value_list := strings.Split(value, "\r\n")
	manlog.Debug(value_list)
	for _, v := range value_list {
		if string([]rune(v)[0:2]) == "db" {

			manlog.Debug(v)
			id := GetIntValue(v, `db(.*?):`)
			keys := GetIntValue(v, `keys=(.*?),`)
			expires := GetIntValue(v, `expires=(.*?),`)
			avg_ttl := GetIntValue(v, `avg_ttl=(.*?)$`)
			manlog.Debug(id, keys, expires, avg_ttl)

			datas = append(datas, &models.RedisKeyspaceDB{
				Hid:     rcid,
				GetTime: get_time,
				DBID:    id,
				Keys:    keys,
				Expires: expires,
				AvgTTL:  avg_ttl,
			})

		}
	}

	manlog.Debug(datas)
	for _, v := range datas {
		manlog.Debug(*v)
	}
	//err = nil
	return
}

func GetStrValue(strs string, reg string) string {
	valuelist := common.FindAllstrlist(reg, strs)
	if len(valuelist) > 0 && len(valuelist[0]) > 1 {
		return valuelist[0][1]
	}
	return ""
}

func GetIntValue(strs string, reg string) int64 {
	valuelist := common.FindAllstrlist(reg, strs)

	// for _, v := range valuelist {
	// 	manlog.Debug(v)
	// }

	if len(valuelist) > 0 && len(valuelist[0]) > 1 {
		val_str := valuelist[0][1]
		return common.Str2Int64(val_str)
	}
	//manlog.Error("获取失败, str = ", strs, "; reg = ", reg)
	return 0
}

func GetFloatValue(strs string, reg string) float64 {
	valuelist := common.FindAllstrlist(reg, strs)
	if len(valuelist) > 0 && len(valuelist[0]) > 1 {
		val_str := valuelist[0][1]
		return common.Str2Float64(val_str)
	}
	//manlog.Error("获取失败")
	return 0
}
