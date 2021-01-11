//
//

package manlog

import (
	"fmt"
	"log"
	"runtime"
)

func init() {
	log.SetFlags(log.Ldate | log.Lmicroseconds)
}

const (
	ColorTpl = "\x1b[%sm%s\x1b[0m"
)

var (
	InfoTpl  = "\033[1;30;44m [Info]\033[0m"
	DebugTpl = "\033[1;30;42m [Debug]\033[0m"
	ErrorTpl = "\033[1;37;41m [Error]\033[0m"
)

func funcnametpl(funName string) string {
	return fmt.Sprintf("\033[1;34;1m FuncName=%s; \033[0m", funName) //蓝色字
}

func filetpl(file string, line int) string {
	return fmt.Sprintf("\033[1;35;1m File=%s ==> Line=%d; \033[0m", file, line) //绿色字
}

func funcnametplerr(funName string) string {
	return fmt.Sprintf("\033[1;31;1m FuncName=%s; \033[0m", funName) //蓝色字
}

func filetplerr(file string, line int) string {
	return fmt.Sprintf("\033[1;31;1m File=%s ==> Line=%d; \033[0m", file, line) //绿色字
}

func filetplinfo(file string, line int) string {
	return fmt.Sprintf("\033[1;32;1m File=%s ==> Line=%d; \033[0m", file, line) //绿色字
}

func Debug(data ...interface{}) {
	pc, file, line, _ := runtime.Caller(1)
	fun := runtime.FuncForPC(pc)
	funName := fun.Name()
	fmt.Print("\n")
	log.Printf("%s|%s|%s|\033[1;32;1m%v\033[0m", DebugTpl, funcnametpl(funName), filetpl(file, line), data)
}

func Panic(data ...interface{}) {
	pc, file, line, _ := runtime.Caller(1)
	fun := runtime.FuncForPC(pc)
	funName := fun.Name()
	fmt.Print("\n")
	log.Printf("%c[7;40;37m [Panic] FuncName=%s; File=%s; Line=%d; \tLog=%v %c[0m", 0x1B, funName, file, line, data, 0x1B)
	log.Panic()
}

func Error(data ...interface{}) {
	pc, file, line, _ := runtime.Caller(1)
	fun := runtime.FuncForPC(pc)
	funName := fun.Name()
	fmt.Print("\n")
	log.Printf("%s|%s|%s|\033[1;31;1m%v\033[0m", ErrorTpl, funcnametplerr(funName), filetplerr(file, line), data)
}

func Info(data ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	fmt.Print("\n")
	//log.Printf("%s|%s|\033[1;32;1m%v\033[0m", InfoTpl, filetplinfo(file, line), data)
	log.Println("%s|%s|\033[1;32;1m%v\033[0m", InfoTpl, filetplinfo(file, line), data)
}

// 字背景颜色范围:40----49
// 40:黑
// 41:深红
// 42:绿
// 43:黄色
// 44:蓝色
// 45:紫色
// 46:深绿
// 47:白色
//
// 字颜色:30-----------39
// 30:黑
// 31:红
// 32:绿
// 33:黄
// 34:蓝色
// 35:紫色
// 36:深绿
// 37:白色
//
// ascii 控制码
// \33[0m 关闭所有属性
// \33[1m 设置高亮度
// \33[4m 下划线
// \33[5m 闪烁
// \33[7m 反显
// \33[8m 消隐
// \33[30m -- \33[37m 设置前景色
// \33[40m -- \33[47m 设置背景色
// \33[nA 光标上移n行
// \33[nB 光标下移n行
// \33[nC 光标右移n行
// \33[nD 光标左移n行
// \33[y;xH设置光标位置
// \33[2J 清屏
// \33[K 清除从光标到行尾的内容
// \33[s 保存光标位置
// \33[u 恢复光标位置
// \33[?25l 隐藏光标
// \33[?25h 显示光标
//
// 代码             意义
//  -------------------------
//  0                 OFF
//  1                 高亮显示
//  4                 underline
//  5                 闪烁
//  7                 反白显示
//  8                 不可见
//
// 序列说 明
// \a ASCII响铃字符（也可以键入 \007）
// \d "Wed Sep 06"格式的日期
// \e ASCII转义字符（也可以键入 \033）
// \h 主机 名的第一部分（如 "mybox"）
// \H 主机 的全称（如 "mybox.mydomain.com"）
// \j 在此 shell中通过按 ^Z挂起的进程数
// \l 此 shell的终端设备名 （如 "ttyp4"）
// \n 换行 符
// \r 回车 符
// \s shell的名称（如 "bash"）
// \t 24小时制时间（如 "23:01:01"）
// \T 12小时制时间（如 "11:01:01"）
// \@ 带有 am/pm的 12小时制时间
// \u 用户 名
// \v bash的版本（如 2.04）
// \V Bash版本（包括补丁级别） ?/td>;
// \w 当前 工作目录（如 "/home/drobbins"）
// \W 当前 工作目录的“基名 (basename)”（如 "drobbins"）
// \! 当前 命令在历史缓冲区中的位置
// \# 命令 编号（只要您键入内容，它就会在每次提示时累加）
// \$ 如果 您不是超级用户 (root)，则插入一个 "$"；如果您是超级用户，则显示一个 "#"
// \xxx 插 入一个用三位数 xxx（用零代替未使用的数字，如 "/007"）表示的 ASCII 字符
// \\ 反斜 杠
// \[这个序列应该出现 在不移动光标的字符序列（如颜色转义序列）之前。它使 bash能够正确计算自动换行。
// \] 这个序列应该出现在非打印字符序列之后。
//
// 颜色的设置公式
// 颜色=\033[代码;前景;背景m
