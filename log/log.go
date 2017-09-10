package log

import (
	"fmt"
	"time"
)

var debug bool = true

type LogLevel string

const (
	DEBUG   = "D"
	INFO    = "I"
	WARNING = "W"
	ERROR   = "E"
)

type LogType string

const (
	LOGIN     = "登录"
	PASSENGER = "获取用户信息"
	OTHER     = "其它"
)

func MyLoginLogI(format string, a ...interface{}) {
	MyLog(INFO, LOGIN, format, a...)
}
func MyLoginLogW(format string, a ...interface{}) {
	MyLog(WARNING, LOGIN, format, a...)
}
func MyLoginLogE(format string, a ...interface{}) {
	MyLog(ERROR, LOGIN, format, a...)
}

func MyLogDebug(format string, a ...interface{}) {
	MyLog(DEBUG, OTHER, format, a...)
}

func MyLog(l LogLevel, t LogType, format string, a ...interface{}) {
	if l == DEBUG && debug == false {
		return
	}
	format = string("[%v][%s][%s]:") + format + string("\n")
	arg := make([]interface{}, 3, len(a)+3)
	arg[0] = time.Now().Format("2006-01-02 15:04:05.999")
	arg[1] = t
	arg[2] = l
	if len(a) > 0 {
		arg = append(arg, a...)
	}
	fmt.Printf(format, arg...)
}
