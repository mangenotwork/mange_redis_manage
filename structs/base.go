//
//	基础结构体,主要包含了输出
//
package structs

type ResponseJson struct {
	Code      int64       `json:"code"`
	Mag       string      `json:"mag"`
	Date      interface{} `json:"data"`
	TimeStamp int64       `json:"timeStamp"`
}

type UserParameter struct {
	UserID    int64  `json:"user_id"`
	Username  string `json:"user_name"`
	UserGropy int64  `json:"user_group"`
}

type ResponseKV struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}


