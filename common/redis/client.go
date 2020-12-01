//
//	与连接redis客户端相关的信息与操作
//

package redis

import (
	"strings"
	_ "sync"
	_ "time"

	"github.com/garyburd/redigo/redis"

	// "github.com/go-ini/ini"
	"github.com/mangenotwork/mange_redis_manage/common"
	"github.com/mangenotwork/mange_redis_manage/common/manlog"
	"github.com/mangenotwork/mange_redis_manage/dao"
)

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
//
//CLIENT SETNAME connection-name
//为当前连接分配一个名字。
//这个名字会显示在 CLIENT LIST 命令的结果中， 用于识别当前正在与服务器进行连接的客户端。
//CLIENT SETNAME ""     来清空
//
//

//为当前连接分配一个名字。
func SetRedisClientName() {}

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
//
// CLIENT LIST
// 以人类可读的格式，返回所有连接到服务器的客户端信息和统计数据。
//参数参考  http://doc.redisfans.com/server/client_list.html
//

//Redis 客户端信息
// type RedisClientInfo struct {
// 	ID       string //
// 	Addr     string // 客户端的地址和端口
// 	Fd       string // 套接字所使用的文件描述符
// 	Name     string //
// 	Age      string // 以秒计算的已连接时长
// 	Idle     string // 以秒计算的空闲时长
// 	Flags    string // 客户端 flag
// 	Db       string // 该客户端正在使用的数据库 ID
// 	Sub      string // 已订阅频道的数量
// 	Psub     string // 已订阅模式的数量
// 	Multi    string // 在事务中被执行的命令数量
// 	Qbuf     string // 查询缓存的长度（ 0 表示没有查询在等待）
// 	QbufFree string // 查询缓存的剩余空间（ 0 表示没有剩余空间）
// 	Obl      string // 输出缓存的长度
// 	Oll      string // 输出列表的长度（当输出缓存没有剩余空间时，回复被入队到这个队列里）
// 	Omem     string // 输出缓存的内存占用量
// 	Events   string // 文件描述符事件
// 	Cmd      string // 最近一次执行的命令
// }

//客户端 flag
var RedisClientFlags = map[string]string{
	"O": "客户端是 MONITOR 模式下的附属节点（slave）",
	"S": "客户端是一般模式下（normal）的附属节点",
	"M": "客户端是主节点（master）",
	"x": "客户端正在执行事务",
	"b": "客户端正在等待阻塞事件",
	"i": "客户端正在等待 VM I/O 操作（已废弃）",
	"d": "一个受监视（watched）的键已被修改， EXEC 命令将失败",
	"c": "在将回复完整地写出之后，关闭链接",
	"u": "客户端未被阻塞（unblocked）",
	"A": "尽可能快地关闭连接",
	"N": "未设置任何 flag",
}

// 文件描述符事件
var RedisClientEvents = map[string]string{
	"r": "客户端套接字（在事件 loop 中）是可读的（readable）",
	"w": "客户端套接字（在事件 loop 中）是可写的（writeable）",
}

//获取所有客户端
func GetAllRedisClient(c redis.Conn) (datas []*dao.RedisClientInfo) {
	manlog.Debug("[Execute redis command]: ", "CLIENT LIST")
	res, err := redis.String(c.Do("client", "list"))
	if err != nil {
		manlog.Error(err)
		return
	}
	manlog.Debug(res)
	reslist := strings.Split(res, "\n")
	for _, v := range reslist {
		manlog.Debug(v)
		v_list := strings.Split(v, " ")
		manlog.Debug(v_list)
		if len(v_list) > 17 {
			_, id_value := SplitEqualKV2Str(v_list[0])
			_, addr_value := SplitEqualKV2Str(v_list[1])
			_, fd_value := SplitEqualKV2Str(v_list[2])
			_, name_value := SplitEqualKV2Str(v_list[3])
			_, age_value := SplitEqualKV2Str(v_list[4])
			_, idle_value := SplitEqualKV2Str(v_list[5])
			_, flags_value := SplitEqualKV2Str(v_list[6])
			_, db_value := SplitEqualKV2Str(v_list[7])
			_, sub_value := SplitEqualKV2Str(v_list[8])
			_, psub_value := SplitEqualKV2Str(v_list[9])
			_, multi_value := SplitEqualKV2Str(v_list[10])
			_, qbuf_value := SplitEqualKV2Str(v_list[11])
			_, qbuf_free_value := SplitEqualKV2Str(v_list[12])
			_, obl_value := SplitEqualKV2Str(v_list[13])
			_, ool_value := SplitEqualKV2Str(v_list[14])
			_, omem_avlue := SplitEqualKV2Str(v_list[15])
			_, events_value := SplitEqualKV2Str(v_list[16])
			_, cmd_value := SplitEqualKV2Str(v_list[17])
			datas = append(datas, &dao.RedisClientInfo{
				ID:       id_value,
				Addr:     addr_value,
				Fd:       fd_value,
				Name:     name_value,
				Age:      age_value,
				Idle:     idle_value,
				Flags:    flags_value,
				Db:       db_value,
				Sub:      sub_value,
				Psub:     psub_value,
				Multi:    multi_value,
				Qbuf:     qbuf_value,
				QbufFree: qbuf_free_value,
				Obl:      obl_value,
				Oll:      ool_value,
				Omem:     omem_avlue,
				Events:   events_value,
				Cmd:      cmd_value,
			})
		}
	}

	manlog.Debug(datas)
	for _, v := range datas {
		manlog.Debug(*v)
	}

	return
}

func SplitEqualKV2Str(strs string) (k string, v string) {
	k = ""
	v = ""
	strs_list := strings.Split(strs, "=")
	if len(strs_list) > 2 {
		manlog.Error("存在多个等于符号")
		return
	}
	if len(strs_list) == 2 {
		k = strs_list[0]
		v = strs_list[1]
	}
	if len(strs_list) == 1 {
		k = strs_list[0]
	}
	return
}

func SplitEqualKV2Int(strs string) (k string, v int64) {
	k = ""
	v_str := ""
	v = 0
	strs_list := strings.Split(strs, "=")
	if len(strs_list) > 2 {
		manlog.Error("存在多个等于符号")
		return
	}
	if len(strs_list) == 2 {
		k = strs_list[0]
		v_str = strs_list[1]
		v = common.Str2Int64(v_str)
	}
	if len(strs_list) == 1 {
		k = strs_list[0]
	}
	return
}

func SplitEqualKV2Float(strs string) (k string, v float64) {
	k = ""
	v_str := ""
	v = 0
	strs_list := strings.Split(strs, "=")
	if len(strs_list) > 2 {
		manlog.Error("存在多个等于符号")
		return
	}
	if len(strs_list) == 2 {
		k = strs_list[0]
		v_str = strs_list[1]
		v = common.Str2Float64(v_str)
	}
	if len(strs_list) == 1 {
		k = strs_list[0]
	}
	return
}

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
//
//CLIENT KILL ip:port
//关闭地址为 ip:port 的客户端。
//ip:port 应该和 CLIENT LIST 命令输出的其中一行匹配。
//因为 Redis 使用单线程设计，所以当 Redis 正在执行命令的时候，不会有客户端被断开连接。
//
//

//关闭指定客户端
func KillRedisClient() {}

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
//
// CLIENT GETNAME
// 返回 CLIENT SETNAME 命令为连接设置的名字。
// 因为新创建的连接默认是没有名字的， 对于没有名字的连接， CLIENT GETNAME 返回空白回复。
//
//

//获取连接设置的名字
func GetRedisClientName() {}
