package models

/*
CREATE TABLE table_redis_servers_stats(
    id INTEGER PRIMARY KEY,
    hid TEXT NOT NULL,
    get_time BIGINT NOT NULL,
	total_connections_received BIGINT NOT NULL,
	total_commands_processed BIGINT NOT NULL,
	instantaneous_ops_per_sec BIGINT NOT NULL,
	total_net_input_bytes BIGINT NOT NULL,
	total_net_output_bytes BIGINT NOT NULL,
	instantaneous_input_kbps DOUBLE NOT NULL,
	instantaneous_output_kbps DOUBLE NOT NULL,
	rejected_connections BIGINT NOT NULL,
	sync_full BIGINT NOT NULL,
	sync_partial_ok BIGINT NOT NULL,
	sync_partial_err BIGINT NOT NULL,
	expired_keys BIGINT NOT NULL,
	expired_stale_perc TEXT NOT NULL,
	expired_time_cap_reached_count BIGINT NOT NULL,
	evicted_keys BIGINT NOT NULL,
	keyspace_hits BIGINT NOT NULL,
	keyspace_misses BIGINT NOT NULL,
	pubsub_channels BIGINT NOT NULL,
	pubsub_patterns BIGINT NOT NULL,
	latest_fork_usec BIGINT NOT NULL,
	migrate_cached_sockets BIGINT NOT NULL,
	slave_expires_tracked_keys BIGINT NOT NULL,
	active_defrag_hits BIGINT NOT NULL,
	active_defrag_misses BIGINT NOT NULL,
	active_defrag_key_hits BIGINT NOT NULL,
	active_defrag_key_misses BIGINT NOT NULL
);
*/

type RedisStatsDB struct {
	ID                         int64   `gorm:"primary_key;column:id" json:"-"`
	Hid                        string  `gorm:"column:hid" json:"host_id"`
	GetTime                    int64   `gorm:"column:get_time" json:"get_time"`
	TotalConnectionsReceived   int64   `gorm:"column:total_connections_received" json:"total_connections_received"`         //新创建的链接个数，如果过多，会影响性能
	TotalCommandsProcessed     int64   `gorm:"column:total_commands_processed" json:"total_commands_processed"`             //redis处理的命令数
	InstantaneousOpsPerSec     int64   `gorm:"column:instantaneous_ops_per_sec" json:"instantaneous_ops_per_sec"`           //redis当前的qps，redis内部较实时的每秒执行命令数
	TotalNetInputBytes         int64   `gorm:"column:total_net_input_bytes" json:"total_net_input_bytes"`                   //redis网络入口流量字节数
	TotalNetOutputBytes        int64   `gorm:"column:total_net_output_bytes" json:"total_net_output_bytes"`                 //redis网络出口流量字节数
	InstantaneousInputKbps     float64 `gorm:"column:instantaneous_input_kbps" json:"instantaneous_input_kbps"`             //redis网络入口kps
	InstantaneousOutputKbps    float64 `gorm:"column:instantaneous_output_kbps" json:"instantaneous_output_kbps"`           //redis网络出口kps
	RejectedConnections        int64   `gorm:"column:rejected_connections" json:"rejected_connections"`                     //拒绝的连接个数，redis连接个数已经达到maxclients限制。
	SyncFull                   int64   `gorm:"column:sync_full" json:"sync_full"`                                           //主从完全同步成功次数
	SyncPartialOk              int64   `gorm:"column:sync_partial_ok" json:"sync_partial_ok"`                               //主从部分同步成功次数
	SyncPartialErr             int64   `gorm:"column:sync_partial_err" json:"sync_partial_err"`                             //主从部分同步失败次数
	ExpiredKeys                int64   `gorm:"column:expired_keys" json:"expired_keys"`                                     //运行以来过期的key的数量
	ExpiredStalePerc           string  `gorm:"column:expired_stale_perc" json:"expired_stale_perc"`                         //
	ExpiredTimeCapReachedCount int64   `gorm:"column:expired_time_cap_reached_count" json:"expired_time_cap_reached_count"` //
	EvictedKeys                int64   `gorm:"column:evicted_keys" json:"evicted_keys"`                                     //运行以来剔除（超过maxmemory）的key的数量s
	KeyspaceHits               int64   `gorm:"column:keyspace_hits" json:"keyspace_hits"`                                   //命中次数
	KeyspaceMisses             int64   `gorm:"column:keyspace_misses" json:"keyspace_misses"`                               //没命中次数
	PubsubChannels             int64   `gorm:"column:pubsub_channels" json:"pubsub_channels"`                               //当前使用中的频道数量
	PubsubPatterns             int64   `gorm:"column:pubsub_patterns" json:"pubsub_patterns"`                               //当前使用的模式数量
	LatestForkUsec             int64   `gorm:"column:latest_fork_usec" json:"latest_fork_usec"`                             //
	MigrateCachedSockets       int64   `gorm:"column:migrate_cached_sockets" json:"migrate_cached_sockets"`                 //
	SlaveExpiresTrackedKeys    int64   `gorm:"column:slave_expires_tracked_keys" json:"slave_expires_tracked_keys"`         //
	ActiveDefragHits           int64   `gorm:"column:active_defrag_hits" json:"active_defrag_hits"`                         //
	ActiveDefragMisses         int64   `gorm:"column:active_defrag_misses" json:"active_defrag_misses"`                     //
	ActiveDefragKeyHits        int64   `gorm:"column:active_defrag_key_hits" json:"active_defrag_key_hits"`                 //
	ActiveDefragKeyMisses      int64   `gorm:"column:active_defrag_key_misses" json:"active_defrag_key_misses"`             //
}
