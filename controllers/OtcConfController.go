package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"tokensky_bg_admin/enums"
	"tokensky_bg_admin/models"
)

//OtcConfController Otc配置
type OtcConfController struct {
	BaseController
}

//Prepare 参考beego官方文档说明
func (c *OtcConfController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	c.checkAuthor()
}

// DataGrid
func (c *OtcConfController) DataGrid() {
	//直接反序化获取json格式的requestbody里的值
	var params models.OtcConfQueryParam
	switch c.Ctx.Request.Method {
	case "POST":
		json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	case "GET":
		params.Sort = c.GetString("sort")
		params.Order = c.GetString("order")
		params.Limit, _ = c.GetInt64("limit")
		params.Offset, _ = c.GetInt64("offset")
	}

	//获取数据列表和总数
	data, total := models.OtcConfPageList(&params)
	//定义返回的数据结构
	mapp := map[string]interface{}{
		"rows":  data,
		"total": total,
	}
	c.jsonResult(enums.JRCodeSucc, "", mapp)
}

// 获取最新配置数据
func (c *OtcConfController) GetConf() {
	obj := models.OtcConfGetLastOne()
	mapp := map[string]interface{}{
		"obj": obj,
	}
	c.jsonResult(enums.JRCodeSucc, "", mapp)
}

//Edit 添加
func (c *OtcConfController) Edit() {
	if c.Ctx.Request.Method == "POST" {
		c.Save()
	}
}

//Save 添加、编辑页面 保存
func (c *OtcConfController) Save() {
	var err error
	m := models.OtcConf{}
	mBak := models.OtcConfBak{}
	body := c.Ctx.Input.RequestBody //接收raw body内容
	if err := json.Unmarshal(body, &m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "数据异常:", 0)
	}
	if err := json.Unmarshal(body, &mBak); err != nil {
		c.jsonResult(enums.JRCodeFailed, "数据异常:", 0)
	}
	if m.OrdersMinQuota > m.OrdersMaxQuota {
		c.jsonResult(enums.JRCodeFailed, "交易额度异常:", 0)
	}

	o := orm.NewOrm()

	////数据校验
	//if ok,msg := m.Check();!ok{
	//	c.jsonResult(enums.JRCodeFailed, msg, 0)
	//}
	m.UserId = c.BaseController.curUser.Id
	mBak.UserId = c.BaseController.curUser.Id
	mBak.Id = 0
	//修改本表
	m.Id = 1
	if _, err = o.Update(&m); err == nil {
		//新增副表
		o.Insert(&mBak)
		c.jsonResult(enums.JRCodeSucc, "编辑成功", 0)
	} else {
		c.jsonResult(enums.JRCodeFailed, "编辑失败", 0)
	}
}
