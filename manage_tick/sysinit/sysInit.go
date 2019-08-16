package sysinit

import "tokensky_bg_admin/utils"

func init() {
	//初始化数据库
	InitDatabase()
	//初始化日志
	utils.InitLogs()
}
