package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"strings"
	"tokensky_bg_admin/enums"
	"tokensky_bg_admin/models"
)

//提币配置
type TokenskyTibiConfigController struct {
	BaseController
}

//Prepare 参考beego官方文档说明
func (c *TokenskyTibiConfigController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	c.checkAuthor()
}

// DataGrid
func (c *TokenskyTibiConfigController) DataGrid() {
	//直接反序化获取json格式的requestbody里的值
	var params models.TokenskyTibiConfigQueryParam
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
	data, total := models.TokenskyTibiConfigPageList(&params)
	//定义返回的数据结构
	mapp := map[string]interface{}{
		"rows":  data,
		"total": total,
	}
	c.jsonResult(enums.JRCodeSucc, "", mapp)
}

//Edit 添加
func (c *TokenskyTibiConfigController) Edit() {
	if c.Ctx.Request.Method == "POST" {
		c.Save()
	}
}

//Save 添加、编辑页面 保存
func (c *TokenskyTibiConfigController) Save() {
	var err error
	m := models.TokenskyTibiConfig{}
	mBak := models.TokenskyTibiConfigBak{}
	body := c.Ctx.Input.RequestBody //接收raw body内容
	if err := json.Unmarshal(body, &m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "数据异常", 0)
	}
	if err := json.Unmarshal(body, &mBak); err != nil {
		c.jsonResult(enums.JRCodeFailed, "数据异常", 0)
	}
	m.CoinType = strings.ToUpper(m.CoinType)
	mBak.AdminId = c.BaseController.curUser.Id
	obj := models.TokenskyTibiConfigGetLastOne(m.CoinType)
	switch m.Id {
	case 0:
		//新增
		if obj != nil {
			c.jsonResult(enums.JRCodeFailed, "该货币配置已存在", 0)
		}
		m.AdminId = c.BaseController.curUser.Id

	default:
		//修改
		mBak.Id = 0
	}
	o := orm.NewOrm()
	m.AdminId = c.BaseController.curUser.Id
	mBak.AdminId = c.BaseController.curUser.Id
	mBak.Id = 0
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
		//提币配置名称不可变
		m.CoinType = obj.CoinType
		mBak.CoinType = obj.CoinType
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
