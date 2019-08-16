package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

//提币配置表

type TokenskyTibiConfigQueryParam struct {
	BaseQueryParam
	StartTime string `json:"startTime"` //开始时间
	EndTime   string `json:"endTime"`   //截止时间

}

func (a *TokenskyTibiConfig) TableName() string {
	return TokenskyTibiConfigTBName()
}

//提币配置表
type TokenskyTibiConfig struct {
	Id       int    `orm:"pk;column(id)"json:"id"form:"id"`
	CoinType string `orm:"column(coin_type)"json:"coinType"form:"coinType"`
	//单次最小提币数量
	Min float64 `orm:"column(min)"json:"min"form:"min"`
	//单次最多提币数量
	Max float64 `orm:"column(max)"json:"max"form:"max"`
	//当日限额
	CurDayQuantity float64 `orm:"column(cur_day_quantity)"json:"curDayQuantity"form:"curDayQuantity"`
	//百分比手续费
	ServiceCharge float64 `orm:"column(service_charge)"json:"serviceCharge"form:"serviceCharge"`
	//基础手续费
	BaseServiceCharge float64   `orm:"column(base_service_charge)"json:"baseServiceCharge"form:"baseServiceCharge"`
	AdminId           int       `orm:"column(admin_id)"json:"-"form:"-"`
	Status            int       `orm:"column(status)"json:"status"form:"status"`
	CreateTime        time.Time `orm:"auto_now_add;type(datetime);column(create_time)"json:"createTime"form:"createTime"`
}

func (a *TokenskyTibiConfigBak) TableName() string {
	return TokenskyTibiConfigBakTBName()

}

type TokenskyTibiConfigBak struct {
	TokenskyTibiConfig
}

//获取分页数据
func TokenskyTibiConfigPageList(params *TokenskyTibiConfigQueryParam) ([]*TokenskyTibiConfig, int64) {
	o := orm.NewOrm()
	query := o.QueryTable(TokenskyTibiConfigTBName())
	data := make([]*TokenskyTibiConfig, 0)
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
	if total > 0 {
		query.OrderBy(sortorder).Limit(params.Limit, (params.Offset-1)*params.Limit).All(&data)
	}
	return data, total
}

//获取最后一条数据 根据货币类型获取单条
func TokenskyTibiConfigGetLastOne(coinType string) *TokenskyTibiConfig {
	o :=orm.NewOrm()
	query := o.QueryTable(TokenskyTibiConfigTBName())
	var obj TokenskyTibiConfig
	if err := query.Filter("coin_type__iexact", coinType).One(&obj); err == nil {
		return &obj
	} else {
		return nil
	}
}
