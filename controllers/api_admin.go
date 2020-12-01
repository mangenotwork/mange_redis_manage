//
//	管理员接口,设置接口
//

package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mangenotwork/mange_redis_manage/service"
)

type AdminAPI struct {
	BaseController
}

func (this *AdminAPI) Say(c *gin.Context) {
	fmt.Printf("this is admin api")
	c.String(http.StatusOK, "this is admin api")
	return
}

//查看所有缓存
func (this *AdminAPI) CacheGetAll(c *gin.Context) {
	this.Context = c

	var CaService service.CacheService = new(service.HostCache)
	CaService.GetAll()

	c.String(http.StatusOK, "this is admin api")
	return
}

//查看指定缓存
func (this *AdminAPI) CacheGet(c *gin.Context) {
	fmt.Printf("this is admin api")
	c.String(http.StatusOK, "this is admin api")
	return
}

//修改指定缓存
func (this *AdminAPI) CacheUpdate(c *gin.Context) {
	fmt.Printf("this is admin api")
	c.String(http.StatusOK, "this is admin api")
	return
}

//删除指定缓存
func (this *AdminAPI) CacheDel(c *gin.Context) {
	fmt.Printf("this is admin api")
	c.String(http.StatusOK, "this is admin api")
	return
}

//新增缓存
func (this *AdminAPI) CacheAdd(c *gin.Context) {
	fmt.Printf("this is admin api")
	c.String(http.StatusOK, "this is admin api")
	return
}
