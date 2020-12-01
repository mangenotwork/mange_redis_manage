// 定时任务
//

package service

import (
	"time"

	"github.com/mangenotwork/mange_redis_manage/common/cache"
	"github.com/mangenotwork/mange_redis_manage/common/manlog"
	"github.com/mangenotwork/mange_redis_manage/dao"
)

func TimingTask() {
	//HealthAnalyzer()
	//DelCollectRedisInfo()
	//BackupCache()
}

//删除采集的redis服务器信息表里的数据
func DelCollectRedisInfo() {
	//获取配置， 删除几天前的
	//多久删除一次

	//执行删除
	//获取当前时间
	now_time := time.Now().Unix()

	//删除10天前的数据
	del_time := now_time - 60*60*24*10

	go new(dao.DaoRedisClients).DelTimeData(del_time)
	go new(dao.DaoRedisCluster).DelTimeData(del_time)
	go new(dao.DaoRedisCPU).DelTimeData(del_time)
	go new(dao.DaoRedisServerInfos).DelTimeData(del_time)
	go new(dao.DaoRedisKeyspace).DelTimeData(del_time)
	go new(dao.DaoRedisMemory).DelTimeData(del_time)
	go new(dao.DaoRedisPersistence).DelTimeData(del_time)
	go new(dao.DaoRedisReplication).DelTimeData(del_time)
	go new(dao.DaoRedisStats).DelTimeData(del_time)

}

//将服务本身的缓存备份
func BackupCache() {
	//获取配置， 多久备份一次

	//执行备份
	cache.Save2File()
}

//redis 健康分析
func HealthAnalyzer() {
	//获取配置， 多久执行一次

	//执行健康分析
	//获取连接的db，连接
	conndata, err := new(dao.DaoRedisInfo).GetAllConn()
	if err != nil {
		manlog.Error("获取所有连接数据出错， err = ", err)
	}

	for _, v := range conndata {
		conn, cid, err := new(RedisConn).RredisConn(v)
		manlog.Debug(conn, cid, err)
		if err == nil {
			RedisHealthAnalyzerRun(conn, cid, v)
		}
	}
}

func EmptyCollectRedisInfo() {
	go new(dao.DaoRedisClients).Empty()
	go new(dao.DaoRedisCluster).Empty()
	go new(dao.DaoRedisCPU).Empty()
	go new(dao.DaoRedisServerInfos).Empty()
	go new(dao.DaoRedisKeyspace).Empty()
	go new(dao.DaoRedisMemory).Empty()
	go new(dao.DaoRedisPersistence).Empty()
	go new(dao.DaoRedisReplication).Empty()
	go new(dao.DaoRedisStats).Empty()
}
