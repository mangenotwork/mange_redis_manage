package models

/*
CREATE TABLE table_redis_servers_replication(
    id INTEGER PRIMARY KEY,
    hid TEXT NOT NULL,
    get_time BIGINT NOT NULL,
    role_value TEXT NOT NULL,
    connected_slaves TEXT NOT NULL,
    master_replid TEXT NOT NULL,
    master_replid2 TEXT NOT NULL,
    master_repl_offset TEXT NOT NULL,
    second_repl_offset TEXT NOT NULL,
    repl_backlog_active TEXT NOT NULL,
    repl_backlog_size BIGINT NOT NULL,
    repl_backlog_first_byteoffset BIGINT NOT NULL,
    repl_backlog_histlen BIGINT NOT NULL
);
*/

type RedisReplicationDB struct {
	ID                         int64  `gorm:"primary_key;column:id" json:"-"`                                            //
	Hid                        string `gorm:"column:hid" json:"host_id"`                                                 //
	GetTime                    int64  `gorm:"column:get_time" json:"get_time"`                                           //
	Role                       string `gorm:"column:role_value" json:"role_value"`                                       //
	ConnectedSlaves            string `gorm:"column:connected_slaves" json:"connected_slaves"`                           //
	MasterReplid               string `gorm:"column:master_replid" json:"master_replid"`                                 //
	MasterReplid2              string `gorm:"column:master_replid2" json:"master_replid2"`                               //
	MasterReplOffset           string `gorm:"column:master_repl_offset" json:"master_repl_offset"`                       //主从同步偏移量（通过主从对比判断主从同步是否一致）
	SecondReplOffset           string `gorm:"column:second_repl_offset" json:"second_repl_offset"`                       //
	ReplBacklogActive          string `gorm:"column:repl_backlog_active" json:"repl_backlog_active"`                     //
	ReplBacklogSize            int64  `gorm:"column:repl_backlog_size" json:"repl_backlog_size"`                         //
	ReplBacklogFirstByteOffset int64  `gorm:"column:repl_backlog_first_byteoffset" json:"repl_backlog_first_byteoffset"` //
	ReplBacklogHistlen         int64  `gorm:"column:repl_backlog_histlen" json:"repl_backlog_histlen"`                   //
}
