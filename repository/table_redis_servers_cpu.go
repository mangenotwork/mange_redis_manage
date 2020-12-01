package models

/*
CREATE TABLE table_redis_servers_cpu(
   id INTEGER PRIMARY KEY,
   hid TEXT NOT NULL,
   get_time BIGINT NOT NULL,
   used_cpu_sys DOUBLE NOT NULL,
   used_cpu_user DOUBLE NOT NULL,
   used_cpu_sys_children DOUBLE NOT NULL,
   used_cpu_user_children DOUBLE NOT NULL
);
*/

type RedisCPUDB struct {
	ID                  int64   `gorm:"primary_key;column:id" json:"-"`                              //
	Hid                 string  `gorm:"column:hid" json:"host_id"`                                   //
	GetTime             int64   `gorm:"column:get_time" json:"get_time"`                             //
	UsedCpuSys          float64 `gorm:"column:used_cpu_sys" json:"used_cpu_sys"`                     //
	UsedCpuUser         float64 `gorm:"column:used_cpu_user" json:"used_cpu_user"`                   //
	UsedCpuSysChildren  float64 `gorm:"column:used_cpu_sys_children" json:"used_cpu_sys_children"`   //
	UsedCpuUserChildren float64 `gorm:"column:used_cpu_user_children" json:"used_cpu_user_children"` //
}
