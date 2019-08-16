package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"strconv"
	"tokensky_bg_admin/enums"
	"tokensky_bg_admin/models"
	"tokensky_bg_admin/utils"
)

//用户
type TokenskyUserController struct {
	BaseController
}

//Prepare 参考beego官方文档说明
func (c *TokenskyUserController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	c.checkAuthor()

}

// DataGrid
func (c *TokenskyUserController) DataGrid() {
	//直接反序化获取json格式的requestbody里的值
	var params models.TokenskyUserQueryParam
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
		params.NickName = c.GetString("nickName")
		params.UserId, _ = c.GetInt("userId")
	}
	//获取数据列表和总数
	data, total := models.TokenskyUserPageList(&params)
	//定义返回的数据结构
	mapp := map[string]interface{}{
		"rows":  data,
		"total": total,
	}
	c.jsonResult(enums.JRCodeSucc, "", mapp)
}

//等级修改
func (c *TokenskyUserController) SetLevel() {
	mapp := make(map[string]int)
	body := c.Ctx.Input.RequestBody //接收raw body内容
	if err := json.Unmarshal(body, &mapp); err != nil {
		c.jsonResult(enums.JRCodeFailed, "数据异常:", 0)
	}
	id := mapp["userId"]
	level := mapp["level"]
	obj, err := models.TokenskyUserOne(id)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "无此用户", id)
	}
	obj.Level = level
	o := orm.NewOrm()
	if _, err := o.Update(obj); err == nil {
		c.jsonResult(enums.JRCodeSucc, "修改成功", id)
	} else {
		c.jsonResult(enums.JRCodeSucc, "修改失败", id)
	}
}

//设置用户是否拥有邀请权限
func (c *TokenskyUserController)SetInvitation(){
	mapp := make(map[string]int)
	body := c.Ctx.Input.RequestBody //接收raw body内容
	if err := json.Unmarshal(body, &mapp); err != nil {
		c.jsonResult(enums.JRCodeFailed, "数据异常:", 0)
	}
	id := mapp["userId"]
	invitation := mapp["invitation"]
	obj, err := models.TokenskyUserOne(id)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "无此用户", id)
	}
	switch invitation {
	case 0,1:
	default:
		c.jsonResult(enums.JRCodeFailed, "无此状态", id)
	}

	obj.Invitation = invitation
	o := orm.NewOrm()
	if _, err := o.Update(obj); err == nil {
		c.jsonResult(enums.JRCodeSucc, "修改成功", id)
	} else {
		c.jsonResult(enums.JRCodeSucc, "修改失败", id)
	}
}

//获取用户url连接
func (c *TokenskyUserController)GetAddr(){
	mapp := make(map[string]int)
	body := c.Ctx.Input.RequestBody //接收raw body内容
	if err := json.Unmarshal(body, &mapp); err != nil {
		c.jsonResult(enums.JRCodeFailed, "数据异常:", 0)
	}
	id := mapp["userId"]
	obj, err := models.TokenskyUserOne(id)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "无此用户", id)
	}
	uid := strconv.Itoa(obj.UserId)
	addr,found := utils.EncryptWithAESUrl(uid)
	if !found{
		c.jsonResult(enums.JRCodeFailed, "无此用户", id)
	}
	data := map[string]string{
		"userId":uid,
		"addr":addr,
	}
	c.jsonResult(enums.JRCodeSucc, "ok", data)
}