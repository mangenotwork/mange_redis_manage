//
//	user接口
//

package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mangenotwork/mange_redis_manage/service"
	"github.com/mangenotwork/mange_redis_manage/structs"
)

type UserAPI struct {
	BaseController
}

func (this *UserAPI) Say(c *gin.Context) {
	fmt.Printf("this is user api")
	c.String(http.StatusOK, "this is user api")
	return
}

func (this *UserAPI) Login(c *gin.Context) {
	this.Context = c

	user_data := new(structs.UserLogin)
	if err := this.GetPostArgs(&user_data); err != nil {
		return
	}
	var userService service.UserService = new(service.User)
	token, err := userService.Login(user_data)
	if err != nil || token == "" {
		this.APIOutPut(1, "登录失败", "")
		return
	}

	this.Context.SetCookie("token", token, 60*60*24*7, "/", "0.0.0.0/24", false, true)
	this.APIOutPut(0, "", "登录成功")
	return
}
