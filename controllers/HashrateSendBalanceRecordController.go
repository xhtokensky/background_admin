package controllers

import (
	"encoding/json"
	"time"

	//"time"
	"tokensky_bg_admin/conf"
	"tokensky_bg_admin/enums"
	"tokensky_bg_admin/models"
)

type HashrateSendBalanceRecordController struct {
	BaseController
}

func (c *HashrateSendBalanceRecordController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	c.checkAuthor()
}

func (c *HashrateSendBalanceRecordController) DataGrid() {
	//直接反序化获取json格式的requestbody里的值
	var params models.HashrateSendBalanceRecordParam
	switch c.Ctx.Request.Method {
	case "POST":
		json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	case "GET":
		params.Sort = c.GetString("sort")
		params.Order = c.GetString("order")
		params.Limit, _ = c.GetInt64("limit")
		params.Offset, _ = c.GetInt64("offset")
		params.StartTime, _ = c.GetInt64("startTime")
		params.EndTime, _ = c.GetInt64("endTime")
		params.Status = c.GetString("status")
		params.CoinType = c.GetString("coinType")
	}
	//获取数据列表和总数
	data, total := models.HashrateSendBalanceRecordPageList(&params)
	//定义返回的数据结构
	mapp := map[string]interface{}{
		"rows":  data,
		"total": total,
	}
	c.jsonResult(enums.JRCodeSucc, "", mapp)
}

//发放资产
func (c *HashrateSendBalanceRecordController) SendBalcnce() {
	if conf.HASHRATE_SEND_SIGN {
		c.jsonResult(enums.JRCodeFailed, "收益发放中", 0)
	}

	//时间校验
	tm, _ := c.GetInt64("tm")
	//这里是因为前端是前端单位是毫秒
	tm = tm/1000
	now := time.Now()
	tz := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Unix()
	if tm <= 0 {
		c.jsonResult(enums.JRCodeFailed, "时间为0", 0)
	}
	if tm > tz {
		c.jsonResult(enums.JRCodeFailed, "只能发送今日往前的", 0)
	}
	t1 := time.Unix(tm, 0)
	hour := t1.Hour()
	if hour >=11 && hour <= 15{
		//防止挖矿数据正在处理中导致的一些异常，也可以加锁处理
		c.jsonResult(enums.JRCodeFailed, "当前时间无法发放", 0)
	}
	conf.HASHRATE_SEND_SIGN = true
	defer func() {
		conf.HASHRATE_SEND_SIGN = false
	}()
	//拉取收益
	HashrateSendBalanceRecords,err := models.HashrateOrderSendBalanceGetProfitRecord(tm, conf.Hashrate_Send_Balance_Allow_Coin_Type)
	if err != nil{
		c.jsonResult(enums.JRCodeFailed, "拉取收益异常:"+err.Error() , 0)
	}
	if len(HashrateSendBalanceRecords) > 0 {
		//创建收益
		models.HashrateOrderSendBalanceCreateProfitTb(tm, HashrateSendBalanceRecords)
		//奖励发放
		models.HashrateOrderSendBalanceProfit()
	} else {
		c.jsonResult(enums.JRCodeFailed, "无需要发放的收益", 0)
	}
	c.jsonResult(enums.JRCodeSucc, "奖励发放完成", 0)
}
