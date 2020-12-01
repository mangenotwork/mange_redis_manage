//
//	controller基础方法
//

package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mangenotwork/mange_redis_manage/common"
	"github.com/mangenotwork/mange_redis_manage/common/cache"
	"github.com/mangenotwork/mange_redis_manage/common/jwt"
	"github.com/mangenotwork/mange_redis_manage/common/manlog"
	"github.com/mangenotwork/mange_redis_manage/structs"
)

type BaseController struct {
	Context *gin.Context
}

//API 输出
func (this *BaseController) APIOutPut(code int64, mag string, data interface{}) {
	this.Context.IndentedJSON(http.StatusOK, structs.ResponseJson{
		Code:      code,
		Mag:       mag,
		Date:      data,
		TimeStamp: time.Now().Unix(),
	})
	return
}

//API 输出 错误
func (this *BaseController) APIOutPutError(code int64, mag string) {
	this.Context.IndentedJSON(http.StatusOK, structs.ResponseJson{
		Code:      code,
		Mag:       mag,
		Date:      "",
		TimeStamp: time.Now().Unix(),
	})
	return
}

//Set Sing
func (this *BaseController) SetSign() {
	this.Context.SetCookie("sign", common.NewMangeSign(), 60*60*24*7, "/", "localhost", false, true)
}

//接收post 传入参数
func (this *BaseController) GetPostArgs(obj interface{}) error {
	err := this.Context.BindJSON(obj)
	if err != nil {
		manlog.Error(err)
		this.APIOutPutError(1, "非法传入参数!")
		return err
	}
	return nil
}

//用户鉴权中间件 - 通用
func (this *BaseController) Authentication() gin.HandlerFunc {
	manlog.Debug("用户鉴权中间件 - 通用")
	return func(c *gin.Context) {
		this.Context = c
		userPara := this.GetUser()
		if userPara != nil {
			//上下文添加解析token的用户信息
			this.Context.Set("userParameter", userPara)
			this.Context.Next()
			return
		} else {
			this.Context.Abort()
			this.APIOutPutError(2, "未登录")
			return
		}
	}
}

//用户鉴权中间件 - 页面
func AuthPG() gin.HandlerFunc {
	return func(c *gin.Context) {
		this := new(BaseController)
		this.Context = c
		userPara := this.GetUser()
		manlog.Error("userPara = ", userPara)
		if userPara != nil {
			//上下文添加解析token的用户信息
			this.Context.Set("userParameter", userPara)
			this.Context.Next()
			return
		} else {
			manlog.Error("notlogin")
			this.Context.Abort()
			this.Context.Redirect(http.StatusFound, "/notlogin")
			return
		}
	}
}

//获取上下文的用户信息
func (this *BaseController) GetUserParameter() (userdata *structs.UserParameter) {
	userdata = &structs.UserParameter{}
	data := this.Context.MustGet("userParameter")
	if data != "" {
		userdata = data.(*structs.UserParameter)
	}
	return userdata
}

//直接获取用户信息,非中间件
func (this *BaseController) GetUser() *structs.UserParameter {
	token, _ := this.Context.Cookie("token")
	manlog.Debug("token = ", token)
	if token != "" {
		//解析token
		userparse, err := jwt.UserParseToken(token)
		if err == nil {
			//在系统缓存匹配token
			cachetoken, _ := cache.Get(string(userparse.UserID))
			if cachetoken == token {
				return userparse.UserParameter
			}
		}
	}
	return nil
}
