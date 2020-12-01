package router

import (
	"encoding/binary"
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mangenotwork/mange_redis_manage/controllers"
	"github.com/mangenotwork/mange_redis_manage/install"
)

var Router *gin.Engine

func init() {
	Router = gin.Default()
}

var secrets = gin.H{
	"foo":    gin.H{"email": "foo@bar.com", "phone": "123433"},
	"austin": gin.H{"email": "austin@example.com", "phone": "666"},
	"lena":   gin.H{"email": "lena@guapa.com", "phone": "523443"},
}

func Routers() *gin.Engine {

	//静态目录配置
	// Router.Static("/static", "static")
	// Router.Static("/install/static", "static")
	Router.StaticFS("/static", http.Dir("./static"))

	//模板
	Router.LoadHTMLGlob("views/**/*")

	//全局跨域
	//Router.Use(CrosHandler())

	//全局中间件 检查安装
	Router.Use(install.CheckInstall())

	//redis相关交互
	API_Redis_Auth := Router.Group("/api/redis/auth", new(controllers.BaseController).Authentication()) //LogMiddleware())
	{
		API_Redis_Auth.GET("/conn/list", new(controllers.RedisAPI).ConnList)
		API_Redis_Auth.GET("/infos/:redisId", new(controllers.RedisAPI).RedisInfos)                     //获取redis服务器的信息
		API_Redis_Auth.GET("/client/list/:redisId", new(controllers.RedisAPI).GetRedisAllClient)        //获取连接redis的所有客户端
		API_Redis_Auth.POST("/conn/new", new(controllers.RedisAPI).ConnRedis)                           //新的redis连接
		API_Redis_Auth.GET("/memory/chart/:redisId", new(controllers.RedisAPI).RedisMemoryChartData)    //redis 服务的memory 图表数据
		API_Redis_Auth.GET("/dbtree/init/:redisId", new(controllers.RedisAPI).RedisDBTree)              //获取Redis的db树
		API_Redis_Auth.GET("/db/info/:redisId", new(controllers.RedisAPI).RedisDBInfo)                  //获取db info
		API_Redis_Auth.GET("/keytree/:redisId", new(controllers.RedisAPI).RedisKeyTree)                 //获取key的树
		API_Redis_Auth.GET("/keysearch/:redisId", new(controllers.RedisAPI).RedisKeySearch)             //获取key的树
		API_Redis_Auth.GET("/realtime/:redisId", new(controllers.RedisAPI).RedisRealtime)               //获取实时监控数据
		API_Redis_Auth.GET("/keys/info/:redisId", new(controllers.RedisAPI).RedisKeyInfo)               //获取key的信息与key的值
		API_Redis_Auth.GET("/keys/modify/name/:redisId", new(controllers.RedisAPI).RedisKeyModifyName)  //修改key名称
		API_Redis_Auth.POST("/keys/update/:redisId", new(controllers.RedisAPI).RedisKeyUpdate)          //更新key的值 TODO
		API_Redis_Auth.GET("/keys/modify/ttl/:redisId", new(controllers.RedisAPI).RedisKeyModifyTTL)    //更新key的ttl
		API_Redis_Auth.GET("/keys/delete/:redisId", new(controllers.RedisAPI).RedisKeyDelete)           //删除key
		API_Redis_Auth.GET("/keys/copy/todb/:redisId", new(controllers.RedisAPI).RedisKeyCopy2DB)       //将key复制到指定db
		API_Redis_Auth.GET("/keys/move/todb/:redisId", new(controllers.RedisAPI).RedisKeyMove2DB)       //将key转移到指定db TODO
		API_Redis_Auth.GET("/keys/copy/toredis/:redisId", new(controllers.RedisAPI).RedisKeyCopy2Redis) //将key复制到指定redis TODO
		API_Redis_Auth.GET("/keys/move/toredis/:redisId", new(controllers.RedisAPI).RedisKeyMove2Redis) //将key转移到制定redis TODO
		API_Redis_Auth.POST("/keys/create/:redisId", new(controllers.RedisAPI).RedisCreateNewKey)       //新建key
		API_Redis_Auth.POST("/console", new(controllers.RedisAPI).RedisConsole)                         //redis 命令终端
	}

	API_Redis_Allow := Router.Group("/api/redis/allow")
	{
		API_Redis_Allow.GET("/say", new(controllers.RedisAPI).Say)
		API_Redis_Allow.POST("/say", new(controllers.RedisAPI).Say)
		API_Redis_Allow.POST("/conn/test", new(controllers.RedisAPI).ConnRedisT)
	}

	//用户接口
	API_User_V1 := Router.Group("/api/user/v1")
	{
		API_User_V1.GET("/say", new(controllers.UserAPI).Say)
		API_User_V1.POST("/login", new(controllers.UserAPI).Login)
	}

	//管理员与系统管理相关接口
	API_Admin_V1 := Router.Group("/api/admit/v1")
	{
		API_Admin_V1.Use(new(controllers.BaseController).Authentication())
		API_Admin_V1.GET("/say", new(controllers.AdminAPI).Say)
		API_Admin_V1.GET("/cache/all", new(controllers.AdminAPI).CacheGetAll)    //查看所有缓存
		API_Admin_V1.GET("/cache/get", new(controllers.AdminAPI).CacheGet)       //查看指定缓存
		API_Admin_V1.GET("/cache/update", new(controllers.AdminAPI).CacheUpdate) //修改指定缓存
		API_Admin_V1.GET("/cache/del", new(controllers.AdminAPI).CacheDel)       //删除指定缓存
		API_Admin_V1.GET("/cache/add", new(controllers.AdminAPI).CacheAdd)       //新增缓存

	}

	Install := Router.Group("/install")
	{
		Install.GET("/say", new(controllers.AdminAPI).Say)
		Install.GET("/index", install.InstallPG)
		Install.GET("/run", install.Run)
	}

	//页面
	{
		Router.GET("/say", new(controllers.ClientPG).Say)                                                          //测试页面
		Router.GET("/", new(controllers.ClientPG).Index)                                                           //首页
		Router.GET("/notlogin", new(controllers.ClientPG).NotLogin)                                                //未登录的跳转
		Router.GET("/echarts/test/:id", new(controllers.ClientPG).Echarts)                                         //测试echar
		Router.GET("/echarts/memory/:redisId", controllers.AuthPG(), new(controllers.ClientPG).EchartsRedisMemory) //内存页面
		Router.GET("/client/list/:redisId", controllers.AuthPG(), new(controllers.ClientPG).RedisClientList)       //当前redis客户端连接列表页面
		Router.GET("/operation/:redisId", controllers.AuthPG(), new(controllers.ClientPG).RedisOperation)          //redis操作交互页面
		Router.GET("/realtime/:redisId", controllers.AuthPG(), new(controllers.ClientPG).RedisRealtime)            //实时监控页面
		Router.GET("/console/:redisId", controllers.AuthPG(), new(controllers.ClientPG).RedisConsole)              //redis 终端页面
		Router.GET("/serverinfo/:redisId", controllers.AuthPG(), new(controllers.ClientPG).RedisServerInfo)        //redis 服务信息页面
		Router.GET("/serverconfig/:redisId", controllers.AuthPG(), new(controllers.ClientPG).RedisServerConfig)    //redis 服务配置页面
		Router.GET("/slowlog/:redisId", controllers.AuthPG(), new(controllers.ClientPG).RedisSlowLog)              //redis 慢日志
	}

	//如下是测试 BasicAuth 中间件
	authorized := Router.Group("/admin", gin.BasicAuth(gin.Accounts{
		"foo":    "bar",
		"austin": "1234",
		"lena":   "hello2",
		"manu":   "4321",
	}), LogMiddleware())

	authorized.GET("/secrets", func(c *gin.Context) {
		// get user, it was set by the BasicAuth middleware
		user := c.MustGet(gin.AuthUserKey).(string)
		if secret, ok := secrets[user]; ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "secret": secret})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
		}
	})

	authorized.GET("/s", func(c *gin.Context) {
		// get user, it was set by the BasicAuth middleware
		user := c.MustGet(gin.AuthUserKey).(string)
		if secret, ok := secrets[user]; ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "secret": secret})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
		}
	})

	//404
	Router.NoRoute(func(ctx *gin.Context) {
		ctx.HTML(http.StatusNotFound, "404.html", "")
	})

	//401
	Router.NoRoute(func(ctx *gin.Context) {
		ctx.String(http.StatusUnauthorized, "未授权的访问")
	})

	//403
	Router.NoRoute(func(ctx *gin.Context) {
		ctx.HTML(http.StatusForbidden, "404.html", "")
	})

	return Router
}

//中间件
func LogMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		fmt.Println("请求log 中间件")
		fmt.Println(c)
		fmt.Println("Request = ", c.Request) // *http.Request
		ips := RemoteIp(c.Request)
		fmt.Println(ips)
		fmt.Println("Method = ", c.Request.Method)
		fmt.Println("URL = ", c.Request.URL)
		fmt.Println("Proto = ", c.Request.Proto)
		fmt.Println("ProtoMajor = ", c.Request.ProtoMajor)
		fmt.Println("ProtoMinor = ", c.Request.ProtoMinor)
		fmt.Println("Header = ", c.Request.Header)
		fmt.Println("Body = ", c.Request.Body)
		fmt.Println("ContentLength = ", c.Request.ContentLength)
		fmt.Println("TransferEncoding = ", c.Request.TransferEncoding)
		fmt.Println("Close = ", c.Request.Close)
		fmt.Println("Host = ", c.Request.Host)
		fmt.Println("PostForm = ", c.Request.PostForm)
		fmt.Println("Form = ", c.Request.Form)
		fmt.Println("MultipartForm = ", c.Request.MultipartForm)
		fmt.Println("Trailer = ", c.Request.Trailer)
		fmt.Println("RemoteAddr = ", c.Request.RemoteAddr)
		fmt.Println("RequestURI = ", c.Request.RequestURI)
		fmt.Println("TLS = ", c.Request.TLS)

		fmt.Println("ClientIP = ", c.ClientIP())       //ClientIP()
		fmt.Println("ContentType = ", c.ContentType()) //ContentType()
		fmt.Println("Params = ", c.Params)             //Params

		fmt.Println("Keys = ", c.Keys)         //Keys map[string]interface{}
		fmt.Println("Accepted = ", c.Accepted) // Accepted []string

		fmt.Println("Handler = ", c.Handler())
		fmt.Println("HandlerName = ", c.HandlerName())
		fmt.Println("HandlerNames = ", c.HandlerNames())

		fmt.Println("Accept = ", c.GetHeader("Accept"))
		fmt.Println("Accept-Encoding = ", c.GetHeader("Accept-Encoding"))
		fmt.Println("Accept-Language = ", c.GetHeader("Accept-Language"))
		fmt.Println("Connection = ", c.GetHeader("Connection"))
		fmt.Println("User-Agent = ", c.GetHeader("User-Agent"))

		c.Next()
	}
}

// RemoteIp 返回远程客户端的 IP，如 192.168.1.1
func RemoteIp(req *http.Request) string {
	remoteAddr := req.RemoteAddr
	if ip := req.Header.Get("X-Real-IP"); ip != "" {
		remoteAddr = ip
	} else if ip = req.Header.Get("X-Forwarded-For"); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}
	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}
	return remoteAddr
}

// Ip2long 将 IPv4 字符串形式转为 uint32
func Ip2long(ipstr string) uint32 {
	ip := net.ParseIP(ipstr)
	if ip == nil {
		return 0
	}
	ip = ip.To4()
	return binary.BigEndian.Uint32(ip)
}

// //跨域访问：cross  origin resource share
// func CrosHandler() gin.HandlerFunc {
// 	return func(context *gin.Context) {

// 		context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
// 		context.Header("Access-Control-Allow-Origin", "*") // 设置允许访问所有域
// 		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
// 		context.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma,token,openid,opentoken")
// 		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
// 		context.Header("Access-Control-Max-Age", "172800")
// 		context.Header("Access-Control-Allow-Credentials", "false")
// 		context.Set("content-type", "application/json")

// 		// method := context.Request.Method
// 		// if method == "OPTIONS" {
// 		// 	//doing...
// 		// }

// 		//处理请求
// 		context.Next()
// 	}
// }
