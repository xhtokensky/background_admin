package utils

import (
	"sync"
	"time"
)

//用户登陆错误次数校验模块
var (
	// map[ip]时间戳
	userDoLogRescource sync.Map
)

const (
	userDoLogMaxTime  int64 = 3600 //单位时间
	userDologMaxTimes int64 = 10   //单位时间的最大偶误登陆次数
	userDologMaxNum   int64 = 200
)

var (
	userDologMaxNum2 int64 = 0
)

//校验IP是有拥有访问登陆权限
func UserDolangCheckIsOk(ip string) bool {
	t, ok := userDoLogRescource.Load(ip)
	if ok {
		now := time.Now().Unix()
		tm := t.(int64)
		if tm > now+userDoLogMaxTime {
			return false
		}
	}
	return true
}

//增加错误次数
func UserDolangCheckAddErr(ip string) {
	t, ok := userDoLogRescource.Load(ip)
	now := time.Now().Unix()
	if ok {
		tm := t.(int64)
		if tm < now {
			userDoLogRescource.Store(ip, now)
		} else {
			userDoLogRescource.Store(ip, tm+userDoLogMaxTime/userDologMaxTimes)
		}
	} else {
		userDoLogRescource.Store(ip, now)
	}

	//每错误多少次清理下空间
	userDologMaxNum2++
	if userDologMaxNum2/userDologMaxNum == 0 {
		UserDolangCheckClear()
	}
}

//清理 可以定时器清理
func UserDolangCheckClear() {
	now := time.Now().Unix()
	delList := make([]string, 0)
	userDoLogRescource.Range(func(key, value interface{}) bool {
		k := key.(string)
		v := value.(int64)
		if v < now {
			delList = append(delList, k)
		}
		return true
	})
	for _, k := range delList {
		userDoLogRescource.Delete(k)
	}
}
