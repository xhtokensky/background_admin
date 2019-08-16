package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"tokensky_bg_admin/conf"
	"tokensky_bg_admin/enums"
	"tokensky_bg_admin/models"
)

//TokenskyMessageController 消息
type TokenskyMessageController struct {
	BaseController
}

//Prepare 参考beego官方文档说明
func (c *TokenskyMessageController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	c.checkAuthor()

}

// DataGrid 消息信息
func (c *TokenskyMessageController) DataGrid() {
	//直接反序化获取json格式的requestbody里的值
	var params models.TokenskyMessageQueryParam
	switch c.Ctx.Request.Method {
	case "POST":
		json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	case "GET":
		params.Sort = c.GetString("sort")
		params.Order = c.GetString("order")
		params.Limit, _ = c.GetInt64("limit")
		params.Offset, _ = c.GetInt64("offset")
		params.Status = c.GetString("status")
		params.StartTime, _ = c.GetInt64("startTime")
		params.EndTime, _ = c.GetInt64("endTime")
		params.Phone = c.GetString("phone")
	}
	if params.Order != conf.QUREY_PARAM_ORDER_ASC {
		params.Order = conf.QUREY_PARAM_ORDER_DESC
	}
	if params.Sort == "" {
		params.Sort = "messageId"
	}

	//获取数据列表和总数
	data, total := models.TokenskyMessagePageList(&params)
	//定义返回的数据结构
	mapp := make(map[string]interface{})
	mapp["rows"] = data
	mapp["total"] = total
	c.jsonResult(enums.JRCodeSucc, "", mapp)
}

//Edit 添加、编辑
func (c *TokenskyMessageController) Edit() {
	if c.Ctx.Request.Method == "POST" {
		c.Save()
	}
}

//Save 添加、编辑页面 保存
func (c *TokenskyMessageController) Save() {
	var err error
	m := models.TokenskyMessage{}
	body := c.Ctx.Input.RequestBody //接收raw body内容
	if err := json.Unmarshal(body, &m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "数据异常:", m.MessageId)
	}
	o := orm.NewOrm()
	//
	if m.Title == "" {
		c.jsonResult(enums.JRCodeFailed, "标题为空", m.MessageId)
	}
	if m.Content == "" {
		c.jsonResult(enums.JRCodeFailed, "消息为空", m.MessageId)
	}
	//默认1
	m.Status = 1
	m.AdminId = c.curUser.Id
	m.EditorId = c.curUser.Id
	//类型校验 个人用户 必须需要电话信息
	switch m.Type {
	case 0:
		//全部
		m.Phone = ""
	case 1:
		//个人
		if m.Phone == "" {
			c.jsonResult(enums.JRCodeFailed, "电话为空", m.MessageId)
		}
		//电话
		if obj := models.TokenskyUserOneByPhone(m.Phone); obj == nil {
			c.jsonResult(enums.JRCodeFailed, "无此用户", m.MessageId)
		}else {
			m.User = obj
		}
	}
	if m.MessageId == 0 {
		if _, err = o.Insert(&m); err == nil {
			c.jsonResult(enums.JRCodeSucc, "添加成功", m.MessageId)
		} else {
			c.jsonResult(enums.JRCodeFailed, "添加失败", m.MessageId)
		}
	} else {
		if _, err = o.Update(&m); err == nil {
			c.jsonResult(enums.JRCodeSucc, "编辑成功", m.MessageId)
		} else {
			c.jsonResult(enums.JRCodeFailed, "编辑失败", m.MessageId)
		}
	}
}

//Delete 批量删除
func (c *TokenskyMessageController) Delete() {
	body := c.Ctx.Input.RequestBody
	mapp := make(map[string]string)
	json.Unmarshal(body, &mapp)
	ids := make([]int, 0)
	for _, str := range strings.Split(mapp["ids"], ",") {
		if id, err := strconv.Atoi(str); err == nil {
			ids = append(ids, id)
		}
	}
	if len(ids) > 0 {
		/*删除操作*/
		//if num, err := models.TokenskyMessageDelete(ids); err == nil {
		//	c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), 0)
		//} else {
		//	c.jsonResult(enums.JRCodeFailed, "删除失败", 0)
		//}

		//状态修改操作
		params := map[string]interface{}{
			"status":   0,
			"admin_id": c.curUser.Id,
		}
		query := orm.NewOrm().QueryTable(models.TokenskyMessageTBName())
		if num, err := query.Filter("message_id__in", ids).Update(params); err == nil {
			c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), 0)
		} else {
			c.jsonResult(enums.JRCodeFailed, "删除失败", 0)
		}
	} else {
		c.jsonResult(enums.JRCodeSucc, "缺少要删除字段ids", 0)
	}

}
