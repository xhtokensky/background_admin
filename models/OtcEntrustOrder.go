package models

import (
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"time"
	"tokensky_bg_admin/conf"
)

// OtcEntrustOrderQueryParam 用于查询的类
type OtcEntrustOrderQueryParam struct {
	BaseQueryParam
	Phone string `json:"phone"` //手机号 模糊查询
	//发布时间
	StartTime int64  `json:"startTime"` //开始时间
	EndTime   int64  `json:"endTime"`   //截止时间
	Status    string `json:"status"`    //状态
}

func (a *OtcEntrustOrder) TableName() string {
	return OtcEntrustOrderTBName()
}

//委托单管理
type OtcEntrustOrder struct {
	KeyId int `orm:"pk;column(key_id)"json:"keyId"form:"keyId"`
	//委托类型 1买单 2卖单
	EntrustType int `orm:"column(entrust_type)"json:"entrustType"form:"entrustType"`
	//货币类型
	CoinType string `orm:"column(coin_type)"json:"coinType"form:"coinType"`
	//单价
	UnitPrice float64 `orm:"column(unit_price)"json:"unitPrice"form:"unitPrice"`
	//金钱单位
	MoneyType string `orm:"column(money_type)"json:"moneyType"form:"moneyType"`

	//数量
	Quantity float64 `orm:"column(quantity)"json:"quantity"form:"quantity"`
	//剩下数量
	QuantityLeft float64 `orm:"column(quantity_left)"json:"quantityLeft"form:"quantityLeft"`

	//最小交易额
	Min float64 `orm:"column(min)"json:"min"form:"min"`
	//最大交易额
	Max float64 `orm:"column(max)"json:"max"form:"max"`
	//买方手续费
	VendeeServiceCharge float64 `orm:"column(vendee_service_charge)"json:"vendeeServiceCharge"form:"vendeeServiceCharge"`
	//卖方手续费
	VendorServiceCharge float64 `orm:"column(vendor_service_charge)"json:"vendorServiceCharge"form:"vendorServiceCharge"`
	//支付方式  1支付宝 2微信 3银行卡
	PayType string `orm:"column(pay_type)"json:"payType"form:"payType"`
	//支付方式2
	PayTypeList []int `orm:"-"json:"payTypeList"form:"payTypeList"`
	//完成时间
	FinishTime time.Time `orm:"type(datetime);column(finish_time)"json:"finishTime"form:"finishTime"`
	//发布时间
	PushTime time.Time `orm:"type(datetime);column(push_time)"json:"pushTime"form:"pushTime"`
	//状态 1发布中 2已完成 0已取消 3系统自动取消
	Status int `orm:"column(status)"json:"status"form:"status"`
	//创建时间
	CretaeTime time.Time `orm:"type(datetime);column(cretae_time)"json:"cretaeTime"form:"cretaeTime"`
	//更新时间
	UpdateTime time.Time `orm:"auto_now;type(datetime);column(update_time)"json:"updateTime"form:"updateTime"`
	//自动完成时间
	AutoCancelTime time.Time `orm:"type(datetime);column(auto_cancel_time)"json:"autoCancelTime"form:"autoCancelTime"`
	//手机号
	Phone    string `orm:"-"json:"phone"form:"phone"`
	NickName string `orm:"-"json:"-"form:"nickName"`
	//用户Uid 连表用户 一对多关系
	User *TokenskyUser `orm:"rel(fk)"json:"-"form:"-"`
}

// OtcEntrustOrderPageList 获取分页数据
func OtcEntrustOrderPageList(params *OtcEntrustOrderQueryParam) ([]*OtcEntrustOrder, int64) {
	o := orm.NewOrm()
	query := o.QueryTable(OtcEntrustOrderTBName())
	data := make([]*OtcEntrustOrder, 0)
	var total int64
	//默认排序
	sortorder := "key_id"
	switch params.Sort {
	case "key_id":
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
	//状态
	if params.Status != "" && params.Status != "-1" {
		query = query.Filter("status__iexact", params.Status)
	}
	//时间段
	if params.StartTime > 0 {
		query = query.Filter("push_time__gte", time.Unix(params.StartTime, 0))
	}
	if params.EndTime > 0 {
		query = query.Filter("push_time__lte", time.Unix(params.EndTime, 0))
	}
	//电话查询
	if params.Phone != "" {
		query = query.Filter("User__Phone__icontains", params.Phone)
	}
	total, _ = query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, (params.Offset-1)*params.Limit).RelatedSel().All(&data)
	for _, obj := range data {
		//电话
		if obj.User != nil {
			obj.Phone = obj.User.Phone
			obj.NickName = obj.User.NickName
		}
		//支付方式
		payList := strings.Split(obj.PayType, ",")
		obj.PayTypeList = make([]int, conf.PAY_TYPE_MAX_NUM)
		for _, str := range payList {
			if con, err := strconv.Atoi(str); err == nil {
				if con > 0 && con <= conf.PAY_TYPE_MAX_NUM {
					obj.PayTypeList[con-1] = 1
				}
			}
		}
	}
	return data, total
}

// OtcEntrustOrder 获取单条
func OtcEntrustOrderOneByKid(kid int) *OtcEntrustOrder {
	m := OtcEntrustOrder{}
	o := orm.NewOrm().QueryTable(OtcEntrustOrderTBName())
	query := o.Filter("key_id__exact", kid)
	err := query.One(&m)
	if err != nil {
		return nil
	}
	return &m
}

func OtcEntrustOrderOneByKid2(o orm.Ormer, kid int) *OtcEntrustOrder {
	m := OtcEntrustOrder{}
	query := o.QueryTable(OtcEntrustOrderTBName())
	query = query.Filter("key_id__exact", kid)
	err := query.One(&m)
	if err != nil {
		return nil
	}
	return &m
}
