package sysinit

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"tokensky_bg_admin/manage_balance/controllers"
)




func InitLogs() {
	controllers.FileLogs = logs.NewLogger(10000)
	level := beego.AppConfig.String("logs::level")
	controllers.FileLogs.SetLogger(logs.AdapterMultiFile, `{"filename":"logs/balance.log",
		"separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"],
		"level":`+level+`,
		"daily":true,
		"maxdays":360}`)
	controllers.FileLogs.Async()                   //异步
	controllers.FileLogs.EnableFuncCallDepth(true) //显示文件名和行号
}