package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

// 提币表

/*

withdraw:
amount 提币数量
withdraw_id 请求发起方填的提现id，不能重复
to 提现目标地址

http://47.244.207.230:8080/service/BTC/address GET
http://47.244.207.230:8080/service/BTC/withdraw POST
*/

//提币审核表
type TokenskyUserTibiQueryParam struct {
	BaseQueryParam
	StartTime  int64  `json:"startTime"`  //开始时间
	EndTime    int64  `json:"endTime"`    //截止时间
	Phone      string `json:"phone"`      //手机号
	CoinType   string `json:"coinType"`   //提币类型
	OutAddress string `json:"outAddress"` //转出地址
	InAddress  string `json:"inAddress"`  //转入地址
	Status     string `json:"status"`     //状态
}

func (a *TokenskyUserTibi) TableName() string {
	return TokenskyUserTibiTBName()
}

//提币审核表
type TokenskyUserTibi struct {
	KeyId    int    `orm:"pk;column(key_id)"json:"keyId"form:"keyId"`
	OrderId  string `orm:"column(order_id)"json:"orderId"form:"orderId"`
	CoinType string `orm:"column(coin_type)"json:"coinType"form:"coinType"` //货币类型
	//交易哈希
	Txid string `orm:"column(txid)"json:"txid"form:"txid"`
	//转出地址
	OutAddress string `orm:"column(out_address)"json:"outAddress"form:"outAddress"`
	//转入地址
	InAddress string  `orm:"column(in_address)"json:"inAddress"form:"inAddress"`
	Quantity  float64 `orm:"column(quantity)"json:"quantity"form:"quantity"`
	//完成时间
	FinishTime time.Time `orm:"type(datetime);column(finish_time)"json:"finishTime"form:"finishTime"`
	//手续费
	ServiceChargeQuantity float64   `orm:"column(service_charge_quantity)"json:"serviceChargeQuantity"form:"serviceChargeQuantity"`
	ServiceCharge         float64   `orm:"column(service_charge)"json:"serviceCharge"form:"serviceCharge"`
	SumQuantity           float64   `orm:"column(sum_quantity)"json:"sumQuantity"form:"sumQuantity"`
	PushTime              time.Time `orm:"type(datetime);column(push_time)"json:"pushTime"form:"pushTime"`
	//审核时间
	VerifyTime time.Time `orm:"type(datetime);column(verify_time)"json:"verifyTime"form:"verifyTime"`
	Status     int       `orm:"column(status)"json:"status"form:"status"` //0未审核 1审核通过 2审核未通过 3处理中 4异常状态
	CreateTime time.Time `orm:"auto_now_add;type(datetime);column(create_time)"json:"createTime"form:"createTime"`
	UpdateTime time.Time `orm:"auto_now;type(datetime);column(update_time)"json:"updateTime"form:"updateTime"`
	AdminId    int       `orm:"column(admin_id)"json:"-"form:"-"`
	//连表
	User     *TokenskyUser `orm:"rel(fk);column(user_id)"json:"-"form:"-"`
	UserId   int           `orm:"-"json:"userId"form:"userId"`
	Phone    string        `orm:"-"json:"phone"form:"phone"`
	NickName string        `orm:"-"json:"nickName"form:"-"` //昵称
}

//获取分页数据
func TokenskyUserTibiPageList(params *TokenskyUserTibiQueryParam) ([]*TokenskyUserTibi, int64) {
	o := orm.NewOrm()
	query := o.QueryTable(TokenskyUserTibiTBName())
	data := make([]*TokenskyUserTibi, 0)
	//默认排序
	sortorder := "key_id"
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
	//电话查询
	if params.Phone != "" {
		query = query.Filter("User__Phone__icontains", params.Phone)
	}
	//提币类型
	if params.CoinType != "" {
		query = query.Filter("coin_type__iexact", params.CoinType)
	}
	//状态
	if params.Status != "" && params.Status != "-1" {
		query = query.Filter("status__exact", params.Status)
	}
	//转入地址 转出地址
	if params.InAddress != "" {
		query = query.Filter("in_address__icontains", params.InAddress)
	}
	if params.OutAddress != "" {
		query = query.Filter("out_address__icontains", params.OutAddress)
	}
	//时间段
	if params.StartTime > 0 {
		query = query.Filter("create_time__gte", time.Unix(params.StartTime, 0))
	}
	if params.EndTime > 0 {
		query = query.Filter("create_time__lte", time.Unix(params.EndTime, 0))
	}
	total, _ := query.Count()
	if total > 0 {
		query.OrderBy(sortorder).Limit(params.Limit, (params.Offset-1)*params.Limit).RelatedSel().All(&data)
	}

	for _, obj := range data {
		if obj.User != nil {
			obj.UserId = obj.User.UserId
			obj.Phone = obj.User.Phone
			obj.NickName = obj.User.NickName
		}
	}
	return data, total
}

//获取指定ids资源
func TokenskyUserTibiByIds(ids []string) ([]*TokenskyUserTibi, int64) {
	o := orm.NewOrm()
	query := o.QueryTable(TokenskyUserTibiTBName())
	data := make([]*TokenskyUserTibi, 0)
	query = query.Filter("key_id__in", ids)
	total, _ := query.Count()
	if total > 0 {
		query.RelatedSel().All(&data)
	}
	return data, total
}

func TokenskyUserTibiByIds2(o orm.Ormer, ids []string) ([]*TokenskyUserTibi, int64) {
	query := o.QueryTable(TokenskyUserTibiTBName())
	data := make([]*TokenskyUserTibi, 0)
	query = query.Filter("key_id__in", ids)
	total, _ := query.Count()
	if total > 0 {
		query.RelatedSel().All(&data)
	}
	return data, total
}

func TokenskyUserTibiById(id int) *TokenskyUserTibi {
	o := orm.NewOrm()
	obj := &TokenskyUserTibi{}
	if err := o.QueryTable(TokenskyUserTibiTBName()).Filter("key_id", id).One(obj); err != nil {
		return nil
	}
	return obj
}

func TokenskyUserTibiById2(o orm.Ormer, id int) *TokenskyUserTibi {
	obj := &TokenskyUserTibi{}
	if err := o.QueryTable(TokenskyUserTibiTBName()).Filter("key_id", id).One(obj); err != nil {
		return nil
	}
	return obj
}
