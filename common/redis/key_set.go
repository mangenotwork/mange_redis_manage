//
//	redis Set相关的所有操作
//

package redis

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
	_ "github.com/mangenotwork/mange_redis_manage/common/manlog"
)

// SMEMBERS key
// 返回集合 key 中的所有成员。
// 获取Set value 返回集合 key 中的所有成员。
func SetSMEMBERS(rc redis.Conn, keyname string) []interface{} {
	fmt.Println("执行redis : ", "SMEMBERS", keyname)
	res, err := redis.Values(rc.Do("SMEMBERS", keyname))
	if err != nil {
		fmt.Println("GET error", err.Error())
	}
	fmt.Println(res)
	return res
}

//新创建Set  将一个或多个 member 元素加入到集合 key 当中，已经存在于集合的 member 元素将被忽略。
func SetSADD(rc redis.Conn, keyname string, values []interface{}) error {
	args := redis.Args{}.Add(keyname)
	for _, value := range values {
		args = args.Add(value)
	}
	fmt.Println("执行redis : ", "SADD", args)
	res, err := rc.Do("SADD", args...)
	if err != nil {
		fmt.Println("GET error", err.Error())
		//key已经存在
		//WRONGTYPE Operation against a key holding the wrong kind of value
		return err
	}
	fmt.Println(res)
	return nil
}

// SCARD key
// 返回集合 key 的基数(集合中元素的数量)。
func SetSCARD(rc redis.Conn, keyname string) (err error) {
	fmt.Println("执行redis : ", "SCARD", keyname)
	res, err := redis.Int64(rc.Do("SCARD", keyname))
	if err != nil {
		fmt.Println("GET error", err.Error())
		return
	}
	fmt.Println(res)
	return
}

// SDIFF key [key ...]
// 返回一个集合的全部成员，该集合是所有给定集合之间的差集。
// 不存在的 key 被视为空集。
func SetSDIFF(rc redis.Conn, keys []string) (res []interface{}, err error) {
	args := redis.Args{}
	for _, key := range keys {
		args = args.Add(key)
	}
	fmt.Println("执行redis : ", "SDIFF", args)
	res, err = redis.Values(rc.Do("SDIFF", args))
	if err != nil {
		fmt.Println("GET error", err.Error())
		return
	}
	fmt.Println(res)
	return
}

// SDIFFSTORE destination key [key ...]
// 这个命令的作用和 SDIFF 类似，但它将结果保存到 destination 集合，而不是简单地返回结果集。
// 如果 destination 集合已经存在，则将其覆盖。
// destination 可以是 key 本身。
func SetSDIFFSTORE(rc redis.Conn, keyname string, keys []string) (res []interface{}, err error) {
	args := redis.Args{}.Add(keyname)
	for _, key := range keys {
		args = args.Add(key)
	}
	fmt.Println("执行redis : ", "SDIFFSTORE", args)
	res, err = redis.Values(rc.Do("SDIFFSTORE", args))
	if err != nil {
		fmt.Println("GET error", err.Error())
		return
	}
	fmt.Println(res)
	return
}

// SINTER key [key ...]
// 返回一个集合的全部成员，该集合是所有给定集合的交集。
// 不存在的 key 被视为空集。
func SetSINTER(rc redis.Conn, keys []string) (res []interface{}, err error) {
	args := redis.Args{}
	for _, key := range keys {
		args = args.Add(key)
	}
	fmt.Println("执行redis : ", "SINTER", args)
	res, err = redis.Values(rc.Do("SINTER", args))
	if err != nil {
		fmt.Println("GET error", err.Error())
		return
	}
	fmt.Println(res)
	return
}

// SINTERSTORE destination key [key ...]
// 这个命令类似于 SINTER 命令，但它将结果保存到 destination 集合，而不是简单地返回结果集。
// 如果 destination 集合已经存在，则将其覆盖。
// destination 可以是 key 本身。
func SetSINTERSTORE(rc redis.Conn, keyname string, keys []string) (res []interface{}, err error) {
	args := redis.Args{}.Add(keyname)
	for _, key := range keys {
		args = args.Add(key)
	}
	fmt.Println("执行redis : ", "SINTERSTORE", args)
	res, err = redis.Values(rc.Do("SINTERSTORE", args))
	if err != nil {
		fmt.Println("GET error", err.Error())
		return
	}
	fmt.Println(res)
	return
}

// SISMEMBER key member
// 判断 member 元素是否集合 key 的成员。
// 返回值:
// 如果 member 元素是集合的成员，返回 1 。
// 如果 member 元素不是集合的成员，或 key 不存在，返回 0 。
func SetSISMEMBER(rc redis.Conn, keyname string, value interface{}) (resdata bool, err error) {
	fmt.Println("执行redis : ", "SISMEMBER", keyname, value)
	resdata = false
	res, err := redis.Int64(rc.Do("SISMEMBER", keyname, value))
	if err != nil {
		fmt.Println("GET error", err.Error())
		return
	}
	if res == 1 {
		resdata = true
		return
	}
	fmt.Println(res)
	return
}

// SMOVE source destination member
// 将 member 元素从 source 集合移动到 destination 集合。
// SMOVE 是原子性操作。
// 如果 source 集合不存在或不包含指定的 member 元素，则 SMOVE 命令不执行任何操作，仅返回 0 。否则，
// member 元素从 source 集合中被移除，并添加到 destination 集合中去。
// 当 destination 集合已经包含 member 元素时， SMOVE 命令只是简单地将 source 集合中的 member 元素删除。
// 当 source 或 destination 不是集合类型时，返回一个错误。
// 返回值: 成功移除，返回 1 。失败0
func SetSMOVE(rc redis.Conn, keyname, destination string, member interface{}) (resdata bool, err error) {
	fmt.Println("执行redis : ", "SMOVE", keyname, destination, member)
	resdata = false
	res, err := redis.Int64(rc.Do("SMOVE", keyname, destination, member))
	if err != nil {
		fmt.Println("GET error", err.Error())
		return
	}
	if res == 1 {
		resdata = true
		return
	}
	fmt.Println(res)
	return
}

// SPOP key
// 移除并返回集合中的一个随机元素。
func SetSPOP(rc redis.Conn, keyname string) (res string, err error) {
	fmt.Println("执行redis : ", "SPOP", keyname)
	res, err = redis.String(rc.Do("SPOP", keyname))
	if err != nil {
		fmt.Println("GET error", err.Error())
		return
	}
	fmt.Println(res)
	return
}

// SRANDMEMBER key [count]
// 如果命令执行时，只提供了 key 参数，那么返回集合中的一个随机元素。
// 如果 count 为正数，且小于集合基数，那么命令返回一个包含 count 个元素的数组，数组中的元素各不相同。如果 count 大于等于集合基数，那么返回整个集合。
// 如果 count 为负数，那么命令返回一个数组，数组中的元素可能会重复出现多次，而数组的长度为 count 的绝对值。
func SetSRANDMEMBER(rc redis.Conn, keyname string, count int64) (res []interface{}, err error) {
	fmt.Println("执行redis : ", "SRANDMEMBER", keyname, count)
	res, err = redis.Values(rc.Do("SRANDMEMBER", keyname, count))
	if err != nil {
		fmt.Println("GET error", err.Error())
		return
	}
	fmt.Println(res)
	return
}

// SREM key member [member ...]
// 移除集合 key 中的一个或多个 member 元素，不存在的 member 元素会被忽略。
func SetSREM(rc redis.Conn, keyname string, member []interface{}) (err error) {
	args := redis.Args{}.Add(keyname)
	for _, v := range member {
		args = args.Add(v)
	}
	fmt.Println("执行redis : ", "SREM", args)
	res, err := rc.Do("SREM", args)
	if err != nil {
		fmt.Println("GET error", err.Error())
		return
	}
	fmt.Println(res)
	return
}

// SUNION key [key ...]
// 返回一个集合的全部成员，该集合是所有给定集合的并集。
func SetSUNION(rc redis.Conn, keys []string) (res []interface{}, err error) {
	args := redis.Args{}
	for _, v := range keys {
		args = args.Add(v)
	}
	fmt.Println("执行redis : ", "SUNION", args)
	res, err = redis.Values(rc.Do("SUNION", args))
	if err != nil {
		fmt.Println("GET error", err.Error())
		return
	}
	fmt.Println(res)
	return
}

// SUNIONSTORE destination key [key ...]
// 这个命令类似于 SUNION 命令，但它将结果保存到 destination 集合，而不是简单地返回结果集。
func SetSUNIONSTORE(rc redis.Conn, keyname string, keys []string) (res []interface{}, err error) {
	args := redis.Args{}.Add(keyname)
	for _, v := range keys {
		args = args.Add(v)
	}
	fmt.Println("执行redis : ", "SUNIONSTORE", args)
	res, err = redis.Values(rc.Do("SUNIONSTORE", args))
	if err != nil {
		fmt.Println("GET error", err.Error())
		return
	}
	fmt.Println(res)
	return
}

//搜索值  SSCAN key cursor [MATCH pattern] [COUNT count]
