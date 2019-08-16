package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type TokenskyChongbiConfigQueryParam struct {
	BaseQueryParam
	StartTime string `json:"startTime"` //开始时间
	EndTime   string `json:"endTime"`   //截止时间

}

func (a *TokenskyChongbiConfig) TableName() string {
	return TokenskyChongbiConfigTBName()
}

//冲币配置
type TokenskyChongbiConfig struct {
	Id int `orm:"pk;column(id)"json:"id"form:"id"`
	//货币类型
	CoinType string `orm:"column(coin_type)"json:"coinType"form:"coinType"`
	//最小充币值
	Min     float64 `orm:"column(min)"json:"min"form:"min"`
	AdminId int     `orm:"column(admin_id)"json:"-"form:"-"`
	//百分比手续费
	ServiceCharge float64 `orm:"column(service_charge)"json:"serviceCharge"form:"serviceCharge"`
	//基础手续费
	BaseServiceCharge float64 `orm:"column(base_service_charge)"json:"baseServiceCharge"form:"baseServiceCharge"`
	//状态[暂未使用]
	Status     int       `orm:"column(status)"json:"status"form:"status"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime);column(create_time)"json:"createTime"form:"createTime"`
}

func (a *TokenskyChongbiConfigBak) TableName() string {
	return TokenskyChongbiConfigBakTBName()

}

type TokenskyChongbiConfigBak struct {
	TokenskyChongbiConfig
}

//获取分页数据
func TokenskyChongbiConfigPageList(params *TokenskyChongbiConfigQueryParam) ([]*TokenskyChongbiConfig, int64) {
	o := orm.NewOrm()
	query := o.QueryTable(TokenskyChongbiConfigTBName())
	data := make([]*TokenskyChongbiConfig, 0)
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
	query.OrderBy(sortorder).Limit(params.Limit, (params.Offset-1)*params.Limit).All(&data)
	return data, total
}

//获取最后一条数据 根据货币类型获取单条
func TokenskyChongbiConfigGetLastOne(coinType string) *TokenskyChongbiConfig {
	o :=orm.NewOrm()
	query := o.QueryTable(TokenskyChongbiConfigTBName())
	var obj TokenskyChongbiConfig
	if err := query.Filter("coin_type__iexact", coinType).One(&obj); err == nil {
		return &obj
	} else {
		return nil
	}
}
