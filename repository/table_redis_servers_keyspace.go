package models

/*
CREATE TABLE table_redis_servers_keyspace(
    id INTEGER PRIMARY KEY,
    hid TEXT NOT NULL,
    get_time BIGINT NOT NULL,
   	db_id BIGINT NOT NULL,
	keys_count BIGINT NOT NULL,
	expires BIGINT NOT NULL,
	avgttl BIGINT NOT NULL
);
*/

type RedisKeyspaceDB struct {
	ID      int64  `gorm:"primary_key;column:id" json:"-"`      //
	Hid     string `gorm:"column:hid" json:"host_id"`           //
	GetTime int64  `gorm:"column:get_time" json:"get_time"`     //
	DBID    int64  `gorm:"column:db_id" json:"db_id"`           //
	Keys    int64  `gorm:"column:keys_count" json:"keys_count"` //
	Expires int64  `gorm:"column:expires" json:"expires"`       //
	AvgTTL  int64  `gorm:"column:avgttl" json:"avgttl"`         //
}
