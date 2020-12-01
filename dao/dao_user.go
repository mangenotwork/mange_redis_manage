package dao

import (
	"fmt"

	"github.com/mangenotwork/mange_redis_manage/common/manlog"
	"github.com/mangenotwork/mange_redis_manage/common/sqlitedb"
	"github.com/mangenotwork/mange_redis_manage/repository"
)

//
type DaoUser struct {
	models.User
}

//用户登录验证
func (this *DaoUser) Login() (isuser bool, user *models.User) {
	user = &models.User{}
	isuser = false

	rows, err := sqlitedb.Query(fmt.Sprintf("SELECT * FROM table_user where uname='%s' and upassword='%s';", this.Uname, this.Upassword))
	if err != nil {
		manlog.Error(err)
		return
	}
	manlog.Debug(rows)
	for rows.Next() {
		err := rows.Scan(&user.Uid, &user.Uname, &user.Upassword, &user.Ugroup)
		if err != nil {
			manlog.Error(err)
		}
		if user.Uid > 0 {
			isuser = true
		}
	}
	return
}
