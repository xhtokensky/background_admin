package controllers

import (
	"encoding/json"
	"strconv"
	"tokensky_bg_admin/enums"
	"tokensky_bg_admin/models"
)

//强平订单表
type BorrowLimitingController struct {
	BaseController
}

//Prepare 参考beego官方文档说明
func (c *BorrowLimitingController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	c.checkAuthor()
}


// DataGrid
func (c *BorrowLimitingController) DataGrid() {
	var params models.BorrowLimitingQueryParam
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

		params.OrderId = c.GetString("orderId")
		params.Name = c.GetString("name")
	}
	//获取数据列表和总数
	data, total := models.BorrowLimitingPageList(&params)
	//定义返回的数据结构
	mapp := map[string]interface{}{
		"rows":  data,
		"total": total,
	}
	c.jsonResult(enums.JRCodeSucc, "", mapp)
}

//售卖 BorrowLimitingsSell
func (c *BorrowLimitingController)Sell(){
	if c.Ctx.Request.Method == "POST"{
		mapp := make(map[string]string)
		err:= json.Unmarshal(c.Ctx.Input.RequestBody, &mapp)
		if err != nil{
			c.jsonResult(enums.JRCodeFailed, "解析错误", "")
		}
		orderId,ok := mapp["orderId"]
		if !ok{
			c.jsonResult(enums.JRCodeFailed, "缺少orderId", "")
		}
		totalPriceStr,ok := mapp["totalPrice"]
		if !ok{
			c.jsonResult(enums.JRCodeFailed, "缺少售卖总额 totalPrice", "")
		}
		totalPrice,err := strconv.ParseFloat(totalPriceStr,64)
		if err != nil{
			c.jsonResult(enums.JRCodeFailed, "售卖总额解析异常 err:"+err.Error(), "")
		}
		err = models.BorrowLimitingsSell(orderId,totalPrice,c.curUser.Id)
		if err != nil{
			c.jsonResult(enums.JRCodeSucc, "售卖完成", "")
		}else {
			c.jsonResult(enums.JRCodeSucc, "售卖异常:"+err.Error(), "")
		}
	}
	c.jsonResult(enums.JRCodeFailed, "请求错误", "")
}