//
//	redis key相关的所有操作
//

package redis

import (
	"fmt"
	"strings"

	"github.com/garyburd/redigo/redis"
	"github.com/mangenotwork/mange_redis_manage/common/manlog"
)

//获取所有的key
func GetALLKeys(c redis.Conn, matchvalue string) (ksyList map[string]int) {
	//初始化拆分值
	match_split := matchvalue

	//matchvalue :匹配值，没有则匹配所有 *
	if matchvalue == "" {
		matchvalue = "*"
	} else {
		matchvalue = fmt.Sprintf("*%s*", matchvalue)
	}
	//youbiao :初始游标为0
	youbiao := "0"

	ksyList = make(map[string]int)
	ksyList, youbiao = addGetKey(ksyList, matchvalue, match_split, c, youbiao)

	//当游标等于0的时候停止获取key
	//线性获取，一直循环获取key,直到游标为0
	if youbiao != "0" {
		for {
			ksyList, youbiao = addGetKey(ksyList, matchvalue, match_split, c, youbiao)
			if youbiao == "0" {
				break
			}
		}
	}

	fmt.Println("ksyList= ", ksyList)
	return
}

//addGetKey 内部方法
//针对分组的key进行分组合并处理
func addGetKey(ksyList map[string]int, matchvalue string, match_split string, conn redis.Conn, youbiao string) (map[string]int, string) {
	//count_number :一次10000
	count_number := "10000"
	res, err := redis.Values(conn.Do("scan", youbiao, "MATCH", matchvalue, "COUNT", count_number))
	manlog.Debug("执行redis : ", "scan", youbiao, "MATCH", matchvalue, "COUNT", count_number)
	if err != nil {
		manlog.Error("GET error", err.Error())
	}

	//获取	matchvalue 含有多少:
	cfnumber := strings.Count(matchvalue, ":")

	//获取新的游标
	newyoubiao := string(res[0].([]byte))
	allkey := res[1]
	allkey_data := allkey.([]interface{})
	for _, v := range allkey_data {
		key_data := string(v.([]byte))
		//manlog.Error("key_data = ", key_data)
		//manlog.Error("matchvalue = ", matchvalue)
		//没有:的key 则不集合
		if strings.Count(key_data, ":") == cfnumber || key_data == matchvalue {
			ksyList[key_data] = 0
			continue
		}

		//有:需要集合
		key_data_new, _ := FenGeYinghaoOne(key_data, match_split)
		ksyList[key_data_new] = ksyList[key_data_new] + 1
	}

	return ksyList, newyoubiao
}

//对查询出来的key进行拆分，集合，分组处理
func FenGeYinghaoOne(str string, match_split string) (string, int) {
	likekey := ""
	if match_split != "" {
		likekey = fmt.Sprintf("%s", match_split)
	}
	//fmt.Println("整理key的分组: ", str, likekey)
	str = strings.Replace(str, likekey, "", 1)
	fg := strings.Split(str, ":")
	//fmt.Println("str = ", str)
	//fmt.Println("fg = ", fg)
	if len(fg) > 0 {
		//fmt.Println(fg[0], len(fg))
		//fmt.Println(fmt.Sprintf("%s%s", likekey, fg[0]), len(fg), "\n\n")
		return fmt.Sprintf("%s%s", likekey, fg[0]), len(fg)
	}
	return "", len(fg)
}

func SearchKeys(c redis.Conn, matchvalue string) (ksyList map[string]int) {
	ksyList = make(map[string]int)
	//matchvalue :匹配值，没有则返回空
	if matchvalue == "" {
		return
	} else {
		matchvalue = fmt.Sprintf("*%s*", matchvalue)
	}

	//youbiao :初始游标为0
	youbiao := "0"

	ksyList = make(map[string]int)
	ksyList, youbiao = addSearchKey(ksyList, matchvalue, c, youbiao)

	//当游标等于0的时候停止获取key
	//线性获取，一直循环获取key,直到游标为0
	if youbiao != "0" {
		for {
			ksyList, youbiao = addSearchKey(ksyList, matchvalue, c, youbiao)
			if youbiao == "0" {
				break
			}
		}
	}

	fmt.Println("ksyList= ", ksyList)
	return
}

//addGetKey 内部方法获取key
func addSearchKey(ksyList map[string]int, matchvalue string, conn redis.Conn, youbiao string) (map[string]int, string) {
	//count_number :一次10000
	count_number := "10000"
	res, err := redis.Values(conn.Do("scan", youbiao, "MATCH", matchvalue, "COUNT", count_number))
	manlog.Debug("执行redis : ", "scan", youbiao, "MATCH", matchvalue, "COUNT", count_number)
	if err != nil {
		manlog.Error("GET error", err.Error())
	}

	//获取新的游标
	newyoubiao := string(res[0].([]byte))
	allkey := res[1]
	allkey_data := allkey.([]interface{})
	for _, v := range allkey_data {
		key_data := string(v.([]byte))
		ksyList[key_data] = 0
	}

	return ksyList, newyoubiao
}

//获取所有key name
// 返回切片
func GetAllKeyName(c redis.Conn) ([]interface{}, int) {

	all_key := make([]interface{}, 0)

	keydatas, youbiao := getallkey(c, "0")
	all_key = append(all_key, keydatas...)

	if youbiao != "0" {
		for {
			keydatas, youbiao = getallkey(c, youbiao)
			all_key = append(all_key, keydatas...)
			if youbiao == "0" {
				break
			}
		}
	}
	//fmt.Println(all_key)
	return all_key, len(all_key)
}

func getallkey(conn redis.Conn, youbiao string) ([]interface{}, string) {
	count_number := "10000"
	res, err := redis.Values(conn.Do("scan", youbiao, "MATCH", "*", "COUNT", count_number))
	manlog.Debug("执行redis : ", "scan", youbiao, "MATCH", "*", "COUNT", count_number)
	if err != nil {
		manlog.Error("GET error", err.Error())
	}
	return res[1].([]interface{}), string(res[0].([]byte))
}

//获取key的信息
func GetKeyInfo(rc redis.Conn, key string) {

}

//GetKeyType 获取key的类型
func GetKeyType(rc redis.Conn, keyname string) string {
	fmt.Println("执行redis : ", "TYPE", keyname)
	res, err := redis.String(rc.Do("TYPE", keyname))
	if err != nil {
		fmt.Println("GET error", err.Error())
	}
	fmt.Println(res)
	return res
}

//GetKeyTTL 获取key的过期时间
func GetKeyTTL(rc redis.Conn, keyname string) int64 {
	fmt.Println("执行redis : ", "TTL", keyname)
	res, err := redis.Int64(rc.Do("TTL", keyname))
	if err != nil {
		fmt.Println("GET error", err.Error())
	}
	fmt.Println(res)
	return res
}

//EXISTSKey 检查给定 key 是否存在。
func EXISTSKey(rc redis.Conn, keyname string) bool {
	fmt.Println("[Execute redis command]: ", "EXISTS", keyname)
	datas, err := redis.String(rc.Do("DUMP", keyname))
	if err != nil || datas == "0" {
		fmt.Println("GET error", err.Error())
		return false
	}
	return true
}

//修改key名称
func RenameKey(rc redis.Conn, keyname, newname string) bool {
	fmt.Println("[Execute redis command]: ", "RENAME", keyname, newname)
	_, err := rc.Do("RENAME", keyname, newname)
	if err != nil {
		fmt.Println("GET error", err.Error())
		return false
	}
	return true
}

//更新key ttl
func UpdateKeyTTL(rc redis.Conn, keyname string, ttlvalue int64) bool {
	fmt.Println("[Execute redis command]: ", "EXPIRE", keyname, ttlvalue)
	_, err := rc.Do("EXPIRE", keyname, ttlvalue)
	if err != nil {
		fmt.Println("GET error", err.Error())
		return false
	}
	return true
}

//指定key多久过期 接收的是unix时间戳
func EXPIREATKey(rc redis.Conn, keyname string, date int64) bool {
	fmt.Println("[Execute redis command]: ", "EXPIREAT", keyname, date)
	_, err := rc.Do("EXPIREAT", keyname, date)
	if err != nil {
		fmt.Println("GET error", err.Error())
		return false
	}
	return true
}

//删除key
func DELKey(rc redis.Conn, keyname string) bool {
	fmt.Println("[Execute redis command]: ", "DEL", keyname)
	_, err := rc.Do("DEL", keyname)
	if err != nil {
		fmt.Println("GET error", err.Error())
		return false
	}
	return true
}
