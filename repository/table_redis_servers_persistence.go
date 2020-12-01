package models

/*
CREATE TABLE table_redis_servers_persistence(
    id INTEGER PRIMARY KEY,
    hid TEXT NOT NULL,
    get_time BIGINT NOT NULL,
   	loading BIGINT NOT NULL,
	rdb_changes_since_last_save BIGINT NOT NULL,
	rdb_bgsave_in_progress BIGINT NOT NULL,
	rdb_last_save_time BIGINT NOT NULL,
	rdb_last_bgsave_status TEXT NOT NULL,
	rdb_last_bgsave_time_sec BIGINT NOT NULL,
	rdb_current_bgsave_time_sec BIGINT NOT NULL,
	rdb_last_cow_size BIGINT NOT NULL,
	aof_enabled BIGINT NOT NULL,
	aof_rewrite_in_progress BIGINT NOT NULL,
	aof_rewrite_scheduled BIGINT NOT NULL,
	aof_last_rewrite_time_sec BIGINT NOT NULL,
	aof_current_rewrite_time_sec BIGINT NOT NULL,
	aof_last_bgrewrite_status TEXT NOT NULL,
	aof_last_write_status TEXT NOT NULL,
	aof_last_cow_size BIGINT NOT NULL,
	aof_current_size BIGINT NOT NULL,
	aof_base_size BIGINT NOT NULL,
	aof_pending_rewrite BIGINT NOT NULL,
	aof_buffer_length BIGINT NOT NULL,
	aof_rewrite_buffer_length BIGINT NOT NULL,
	aof_pending_bio_fsync BIGINT NOT NULL,
	aof_delayed_fsync BIGINT NOT NULL
);
*/

type RedisPersistenceDB struct {
	ID                       int64  `gorm:"primary_key;column:id" json:"-"`
	Hid                      string `gorm:"column:hid" json:"host_id"`
	GetTime                  int64  `gorm:"column:get_time" json:"get_time"`
	Loading                  int64  `gorm:"column:loading" json:"loading"`                                           //服务器是否正在载入持久化文件
	RdbChangesSinceLastSave  int64  `gorm:"column:rdb_changes_since_last_save" json:"rdb_changes_since_last_save"`   //有多少个已经写入的命令还未被持久化
	RdbBgsaveInProgress      int64  `gorm:"column:rdb_bgsave_in_progress" json:"rdb_bgsave_in_progress"`             //服务器是否正在创建rdb文件
	RdbLastSaveTime          int64  `gorm:"column:rdb_last_save_time" json:"rdb_last_save_time"`                     //已经有多长时间没有进行持久化了
	RdbLastBgsaveStatus      string `gorm:"column:rdb_last_bgsave_status" json:"rdb_last_bgsave_status"`             //最后一次的rdb持久化是否成功
	RdbLastBgsaveTimeSec     int64  `gorm:"column:rdb_last_bgsave_time_sec" json:"rdb_last_bgsave_time_sec"`         //最后一次生成rdb文件耗时秒数
	RdbCurrentBgsaveTimeSec  int64  `gorm:"column:rdb_current_bgsave_time_sec" json:"rdb_current_bgsave_time_sec"`   //如果服务器正在创建rdb文件，那么当前这个记录就是创建操作耗时秒数
	RdbLastCowSize           int64  `gorm:"column:rdb_last_cow_size" json:"rdb_last_cow_size"`                       //
	AofEnabled               int64  `gorm:"column:aof_enabled" json:"aof_enabled"`                                   //是否开启了aof
	AofRewriteInProgress     int64  `gorm:"column:aof_rewrite_in_progress" json:"aof_rewrite_in_progress"`           //标识aof的rewrite操作是否进行中
	AofRewriteScheduled      int64  `gorm:"column:aof_rewrite_scheduled" json:"aof_rewrite_scheduled"`               //
	AofLastRewriteTimeSec    int64  `gorm:"column:aof_last_rewrite_time_sec" json:"aof_last_rewrite_time_sec"`       //
	AofCurrentRewriteTimeSec int64  `gorm:"column:aof_current_rewrite_time_sec" json:"aof_current_rewrite_time_sec"` //
	AofLastBgrewriteStatus   string `gorm:"column:aof_last_bgrewrite_status" json:"aof_last_bgrewrite_status"`       //上次bgrewriteaof操作的状态
	AofLastWriteStatus       string `gorm:"column:aof_last_write_status" json:"aof_last_write_status"`               //上一次aof写入状态
	AofLastCowSize           int64  `gorm:"column:aof_last_cow_size" json:"aof_last_cow_size"`                       //
	AofCurrentSize           int64  `gorm:"column:aof_current_size" json:"aof_current_size"`                         //
	AofBaseSize              int64  `gorm:"column:aof_base_size" json:"aof_base_size"`                               //
	AofPendingRewrite        int64  `gorm:"column:aof_pending_rewrite" json:"aof_pending_rewrite"`                   //
	AofBufferLength          int64  `gorm:"column:aof_buffer_length" json:"aof_buffer_length"`                       //
	AofRewriteBufferLength   int64  `gorm:"column:aof_rewrite_buffer_length" json:"aof_rewrite_buffer_length"`       //
	AofPendingBioFsync       int64  `gorm:"column:aof_pending_bio_fsync" json:"aof_pending_bio_fsync"`               //
	AofDelayedFsync          int64  `gorm:"column:aof_delayed_fsync" json:"aof_delayed_fsync"`                       //
}
