//
//
//
package dao

import (
	"database/sql"
	"fmt"
	_ "reflect"

	"github.com/mangenotwork/mange_redis_manage/common/manlog"
	"github.com/mangenotwork/mange_redis_manage/common/sqlitedb"
	"github.com/mangenotwork/mange_redis_manage/repository"
	_ "github.com/mattn/go-sqlite3"
)

type DaoRedisInfo struct {
	models.RedisInfoDB
}

//提取连接数据
func (this *DaoRedisInfo) exportdatas(rows *sql.Rows) (datas []*models.RedisInfoDB, err error) {
	for rows.Next() {
		data := &models.RedisInfoDB{}
		err := rows.Scan(&data.ID, &data.UID, &data.ConnName, &data.ConnHost, &data.ConnPort, &data.ConnPassword,
			&data.IsSSH, &data.SSHUrl, &data.SSHUser, &data.SSHPassword, &data.ConnCreate)
		if err != nil {
			manlog.Error(err)
		}
		//manlog.Debug(*data)
		datas = append(datas, data)
	}
	return
}

//查询指定用户的所有连接
func (this *DaoRedisInfo) GetAll(uid int64) (datas []*models.RedisInfoDB, err error) {
	rows, err := sqlitedb.Query(fmt.Sprintf("SELECT * FROM table_redis_info where uid=%d", uid))
	defer rows.Close()
	if err != nil {
		manlog.Error(err)
	}
	datas, err = this.exportdatas(rows)
	return
}

//查询所有连接,并去重过
func (this *DaoRedisInfo) GetAllConn() (datas []*models.RedisInfoDB, err error) {
	rows, err := sqlitedb.Query(fmt.Sprintf("SELECT * FROM table_redis_info group by conn_host,conn_port"))
	defer rows.Close()
	if err != nil {
		manlog.Error(err)
	}
	datas, err = this.exportdatas(rows)
	return
}

//新建连接
func (this *DaoRedisInfo) Create() error {
	db := sqlitedb.GetDBConn()
	stmt, err := db.Prepare("INSERT INTO table_redis_info (uid,conn_name,conn_host,conn_port,conn_password,is_ssh,ssh_url,ssh_user," +
		"ssh_password,conn_create) values(?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		manlog.Error(err)
		return err
	}
	res, err := stmt.Exec(this.UID, this.ConnName, this.ConnHost, this.ConnPort, this.ConnPassword,
		this.IsSSH, this.SSHUrl, this.SSHUser, this.SSHPassword, this.ConnCreate)
	//fmt.Println(&res)
	manlog.Debug(&res)
	if err != nil {
		manlog.Error(err)
		return err
	}
	return nil
}

//通过uid 与redis conn id 获取连接信息
func (this *DaoRedisInfo) GetConnInfo(uid, redisid int64) (data *models.RedisInfoDB, err error) {
	data = &models.RedisInfoDB{}
	rows, err := sqlitedb.Query(fmt.Sprintf("SELECT * FROM table_redis_info where uid=%d and id=%d", uid, redisid))
	defer rows.Close()
	if err != nil {
		manlog.Error(err)
	}
	for rows.Next() {
		err := rows.Scan(&data.ID, &data.UID, &data.ConnName, &data.ConnHost, &data.ConnPort, &data.ConnPassword,
			&data.IsSSH, &data.SSHUrl, &data.SSHUser, &data.SSHPassword, &data.ConnCreate)
		if err != nil {
			manlog.Error(err)
		}
	}
	return
}
