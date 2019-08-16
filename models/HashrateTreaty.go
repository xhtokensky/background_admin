package models

import (
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"time"
	"tokensky_bg_admin/utils"
)

//查询的类
type HashrateTreatyQueryParam struct {
	BaseQueryParam
	Title     string `json:"title"`
	StartTime int64  `json:"startTime"` //开始时间
	EndTime   int64  `json:"endTime"`   //截止时间
	Status    string `json:"status"`    //状态
}

func (a *HashrateTreaty) TableName() string {
	return HashrateTreatyTBName()
}

//算力合约表
type HashrateTreaty struct {
	KeyId int `orm:"pk;column(key_id)"json:"keyId"form:"keyId"`
	//标题
	Title string `orm:"column(title)"json:"title"form:"title"`
	//价格
	Price float64 `orm:"column(price)"json:"price"form:"price"`
	//现价
	OriginalPrice float64 `orm:"column(original_price)"json:"originalPrice"form:"originalPrice"`
	//标签 0是普通版 1优惠版
	Tag int `orm:"column(tag)"json:"tag"form:"tag"`
	//库存
	Inventory int `orm:"column(inventory)"json:"inventory"form:"inventory"`
	//电费
	ElectricBill float64 `orm:"column(electric_bill)"json:"electricBill"form:"electricBill"`
	//剩余库存
	InventoryLeft int `orm:"column(inventory_left)"json:"inventoryLeft"form:"inventoryLeft"`
	//运行周期
	RunCycle int `orm:"column(run_cycle)"json:"runCycle"form:"runCycle"`
	//排序
	Sort int `orm:"column(sort)"json:"sort"form:"sort"`
	//状态 0待上架 1上架 2后台下架(不可重新上架)
	Status int `orm:"column(status)"json:"status"form:"status"`
	//限购 0无限制
	Restriction int `orm:"column(restriction)"json:"restriction"form:"restriction"`
	//简介
	Intro string `orm:"column(intro)"json:"intro"form:"intro"`
	//管理费
	Management float64 `orm:"column(management)"json:"management"form:"management"`
	//预计收益
	EarningsRate float64 `orm:"column(earnings_rate)"json:"earningsRate"form:"earningsRate"`
	//开挖时间
	//ExcavateTime time.Time `orm:"column(excavate_time)"json:"excavateTime"form:"excavateTime"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime);column(create_time)"json:"createTime"form:"createTime"`
	UpdateTime time.Time `orm:"auto_now;type(datetime);column(update_time)"json:"updateTime"form:"updateTime"`
	//货币
	UnitMoney string `orm:"column(unit_money)"json:"unitMoney"form:"unitMoney"`
	//创建人id
	AdminUserId int `orm:"column(admin_user_id)"json:"-"form:"-"`
	//算力合约分类表分类
	HashrateCategoryObj  *HashrateCategory `orm:"rel(fk);column(hashrate_category)"json:"-"form:"-"`
	ImgUrl               string            `orm:"-"json:"imgUrl"` //合约类型表
	HashrateCategory     int               `orm:"-"json:"hashrateCategory"form:"hashrateCategory"`
	HashrateCategoryName string            `orm:"-"json:"hashrateCategoryName"form:"hashrateCategoryName"`
	//连表 算力订单表 一对多反向关系
	HashrateOrders []*HashrateOrder `orm:"reverse(many)"json:"-"form:"-"`

	//期货类型 0非期货 1期货类
	FuturesType int `orm:"column(futures_type)"json:"futuresType"form:"futuresType"`
	//生效日期
	EffectiveDate  time.Time `orm:"column(effective_date)"json:"effectiveDate"form:"effectiveDate"`
	EffectiveDateTm int64     `orm:"-"json:"effectiveDateTm"form:"effectiveDateTm"`
}

func (this *HashrateTreaty)Json() string{
	js,err := json.Marshal(this)
	if err != nil{
		return "err"
	}
	return string(js)
}

//获取分页数据
func HashrateTreatyPageList(params *HashrateTreatyQueryParam) ([]*HashrateTreaty, int64) {
	o := orm.NewOrm()
	query := o.QueryTable(HashrateTreatyTBName())
	data := make([]*HashrateTreaty, 0)
	//默认排序
	sortorder := "key_id"
	switch params.Sort {
	case "keyId":
		sortorder = "key_id"
	}
	switch params.Order {
	case "":
		sortorder = "-" + sortorder
	case "desc":
		sortorder = "-" + sortorder
	default:
		sortorder = sortorder
	}
	//姓名模糊查询
	if params.Title != "" {
		query = query.Filter("title__icontains", params.Title)
	}
	//状态
	if params.Status != "" && params.Status != "-1" {
		query = query.Filter("status__iexact", params.Status)
	}
	//时间段
	if params.StartTime > 0 {
		query = query.Filter("create_time__gte", time.Unix(params.StartTime, 0))
	}
	if params.EndTime > 0 {
		query = query.Filter("create_time__lte", time.Unix(params.EndTime, 0))
	}
	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, (params.Offset-1)*params.Limit).RelatedSel().All(&data)

	deadline := time.Now().Add(time.Second * 3600).Unix()
	for _, obj := range data {
		//获取分类表
		if obj.HashrateCategoryObj != nil {
			obj.ImgUrl = utils.QiNiuDownload(obj.HashrateCategoryObj.ImgKey, deadline)
			obj.HashrateCategory = obj.HashrateCategoryObj.KeyId
			obj.HashrateCategoryName = obj.HashrateCategoryObj.Name
		}
		obj.EffectiveDateTm = obj.EffectiveDate.Unix()
	}
	return data, total
}

//批量删除 只有订单号为0可以编辑
func HashrateTreatyDelete(ids []int) (int64, error) {
	o := orm.NewOrm()
	query := o.QueryTable(HashrateTreatyTBName())
	num, err := query.Filter("status__exact", 0).Filter("key_id__in", ids).Delete()
	return num, err
}

//批量修改状态
func HashrateTreatyStatusUpdate(o orm.Ormer,ids []int, status int) (int64, error) {
	if len(ids) <= 0 {
		return 0, nil
	}
	query := o.QueryTable(HashrateTreatyTBName())
	params := map[string]interface{}{
		"status": status,
	}
	var num int64
	var err error
	switch status {
	case 1:
		//上架
		num, err = query.Filter("key_id__in", ids).Filter("status__exact", 0).Update(params)
	case 2:
		//后台下架
		num, err = query.Filter("key_id__in", ids).Filter("status__exact", 1).Update(params)
	default:
		return 0, err
	}
	return num, err
}

//获取指定
func HashrateTreatysByIds(ids []int) map[int]*HashrateTreaty {
	mapp := make(map[int]*HashrateTreaty)
	data := make([]*HashrateTreaty, 0)
	if len(ids) > 0 {
		query := orm.NewOrm().QueryTable(HashrateTreatyTBName())
		query = query.Filter("hashrate_category__in", ids)
		query.All(&data)
	}
	for _, obj := range data {
		mapp[obj.HashrateCategoryObj.KeyId] = obj
		obj.HashrateCategory = obj.HashrateCategoryObj.KeyId
		obj.HashrateCategoryName = obj.HashrateCategoryObj.Name
	}
	return mapp
}

//获取单条
func HashrateTreatyOneById(id int) *HashrateTreaty {
	obj := HashrateTreaty{}
	o := orm.NewOrm().QueryTable(HashrateTreatyTBName())
	if err := o.Filter("key_id__exact", id).RelatedSel().One(&obj); err != nil {
		return nil
	}
	return &obj
}
