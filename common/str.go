//
//	字符串操作
//
package common

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

func Str2Int64(str string) int64 {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

func Str2Int(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return i
}

func Str2Float64(str string) float64 {
	i, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}
	return i
}

//[]uint8 转 string
func Uint82Str(bs []uint8) string {
	ba := []byte{}
	for _, b := range bs {
		ba = append(ba, byte(b))
	}
	return string(ba)
}

//返回一个32位md5加密后的字符串
func GetMD5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

//返回一个16位md5加密后的字符串
func Get16MD5Encode(data string) string {
	return GetMD5Encode(data)[8:24]
}

// 任何类型返回值字符串形式
func StringValue(i interface{}) string {
	var buf bytes.Buffer
	stringValue(reflect.ValueOf(i), 0, &buf)
	return buf.String()
}

// 任何类型返回值字符串形式的实现方法，私有
func stringValue(v reflect.Value, indent int, buf *bytes.Buffer) {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		buf.WriteString("{\n")

		for i := 0; i < v.Type().NumField(); i++ {
			ft := v.Type().Field(i)
			fv := v.Field(i)

			if ft.Name[0:1] == strings.ToLower(ft.Name[0:1]) {
				continue // ignore unexported fields
			}
			if (fv.Kind() == reflect.Ptr || fv.Kind() == reflect.Slice) && fv.IsNil() {
				continue // ignore unset fields
			}

			buf.WriteString(strings.Repeat(" ", indent+2))
			buf.WriteString(ft.Name + ": ")

			if tag := ft.Tag.Get("sensitive"); tag == "true" {
				buf.WriteString("<sensitive>")
			} else {
				stringValue(fv, indent+2, buf)
			}

			buf.WriteString(",\n")
		}

		buf.WriteString("\n" + strings.Repeat(" ", indent) + "}")
	case reflect.Slice:
		nl, id, id2 := "", "", ""
		if v.Len() > 3 {
			nl, id, id2 = "\n", strings.Repeat(" ", indent), strings.Repeat(" ", indent+2)
		}
		buf.WriteString("[" + nl)
		for i := 0; i < v.Len(); i++ {
			buf.WriteString(id2)
			stringValue(v.Index(i), indent+2, buf)

			if i < v.Len()-1 {
				buf.WriteString("," + nl)
			}
		}

		buf.WriteString(nl + id + "]")
	case reflect.Map:
		buf.WriteString("{\n")

		for i, k := range v.MapKeys() {
			buf.WriteString(strings.Repeat(" ", indent+2))
			buf.WriteString(k.String() + ": ")
			stringValue(v.MapIndex(k), indent+2, buf)

			if i < v.Len()-1 {
				buf.WriteString(",\n")
			}
		}

		buf.WriteString("\n" + strings.Repeat(" ", indent) + "}")
	default:
		format := "%v"
		switch v.Interface().(type) {
		case string:
			format = "%q"
		}
		fmt.Fprintf(buf, format, v.Interface())
	}
}

//字节大小转换
func humanSize(value float64) string {
	switch {
	case value > 1<<30:
		return fmt.Sprintf("%.2f GB", value/(1<<30))
	case value > 1<<20:
		return fmt.Sprintf("%.2f MB", value/(1<<20))
	case value > 1<<10:
		return fmt.Sprintf("%.2f kB", value/(1<<10))
	}
	return fmt.Sprintf("%.2f B", value)
}

/**
 * 对象转换为string
 * 支持类型：int,float64,string,bool(true:"1";false:"0")
 * 其他类型报错
 */
func ToString(obj interface{}) string {
	switch obj.(type) {
	case int:
		return strconv.Itoa(obj.(int))
	case float64:
		return strconv.FormatFloat(obj.(float64), 'f', -1, 64)
	case string:
		return obj.(string)
	case bool:
		if obj.(bool) {
			return "1"
		} else {
			return "0"
		}
	default:
		panic("ToString出错")
	}
}

/**
 * 对象转换为bool
 * 支持类型：int,float64,string,bool
 * 其他类型报错
 */
func ToBool(obj interface{}) bool {
	switch obj.(type) {
	case int:
		if obj.(int) == 0 {
			return false
		} else {
			return true
		}
	case float64:
		if obj.(float64) == 0 {
			return false
		} else {
			return true
		}
	case string:
		trues := map[string]int{"true": 1, "是": 1, "1": 1, "真": 1}
		if _, ok := trues[strings.ToLower(obj.(string))]; ok {
			return true
		} else {
			return false
		}
	case bool:
		return obj.(bool)
	default:
		panic("ToBool出错")
	}
}

/**
 * 对象转换为int
 * 支持类型：int,float64,string,bool(true:1;false:0)
 * 其他类型报错
 */
func ToInt(obj interface{}) int {
	switch obj.(type) {
	case int:
		return obj.(int)
	case float64:
		return int(obj.(float64))
	case string:
		ret, _ := strconv.Atoi(obj.(string))
		return ret
	case bool:
		if obj.(bool) {
			return 1
		} else {
			return 0
		}
	default:
		panic("ToInt出错")
	}
}

/**
 * 对象转换为float64
 * 支持类型：int,float64,string,bool(true:1;false:0)
 * 其他类型报错
 */
func ToFloat(obj interface{}) float64 {
	switch obj.(type) {
	case int:
		return float64(obj.(int))
	case float64:
		return obj.(float64)
	case string:
		ret, _ := strconv.ParseFloat(obj.(string), 64)
		return ret
	case bool:
		if obj.(bool) {
			return float64(1)
		} else {
			return float64(0)
		}
	default:
		panic("ToFloat出错")
	}
}

// ToString 将 []byte 转换为 string
func ByteToString(p []byte) string {
	return *(*string)(unsafe.Pointer(&p))
}

// ToBytes 将 string 转换为 []byte
func StringToBytes(str string) []byte {
	return *(*[]byte)(unsafe.Pointer(&str))
}

// IntToBool int 类型转换为 bool
// 0:false
// !0 : true
func IntToBool(i int) bool {
	if i == 0 {
		return false
	}
	return true
}

// FormatTime 将 Unix 时间戳, 转换为字符串
func FormatTime(t int64) string {
	return time.Unix(t, 0).Format("2006-01-02 03:04:05")
}
