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

type OperationBannerController struct {
	BaseController
}

//Prepare 参考beego官方文档说明
func (c *OperationBannerController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	c.checkAuthor()
}

// DataGrid
func (c *OperationBannerController) DataGrid() {
	//直接反序化获取json格式的requestbody里的值
	var params models.OperationBannerQueryParam
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
		params.Status, _ = c.GetInt("status")
	}
	//获取数据列表和总数
	data, total := models.OperationBannerPageList(&params)
	//定义返回的数据结构
	mapp := map[string]interface{}{
		"rows":  data,
		"total": total,
	}
	c.jsonResult(enums.JRCodeSucc, "", mapp)
}

//Edit 添加、编辑
func (c *OperationBannerController) Edit() {
	if c.Ctx.Request.Method == "POST" {
		c.Save()
	}
}

//Save 添加、编辑页面 保存
func (c *OperationBannerController) Save() {
	var err error
	m := models.OperationBanner{}
	body := c.Ctx.Input.RequestBody //接收raw body内容
	if err := json.Unmarshal(body, &m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "数据异常:", m.Bid)
	}
	switch m.Status {
	case conf.OPERATION_BANNER_STATUS_OFF:
	case conf.OPERATION_BANNER_STATUS_NO:
	default:
		c.jsonResult(enums.JRCodeFailed, "未知状态:", m.Bid)
	}
	o := orm.NewOrm()
	obj := models.OperationBannerOneById(m.Bid)
	if obj == nil{
		c.jsonResult(enums.JRCodeFailed, "Banner不存在", m.Bid)
	}
	m.AdminId = c.curUser.Id
	if m.Bid <= 0 {
		//新增
		if _, err = o.Insert(&m); err == nil {
			c.jsonResult(enums.JRCodeSucc, "添加成功", m.Bid)
		} else {
			c.jsonResult(enums.JRCodeFailed, "添加失败，可能原因："+err.Error(), m.Bid)
		}
	} else {
		obj.AdminId = m.AdminId
		obj.Name = m.Name
		obj.Seq = m.Seq
		obj.Url = m.Url
		obj.Status = m.Status
		obj.ImgKey = m.ImgKey
		obj.ImgUrl = m.ImgUrl
		if _, err = o.Update(obj); err == nil {
			c.jsonResult(enums.JRCodeSucc, "编辑成功", m.Bid)
		} else {
			c.jsonResult(enums.JRCodeFailed, "编辑失败", m.Bid)
		}
	}
}

//Delete 批量删除
func (c *OperationBannerController) Delete() {
	body := c.Ctx.Input.RequestBody
	mapp := make(map[string]string)
	json.Unmarshal(body, &mapp)
	ids := make([]int, 0)
	for _, str := range strings.Split(mapp["ids"], ",") {
		if id, err := strconv.Atoi(str); err == nil {
			ids = append(ids, id)
		}
	}
	if num, err := models.OperationBannerDelete(ids); err == nil {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), 0)
	} else {
		c.jsonResult(enums.JRCodeFailed, "删除失败", 0)
	}
}
