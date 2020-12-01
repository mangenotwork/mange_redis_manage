package models

/*
CREATE TABLE table_redis_servers_clients(
   id INTEGER PRIMARY KEY,
   hid TEXT NOT NULL,
   get_time BIGINT NOT NULL,
   connected_clients BIGINT NOT NULL,
   client_recent_max_input_buffer BIGINT NOT NULL,
   client_recent_max_output_buffer BIGINT NOT NULL,
   blocked_clients BIGINT NOT NULL
);
*/

type RedisClientsDB struct {
	ID                          int64  `gorm:"primary_key;column:id" json:"-"`
	Hid                         string `gorm:"column:hid" json:"host_id"`
	GetTime                     int64  `gorm:"column:get_time" json:"get_time"`
	ConnectedClients            int64  `gorm:"column:connected_clients" json:"connected_clients"`                             //已经连接客户端数量（不包括slave连接的客户端）
	ClientRecentMaxInputBuffer  int64  `gorm:"column:client_recent_max_input_buffer" json:"client_recent_max_input_buffer"`   //客户端最近最大输入缓冲区
	ClientRecentMaxOutputBuffer int64  `gorm:"column:client_recent_max_output_buffer" json:"client_recent_max_output_buffer"` //客户端最近最大输出缓冲区
	BlockedClients              int64  `gorm:"column:clocked_clients" json:"clocked_clients"`                                 //正在等待阻塞命令的客户端数量
}
