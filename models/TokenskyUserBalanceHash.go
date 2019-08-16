package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

func (a *TokenskyUserBalanceHash) TableName() string {
	return TokenskyUserBalanceHashTBName()
}


//用户资产哈希表
type TokenskyUserBalanceHash struct {
	HashId string  `orm:"pk;column(hash_id)"json:"hashId"form:"hashId"`
	//资产服务处理状态
	BalanceStatus int `orm:"column(balance_status)"json:"balanceStatus"form:"balanceStatus"`
	//来源
	Source int `orm:"column(source)"json:"source"form:"source"`
	//状态
	ModelStatus int `orm:"column(model_status)"json:"modelStatus"form:"modelStatus"`
	//创建时间
	CreateTime time.Time `orm:"auto_now_add;type(datetime);column(create_time)"json:"createTime"form:"createTime"`
}

func TokenskyUserBalanceHashOne(o orm.Ormer,hashId string)*TokenskyUserBalanceHash {
	query := o.QueryTable(TokenskyUserBalanceHashTBName())
	obj := &TokenskyUserBalanceHash{}
	err := query.Filter("hash_id__exact",hashId).One(obj)
	if err !=nil{
		return nil
	}
	return obj
}

func TokenskyUserBalanceHashSetStatus(o orm.Ormer,hashId string)bool {
	o2 := orm.NewOrm()
	obj := &TokenskyUserBalanceHash{
		HashId:hashId,
	}
	err := o2.Read(obj)
	if err != nil{
		return false
	}
	obj.ModelStatus = 1
	_,err = o.InsertOrUpdate(obj)
	if err != nil{
		return false
	}
	return true
}