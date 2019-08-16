package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"tokensky_bg_admin/enums"
	"tokensky_bg_admin/models"
)

//提币配置
type TokenskyChongbiConfigController struct {
	BaseController
}

//Prepare 参考beego官方文档说明
func (c *TokenskyChongbiConfigController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	c.checkAuthor()

}

// DataGrid
func (c *TokenskyChongbiConfigController) DataGrid() {
	//直接反序化获取json格式的requestbody里的值
	var params models.TokenskyChongbiConfigQueryParam
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
	data, total := models.TokenskyChongbiConfigPageList(&params)
	//定义返回的数据结构
	mapp := map[string]interface{}{
		"rows":  data,
		"total": total,
	}
	c.jsonResult(enums.JRCodeSucc, "", mapp)
}

//Edit 添加
func (c *TokenskyChongbiConfigController) Edit() {
	if c.Ctx.Request.Method == "POST" {
		c.Save()
	}
}

//Save 添加、编辑页面 保存
func (c *TokenskyChongbiConfigController) Save() {
	var err error
	m := models.TokenskyChongbiConfig{}
	mBak := models.TokenskyChongbiConfigBak{}
	body := c.Ctx.Input.RequestBody //接收raw body内容
	if err := json.Unmarshal(body, &m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "数据异常", 0)
	}
	if err := json.Unmarshal(body, &mBak); err != nil {
		c.jsonResult(enums.JRCodeFailed, "数据异常", 0)
	}
	mBak.AdminId = c.BaseController.curUser.Id
	obj := models.TokenskyChongbiConfigGetLastOne(m.CoinType)
	switch m.Id {
	case 0:
		//校验该配置是否存在
		if obj != nil {
			c.jsonResult(enums.JRCodeFailed, "该货币配置已存在", 0)
		}
		m.AdminId = c.BaseController.curUser.Id
	default:
		mBak.Id = 0
	}
	o := orm.NewOrm()
	//修改本表
	switch m.Id {
	case 0:
		//新增
		if _, err = o.Insert(&m); err == nil {
			o.Insert(&mBak)
			c.jsonResult(enums.JRCodeSucc, "新增成功", 0)
		} else {
			c.jsonResult(enums.JRCodeFailed, "新增失败", 0)
		}
	default:
		//编辑
		if obj == nil{
			c.jsonResult(enums.JRCodeFailed, "编辑对象不存在", 0)
		}
		m.CoinType = obj.CoinType
		m.CreateTime = obj.CreateTime
		if _, err = o.Update(&m); err == nil {
			//新增副表
			o.Insert(&mBak)
			c.jsonResult(enums.JRCodeSucc, "编辑成功", 0)
		} else {
			c.jsonResult(enums.JRCodeFailed, "编辑失败", 0)
		}
	}

}
