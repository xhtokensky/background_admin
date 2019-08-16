package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

func (a *SpiderCoinMarket) TableName() string {
	return SpiderCoinMarketTBName()
}

//货币行情数据
type SpiderCoinMarket struct {
	Coin string `orm:"pk;column(coin)"json:"coin"form:"coin"`
	//区块儿奖励
	BlockReward string `json:"block_reward"`
	//出块儿时间[秒]
	BlockTime int `json:"block_time"`
	CoinPrice string `json:"coin_price"`
	CurrConnections int `json:"curr_connections"`
	CurrDiff string `json:"curr_diff"`
	HashUnit string `json:"hash_unit"`
	MinPaymentAmount string `json:"min_payment_amount"`
	MiningAlgorithm string `json:"mining_algorithm"`
	NetworkHashrate string `json:"network_hashrate"`
	PaymentEndTime string `json:"payment_end_time"`
	PaymentStartTime string `json:"payment_start_time"`
	PoolHashrate string `json:"pool_hashrate"`
	PricingCurrency string `json:"pricing_currency"`
	PricingCurrencySymbol string `json:"pricing_currency_symbol"`
	//每T收益
	UnitOutput string `json:"unit_output"`
	UnitOutputCurrency string `json:"unit_output_currency"`
	//更新时间
	IsDate time.Time `json:"is_date"`
}

//获取单条
func SpiderCoinMarketOne(coin string)*SpiderCoinMarket{
	o := orm.NewOrm()
	obj := &SpiderCoinMarket{}
	query := o.QueryTable(SpiderCoinMarketTBName())
	err := query.Filter("coin__exact",coin).One(obj)
	if err != nil{
		return nil
	}
	return obj
}