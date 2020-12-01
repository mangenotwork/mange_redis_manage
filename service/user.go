//
//	user 用户服务
//
package service

import (
	"github.com/mangenotwork/mange_redis_manage/common/cache"
	"github.com/mangenotwork/mange_redis_manage/common/jwt"
	"github.com/mangenotwork/mange_redis_manage/common/manlog"
	"github.com/mangenotwork/mange_redis_manage/dao"
	"github.com/mangenotwork/mange_redis_manage/structs"
)

type UserService interface {
	Login(logindata *structs.UserLogin) (userJwt string, err error) //用户登录
	Info()                                                          //用户信
}

type User struct {
}

func (this *User) Login(logindata *structs.UserLogin) (userJwt string, err error) {
	manlog.Error("test err log")
	manlog.Info("test info log")
	user := new(dao.DaoUser)
	user.Uname = logindata.User
	user.Upassword = logindata.Password
	islogin, userinfo := user.Login()

	userJwt = ""

	//生成jwt
	if islogin {
		usertoken := &structs.UserParameter{userinfo.Uid, userinfo.Uname, userinfo.Ugroup}
		userJwt, err = jwt.UserToken(usertoken)
		cache.Set(string(userinfo.Uid), userJwt)
		manlog.Debug(userJwt, err)
	}
	return
}

func (this *User) GetUserToken(username string) string {
	return ""
}

func (this *User) Info() {}
