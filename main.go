package main

import (
	// "fmt"
	// "net/http"

	"github.com/gin-gonic/gin"
	"github.com/mangenotwork/mange_redis_manage/install"
	"github.com/mangenotwork/mange_redis_manage/routers"
	"github.com/mangenotwork/mange_redis_manage/service"
)

func init() {
	//检查是否安装
	is_install := install.IsInstall()
	//is_install := false

	//如果已经安装做如下启动
	if is_install {
		//开启队列生成者
		go service.QueueProducerStart()
		//开启队列消费则
		//go service.ExpenserStrat()
		//开启定时任务
		//go service.TimingTask()
	}
}

func main() {
	gin.SetMode(gin.DebugMode)
	s := router.Routers()
	port := "8334"
	s.Run(":" + port)
}
