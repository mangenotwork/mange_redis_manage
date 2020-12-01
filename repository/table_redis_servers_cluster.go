package models

/*
CREATE TABLE table_redis_servers_cluster(
   id INTEGER PRIMARY KEY,
   hid TEXT NOT NULL,
   get_time BIGINT NOT NULL,
   cluster_enabled TEXT NOT NULL
);
*/

type RedisClusterDB struct {
	ID             int64  `gorm:"primary_key;column:id" json:"-"`                //
	Hid            string `gorm:"column:hid" json:"host_id"`                     //
	GetTime        int64  `gorm:"column:get_time" json:"get_time"`               //
	ClusterEnabled string `gorm:"column:cluster_enabled" json:"cluster_enabled"` //
}
