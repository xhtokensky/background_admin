package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

// TableName 设置OctConf表名
func (a *OtcConf) TableName() string {
	return OtcConfTBName()
}

// OtcConfQueryParam 用于查询的类
type OtcConfQueryParam struct {
	BaseQueryParam
	Phone     string `json:"phone"`     //手机号 模糊查询
	StartTime string `json:"startTime"` //开始时间
	EndTime   string `json:"endTime"`   //截止时间

}

// OtcConf 实体类
type OtcConf struct {
	Id int `orm:"pk;column(id)"json:"id"form:"id"`
	//委托单最大交易额度
	OrdersMaxQuota float64 `orm:"column(orders_max_quota)"json:"ordersMaxQuota"form:"ordersMaxQuota"`
	//委托单最小交易额度
	OrdersMinQuota float64 `orm:"column(orders_min_quota)"json:"ordersMinQuota"form:"ordersMinQuota"`
	//买方手续费
	BuyerCost float64 `orm:"column(buyer_cost)"json:"buyerCost"form:"buyerCost"`
	//卖方手续费
	SellerCost float64 `orm:"column(seller_cost)"json:"sellerCost"form:"sellerCost"`
	//取消次数
	CancelNumber int `orm:"column(cancel_number)"json:"cancelNumber"form:"cancelNumber"`
	//买方超时时间
	BuyerOvertime int `orm:"column(buyer_overtime)"json:"buyerOvertime"form:"buyerOvertime"`
	//卖方超时时间
	SellerOvertime int `orm:"column(seller_overtime)"json:"sellerOvertime"form:"sellerOvertime"`
	//卖家申诉超时时间[卖家发起申诉]
	AppealOvertime int `orm:"column(appeal_overtime)"json:"appealOvertime"form:"appealOvertime"`
	//买家申诉超时时间
	VendeeAppealOvertime int `orm:"column(vendee_appeal_overtime)"json:"vendeeAppealOvertime"form:"vendeeAppealOvertime"`
	//委托挂单时间
	OrderEntryTime int `orm:"column(order_entry_time)"json:"orderEntryTime"form:"orderEntryTime"`
	//描述说明
	Msg string `orm:"column(msg)"json:"msg"form:"msg"`
	//交易规则
	Content string `orm:"column(content)"json:"content"form:"content"`
	//修改人uid
	UserId int `orm:"column(user_id)"json:"userId"form:"userId"`
	//创建时间
	CreateTime time.Time `orm:"auto_now_add;type(datetime);column(create_time)"json:"createTime"form:"createTime"`
}

func (a *OtcConfBak) TableName() string {
	return OtcConfBakTBName()
}

//副表 只做记录
type OtcConfBak struct {
	Id int `orm:"pk;column(id)"json:"id"form:"id"`
	//委托单最大交易额度
	OrdersMaxQuota float64 `orm:"column(orders_max_quota)"json:"ordersMaxQuota"form:"ordersMaxQuota"`
	//委托单最小交易额度
	OrdersMinQuota float64 `orm:"column(orders_min_quota)"json:"ordersMinQuota"form:"ordersMinQuota"`
	//买方手续费
	BuyerCost float32 `orm:"column(buyer_cost)"json:"buyerCost"form:"buyerCost"`
	//卖方手续费
	SellerCost float32 `orm:"column(seller_cost)"json:"sellerCost"form:"sellerCost"`
	//取消次数
	CancelNumber int `orm:"column(cancel_number)"json:"cancelNumber"form:"cancelNumber"`
	//买方超时时间
	BuyerOvertime int `orm:"column(buyer_overtime)"json:"buyerOvertime"form:"buyerOvertime"`
	//卖方超时时间
	SellerOvertime int `orm:"column(seller_overtime)"json:"sellerOvertime"form:"sellerOvertime"`
	//卖家申诉超时时间
	AppealOvertime int `orm:"column(appeal_overtime)"json:"appealOvertime"form:"appealOvertime"`
	//买家申诉超时时间
	VendeeAppealOvertime int `orm:"column(vendee_appeal_overtime)"json:"vendeeAppealOvertime"form:"vendeeAppealOvertime"`
	//委托挂单时间
	OrderEntryTime int `orm:"column(order_entry_time)"json:"orderEntryTime"form:"orderEntryTime"`
	//描述说明
	Msg string `orm:"column(msg)"json:"msg"form:"msg"`
	//交易规则
	Content string `orm:"column(content)"json:"content"form:"content"`
	//修改人uid
	UserId int `orm:"column(user_id)"json:"userId"form:"userId"`
	//创建时间
	CreateTime time.Time `orm:"auto_now_add;type(datetime);column(create_time)"json:"createTime"form:"createTime"`
}

// OtcConfPageList 获取分页数据
func OtcConfPageList(params *OtcConfQueryParam) ([]*OtcConf, int64) {
	o := orm.NewOrm()
	query := o.QueryTable(OtcConfBakTBName())
	data := make([]*OtcConf, 0)
	//默认排序
	sortorder := "id"
	switch params.Sort {
	case "id":
		sortorder = "id"
	}
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}
	total, _ := query.Count()
	if total > 0 {
		query.OrderBy(sortorder).Limit(params.Limit, (params.Offset-1)*params.Limit).All(&data)
	}
	return data, total
}

//获取最后一条数据
func OtcConfGetLastOne() OtcConf {
	o := orm.NewOrm()
	query := o.QueryTable(OtcConfTBName())
	var obj OtcConf
	query.OrderBy("-id").One(&obj)
	return obj
}
