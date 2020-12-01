package models

/*
CREATE TABLE table_redis_servers_infos(
    id INTEGER PRIMARY KEY,
    hid TEXT NOT NULL,
    get_time BIGINT NOT NULL,
    redis_version TEXT NOT NULL,
	redis_git_sha1 TEXT NOT NULL,
	redis_git_dirty TEXT NOT NULL,
	redis_build_id TEXT NOT NULL,
	redis_mode TEXT NOT NULL,
	os TEXT NOT NULL,
	arch_bits TEXT NOT NULL,
	multiplexing_api TEXT NOT NULL,
	atomicvar_api TEXT NOT NULL,
	gcc_version TEXT NOT NULL,
	process_id TEXT NOT NULL,
	run_id TEXT NOT NULL,
	tcp_port BIGINT NOT NULL,
	uptime_in_seconds BIGINT NOT NULL,
	uptime_in_days BIGINT NOT NULL,
	hz BIGINT NOT NULL,
	configured_hz BIGINT NOT NULL,
	lru_clock BIGINT NOT NULL,
	executable TEXT NOT NULL,
	config_file TEXT NOT NULL
);
*/

type RedisServerInfosDB struct {
	ID              int64  `gorm:"primary_key;column:id" json:"-"`
	Hid             string `gorm:"column:hid" json:"host_id"`
	GetTime         int64  `gorm:"column:get_time" json:"get_time"`
	RedisVersion    string `gorm:"column:redis_version" json:"redis_version"`         //Redis 服务器版本
	RedisGitSha1    string `gorm:"column:redis_git_sha1" json:"redis_git_sha1"`       //Git SHA1
	RedisGitDirty   string `gorm:"column:redis_git_dirty" json:"redis_git_dirty"`     //Git dirty flag
	RedisBuildId    string `gorm:"column:redis_build_id" json:"redis_build_id"`       //Redis build id
	RedisMode       string `gorm:"column:redis_mode" json:"redis_mode"`               //运行模式，单机或集群
	Os              string `gorm:"column:os" json:"os"`                               //Redis 服务器的宿主操作系统
	ArchBits        string `gorm:"column:arch_bits" json:"arch_bits"`                 //架构64位
	MultiplexingApi string `gorm:"column:multiplexing_api" json:"multiplexing_api"`   //redis所使用的事件处理模型
	AtomicvarApi    string `gorm:"column:atomicvar_api" json:"atomicvar_api"`         //
	GccVersion      string `gorm:"column:gcc_version" json:"gcc_version"`             //编译redis时gcc版本
	ProcessId       string `gorm:"column:process_id" json:"process_id"`               //redis服务器进程的pid
	RunId           string `gorm:"column:run_id" json:"run_id"`                       //redis服务器的随机标识符（sentinel和集群）
	TcpPort         int64  `gorm:"column:tcp_port" json:"tcp_port"`                   //
	UptimeInSeconds int64  `gorm:"column:uptime_in_seconds" json:"uptime_in_seconds"` //redis服务器启动总时间，单位秒
	UptimeInDays    int64  `gorm:"column:uptime_in_days" json:"uptime_in_days"`       //redis服务器启动总时间，单位天
	Hz              int64  `gorm:"column:hz" json:"hz"`                               //redis内部调度频率（关闭timeout客户端，删除过期key）
	ConfiguredHz    int64  `gorm:"column:configured_hz" json:"configured_hz"`         //
	Lru_clock       int64  `gorm:"column:lru_clock" json:"lru_clock"`                 //自增时间，用于LRU管理
	Executable      string `gorm:"column:executable" json:"executable"`               //
	ConfigFile      string `gorm:"column:config_file" json:"config_file"`             //配置文件路径
}
