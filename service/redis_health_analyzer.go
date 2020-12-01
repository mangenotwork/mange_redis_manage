// redis的健康分析
// 主要采用具体到每个key和value的全面分析
// 注意: 该功能IO占用大，会占用服务资源， 使用单线，不使用并发
// TODO : 加入慢日志

package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
	"unicode/utf8"

	"github.com/garyburd/redigo/redis"
	"github.com/mangenotwork/mange_redis_manage/common/manlog"
	manredis "github.com/mangenotwork/mange_redis_manage/common/redis"
	"github.com/mangenotwork/mange_redis_manage/repository"
	"github.com/mangenotwork/mange_redis_manage/structs"
)

type RedisHealthData struct {
	RedisName  string `json:"redis_name"`
	RedisHost  string `json:"redis_host"`
	HealthDate string `json:"time"`

	//redis 宿主服务器内存
	SysMemory string `json:"sys_memory"`

	//redis 已用内存大小
	USEDMemory string `json:"used_memory"`

	//内存使用率
	RedisMemoryUSED      float64 `json:"redis_memory_used"`
	RedisMemoryUSEDHuman string  `json:"redis_memory_used_human"`

	//ping redis 服务的结果
	RedisPingList []*structs.RedisPingData `json:"redis_ping_list"`
	PingAvgTime   string                   `json:"ping_avg"`

	//命中率
	RedisHitRate      float64 `json:"redis_hit_rate"`
	RedisHitRateHuman string  `json:"redis_hit_rate_human"`

	//设置了ttl在所有key的占比
	RedisTTLRatio      float64 `json:"redis_ttl_ratio"`
	RedisTTLRatioHuman string  `json:"redis_ttl_ratio_human"`

	//key数量，各种类型key数量 ....
	CountNumber *SUMs `json:"count_data"`

	//每个db的分析数据, 只有使用过的db才做分析
	RedisDB []*RedisDBHealthData `json:"redis_db_list"`
}

//redis key相关的总计数量
type SUMs struct {
	AllkeyNumber    int64 `json:"all_key_number"`    //所有key数量
	TTLKeyNumber    int64 `json:"ttl_key_number"`    //设置了ttl的key数量
	StringKeyNumber int64 `json:"string_key_number"` //类型为string的key数量
	HashKeyNumber   int64 `json:"hash_key_number"`   //类型为hash的key数量
	ListKeyNumber   int64 `json:"list_key_number"`   //类型为list的key数量
	SetKeyNumber    int64 `json:"set_key_number"`    //类型为set的key数量
	ZSetKeyNumber   int64 `json:"zset_key_number"`   //类型为zset的key数量
	ALLSize         int64 `json:"all_key_size"`      //统计的Key的value占字节的大小
}

//redis 各db的相关健康分析数据
type RedisDBHealthData struct {
	DBName string
	DBID   int64

	//redis db key相关的总计数量
	Number *SUMs

	//统计到key命名太长的key
	KeyNameMaxList []*KeyNameTooLong

	//统计到key value 太大的key
	KeyValueMaxList []*KeyValueToolLoang
}

type KeyNameTooLong struct {
	KeyName string
	KeyType string
	KeyTTL  int64
	Count   int64
}

type KeyValueToolLoang struct {
	KeyName string
	KeyType string
	KeyTTL  int64
	Size    int64
}

//redis健康分析的具体实现
func RedisHealthAnalyzerRun(rc redis.Conn, rcid string, redisconninfo *models.RedisInfoDB) {

	redisHealth := &RedisHealthData{}

	redisHealth.RedisName = redisconninfo.ConnName
	redisHealth.RedisHost = fmt.Sprintf("%s:%d", redisconninfo.ConnHost, redisconninfo.ConnPort)
	redisHealth.HealthDate = time.Now().Format("2006-01-02 15:04:05")

	//分析指标:
	//1. 容量占比
	//redis宿主的内存大小
	// 如果设置了 maxmemory 就用 maxmemory，否则用total_system_memory
	//"maxmemory": 0,
	//	"maxmemory_human": "0B",
	//"total_system_memory": 33513107456,
	//	"total_system_memory_human": "31.21G",
	//公式:   used_memory /  total_system_memory
	var redis_memory_used float64
	//获取服务信息
	redis_info_data := manredis.GetRedisServersInfos(rc, rcid)
	if redis_info_data.Memory.Maxmemory == 0 {
		redis_memory_used = float64(redis_info_data.Memory.UsedMemory) / float64(redis_info_data.Memory.TotalSystemMemory)
	} else {
		redis_memory_used = float64(redis_info_data.Memory.UsedMemory) / float64(redis_info_data.Memory.Maxmemory)
	}
	manlog.Debug("redis_memory_used = ", redis_memory_used)
	redisHealth.RedisMemoryUSED = redis_memory_used
	redisHealth.SysMemory = redis_info_data.Memory.TotalSystemMemoryHuman
	redisHealth.USEDMemory = redis_info_data.Memory.UsedMemoryHuman

	//2. ping速度
	//ping 10次无统计时间
	var sum_tiem float64 = 0
	redis_pings := make([]*structs.RedisPingData, 0)
	for i := 0; i < 10; i++ {
		a_time := time.Now().UnixNano()
		isok := manredis.Ping(rc)
		t := float64(time.Now().UnixNano()-a_time) / 1000000
		manlog.Debug("t, isok = ", t, isok)
		sum_tiem = sum_tiem + t
		redis_pings = append(redis_pings, &structs.RedisPingData{
			Number: int64(i + 1),
			Time:   t,
			IsOK:   isok,
		})
	}
	redisHealth.RedisPingList = redis_pings
	avgping := sum_tiem / 10
	redisHealth.PingAvgTime = fmt.Sprintf("%g ms", avgping)

	//3. 命中率
	//keyspace_hits：命中的次数
	//keyspace_misses：没有命中的次数
	//缓存命中率 = keyspace_hits / (keyspace_hits + keyspace_misses)
	var hits_misses float64 = 0
	if redis_info_data.Stats.KeyspaceHits != 0 && redis_info_data.Stats.KeyspaceMisses != 0 {
		hits_misses = float64(redis_info_data.Stats.KeyspaceHits) / float64(redis_info_data.Stats.KeyspaceHits+redis_info_data.Stats.KeyspaceMisses)
	}
	manlog.Debug(hits_misses)
	redisHealth.RedisHitRate = hits_misses

	//4. 各个db的单项值
	//遍历每个db获取所有的keyname
	//先检查keyname名称是否太长，将太长的记录
	//查询值，值的大小判断，记录太大的key
	//总计db大小,key数量
	//统计有设置ttl 与没有设置ttl的key
	//key的设置过期的占比
	dbsDatas := make([]*RedisDBHealthData, 0)
	sums := &SUMs{}

	for _, v := range redis_info_data.Keyspace {
		dbdata := RedisDBHealthAnalyzer(rc, v.DBID, sums)
		manlog.Debug(*dbdata, *dbdata.Number)
		dbsDatas = append(dbsDatas, dbdata)
	}

	var key_ttl_ratio float64
	key_ttl_ratio = float64(sums.TTLKeyNumber) / float64(sums.AllkeyNumber)
	redisHealth.RedisTTLRatio = key_ttl_ratio

	manlog.Debug(*sums)
	manlog.Debug(dbsDatas)
	manlog.Debug(key_ttl_ratio)

	redisHealth.CountNumber = sums
	redisHealth.RedisDB = dbsDatas

	b, err := json.MarshalIndent(redisHealth, "", "\t")
	if err != nil {
		fmt.Println("Umarshal failed:", err)
	}
	healthdata := string(b)
	manlog.Debug(healthdata)

	//最后生成md文档到temp中
	RedisHealthMarkdown(redisHealth)
}

//redis db 的健康分析
func RedisDBHealthAnalyzer(rc redis.Conn, dbid int64, sums *SUMs) (dbdata *RedisDBHealthData) {
	var err error
	toolangkeyname := make([]*KeyNameTooLong, 0)
	toolangkeyvalue := make([]*KeyValueToolLoang, 0)

	//切换db id
	rc, err = manredis.SelectDB(rc, dbid)
	if err != nil {
		manlog.Error("切换db失败， err = ", err)
	}
	dbdata = &RedisDBHealthData{}
	number := &SUMs{}

	dbdata.DBName = fmt.Sprintf("DB%d", dbid)
	dbdata.DBID = dbid

	//获取每个db获取所有的keyname
	allkey, key_number := manredis.GetAllKeyName(rc)
	sums.AllkeyNumber = sums.AllkeyNumber + int64(key_number)
	number.AllkeyNumber = int64(key_number)
	for _, v := range allkey {
		keyname := string(v.([]byte))

		//获取key类型
		key_type := manredis.GetKeyType(rc, keyname)
		switch key_type {
		case "string":
			sums.StringKeyNumber++
			number.StringKeyNumber++
		case "hash":
			sums.HashKeyNumber++
			number.HashKeyNumber++
		case "list":
			sums.ListKeyNumber++
			number.ListKeyNumber++
		case "set":
			sums.SetKeyNumber++
			number.SetKeyNumber++
		case "zset":
			sums.ZSetKeyNumber++
			number.ZSetKeyNumber++
		}

		//获取key ttl
		key_ttl := manredis.GetKeyTTL(rc, keyname)
		if key_ttl > 1 {
			sums.TTLKeyNumber++
			number.TTLKeyNumber++
		}

		//key命名超过50个字符
		if utf8.RuneCountInString(keyname) > 50 {
			manlog.Debug("key命名超过50个字符 = ", keyname)
			toolangkeyname = append(toolangkeyname, &KeyNameTooLong{
				KeyName: keyname,
				KeyType: key_type,
				KeyTTL:  key_ttl,
				Count:   int64(utf8.RuneCountInString(keyname)),
			})
		}

		//获取key的值
		var factory = new(RedisTypeOperateFactory)
		keyObj := factory.Operate(key_type)
		key_size := keyObj.ValueSize(rc, keyname)
		//如果key的值大于10kb就统计     1024*10
		if key_size > 1024*10 {
			manlog.Debug("key的值大于10kb , key = ", keyname)
			toolangkeyvalue = append(toolangkeyvalue, &KeyValueToolLoang{
				KeyName: keyname,
				KeyType: key_type,
				KeyTTL:  key_ttl,
				Size:    key_size,
			})
		}
		number.ALLSize = number.ALLSize + key_size
		sums.ALLSize = sums.ALLSize + key_size
	}

	dbdata.Number = number
	dbdata.KeyNameMaxList = toolangkeyname
	dbdata.KeyValueMaxList = toolangkeyvalue
	return
}

func RedisHealthMarkdown(data *RedisHealthData) {

	//头
	title_template := "# Redis： %s 的健康分析\n"
	title_s := fmt.Sprintf(title_template, data.RedisName)

	//文件名
	file_name := fmt.Sprintf("%s的健康分析-%s.md", data.RedisName, data.HealthDate)

	//基础
	redisinfo_template := `
## 基础信息
- host: %s
- 分析时间：%s
`
	redisinfo_s := fmt.Sprintf(redisinfo_template, data.RedisHost, data.HealthDate)

	//容量
	redismemory_template := `
## 容量分析
- 宿主系统内存 : %s
- 已用内存: %s
- 内存使用占比: %v %%
`
	redismemory_s := fmt.Sprintf(redismemory_template, data.SysMemory, data.USEDMemory, data.RedisMemoryUSED)

	//ping
	redisping_template := `
## Ping
> 平均值 : %s
%s
`
	redispingdata_template := `
%d. 是否成功: %v; 用时: %v ms;
`
	redispingdata_s := ""
	for _, v := range data.RedisPingList {
		redispingdata_s = redispingdata_s + fmt.Sprintf(redispingdata_template, v.Number, v.IsOK, v.Time)
	}
	redisping_s := fmt.Sprintf(redisping_template, data.PingAvgTime, redispingdata_s)

	//命中率
	redishitrate_template := `
## 命中率
> %v %%
`
	redishitrate_s := fmt.Sprintf(redishitrate_template, data.RedisHitRate)

	//ttl设置率
	redisttlratio_template := `
## 设置了ttl的key所在总量的占比
> %v %%
`
	redisttlratio_s := fmt.Sprintf(redisttlratio_template, data.RedisTTLRatio)

	//统计的总量
	rediscountnumber_template := `
## 统计总数
- key总数 : %d
- 设置了ttl key总数 : %d
- string类型的key总数 : %d
- hash类型的key总数 : %d
- list类型的key总数 : %d
- set类型的key总数 : %d
- zset类型key总数 : %d
- key大小 : %d
`

	rediscountnumber_s := fmt.Sprintf(rediscountnumber_template, data.CountNumber.AllkeyNumber, data.CountNumber.TTLKeyNumber,
		data.CountNumber.StringKeyNumber, data.CountNumber.HashKeyNumber, data.CountNumber.ListKeyNumber, data.CountNumber.SetKeyNumber,
		data.CountNumber.ZSetKeyNumber, data.CountNumber.ALLSize)

	//DB具体数据
	redisdbtitle_template := `
## 各个DB的数据	
`
	//db详细数据
	redisdbinfo_template := `
### DB Name : %s

#### 统计总数
%s

#### 统计到key命名太长的key
%s

#### 统计到key value 太大的key
%s

`
	redisdbcountnumber_template := `
- key总数 : %d
- 设置了ttl key总数 : %d
- string类型的key总数 : %d
- hash类型的key总数 : %d
- list类型的key总数 : %d
- set类型的key总数 : %d
- zset类型key总数 : %d
- key大小 : %d

`

	redisdbnamemaxlist_template := `
> KeyName : %s
- KeyType : %s
- KeyTTL  : %d
- Count   : %d

`

	redisdbvaluemaxlist_template := `
> KeyName  : %s
- KeyType : %s
- KeyTTL  : %d
- Size    : %d

`
	redisdb_s := redisdbtitle_template
	for _, v := range data.RedisDB {

		redisdbsum_s := fmt.Sprintf(redisdbcountnumber_template, v.Number.AllkeyNumber, v.Number.TTLKeyNumber,
			v.Number.StringKeyNumber, v.Number.HashKeyNumber, v.Number.ListKeyNumber, v.Number.SetKeyNumber,
			v.Number.ZSetKeyNumber, v.Number.ALLSize)

		redismaxnamelist_s := ""
		for _, n := range v.KeyNameMaxList {
			redismaxnamelist_s = redismaxnamelist_s + fmt.Sprintf(redisdbnamemaxlist_template, n.KeyName,
				n.KeyType, n.KeyTTL, n.Count)
		}

		redismaxvaluelist_s := ""
		for _, m := range v.KeyValueMaxList {
			redismaxvaluelist_s = redismaxvaluelist_s + fmt.Sprintf(redisdbvaluemaxlist_template, m.KeyName,
				m.KeyType, m.KeyTTL, m.Size)
		}

		redisdbinfo_s := fmt.Sprintf(redisdbinfo_template, v.DBName, redisdbsum_s, redismaxnamelist_s, redismaxvaluelist_s)
		redisdb_s = redisdb_s + redisdbinfo_s
	}

	s := title_s + redisinfo_s + redismemory_s + redisping_s + redishitrate_s + redisttlratio_s + rediscountnumber_s + redisdb_s
	manlog.Debug(s)
	manlog.Debug(file_name)

	//当前路径
	pwd, _ := os.Getwd()
	file_path := pwd + "/report/" + file_name
	isPath := checkFileIsExist(file_path)
	if !isPath {
		var wireteSb = []byte(s)
		err := ioutil.WriteFile(file_path, wireteSb, 0666) //写入文件(字节数组)
		manlog.Debug(err)
	}
}

// 判断文件是否存在  存在返回 true 不存在返回false
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
