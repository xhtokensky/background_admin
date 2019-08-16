package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

//查询的类
type HashrateOrderTransactionQueryParam struct {
	BaseQueryParam
	Name      string `json:"name"`
	StartTime int64  `json:"startTime"` //开始时间
	EndTime   int64  `json:"endTime"`   //截止时间
	Status    string `json:"status"`    //状态
}

func (a *HashrateOrderTransaction) TableName() string {
	return HashrateOrderTransactionTBName()
}

//算力合约关联表
type HashrateOrderTransaction struct {
	KeyId int `orm:"pk;column(key_id)"json:"keyId"form:"keyId"`
	//订单ID
	OrderId string `orm:"column(order_id)"json:"orderId"form:"orderId"`
	//支付方式 货币
	PayType string `orm:"column(pay_type)"json:"payType"form:"payType"`
	//交易金额
	TransactionMoney float64 `orm:"column(transaction_money)"json:"transactionMoney"form:"transactionMoney"`
	//创建时间
	CreateTime time.Time `orm:"auto_now_add;type(datetime);column(create_time)"json:"createTime"form:"createTime"`
}

func HashrateOrderTransactionsByIds(ids []string) map[string]map[string]*HashrateOrderTransaction {
	mapp := make(map[string]map[string]*HashrateOrderTransaction)
	data := make([]*HashrateOrderTransaction, 0)
	if len(ids) > 0 {
		query := orm.NewOrm().QueryTable(HashrateOrderTransactionTBName())
		query = query.Filter("order_id__in", ids)
		query.All(&data)
	}
	for _, obj := range data {
		if _, found := mapp[obj.OrderId]; !found {
			mapp[obj.OrderId] = make(map[string]*HashrateOrderTransaction)
		}
		mapp[obj.OrderId][obj.PayType] = obj
	}
	return mapp
}
