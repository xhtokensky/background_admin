package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"tokensky_bg_admin/conf"
	"tokensky_bg_admin/enums"
	"tokensky_bg_admin/models"
)

//身份审核
type TokenskyRealAuthController struct {
	BaseController
}

//Prepare 参考beego官方文档说明
func (c *TokenskyRealAuthController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	c.checkAuthor()

}

// DataGrid
func (c *TokenskyRealAuthController) DataGrid() {
	//直接反序化获取json格式的requestbody里的值
	var params models.TokenskyRealAuthQueryParam
	switch c.Ctx.Request.Method {
	case "POST":
		json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	case "GET":
		params.Sort = c.GetString("sort")
		params.Order = c.GetString("order")
		params.Limit, _ = c.GetInt64("limit")
		params.Offset, _ = c.GetInt64("offset")

		params.Phone = c.GetString("phone")
		params.StartTime, _ = c.GetInt64("startTime")
		params.EndTime, _ = c.GetInt64("endTime")
		params.Name = c.GetString("name")
		params.IdentityCard = c.GetString("identityCard")
		params.Status = c.GetString("status")
	}
	//获取数据列表和总数
	data, total := models.TokenskyRealAuthPageList(&params)
	//定义返回的数据结构
	mapp := map[string]interface{}{
		"rows":  data,
		"total": total,
	}
	c.jsonResult(enums.JRCodeSucc, "", mapp)
}

//身份审核
func (c *TokenskyRealAuthController) Auditing() {
	if c.Ctx.Request.Method == "POST" {
		var err error
		m := models.TokenskyRealAuth{}
		body := c.Ctx.Input.RequestBody
		if err := json.Unmarshal(body, &m); err != nil {
			c.jsonResult(enums.JRCodeFailed, "数据异常:", m.KeyId)
		}
		obj, err := models.TokenskyRealAuthOneById(m.KeyId)
		if err != nil {
			c.jsonResult(enums.JRCodeFailed, "获取记录失败", m.KeyId)
		}
		switch m.Status {
		case conf.TOKENSKY_REAL_AUTH_STATUS_PASSED:
		case conf.TOKENSKY_REAL_AUTH_STATUS_FAILED:
		default:
			c.jsonResult(enums.JRCodeFailed, "未知状态", m.KeyId)
		}
		switch obj.Status {
		case conf.TOKENSKY_REAL_AUTH_STATUS_PASSED, conf.TOKENSKY_REAL_AUTH_STATUS_FAILED:
			c.jsonResult(enums.JRCodeFailed, "重复审核", obj.KeyId)
		}
		obj.Status = m.Status
		o := orm.NewOrm()
		if _, err = o.Update(obj); err == nil {
			c.jsonResult(enums.JRCodeSucc, "审核成功", obj.KeyId)
		} else {
			c.jsonResult(enums.JRCodeFailed, "更新数据失败", obj.KeyId)
		}
	}
}
