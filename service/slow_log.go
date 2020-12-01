// 定时任务
//

package service

import (
	"github.com/mangenotwork/mange_redis_manage/common/manlog"
	"github.com/mangenotwork/mange_redis_manage/common/redis"
	"github.com/mangenotwork/mange_redis_manage/structs"
)

//在Redis中，关于慢查询有两个设置--慢查询最大超时时间和慢查询最大日志数。
// CONFIG  SET  slowlog-log-slower-than  num
// 设置超过多少微妙的查询为慢查询，并且将这些慢查询加入到日志文件中，num的单位为毫秒，windows下redis的默认慢查询时10000微妙即10毫秒。
// CONFIG  SET  slowlog-max-len  num
// 设置日志的最大数量，num无单位值，windows下redis默认慢查询日志的记录数量为128条。

//慢日志值的解析
// 1)    1)   "5939"
//  2)   "1600311723"
//  3)   "16784"
//  4)      1)    "SSCAN"
//   2)    "search:index"
//   3)    "77629"
//   4)    "COUNT"
//   5)    "10000"
//  5)   "192.168.0.101:61819"
//  6)   ""
//1 : 日志的唯一标识符
//2 : 命令执行时系统的时间戳
//3 : 命令执行的时长，以微妙来计算
//4 : 命令
//5 : 客户端

//redis 慢日志功能对外提供接口
type RedisSlowLogService interface {
	Get(user *structs.UserParameter, redisId int64) //获取慢日志
}

//redis的操作
type RedisSlowLog struct {
}

//慢日志
func (this *RedisSlowLog) Get(user *structs.UserParameter, redisId int64) {

	//获取连接
	rc, rcid, err := new(RedisConn).DBconn(user, redisId, 0)
	if err != nil {
		manlog.Error(err)
		//return false, err
	}

	redis.GetSlowlog(rc, rcid)
}
