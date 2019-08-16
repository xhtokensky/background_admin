package controllers

import (
	"encoding/json"
	"tokensky_bg_admin/enums"
	"tokensky_bg_admin/models"
)

//用户资产表

type TokenskyUserBalanceController struct {
	BaseController
}

//Prepare 参考beego官方文档说明
func (c *TokenskyUserBalanceController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	c.checkAuthor()

}

//获取某一用户的所有资产
func (c *TokenskyUserBalanceController) GetBalances() {
	var params models.TokenskyUserBalanceRecordParam
	switch c.Ctx.Request.Method {
	case "POST":
		json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	case "GET":
		params.Sort = c.GetString("sort")
		params.Order = c.GetString("order")
		params.Limit, _ = c.GetInt64("limit")
		params.Offset, _ = c.GetInt64("offset")
		params.Uid, _ = c.GetInt("uid")
	}
	if params.Uid > 0 {
		data, total := models.GetTokenskyUserBalancesByUid(params.Uid)
		//定义返回的数据结构
		mapp := map[string]interface{}{
			"rows":  data,
			"total": total,
		}
		c.jsonResult(enums.JRCodeSucc, "", mapp)
	} else {
		c.jsonResult(enums.JRCodeFailed, "用户uid不能为0", 0)
	}
}
