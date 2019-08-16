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

//理财配置
type FinancialProductController struct {
	BaseController
}

//Prepare 参考beego官方文档说明
func (c *FinancialProductController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	c.checkAuthor()
}

// DataGrid
func (c *FinancialProductController) DataGrid() {
	var params models.FinancialProductQueryParam
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
		params.Status = c.GetString("status")
		params.FinancialCategoryId = c.GetString("financialCategoryId")
		params.Category = c.GetString("category")
	}
	//获取数据列表和总数
	data, total := models.FinancialProductPageList(&params)
	//定义返回的数据结构
	mapp := map[string]interface{}{
		"rows":  data,
		"total": total,
	}
	c.jsonResult(enums.JRCodeSucc, "", mapp)
}

//编辑
func (c *FinancialProductController) Edit() {
	if c.Ctx.Request.Method == "POST" {
		var err error
		m := models.FinancialProduct{}
		body := c.Ctx.Input.RequestBody
		if err = json.Unmarshal(body, &m); err != nil {
			c.jsonResult(enums.JRCodeFailed, "数据异常:", m.Id)
		}
		o := orm.NewOrm()
		//关联
		m.FinancialCategoryObj = &models.FinancialCategory{Id: m.FinancialCategoryId}
		switch m.Id {
		case 0:
			//新增
			switch m.Category {
			case 1:
				//活期重复校验
				if !models.FinancialProductIsAddObj(m.FinancialCategoryId) {
					c.jsonResult(enums.JRCodeFailed, "该活息配置重复添加", m.Id)
				}
			case 2:
			default:
				c.jsonResult(enums.JRCodeFailed, "未知类型", m.Id)
			}
			m.Admin = &models.AdminBackendUser{Id:c.curUser.Id}
			//事务
			o.Begin()
			id, err := o.Insert(&m)
			if err != nil {
				o.Rollback()
				c.jsonResult(enums.JRCodeFailed, "添加失败", m.Id)
			}
			//
			newRecord := &models.FinancialProductHistoricalRecord{
				Config:     &models.FinancialProduct{Id: int(id)},
				Admin:      &models.AdminBackendUser{Id: c.curUser.Id},
				Msg:        m.Msg,
				NewRate:    m.YearProfit,
				OldRate:    0,
				Category:   m.Category,
				RecordType: "新增",
			}
			_, err = o.Insert(newRecord)
			if err != nil {
				o.Rollback()
				c.jsonResult(enums.JRCodeFailed, "添加历史记录失败", m.Id)
			}
			err = o.Commit()
			if err != nil {
				o.Rollback()
				c.jsonResult(enums.JRCodeFailed, "事务失败", m.Id)
			}
			c.jsonResult(enums.JRCodeSucc, "添加成功", m.Id)
		default:
			//编辑
			obj := models.FinancialProductOne(m.Id)
			if obj == nil {
				c.jsonResult(enums.JRCodeFailed, "编辑记录不存在", m.Id)
			}
			//状态判断[仅针对定期]
			if obj.Status != 0 && obj.Category == 2 {
				c.jsonResult(enums.JRCodeFailed, "只有待上架状态可以编辑", m.Id)
			}
			//不可修改
			m.FinancialCategoryObj = obj.FinancialCategoryObj
			m.Category = obj.Category
			//类型
			switch obj.Category {
			case 1:
				//活期状态不能修改,必须上架
				m.Status = 1
			}
			newRecord := &models.FinancialProductHistoricalRecord{
				Config:     &models.FinancialProduct{Id: obj.Id},
				Admin:      &models.AdminBackendUser{Id: c.curUser.Id},
				Msg:        m.Msg,
				NewRate:    m.YearProfit,
				OldRate:    obj.YearProfit,
				Category:   obj.Category,
				RecordType: "编辑",
			}
			o.Begin()
			if _, err := o.Insert(newRecord); err != nil {
				o.Rollback()
				c.jsonResult(enums.JRCodeFailed, "新增记录异常", m.Id)
			}
			if _, err := o.Update(&m); err != nil {
				o.Rollback()
				c.jsonResult(enums.JRCodeFailed, "编辑配置异常", m.Id)
			}
			if err := o.Commit(); err != nil {
				o.Rollback()
				c.jsonResult(enums.JRCodeFailed, "事务异常", m.Id)
			}
			c.jsonResult(enums.JRCodeSucc, "编辑成功", m.Id)
		}
	}
	c.jsonResult(enums.JRCodeFailed, "请求错误", nil)
}

//上下架
func (c *FinancialProductController) TheUpper() {
	if c.Ctx.Request.Method == "POST" {
		body := c.Ctx.Input.RequestBody
		mapp := make(map[string]string)
		err := json.Unmarshal(body,&mapp)
		if err != nil{
			c.jsonResult(enums.JRCodeFailed, "数据解析异常", 0)
		}

		var status int
		ids := make([]int,0)
		if con,err := strconv.Atoi(mapp["status"]);err!=nil{
			c.jsonResult(enums.JRCodeFailed, "status数据解析异常", 0)
		}else {
			status = con
		}
		for _,str := range strings.Split(mapp["ids"],","){
			if con,err := strconv.Atoi(str);err !=nil{
				c.jsonResult(enums.JRCodeFailed, "ids数据解析异常", 0)
			}else {
				ids = append(ids, con)
			}
		}
		msg := ""
		switch status {
		case 1:
			msg = "上架"
		case 2:
			msg = "下架"
		default:
			c.jsonResult(enums.JRCodeFailed, "不支持的操作", 0)
		}
		num, err := models.FinancialProductTheUppers(status, ids,c.curUser.Id)
		if err == nil {
			c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功"+msg+" %d 项", num), 0)
		} else {
			c.jsonResult(enums.JRCodeFailed, msg+"失败", 0)
		}
	}
	c.jsonResult(enums.JRCodeFailed, "请求错误", nil)
}

//删除
func (c *FinancialProductController) Delete() {
	if c.Ctx.Request.Method == "POST" {
		body := c.Ctx.Input.RequestBody
		mapp := make(map[string]string)
		var id int
		err := json.Unmarshal(body,&mapp)
		if err !=nil{
			c.jsonResult(enums.JRCodeFailed, "解析异常", nil)
		}
		if con,err := strconv.Atoi(mapp["id"]);err !=nil{
			c.jsonResult(enums.JRCodeFailed, "解析异常", nil)
		}else {
			id = con
		}
		ok, msg := models.FinancialProductDelete(id)
		if ok {
			c.jsonResult(enums.JRCodeSucc, "ok", nil)
		} else {
			c.jsonResult(enums.JRCodeFailed, msg, nil)
		}
	}
	c.jsonResult(enums.JRCodeFailed, "请求错误", nil)
}
