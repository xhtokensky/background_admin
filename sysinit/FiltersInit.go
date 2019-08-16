package sysinit

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

/*中间件*/

/*
BeforeStatic 静态地址之前
BeforeRouter 寻找路由之前
BeforeExec 找到路由之后，开始执行相应的 Controller 之前
AfterExec 执行完 Controller 逻辑之后执行的过滤器
FinishRouter 执行完逻辑之后执行的过滤器
*/

//设置cookie 为获取session做准备
func setCookieFunc(ctx *context.Context) {
	token := ctx.Request.Header.Get("token")
	if token != "" {
		ctx.SetCookie("token", token)
	}
}

//跨域
func setAllowCrossFunc(ctx *context.Context) {
	ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")                           //允许访问源
	ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "*")                          //允许post访问
	ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization") //header的类型
	ctx.ResponseWriter.Header().Set("Access-Control-Max-Age", "1728000")
	ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	ctx.ResponseWriter.Header().Set("content-type", "application/json") //返回数据格式是json
}

func init() {
	//静态地址之前
	//setCookis := setCookieFunc
	//beego.InsertFilter("*",beego.BeforeStatic,setCookis)
	//开始执行相应的 Controller 之前
	setAllowCross := setAllowCrossFunc
	beego.InsertFilter("*", beego.BeforeRouter, setAllowCross)
}
