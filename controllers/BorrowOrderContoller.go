package controllers

import (
	"encoding/json"
	"tokensky_bg_admin/enums"
	"tokensky_bg_admin/models"
)

//
type BorrowOrderContoller struct {
	BaseController
}

//Prepare 参考beego官方文档说明
func (c *BorrowOrderContoller) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	c.checkAuthor()
}

// DataGrid
func (c *BorrowOrderContoller) DataGrid() {
	var params models.BorrowOrderQueryParam
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
		params.Name = c.GetString("name")
		params.Status = c.GetString("status")
		params.Symbol = c.GetString("symbol")
	}
	//获取数据列表和总数
	data, total := models.BorrowOrderPageList(&params)
	//定义返回的数据结构
	mapp := map[string]interface{}{
		"rows":  data,
		"total": total,
	}
	c.jsonResult(enums.JRCodeSucc, "", mapp)
}

//最高质押
func (c *BorrowOrderContoller) GetMaxPledge() {

	var maxPledge float64
	switch c.Ctx.Request.Method {
	case "POST":
		mapp := make(map[string]int)
		err := json.Unmarshal(c.Ctx.Input.RequestBody, &mapp)
		if err != nil{
			c.jsonResult(enums.JRCodeFailed, "解析异常", nil)
		}
		maxPledge = models.BorrowOrderGetMaxPledge(mapp["id"])
	case "GET":
		id, err := c.GetInt("id")
		if err !=nil{
			c.jsonResult(enums.JRCodeFailed, "没有传id字段", nil)
		}
		maxPledge = models.BorrowOrderGetMaxPledge(id)
	}
	data := map[string]float64{
		"maxPledge": maxPledge,
	}
	c.jsonResult(enums.JRCodeSucc, "", data)
}
