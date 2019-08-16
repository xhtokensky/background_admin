package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

//订单提币表

func (a *FinancialOrderWithdrawal) TableName() string {
	return FinancialOrderWithdrawalTBName()
}

//查询的类
type FinancialOrderWithdrawalQueryParam struct {
	BaseQueryParam
	OrderId   string `json:"orderId"` //订单号
}

//理财收益表
type FinancialOrderWithdrawal struct {
	Id int `orm:"pk;column(id)"json:"id"form:"id"`
	//用户id
	OrderId string `orm:"column(order_id)"json:"orderId"form:"orderId"`
	//取出数量
	Quantity float64 `orm:"column(quantity)"json:"quantity"form:"quantity"`
	//取出时间
	WithdrawalTime int64 `orm:"column(withdrawal_time)"json:"withdrawalTime"form:"withdrawalTime"`
	//收益
	Profit float64 `orm:"column(profit)"json:"profit"form:"profit"`
	//年化收益
	YearProfit float64 `orm:"column(year_profit)"json:"yearProfit"form:"yearProfit"`
	//创建时间
	CreateTime time.Time `orm:"type(datetime);column(create_time)"json:"createTime"form:"createTime"`

}


//获取分页数据
func FinancialOrderWithdrawalPageList(params *FinancialOrderWithdrawalQueryParam) ([]*FinancialOrderWithdrawal, int64) {
	o := orm.NewOrm()
	query := o.QueryTable(FinancialOrderWithdrawalTBName())
	data := make([]*FinancialOrderWithdrawal, 0)
	//默认排序
	sortorder := "id"
	switch params.Sort {
	case "id":
		sortorder = "id"
	}
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}
	if params.OrderId != "" {
		query = query.Filter("order_id__exact", params.OrderId)
	}
	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, (params.Offset-1)*params.Limit).All(&data)
	return data, total
}