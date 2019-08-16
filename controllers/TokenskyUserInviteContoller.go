package controllers

import (
	"tokensky_bg_admin/enums"
	"tokensky_bg_admin/models"
)

//用户邀请表
type TokenskyUserInviteContoller struct {
	BaseController
}

//Prepare 参考beego官方文档说明
func (c *TokenskyUserInviteContoller) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	c.checkAuthor()
}

//次数统计
func (c *TokenskyUserInviteContoller) FormToAmount() {
	//
	data := models.TokenskyUserInviteFormToAmount()
	c.jsonResult(enums.JRCodeSucc, "ok", data)
}