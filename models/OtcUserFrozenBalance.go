package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

//查询的类
type OtcUserFrozenBalanceQueryParam struct {
	BaseQueryParam
}

func (a *OtcUserFrozenBalance) TableName() string {
	return OtcUserFrozenBalanceTBName()
}

//卖出委托订单表
type OtcUserFrozenBalance struct {
	KeyId int `orm:"pk;column(key_id)"json:"keyId"form:"keyId"`
	//1委托卖出 2订单卖出'
	Type   int `orm:"column(type)"json:"type"form:"type"`
	UserId int `orm:"column(user_id)"json:"userId"form:"userId"`
	//关联ID
	RelevanceId string `orm:"column(relevance_id)"json:"relevanceId"form:"relevanceId"`
	//冻结金额 包含手续费
	FrozenBalance float64 `orm:"column(frozen_balance)"json:"frozenBalance"form:"frozenBalance"`
	//手续费
	ServiceChargeBalance float64 `orm:"column(service_charge_balance)"json:"serviceChargeBalance"form:"serviceChargeBalance"`
	//
	Status int `orm:"column(status)"json:"status"form:"status"`
	//创建时间
	CreateTime time.Time `orm:"type(datetime);column(create_time)"json:"createTime"form:"createTime"`
}

func OtcUserFrozenBalanceOneByRelevanceIdAndStatusAndType(stype, status int, relevanceId string) *OtcUserFrozenBalance {
	obj := OtcUserFrozenBalance{}
	query := orm.NewOrm().QueryTable(OtcUserFrozenBalanceTBName())
	query = query.Filter("relevance_id__exact", relevanceId)
	query = query.Filter("type__exact", stype)
	query = query.Filter("status__exact", status)
	err := query.One(&obj)
	if err != nil {
		return nil
	}
	return &obj
}

func OtcUserFrozenBalanceOneByRelevanceIdAndStatusAndType2(o orm.Ormer, stype, status int, relevanceId string) *OtcUserFrozenBalance {
	obj := OtcUserFrozenBalance{}
	query := o.QueryTable(OtcUserFrozenBalanceTBName())
	query = query.Filter("relevance_id__exact", relevanceId)
	query = query.Filter("type__exact", stype)
	query = query.Filter("status__exact", status)
	err := query.One(&obj)
	if err != nil {
		return nil
	}
	return &obj
}

func OtcUserFrozenBalanceOneByRelevanceIdAndType2(o orm.Ormer, relevanceId string, Type int) *OtcUserFrozenBalance {
	m := OtcUserFrozenBalance{}
	err := o.QueryTable(OtcUserFrozenBalanceTBName()).Filter("relevance_id__exact", relevanceId).Filter("type__exact", Type).One(&m)
	if err != nil {
		return nil
	}
	return &m
}
