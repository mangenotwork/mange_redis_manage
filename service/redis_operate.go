//
//	redis 操作服务
//
package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/garyburd/redigo/redis"
	"github.com/mangenotwork/mange_redis_manage/common"
	"github.com/mangenotwork/mange_redis_manage/common/manlog"
	manredis "github.com/mangenotwork/mange_redis_manage/common/redis"
	"github.com/mangenotwork/mange_redis_manage/structs"
)

//redis操作服务，对外提供的接口
type RedisOperateService interface {
	AllKeyList()                                                                                           //redis 所有key列表
	ModifyKeyName(user *structs.UserParameter, redisId, dbid int64, key, newkey string) (bool, error)      //修改key名称
	ModifyKeyTTL(user *structs.UserParameter, redisId, dbid int64, key, ttltype, ttl string) (bool, error) //修改key过期时间
	DeleteKey(user *structs.UserParameter, redisId, dbid int64, key string) (bool, error)                  //删除key
	KeyCopy2DB(user *structs.UserParameter, redisId, dbid, todb int64, key string) (bool, error)           //将key复制到指定db
	CreateKey(user *structs.UserParameter, redisId int64, newkeydata *structs.CreateKeyPostData) error     //新建key
	Console(user *structs.UserParameter, cmddata *structs.RedisConsoleData) (rdata interface{})            //redis 命令行终端
}

//redis的操作
type RedisOperate struct {
}

//redis 基础类型操作工厂
type RedisTypeOperateFactory struct {
}

//redis 基础类型操作接口， 主要对值的操作
type RedisTypeOperate interface {
	Get(redis.Conn, string) (interface{}, string, string, error) //获取key的值
	Create(redis.Conn, string, interface{}) error                //创建一个key
	Append(redis.Conn, string, interface{}) error                //在原有值上追加值
	Del()                                                        //删除key的值
	Update()                                                     //更新key的值
	ValueSize(redis.Conn, string) int64                          //获取key value的大小
}

//redis基础类型操作
func (this *RedisTypeOperateFactory) Operate(redistype string) RedisTypeOperate {
	switch redistype {
	case "string":
		return new(RedisString)
	case "hash":
		return new(RedisHash)
	case "list":
		return new(RedisList)
	case "set":
		return new(RedisSet)
	case "zset":
		return new(RedisZSet)
	default:
		//默认为string操作
		return new(RedisString)
	}
	return nil
}

func (this *RedisOperate) AllKeyList() {}

//修改key名称
func (this *RedisOperate) ModifyKeyName(user *structs.UserParameter, redisId, dbid int64, key, newkey string) (bool, error) {
	//获取连接
	rc, _, err := new(RedisConn).DBconn(user, redisId, dbid)
	if err != nil {
		manlog.Error(err)
		return false, err
	}

	return manredis.RenameKey(rc, key, newkey), nil

}

//修改key 的过期时间
func (this *RedisOperate) ModifyKeyTTL(user *structs.UserParameter, redisId, dbid int64, key, ttltype, ttl string) (bool, error) {
	//获取连接
	rc, _, err := new(RedisConn).DBconn(user, redisId, dbid)
	if err != nil {
		manlog.Error(err)
		return false, err
	}

	//0:按秒   1:按日期
	var ttl_value int64
	isok := false
	switch ttltype {
	case "0":
		//将ttl解析为秒
		ttl_value = common.Str2Int64(ttl)
		isok = manredis.UpdateKeyTTL(rc, key, ttl_value)
	case "1":
		//将tll解析为时间戳

		ttl = common.ReverseDate(ttl)
		manlog.Debug(ttl)

		ttl_value = common.Date2Unix(ttl)
		manlog.Debug(ttl_value)
		isok = manredis.EXPIREATKey(rc, key, ttl_value)
	}
	return isok, nil

}

func (this *RedisOperate) DeleteKey(user *structs.UserParameter, redisId, dbid int64, key string) (bool, error) {
	//获取连接
	rc, _, err := new(RedisConn).DBconn(user, redisId, dbid)
	if err != nil {
		manlog.Error(err)
		return false, err
	}

	return manredis.DELKey(rc, key), nil
}

func (this *RedisOperate) KeyCopy2DB(user *structs.UserParameter, redisId, dbid, todb int64, key string) (bool, error) {
	//获取连接
	rc, _, err := new(RedisConn).DBconn(user, redisId, dbid)
	if err != nil {
		manlog.Error(err)
		return false, err
	}
	//获取key类型和值
	value, keytype, _, _ := new(RedisInfo).GetKeyValue(rc, key)

	//切换db
	newrc, err := manredis.SelectDB(rc, todb)
	if err != nil {
		return false, err
	}

	//写获取的值
	manlog.Debug(value, keytype, newrc)

	return false, nil
}

func (this *RedisOperate) CreateKey(user *structs.UserParameter, redisId int64, newkeydata *structs.CreateKeyPostData) error {

	//获取连接
	rc, _, err := new(RedisConn).DBconn(user, redisId, newkeydata.DBID)
	if err != nil {
		manlog.Error(err)
		return err
	}

	//创建key
	var factory = new(RedisTypeOperateFactory)
	keyObj := factory.Operate(newkeydata.KeyType)
	adderr := keyObj.Create(rc, newkeydata.Key, newkeydata.Value)

	//设置过期时间
	if newkeydata.TTL > 0 {
		manredis.UpdateKeyTTL(rc, newkeydata.Key, newkeydata.TTL)
	}

	//效验错误
	if adderr != nil && adderr.Error() == "WRONGTYPE Operation against a key holding the wrong kind of value" {
		adderr = errors.New("key已存在")
	}

	return adderr
}

func (this *RedisOperate) Console(user *structs.UserParameter, cmddata *structs.RedisConsoleData) (rdata interface{}) {
	manlog.Debug(*user)
	manlog.Debug(*cmddata)
	rdata = "未知命令！"
	if cmddata.CMD == "help" {
		rdata = `<br> === redis 命令终端 help === <br>
				1. 支持所有redis原生命令,使用方式与redis-cli一样;<br>
				2. 支持以下特有命令:<br>
					- help : 打开帮助文档<br>
					- db : 切换db与select一样，使用db&lt;id&gt;<br>
					- clear : 清空<br>
				`
	}

	if cmddata.CMD == "mange" {
		rdata = `<br> ××× redis 命令终端作者：ManGe ××× <br>`
	}

	cmdArry := strings.Split(cmddata.CMD, " ")
	manlog.Debug(cmdArry)

	if len(cmdArry) == 0 {
		return
	}

	if cmdArry[0] == "man" {
		rdata = `使用命令man参数错误，正确使用： man help`
		if len(cmdArry) > 1 {
			rdata = ManCMD[cmdArry[1]]
		}
		if rdata == "" {
			rdata = `未知命令`
		}
	}

	if cmddata.CMD == "allcmd" {
		rdata = ""
		for k, _ := range ManCMD {
			rdata = fmt.Sprintf("%v>%s<br>", rdata, k)
		}
	}

	return

}

var ManCMDDisTemp = "%s<hr>---说明:<br>%s<hr>---使用(实例):<br>%s<hr>"

var ManCMD = map[string]string{
	"help":         fmt.Sprintf(ManCMDDisTemp, "help", "打开帮助文档", "help"),
	"db":           fmt.Sprintf(ManCMDDisTemp, "db", "切换当前redis的db", "db<id>"),
	"clear":        fmt.Sprintf(ManCMDDisTemp, "clear", "清空当前终端", "clear"),
	"allcmd":       fmt.Sprintf(ManCMDDisTemp, "allcmd", "列出所有redis支持的命令，可输入类别: key,server,string,hash,list,set,zet,pub", "allcmd <类别>"),
	"bgrewriteaof": fmt.Sprintf(ManCMDDisTemp, "BGREWRITEAOF", Doc_Bgrewriteaof, "BGREWRITEAOF"),
	"bgsave":       fmt.Sprintf(ManCMDDisTemp, "BGSAVE", Doc_Bgrewriteaof, "BGSAVE"),
	"client":       fmt.Sprintf(ManCMDDisTemp, "CLIENT", Doc_Client, "CLIENT GETNAME<br>CLIENT KILL ip:port<br>CLIENT LIST<br>CLIENT SETNAME connection-name<br>"),
	"config":       fmt.Sprintf(ManCMDDisTemp, "CONFIG", Doc_Config, "CONFIG GET parameter<br>CONFIG RESETSTAT<br>CONFIG REWRITE<br>CONFIG SET parameter value<br>"),
	"dbsize":       fmt.Sprintf(ManCMDDisTemp, "DBSIZE", Doc_Dbsize, "DBSIZE<br>"),
	"debug":        fmt.Sprintf(ManCMDDisTemp, "DEBUG", Doc_Debug, "DEBUG OBJECT key<br>DEBUG SEGFAULT<br>"),
	"flushall":     fmt.Sprintf(ManCMDDisTemp, "FLUSHALL", Doc_Flushall, "FLUSHALL"),
	"flushdb":      fmt.Sprintf(ManCMDDisTemp, "FLUSHDB", Doc_Flushdb, "FLUSHDB"),
	"info":         fmt.Sprintf(ManCMDDisTemp, "INFO", Doc_Info, "INFO [section]"),
	"lastsave":     fmt.Sprintf(ManCMDDisTemp, "LASTSAVE", Doc_Lastsave, "LASTSAVE"),
	"monitor":      fmt.Sprintf(ManCMDDisTemp, "MONITOR", Doc_Monitor, "MONITOR"),
	"psync":        fmt.Sprintf(ManCMDDisTemp, "PSYNC", Doc_Psync, "PSYNC ? -1"),
	"save":         fmt.Sprintf(ManCMDDisTemp, "SAVE", Doc_Save, "SAVE"),
	"shutdown":     fmt.Sprintf(ManCMDDisTemp, "SHUTDOWN", Doc_Shutdown, "SHUTDOWN"),
	"slaveof":      fmt.Sprintf(ManCMDDisTemp, "SLAVEOF", Doc_Slaveof, "SLAVEOF 127.0.0.1 6379"),
	"slowlog":      fmt.Sprintf(ManCMDDisTemp, "SLOWLOG", Doc_Slowlog, "SLOWLOG GET"),
	"sync":         fmt.Sprintf(ManCMDDisTemp, "SYNC", Doc_Sync, "SYNC"),
	"time":         fmt.Sprintf(ManCMDDisTemp, "TIME", Doc_Time, "TIME"),
}

//bgrewriteaof 描述
var Doc_Bgrewriteaof string = `执行一个 AOF文件 重写操作。重写会创建一个当前 AOF 文件的体积优化版本。<br>
			即使 BGREWRITEAOF 执行失败，也不会有任何数据丢失，因为旧的 AOF 文件在 BGREWRITEAOF 成功之前不会被修改。<br>
			重写操作只会在没有其他持久化工作在后台执行时被触发，也就是说：<br>
			&ensp;&ensp;&ensp;&ensp; 如果 Redis 的子进程正在执行快照的保存工作，那么 AOF 重写的操作会被预定(scheduled)，等到保存工作完成之后再执行 AOF 重写。在这种情况下，<br>
		    BGREWRITEAOF 的返回值仍然是 OK ，但还会加上一条额外的信息，说明 BGREWRITEAOF 要等到保存操作完成之后才能执行。在 Redis 2.6 或以上的版本，<br>
			可以使用 INFO 命令查看 BGREWRITEAOF 是否被预定。<br>
			&ensp;&ensp;&ensp;&ensp; 如果已经有别的 AOF 文件重写在执行，那么 BGREWRITEAOF 返回一个错误，并且这个新的 BGREWRITEAOF 请求也不会被预定到下次执行。<br>
			从 Redis 2.4 开始， AOF 重写由 Redis 自行触发， BGREWRITEAOF 仅仅用于手动触发重写操作。`

var Doc_Bgsave string = `在后台异步(Asynchronously)保存当前数据库的数据到磁盘。<br>
			BGSAVE 命令执行之后立即返回 OK ，然后 Redis fork 出一个新子进程，原来的 Redis 进程(父进程)继续处理客户端请求，而子进程则负责将数据保存到磁盘，然后退出。<br>
			客户端可以通过 LASTSAVE 命令查看相关信息，判断 BGSAVE 命令是否执行成功。`

var Doc_Client string = `CLIENT GETNAME<br>
			返回 CLIENT SETNAME 命令为连接设置的名字。<br>
			因为新创建的连接默认是没有名字的， 对于没有名字的连接， CLIENT GETNAME 返回空白回复。<br><br>
			CLIENT KILL<br>
			关闭地址为 ip:port 的客户端。<br>
			ip:port 应该和 CLIENT LIST 命令输出的其中一行匹配。<br>
			因为 Redis 使用单线程设计，所以当 Redis 正在执行命令的时候，不会有客户端被断开连接。<br><br>
			CLIENT LIST<br>
			以人类可读的格式，返回所有连接到服务器的客户端信息和统计数据<br><br>
			CLIENT SETNAME<br>
			为当前连接分配一个名字。<br>
			这个名字会显示在 CLIENT LIST 命令的结果中， 用于识别当前正在与服务器进行连接的客户端<br><br>。
`

var Doc_Config string = `CONFIG GET parameter<br>
			CONFIG GET 命令用于取得运行中的 Redis 服务器的配置参数(configuration parameters)，在 Redis 2.4 版本中， 有部分参数没有办法用 <br>
			CONFIG GET 访问，但是在最新的 Redis 2.6 版本中，所有配置参数都已经可以用 CONFIG GET 访问了。<br>
			CONFIG GET 接受单个参数 parameter 作为搜索关键字，查找所有匹配的配置参数，其中参数和值以“键-值对”(key-value pairs)的方式排列。<br><br>
			重置 INFO 命令中的某些统计数据，包括：<br>
			Keyspace hits (键空间命中次数)<br>
			Keyspace misses (键空间不命中次数)<br>
			Number of commands processed (执行命令的次数)<br>
			Number of connections received (连接服务器的次数)<br>
			Number of expired keys (过期key的数量)<br>
			Number of rejected connections (被拒绝的连接数量)<br>
			Latest fork(2) time(最后执行 fork(2) 的时间)<br>
			The aof_delayed_fsync counter(aof_delayed_fsync 计数器的值<br><br>
			CONFIG REWRITE<br>
			CONFIG REWRITE 命令对启动 Redis 服务器时所指定的 redis.conf 文件进行改写： 因为 CONFIG SET 命令可以对服务器的当前配置进行修改， <br>
			而修改后的配置可能和 redis.conf 文件中所描述的配置不一样， CONFIG REWRITE 的作用就是通过尽可能少的修改， 将服务器当前所使用的配置记录到 redis.conf 文件中。<br><br>
			CONFIG SET parameter value<br>
			CONFIG SET 命令可以动态地调整 Redis 服务器的配置(configuration)而无须重启。<br>
			你可以使用它修改配置参数，或者改变 Redis 的持久化(Persistence)方式。<br>
			CONFIG SET 可以修改的配置参数可以使用命令 CONFIG GET * 来列出，所有被 CONFIG SET 修改的配置参数都会立即生效。<br>
			关于 CONFIG SET 命令的更多消息，请参见命令 CONFIG GET 的说明。<br>
			关于如何使用 CONFIG SET 命令修改 Redis 持久化方式，请参见 Redis Persistence 。<br><br>
`

var Doc_Dbsize string = `DBSIZE<br>
			返回当前数据库的 key 的数量。<br>
`

var Doc_Debug string = `DEBUG OBJECT key<br>
			DEBUG OBJECT 是一个调试命令，它不应被客户端所使用。<br><br>
			DEBUG SEGFAULT<br>
			执行一个不合法的内存访问从而让 Redis 崩溃，仅在开发时用于 BUG 模拟。<br><br>
`

var Doc_Flushall string = `FLUSHALL<br>
			清空整个 Redis 服务器的数据(删除所有数据库的所有 key )。<br>
			此命令从不失败。<br>
`

var Doc_Flushdb string = `FLUSHDB<br>
			清空当前数据库中的所有 key。<br>
			此命令从不失败。<br>
`

var Doc_Info string = `INFO [section]<br>
			以一种易于解释（parse）且易于阅读的格式，返回关于 Redis 服务器的各种信息和统计数值。<br>
`

var Doc_Lastsave string = `LASTSAVE<br>
			返回最近一次 Redis 成功将数据保存到磁盘上的时间，以 UNIX 时间戳格式表示。<br>

`

var Doc_Monitor string = `MONITOR<br>
			实时打印出 Redis 服务器接收到的命令，调试用。<br>
			总是返回 OK 。<br>
`

var Doc_Psync string = `PSYNC <MASTER_RUN_ID> <OFFSET><br>
			用于复制功能(replication)的内部命令。<br>
`

var Doc_Save string = `
			SAVE 命令执行一个同步保存操作，将当前 Redis 实例的所有数据快照(snapshot)以 RDB 文件的形式保存到硬盘。<br>
			一般来说，在生产环境很少执行 SAVE 操作，因为它会阻塞所有客户端，保存数据库的任务通常由 BGSAVE 命令异步地执行。<br>
			然而，如果负责保存数据的后台子进程不幸出现问题时， SAVE 可以作为保存数据的最后手段来使用。<br>
`

var Doc_Shutdown string = `SHUTDOWN 命令执行以下操作：<br>
			停止所有客户端<br>
			如果有至少一个保存点在等待，执行 SAVE 命令<br>
			如果 AOF 选项被打开，更新 AOF 文件<br>
			关闭 redis 服务器(server)<br>
			如果持久化被打开的话， SHUTDOWN 命令会保证服务器正常关闭而不丢失任何数据。<br>
`

var Doc_Slaveof string = `SLAVEOF host port<br>
			SLAVEOF 命令用于在 Redis 运行时动态地修改复制(replication)功能的行为。<br>
			通过执行 SLAVEOF host port 命令，可以将当前服务器转变为指定服务器的从属服务器(slave server)。<br>
`

var Doc_Slowlog string = `
			Slow log 是 Redis 用来记录查询执行时间的日志系统。<br>
			查询执行时间指的是不包括像客户端响应(talking)、发送回复等 IO 操作，而单单是执行一个查询命令所耗费的时间。<br>
			另外，slow log 保存在内存里面，读写速度非常快，因此你可以放心地使用它，不必担心因为开启 slow log 而损害 Redis 的速度。<br>
`

var Doc_Sync string = `SYNC<br>
			用于复制功能(replication)的内部命令。<br>

`

var Doc_Time string = `TIME<br>
			返回当前服务器时间。<br>
			返回值：一个包含两个字符串的列表： 第一个字符串是当前时间(以 UNIX 时间戳格式表示)，而第二个字符串是当前这一秒钟已经逝去的微秒数。<br>
`
