package sysinit



func init() {
	//初始化数据库
	InitDatabase()
	//初始化日志
	InitLogs()
}
