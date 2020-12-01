//
//	页面
//

package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mangenotwork/mange_redis_manage/common"
	_ "github.com/mangenotwork/mange_redis_manage/common/cache"
	"github.com/mangenotwork/mange_redis_manage/common/manlog"
	"github.com/mangenotwork/mange_redis_manage/service"
	"github.com/mangenotwork/mange_redis_manage/structs"
)

type ClientPG struct {
	BaseController
}

func (this *ClientPG) Say(c *gin.Context) {
	fmt.Printf("this is ClientPG")

	token, _ := c.Cookie("token")
	manlog.Debug(token)

	c.String(http.StatusOK, token)

	return
}

func (this *ClientPG) Index(c *gin.Context) {
	this.Context = c
	this.Context.SetCookie("sign", common.NewMangeSign(), 60*60*24*7, "/", "0.0.0.0/24", false, true)

	isShow := 1
	// welcome, err := this.Context.Cookie("welcome")
	// if welcome == "" || err != nil {
	// 	isShow = 1
	// }

	sign, err := this.Context.Cookie("sign")
	manlog.Debug(sign, err)

	user := this.GetUser()
	manlog.Debug(user)

	isLogin := 0 //未登录
	username := ""
	if user != nil {
		isLogin = 1 //已登录
		isShow = 0
		username = user.Username
	}

	this.Context.HTML(http.StatusOK, "home.html", gin.H{
		"is_show":   isShow,
		"title":     "Ymzy Redis 工具v0.1",
		"welcome":   "Ymzy Redis 工具v0.1",
		"thank":     "感谢圆梦时刻提供技术支持!",
		"thank_url": "https://www.ymzy.cn",
		"author":    "ManGe (2912882908@qq.com)",
		"isLogin":   isLogin,
		"username":  username,
	})

	return
}

func (this *ClientPG) NotLogin(c *gin.Context) {
	this.Context = c
	this.Context.HTML(http.StatusOK, fmt.Sprintf("notlogin.html"), gin.H{})
}

func (this *ClientPG) Echarts(c *gin.Context) {
	this.Context = c
	id := this.Context.Param("id")
	this.Context.HTML(http.StatusOK, fmt.Sprintf("echarts%s.html", id), gin.H{})
}

func (this *ClientPG) EchartsRedisMemory(c *gin.Context) {
	this.Context = c
	user := this.GetUserParameter()

	redisId := this.Context.Param("redisId")
	rid := common.Str2Int64(redisId)

	h := this.Context.Query("h")
	d := this.Context.Query("d")
	if h == "" {
		h = "0"
	}
	if d == "" {
		d = "0"
	}

	hours := common.Str2Int64(h)
	day := common.Str2Int64(d)

	var rcServers service.RedisConnService = new(service.RedisConn)
	datas, _ := rcServers.GetEchartsRedisMemoryData(user, rid, hours, day)

	this.Context.HTML(http.StatusOK, "redis_index.html", gin.H{
		"rid":              rid,
		"time_list":        datas.TimeList,
		"used_memory":      datas.UsedMemory,
		"used_memory_rss":  datas.UsedMemoryRss,
		"used_memory_lua":  datas.UsedMemoryLua,
		"used_memory_peak": datas.UsedMemoryPeak,
		"used_memory_str":  datas.UsedMemoryStr,
		"clinet_number":    datas.ClinetNumber,
		"cmder_number":     datas.CmderNumber,
		"run_time":         datas.RunTime,
		"dbs":              datas.RedisDB,
		"h":                h,
		"d":                d,
	})
}

func (this *ClientPG) RedisClientList(c *gin.Context) {
	this.Context = c
	user := this.GetUser()
	manlog.Debug(*user)

	redisId := this.Context.Param("redisId")
	rid := common.Str2Int64(redisId)
	this.Context.HTML(http.StatusOK, "redis_clinet.html", gin.H{
		"rid": rid,
	})

}

func (this *ClientPG) RedisOperation(c *gin.Context) {
	this.Context = c
	user := this.GetUser()
	manlog.Debug(*user)

	redisId := this.Context.Param("redisId")
	rid := common.Str2Int64(redisId)

	var rcServers service.RedisConnService = new(service.RedisConn)
	dblist := rcServers.GetRedisDBList(user, rid)
	initdb := &structs.RedisServersDBInfo{}
	if len(dblist) > 0 {
		initdb = dblist[0]
	}

	this.Context.HTML(http.StatusOK, "db_doing.html", gin.H{
		"rid":    rid,
		"dblist": dblist,
		"initdb": initdb,
	})
}

func (this *ClientPG) RedisRealtime(c *gin.Context) {
	this.Context = c
	user := this.GetUser()
	manlog.Debug(*user)

	redisId := this.Context.Param("redisId")
	rid := common.Str2Int64(redisId)

	var rinfoServers service.RedisInfoService = new(service.RedisInfo)
	data, err := rinfoServers.RedisRealTimeInit(user, rid)
	if err != nil {
		manlog.Error(err)
	}

	this.Context.HTML(http.StatusOK, "real_time.html", gin.H{
		"real_time_id":                   data.RealTimeId,
		"rid":                            rid,
		"option_x":                       data.Xdata,
		"cpu_option_series_data":         data.CPUData,
		"cpu_option_name":                "cpu",
		"memory_option_series_data":      data.MemoryData,
		"memory_option_name":             "memory",
		"memory_dw":                      data.MemoryDW,
		"qps_option_series_data":         data.QpsData,
		"qps_option_name":                "qps",
		"qps_dw":                         "(1/s)",
		"conn_option_series_data":        data.ConnData,
		"conn_option_name":               "conn",
		"keys_option_series_data":        data.KeysData,
		"keys_option_name":               "keys",
		"kbps_input_option_series_data":  data.InputKbpsData,
		"kbps_input_option_name":         "input",
		"kbps_output_option_series_data": data.OutputKbpsData,
		"kbps_output_option_name":        "output",
		"hitrate_option_series_data":     []int{150, 200, 259, 360, 378, 450},
		"hitrate_option_name":            "hitrate",
	})
}

func (this *ClientPG) RedisConsole(c *gin.Context) {
	this.Context = c
	//user := this.GetUser()
	redisId := this.Context.Param("redisId")

	this.Context.HTML(http.StatusOK, "console.html", gin.H{
		"redisId": redisId,
	})
}

//redis 服务信息页面
func (this *ClientPG) RedisServerInfo(c *gin.Context) {

	this.Context = c
	user := this.GetUser()
	redisId := this.Context.Param("redisId")
	rid := common.Str2Int64(redisId)

	var rinfoServers service.RedisInfoService = new(service.RedisInfo)
	data, err := rinfoServers.GetInfo(user, rid)
	if err != nil {
		manlog.Error(err)
	}

	this.Context.HTML(http.StatusOK, "redis_servers.html", gin.H{
		"redisId": redisId,
		"user":    user,
		"title":   "redis服务器信息",
		"data":    data,
	})
}

//redis 服务配置页面
func (this *ClientPG) RedisServerConfig(c *gin.Context) {

	this.Context = c
	user := this.GetUser()
	redisId := this.Context.Param("redisId")
	rid := common.Str2Int64(redisId)

	var rinfoServers service.RedisInfoService = new(service.RedisInfo)
	datas, err := rinfoServers.GetConfig(user, rid)
	if err != nil {
		manlog.Error(err)
	}

	this.Context.HTML(http.StatusOK, "redis_servers_config.html", gin.H{
		"redisId": redisId,
		"user":    user,
		"title":   "redis服务器配置",
		"data":    datas,
	})
}

//redis 慢日志
func (this *ClientPG) RedisSlowLog(c *gin.Context) {

	this.Context = c
	user := this.GetUser()
	redisId := this.Context.Param("redisId")
	rid := common.Str2Int64(redisId)

	var slowlogServers service.RedisSlowLogService = new(service.RedisSlowLog)
	slowlogServers.Get(user, rid)
	// if err != nil {
	// 	manlog.Error(err)
	// }

	// this.Context.HTML(http.StatusOK, "redis_servers_config.html", gin.H{
	// 	"redisId": redisId,
	// 	"user":    user,
	// 	"title":   "redis服务器配置",
	// 	"data":    datas,
	// })
	this.APIOutPut(0, "", "ok")
}
