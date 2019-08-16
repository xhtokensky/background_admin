package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"tokensky_bg_admin/enums"
	"tokensky_bg_admin/models"
)

//理财分类
type FinancialCategoryController struct {
	BaseController
}

//Prepare 参考beego官方文档说明
func (c *FinancialCategoryController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	c.checkAuthor()
}

// DataGrid
func (c *FinancialCategoryController) DataGrid() {
	var params models.FinancialCategoryQueryParam
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
	}
	//获取数据列表和总数
	data, total := models.FinancialCategoryPageList(&params)
	//定义返回的数据结构
	mapp := map[string]interface{}{
		"rows":  data,
		"total": total,
	}
	c.jsonResult(enums.JRCodeSucc, "", mapp)
}

//编辑
func (c *FinancialCategoryController)Edit() {
	if c.Ctx.Request.Method == "POST" {
		var err error
		m := models.FinancialCategory{}
		body := c.Ctx.Input.RequestBody
		if err = json.Unmarshal(body, &m); err != nil {
			c.jsonResult(enums.JRCodeFailed, "数据异常:", m.Id)
		}
		m.AdminId = c.curUser.Id
		o := orm.NewOrm()
		switch m.Id {
		case 0:
			//新增
			if _, err = o.Insert(&m); err == nil {
				c.jsonResult(enums.JRCodeSucc, "添加成功", m.Id)
			} else {
				c.jsonResult(enums.JRCodeFailed, "添加失败", m.Id)
			}
		default:
			////编辑
			//obj := &models.FinancialCategory{Id:m.Id}
			//err =o.Read(obj)
			//m.CreateTime = obj.CreateTime
			//m.Symbol = obj.Symbol
			//if m.Avatar == ""{
			//	c.jsonResult(enums.JRCodeFailed, "图片不存在", m.Id)
			//}
			//if m.Symbol == ""{
			//	c.jsonResult(enums.JRCodeFailed, "货币不存在", m.Id)
			//}
			//if !models.TokenskyUserBalanceCoinIsFound(m.Symbol){
			//	c.jsonResult(enums.JRCodeFailed, "货币不存在", m.Id)
			//}
			//if err != nil{
			//	c.jsonResult(enums.JRCodeFailed, "编辑数据不存在", m.Id)
			//}
			//if _, err = o.Update(&m); err == nil {
			//	c.jsonResult(enums.JRCodeSucc, "编辑成功", m.Id)
			//} else {
			//	c.jsonResult(enums.JRCodeFailed, "编辑失败", m.Id)
			//}
		}
	}
	c.jsonResult(enums.JRCodeFailed, "请求错误", nil)
}

//删除
func (c *FinancialCategoryController)Delete() {
	if c.Ctx.Request.Method == "POST" {
		var id int
		var err error
		body := c.Ctx.Input.RequestBody
		mapp := make(map[string]int)
		err = json.Unmarshal(body,&mapp)
		if err != nil{
			c.jsonResult(enums.JRCodeFailed, "数据解析失败",0)
		}
		id = mapp["id"]
		if id <=0{
			c.jsonResult(enums.JRCodeFailed, "ID数据错误",0)
		}
		ok,msg := models.FinancialCategoryDelete(id)
		if ok{
			c.jsonResult(enums.JRCodeSucc, msg, nil)
		}else {
			c.jsonResult(enums.JRCodeFailed, msg, nil)
		}
	}
	c.jsonResult(enums.JRCodeFailed, "请求错误", nil)
}