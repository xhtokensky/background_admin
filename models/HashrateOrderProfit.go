package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

//算力收益表记录表
func (a *HashrateOrderProfit) TableName() string {
	return HashrateOrderProfitTBName()
}

//查询的类
type HashrateOrderProfitQueryParam struct {
	BaseQueryParam
	Phone     string `json:"phone"`     //手机号
	StartTime int64  `json:"startTime"` //开始时间
	EndTime   int64  `json:"endTime"`   //截止时间
	Status    string `json:"status"`    //状态
	OrderId   string `json:"orderId"`   //订单号
	UserId    string `json:"userId"`
}

type HashrateOrderProfit struct {
	Id int `orm:"pk;column(id)"json:"id"form:"id"`
	//名称
	CategoryName string `orm:"column(category_name)"json:"categoryName"form:"categoryName"`
	//创建时间
	CreateTime time.Time `orm:"auto_now_add;type(datetime);column(create_time)"json:"createTime"form:"createTime"`
	//更新时间
	UpdateTime time.Time `orm:"auto_now;type(datetime);column(update_time)"json:"updateTime"form:"updateTime"`
	//收益
	Profit float64 `orm:"column(profit)"json:"profit"form:"profit"`
	//货币类型
	CoinType string `orm:"column(coin_type)"json:"coinType"form:"coinType"`
	//电费
	Electricity float64 `orm:"column(electricity)"json:"electricity"form:"electricity"`
	//资源日期
	Isdate time.Time `orm:"type(date);column(isdate)"json:"isdate"form:"isdate"`
	//资产收益表id
	RecordId int `orm:"column(record_id)"json:"recordId"form:"recordId"`
	//状态 0收益未发放 1收益已发放 2电费不足
	Status int `orm:"column(status)"json:"status"form:"status"`

	//用户
	User     *TokenskyUser `orm:"rel(fk);column(user_id)"json:"-"form:"-"`
	UserId   int           `orm:"-"json:"userId"form:"userId"`
	NickName string        `orm:"-"json:"nickName"form:"nickName"`
	Phone    string        `orm:"-"json:"phone"form:"phone"`

	//算力订单表
	Order       *HashrateOrder `orm:"rel(fk);column(order_id)"json:"-"form:"-"`
	BuyQuantity int            `orm:"-"json:"buyQuantity"form:"buyQuantity"`
	OrderId     string         `orm:"-"json:"orderId"form:"orderId"`
	Title       string         `orm:"-"json:"title"form:"title"`
}

//获取分页数据
func HashrateOrderProfitPageList(params *HashrateOrderProfitQueryParam) ([]*HashrateOrderProfit, int64) {
	o := orm.NewOrm()
	query := o.QueryTable(HashrateOrderProfitTBName())
	data := make([]*HashrateOrderProfit, 0)
	//默认排序
	sortorder := "create_time"
	switch params.Sort {
	case "createTime":
		sortorder = "create_time"
	}
	switch params.Order {
	case "":
		sortorder = "-" + sortorder
	case "desc":
		sortorder = "-" + sortorder
	default:
		sortorder = sortorder
	}

	//时间段
	if params.StartTime > 0 {
		query = query.Filter("isdate__gte", time.Unix(params.StartTime, 0))
	}
	if params.EndTime > 0 {
		query = query.Filter("isdate__lte", time.Unix(params.EndTime, 0))
	}

	//订单号查询
	if params.OrderId != "" {
		query = query.Filter("order_id__iexact", params.OrderId)
	}
	//手机查询
	if params.Phone != "" {
		query = query.Filter("User__Phone__icontains", params.Phone)
	}
	if params.UserId != "" {
		query = query.Filter("User__user_id__iiexact", params.UserId)
	}
	total, _ := query.Count()
	query = query.OrderBy(sortorder).Limit(params.Limit, (params.Offset-1)*params.Limit)
	query.RelatedSel("User", "Order", "Order__HashrateTreaty").All(&data)
	for _, obj := range data {
		if obj.User != nil {
			//用户UID
			obj.UserId = obj.User.UserId
			//电话
			obj.Phone = obj.User.Phone
			obj.NickName = obj.User.NickName
		}
		if obj.Order != nil {
			obj.OrderId = obj.Order.OrderId
			if obj.Order.HashrateTreaty != nil {
				obj.Title = obj.Order.HashrateTreaty.Title
				obj.BuyQuantity = obj.Order.BuyQuantity
			}
		}
	}
	return data, total
}

func HashrateOrderProfitIteration(num int) func() []*HashrateOrderProfit {
	query := orm.NewOrm().QueryTable(HashrateOrderProfitTBName())
	query = query.OrderBy("create_time").Filter("status__in", []int{0, 2})
	page := 1
	count,_ := query.Count()
	return func() []*HashrateOrderProfit {
		data := make([]*HashrateOrderProfit, 0)
		if count>0{
			query.Limit(num, (page-1)*num).RelatedSel().All(&data)
			count -= int64(len(data))
			page++
		}
		return data
	}
}

//获取算力收益记录表 某个时间段 map[订单表]用户id
func HashrateOrderProfitMapByIds(ids []string, tm time.Time) map[string]struct{} {
	t1 := time.Date(tm.Year(), tm.Month(), tm.Day(), 0, 0, 0, 0, tm.Location())
	query := orm.NewOrm().QueryTable(HashrateOrderProfitTBName())
	query = query.Filter("order_id__in", ids).Filter("isdate__exact", t1)
	data := make([]*HashrateOrderProfit, 0)
	query.All(&data, "order_id")
	mapp := make(map[string]struct{})
	for _, obj := range data {
		mapp[obj.Order.OrderId] = struct{}{}
	}
	return mapp
}
