package controllers

import (
	"encoding/json"
	"tokensky_bg_admin/enums"
	"tokensky_bg_admin/models"
)

//算力合约订单表
type HashrateOrderController struct {
	BaseController
}

//Prepare 参考beego官方文档说明
func (c *HashrateOrderController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	c.checkAuthor()
}

// DataGrid
func (c *HashrateOrderController) DataGrid() {
	//直接反序化获取json格式的requestbody里的值
	var params models.HashrateOrderQueryParam
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
		params.Status = c.GetString("status")
		params.Phone = c.GetString("phone")
	}
	//获取数据列表和总数
	data, total := models.HashrateOrderPageList(&params)
	//定义返回的数据结构
	mapp := map[string]interface{}{
		"rows":  data,
		"total": total,
	}
	c.jsonResult(enums.JRCodeSucc, "", mapp)
}
