//
//	redis 连接服务
//
package service

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/mangenotwork/mange_redis_manage/common"
	"github.com/mangenotwork/mange_redis_manage/common/manlog"
	manredis "github.com/mangenotwork/mange_redis_manage/common/redis"
	"github.com/mangenotwork/mange_redis_manage/dao"
	"github.com/mangenotwork/mange_redis_manage/repository"
	"github.com/mangenotwork/mange_redis_manage/structs"
)

type RedisConnService interface {
	New(data *structs.RedisConnData, uid int64) string                                                                                   //新的连接
	Detection(data *structs.RedisConnData) string                                                                                        //测试连接
	GetAll(user *structs.UserParameter) (datas *structs.RedisConnList)                                                                   //所有连接列表
	Info()                                                                                                                               //连接信息
	Modify()                                                                                                                             //修改连接
	GetRedisInfos(user *structs.UserParameter, redisId int64) (data []*structs.ResponseKV)                                               //获取服务的信息
	GetAllClient(user *structs.UserParameter, redisId int64) (datas []*dao.RedisClientInfo)                                              //获取所有连接redis的客户端
	GetMemoryChartData(user *structs.UserParameter, redisId int64) (memory_show *structs.RedisMemoryShow, err error)                     //获取内存图表数据
	GetEchartsRedisMemoryData(user *structs.UserParameter, redisId, hours, day int64) (datas *structs.EchartsRedisMemoryData, err error) //获取Redis服务信息图表数据
	GetRedisDBTree(user *structs.UserParameter, redisId int64) (dblist []*structs.Tree)                                                  //获取db树
	GetRedisDBList(user *structs.UserParameter, redisId int64) (dblist []*structs.RedisServersDBInfo)                                    //获取db列表
	GetRedisKeyTree(user *structs.UserParameter, redisId, dbId int64, match string) (dblist []*structs.Tree)                             //获取keys数
	GetRedisKeySearch(user *structs.UserParameter, redisId, dbId int64, match string) (dblist []*structs.Tree)                           //搜索keys
}

func RedisConnServiceFunc() RedisConnService {
	return new(RedisConn)
}

type RedisConn struct {
}

//私有方法 redis连接
func (this *RedisConn) conn(data *structs.RedisConnData) (rc redis.Conn, connErr error) {
	var isSSH bool = false

	if data.RedisHost == "" {
		connErr = errors.New("redis host 为空")
		return
	}

	if data.RedisPort == 0 {
		data.RedisPort = 6379
	}

	if data.SSHHost != "" {
		isSSH = true
	}

	manlog.Debug("ssh = ", data.SSHHost, data.SSHPassword, data.SSHUser)

	if isSSH && data.SSHUser == "" {
		data.SSHUser = "root"
	}

	if isSSH {
		rc, connErr = manredis.RSSHConn(data.SSHUser, data.SSHPassword, data.SSHHost, data.RedisHost, data.RedisPort, data.RedisPassword)
	}

	rc, connErr = manredis.RConn(data.RedisHost, data.RedisPort, data.RedisPassword)
	return
}

//获取指定db的连接
func (this *RedisConn) DBconn(user *structs.UserParameter, redisId, dbId int64) (rc redis.Conn, rcid string, err error) {
	//通过uid,redisid 取连接信息
	redisconn, err := new(dao.DaoRedisInfo).GetConnInfo(user.UserID, redisId)
	if err != nil {
		manlog.Error("获取连接信息错误", err)
		return
	}
	rcid = ""

	rc, err = this.conn(&structs.RedisConnData{
		RedisHost:     redisconn.ConnHost,
		RedisPort:     redisconn.ConnPort,
		RedisPassword: redisconn.ConnPassword,
		SSHHost:       redisconn.SSHUrl,
		SSHUser:       redisconn.SSHUser,
		SSHPassword:   redisconn.SSHPassword,
	})
	rcid = common.Get16MD5Encode(fmt.Sprintf("%s:%d", redisconn.ConnHost, redisconn.ConnPort))
	if dbId != 0 {
		rc, err = manredis.SelectDB(rc, dbId)
	}

	return
}

//通过用户信息与连接id获取redis连接
func (this *RedisConn) getconn(user *structs.UserParameter, redisId int64) (rc redis.Conn, err error) {
	//通过uid,redisid 取连接信息
	redisconn, err := new(dao.DaoRedisInfo).GetConnInfo(user.UserID, redisId)
	if err != nil {
		manlog.Error("获取连接信息错误", err)
		return
	}
	rc, err = this.conn(&structs.RedisConnData{
		RedisHost:     redisconn.ConnHost,
		RedisPort:     redisconn.ConnPort,
		RedisPassword: redisconn.ConnPassword,
		SSHHost:       redisconn.SSHUrl,
		SSHUser:       redisconn.SSHUser,
		SSHPassword:   redisconn.SSHPassword,
	})
	if err != nil {
		manlog.Error("连接redis失败", err)
		return
	}
	return
}

//通过用户信息与连接id获取redis连接 并返回连接id
func (this *RedisConn) getconnID(user *structs.UserParameter, redisId int64) (rc redis.Conn, rcid string, err error) {
	//通过uid,redisid 取连接信息
	redisconn, err := new(dao.DaoRedisInfo).GetConnInfo(user.UserID, redisId)
	if err != nil {
		manlog.Error("获取连接信息错误", err)
		return
	}
	rcid = ""
	rc, err = this.conn(&structs.RedisConnData{
		RedisHost:     redisconn.ConnHost,
		RedisPort:     redisconn.ConnPort,
		RedisPassword: redisconn.ConnPassword,
		SSHHost:       redisconn.SSHUrl,
		SSHUser:       redisconn.SSHUser,
		SSHPassword:   redisconn.SSHPassword,
	})
	if err != nil {
		manlog.Error("连接redis失败", err)
		return
	}
	//由redis host+port组成redis id
	rcid = common.Get16MD5Encode(fmt.Sprintf("%s:%d", redisconn.ConnHost, redisconn.ConnPort))
	return
}

//提供给其他方法使用的连接
func (this *RedisConn) RredisConn(conn *models.RedisInfoDB) (rc redis.Conn, rcid string, err error) {
	rc, err = this.conn(&structs.RedisConnData{
		RedisHost:     conn.ConnHost,
		RedisPort:     conn.ConnPort,
		RedisPassword: conn.ConnPassword,
		SSHHost:       conn.SSHUrl,
		SSHUser:       conn.SSHUser,
		SSHPassword:   conn.SSHPassword,
	})
	rcid = ""
	if err == nil {
		rcid = common.Get16MD5Encode(fmt.Sprintf("%s:%d", conn.ConnHost, conn.ConnPort))
	}
	return
}

//创建新的连接
func (this *RedisConn) New(data *structs.RedisConnData, uid int64) string {
	manlog.Debug("this is RedisConn.New()")
	rc, err := this.conn(data)
	if err != nil {
		manlog.Error(err)
		return "连接失败! err=" + err.Error()
	}

	//保存连接信息
	var redisinfo dao.DaoRedisInfo
	redisinfo.UID = uid
	redisinfo.ConnName = data.ConnName
	manlog.Debug(data.ConnName, redisinfo.ConnName)
	redisinfo.ConnHost = data.RedisHost
	redisinfo.ConnPort = data.RedisPort
	redisinfo.ConnPassword = data.RedisPassword
	redisinfo.IsSSH = false
	redisinfo.SSHUrl = ""
	redisinfo.SSHUser = ""
	redisinfo.SSHPassword = ""
	redisinfo.ConnCreate = time.Now().Unix()
	err = redisinfo.Create()
	if err != nil {
		return "连接失败! err=" + err.Error()
	}
	manlog.Debug(rc)
	return "连接成功"
}

//测试连接
func (this *RedisConn) Detection(data *structs.RedisConnData) string {
	if _, err := this.conn(data); err != nil {
		manlog.Error(err)
		return "连接失败! err=" + err.Error()
	}

	return "连接成功"
}

//获取所有连接
func (this *RedisConn) GetAll(user *structs.UserParameter) (datas *structs.RedisConnList) {
	datas = new(structs.RedisConnList)
	manlog.Debug("this is RedisConn.GetAll()")
	var redisinfo dao.DaoRedisInfo
	conndatas := make([]*structs.RedisConnInfo, 0)
	connlist, _ := redisinfo.GetAll(user.UserID)
	datas.Count = 0
	for _, v := range connlist {
		conndatas = append(conndatas, &structs.RedisConnInfo{
			ConnId:     v.ID,
			ConnName:   v.ConnName,
			RedisConn:  fmt.Sprintf("%s:%d", v.ConnHost, v.ConnPort),
			ConnCreate: time.Unix(v.ConnCreate, 0).Format("2006-01-02 15:04:05"),
		})
		datas.Count++
	}
	datas.List = conndatas
	return
}

func (this *RedisConn) Info() {}

func (this *RedisConn) Modify() {}

//获取redis 服务器信息
func (this *RedisConn) GetRedisInfos(user *structs.UserParameter, redisId int64) (data []*structs.ResponseKV) {
	rc, rcid, err := this.getconnID(user, redisId)
	if err != nil {
		manlog.Error(err)
	}
	datas := manredis.GetRedisServersInfos(rc, rcid)

	//Redis 服务器运行模式
	redis_mode_value := "集群"
	if datas.Server.RedisMode == "standalone" {
		redis_mode_value = "单机"
	}

	// data = new(structs.RedisServersInfo)
	// data.BaseInfo = &structs.RedisServersBaseInfo{
	// 	RedisVersion:                fmt.Sprintf("Redis 服务器版本 ： %s", datas.Server.RedisVersion),
	// 	RedisMode:                   fmt.Sprintf("Redis 服务器运行模式 ： %s", redis_mode_value),
	// 	Os:                          fmt.Sprintf("宿主操作系统 ： %s", datas.Server.Os),
	// 	RedisBuildId:                fmt.Sprintf("Redis build id ： %s", datas.Server.RedisBuildId),
	// 	ArchBits:                    fmt.Sprintf("架构 ： %s位", datas.Server.ArchBits),
	// 	MultiplexingApi:             fmt.Sprintf("redis所使用的事件处理模型 ： %s", datas.Server.MultiplexingApi),
	// 	GccVersion:                  fmt.Sprintf("编译redis时gcc版本 ： %s", datas.Server.GccVersion),
	// 	ProcessId:                   fmt.Sprintf("redis服务器进程的pid ： %s", datas.Server.ProcessId),
	// 	UptimeInSeconds:             fmt.Sprintf("redis服务器启动总时间 ： %d秒", datas.Server.UptimeInSeconds),
	// 	UptimeInDays:                fmt.Sprintf("redis服务器启动总时间 ： %d天", datas.Server.UptimeInDays),
	// 	Hz:                          fmt.Sprintf("redis内部调度频率 ： %dHz", datas.Server.Hz),
	// 	ConfigFile:                  fmt.Sprintf("redis配置路径 ： %s", datas.Server.ConfigFile),
	// 	ConnectedClients:            fmt.Sprintf("已经连接客户端数量 ： %d", datas.Clients.ConnectedClients),
	// 	ClientRecentMaxInputBuffer:  fmt.Sprintf("客户端最近最大输入缓冲区 ： %d", datas.Clients.ClientRecentMaxInputBuffer),
	// 	ClientRecentMaxOutputBuffer: fmt.Sprintf("客户端最近最大输出缓冲区 ： %d", datas.Clients.ClientRecentMaxOutputBuffer),
	// 	BlockedClients:              fmt.Sprintf("正在等待阻塞命令的客户端数量: %d", datas.Clients.BlockedClients),
	// 	Loading:                     fmt.Sprintf("服务器是否正在载入持久化文件: %d", datas.Persistence.Loading),
	// 	RdbChangesSinceLastSave:     fmt.Sprintf("已经写入的命令还未被持久化: %d", datas.Persistence.RdbChangesSinceLastSave),
	// 	RdbBgsaveInProgress:         fmt.Sprintf("服务器是否正在创建rdb文件: %d", datas.Persistence.RdbBgsaveInProgress),
	// 	RdbLastSaveTime:             fmt.Sprintf("多长时间未进行持久化: %d", datas.Persistence.RdbLastSaveTime),
	// 	RdbLastBgsaveStatus:         fmt.Sprintf("最后一次的rdb持久化是否成功: %s", datas.Persistence.RdbLastBgsaveStatus),
	// 	RdbLastBgsaveTimeSec:        fmt.Sprintf("最后一次生成rdb文件耗时秒数: %d", datas.Persistence.RdbLastBgsaveTimeSec),
	// 	AofEnabled:                  fmt.Sprintf("是否开启了aof: %d", datas.Persistence.AofEnabled),
	// 	AofRewriteInProgress:        fmt.Sprintf("标识aof的rewrite操作是否进行中: %d", datas.Persistence.AofRewriteInProgress),
	// 	AofLastWriteStatus:          fmt.Sprintf("上一次aof写入状态: %s", datas.Persistence.AofLastWriteStatus),
	// 	TotalConnectionsReceived:    fmt.Sprintf("新创建的链接个数: %d", datas.Stats.TotalConnectionsReceived),
	// 	TotalCommandsProcessed:      fmt.Sprintf("redis处理的命令数: %d", datas.Stats.TotalCommandsProcessed),
	// 	InstantaneousOpsPerSec:      fmt.Sprintf("redis当前的qps(每秒执行命令数): %d", datas.Stats.InstantaneousOpsPerSec),
	// 	TotalNetInputBytes:          fmt.Sprintf("redis网络入口流量字节数: %d", datas.Stats.TotalNetInputBytes),
	// 	TotalNetOutputBytes:         fmt.Sprintf("redis网络出口流量字节数: %d", datas.Stats.TotalNetOutputBytes),
	// 	InstantaneousInputKbps:      fmt.Sprintf("redis网络入口kps: %v", datas.Stats.InstantaneousInputKbps),
	// 	InstantaneousOutputKbps:     fmt.Sprintf("redis网络出口kps: %v", datas.Stats.InstantaneousOutputKbps),
	// 	RejectedConnections:         fmt.Sprintf("拒绝的连接个数: %d", datas.Stats.RejectedConnections),
	// 	SyncFull:                    fmt.Sprintf("主从完全同步成功次数: %d", datas.Stats.SyncFull),
	// 	SyncPartialOk:               fmt.Sprintf("主从部分同步成功次数: %d", datas.Stats.SyncPartialOk),
	// 	SyncPartialErr:              fmt.Sprintf("主从部分同步失败次数: %d", datas.Stats.SyncPartialErr),
	// 	ExpiredKeys:                 fmt.Sprintf("运行以来过期的key的数量: %d", datas.Stats.ExpiredKeys),
	// 	EvictedKeys:                 fmt.Sprintf("运行以来剔除（超过maxmemory）的key的数量: %d", datas.Stats.EvictedKeys),
	// 	KeyspaceHits:                fmt.Sprintf("命中次数: %d", datas.Stats.KeyspaceHits),
	// 	KeyspaceMisses:              fmt.Sprintf("没命中次数: %d", datas.Stats.KeyspaceMisses),
	// 	PubsubChannels:              fmt.Sprintf("当前使用中的频道数量: %d", datas.Stats.PubsubChannels),
	// 	PubsubPatterns:              fmt.Sprintf("当前使用的模式数量: %d", datas.Stats.PubsubPatterns),
	// }

	kv_datas := make([]*structs.ResponseKV, 0)
	kv_datas = append(kv_datas, &structs.ResponseKV{"Redis 服务器版本", datas.Server.RedisVersion})
	kv_datas = append(kv_datas, &structs.ResponseKV{"Redis 服务器运行模式", redis_mode_value})
	kv_datas = append(kv_datas, &structs.ResponseKV{"宿主操作系统", datas.Server.Os})
	kv_datas = append(kv_datas, &structs.ResponseKV{"Redis build id", datas.Server.RedisBuildId})
	kv_datas = append(kv_datas, &structs.ResponseKV{"架构(位)", datas.Server.ArchBits})
	kv_datas = append(kv_datas, &structs.ResponseKV{"redis所使用的事件处理模型", datas.Server.MultiplexingApi})
	kv_datas = append(kv_datas, &structs.ResponseKV{"编译redis时gcc版本", datas.Server.GccVersion})
	kv_datas = append(kv_datas, &structs.ResponseKV{"redis服务器进程的pid", datas.Server.ProcessId})
	//kv_datas = append(kv_datas, &structs.ResponseKV{"redis服务器启动总时间(秒)", datas.Server.UptimeInSeconds})
	kv_datas = append(kv_datas, &structs.ResponseKV{"redis服务器启动总时间(天)", datas.Server.UptimeInDays})
	kv_datas = append(kv_datas, &structs.ResponseKV{"redis内部调度频率(Hz)", datas.Server.Hz})
	kv_datas = append(kv_datas, &structs.ResponseKV{"redis配置路径", datas.Server.ConfigFile})
	//kv_datas = append(kv_datas, &structs.ResponseKV{"已经连接客户端数量", datas.Clients.ConnectedClients})
	kv_datas = append(kv_datas, &structs.ResponseKV{"客户端最近最大输入缓冲区", datas.Clients.ClientRecentMaxInputBuffer})
	kv_datas = append(kv_datas, &structs.ResponseKV{"客户端最近最大输出缓冲区", datas.Clients.ClientRecentMaxOutputBuffer})
	kv_datas = append(kv_datas, &structs.ResponseKV{"正在等待阻塞命令的客户端数量", datas.Clients.BlockedClients})
	kv_datas = append(kv_datas, &structs.ResponseKV{"服务器是否正在载入持久化文件", datas.Persistence.Loading})
	kv_datas = append(kv_datas, &structs.ResponseKV{"已经写入的命令还未被持久化", datas.Persistence.RdbChangesSinceLastSave})
	kv_datas = append(kv_datas, &structs.ResponseKV{"服务器是否正在创建rdb文件", datas.Persistence.RdbBgsaveInProgress})
	kv_datas = append(kv_datas, &structs.ResponseKV{"多长时间未进行持久化", datas.Persistence.RdbLastSaveTime})
	kv_datas = append(kv_datas, &structs.ResponseKV{"最后一次的rdb持久化是否成功", datas.Persistence.RdbLastBgsaveStatus})
	kv_datas = append(kv_datas, &structs.ResponseKV{"最后一次生成rdb文件耗时秒数", datas.Persistence.RdbLastBgsaveTimeSec})
	kv_datas = append(kv_datas, &structs.ResponseKV{"是否开启了aof", datas.Persistence.AofEnabled})
	kv_datas = append(kv_datas, &structs.ResponseKV{"标识aof的rewrite\n操作是否进行中", datas.Persistence.AofRewriteInProgress})
	kv_datas = append(kv_datas, &structs.ResponseKV{"上一次aof写入状态", datas.Persistence.AofLastWriteStatus})
	kv_datas = append(kv_datas, &structs.ResponseKV{"新创建的链接个数", datas.Stats.TotalConnectionsReceived})
	//kv_datas = append(kv_datas, &structs.ResponseKV{"redis处理的命令数", datas.Stats.TotalCommandsProcessed})
	kv_datas = append(kv_datas, &structs.ResponseKV{"redis当前的qps\n(每秒执行命令数)", datas.Stats.InstantaneousOpsPerSec})
	kv_datas = append(kv_datas, &structs.ResponseKV{"redis网络入口流量字节数", datas.Stats.TotalNetInputBytes})
	kv_datas = append(kv_datas, &structs.ResponseKV{"redis网络出口流量字节数", datas.Stats.TotalNetOutputBytes})
	kv_datas = append(kv_datas, &structs.ResponseKV{"redis网络入口kps", datas.Stats.InstantaneousInputKbps})
	kv_datas = append(kv_datas, &structs.ResponseKV{"redis网络出口kps", datas.Stats.InstantaneousOutputKbps})
	kv_datas = append(kv_datas, &structs.ResponseKV{"拒绝的连接个数", datas.Stats.RejectedConnections})
	kv_datas = append(kv_datas, &structs.ResponseKV{"主从完全同步成功次数", datas.Stats.SyncFull})
	kv_datas = append(kv_datas, &structs.ResponseKV{"主从部分同步成功次数", datas.Stats.SyncPartialOk})
	kv_datas = append(kv_datas, &structs.ResponseKV{"主从部分同步失败次数", datas.Stats.SyncPartialErr})
	kv_datas = append(kv_datas, &structs.ResponseKV{"运行以来过期的key的数量", datas.Stats.ExpiredKeys})
	kv_datas = append(kv_datas, &structs.ResponseKV{"运行以来剔除\n（超过maxmemory）的key的数量", datas.Stats.EvictedKeys})
	kv_datas = append(kv_datas, &structs.ResponseKV{"命中次数", datas.Stats.KeyspaceHits})
	kv_datas = append(kv_datas, &structs.ResponseKV{"没命中次数", datas.Stats.KeyspaceMisses})
	kv_datas = append(kv_datas, &structs.ResponseKV{"当前使用中的频道数量", datas.Stats.PubsubChannels})
	kv_datas = append(kv_datas, &structs.ResponseKV{"当前使用的模式数量", datas.Stats.PubsubPatterns})

	// dbinfos := make([]*structs.RedisServersDBInfo, 0)
	// for _, v := range datas.Keyspace {
	// 	dbinfos = append(dbinfos, &structs.RedisServersDBInfo{
	// 		DBID:    fmt.Sprintf("DB: %d", v.DBID),
	// 		Keys:    fmt.Sprintf("Key数量: %d", v.Keys),
	// 		Expires: fmt.Sprintf("过期数量: %d", v.Expires),
	// 		AvgTTL:  fmt.Sprintf("平均TTL: %d", v.AvgTTL),
	// 	})
	// }
	// data.DBInfo = dbinfos
	data = kv_datas
	return
}

func (this *RedisConn) RedisServserInfo2DB(datas *manredis.RedisServersInfo) {
	var mutex sync.Mutex

	dao_replication := &dao.DaoRedisReplication{datas.Replication, mutex}
	dao_replication.Create()

	dao_clients := &dao.DaoRedisClients{datas.Clients, mutex}
	dao_clients.Create()

	dao_cluster := &dao.DaoRedisCluster{datas.Cluster, mutex}
	dao_cluster.Create()

	dao_cpu := &dao.DaoRedisCPU{datas.CPU, mutex}
	dao_cpu.Create()

	dao_serverinfo := &dao.DaoRedisServerInfos{datas.Server, mutex}
	dao_serverinfo.Create()

	dao_memory := &dao.DaoRedisMemory{datas.Memory, mutex}
	dao_memory.Create()

	dao_persistence := &dao.DaoRedisPersistence{datas.Persistence, mutex}
	dao_persistence.Create()

	dao_stats := &dao.DaoRedisStats{datas.Stats, mutex}
	dao_stats.Create()

	for _, v := range datas.Keyspace {
		dao_db := &dao.DaoRedisKeyspace{v, mutex}
		dao_db.Create()
	}
}

func (this *RedisConn) GetAllClient(user *structs.UserParameter, redisId int64) (datas []*dao.RedisClientInfo) {
	rc, err := this.getconn(user, redisId)
	if err != nil {
		manlog.Error(err)
	}
	datas = manredis.GetAllRedisClient(rc)
	return
}

//获取内存图表数据
func (this *RedisConn) GetMemoryChartData(user *structs.UserParameter, redisId int64) (memory_show *structs.RedisMemoryShow, err error) {
	memory_show = new(structs.RedisMemoryShow)
	//通过user redisId 获取连接id
	_, rcid, err := this.getconnID(user, redisId)
	if err != nil {
		manlog.Error(err)
		return
	}

	//默认为1小时的周期
	nowtime := time.Now().Unix()
	ago_hour := nowtime - 3600

	//获取ago_hour 之后的所有get_time
	all_get_time, err := new(dao.DaoRedisMemory).GetAllGetTime(rcid, ago_hour)
	if all_get_time == nil || err != nil {
		return
	}

	get_time_number := len(all_get_time)
	add_number := get_time_number / 20
	var show_time []int64

	for k, v := range all_get_time {
		if k%add_number == 0 {
			//manlog.Error(v)
			show_time = append(show_time, v)
		}
	}

	//获取要显示的Memory数据
	show_memory_data, err := new(dao.DaoRedisMemory).GetShowMemory(rcid, show_time)
	if err != nil {
		return
	}

	memory_data := make([]*structs.RedisMemoryData, 0)
	for _, m := range show_memory_data {
		memory_data = append(memory_data, &structs.RedisMemoryData{
			Time:                  m.GetTime,
			TimeStr:               time.Unix(m.GetTime, 0).Format("2006-01-02 15:04:05"),
			UsedMemory:            m.UsedMemory,            //由redis分配器分配的内存总量，单位字节
			UsedMemoryHuman:       m.UsedMemoryHuman,       //
			UsedMemoryRss:         m.UsedMemoryRss,         //从操作系统角度，返回redis已分配内存总量
			UsedMemoryRssHuman:    m.UsedMemoryRssHuman,    //
			UsedMemoryPeak:        m.UsedMemoryPeak,        //redis的内存消耗峰值（以字节为单位）
			UsedMemoryPeakHuman:   m.UsedMemoryRssHuman,    //
			UsedMemoryLua:         m.UsedMemoryLua,         //lua引擎所使用的内存大小（单位字节）
			UsedMemoryLuaHuman:    m.UsedMemoryLuaHuman,    //
			MemFragmentationRatio: m.MemFragmentationRatio, //used_memory_rss 和 used_memory 之间的比率
		})
	}
	memory_base := &structs.RedisMemoryInfo{}
	if len(show_memory_data) > 0 {
		memory_base.MemAllocator = show_memory_data[0].MemAllocator
	}

	memory_show.BaseInfo = memory_base
	memory_show.Memory = memory_data

	return
}

func (this *RedisConn) GetEchartsRedisMemoryData(user *structs.UserParameter, redisId, hours, day int64) (datas *structs.EchartsRedisMemoryData, err error) {
	datas = new(structs.EchartsRedisMemoryData)
	//通过user redisId 获取连接id
	_, rcid, err := this.getconnID(user, redisId)
	if err != nil {
		manlog.Error(err)
		return
	}
	nowtime := time.Now().Unix()

	//时间差默认为 3600 1小时
	var sjc int64 = 3600

	if hours != 0 {
		sjc = sjc * hours
	}

	if day != 0 {
		sjc = sjc * day * 24
	}

	ago_hour := nowtime - sjc

	//获取ago_hour 之后的所有get_time
	all_get_time, err := new(dao.DaoRedisMemory).GetAllGetTime(rcid, ago_hour)
	if all_get_time == nil || err != nil {
		return
	}

	get_time_number := len(all_get_time)
	var show_time []int64
	if get_time_number > 20 {
		add_number := get_time_number / 20
		for k, v := range all_get_time {
			if k%add_number == 0 {
				//manlog.Error(v)
				show_time = append(show_time, v)
			}
		}
	} else {
		for _, v := range all_get_time {
			show_time = append(show_time, v)
		}
	}

	//获取要显示的Memory数据
	show_memory_data, err := new(dao.DaoRedisMemory).GetShowMemory(rcid, show_time)
	if err != nil {
		return
	}

	for _, m := range show_memory_data {
		datas.TimeList = append(datas.TimeList, time.Unix(m.GetTime, 0).Format("2006-01-02 15:04:05"))
		datas.UsedMemory = append(datas.UsedMemory, m.UsedMemory)
		datas.UsedMemoryLua = append(datas.UsedMemoryLua, m.UsedMemoryLua)
		datas.UsedMemoryPeak = append(datas.UsedMemoryPeak, m.UsedMemoryPeak)
		datas.UsedMemoryRss = append(datas.UsedMemoryRss, m.UsedMemoryRss)
	}

	if len(show_memory_data) > 0 {
		datas.UsedMemoryStr = show_memory_data[len(show_memory_data)-1].UsedMemoryHuman
	}

	//获取客户端数据
	clientsdata, err := new(dao.DaoRedisClients).GetNewData(rcid)
	if err != nil {
		manlog.Error(err)
	}
	datas.ClinetNumber = fmt.Sprintf("%d个", clientsdata.ConnectedClients)

	//获取执行命令的总数
	statsdata, err := new(dao.DaoRedisStats).GetNewData(rcid)
	if err != nil {
		manlog.Error(err)
	}
	datas.CmderNumber = fmt.Sprintf("%d条", statsdata.TotalCommandsProcessed)

	//获取运行时长
	serversinfo, err := new(dao.DaoRedisServerInfos).GetNewData(rcid)
	if err != nil {
		manlog.Error(err)
	}
	datas.RunTime = fmt.Sprintf("%d秒", serversinfo.UptimeInSeconds)

	//获取redisdb
	redisdbs := make([]*structs.RedisServersDBInfo, 0)
	maxgietime := show_time[len(show_time)-1]
	dbs, err := new(dao.DaoRedisKeyspace).GetNewData(rcid, maxgietime)
	if err != nil {
		manlog.Error(err)
	}
	for _, d := range dbs {
		redisdbs = append(redisdbs, &structs.RedisServersDBInfo{
			DBID:    d.DBID,
			Keys:    d.Keys,
			Expires: d.Expires,
			AvgTTL:  d.AvgTTL,
		})
	}
	datas.RedisDB = redisdbs

	return
}

func (this *RedisConn) GetRedisDBTree(user *structs.UserParameter, redisId int64) (dblist []*structs.Tree) {
	rc, rcid, err := this.getconnID(user, redisId)
	if err != nil {
		manlog.Error(err)
	}
	datas := manredis.GetRedisServersInfos(rc, rcid)
	manlog.Debug(datas.Keyspace)

	//dblist := make([]*structs.Tree, 0)

	for _, v := range datas.Keyspace {
		dblist = append(dblist, &structs.Tree{
			Text:     fmt.Sprintf("DB%d(%d)", v.DBID, v.Keys),
			State:    "closed",
			Children: make([]*structs.Tree, 0),
		})
	}

	return
}

func (this *RedisConn) GetRedisDBList(user *structs.UserParameter, redisId int64) (dblist []*structs.RedisServersDBInfo) {
	rc, rcid, err := this.getconnID(user, redisId)
	if err != nil {
		manlog.Error(err)
	}
	datas := manredis.GetRedisServersInfos(rc, rcid)
	manlog.Debug(datas.Keyspace)

	ismap := make(map[int64]bool, 0)

	for _, v := range datas.Keyspace {
		dblist = append(dblist, &structs.RedisServersDBInfo{
			DBID:    v.DBID,
			Keys:    v.Keys,
			Expires: v.Expires,
			AvgTTL:  v.AvgTTL,
		})
		ismap[v.DBID] = true
	}

	dbnum := manredis.GetDatabasesCount(rc)
	manlog.Debug("dbnum = ", dbnum)

	for i := 0; i <= dbnum; i++ {
		if !ismap[int64(i)] {
			dblist = append(dblist, &structs.RedisServersDBInfo{
				DBID: int64(i),
			})
		}
	}

	return

}

func (this *RedisConn) GetRedisKeyTree(user *structs.UserParameter, redisId, dbId int64, match string) (dblist []*structs.Tree) {
	rc, rcid, err := this.DBconn(user, redisId, dbId)
	if err != nil {
		manlog.Error(err)
	}
	manlog.Debug(rcid)
	data := manredis.GetALLKeys(rc, match)
	manlog.Debug(data)

	for k, v := range data {

		if v == 0 {
			dblist = append(dblist, &structs.Tree{
				Text: fmt.Sprintf("%s", k),
			})
		} else {
			dblist = append(dblist, &structs.Tree{
				Text:     fmt.Sprintf("%s(%d)", k, v),
				State:    "closed",
				Children: make([]*structs.Tree, 0),
			})
		}

	}
	return
}

func (this *RedisConn) GetRedisKeySearch(user *structs.UserParameter, redisId, dbId int64, match string) (dblist []*structs.Tree) {
	rc, rcid, err := this.DBconn(user, redisId, dbId)
	if err != nil {
		manlog.Error(err)
	}
	manlog.Debug(rcid)
	data := manredis.SearchKeys(rc, match)
	manlog.Debug(data)

	for k, v := range data {

		if v == 0 {
			dblist = append(dblist, &structs.Tree{
				Text: fmt.Sprintf("%s", k),
			})
		} else {
			dblist = append(dblist, &structs.Tree{
				Text:     fmt.Sprintf("%s(%d)", k, v),
				State:    "closed",
				Children: make([]*structs.Tree, 0),
			})
		}

	}
	return
}
