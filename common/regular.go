// 包含了正则的常用方法

package common

import (
	_ "fmt"
	"regexp"
)

// td:=`<td>(.*?)</td>`
// tdreg := regexp.MustCompile(td)
// tdList := tdreg.FindAllStringSubmatch(rest,-1)

func FindAllstrlist(regstr, rest string) [][]string {
	reg := regexp.MustCompile(regstr)
	List := reg.FindAllStringSubmatch(rest, -1)
	// for _, v := range List {
	// 	fmt.Println("r = ", v)
	// }
	return List
}
