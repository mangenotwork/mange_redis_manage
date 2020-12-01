//
//	redis交互接口
//

package controllers

import (
	_ "net/http"
	_ "time"

	"github.com/gin-gonic/gin"
	"github.com/mangenotwork/mange_redis_manage/common"
	"github.com/mangenotwork/mange_redis_manage/common/manlog"
	"github.com/mangenotwork/mange_redis_manage/service"
	"github.com/mangenotwork/mange_redis_manage/structs"
)

type RedisAPI struct {
	BaseController
}

func (this *RedisAPI) Say(c *gin.Context) {
	// manlog.Debug("aaaaaaaaaaaaaa", 3, true, "aa")
	// manlog.Error("aaaaaaaaaaaaaa", 3, true, "aa", []string{"a"})
	//manlog.Panic("aaaaaaaaaaaaaa", 3, true, "aa", []string{"a"})
	this.Context = c

	// var rcServers service.RedisConnService
	// rc := new(service.RedisConn)
	// rcServers = rc

	// rcServers.GetAll()

	//c.String(http.StatusOK, "this is redis api")

	//输出json
	// c.JSON(http.StatusOK, structs.ResponseJson{
	// 	Code:      0,
	// 	Mag:       "pass",
	// 	Date:      []string{"a", "1"},
	// 	TimeStamp: time.Now().UnixNano(),
	// })

	//设置cookie
	c.SetCookie("sign", common.NewMangeSign(), 60*60*24*7, "/", "localhost", false, true)
	c.SetAccepted("Man Ge Redis Manage")
	//func (c *Context) Header(key, value string)
	c.Header("ManGe", "ManGe Redis Manage")

	// //输出美化的json
	// c.IndentedJSON(http.StatusOK, structs.ResponseJson{
	// 	Code:      0,
	// 	Mag:       "pass",
	// 	Date:      []string{"a", "1"},
	// 	TimeStamp: time.Now().Unix(),
	// })

	d := new(Data)
	err := this.GetPostArgs(&d)
	if err != nil {
		return
	}

	manlog.Debug(*d)
	this.APIOutPut(0, "pass", []string{"a", "1"})
	return
}

type Data struct {
	AA string `json:"a"`
	BB int64  `json:"b"`
	CC bool   `json"c"`
}

func (this *RedisAPI) ConnRedis(c *gin.Context) {
	//接收上下文
	this.Context = c
	//接收参数
	redis_conn_data := new(structs.RedisConnData)
	if err := this.GetPostArgs(&redis_conn_data); err != nil {
		return
	}

	//获取user
	user := this.GetUserParameter()

	//实例化接口
	var rcServers service.RedisConnService = new(service.RedisConn)
	callback := rcServers.New(redis_conn_data, user.UserID)

	this.APIOutPut(0, "", callback)
	return
}

func (this *RedisAPI) ConnRedisT(c *gin.Context) {
	this.Context = c
	redis_conn_data := new(structs.RedisConnData)
	if err := this.GetPostArgs(&redis_conn_data); err != nil {
		return
	}

	rcServers := service.RedisConnServiceFunc()
	callback := rcServers.Detection(redis_conn_data)

	// var rcServers service.RedisConnService = new(service.RedisConn)
	// callback := rcServers.Detection(redis_conn_data)

	this.APIOutPut(0, "", callback)
	return
}

func (this *RedisAPI) ConnList(c *gin.Context) {
	this.Context = c
	user := this.GetUserParameter()
	manlog.Debug(*user)
	//非登录状态返回状态
	var rcServers service.RedisConnService = new(service.RedisConn)
	datas := rcServers.GetAll(user)
	this.APIOutPut(0, "", datas)
}

func (this *RedisAPI) RedisInfos(c *gin.Context) {
	this.Context = c
	user := this.GetUserParameter()
	manlog.Debug(*user)

	redisId := this.Context.Param("redisId")
	rid := common.Str2Int64(redisId)
	//非登录状态返回状态
	var rcServers service.RedisConnService = new(service.RedisConn)
	datas := rcServers.GetRedisInfos(user, rid)
	this.APIOutPut(0, "", datas)
}

func (this *RedisAPI) GetRedisAllClient(c *gin.Context) {
	this.Context = c
	user := this.GetUserParameter()
	manlog.Debug(*user)

	redisId := this.Context.Param("redisId")
	rid := common.Str2Int64(redisId)
	var rcServers service.RedisConnService = new(service.RedisConn)
	datas := rcServers.GetAllClient(user, rid)
	this.APIOutPut(0, "", datas)
}

func (this *RedisAPI) RedisMemoryChartData(c *gin.Context) {
	this.Context = c
	user := this.GetUserParameter()
	manlog.Debug(*user)

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
	//datas, _ := rcServers.GetMemoryChartData(user, rid)
	datas, _ := rcServers.GetEchartsRedisMemoryData(user, rid, hours, day)
	this.APIOutPut(0, "", datas)
}

func (this *RedisAPI) RedisDBTree(c *gin.Context) {
	this.Context = c
	user := this.GetUserParameter()
	manlog.Debug(*user)

	redisId := this.Context.Param("redisId")
	rid := common.Str2Int64(redisId)

	var rcServers service.RedisConnService = new(service.RedisConn)
	data := rcServers.GetRedisDBTree(user, rid)

	this.APIOutPut(0, "", data)
}

func (this *RedisAPI) RedisDBInfo(c *gin.Context) {
	this.Context = c
	user := this.GetUserParameter()
	manlog.Debug(*user)

	redisId := this.Context.Param("redisId")
	rid := common.Str2Int64(redisId)
	//获取db信息数据

	var rcServers service.RedisConnService = new(service.RedisConn)
	dblist := rcServers.GetRedisDBList(user, rid)

	//验证参数db, 如果没有参数取所有
	db := this.Context.Query("db")
	if db == "" {
		this.APIOutPut(0, "", dblist)
		return
	}

	//如果有传入 db则只返回db的信息
	dbid := common.Str2Int64(db)
	for _, v := range dblist {
		if v.DBID == dbid {
			this.APIOutPut(0, "", v)
			return
		}
	}

	this.APIOutPutError(1, "未获取到db信息")
	return
}

func (this *RedisAPI) RedisKeyTree(c *gin.Context) {
	this.Context = c
	user := this.GetUserParameter()
	manlog.Debug(*user)

	redisId := this.Context.Param("redisId")
	rid := common.Str2Int64(redisId)
	var dbid int64 = 0
	db := this.Context.Query("db")
	if db != "" {
		dbid = common.Str2Int64(db)
	}

	match := this.Context.Query("match")

	var rcServers service.RedisConnService = new(service.RedisConn)
	data := rcServers.GetRedisKeyTree(user, rid, dbid, match)

	this.APIOutPut(0, "", data)
}

func (this *RedisAPI) RedisKeySearch(c *gin.Context) {
	this.Context = c
	user := this.GetUserParameter()
	manlog.Debug(*user)

	redisId := this.Context.Param("redisId")
	rid := common.Str2Int64(redisId)
	var dbid int64 = 0
	db := this.Context.Query("db")
	if db != "" {
		dbid = common.Str2Int64(db)
	}

	match := this.Context.Query("match")
	var rcServers service.RedisConnService = new(service.RedisConn)
	data := rcServers.GetRedisKeySearch(user, rid, dbid, match)

	this.APIOutPut(0, "", data)

}

func (this *RedisAPI) RedisRealtime(c *gin.Context) {
	this.Context = c
	user := this.GetUserParameter()
	redisId := this.Context.Param("redisId")
	rid := common.Str2Int64(redisId)
	rtid := this.Context.Query("rtid")
	var rinfoServers service.RedisInfoService = new(service.RedisInfo)
	data, err := rinfoServers.RedisRealTime(user, rid, rtid)
	if err != nil {
		manlog.Error(err)
	}
	this.APIOutPut(0, "", data)
}

func (this *RedisAPI) RedisKeyInfo(c *gin.Context) {
	this.Context = c
	user := this.GetUserParameter()
	redisId := this.Context.Param("redisId")
	rid := common.Str2Int64(redisId)
	var dbid int64 = 0
	db := this.Context.Query("db")
	if db != "" {
		dbid = common.Str2Int64(db)
	}
	key := this.Context.Query("key")
	var rinfoServers service.RedisInfoService = new(service.RedisInfo)
	data, err := rinfoServers.GetKeyInfo(user, rid, dbid, key)
	if err != nil {
		manlog.Error(err)
	}
	this.APIOutPut(0, "", data)
}

//修改key名称
func (this *RedisAPI) RedisKeyModifyName(c *gin.Context) {
	this.Context = c
	user := this.GetUserParameter()
	redisId := this.Context.Param("redisId")
	rid := common.Str2Int64(redisId)
	var dbid int64 = 0
	db := this.Context.Query("db")
	if db != "" {
		dbid = common.Str2Int64(db)
	}
	key := this.Context.Query("key")
	newname := this.Context.Query("newname")
	if key == "" || newname == "" {
		this.APIOutPut(2, "", "key与新命名不能为空")
		return
	}

	if key == newname {
		this.APIOutPut(0, "", "ok")
		return
	}

	var ropeServers service.RedisOperateService = new(service.RedisOperate)
	data, err := ropeServers.ModifyKeyName(user, rid, dbid, key, newname)
	if err != nil {
		manlog.Error(err)
	}
	if data {
		this.APIOutPut(0, "", "ok")
		return
	}
	this.APIOutPut(2, "", "falied")
}

//更新key的值 TODO
func (this *RedisAPI) RedisKeyUpdate(c *gin.Context) {
	this.APIOutPut(0, "", "ok")
}

//更新key的ttl
func (this *RedisAPI) RedisKeyModifyTTL(c *gin.Context) {

	this.Context = c
	user := this.GetUserParameter()
	redisId := this.Context.Param("redisId")
	rid := common.Str2Int64(redisId)
	var dbid int64 = 0
	db := this.Context.Query("db")
	if db != "" {
		dbid = common.Str2Int64(db)
	}
	key := this.Context.Query("key")
	ttltype := this.Context.Query("type") //0:按秒   1:按日期
	ttl := this.Context.Query("ttl")
	if key == "" || ttltype == "" || ttl == "" {
		this.APIOutPut(2, "", "传入参数不能为空")
		return
	}

	manlog.Debug("ttl = ", ttl)

	var ropeServers service.RedisOperateService = new(service.RedisOperate)
	data, err := ropeServers.ModifyKeyTTL(user, rid, dbid, key, ttltype, ttl)
	if err != nil {
		manlog.Error(err)
	}
	if data {
		this.APIOutPut(0, "", "ok")
		return
	}
	this.APIOutPut(2, "", "falied")
}

//删除key
func (this *RedisAPI) RedisKeyDelete(c *gin.Context) {

	this.Context = c
	user := this.GetUserParameter()
	redisId := this.Context.Param("redisId")
	rid := common.Str2Int64(redisId)
	var dbid int64 = 0
	db := this.Context.Query("db")
	if db != "" {
		dbid = common.Str2Int64(db)
	}
	key := this.Context.Query("key")
	if key == "" {
		this.APIOutPut(2, "", "传入参数不能为空")
		return
	}

	var ropeServers service.RedisOperateService = new(service.RedisOperate)
	data, err := ropeServers.DeleteKey(user, rid, dbid, key)
	if err != nil {
		manlog.Error(err)
	}
	if data {
		this.APIOutPut(0, "", "ok")
		return
	}
	this.APIOutPut(2, "", "falied")
}

//将key复制到指定db TODO
func (this *RedisAPI) RedisKeyCopy2DB(c *gin.Context) {

	this.Context = c
	user := this.GetUserParameter()
	redisId := this.Context.Param("redisId")
	rid := common.Str2Int64(redisId)
	var dbid int64 = 0
	db := this.Context.Query("db")
	if db != "" {
		dbid = common.Str2Int64(db)
	}
	key := this.Context.Query("key")
	todb := this.Context.Query("todb")
	if key == "" || todb == "" {
		this.APIOutPut(2, "", "传入参数不能为空")
		return
	}
	todb_value := common.Str2Int64(todb)
	if dbid == todb_value {
		this.APIOutPut(2, "", "无法复制到同db")
		return
	}

	var ropeServers service.RedisOperateService = new(service.RedisOperate)
	data, err := ropeServers.KeyCopy2DB(user, rid, dbid, todb_value, key)
	if err != nil {
		manlog.Error(err)
	}
	if data {
		this.APIOutPut(0, "", "ok")
		return
	}
	this.APIOutPut(2, "", "falied")
}

//将key转移到指定db TODO
func (this *RedisAPI) RedisKeyMove2DB(c *gin.Context) {
	this.APIOutPut(0, "", "ok")
}

//将key复制到指定redis TODO
func (this *RedisAPI) RedisKeyCopy2Redis(c *gin.Context) {
	this.APIOutPut(0, "", "ok")
}

//将key转移到制定redis TODO
func (this *RedisAPI) RedisKeyMove2Redis(c *gin.Context) {
	this.APIOutPut(0, "", "ok")
}

func (this *RedisAPI) RedisCreateNewKey(c *gin.Context) {

	manlog.Debug("RedisCreateNewKey")

	this.Context = c
	//接收参数
	newkeydata := new(structs.CreateKeyPostData)
	if err := this.GetPostArgs(&newkeydata); err != nil {
		return
	}

	//获取user
	user := this.GetUserParameter()
	manlog.Debug(user)

	redisId := this.Context.Param("redisId")
	rid := common.Str2Int64(redisId)

	//新建key
	var ropeServers service.RedisOperateService = new(service.RedisOperate)
	err := ropeServers.CreateKey(user, rid, newkeydata)
	if err != nil {
		this.APIOutPutError(1, err.Error())
	}

	this.APIOutPut(0, "", "创建key成功")
}

func (this *RedisAPI) RedisConsole(c *gin.Context) {
	this.Context = c
	cmddata := new(structs.RedisConsoleData)
	if err := this.GetPostArgs(&cmddata); err != nil {
		return
	}
	//获取user
	user := this.GetUserParameter()
	var ropeServers service.RedisOperateService = new(service.RedisOperate)
	data := ropeServers.Console(user, cmddata)
	this.APIOutPut(0, "", data)
}
