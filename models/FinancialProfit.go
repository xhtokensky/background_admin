package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

func (a *FinancialProfit) TableName() string {
	return FinancialProfitTBName()
}

//查询的类
type FinancialProfitQueryParam struct {
	BaseQueryParam
	Name      string `json:"name"`
	StartTime int64  `json:"startTime"` //开始时间
	EndTime   int64  `json:"endTime"`   //截止时间
	OrderId   string `json:"orderId"` //订单号
	Status    string `json:"status"`    //状态
	Category string `json:"category"`
}

//理财收益表
type FinancialProfit struct {
	Id int `orm:"pk;column(id)"json:"id"form:"id"`
	//用户id
	User   *TokenskyUser `orm:"rel(fk);column(user_id)"json:"-"form:"-"`
	UserId int           `orm:"-"json:"userId"form:"userId"`
	//配置表id
	ProductId int `orm:"column(product_id)"json:"productId"form:"productId"`
	//类型 1活期 2定期
	Category int `orm:"column(category)"json:"category"form:"category"`
	//关联id
	RelevanceId int `orm:"column(relevance_id)"json:"relevanceId"form:"relevanceId"`
	//子类 live 活期  dead 定期 withdrawal 提前取出
	Product string  `orm:"column(product)"json:"product"form:"product"`
	//货币类x型
	Symbol string `orm:"column(symbol)"json:"symbol"form:"symbol"`
	//资产数量
	Balance float64 `orm:"column(balance)"json:"balance"form:"balance"`
	//支付资产
	PayBalance float64 `orm:"column(pay_balance)"json:"payBalance"form:"payBalance"`

	//年化利率
	YearProfit float64 `orm:"column(year_profit)"json:"yearProfit"form:"yearProfit"`
	//收益
	Profit float64 `orm:"column(profit)"json:"profit"form:"profit"`
	//收益日期
	IsDate time.Time `orm:"type(date);column(is_date)"json:"isDate"form:"isDate"`
	//创建时间
	CreateTime time.Time `orm:"auto_now_add;type(datetime);column(create_time)"json:"createTime"form:"createTime"`
	//状态 0收益未发放 1收益已发放
	Status int `orm:"column(status)"json:"status"form:"status"`
	//
	Name string `orm:"-"json:"name"form:"name"`
}

func FinancialProfitPageList(params *FinancialProfitQueryParam) ([]*FinancialProfit, int64) {
	o := orm.NewOrm()
	query := o.QueryTable(FinancialProfitTBName())
	data := make([]*FinancialProfit, 0)
	//默认排序
	sortorder := "id"
	switch params.Sort {
	case "id":
		sortorder = "id"
	}
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}
	//姓名
	if params.Name != ""{
		query = query.Filter("User__nick_name__exact",params.Name)
	}
	//订单号
	if params.OrderId != ""{
		query = query.Filter("financial_order_id__exact",params.OrderId)
	}
	if params.Category != ""{
		query = query.Filter("category__exact",params.Category)
	}
	//时间段
	if params.StartTime > 0 {
		query = query.Filter("is_date__gte", time.Unix(params.StartTime, 0))
	}
	if params.EndTime > 0 {
		query = query.Filter("is_date__lte", time.Unix(params.EndTime, 0))
	}
	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, (params.Offset-1)*params.Limit).RelatedSel().All(&data)
	for _,obj := range data{
		if obj.User != nil{
			obj.UserId = obj.User.UserId
			obj.Name = obj.User.NickName
		}

	}
	return data, total
}
