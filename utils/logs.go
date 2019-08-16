package utils

import (
	"strings"
	"tokensky_bg_admin/conf"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

// consoleLogs开发模式下日志
var consoleLogs *logs.BeeLogger

// fileLogs 生产环境下日志
var fileLogs *logs.BeeLogger

//运行方式
var runmode string

func InitLogs() {
	consoleLogs = logs.NewLogger(1)
	consoleLogs.SetLogger(logs.AdapterConsole)
	consoleLogs.Async() //异步
	fileLogs = logs.NewLogger(10000)
	level := beego.AppConfig.String("logs::level")
	fileLogs.SetLogger(logs.AdapterMultiFile, `{"filename":"logs/rms.log",
		"separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"],
		"level":`+level+`,
		"daily":true,
		"maxdays":360}`)
	fileLogs.Async()                   //异步
	fileLogs.EnableFuncCallDepth(true) //显示文件名和行号
	runmode = strings.TrimSpace(strings.ToLower(beego.AppConfig.String("runmode")))
	if runmode == "" {
		runmode = "dev"
	}
}
func LogEmergency(v interface{}) {
	log("emergency", v)
}
func LogAlert(v interface{}) {
	log("alert", v)
}
func LogCritical(v interface{}) {
	//邮件
	EmailAdd(conf.EMAIL_ERROR_LEVEL_CRITICAL, "异常", "", v)
	log("critical", v)
}
func LogError(v interface{}) {
	//邮件
	EmailAdd(conf.EMAIL_ERROR_LEVEL_ERROR, "异常", "", v)
	log("error", v)
}
func LogWarning(v interface{}) {
	log("warning", v)
}
func LogNotice(v interface{}) {
	log("notice", v)
}
func LogInfo(v interface{}) {
	log("info", v)
}
func LogDebug(v interface{}) {
	log("debug", v)
}

func LogTrace(v interface{}) {
	log("trace", v)
}

//Log 输出日志
func log(level, v interface{}) {
	format := "%s"
	if level == "" {
		level = "debug"
	}

	/* 写入 */
	if runmode == "dev" || runmode== "prod" {
		switch level {
		case "emergency":
			fileLogs.Emergency(format, v)
		case "alert":
			fileLogs.Alert(format, v)
		case "critical":
			fileLogs.Critical(format, v)
		case "error":
			fileLogs.Error(format, v)
		case "warning":
			fileLogs.Warning(format, v)
		case "notice":
			fileLogs.Notice(format, v)
		case "info":
			fileLogs.Info(format, v)
		case "debug":
			fileLogs.Debug(format, v)
		case "trace":
			fileLogs.Trace(format, v)
		default:
			fileLogs.Debug(format, v)
		}
	}

	/* 屏幕打印 */
	switch level {
	case "emergency":
		consoleLogs.Emergency(format, v)
	case "alert":
		consoleLogs.Alert(format, v)
	case "critical":
		consoleLogs.Critical(format, v)
	case "error":
		consoleLogs.Error(format, v)
	case "warning":
		consoleLogs.Warning(format, v)
	case "notice":
		consoleLogs.Notice(format, v)
	case "info":
		consoleLogs.Info(format, v)
	case "debug":
		consoleLogs.Debug(format, v)
	case "trace":
		consoleLogs.Trace(format, v)
	default:
		consoleLogs.Debug(format, v)
	}
}
