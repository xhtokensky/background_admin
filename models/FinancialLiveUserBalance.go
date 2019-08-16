package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

/*
CREATE TABLE `financial_live_user_balance` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `symbol` varchar(50) NOT NULL,
  `balance` double(255,8) NOT NULL,
  `push_time` int(11) NOT NULL COMMENT '创建时间  时间戳',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='理财活期资产表'
*/

//查询类
type FinancialLiveUserBalanceQueryParam struct {
	BaseQueryParam
	StartTime int64 `json:"startTime"` //开始时间
	EndTime   int64 `json:"endTime"`   //截止时间
	UserId    int `json:"userId"`    //用户id
}

func (a *FinancialLiveUserBalance) TableName() string {
	return FinancialLiveUserBalanceTBName()
}

//理财活期资产表
type FinancialLiveUserBalance struct {
	Id      int           `orm:"pk;column(id)"json:"id"form:"id"`
	User    *TokenskyUser `orm:"rel(fk);column(user_id)"json:"userId"form:"userId"`
	Symbol  string        `orm:"column(symbol)"json:"symbol"form:"symbol"`
	Balance float64       `orm:"column(balance)"json:"balance"form:"balance"`
	//推送时间
	PushTime int64 `orm:"column(push_time)"json:"pushTime"form:"pushTime"`
	//创建时间
	CreateTime time.Time `orm:"auto_now_add;type(datetime);column(create_time)"json:"createTime"form:"createTime"`
}

func FinancialLiveUserBalancePageList(params *FinancialLiveUserBalanceQueryParam) ([]*FinancialLiveUserBalance, int64) {
	data := make([]*FinancialLiveUserBalance, 0)
	o := orm.NewOrm()
	query := o.QueryTable(FinancialLiveUserBalanceTBName())
	//默认排序
	sortorder := "id"
	switch params.Sort {
	case "id":
		sortorder = "id"
	}
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}
	if params.UserId >0{
		query = query.Filter("User__user_id__exact",params.UserId)
	}
	query.OrderBy(sortorder).Limit(params.Limit, (params.Offset-1)*params.Limit).RelatedSel().All(&data)
	return data, 0
}

//获取用户资产记录
func  FinancialLiveUserBalanceOne(o orm.Ormer,uid int,symbol string)(*FinancialLiveUserBalance){
	obj := &FinancialLiveUserBalance{}
	query := o.QueryTable(FinancialLiveUserBalanceTBName())
	if err := query.Filter("User__user_id__exact",uid).Filter("symbol__exact",symbol).One(obj);err!=nil{
		return nil
	}
	return obj
}