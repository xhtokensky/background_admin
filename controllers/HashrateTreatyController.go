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

//算力合约表
type HashrateTreatyController struct {
	BaseController
}

//Prepare 参考beego官方文档说明
func (c *HashrateTreatyController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	c.checkAuthor()
}

// DataGrid
func (c *HashrateTreatyController) DataGrid() {
	//直接反序化获取json格式的requestbody里的值
	var params models.HashrateTreatyQueryParam
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
		params.Title = c.GetString("title")
		params.Status = c.GetString("status")
	}
	//获取数据列表和总数
	data, total := models.HashrateTreatyPageList(&params)
	//定义返回的数据结构
	mapp := map[string]interface{}{
		"rows":  data,
		"total": total,
	}
	c.jsonResult(enums.JRCodeSucc, "", mapp)
}

//编辑
func (c *HashrateTreatyController) Edit() {
	if c.Ctx.Request.Method == "POST" {
		var err error
		m := models.HashrateTreaty{}
		body := c.Ctx.Input.RequestBody
		if err = json.Unmarshal(body, &m); err != nil {
			c.jsonResult(enums.JRCodeFailed, "数据异常", m.KeyId)
		}
		m.EffectiveDate = time.Unix(m.EffectiveDateTm, 0)
		o := orm.NewOrm()

		switch m.KeyId {
		case 0:
			//新增
			m.AdminUserId = c.curUser.Id  //创建人id
			m.InventoryLeft = m.Inventory //新增数据一致
			m.UnitMoney = "CNY"           //默认
			m.Status = 0                  //默认待上架
			switch m.Tag {
			case 0:
			//普通版
				if m.OriginalPrice != m.Price{
					m.Price = m.OriginalPrice
					//c.jsonResult(enums.JRCodeFailed, "原价现价 价格需一直", m.KeyId)
				}
			case 1:
			//优惠版
				if m.OriginalPrice > m.Price{
					c.jsonResult(enums.JRCodeFailed, "原价高于现价", m.KeyId)
				}
				//优惠版和期货 不可同时出现
				if m.FuturesType == 1{
					c.jsonResult(enums.JRCodeFailed, "优惠版和期货 不可同时出现", m.KeyId)
				}
			}

			//获取合约分类表
			category := models.HashrateCategoryById(m.HashrateCategory)
			if category == nil {
				c.jsonResult(enums.JRCodeFailed, "无此分类", 0)
			}
			m.HashrateCategoryObj = &models.HashrateCategory{KeyId: m.HashrateCategory}
			//事务
			err = o.Begin()
			if err != nil{
				c.jsonResult(enums.JRCodeFailed, "事务异常", 0)
			}
			id,err := o.Insert(&m)
			if err != nil{
				o.Rollback()
				c.jsonResult(enums.JRCodeFailed, "添加失败", m.KeyId)
			}
			adminModelRecord := &models.AdminModelRecord{
				Uid:c.curUser.Id,
				Handle:"新增",
				Model:"HashrateTreaty",
				Tbid:strconv.Itoa(int(id)),
				NewData:m.Json(),
			}
			_,err = o.Insert(adminModelRecord)
			if err != nil{
				o.Rollback()
				c.jsonResult(enums.JRCodeFailed, "添加记录表失败", m.KeyId)
			}
			err = o.Commit()
			if err != nil{
				c.jsonResult(enums.JRCodeFailed, "开启事务失败", m.KeyId)
			}
			c.jsonResult(enums.JRCodeSucc, "添加成功", m.KeyId)
		default:
			//编辑
			obj := models.HashrateTreatyOneById(m.KeyId)
			if obj == nil {
				c.jsonResult(enums.JRCodeFailed, "算力合约表 编辑数据不存在", 0)
			}
			if obj.Status != 0 {
				//上架只能编辑电费信息
				obj.ElectricBill = m.ElectricBill
				obj.Management = m.Management
				obj.EarningsRate = m.EarningsRate
			}
			//事务
			err = o.Begin()
			if err != nil{
				c.jsonResult(enums.JRCodeFailed, "事务开启失败", 0)
			}
			adminModelRecord := &models.AdminModelRecord{
				Uid:c.curUser.Id,
				Handle:"编辑",
				Model:"HashrateTreaty",
				Tbid:strconv.Itoa(obj.KeyId),
				OldData:obj.Json(),
			}
			obj.EffectiveDateTm = obj.EffectiveDate.Unix()
			if obj.EffectiveDateTm != m.EffectiveDateTm {
				if obj.FuturesType == 1 && m.FuturesType == 1 {
					//期货类型只允许修改生效日期
					now := time.Now()
					tm1 := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
					newTm := time.Unix(m.EffectiveDateTm, 0)
					tm3 := time.Date(newTm.Year(), newTm.Month(), newTm.Day(), 0, 0, 0, 0, newTm.Location())
					if tm1.Unix()+7200 >= tm3.Unix() {
						o.Rollback()
						c.jsonResult(enums.JRCodeFailed, "当前距离挖矿时间过近，无法修改", m.KeyId)
					}
					obj.EffectiveDate = newTm
					//生效时间更新 购买订单的生效时间修改
					_, err := models.HashrateOrderIsNotFutures(o, obj)
					if err != nil {
						o.Rollback()
						c.jsonResult(enums.JRCodeFailed, "事务更新失败", m.KeyId)
					}
					obj.AdminUserId = c.curUser.Id
					_, err = o.Update(&obj)
					if err != nil {
						o.Rollback()
						c.jsonResult(enums.JRCodeSucc, "编辑失败", m.KeyId)
					}
					adminModelRecord.NewData = obj.Json()
					_,err = o.Insert(adminModelRecord)
					if err != nil{
						c.jsonResult(enums.JRCodeSucc, "新增记录表失败", m.KeyId)
					}
					err = o.Commit()
					if err != nil {
						c.jsonResult(enums.JRCodeSucc, "事务执行异常", m.KeyId)
					}
					c.jsonResult(enums.JRCodeSucc, "修改成功", m.KeyId)
				}
			}
			adminModelRecord.NewData = obj.Json()
			_,err = o.Insert(adminModelRecord)
			if err != nil{
				o.Rollback()
				c.jsonResult(enums.JRCodeSucc, "新增记录表失败", m.KeyId)
			}
			//编辑人信息
			obj.AdminUserId = c.curUser.Id
			_,err = o.Update(obj)
			if err != nil{
				o.Rollback()
				c.jsonResult(enums.JRCodeFailed, "编辑失败", m.KeyId)
			}
			err = o.Commit()
			if err != nil {
				c.jsonResult(enums.JRCodeSucc, "事务执行异常", m.KeyId)
			}
			c.jsonResult(enums.JRCodeSucc, "修改成功", m.KeyId)
		}
	}
	c.jsonResult(enums.JRCodeFailed, "请求错误", 0)
}

//
//算力合约期货上架漏订单补救措施[期货修改过程中，提交的订单]
func (c *HashrateTreatyController) IsNotFutures() {
	body := c.Ctx.Input.RequestBody
	mapp := make(map[string]string)
	if err := json.Unmarshal(body, &mapp); err != nil {
		c.jsonResult(enums.JRCodeSucc, "解析异常", 0)
	}
	var id int
	id, err := strconv.Atoi(mapp["id"])
	if err != nil {
		c.jsonResult(enums.JRCodeSucc, "解析异常", 0)
	}
	//获取单条
	obj := models.HashrateTreatyOneById(id)
	if obj == nil {
		c.jsonResult(enums.JRCodeSucc, "数据不存在", 0)
	}
	if obj.FuturesType == 1 {
		c.jsonResult(enums.JRCodeFailed, "算力合约表 只有上架可以修改", 0)
	}
	o := orm.NewOrm()
	//1事务校验
	err = o.Begin()
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "开启事务失败", obj.KeyId)
	}
	//2时间状态更新
	_, err = models.HashrateOrderIsNotFutures(o, obj)
	if err != nil {
		o.Rollback()
		c.jsonResult(enums.JRCodeFailed, "事务更新失败", obj.KeyId)
	}
	err = o.Commit()
	if err != nil {
		c.jsonResult(enums.JRCodeSucc, "事务执行成功", obj.KeyId)
	}
}

//删除
func (c *HashrateTreatyController) Delete() {
	body := c.Ctx.Input.RequestBody
	mapp := make(map[string]string)
	if err := json.Unmarshal(body, &mapp); err != nil {
		c.jsonResult(enums.JRCodeSucc, "解析异常", 0)
	}
	ids := make([]int, 0)
	for _, str := range strings.Split(mapp["ids"], ",") {
		if id, err := strconv.Atoi(str); err == nil {
			ids = append(ids, id)
		}
	}
	if len(ids) > 0 {
		if num, err := models.HashrateTreatyDelete(ids); err == nil {
			c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项 ps:只有待上架待可以删除", num), 0)
		} else {
			c.jsonResult(enums.JRCodeFailed, "删除失败", 0)
		}
	} else {
		c.jsonResult(enums.JRCodeSucc, "缺少要删除字段ids", 0)
	}
}

//算力合约上架下架
func (c *HashrateTreatyController) Shelves() {
	body := c.Ctx.Input.RequestBody
	mapp := make(map[string]string)
	if err := json.Unmarshal(body, &mapp); err != nil {
		c.jsonResult(enums.JRCodeFailed, "解析异常", 0)
	}
	ids := make([]int, 0)
	var status int
	if con, err := strconv.Atoi(mapp["status"]); err == nil {
		status = con
	} else {
		c.jsonResult(enums.JRCodeFailed, "状态码异常", 0)
	}
	for _, str := range strings.Split(mapp["ids"], ",") {
		if id, err := strconv.Atoi(str); err == nil {
			ids = append(ids, id)
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
		Model:"HashrateTreaty",
		Tbid:mapp["ids"],
		NewData:strconv.Itoa(status),
	}
	o:=orm.NewOrm()
	o.Insert(record)
	if len(ids) > 0 {
		if num, err := models.HashrateTreatyStatusUpdate(o,ids, status); err == nil {
			c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功"+msg+" %d 项", num), 0)
		} else {
			c.jsonResult(enums.JRCodeFailed, msg+"失败", 0)
		}
	} else {
		c.jsonResult(enums.JRCodeFailed, "缺少参数", 0)
	}
}
