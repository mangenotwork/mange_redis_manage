// 消息队列生产者服务
//

package service

import (
	"sync"
	"time"

	"github.com/mangenotwork/mange_redis_manage/common/manlog"
	manredis "github.com/mangenotwork/mange_redis_manage/common/redis"
	"github.com/mangenotwork/mange_redis_manage/dao"
	"github.com/mangenotwork/mange_redis_manage/repository"
)

var redisinfoch = make(chan interface{}, 10)
var redisch = make(chan interface{}, 10)
var wg = new(sync.WaitGroup)

func QueueProducerStart() {

	go func() {
		for {
			select {
			//消费方
			case conn := <-redisinfoch:
				conndata := conn.(*models.RedisInfoDB)
				manlog.Debug(*conndata)
				RedisSrtverInfo2DB(conndata)
			case conn := <-redisch:
				conndata := conn.(*models.RedisInfoDB)
				manlog.Debug(*conndata)
			}
		}
	}()

	//生产方
	for {
		// go ADDRedisInfoMessage()
		// go Add()
		// wg.Wait()
		ADDRedisInfoMessage()
	}

}

//添加获取有效连接的redis服务器连接信息
func ADDRedisInfoMessage() {
	//wg.Add(1)
	//defer wg.Done()
	//获取所有有效连接，并去重
	conndata, err := new(dao.DaoRedisInfo).GetAllConn()
	if err != nil {
		manlog.Error("获取所有连接数据出错， err = ", err)
	}

	for _, v := range conndata {
		redisinfoch <- v
	}
	//这里取设置的时间，默认为10秒
	time.Sleep(10 * time.Second)
}

func Add() {
	wg.Add(1)
	defer wg.Done()
	conndata, err := new(dao.DaoRedisInfo).GetAllConn()
	if err != nil {
		manlog.Error("获取所有连接数据出错， err = ", err)
	}

	for _, v := range conndata {
		redisch <- v
	}
	time.Sleep(10 * time.Second)
}

func RedisSrtverInfo2DB(conndata *models.RedisInfoDB) {
	conn, cid, err := new(RedisConn).RredisConn(conndata)
	if err == nil {
		redisinfos := manredis.GetRedisServersInfos(conn, cid)
		new(RedisConn).RedisServserInfo2DB(redisinfos)
	}

}
