package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"time"
	"tokensky_bg_admin/enums"
	"tokensky_bg_admin/models"
)

//
type BorrowConfController struct {
	BaseController
}

//Prepare 参考beego官方文档说明
func (c *BorrowConfController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	c.checkAuthor()
}

// DataGrid
func (c *BorrowConfController) DataGrid() {
	var params models.BorrowConfQueryParam
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

		params.CoinType = c.GetString("coinType")
		params.LoanSymbol = c.GetString("loanSymbol")
		params.Title = c.GetString("title")
	}
	//获取数据列表和总数
	data, total := models.BorrowConfPageList(&params)
	//定义返回的数据结构
	mapp := map[string]interface{}{
		"rows":  data,
		"total": total,
	}
	c.jsonResult(enums.JRCodeSucc, "", mapp)
}

//编辑
func (c *BorrowConfController) Edit() {
	if c.Ctx.Request.Method == "POST" {
		var err error
		m := models.BorrowConf{}
		body := c.Ctx.Input.RequestBody
		if err = json.Unmarshal(body, &m); err != nil {
			c.jsonResult(enums.JRCodeFailed, "数据异常:", m.Id)
		}
		m.AdminId = c.curUser.Id
		o := orm.NewOrm()
		switch m.Id {
		case 0:
			//新增只能是待上架状态
			m.IsPutaway = 0
			//数据校验
			if !models.TokenskyUserBalanceCoinIsFound(m.CoinType){
				c.jsonResult(enums.JRCodeFailed, "质押货币类型 不存在", m.Id)
			}
			if !models.TokenskyUserBalanceCoinIsFound(m.LoanSymbol){
				c.jsonResult(enums.JRCodeFailed, "借贷货币类型 不存在", m.Id)
			}
			m.CreateTime = time.Now().Unix()*1000

			//事务
			err = o.Begin()
			if err != nil{
				c.jsonResult(enums.JRCodeFailed, "事务异常", 0)
			}
			id,err := o.Insert(&m)
			if err != nil{
				o.Rollback()
				c.jsonResult(enums.JRCodeFailed, "添加失败", m.Id)
			}
			adminModelRecord := &models.AdminModelRecord{
				Uid:c.curUser.Id,
				Handle:"新增",
				Model:"BorrowConf",
				Tbid:strconv.Itoa(int(id)),
				NewData:m.Json(),
			}
			_,err = o.Insert(adminModelRecord)
			if err != nil{
				o.Rollback()
				c.jsonResult(enums.JRCodeFailed, "添加记录表失败", m.Id)
			}
			err = o.Commit()
			if err != nil{
				c.jsonResult(enums.JRCodeFailed, "开启事务失败", m.Id)
			}
			c.jsonResult(enums.JRCodeSucc, "添加成功", m.Id)
		default:
			//编辑
			obj := models.BorrowConfOne(m.Id)
			if obj == nil{
				c.jsonResult(enums.JRCodeFailed, "编辑数据不存在", nil)
			}
			//非0
			if obj.IsPutaway!= 0{
				c.jsonResult(enums.JRCodeFailed, "只有待上架状态数据可以编辑", nil)
			}
			m.CoinType = obj.CoinType
			m.LoanSymbol = obj.LoanSymbol
			m.CreateTime = obj.CreateTime
			//事务
			err = o.Begin()
			if err != nil{
				c.jsonResult(enums.JRCodeFailed, "事务开启失败", 0)
			}
			_,err := o.Update(&m)
			if err != nil{
				o.Rollback()
				c.jsonResult(enums.JRCodeFailed, "编辑错误", nil)
			}
			adminModelRecord := &models.AdminModelRecord{
				Uid:c.curUser.Id,
				Handle:"编辑",
				Model:"BorrowConf",
				Tbid:strconv.Itoa(obj.Id),
				OldData:obj.Json(),
			}
			_,err = o.Insert(adminModelRecord)
			if err != nil{
				o.Rollback()
				c.jsonResult(enums.JRCodeFailed, "新增记录异常", nil)
			}
			err = o.Commit()
			if err != nil {
				c.jsonResult(enums.JRCodeSucc, "事务执行异常", m.Id)
			}
			c.jsonResult(enums.JRCodeSucc, "修改成功", m.Id)
		}
	}
	c.jsonResult(enums.JRCodeFailed, "请求错误", nil)
}

//上架下架
func (c *BorrowConfController) TheUpper() {
	if c.Ctx.Request.Method == "POST" {
		body := c.Ctx.Input.RequestBody
		mapp := make(map[string]string)
		err := json.Unmarshal(body, &mapp)
		if err != nil {
			c.jsonResult(enums.JRCodeFailed, "数据解析异常", 0)
		}
		var status int
		ids := make([]int,0)
		if con,err := strconv.Atoi(mapp["isPutaway"]);err!=nil{
			c.jsonResult(enums.JRCodeFailed, "isPutaway数据解析异常", 0)
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
		record := &models.AdminModelRecord{
			Uid:c.curUser.Id,
			Handle:msg,
			Model:"BorrowConf",
			Tbid:mapp["ids"],
			NewData:strconv.Itoa(status),
		}
		o:=orm.NewOrm()
		o.Insert(record)
		num, err := models.BorrowConfTheUppers(o,status, ids,c.curUser.Id)
		if err == nil {
			c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功"+msg+" %d 项", num), 0)
		} else {
			c.jsonResult(enums.JRCodeFailed, msg+"失败", 0)
		}

	}
	c.jsonResult(enums.JRCodeFailed, "请求错误", nil)
}