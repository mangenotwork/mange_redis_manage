//
//	redis 相关的结构体
//
package structs

//redis 连接接收参数
type RedisConnData struct {
	ConnName      string `json:"conn_name"`
	RedisHost     string `json:"redis_host"`
	RedisPort     int    `json:"redis_port"`
	RedisPassword string `json:"redis_password"`
	SSHHost       string `json:"ssh_host"`
	SSHUser       string `json:"ssh_user"`
	SSHPassword   string `json:"ssh_password"`
}

//redis 连接信息列表
type RedisConnList struct {
	Count int              `json:"count"`
	List  []*RedisConnInfo `json:"list"`
}

//redis 连接信息
type RedisConnInfo struct {
	ConnId     int64  `json:"conn_id"`
	ConnName   string `json:"conn_name"`
	RedisConn  string `json:"redis_conn"`
	ConnCreate string `json:"create"`
}

//redis 服务信息
type RedisServersInfo struct {
	BaseInfo *RedisServersBaseInfo `json:"redis_servers_baseinfo"`
	DBInfo   []*RedisServersDBInfo `json:"redis_servers_dbinfo"`
}

//redis 服务基础信息
type RedisServersBaseInfo struct {
	RedisVersion                string `json:"redis_version"`                   //Redis 服务器版本
	RedisMode                   string `json:"redis_mode"`                      //运行模式，单机或集群
	Os                          string `json:"os"`                              //Redis 服务器的宿主操作系统
	RedisBuildId                string `json:"redis_build_id"`                  //Redis build id
	ArchBits                    string `json:"arch_bits"`                       //架构64位
	MultiplexingApi             string `json:"multiplexing_api"`                //redis所使用的事件处理模型
	GccVersion                  string `json:"gcc_version"`                     //编译redis时gcc版本
	ProcessId                   string `json:"process_id"`                      //redis服务器进程的pid
	UptimeInSeconds             string `json:"uptime_in_seconds"`               //redis服务器启动总时间，单位秒
	UptimeInDays                string `json:"uptime_in_days"`                  //redis服务器启动总时间，单位天
	Hz                          string `json:"hz"`                              //redis内部调度频率（关闭timeout客户端，删除过期key）
	ConfigFile                  string `json:"config_file"`                     //配置文件路径
	ConnectedClients            string `json:"connected_clients"`               //已经连接客户端数量（不包括slave连接的客户端）
	ClientRecentMaxInputBuffer  string `json:"client_recent_max_input_buffer"`  //客户端最近最大输入缓冲区
	ClientRecentMaxOutputBuffer string `json:"client_recent_max_output_buffer"` //客户端最近最大输出缓冲区
	BlockedClients              string `json:"clocked_clients"`                 //正在等待阻塞命令的客户端数量
	Loading                     string `json:"loading"`                         //服务器是否正在载入持久化文件
	RdbChangesSinceLastSave     string `json:"rdb_changes_since_last_save"`     //有多少个已经写入的命令还未被持久化
	RdbBgsaveInProgress         string `json:"rdb_bgsave_in_progress"`          //服务器是否正在创建rdb文件
	RdbLastSaveTime             string `json:"rdb_last_save_time"`              //已经有多长时间没有进行持久化了
	RdbLastBgsaveStatus         string `json:"rdb_last_bgsave_status"`          //最后一次的rdb持久化是否成功
	RdbLastBgsaveTimeSec        string `json:"rdb_last_bgsave_time_sec"`        //最后一次生成rdb文件耗时秒数
	AofEnabled                  string `json:"aof_enabled"`                     //是否开启了aof
	AofRewriteInProgress        string `json:"aof_rewrite_in_progress"`         //标识aof的rewrite操作是否进行中
	AofLastWriteStatus          string `json:"aof_last_write_status"`           //上一次aof写入状态
	TotalConnectionsReceived    string `json:"total_connections_received"`      //新创建的链接个数，如果过多，会影响性能
	TotalCommandsProcessed      string `json:"total_commands_processed"`        //redis处理的命令数
	InstantaneousOpsPerSec      string `json:"instantaneous_ops_per_sec"`       //redis当前的qps，redis内部较实时的每秒执行命令数
	TotalNetInputBytes          string `json:"total_net_input_bytes"`           //redis网络入口流量字节数
	TotalNetOutputBytes         string `json:"total_net_output_bytes"`          //redis网络出口流量字节数
	InstantaneousInputKbps      string `json:"instantaneous_input_kbps"`        //redis网络入口kps
	InstantaneousOutputKbps     string `json:"instantaneous_output_kbps"`       //redis网络出口kps
	RejectedConnections         string `json:"rejected_connections"`            //拒绝的连接个数，redis连接个数已经达到maxclients限制。
	SyncFull                    string `json:"sync_full"`                       //主从完全同步成功次数
	SyncPartialOk               string `json:"sync_partial_ok"`                 //主从部分同步成功次数
	SyncPartialErr              string `json:"sync_partial_err"`                //主从部分同步失败次数
	ExpiredKeys                 string `json:"expired_keys"`                    //运行以来过期的key的数量
	EvictedKeys                 string `json:"evicted_keys"`                    //运行以来剔除（超过maxmemory）的key的数量s
	KeyspaceHits                string `json:"keyspace_hits"`                   //命中次数
	KeyspaceMisses              string `json:"keyspace_misses"`                 //没命中次数
	PubsubChannels              string `json:"pubsub_channels"`                 //当前使用中的频道数量
	PubsubPatterns              string `json:"pubsub_patterns"`                 //当前使用的模式数量
}

//redis 服务DB数据信息
type RedisServersDBInfo struct {
	DBID    int64 `json:"db"`
	Keys    int64 `json:"keys_count"`
	Expires int64 `json:"expires"`
	AvgTTL  int64 `json:"avgttl"`
}

//redis 内存使用数据
type RedisMemoryShow struct {
	Memory   []*RedisMemoryData `json:"memory"`
	BaseInfo *RedisMemoryInfo   `json:"base_info"`
}

//redis 内存数据
type RedisMemoryData struct {
	Time                  int64   `json:"time_unix"`
	TimeStr               string  `json:"time"`
	UsedMemory            int64   `json:"used_memory"`             //由redis分配器分配的内存总量，单位字节
	UsedMemoryHuman       string  `json:"used_memory_human"`       //
	UsedMemoryRss         int64   `json:"used_memory_rss"`         //从操作系统角度，返回redis已分配内存总量
	UsedMemoryRssHuman    string  `json:"used_memory_rss_human"`   //
	UsedMemoryPeak        int64   `json:"used_memory_peak"`        //redis的内存消耗峰值（以字节为单位）
	UsedMemoryPeakHuman   string  `json:"used_memory_peak_human"`  //
	UsedMemoryLua         int64   `json:"used_memory_lua"`         //lua引擎所使用的内存大小（单位字节）
	UsedMemoryLuaHuman    string  `json:"used_memory_lua_human"`   //
	MemFragmentationRatio float64 `json:"mem_fragmentation_ratio"` //used_memory_rss 和 used_memory 之间的比率
}

//redis 内存基本信息
type RedisMemoryInfo struct {
	MemAllocator string `json:"mem_allocator"` //编译时指定的redis的内存分配器。越好的分配器内存碎片化率越低，低版本建议升级
}

type EchartsRedisMemoryData struct {
	TimeList       []string              `json:"time_list"`
	UsedMemory     []int64               `json:"used_memory"`
	UsedMemoryRss  []int64               `json:"used_memory_rss"`
	UsedMemoryLua  []int64               `json:"used_memory_lua"`
	UsedMemoryPeak []int64               `json:"used_memory_peak"`
	UsedMemoryStr  string                `json:"used_memory_str"`
	ClinetNumber   string                `json:"clinet_number"`
	CmderNumber    string                `json:"cmder_number"`
	RunTime        string                `json:"run_time"`
	RedisDB        []*RedisServersDBInfo `json:"redis_dbs"`
}

type Tree struct {
	Text     string  `json:"text"`
	State    string  `json:"state,omitempty"`
	Children []*Tree `json:"children,omitempty"`
}

type DBTree struct {
	DBT []*Tree
}

//输出给监控的数据
type RealTime struct {
	RealTimeId     string    `json:"real_time_id"`
	Xdata          []string  `json:"time_data"`
	CPUTip         string    `json:"cpu_tip"`
	CPUData        []float64 `json:"cpu_data"`
	MemoryTip      string    `json:"memory_tip"`
	MemoryData     []float64 `json:"memory_data"`
	MemoryDW       string    `json:"memory_dw"`
	QpsTip         string    `json:"qps_tip"`
	QpsData        []int64   `json:"qps_data"`
	QpsDW          string    `json:"qps_dw"`
	ConnTip        string    `json:"conn_tip"`
	ConnData       []int64   `json:"conn_data"`
	KeysTip        string    `json:"keys_tip"`
	KeysData       []int64   `json:"keys_data"`
	InputKbpsTip   string    `json:"kbps_input_tip"`
	InputKbpsData  []float64 `json:"kbps_input_data"`
	OutputKbpsTip  string    `json:"kbps_output_tip"`
	OutputKbpsData []float64 `json:"kbps_output_data"`
	HitRateTip     string    `json:"hitrate_tip"`
	HitRateData    []int64   `json:"hitrate_data"`
}

//输出的keyinfo
type KeyInfo struct {
	KeyName string      `json:"key_name"`
	TTL     int64       `json:"ttl"`
	KeyType string      `json:"key_type"`
	KeySize string      `json:"size"`
	Value   interface{} `json:"value"`
	DBID    int64       `json:"key_db"`
}

//新建key,接收post请求的参数
type CreateKeyPostData struct {
	DBID    int64       `json:"db_id"`
	Key     string      `json:"key"`
	KeyType string      `json:"key_type"`
	Value   interface{} `json:"value"`
	TTL     int64       `json:"ttl"`
}

//redis 命令终端post传参
type RedisConsoleData struct {
	RedisID int64  `json:"rid"`
	DBID    int64  `json:"db_id"`
	CMD     string `json:"cmd"`
}

//redis 输出服务配置信息
type RedisConfigData struct {
	ConfigName  string `json:"config_name"`
	ConfigValue string `json:"config_value"`
	ConfigDoc   string `json:"config_doc"`
}

//redis ping 数据结构体
type RedisPingData struct {
	Number int64   `json:"number"`
	Time   float64 `json:"time"` //单位ms
	IsOK   bool    `json:"ping"`
}
