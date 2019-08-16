package models

import "time"

func (a *OtcEntrustAutoCancelRecord) TableName() string {
	return OtcEntrustAutoCancelRecordTBName()
}

//订单自动取消记录表
type OtcEntrustAutoCancelRecord struct {
	KeyId int `orm:"pk;column(key_id)"json:"keyId"form:"keyId"`
	//订单类型
	EntrustType int `orm:"column(entrust_type)"json:"entrustType"form:"entrustType"`
	//订单ID
	EntrustOrderId int     `orm:"column(entrust_order_id)"json:"entrustOrderId"form:"entrustOrderId"`
	Money          float64 `orm:"column(money)"json:"money"form:"money"`
	//手续费
	ServiceCharge float64 `orm:"column(service_charge)"json:"serviceCharge"form:"serviceCharge"`
	//总金额
	SumMoney float64 `orm:"column(sum_money)"json:"sumMoney"form:"sumMoney"`
	UserId   int     `orm:"column(user_id)"json:"userId"form:"userId"`
	//创建时间
	CreateTime time.Time `orm:"auto_now_add;type(datetime);column(create_time)"json:"createTime"form:"createTime"`
}
