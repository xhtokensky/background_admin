package sysinit

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

//Beego定时任务

func init() {

}
