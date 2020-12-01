package dao

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/mangenotwork/mange_redis_manage/common/manlog"
	"github.com/mangenotwork/mange_redis_manage/common/sqlitedb"
	"github.com/mangenotwork/mange_redis_manage/repository"
	_ "github.com/mattn/go-sqlite3"
)

type DaoRedisCPU struct {
	Data *models.RedisCPUDB
	Mu   sync.Mutex
}

//提取查询数据
func (this *DaoRedisCPU) exportdatas(rows *sql.Rows) (datas []*models.RedisCPUDB, err error) {
	for rows.Next() {
		data := &models.RedisCPUDB{}
		err := rows.Scan(&data.ID, &data.Hid, &data.GetTime, &data.UsedCpuSys, &data.UsedCpuUser, &data.UsedCpuSysChildren, &data.UsedCpuUserChildren)
		if err != nil {
			manlog.Error(err)
		}
		//manlog.Debug(*data)
		datas = append(datas, data)
	}
	return
}

//提取查询数据
func (this *DaoRedisCPU) exportdatas1(rows *sql.Rows) (data *models.RedisCPUDB, err error) {
	for rows.Next() {
		data = &models.RedisCPUDB{}
		err := rows.Scan(&data.ID, &data.Hid, &data.GetTime, &data.UsedCpuSys, &data.UsedCpuUser, &data.UsedCpuSysChildren, &data.UsedCpuUserChildren)
		if err != nil {
			manlog.Error(err)
		}
	}
	return
}

func (this *DaoRedisCPU) Create() error {
	db := sqlitedb.GetDBConn()
	this.Mu.Lock()
	defer this.Mu.Unlock()
	stmt, err := db.Prepare("INSERT INTO table_redis_servers_cpu (hid,get_time,used_cpu_sys,used_cpu_user,used_cpu_sys_children,used_cpu_user_children)" +
		" values(?,?,?,?,?,?)")
	defer stmt.Close()
	if err != nil {
		manlog.Error(err)
		return err
	}
	res, err := stmt.Exec(this.Data.Hid, this.Data.GetTime, this.Data.UsedCpuSys, this.Data.UsedCpuUser, this.Data.UsedCpuSysChildren, this.Data.UsedCpuUserChildren)
	manlog.Debug(&res)
	if err != nil {
		manlog.Error(err)
		return err
	}
	return nil
}

//查询当前时间之前的n条数据
func (this *DaoRedisCPU) GetLastTimeData(rid string, last_time, n int64) (data []*models.RedisCPUDB, err error) {
	sql := fmt.Sprintf("SELECT * FROM table_redis_servers_cpu where hid='%s' and get_time<%d Order by id desc LIMIT %d;", rid, last_time, n)
	rows, err := sqlitedb.Query(sql)
	defer rows.Close()
	if err != nil {
		manlog.Error(err)
	}
	return this.exportdatas(rows)
}

//获取最新的一条数据
func (this *DaoRedisCPU) GetNewData(rid string) (data *models.RedisCPUDB, err error) {
	sql := fmt.Sprintf("SELECT * FROM table_redis_servers_cpu where hid='%s' Order by get_time desc LIMIT 1;", rid)
	rows, err := sqlitedb.Query(sql)
	defer rows.Close()
	if err != nil {
		manlog.Error(err)
	}
	return this.exportdatas1(rows)
}

//删除多久之前的数据
func (this *DaoRedisCPU) DelTimeData(del_time int64) error {
	this.Mu.Lock()
	defer this.Mu.Unlock()
	sql := "delete from table_redis_servers_cpu where get_time<?"
	return sqlitedb.Del(sql, del_time)
}

//清空数据
func (this *DaoRedisCPU) Empty() error {
	this.Mu.Lock()
	defer this.Mu.Unlock()
	sql := "delete from table_redis_servers_cpu"
	return sqlitedb.Del(sql)
}
