package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

//查询的类
type OtcOrderQueryParam struct {
	BaseQueryParam
	//发布时间
	StartTime   int64  `json:"startTime"`   //开始时间
	EndTime     int64  `json:"endTime"`     //截止时间
	VendorPhone string `json:"vendorPhone"` //买方手机号
	VendeePhone string `json:"vendeePhone"` //卖方手机号
	OrderId     string `json:"orderId"`     //订单号
	Status      int    `json:"status"`      //状态
}

func (a *OtcOrder) TableName() string {
	return OtcOrderTBName()
}

//下单订单表
type OtcOrder struct {
	OrderId string `orm:"pk;column(order_id)"json:"orderId"form:"orderId"`
	//订单类型 订单类型  1买入 2卖出
	OrderType int `orm:"column(order_type)"json:"orderType"form:"orderType"`
	//支付方式 1银行卡 2支付宝 3微信
	PayType int `orm:"column(pay_type)"json:"payType"form:"payType"`
	//委托单ID
	EntrustOrderId int `orm:"column(entrust_order_id)"json:"entrustOrderId"form:"entrustOrderId"`
	//单价
	UnitPrice float64 `orm:"column(unit_price)"json:"unitPrice"form:"unitPrice"`
	//数量
	Quantity float64 `orm:"column(quantity)"json:"quantity"form:"quantity"`
	//总额
	TotalAmount float64 `orm:"column(total_amount)"json:"totalAmount"form:"totalAmount"`
	//下单时间
	BuyOrderTime time.Time `orm:"type(datetime);column(buy_order_time)"json:"buyOrderTime"form:"buyOrderTime"`
	//付款时间
	PayOrderTime time.Time `orm:"type(datetime);column(pay_order_time)"json:"payOrderTime"form:"payOrderTime"`
	//放币时间
	SendCoinTime time.Time `orm:"type(datetime);column(send_coin_time)"json:"sendCoinTime"form:"sendCoinTime"`
	//状态 0:待支付 等待对方支付 0:待支付 等待对方支付 1:已完成 已完成 已完成 2:已支付 等待对方放币 对方已支付 3.已申诉 卖方已申诉 已申诉 4:已取消 已取消 对方已取消
	Status int `orm:"column(status)"json:"status"form:"status"`
	//创建时间
	CreateTime time.Time `orm:"type(datetime);column(create_time)"json:"createTime"form:"createTime"`
	//取消时间
	CancelOrderTime time.Time `orm:"type(datetime);column(cancel_order_time)"json:"cancelOrderTime"form:"cancelOrderTime"`

	//关联卖方用户表
	VendorUser *TokenskyUser `orm:"rel(fk);column(vendor_user_id)"json:"-"form:"-"`
	//关联买方用户表
	VendeeUser *TokenskyUser `orm:"rel(fk);column(vendee_user_id)"json:"-"form:"-"`
	//
	VendorNickName string `orm:"-"json:"vendorNickName"`
	VendeeNickName string `orm:"-"json:"vendeeNickName"`
	//卖方用户ID
	VendorUserId int `orm:"-"json:"vendorUserId"form:"vendorUserId"`
	//买方用户ID
	VendeeUserId int `orm:"-"json:"vendeeUserId"form:"vendeeUserId"`
	//卖方电话
	VendorPhone string `orm:"-"json:"vendorPhone"`
	//买方电话
	VendeePhone string `orm:"-"json:"vendeePhone"`

	//连表 用户订单表 一对多反向关系
	OtcAppeals []*OtcAppeal `orm:"reverse(many)"json:"-"form:"-"`
}

//获取分页数据
func OtcOrderPageList(params *OtcOrderQueryParam) ([]*OtcOrder, int64) {
	o := orm.NewOrm()
	query := o.QueryTable(OtcOrderTBName())
	data := make([]*OtcOrder, 0)
	//默认排序
	sortorder := "buy_order_time"
	switch params.Sort {
	case "buy_order_time":
		sortorder = "buy_order_time"
	}
	switch params.Order {
	case "":
		sortorder = "-" + sortorder
	case "desc":
		sortorder = "-" + sortorder
	default:
		sortorder = sortorder
	}
	//状态
	if params.Status > 0 {
		query = query.Filter("status__iexact", params.Status)
	}

	//时间段
	if params.StartTime > 0 {
		query = query.Filter("buy_order_time__gte", time.Unix(params.StartTime, 0))
	}
	if params.EndTime > 0 {
		query = query.Filter("buy_order_time__lte", time.Unix(params.EndTime, 0))
	}
	//订单号
	if params.OrderId != "" {
		query = query.Filter("order_id__icontains", params.OrderId)
	}
	//卖方电话查询
	if params.VendorPhone != "" {
		query = query.Filter("VendorUser__Phone__icontains", params.VendorPhone)
	}
	//买方电话查询
	if params.VendeePhone != "" {
		query = query.Filter("VendeeUser__Phone__icontains", params.VendeePhone)
	}
	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, (params.Offset-1)*params.Limit).RelatedSel().All(&data)
	//获取用户电话
	for _, obj := range data {
		if obj.VendeeUser != nil {
			obj.VendeePhone = obj.VendeeUser.Phone
			obj.VendeeUserId = obj.VendeeUser.UserId
			obj.VendeeNickName = obj.VendeeUser.NickName
		}
		if obj.VendorUser != nil {
			obj.VendorPhone = obj.VendorUser.Phone
			obj.VendorUserId = obj.VendorUser.UserId
			obj.VendorNickName = obj.VendorUser.NickName
		}
	}
	return data, total

}

//获取买方卖方uids根据订单号
func OtcOrderGetUidsByIds(ids []string) (map[string]int, map[string]int) {
	//返回值，卖方用户 买方用户
	vendorUids := make(map[string]int)
	vendeeUids := make(map[string]int)
	query := orm.NewOrm().QueryTable(OtcOrderTBName())
	if len(ids) > 0 {
		data := make([]*OtcOrder, 0)
		query.Filter("order_id__in", ids).All(&data, "order_id", "vendor_user_id", "vendee_user_id")

		for _, v := range data {
			vendeeUids[v.OrderId] = v.VendeeUser.UserId
			vendorUids[v.OrderId] = v.VendorUser.UserId
		}
	}
	return vendorUids, vendeeUids
}

//根据电话信息模糊搜索 并返回卖方 买方 电话
func OtcOrderIdsByVendorPhoneAndVendeePhone(vendorPhone, vendeePhone string) (map[string]int, map[string]int, map[int]string) {
	//返回值，卖方用户 买方用户
	vendorUids := make(map[string]int)
	vendeeUids := make(map[string]int)
	phones := make(map[int]string)
	query := orm.NewOrm().QueryTable(OtcOrderTBName())
	if vendorPhone != "" {
		mapp := TokenskyUserGetIdsByPhone(vendorPhone)
		if len(mapp) > 0 {
			list := make([]int, 0, len(mapp))
			for uid, phone := range mapp {
				list = append(list, uid)
				phones[uid] = phone
			}
			query = query.Filter("vendor_user_id__in", list)
		} else {
			return vendorUids, vendeeUids, phones
		}
	}
	if vendeePhone != "" {
		mapp := TokenskyUserGetIdsByPhone(vendeePhone)
		if len(mapp) > 0 {
			list := make([]int, 0, len(mapp))
			for uid, phone := range mapp {
				list = append(list, uid)
				phones[uid] = phone
			}
			query = query.Filter("vendee_user_id__in", list)
		} else {
			return vendorUids, vendeeUids, phones
		}
	}
	data := make([]*OtcOrder, 0)
	query.All(&data, "order_id", "vendor_user_id", "vendee_user_id")
	for _, v := range data {
		vendorUids[v.OrderId] = v.VendorUser.UserId
		vendeeUids[v.OrderId] = v.VendeeUser.UserId
	}
	return vendorUids, vendeeUids, phones
}

func OtcOrderGetOrdersByIds(ids []string) map[string]*OtcOrder {
	data := make([]*OtcOrder, 0)
	o := orm.NewOrm().QueryTable(OtcOrderTBName())
	o.Filter("order_id__in", ids).All(&data)
	mapp := make(map[string]*OtcOrder)
	for _, v := range data {
		mapp[v.OrderId] = v
	}
	return mapp
}

//获取订单
func OtcOrderGetOrderById(id string) *OtcOrder {
	obj := OtcOrder{}
	query := orm.NewOrm().QueryTable(OtcOrderTBName())
	if err := query.Filter("order_id__exact", id).One(&obj); err != nil {
		return nil
	}
	return &obj
}

//
func OtcOrderGetOrderByOidAndBidIsStatusNot1or4(o orm.Ormer, orderId string, entrusrId int) *OtcOrder {
	obj := &OtcOrder{}
	query := o.QueryTable(OtcOrderTBName()).Exclude("order_id__exact", orderId)
	query = query.Filter("entrust_order_id__exact", entrusrId)
	query = query.Exclude("status__in", []int{1, 4})
	if err := query.One(obj); err != nil {
		return nil
	}
	return obj
}
