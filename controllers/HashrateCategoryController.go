package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"tokensky_bg_admin/enums"
	"tokensky_bg_admin/models"
)

//算力合约分类表
type HashrateCategoryController struct {
	BaseController
}

//Prepare 参考beego官方文档说明
func (c *HashrateCategoryController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	c.checkAuthor()

}

// DataGrid
func (c *HashrateCategoryController) DataGrid() {
	//直接反序化获取json格式的requestbody里的值
	var params models.HashrateCategoryQueryParam
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
	data, total := models.HashrateCategoryPageList(&params)
	//定义返回的数据结构
	mapp := map[string]interface{}{
		"rows":  data,
		"total": total,
	}
	c.jsonResult(enums.JRCodeSucc, "", mapp)
}

//编辑
func (c *HashrateCategoryController) Edit() {
	if c.Ctx.Request.Method == "POST" {
		var err error
		m := models.HashrateCategory{}
		body := c.Ctx.Input.RequestBody
		if err = json.Unmarshal(body, &m); err != nil {
			c.jsonResult(enums.JRCodeFailed, "数据异常:", m.KeyId)
		}
		o := orm.NewOrm()

		switch m.Unit {
		case "H", "K", "M","G", "T", "P", "E":
		default:
			c.jsonResult(enums.JRCodeFailed, "未知算力合约单位,算力合约支持单位 'H','G','T','P','E'", m.KeyId)
		}
		m.AdminUserId = c.curUser.Id
		switch m.KeyId {
		case 0:
			//新增
			if _, err = o.Insert(&m); err == nil {
				c.jsonResult(enums.JRCodeSucc, "添加成功", m.KeyId)
			} else {
				c.jsonResult(enums.JRCodeFailed, "添加失败", m.KeyId)
			}
		default:
			//编辑
			if _, err = o.InsertOrUpdate(&m); err == nil {
				c.jsonResult(enums.JRCodeSucc, "编辑成功", m.KeyId)
			} else {
				c.jsonResult(enums.JRCodeFailed, "编辑失败", m.KeyId)
			}
		}
	}
}

//删除
func (c *HashrateCategoryController) Delete() {
	c.jsonResult(enums.JRCodeSucc, "算力合约分类 不支持删除", 0)

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
		if num, err := models.HashrateCategoryDelete(ids); err == nil {
			c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), 0)
		} else {
			c.jsonResult(enums.JRCodeFailed, "删除失败", 0)
		}
	} else {
		c.jsonResult(enums.JRCodeFailed, "缺少要删除字段ids", 0)
	}
}
