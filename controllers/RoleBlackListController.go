package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"time"
	"tokensky_bg_admin/conf"
	"tokensky_bg_admin/enums"
	"tokensky_bg_admin/models"
)

//RoleBlackListController 黑名单
type RoleBlackListController struct {
	BaseController
}

//Prepare 参考beego官方文档说明
func (c *RoleBlackListController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	c.checkAuthor()

}

// DataGrid
func (c *RoleBlackListController) DataGrid() {
	//直接反序化获取json格式的requestbody里的值
	var params models.RoleBlackListQueryParam
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
	}
	//获取数据列表和总数
	data, total := models.RoleBlackListPageList(&params)
	//定义返回的数据结构
	mapp := map[string]interface{}{
		"rows":  data,
		"total": total,
	}
	c.jsonResult(enums.JRCodeSucc, "", mapp)
}

//Edit 添加、编辑
func (c *RoleBlackListController) Edit() {
	if c.Ctx.Request.Method == "POST" {
		c.Save()
	}
}

//Save 添加、编辑页面 保存
func (c *RoleBlackListController) Save() {
	var err error
	m := models.RoleBlackList{}
	body := c.Ctx.Input.RequestBody //接收raw body内容
	if err := json.Unmarshal(body, &m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "数据异常:", m.Id)
	}
	//状态校验
	switch m.BalckType {
	case conf.ROLE_BLACK_LIST_STATUS_BAN_LANDING:
	case conf.ROLE_BLACK_LIST_STATUS_BAN_TRADING:
	default:
		//无此状态
		c.jsonResult(enums.JRCodeFailed, "未知的封号类型", m.Id)
	}
	//判断是否存在用户
	user := models.TokenskyUserOneByPhone(m.Phone)
	if user == nil {
		c.jsonResult(enums.JRCodeFailed, "无该手机号用户", m.Id)
	}
	o := orm.NewOrm()
	switch m.Id {
	case 0:
		//开始时间
		m.StartTime = time.Now()
		//结束时间
		m.EndTime = time.Unix(m.StartTime.Unix()+m.PeriodTime, 0)
		m.User = user
		if _, err = o.Insert(&m); err == nil {
			c.jsonResult(enums.JRCodeSucc, "添加成功", m.Id)
		} else {
			c.jsonResult(enums.JRCodeFailed, "添加失败，可能原因："+err.Error(), m.Id)
		}
	default:
		//编辑
		obj, err := models.RoleBlackListOneById(m.Id)
		if err != nil {
			c.jsonResult(enums.JRCodeFailed, "获取记录失败", m.Id)
		}
		obj.BalckType = m.BalckType
		obj.PeriodTime = m.PeriodTime
		//手机号不可编辑
		//obj.Phone = m.Phone
		endTime := time.Unix(obj.StartTime.Unix()+m.PeriodTime, 0)
		obj.EndTime = endTime
		if _, err = o.Update(obj); err == nil {
			c.jsonResult(enums.JRCodeSucc, "编辑成功", obj.Id)
		} else {
			c.jsonResult(enums.JRCodeFailed, "编辑失败", obj.Id)
		}
	}
}

//Delete 批量删除
func (c *RoleBlackListController) Delete() {
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
		if num, err := models.RoleBlackListDelete(ids); err == nil {
			c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), 0)
		} else {
			c.jsonResult(enums.JRCodeFailed, "删除失败", 0)
		}
	}
	c.jsonResult(enums.JRCodeFailed, "没有删除对象", 0)
}
