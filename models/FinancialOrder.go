package models

import (
	"github.com/astaxie/beego/orm"
	"strconv"
	"time"
)

/*
CREATE TABLE `financial_order` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `order_id` varchar(100) NOT NULL COMMENT '业务号07',
  `user_id` int(11) NOT NULL,
  `symbol` varchar(50) NOT NULL,
  `financial_category_id` int(11) NOT NULL,
  `year_profit` double NOT NULL,
  `quantity` double(255,8) NOT NULL COMMENT '数量',
  `quantity_left` double(255,8) NOT NULL COMMENT '剩余数量',
  `buy_time` int(11) NOT NULL DEFAULT '0' COMMENT '购买日期',
  `effective_time` int(11) NOT NULL DEFAULT '0' COMMENT '生效时间',
  `maturity_time` int(11) NOT NULL COMMENT '到期时间',
  `cycle` int(11) NOT NULL COMMENT '周期',
  `status` int(11) NOT NULL DEFAULT '1',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='理财订单表'
*/

//查询的类
type FinancialOrderQueryParam struct {
	BaseQueryParam
	Name      string `json:"name"`
	UserId    string `json:"userId"`
	StartTime int64  `json:"startTime"` //开始时间
	EndTime   int64  `json:"endTime"`   //截止时间
	OrderId   string `json:"orderId"`   //订单号
	ConfId    string `json:"confId"`
	Symbol    string `json:"symbol"`
}

func (a *FinancialOrder) TableName() string {
	return FinancialOrderTBName()
}

//理财订单表
type FinancialOrder struct {
	Id int `orm:"pk;column(id)"json:"id"form:"id"`
	//订单号 业务号07
	OrderId string `orm:"column(order_id)"json:"orderId"form:"orderId"`
	//用户
	User   *TokenskyUser `orm:"rel(fk);column(user_id)"json:"-"form:"-"`
	UserId int           `orm:"-"json:"userId"form:"userId"`
	//货币
	Symbol string `orm:"column(symbol)"json:"symbol"form:"symbol"`
	//配置id
	ProductId int `orm:"column(product_id)"json:"productId"form:"productId"`
	//年化收益
	YearProfit float64 `orm:"column(year_profit)"json:"yearProfit"form:"yearProfit"`
	//数量
	Quantity float64 `orm:"column(quantity)"json:"quantity"form:"quantity"`
	//剩余数量
	QuantityLeft float64 `orm:"column(quantity_left)"json:"quantityLeft"form:"quantityLeft"`
	//购买时间
	BuyTime int64 `orm:"column(buy_time)"json:"buyTime"form:"buyTime"`
	//生效时间
	EffectiveTime int64 `orm:"column(effective_time)"json:"effectiveTime"form:"effectiveTime"`
	//到期时间
	MaturityTime int64 `orm:"column(maturity_time)"json:"maturityTime"form:"maturityTime"`
	//周期
	Cycle int `orm:"column(cycle)"json:"cycle"form:"cycle"`
	//状态 1进行中 2已完成 3质押强平
	Status int `orm:"column(status)"json:"status"form:"status"`
	//创建时间
	CreateTime time.Time `orm:"auto_now_add;type(datetime);column(create_time)"json:"createTime"form:"createTime"`
	//关联字段
	Name  string `orm:"-"json:"name"form:"name"`
	Title string `orm:"-"json:"title"form:"title"`
}

func FinancialOrderPageList(params *FinancialOrderQueryParam) ([]*FinancialOrder, int64) {
	data := make([]*FinancialOrder, 0)
	o := orm.NewOrm()
	query := o.QueryTable(FinancialOrderTBName())
	//默认排序
	sortorder := "id"
	switch params.Sort {
	case "id":
		sortorder = "id"
	}
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}
	if params.UserId != "" {
		query = query.Filter("User__user_id__exact", params.UserId)
	}
	if params.ConfId != "" {
		query = query.Filter("conf_id__exact", params.ConfId)
	}
	if params.OrderId != "" {
		query = query.Filter("order_id__exact", params.OrderId)
	}
	if params.Symbol != "" {
		query = query.Filter("symbol__exact", params.Symbol)
	}
	//时间段
	if params.StartTime > 0 {
		query = query.Filter("create_time__gte", time.Unix(params.StartTime, 0))
	}
	if params.EndTime > 0 {
		query = query.Filter("create_time__lte", time.Unix(params.EndTime, 0))
	}
	count,_ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, (params.Offset-1)*params.Limit).RelatedSel().All(&data)
	for _, obj := range data {
		if obj.User != nil {
			obj.UserId = obj.User.UserId
			obj.Name = obj.User.NickName
		}
		obj.Title = strconv.Itoa(obj.Cycle) + "天定期"
	}
	return data, count
}

//获取到期的
func FinancialOrderMaturityByTm(tm time.Time)[]*FinancialOrder{
	o := orm.NewOrm()
	query := o.QueryTable(FinancialOrderTBName())
	query = query.Filter("maturity_time__exact",tm)
	query = query.Filter("status__exact",1)

	data := make([]*FinancialOrder,0)
	query.All(&data)
	return data
}
func FinancialOrderMaturityIdsByTm(tm time.Time)[]int{
	ids := make([]int,0)
	o := orm.NewOrm()
	query := o.QueryTable(FinancialOrderTBName())
	query = query.Filter("maturity_time__exact",tm.Unix()*1000)
	query = query.Filter("status__exact",1)
	data := make([]*FinancialOrder,0)
	query.All(&data,"id")
	for _,obj := range data{
		ids = append(ids, obj.Id)
	}
	return ids
}
func FinancialOrderOne(o orm.Ormer,id int)*FinancialOrder{
	obj := &FinancialOrder{Id:id}
	if err := o.Read(obj,"id");err!=nil{
		return nil
	}
	return obj
}
func FinancialOrderObjsByIds(ids []string)map[int]*FinancialOrder{
	mapp := make(map[int]*FinancialOrder)
	if len(ids)>0{
		o := orm.NewOrm()
		data := make([]*FinancialOrder,0)
		query := o.QueryTable(FinancialOrderTBName())
		query.Filter("id__in",ids).All(&data)
		for _,obj := range data{
			mapp[obj.Id] = obj
		}
	}
	return mapp
}