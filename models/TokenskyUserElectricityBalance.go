package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

//用户电力资产表
func (a *TokenskyUserElectricityBalance) TableName() string {
	return TokenskyUserElectricityBalanceTBName()
}

type TokenskyUserElectricityBalance struct {
	KeyId  int `orm:"pk;column(key_id)"json:"keyId"form:"keyId"`
	UserId int `orm:"column(user_id)"json:"userId"form:"userId"`
	//资产类型
	CoinType string `orm:"column(coin_type)"json:"coinType"form:"coinType"`
	//资产值
	Balance    float64   `orm:"column(balance)"json:"frozenBalance"form:"frozenBalance"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime);column(create_time)"json:"createTime"form:"createTime"`
}

//获取算力资产
func TokenskyUserElectricityBalanceByUids(uids []int) map[int]*TokenskyUserElectricityBalance {
	o := orm.NewOrm()
	query := o.QueryTable(TokenskyUserElectricityBalanceTBName())
	data := make([]*TokenskyUserElectricityBalance, 0)
	query.Filter("user_id__in", uids).All(&data)
	mapp := make(map[int]*TokenskyUserElectricityBalance)
	for _, obj := range data {
		mapp[obj.UserId] = obj
	}
	return mapp
}

func TokenskyUserElectricityBalanceByUid(o orm.Ormer, uid int) *TokenskyUserElectricityBalance {
	query := o.QueryTable(TokenskyUserElectricityBalanceTBName())
	obj := &TokenskyUserElectricityBalance{}
	if err := query.Filter("user_id__exact", uid).One(obj); err != nil {
		return nil
	}
	return obj
}
