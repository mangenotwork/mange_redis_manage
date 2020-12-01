//
//	redis 信息服务
//
package service

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"
	_ "unicode/utf8"
	_ "unsafe"

	"github.com/garyburd/redigo/redis"
	_ "github.com/mangenotwork/mange_redis_manage/common"
	"github.com/mangenotwork/mange_redis_manage/common/cache"
	"github.com/mangenotwork/mange_redis_manage/common/id"
	"github.com/mangenotwork/mange_redis_manage/common/manlog"
	manredis "github.com/mangenotwork/mange_redis_manage/common/redis"
	"github.com/mangenotwork/mange_redis_manage/dao"
	"github.com/mangenotwork/mange_redis_manage/repository"
	"github.com/mangenotwork/mange_redis_manage/structs"
)

type RedisInfoService interface {
	GetInfo(user *structs.UserParameter, redisId int64) (datas string, err error)                              //redis 服务信息
	RedisRealTimeInit(user *structs.UserParameter, redisId int64) (data *structs.RealTime, err error)          //第一次请求监控获取监听id与历史数据
	RedisRealTime(user *structs.UserParameter, redisId int64, rtid string) (data *structs.RealTime, err error) //实时获取当前时间的监听数据
	GetKeyInfo(user *structs.UserParameter, rid, dbid int64, key string) (data *structs.KeyInfo, err error)    //获取key信息
	GetConfig(user *structs.UserParameter, redisId int64) (datas []*structs.RedisConfigData, err error)        //获取服务器配置信息
}

type RedisInfo struct {
}

func (this *RedisInfo) GetInfo(user *structs.UserParameter, redisId int64) (datas string, err error) {
	//获取连接
	rc, rcid, err := new(RedisConn).DBconn(user, redisId, 0)
	if err != nil {
		manlog.Error(err)
	}

	//获取服务信息
	data := manredis.GetRedisServersInfos(rc, rcid)
	manlog.Debug(data)
	b, err := json.MarshalIndent(data, "", "\t")
	if err != nil {

		fmt.Println("Umarshal failed:", err)
		return
	}
	datas = string(b)
	return
}

func (this *RedisInfo) GetConfig(user *structs.UserParameter, redisId int64) (datas []*structs.RedisConfigData, err error) {
	datas = make([]*structs.RedisConfigData, 0)

	//获取连接
	rc, _, err := new(RedisConn).DBconn(user, redisId, 0)
	if err != nil {
		manlog.Error(err)
		return
	}
	//获取服务信息
	data, err := manredis.GetRedisServersAllConfig(rc)
	if err != nil {
		manlog.Error(err)
		return
	}

	for i := 0; i < len(data); i++ {
		if i%2 != 0 {
			datas = append(datas, &structs.RedisConfigData{
				ConfigName:  data[i-1],
				ConfigValue: data[i],
				ConfigDoc:   RedisConfigDoc[data[i-1]],
			})
		}
	}

	return
}

//第一次请求
func (this *RedisInfo) RedisRealTimeInit(user *structs.UserParameter, redisId int64) (data *structs.RealTime, err error) {
	data = &structs.RealTime{}

	//第一次请求时分发一个id
	iddata, err := id.IdInt64()
	if err != nil {
		manlog.Error(err)
	}
	data.RealTimeId = fmt.Sprintf("%d", iddata)

	//获取当前时间
	nowtime := time.Now().Unix()
	//获取连接
	_, rcid, err := new(RedisConn).DBconn(user, redisId, 0)
	if err != nil {
		manlog.Error(err)
	}

	time_list := []string{}

	//获取cpu
	used_cpu_sys_list := []float64{}
	used_cpu_user_list := []float64{}
	//查询当前时间之前的11条数据
	cpu_datas, err := new(dao.DaoRedisCPU).GetLastTimeData(rcid, nowtime, 21)
	if err != nil {
		manlog.Error(err)
	}
	sort.Slice(cpu_datas, func(i, j int) bool {
		return cpu_datas[i].ID < cpu_datas[j].ID
	})
	for i := 0; i < len(cpu_datas)-1; i++ {
		manlog.Debug("查询当前时间之前的10条数据 1= ", *cpu_datas[i])
		manlog.Debug("查询当前时间之前的10条数据 2= ", *cpu_datas[i+1])
		now_data := cpu_datas[i+1]
		before_data := cpu_datas[i]

		used_cpu_sys_value, used_cpu_user_value := this.RealTimeCPU(now_data, before_data)
		used_cpu_sys_list = append(used_cpu_sys_list, used_cpu_sys_value)
		used_cpu_user_list = append(used_cpu_user_list, used_cpu_user_value)

		time_list = append(time_list, time.Unix(cpu_datas[i+1].GetTime, 0).Format("01-02 15:04:05"))
	}
	manlog.Debug(used_cpu_sys_list, used_cpu_user_list)
	data.Xdata = time_list
	data.CPUData = used_cpu_sys_list
	//存入最新时间的DaoRedisCPU到缓存
	cpu_key := fmt.Sprintf("rt_cpu_%s", data.RealTimeId)
	cache.Set(cpu_key, cpu_datas[len(cpu_datas)-1])

	//获取内存
	memory_list := []float64{}
	memory_datas, err := new(dao.DaoRedisMemory).GetLastTimeData(rcid, nowtime, 20)
	if err != nil {
		manlog.Error(err)
	}
	sort.Slice(memory_datas, func(i, j int) bool {
		return memory_datas[i].ID < memory_datas[j].ID
	})
	memory_dw := "B"
	if len(memory_datas) > 0 {
		if memory_datas[0].UsedMemory > 1*1204 {
			memory_dw = "KB"
		}
		if memory_datas[0].UsedMemory > 1*1204*1024 {
			memory_dw = "MB"
		}
		if memory_datas[0].UsedMemory > 1*1204*1024*1204 {
			memory_dw = "GB"
		}
	}
	for _, v := range memory_datas {
		var used_memory float64
		switch memory_dw {
		case "B":
			used_memory = float64(v.UsedMemory)
		case "KB":
			used_memory = float64(v.UsedMemory) / (1 * 1024)
		case "MB":
			used_memory = float64(v.UsedMemory) / (1 * 1024 * 1024)
		case "GB":
			used_memory = float64(v.UsedMemory) / (1 * 1024 * 1024 * 1024)
		}
		used_memory, _ = this.threefloat64(used_memory)
		memory_list = append(memory_list, used_memory)
	}
	manlog.Debug(memory_list)
	data.MemoryData = memory_list
	data.MemoryDW = memory_dw
	//存入最新时间的DaoRedisMemory到缓存
	memory_key := fmt.Sprintf("rt_memory_%s", data.RealTimeId)
	cache.Set(memory_key, memory_datas[len(memory_datas)-1])

	//获取Stats   QPS,流量
	qps_list := []int64{}
	keys_list := []int64{}
	input_kbps_list := []float64{}
	output_kbps_list := []float64{}
	stats_datas, err := new(dao.DaoRedisStats).GetLastTimeData(rcid, nowtime, 21)
	if err != nil {
		manlog.Error(err)
	}
	sort.Slice(stats_datas, func(i, j int) bool {
		return stats_datas[i].ID < stats_datas[j].ID
	})
	for i := 0; i < len(stats_datas)-1; i++ {
		now_data := stats_datas[i+1].TotalCommandsProcessed
		before_data := stats_datas[i].TotalCommandsProcessed
		now_time := stats_datas[i+1].GetTime
		before_time := stats_datas[i].GetTime
		qps := (now_data - before_data) / (now_time - before_time)
		qps_list = append(qps_list, int64(qps))
		keys_list = append(keys_list, stats_datas[i+1].ExpiredKeys)
		input_kbps_list = append(input_kbps_list, stats_datas[i+1].InstantaneousInputKbps)
		output_kbps_list = append(output_kbps_list, stats_datas[i+1].InstantaneousOutputKbps)
	}
	data.QpsData = qps_list
	data.KeysData = keys_list
	data.InputKbpsData = input_kbps_list
	data.OutputKbpsData = output_kbps_list
	//存入最新时间的DaoRedisStats到缓存
	stats_key := fmt.Sprintf("rt_stats_%s", data.RealTimeId)
	cache.Set(stats_key, stats_datas[len(stats_datas)-1])

	//获取Clients
	clients_list := []int64{}
	clients_datas, err := new(dao.DaoRedisClients).GetLastTimeData(rcid, nowtime, 20)
	if err != nil {
		manlog.Error(err)
	}
	sort.Slice(clients_datas, func(i, j int) bool {
		return clients_datas[i].ID < clients_datas[j].ID
	})
	for _, v := range clients_datas {
		clients_list = append(clients_list, v.ConnectedClients)
	}
	data.ConnData = clients_list
	//存入最新时间的DaoRedisClients到缓存
	clients_key := fmt.Sprintf("rt_clients_%s", data.RealTimeId)
	cache.Set(clients_key, clients_datas[len(clients_datas)-1])

	return
}

//实时获取
//由于需要上下秒计算，需要缓存一次取的数据
//缓存的key,进入监控页分发一个随机id,
//缓存的值，计算的原始值
func (this *RedisInfo) RedisRealTime(user *structs.UserParameter, redisId int64, rtid string) (data *structs.RealTime, err error) {

	//1.获取数据
	data = &structs.RealTime{}
	data.RealTimeId = rtid

	rc, rcid, err := new(RedisConn).DBconn(user, redisId, 0)
	if err != nil {
		manlog.Error(err)
	}

	cpu_last_data := &models.RedisCPUDB{}
	memory_last_data := &models.RedisMemoryDB{}
	stats_last_data := &models.RedisStatsDB{}
	clients_last_data := &models.RedisClientsDB{}

	//2.获取缓存，如果没有缓存id分发缓存id并存入缓存,取数据里最新一条数据
	cpu_key := fmt.Sprintf("rt_cpu_%s", rtid)
	memory_key := fmt.Sprintf("rt_memory_%s", rtid)
	stats_key := fmt.Sprintf("rt_stats_%s", rtid)
	clients_key := fmt.Sprintf("rt_clients_%s", rtid)
	manlog.Error("cpu_key = ", cpu_key)

	//cpu
	cpu_cache, cpu_is := cache.Get(cpu_key)
	manlog.Error(cpu_cache, cpu_is)
	if cpu_is {
		cpu_last_data = cpu_cache.(*models.RedisCPUDB)
	} else {
		//如果没有缓存，取数据里最新一条数据
		cpu_last_data, err = new(dao.DaoRedisCPU).GetNewData(rcid)
		if err != nil {
			manlog.Error(err)
		}
	}
	manlog.Debug(*cpu_last_data)

	//memory
	memory_cache, memory_is := cache.Get(memory_key)
	manlog.Error(memory_cache, memory_is)
	if memory_is {
		memory_last_data = memory_cache.(*models.RedisMemoryDB)
	} else {
		memory_last_data, err = new(dao.DaoRedisMemory).GetNewData(rcid)
		if err != nil {
			manlog.Error(err)
		}
	}
	manlog.Debug(*memory_last_data)

	//stats
	stats_cache, stats_is := cache.Get(stats_key)
	manlog.Error(stats_cache, stats_is)
	if stats_is {
		stats_last_data = stats_cache.(*models.RedisStatsDB)
	} else {
		stats_last_data, err = new(dao.DaoRedisStats).GetNewData(rcid)
		if err != nil {
			manlog.Error(err)
		}
	}
	manlog.Debug(*stats_last_data)

	//clients
	clients_cache, clients_is := cache.Get(clients_key)
	manlog.Error(clients_cache, clients_is)
	if clients_is {
		clients_last_data = clients_cache.(*models.RedisClientsDB)
	} else {
		clients_last_data, err = new(dao.DaoRedisClients).GetNewData(rcid)
		if err != nil {
			manlog.Error(err)
		}
	}
	manlog.Debug(*clients_last_data)
	//manlog.Panic(0)
	//3.获取当前值
	now_data := manredis.GetRedisServersInfos(rc, rcid)

	//4.计算值
	time_list := []string{}

	//cpu
	used_cpu_sys_list := []float64{}
	used_cpu_user_list := []float64{}
	used_cpu_sys_value, used_cpu_user_value := this.RealTimeCPU(now_data.CPU, cpu_last_data)
	used_cpu_sys_list = append(used_cpu_sys_list, used_cpu_sys_value)
	used_cpu_user_list = append(used_cpu_user_list, used_cpu_user_value)
	time_list = append(time_list, time.Unix(now_data.CPU.GetTime, 0).Format("01-02 15:04:05"))
	data.Xdata = time_list
	data.CPUData = used_cpu_sys_list
	data.CPUTip = "cpu"
	//存入最新时间的RedisCPUDB到缓存
	cache.Set(cpu_key, now_data.CPU)

	//获取内存
	memory_list := []float64{}
	memory_dw := "B"
	if now_data.Memory.UsedMemory > 1*1204 {
		memory_dw = "KB"
	}
	if now_data.Memory.UsedMemory > 1*1204*1024 {
		memory_dw = "MB"
	}
	if now_data.Memory.UsedMemory > 1*1204*1024*1204 {
		memory_dw = "GB"
	}
	var used_memory float64
	switch memory_dw {
	case "B":
		used_memory = float64(now_data.Memory.UsedMemory)
	case "KB":
		used_memory = float64(now_data.Memory.UsedMemory) / (1 * 1024)
	case "MB":
		used_memory = float64(now_data.Memory.UsedMemory) / (1 * 1024 * 1024)
	case "GB":
		used_memory = float64(now_data.Memory.UsedMemory) / (1 * 1024 * 1024 * 1024)
	}
	used_memory, _ = this.threefloat64(used_memory)
	memory_list = append(memory_list, used_memory)
	data.MemoryData = memory_list
	data.MemoryDW = memory_dw
	data.MemoryTip = "memory"
	//存入最新时间的DaoRedisMemory到缓存
	cache.Set(memory_key, now_data.Memory)

	//获取Stats   QPS,流量
	qps_list := []int64{}
	keys_list := []int64{}
	input_kbps_list := []float64{}
	output_kbps_list := []float64{}
	stats_now_data := now_data.Stats.TotalCommandsProcessed
	stats_before_data := stats_last_data.TotalCommandsProcessed
	stats_now_time := now_data.Stats.GetTime
	stats_before_time := stats_last_data.GetTime
	qps := (stats_now_data - stats_before_data) / (stats_now_time - stats_before_time)
	qps_list = append(qps_list, int64(qps))
	keys_list = append(keys_list, now_data.Stats.ExpiredKeys)
	input_kbps_list = append(input_kbps_list, now_data.Stats.InstantaneousInputKbps)
	output_kbps_list = append(output_kbps_list, now_data.Stats.InstantaneousOutputKbps)
	data.QpsData = qps_list
	data.QpsTip = "qps"
	data.KeysData = keys_list
	data.KeysTip = "keys"
	data.InputKbpsData = input_kbps_list
	data.InputKbpsTip = "input"
	data.OutputKbpsData = output_kbps_list
	data.OutputKbpsTip = "output"
	//存入最新时间的DaoRedisStats到缓存
	cache.Set(stats_key, now_data.Stats)

	//获取Clients
	clients_list := []int64{}
	clients_list = append(clients_list, now_data.Clients.ConnectedClients)
	data.ConnData = clients_list
	data.ConnTip = "conn"
	//存入最新时间的DaoRedisClients到缓存
	cache.Set(clients_key, now_data.Clients)

	return
}

//实时cpu 监控
// redis进程单cpu的消耗率可以通过如下公式计算:
// ((used_cpu_sys_now-used_cpu_sys_before)/(now-before))*100
// 其中
// used_cpu_sys_now:now时间点的used_cpu_sys值
// used_cpu_sys_before:before时间点的used_cpu_sys值
func (this *RedisInfo) RealTimeCPU(now_data, before_data *models.RedisCPUDB) (used_cpu_sys_value, used_cpu_user_value float64) {
	used_cpu_sys_value = ((now_data.UsedCpuSys - before_data.UsedCpuSys) / float64(now_data.GetTime-before_data.GetTime)) * 100
	used_cpu_user_value = ((now_data.UsedCpuUser - before_data.UsedCpuUser) / float64(now_data.GetTime-before_data.GetTime)) * 100
	used_cpu_sys_value, _ = this.threefloat64(used_cpu_sys_value)
	used_cpu_user_value, _ = this.threefloat64(used_cpu_user_value)
	return
}

//float64 只保留3位小数
func (this *RedisInfo) threefloat64(v float64) (float64, error) {
	return strconv.ParseFloat(fmt.Sprintf("%.3f", v), 3)
}

//实时内存
func (this *RedisInfo) RealTimeMemory() {}

//实时QPS
func (this *RedisInfo) RealTimeQPS() {
	// 使用redis-cli中info统计信息计算差值；
	// redis-cli的info命令中有一项total_commands_processed表示：从启动到现在处理的所有命令总数，可以通过统计两次info指令间的差值来计算QPS：
}

//实时连接数
func (this *RedisInfo) RealTimeConnCount() {}

//实时key数量
func (this *RedisInfo) RealTimeKeyCount() {}

//实时流量
func (this *RedisInfo) RealTimeNetKbps() {}

//实时命中
func (this *RedisInfo) RealTimeHitRate() {
	//https://www.jianshu.com/p/6f3ee2ca599f
	//keyspace_hit / ( keyspace_hit + keyspace_misses ) = hit_rate

	//请求键的命中率 (keyspace_hit_ratio):使用keyspace_hits/(keyspace_hits+keyspace_misses)计算所得，命中率低于50%告警

}

//key 信息
func (this *RedisInfo) GetKeyInfo(user *structs.UserParameter, redisId, dbid int64, key string) (data *structs.KeyInfo, err error) {
	data = new(structs.KeyInfo)
	data.KeyName = key

	//获取连接
	rc, _, err := new(RedisConn).DBconn(user, redisId, dbid)
	if err != nil {
		manlog.Error(err)
		return
	}

	//获取值和key类型与值的大小
	value, keytype, size, _ := this.GetKeyValue(rc, key)
	data.Value = value
	data.KeyType = keytype
	data.KeySize = size
	data.DBID = dbid

	//获取key的ttl
	ttl := manredis.GetKeyTTL(rc, key)
	manlog.Debug(ttl)
	data.TTL = ttl

	return

}

//内部调用方法，获取key的value
func (this *RedisInfo) GetKeyValue(rc redis.Conn, key string) (value interface{}, keytype, size string, err error) {
	key_type := manredis.GetKeyType(rc, key)
	manlog.Debug(key_type)

	var factory = new(RedisTypeOperateFactory)
	keyObj := factory.Operate(key_type)
	value, keytype, size, err = keyObj.Get(rc, key)

	// switch key_type {
	// case "string":
	// 	keytype = "string"
	// 	string_value := manredis.StringGet(rc, key)
	// 	value = string_value
	// 	//计算字节大小,单位b,  1kb=1024b
	// 	string_value_size := utf8.RuneCount([]byte(string_value))
	// 	size = fmt.Sprintf("%dByte", string_value_size)

	// case "hash":
	// 	manlog.Debug("hash")
	// 	keytype = "hash"
	// 	hash_value := manredis.HashHGETALL(rc, key)
	// 	hash_value_size := unsafe.Sizeof(hash_value)
	// 	hash_value_result, _ := json.Marshal(hash_value)
	// 	hash_value_str := string(hash_value_result)
	// 	value = hash_value_str
	// 	size = fmt.Sprintf("%dByte", hash_value_size)

	// case "list":
	// 	manlog.Debug("list")
	// 	keytype = "list"
	// 	list_value := manredis.ListLRANGE(rc, key)
	// 	list_value_map := make(map[int]string, 0)
	// 	list_size := 0
	// 	for v, k := range list_value {
	// 		kStr := common.Uint82Str(k.([]uint8))
	// 		list_value_map[v] = kStr
	// 		list_size = list_size + utf8.RuneCount([]byte(kStr))
	// 	}
	// 	value = list_value_map
	// 	size = fmt.Sprintf("%dByte", list_size)

	// case "set":
	// 	manlog.Debug("set")
	// 	keytype = "set"
	// 	set_value := manredis.SetSMEMBERS(rc, key)
	// 	set_value_map := make(map[int]string, 0)
	// 	set_size := 0
	// 	for v, k := range set_value {
	// 		kStr := common.Uint82Str(k.([]uint8))
	// 		set_value_map[v] = kStr
	// 		set_size = set_size + utf8.RuneCount([]byte(kStr))
	// 	}
	// 	value = set_value_map
	// 	size = fmt.Sprintf("%dByte", set_size)

	// case "zset":
	// 	manlog.Debug("zset")
	// 	keytype = "zset"
	// 	zset_value := manredis.ZSetZRANGE(rc, key)
	// 	zset_value_map := make(map[int]string, 0)
	// 	zset_size := 0
	// 	for v, k := range zset_value {
	// 		kStr := common.Uint82Str(k.([]uint8))
	// 		zset_value_map[v] = kStr
	// 		zset_size = zset_size + utf8.RuneCount([]byte(kStr))
	// 	}
	// 	value = zset_value_map
	// 	size = fmt.Sprintf("%dByte", zset_size)
	// }
	err = nil
	return
}

//redis 监控报警， 存活
func (this *RedisInfo) AlarmPing() {
	//使用ping 命令
}

//redis 监控报警，cpu
func (this *RedisInfo) AlarmCPU() {
	//同上CPU
}

//redis 监控报警，Memory
func (this *RedisInfo) AlarmMemory() {
	//同上Memory
}

var RedisConfigDoc = map[string]string{
	"daemonize":      "Redis 默认不是以守护进程的方式运行，可以通过该配置项修改，使用 yes 启用守护进程（Windows 不支持守护线程的配置为 no ）",
	"pidfile":        "当 Redis 以守护进程方式运行时，Redis 默认会把 pid 写入 /var/run/redis.pid 文件，可以通过 pidfile 指定",
	"port":           "指定 Redis 监听端口，默认端口为 6379，作者在自己的一篇博文中解释了为什么选用 6379 作为默认端口，因为 6379 在手机按键上 MERZ 对应的号码，而 MERZ 取自意大利歌女 Alessia Merz 的名字",
	"bind":           "绑定的主机地址",
	"timeout":        "当客户端闲置多长秒后关闭连接，如果指定为 0 ，表示关闭该功能",
	"loglevel":       "指定日志记录级别，Redis 总共支持四个级别：debug、verbose、notice、warning，默认为 notice",
	"logfile":        "日志记录方式，默认为标准输出，如果配置 Redis 为守护进程方式运行，而这里又配置为日志记录方式为标准输出，则日志将会发送给 /dev/null",
	"databases":      "设置数据库的数量，默认数据库为0，可以使用SELECT 命令在连接上指定数据库id",
	"save":           "指定在多长时间内，有多少次更新操作，就将数据同步到数据文件，可以多个条件配合",
	"rdbcompression": "指定存储至本地数据库时是否压缩数据，默认为 yes，Redis 采用 LZF 压缩，如果为了节省 CPU 时间，可以关闭该选项，但会导致数据库文件变的巨大",
	"dbfilename":     "指定本地数据库文件名，默认值为 dump.rdb",
	"dir":            "指定本地数据库存放目录",
	"slaveof":        "设置当本机为 slave 服务时，设置 master 服务的 IP 地址及端口，在 Redis 启动时，它会自动从 master 进行数据同步",
	"masterauth":     "当 master 服务设置了密码保护时，slav 服务连接 master 的密码",
	"requirepass":    "设置 Redis 连接密码，如果配置了连接密码，客户端在连接 Redis 时需要通过 AUTH <password> 命令提供密码，默认关闭",
	"maxclients": "设置同一时间最大客户端连接数，默认无限制，Redis 可以同时打开的客户端连接数为 Redis 进程可以打开的最大文件描述符数，如果设置 maxclients 0，" +
		"表示不作限制。当客户端连接数到达限制时，Redis 会关闭新的连接并向客户端返回 max number of clients reached 错误信息",
	"maxmemory": "指定 Redis 最大内存限制，Redis 在启动时会把数据加载到内存中，达到最大内存后，Redis 会先尝试清除已到期或即将到期的 Key，当此方法处理 后，" +
		"仍然到达最大内存设置，将无法再进行写入操作，但仍然可以进行读取操作。Redis 新的 vm 机制，会把 Key 存放内存，Value 会存放在 swap 区",
	"appendonly": "指定是否在每次更新操作后进行日志记录，Redis 在默认情况下是异步的把数据写入磁盘，如果不开启，可能会在断电时导致一段时间内的数据丢失。" +
		"因为 redis 本身同步数据文件是按上面 save 条件来同步的，所以有的数据会在一段时间内只存在于内存中。默认为 no",
	"appendfilename": "指定更新日志文件名，默认为 appendonly.aof",
	"appendfsync": "指定更新日志条件，共有 3 个可选值：no：表示等操作系统进行数据缓存同步到磁盘（快）;always：表示每次更新操作后手动调用 fsync() 将数据写到磁盘（慢，安全）;" +
		"everysec：表示每秒同步一次（折中，默认值）",
	"vm-enabled": "指定是否启用虚拟内存机制，默认值为 no，简单的介绍一下，VM 机制将数据分页存放，由 Redis 将访问量较少的页即冷数据 swap 到磁盘上，" +
		"访问多的页面由磁盘自动换出到内存中（在后面的文章我会仔细分析 Redis 的 VM 机制）",
	"vm-swap-file": "虚拟内存文件路径，默认值为 /tmp/redis.swap，不可多个 Redis 实例共享",
	"vm-max-memory": "将所有大于 vm-max-memory 的数据存入虚拟内存，无论 vm-max-memory 设置多小，所有索引数据都是内存存储的(Redis 的索引数据 就是 keys)，" +
		"也就是说，当 vm-max-memory 设置为 0 的时候，其实是所有 value 都存在于磁盘。默认值为 0",
	"vm-page-size": "Redis swap 文件分成了很多的 page，一个对象可以保存在多个 page 上面，但一个 page 上不能被多个对象共享，vm-page-size 是要根据存储的 数据大小来设定的，" +
		"作者建议如果存储很多小对象，page 大小最好设置为 32 或者 64bytes；如果存储很大大对象，则可以使用更大的 page，如果不确定，就使用默认值",
	"vm-pages":                "设置 swap 文件中的 page 数量，由于页表（一种表示页面空闲或使用的 bitmap）是在放在内存中的，，在磁盘上每 8 个 pages 将消耗 1byte 的内存。",
	"vm-max-threads":          "设置访问swap文件的线程数,最好不要超过机器的核数,如果设置为0,那么所有对swap文件的操作都是串行的，可能会造成比较长时间的延迟。默认值为4",
	"glueoutputbuf":           "设置在向客户端应答时，是否把较小的包合并为一个包发送，默认为开启",
	"hash-max-zipmap-entries": "指定在超过一定的数量或者最大的元素超过某一临界值时，采用一种特殊的哈希算法",
	"hash-max-zipmap-value":   "指定在超过一定的数量或者最大的元素超过某一临界值时，采用一种特殊的哈希算法",
	"activerehashing":         "指定是否激活重置哈希，默认为开启（后面在介绍 Redis 的哈希算法时具体介绍）",
	"include":                 "指定包含其它的配置文件，可以在同一主机上多个Redis实例之间使用同一份配置文件，而同时各个实例又拥有自己的特定配置文件",
	"tcp-keepalive":           "指定TCP连接是否为长连接,'侦探'信号有server端维护。默认为0表示禁用",
	"stop-writes-on-bgsave-error": "当持久化出现错误时，是否依然继续进行工作，是否终止所有的客户端write请求。默认设置'yes'表示终止，一旦snapshot数据保存故障，" +
		"那么此server为只读服务。如果为'no'，那么此次snapshot将失败，但下一次snapshot不会受到影响，不过如果出现故障,数据只能恢复到'最近一个成功点'",
	"rdbchecksum": "是否进行校验和，是否对rdb文件使用CRC64校验和,默认为'yes'，那么每个rdb文件内容的末尾都会追加CRC校验和，利于第三方校验工具检测文件完整性",
	"slave-serve-stale-data": "当主master服务器挂机或主从复制在进行时，是否依然可以允许客户访问可能过期的数据。在'yes'情况下,slave继续向客户端提供只读服务," +
		"有可能此时的数据已经过期；在'no'情况下，任何向此server发送的数据请求服务(包括客户端和此server的slave)都将被告知'error'",
	"slave-read-only":          "slave是否为'只读'，强烈建议为'yes'",
	"repl-ping-slave-period":   "slave向指定的master发送ping消息的时间间隔(秒)，默认为10",
	"repl-timeout":             "slave与master通讯中,最大空闲时间,默认60秒.超时将导致连接关闭",
	"repl-disable-tcp-nodelay": "slave与master的连接,是否禁用TCP nodelay选项。'yes'表示禁用,那么socket通讯中数据将会以packet方式发送(packet大小受到socket buffer限制)。",
	"slave-priority": "适用Sentinel模块(unstable,M-S集群管理和监控),需要额外的配置文件支持。slave的权重值,默认100.当master失效后,Sentinel将会从slave列表中找到权重" +
		"值最低(>0)的slave,并提升为master。如果权重值为0,表示此slave为'观察者',不参与master选举",
	"rename-command": "重命名指令,对于一些与'server'控制有关的指令,可能不希望远程客户端(非管理员用户)链接随意使用,那么就可以把这些指令重命名为'难以阅读'的其他字符串",
	"maxmemory-policy": "内存不足时,数据清除策略,默认为'volatile-lru'。volatile-lru  ->对'过期集合'中的数据采取LRU(近期最少使用)算法.如果对key使用'expire'指令指" +
		"定了过期时间,那么此key将会被添加到'过期集合'中。将已经过期/LRU的数据优先移除.如果'过期集合'中全部移除仍不能满足内存需求,将OOM; allkeys-lru ->对所有的数据,采用LRU算法; " +
		"volatile-random ->对'过期集合'中的数据采取'随即选取'算法,并移除选中的K-V,直到'内存足够'为止. 如果如果'过期集合'中全部移除全部移除仍不能满足,将OOM; " +
		"allkeys-random ->对所有的数据,采取'随机选取'算法,并移除选中的K-V,直到'内存足够'为止; volatile-ttl ->对'过期集合'中的数据采取TTL算法(最小存活时间),移除即将过期的数据; " +
		"noeviction ->不做任何干扰操作,直接返回OOM异常",
	"maxmemory-samples":           "默认值3，上面LRU和最小TTL策略并非严谨的策略，而是大约估算的方式，因此可以选择取样值以便检查",
	"no-appendfsync-on-rewrite":   "在aof rewrite期间,是否对aof新记录的append暂缓使用文件同步策略,主要考虑磁盘IO开支和请求阻塞时间。默认为no,表示'不暂缓',新的aof记录仍然会被立即同步",
	"auto-aof-rewrite-percentage": "当Aof log增长超过指定比例时，重写log file， 设置为0表示不自动重写Aof 日志，重写是为了使aof体积保持最小，而确保保存最完整的数据.",
	"auto-aof-rewrite-min-size":   "触发aof rewrite的最小文件尺寸",
	"lua-time-limit":              "lua脚本运行的最大时间",
	"slowlog-log-slower-than": "'慢操作日志'记录,单位:微秒(百万分之一秒,1000 * 1000),如果操作时间超过此值,将会把command信息'记录'起来.(内存,非文件)。其中'操作时间'不包括网络IO开支," +
		"只包括请求达到server后进行'内存实施'的时间.'0'表示记录全部操作",
	"slowlog-max-len": "'慢操作日志'保留的最大条数,'记录'将会被队列化,如果超过了此长度,旧记录将会被移除。可以通过'SLOWLOG <subcommand> args'查看慢记录的信息(SLOWLOG get 10,SLOWLOG reset)",
	"hash-max-ziplist-entries": "hash类型的数据结构在编码上可以使用ziplist和hashtable。ziplist的特点就是文件存储(以及内存存储)所需的空间较小,在内容较小时,性能和hashtable几乎一样." +
		"因此redis对hash类型默认采取ziplist。如果hash中条目的条目个数或者value长度达到阀值,将会被重构为hashtable。",
	"hash-max-ziplist-value":   "ziplist中允许条目value值最大字节数，默认为64，建议为1024",
	"list-max-ziplist-entries": "对于list类型,将会采取ziplist,linkedlist两种编码类型。",
	"list-max-ziplist-value":   "对于list类型,将会采取ziplist,linkedlist两种编码类型。",
	"set-max-intset-entries":   "intset中允许保存的最大条目个数,如果达到阀值,intset将会被重构为hashtable",
	"zset-max-ziplist-entries": "zset为有序集合,有2中编码类型:ziplist,skiplist。因为'排序'将会消耗额外的性能,当zset中数据较多时,将会被重构为skiplist。",
	"zset-max-ziplist-value":   "zset为有序集合,有2中编码类型:ziplist,skiplist。因为'排序'将会消耗额外的性能,当zset中数据较多时,将会被重构为skiplist。",
	"client-output-buffer-limit": "客户端buffer控制。在客户端与server进行的交互中,每个连接都会与一个buffer关联,此buffer用来队列化等待被client接受的响应信息。如果client不能及时的消费响应信息," +
		"那么buffer将会被不断积压而给server带来内存压力.如果buffer中积压的数据达到阀值,将会导致连接被关闭,buffer被移除。",
	"cluster-announce-ip":           "群集公告IP",
	"proto-max-bulk-len":            "客户端查询的缓存极限值大小",
	"client-query-buffer-limit":     "对于pubsub client，如果client-output-buffer一旦超过32mb，又或者超过8mb持续60秒，那么服务器就会立即断开客户端连接",
	"active-defrag-threshold-lower": "动活动碎片整理的最小碎片百分比",
	"active-defrag-threshold-upper": "使用最大努力的最大碎片百分比",
	"active-defrag-ignore-bytes":    "启动活动碎片整理的最小碎片消费量",
	"active-defrag-cycle-min":       "以CPU百分比表示的碎片整理的最小工作量",
	"active-defrag-cycle-max":       "在CPU的百分比最大的努力和碎片整理",
	"active-defrag-max-scan-fields": "将从中处理的set/hash/zset/list字段的最大数目,主词典扫描",
	"stream-node-max-bytes": "宏观节点的最大流。在流数据结构是一个基数树节点编码在这项大的多。利用这个配置它是如何可能#大节点配置是单字节和最大项目数，这可能包含了在切换到新节点的时候appending新的流条目。" +
		"如果任何以下设置来设置ignored极限是零，例如，操作系统，它有可能只是一集通过设置限制最大#纪录到最大字节0和最大输入到所需的值",
	"stream-node-max-entries": "宏观节点的最小流",
	"list-max-ziplist-size": "压缩列表大小 	-5:最大大小：64 KB<--不建议用于正常工作负载	-4:最大大小：32 KB<--不推荐" +
		"-3:最大大小：16 KB<--可能不推荐		-2:最大大小：8kb<--良好	-1:最大大小：4kb<--良好",
	"list-compress-depth": "0:禁用所有列表压缩		1：深度1表示“在列表中的1个节点之后才开始压缩，从头部或尾部所以：【head】->node->node->…->node->【tail】[头部]，[尾部]将始终未压缩；内部节点将压缩。" +
		"2:[头部]->[下一步]->节点->节点->…->节点->[上一步]->[尾部]	2这里的意思是：不要压缩头部或头部->下一个或尾部->上一个或尾部，但是压缩它们之间的所有节点。" +
		"3:[头部]->[下一步]->[下一步]->节点->节点->…->节点->[上一步]->[上一步]->[尾部]",
	"hll-sparse-max-bytes": "value大小小于等于hll-sparse-max-bytes使用稀疏数据结构（sparse），大于hll-sparse-max-bytes使用稠密的数据结构（dense）。一个比16000大的value是几乎没用的，" +
		"建议的value大概为3000。如果对CPU要求不高，对空间要求较高的，建议设置到10000左右",
	"latency-monitor-threshold": "延迟监控功能是用来监控redis中执行比较缓慢的一些操作，用LATENCY打印redis实例在跑命令时的耗时图表。" +
		"只记录大于等于下边设置的值的操作。0的话，就是关闭监视。默认延迟监控功能是关闭的，如果你需要打开，也可以通过CONFIG SET命令动态设置",
	"cluster-announce-bus-port": "群集公告总线端口",
	"hz":                        "redis执行任务的频率为1s除以hz",
	"cluster-node-timeout":      "节点互连超时的阀值。集群节点超时毫秒数",
	"cluster-migration-barrier": "master的slave数量大于该值，slave才能迁移到其他孤立master上，如这个参数若被设为2，那么只有当一个主节点拥有2个可工作的从节点时，它的一个从节点会尝试迁移",
	"cluster-replica-validity-factor": "在进行故障转移的时候，全部slave都会请求申请为master，但是有些slave可能与master断开连接一段时间了，导致数据过于陈旧，这样的slave不应该被提升为master。" +
		"该参数就是用来判断slave节点与master断线的时间是否过长。判断方法是：比较slave断开连接的时间和(node-timeout * slave-validity-factor) + repl-ping-slave-period",
	"cluster-require-full-coverage": "默认情况下，集群全部的slot有节点负责，集群状态才为ok，才能提供服务。设置为no，可以在slot没有全部分配的时候提供服务。不建议打开该配置，这样会造成分区的时候，" +
		"小分区的master一直在接受写请求，而造成很长时间数据不一致",
	"activedefrag":                  "已启用活动碎片整理",
	"aof-rewrite-incremental-fsync": "在aof重写的时候，如果打开了aof-rewrite-incremental-fsync开关，系统会每32MB执行一次fsync。这对于把文件写入磁盘是有帮助的，可以避免过大的延迟峰值",
	"rdb-save-incremental-fsync":    "在rdb保存的时候，如果打开了rdb-save-incremental-fsync开关，系统会每32MB执行一次fsync。这对于把文件写入磁盘是有帮助的，可以避免过大的延迟峰值",
	"aof-load-truncated": "aof文件可能在尾部是不完整的，当redis启动的时候，aof文件的数据被载入内存。重启可能发生在redis所在的主机操作系统宕机后，尤其在ext4文件系统没有加上data=ordered选项" +
		"（redis宕机或者异常终止不会造成尾部不完整现象。）出现这种现象，可以选择让redis退出，或者导入尽可能多的数据。如果选择的是yes，当截断的aof文件被导入的时候，会自动发布一个log给客户端然后load。" +
		"如果是no，用户必须手动redis-check-aof修复AOF文件才可以",
	"aof-use-rdb-preamble":          "加载redis时，可以识别AOF文件以“redis”开头。字符串并加载带前缀的RDB文件，然后继续加载AOF尾巴",
	"dynamic-hz":                    "当启用动态执行任务频率时，实际配置的hz将用作作为基线，但实际配置的hz值的倍数，在连接更多客户端后根据需要使用。这样一个闲置的实例将占用很少的CPU时间，而繁忙的实例将反应更灵敏",
	"notify-keyspace-events":        "键空间通知使得客户端可以通过订阅频道或模式，来接收那些以某种方式改动了Redis数据集的事件。因为开启键空间通知功能需要消耗一些CPU，所以在默认配置下，该功能处于关闭状态。",
	"unixsocket":                    "unix socket通信",
	"slave-announce-ip":             "slave 播报ip",
	"replica-announce-ip":           "replica 播报ip",
	"lfu-log-factor":                "LFU算法  log factor",
	"lfu-decay-time":                "LFU算法 decay time",
	"cluster-announce-port":         "集群公告端口",
	"tcp-backlog":                   "tcp backlog",
	"repl-ping-replica-period":      "复制周期",
	"repl-backlog-size":             "repl backlog 大小",
	"repl-backlog-ttl":              "repl backlog ttl",
	"watchdog-period":               "监督期",
	"replica-priority":              "复制优先级",
	"slave-announce-port":           "slave 播报prot",
	"replica-announce-port":         "replica 播报prot",
	"min-slaves-to-write":           "最小写入slaves",
	"min-replicas-to-write":         "最小写入replicas",
	"min-slaves-max-lag":            "min-slaves-max-lag",
	"min-replicas-max-lag":          "min-replicas-max-lag",
	"cluster-slave-validity-factor": "slave集群流动性系数",
	"repl-diskless-sync-delay":      "repl 磁盘同步延迟",
	"cluster-slave-no-failover":     "slave集群无故障转移",
	"cluster-replica-no-failover":   "replica集群无故障转移",
	"replica-serve-stale-data":      "replica提供陈旧数据",
	"replica-read-only":             "replica只读",
	"slave-ignore-maxmemory":        "slave最大存储容量",
	"replica-ignore-maxmemory":      "replica最大存储容量",
	"protected-mode":                "保护模式",
	"repl-diskless-sync":            "repl磁盘同步",
	"lazyfree-lazy-eviction":        "lazyfree懒惰驱逐",
	"lazyfree-lazy-expire":          "lazyfree懒惰到期",
	"lazyfree-lazy-server-del":      "lazyfree懒惰服务器del",
	"slave-lazy-flush":              "slave懒惰冲洗",
	"replica-lazy-flush":            "replica懒惰冲洗",
	"supervised":                    "监督",
	"syslog-facility":               "系统日志记录的级别",
	"unixsocketperm":                "unix socket 随机存取",
}
