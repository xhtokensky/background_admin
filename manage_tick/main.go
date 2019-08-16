package main

import (
	"fmt"
	"github.com/astaxie/beego/toolbox"
	_ "tokensky_bg_admin/manage_tick/sysinit"
	"tokensky_bg_admin/manage_tick/tick"
)

//定时任务

/*

符号	含义	示例
*	   表示任何时刻
,	   表示分割	如第三段里：2,4，表示 2 点和 4 点执行
－	   表示一个段	如第三端里： 1-5，就表示 1 到 5 点
/n	   表示每个n的单位执行一次	如第三段里，1, 就表示每隔 1 个小时执行一次命令。也可以写成1-23/1

示例	详细含义
0/30 * * * * *	    每 30 秒 执行
0 43 21 * * *	    21:43 执行
0 0 17 * * 1	    每周一的 17:00 执行
0 0,10 17 * * 0,2,3	每周日,周二,周三的 17:00和 17:10 执行
0 0 21 * * 1-6	    周一到周六 21:00 执行
0 0/10 * * *	    每隔 10 分 执行
*/

const (
	//用户地址簿维护时间间隔
	TBI_SERVER_ADDRESS_TICK string = "300"
)

func main() {
	defer func() { select {} }()
	//算计资源发放
	sendBalance := toolbox.NewTask("sendBalance", "0 30 12 * * *", tick.TickHashrateOrderSendBalance)
	toolbox.AddTask("sendBalance", sendBalance)
	//用户地址簿维护
	tickTokenskyUserAddressUp := toolbox.NewTask("tickTokenskyUserAddressUp", "0/"+TBI_SERVER_ADDRESS_TICK+" * * * * *", tick.TickTokenskyUserAddressUp)
	toolbox.AddTask("tickTokenskyUserAddressUp", tickTokenskyUserAddressUp)
	//理财
	financialTick := toolbox.NewTask("financialTick", "0 0 0 * * *", tick.FinancialTick)
	toolbox.AddTask("financialTick", financialTick)
	//行情
	coinGlobalSoider := toolbox.NewTask("coinGlobalSoider", "0 0 0/1 * * *", tick.CoinGlobalSoider)
	toolbox.AddTask("coinGlobalSoider", coinGlobalSoider)
	//借贷风控系统
	borrowWindControlSystem := toolbox.NewTask("borrowWindControlSystem", "0 0/10 * * *", tick.BorrowWindControlSystem)
	toolbox.AddTask("borrowWindControlSystem", borrowWindControlSystem)

	//开始
	toolbox.StartTask()
	fmt.Println("定时器服务开启")

	//开启自动执行操作
	tick.CoinGlobalSoider()
}
