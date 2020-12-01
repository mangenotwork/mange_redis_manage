//
//	签名
//
package common

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

func NewMangeSign() string {
	h := md5.New()
	h.Write([]byte(fmt.Sprintf("ManGe-Redis-%d", time.Now().UnixNano()))) // 需要加密的字符串为 123456
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func Md5File(reader io.Reader) string {
	var buf = make([]byte, 4096)
	hashMd5 := md5.New()
	for {
		n, err := reader.Read(buf)
		if err == io.EOF && n == 0 {
			break
		}
		if err != nil && err != io.EOF {
			break
		}
		hashMd5.Write(buf[:n])
	}

	return fmt.Sprintf("%x", hashMd5.Sum(nil))
}

func Base64Encode(data string) string {
	return base64.URLEncoding.EncodeToString([]byte(data))
}

func Base64Decode(data string) string {
	b, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		return ""
	}
	return string(b)
}

// 下划线写法转为驼峰写法
func CamelName(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = strings.Title(name)
	return strings.Replace(name, " ", "", -1)
}

func SearchString(slice []string, s string) int {
	for i, v := range slice {
		if s == v {
			return i
		}
	}

	return -1
}

//任意类型转字符串
func ConvertString(inter interface{}, precs ...int) string {
	switch v := inter.(type) {
	case string:
		return v
	case float64:
		prec := 0
		if len(precs) > 0 {
			prec = precs[0]
		}
		return strconv.FormatFloat(v, 'f', prec, 64)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case int:
		return strconv.Itoa(v)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	default:
		return fmt.Sprintf("%v", inter)
	}
}

func SafeHtml(s string) string {
	r := strings.NewReplacer("<input", "&lt;input", "<a ", "&lt; a")
	return r.Replace(s)
}
