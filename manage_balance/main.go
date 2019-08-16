package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"net/http"
	"os"
	"tokensky_bg_admin/manage_balance/controllers"
	_ "tokensky_bg_admin/manage_balance/sysinit"
)

func main() {
	defer func() { select {} }()
	http.HandleFunc("/balance/one", controllers.SetBalanceOne)
	http.HandleFunc("/balance/multi", controllers.SetBalanceMulti)
	port := beego.AppConfig.String("httpport")
	fmt.Println("资产服务已启动")
	for {
		err := http.ListenAndServe(":"+port, nil) //监听
		if err != nil {
			controllers.FileLogs.Error("资产服务器错误:" + err.Error())
			os.Exit(0)
		}
	}
}
