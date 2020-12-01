package dao

//Redis 客户端信息
type RedisClientInfo struct {
	ID       string `json:"id"`        //
	Addr     string `json:"addr"`      // 客户端的地址和端口
	Fd       string `json:"fd"`        // 套接字所使用的文件描述符
	Name     string `json:"name"`      //
	Age      string `json:"age"`       // 以秒计算的已连接时长
	Idle     string `json:"idle"`      // 以秒计算的空闲时长
	Flags    string `json:"flags"`     // 客户端 flag
	Db       string `json:"db"`        // 该客户端正在使用的数据库 ID
	Sub      string `json:"sub"`       // 已订阅频道的数量
	Psub     string `json:"psub"`      // 已订阅模式的数量
	Multi    string `json:"multi"`     // 在事务中被执行的命令数量
	Qbuf     string `json:"qbuf"`      // 查询缓存的长度（ 0 表示没有查询在等待）
	QbufFree string `json:"qbuf_free"` // 查询缓存的剩余空间（ 0 表示没有剩余空间）
	Obl      string `json:"obl"`       // 输出缓存的长度
	Oll      string `json:"oll"`       // 输出列表的长度（当输出缓存没有剩余空间时，回复被入队到这个队列里）
	Omem     string `json:"omem"`      // 输出缓存的内存占用量
	Events   string `json:"events"`    // 文件描述符事件
	Cmd      string `json:"cmd"`       // 最近一次执行的命令
}
