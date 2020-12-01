package models

/*
CREATE TABLE table_redis_slowlog(
   id INTEGER PRIMARY KEY,
   hid TEXT NOT NULL,
   get_time BIGINT NOT NULL,
   only_id BIGINT NOT NULL,
   time BIGINT NOT NULL,
   duration BIGINT NOT NULL,
   cmd TEXT NOT NULL,
   client TEXT NOT NULL,
);

only_id   唯一标识符
time  命令执行时的时间，格式为 UNIX 时间戳
duration   执行命令消耗的时间，以微秒为单位
cmd   执行命令
client   执行命令的客户端
*/

type RedisSlowLog struct {
	ID       int64  `gorm:"primary_key;column:id" json:"-"`  //
	Hid      string `gorm:"column:hid" json:"host_id"`       //
	GetTime  int64  `gorm:"column:get_time" json:"get_time"` //
	OnlyId   int64  `gorm:"column:only_id" json:"only_id"`   //唯一标识符
	Time     int64  `gorm:"column:time" json:"time"`         //命令执行时的时间，格式为 UNIX 时间戳
	Duration int64  `gorm:"column:duration" json:"duration"` //执行命令消耗的时间，以微秒为单位
	Cmd      string `gorm:"column:cmd" json:"cmd"`           //执行命令
	Client   string `gorm:"column:client" json:"client"`     //执行命令的客户端
}
