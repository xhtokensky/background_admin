package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

//查询的类
type BorrowOrdeLogQueryParam struct {
	BaseQueryParam
	Name      string `json:"name"`
	StartTime int64  `json:"startTime"` //开始时间
	EndTime   int64  `json:"endTime"`   //截止时间
	ConfName  string `json:"confName"`
	OrderId string `json:"orderId"`
}

func (a *BorrowOrdeLog) TableName() string {
	return BorrowOrdeLogTBName()
}

type BorrowOrdeLog struct {
	Id      int           `orm:"pk;column(id)"json:"id"form:"id"`
	User    *TokenskyUser `orm:"rel(fk);column(user_id)"json:"-"form:"-"`
	OrderId string        `orm:"column(order_id)"json:"orderId"form:"-"`
	//货币质押方式 1活期 2理财
	PledgeWay int `orm:"column(pledge_way)"json:"pledgeWay"form:"pledgeWay"`
	//操作类型
	Operation string `orm:"column(operation)"json:"operation"form:"operation"`
	//参数
	Params1    string    `orm:"column(params1)"json:"params1"form:"params1"`
	Params2    string    `orm:"column(params2)"json:"params2"form:"params2"`
	Params3    string    `orm:"column(params3)"json:"params3"form:"params3"`
	Params4    string    `orm:"column(params4)"json:"params4"form:"params4"`
	Params5    string    `orm:"column(params5)"json:"params5"form:"params5"`
	Params6    string    `orm:"column(params6)"json:"params6"form:"params6"`
	CreateTime time.Time `orm:"ype(datetime);column(create_time)"json:"createTime"form:"createTime"`

	Name     string `orm:"-"json:"name"form:"-"`
	ConfName string `orm:"-"json:"confName"form:"-"`
}

//获取追加订单记录表
func BorrowOrdeLogAddPledgePageList(params *BorrowOrdeLogQueryParam) ([]*BorrowOrdeLog, int64) {
	o := orm.NewOrm()
	//默认排序
	sortorder := "id"
	switch params.Sort {
	case "id":
		sortorder = "id"
	}
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}
	query := o.QueryTable(BorrowOrdeLogTBName())
	if params.OrderId != ""{
		query = query.Filter("order_id__icontains", params.OrderId)
	}
	//if params.ConfName !=""{
	//	query = query.Filter("Order__Conf__title__icontains",params.ConfName)
	//}
	//时间段
	if params.StartTime > 0 {
		query = query.Filter("create_time__gte", time.Unix(params.StartTime, 0))
	}
	if params.EndTime > 0 {
		query = query.Filter("create_time__lte", time.Unix(params.EndTime, 0))
	}
	if params.Name != "" {
		query = query.Filter("User__nick_name__icontains", params.Name)
	}
	query.Filter("operation__exact", "addPledge")
	data := make([]*BorrowOrdeLog, 0)
	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, (params.Offset-1)*params.Limit).RelatedSel().All(&data)
	for _, obj := range data {
		if obj.User != nil {
			obj.Name = obj.User.NickName
		}
		//if obj.Order !=nil{
		//	if obj.Order.Conf !=nil{
		//		obj.ConfName = obj.Order.Conf.Title
		//	}
		//}
	}
	return data, total
}

/*
数据库的borrow_order_log表的解析数据规则：
1、
{
 pledgeWay: 0,  /// 质押的方式 1 活期  2 理财
 operation: “addPledge”; /// 操作为增加质押数量
 params1: ‘’,              /// 增押的货币类型
         params2: ‘’,        /// 增押前质押数
         params3: ‘’,                  /// 增押后质押数
         params4: ‘’,                  /// 增加的具体值
         params5: ‘’,              /// 增押后当前的质押率
         params6: ‘’,     /// pledgeWay=2时 ,为理财包的id(如: ‘2,4,6’)
}

2、
{
 pledgeWay: 0,  /// 质押的方式 1 活期  2 理财
 operation: “create”;  /// 操作为创建借贷记录
 params1: ‘’,              /// 质押货币类型
         params2: ‘’,        /// 质押货币的数量
         params3: ‘’,                  /// 借贷货币类型
         params4: ‘’,                  /// 借贷货币的数量
         params5: ‘’,              /// 借贷当前的质押率
         params6: ‘’,     /// pledgeWay=2时 ,为理财包的id(如: ‘2,4,6’)
}

3、
{
 pledgeWay: 0,  /// 质押的方式 1 活期  2 理财
 operation: “repay”;  /// 操作为还款
 params1: ‘’,              /// 扣除的货币类型
         params2: ‘’,        /// 扣除的本金金额
         params3: ‘’,                  /// 扣除的利息金额
         params4: ‘’,                  /// 当前日期与借贷日期相差的天数
         params5: ‘’,              /// 归还的质押货币数量
         params6: ‘’,     /// pledgeWay=2时 ,为理财包的id(如: ‘2,4,6’)
}
*/
