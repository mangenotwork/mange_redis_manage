//
//	redis Databases相关的所有操作
//

package redis

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
	"github.com/mangenotwork/mange_redis_manage/common"
	_ "github.com/mangenotwork/mange_redis_manage/common/manlog"
)

//获取redis db 数量
func GetDatabasesCount(rc redis.Conn) int {
	fmt.Println("执行redis : ", "config get databases")
	res, err := redis.StringMap(rc.Do("config", "get", "databases"))
	if err != nil {
		fmt.Println("GET error", err.Error())
	}
	fmt.Println("获取redis db 数量 = ", res)
	return common.Str2Int(res["databases"])
}
