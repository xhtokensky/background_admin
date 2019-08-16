package models

import (
	"time"
)

//用户资产变化记录
func (a *TokenskyUserBalancesRecord) TableName() string {
	return TokenskyUserBalancesRecordTBName()
}

type TokenskyUserBalancesRecord struct {
	Id  int `orm:"pk;column(id)"json:"id"form:"id"`
	User *TokenskyUser `orm:"rel(fk);column(user_id)"json:"user"form:"user"`

	//资产类型
	Symbol string `orm:"column(symbol)"json:"symbol"form:"symbol"`
	//资产值
	OldBalance float64 `orm:"column(old_balance)"json:"oldBalance"form:"oldBalance"`
	OldFrozenBalance float64   `orm:"column(old_frozen_balance)"json:"oldFrozenBalance"form:"oldFrozenBalance"`
	//
	NewBalance float64 `orm:"column(new_balance)"json:"newBalance"form:"newBalance"`
	NewFrozenBalance float64   `orm:"column(new_frozen_balance)"json:"newFrozenBalance"form:"newFrozenBalance"`
	CreateTime    time.Time `orm:"auto_now_add;type(datetime);column(create_time)"json:"createTime"form:"createTime"`
	//说明
	Cont string `orm:"column(cont)"json:"cont"form:"cont"`
	//来源  1后端 2后台 3定时服务
	Source int `orm:"column(source)"json:"source"form:"source"`
	//时间[单位毫秒]
	PushTime int64 `orm:"column(push_time)"json:"pushTime"form:"pushTime"`
	//模块
	Mold string `orm:"column(mold)"json:"mold"form:"mold"`
	//标示id
	SignId string `orm:"column(sign_id)"json:"signId"form:"signId"`
	//哈希
	HashId string `orm:"column(hash_id)"json:"hashId"form:"hashId"`
	//操作资产方法
	MethodBalance string `orm:"column(method_balance)"json:"methodBalance"form:"methodBalance"`
	//操作资产数量
	Balance string `orm:"column(balance)"json:"balance"form:"balance"`
	//操作冻结资产方法
	MethodFrozenBalance string `orm:"column(method_frozen_balance)"json:"methodFrozenBalance"form:"methodFrozenBalance"`
	//操作冻结资产数量
	FrozenBalance string `orm:"column(frozen_balance)"json:"frozenBalance"form:"frozenBalance"`
}