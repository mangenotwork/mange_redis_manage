//
//	redis 监控服务
//
package service

type RedisMonitoringService interface {
	Monitoring() //redis 服务监控
}

type RedisMonitoring struct {
}

func (this *RedisMonitoring) Monitoring() {}
