//
//	table_redis_info表 主要保存连接信息
//
package models

/*
CREATE TABLE table_redis_info(
   id INTEGER PRIMARY KEY,
   uid INT NOT NULL,
   conn_name TEXT NOT NULL,
   conn_host TEXT NOT NULL,
   conn_port INT NOT NULL,
   conn_password TEXT NOT NULL,
   is_ssh BOOLEAN NOT NULL,
   ssh_url TEXT NOT NULL,
   ssh_user TEXT NOT NULL,
   ssh_password TEXT NOT NULL,
   conn_create BIGINT NOT NULL
);
*/

type RedisInfoDB struct {
	ID           int64  `gorm:"primary_key;column:id"`
	UID          int64  `gorm:"column:uid"`
	ConnName     string `gorm:"column:conn_name"`     //连接名 主键不能重复
	ConnHost     string `gorm:"column:conn_host"`     //连接的主机地址
	ConnPort     int    `gorm:"column:conn_port"`     //连接的端口
	ConnPassword string `gorm:"column:conn_password"` //连接的密码
	IsSSH        bool   `gorm:"column:is_ssh"`        //是否使用ssh连接 0不使用 1使用
	SSHUrl       string `gorm:"column:ssh_url"`       //ssh地址
	SSHUser      string `gorm:"column:ssh_user"`      //ssh 账号
	SSHPassword  string `gorm:"column:ssh_password"`  //ssh密码
	ConnCreate   int64  `gorm:"column:conn_create"`   //创建连接的时间戳
}
