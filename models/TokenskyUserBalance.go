package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

//查询的类
type TokenskyUserBalanceRecordParam struct {
	BaseQueryParam
	Uid int `json:"uid"`
}

func (a *TokenskyUserBalance) TableName() string {
	return TokenskyUserBalanceTBName()
}

//用户资产表
type TokenskyUserBalance struct {
	KeyId  int `orm:"pk;column(key_id)"json:"keyId"form:"keyId"`
	UserId int `orm:"column(user_id)"json:"userId"form:"userId"`
	//资产类型
	CoinType string `orm:"column(coin_type)"json:"coinType"form:"coinType"`
	//资产值
	Balance float64 `orm:"column(balance)"json:"frozenBalance"form:"frozenBalance"`
	//冻结资产
	FrozenBalance float64   `orm:"column(frozen_balance)"json:"freezeCoin"form:"freezeCoin"`
	CreateTime    time.Time `orm:"auto_now_add;type(datetime);column(create_time)"json:"createTime"form:"createTime"`
}

func GetTokenskyUserBalanceByUidCoinType2(o orm.Ormer, uid int, coinType string) *TokenskyUserBalance {
	//o := orm.NewOrm().QueryTable(TokenskyUserBalanceTBName())
	m := TokenskyUserBalance{}
	query := o.QueryTable(TokenskyUserBalanceTBName())
	query = query.Filter("user_id__exact", uid)
	query = query.Filter("coin_type__exact", coinType)
	err := query.One(&m)
	if err != nil {
		return nil
	}
	return &m
}

//获取用户所有资产
func GetTokenskyUserBalancesByUid(uid int) ([]*TokenskyUserBalance, int64) {
	o := orm.NewOrm()
	query := o.QueryTable(TokenskyUserBalanceTBName())
	query = query.Filter("user_id__iexact", uid)
	total, _ := query.Count()
	data := make([]*TokenskyUserBalance, 0)
	query.All(&data)
	return data, total
}

//循环用户
func TokenskyUserBalancesIteration(num int) func() []*TokenskyUserBalance {
	query := orm.NewOrm().QueryTable(TokenskyUserBalanceTBName())
	page := 1
	return func() []*TokenskyUserBalance {
		data := make([]*TokenskyUserBalance, 0)
		query.Limit(num, (page-1)*num).All(&data, "user_id", "coin_type")
		page++
		return data
	}
}
func TokenskyUserBalancesAllIteration(num int) func() []*TokenskyUserBalance {
	query := orm.NewOrm().QueryTable(TokenskyUserBalanceTBName())
	page := 1
	return func() []*TokenskyUserBalance {
		data := make([]*TokenskyUserBalance, 0)
		query.Limit(num, (page-1)*num).All(&data)
		page++
		return data
	}
}