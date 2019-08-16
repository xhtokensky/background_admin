package models

import (
	"github.com/astaxie/beego/orm"
	"strings"
	"time"
	"tokensky_bg_admin/utils"
)

//查询的类
type OtcAppealQueryParam struct {
	BaseQueryParam
	StartTime   int64  `json:"startTime"`   //开始时间
	EndTime     int64  `json:"endTime"`     //截止时间
	VendorPhone string `json:"vendorPhone"` //买方手机号
	VendeePhone string `json:"vendeePhone"` //卖方手机号
	OrderId     string `json:"orderId"`     //订单号
	Status      string `json:"status"`      //状态
}

func (a *OtcAppeal) TableName() string {
	return OtcAppealTBName()
}

//申诉列表
type OtcAppeal struct {
	KeyId int `orm:"pk;column(key_id)"json:"keyId"form:"keyId"`
	//申诉时间
	AppealTime time.Time `orm:"type(datetime);column(appeal_time)"json:"appealTime"form:"appealTime"`
	//上传凭证时间
	UpVoucherTime time.Time `orm:"type(datetime);column(up_voucher_time)"json:"upVoucherTime"form:"upVoucherTime"`
	//卖方申诉原因
	VendorCause string `orm:"column(vendor_cause)"json:"vendorCause"form:"vendorCause"`
	//卖方上传凭证
	VendorVoucher  string   `orm:"column(vendor_voucher)"json:"vendorVoucher"form:"vendorVoucher"`
	VendorVouchers []string `orm:"-"json:"vendorVouchers"form:"vendorVouchers"`
	//买方凭证备注
	VendeeRemark string `orm:"column(vendee_remark)"json:"vendeeRemark"form:"vendeeRemark"`
	//买方凭证
	VendeeVoucher  string   `orm:"column(vendee_voucher)"json:"vendeeVoucher"form:"vendeeVoucher"`
	VendeeVouchers []string `orm:"-"json:"vendeeVouchers"form:"vendeeVouchers"`
	//状态 0未处理 1 确认放币 2取消放币
	Status int `orm:"column(status)"json:"status"form:"status"`
	//创建时间
	CreateTime time.Time `orm:"type(datetime);column(create_time)"json:"createTime"form:"createTime"`
	//关联表
	Order *OtcOrder `orm:"rel(fk);column(order_id)"json:"-"form:"-"`

	/*连表字段*/
	//订单号
	OrderId string `orm:"-"json:"orderId"form:"orderId"`
	//卖方电话
	VendorPhone    string `orm:"-"json:"vendorPhone"`
	VendorNickName string `orm:"-"json:"vendorNickName"`
	//买方电话
	VendeePhone    string `orm:"-"json:"vendeePhone"`
	VendeeNickName string `orm:"-"json:"vendeeNickName"`

	//支付方式 1银行卡 2支付宝 3微信
	PayType int `orm:"-"json:"payType"form:"payType"`
	//单价
	UnitPrice float64 `orm:"-"json:"unitPrice"form:"unitPrice"`
	//数量
	Quantity float64 `orm:"-"json:"quantity"form:"quantity"`
	//总额
	TotalAmount float64 `orm:"-"json:"totalAmount"form:"totalAmount"`
	//下单时间
	BuyOrderTime time.Time `orm:"-"json:"buyOrderTime"form:"buyOrderTime"`
	//付款时间
	PayOrderTime time.Time `orm:"-"json:"payOrderTime"form:"payOrderTime"`
}

//获取分页数据
func OtcOrderAppealPageList(params *OtcAppealQueryParam) ([]*OtcAppeal, int64) {
	o := orm.NewOrm()
	query := o.QueryTable(OtcAppealTBName())
	data := make([]*OtcAppeal, 0)
	//默认排序
	sortorder := "keyId"
	switch params.Sort {
	case "keyId":
		sortorder = "key_id"
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
	if params.Status != "" && params.Status != "-1" {
		query = query.Filter("status__iexact", params.Status)
	}
	//时间段
	if params.StartTime > 0 {
		query = query.Filter("appeal_time__gte", time.Unix(params.StartTime, 0))
	}
	if params.EndTime > 0 {
		query = query.Filter("appeal_time__lte", time.Unix(params.EndTime, 0))
	}
	//订单号
	if params.OrderId != "" {
		query = query.Filter("order_id__icontains", params.OrderId)
	}

	//卖方电话查询
	if params.VendorPhone != "" {
		query = query.Filter("Order__VendorUser__Phone__icontains", params.VendorPhone)
	}
	//买方电话查询
	if params.VendeePhone != "" {
		query = query.Filter("Order__VendeeUser__Phone__icontains", params.VendeePhone)
	}
	//非手机号查询版
	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, (params.Offset-1)*params.Limit).RelatedSel().All(&data)
	for _, obj := range data {
		if obj.Order != nil {
			obj.OrderId = obj.Order.OrderId
			obj.PayType = obj.Order.PayType
			obj.UnitPrice = obj.Order.UnitPrice
			obj.Quantity = obj.Order.Quantity
			obj.TotalAmount = obj.Order.TotalAmount
			obj.BuyOrderTime = obj.Order.BuyOrderTime
			obj.PayOrderTime = obj.Order.PayOrderTime
			//电话
			if obj.Order.VendorUser != nil {
				obj.VendorPhone = obj.Order.VendorUser.Phone
				obj.VendorNickName = obj.Order.VendorUser.NickName
			}
			if obj.Order.VendeeUser != nil {
				obj.VendeePhone = obj.Order.VendeeUser.Phone
				obj.VendeeNickName = obj.Order.VendeeUser.NickName
			}
		}
		//申诉图片
		if list := strings.Split(obj.VendeeVoucher, ","); len(list) > 0 {
			for _, str := range list {
				obj.VendeeVouchers = append(obj.VendeeVouchers, utils.QiNiuDownload(str, 0))
			}
		}
		if list := strings.Split(obj.VendorVoucher, ","); len(list) > 0 {
			for _, str := range list {
				obj.VendorVouchers = append(obj.VendorVouchers, utils.QiNiuDownload(str, 0))
			}
		}
	}
	return data, total
}

//根据id获取单条记录
func OtcAppealOneById(KeyId int) (*OtcAppeal, error) {
	o := orm.NewOrm()
	query := o.QueryTable(OtcAppealTBName())
	m := &OtcAppeal{}
	err := query.Filter("key_id__exact",KeyId).RelatedSel().One(m)
	if err != nil {
		return nil, err
	}
	m.OrderId = m.Order.OrderId
	return m, nil
}
