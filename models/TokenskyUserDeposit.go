package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

//充值记录表

// 用于搜索的类
type TokenskyUserDepositQueryParam struct {
	BaseQueryParam
	Status    string `json:"status"`
	StartTime int64  `json:"startTime"` //开始时间
	EndTime   int64  `json:"endTime"`   //截止时间
	Phone     string `json:"phone"`     //手机号
	OrderId   string `json:"orderId"`   //订单号
	CoinType  string `json:"coinType"`  //货币类型
}

func (a *TokenskyUserDeposit) TableName() string {
	return TokenskyUserDepositTBName()
}

type TokenskyUserDeposit struct {
	Id int `orm:"pk;column(id)"json:"id"form:"id"`
	//订单号
	OrderId string `orm:"column(order_id)"json:"orderId"form:"orderId"`
	//订单id
	DepositId int `orm:"column(deposit_id)"json:"depositId"form:"depositId"` //钱包服务生成的充值id
	//货币类型
	CoinType string `orm:"column(coin_type)"json:"coinType"form:"coinType"` //货币类型
	//交易哈希
	Txid string `orm:"column(txid)"json:"txid"form:"txid"` //充值订单号
	//交易所在区块高度
	Height int `orm:"column(height)"json:"height"form:"height"` //充值订单所在区块高度
	//金额
	Amount float64 `orm:"column(amount)"json:"amount"form:"amount"` //充值数额
	//确认区块高度
	ChainHeight int `orm:"column(chain_height)"json:"chainHeight"form:"chainHeight"` //保留字段
	//充值地址
	ToAddress string `orm:"column(to_address)"json:"toAddress"form:"toAddress"` //充值目标地址
	//转出地址
	InAddress string `orm:"-"json:"inAddress"form:"inAddress"`
	//状态 1成功
	Status int `orm:"column(status)"json:"status"form:"status"` //状态 1充值 2提币
	//创建时间
	CreateTime time.Time `orm:"auto_now_add;type(datetime);column(create_time)"json:"createTime"form:"createTime"`
	//手续费
	ServiceCharge float64 `orm:"column(service_charge)"json:"serviceCharge"form:"serviceCharge"`
	//到账时间
	FinishTime time.Time `orm:"type(datetime);column(finish_time)"json:"finishTime"form:"finishTime"`
	//用户
	User     *TokenskyUser `orm:"rel(fk);column(user_id)"json:"-"form:"-"`
	UserId   int           `orm:"-"json:"userId"form:"userId"`
	NickName string        `orm:"-"json:"nickName"form:"nickName"`
	Phone    string        `orm:"-"json:"phone"form:"phone"`
}

// 获取分页数据
func TokenskyUserDepositPageList(params *TokenskyUserDepositQueryParam) ([]*TokenskyUserDeposit, int64) {
	o := orm.NewOrm()
	query := o.QueryTable(TokenskyUserDepositTBName())
	data := make([]*TokenskyUserDeposit, 0)
	//默认排序
	sortorder := "id"
	switch params.Sort {
	case "id":
		sortorder = "id"
	}
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}
	//状态
	if params.Status != "" && params.Status != "-1" {
		query = query.Filter("status__iexact", params.Status)
	}
	//时间段
	if params.StartTime > 0 {
		query = query.Filter("create_time__gte", time.Unix(params.StartTime, 0))
	}
	if params.EndTime > 0 {
		query = query.Filter("create_time__lte", time.Unix(params.EndTime, 0))
	}
	//手机号
	if params.Phone != "" {
		query = query.Filter("User__Phone__icontains", params.Phone)
	}
	if params.OrderId != "" {
		query = query.Filter("order_id__icontains", params.OrderId)
	}
	if params.CoinType != "" {
		query = query.Filter("coin_type__ iexact", params.CoinType)
	}
	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, (params.Offset-1)*params.Limit).RelatedSel().All(&data)

	for _, obj := range data {
		//用户信息
		if obj.User != nil {
			obj.UserId = obj.User.UserId
			obj.Phone = obj.User.Phone
			obj.NickName = obj.User.NickName
		}
	}
	return data, total
}

func TokenskyUserDepositByDid(did int) *TokenskyUserDeposit {
	obj := &TokenskyUserDeposit{}

	o := orm.NewOrm()
	if err := o.QueryTable(TokenskyUserDepositTBName()).Filter("deposit_id__exact", did).One(obj); err != nil {
		return nil
	}
	return obj
}
