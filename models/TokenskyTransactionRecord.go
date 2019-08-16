package models

import (
	"github.com/astaxie/beego/orm"
	"time"
	"tokensky_bg_admin/utils"
)

func (a *TokenskyTransactionRecord) TableName() string {
	return TokenskyTransactionRecordTBName()
}

//查询类
type TokenskyTransactionRecordQueryParam struct {
	BaseQueryParam
	StartTime int64  `json:"startTime"` //开始时间
	EndTime   int64  `json:"endTime"`   //截止时间
	Phone     string `json:"phone"`     //电话
	Status    string `json:"status"`    //状态
	TranType  string `json:"tranType"`  //名称类型
	CoinType  string `json:"coinType"`  //货币类型
	TranNum   string `json:"tranNum"`   //交易单号
}

//交易明细表
type TokenskyTransactionRecord struct {
	KeyId int `orm:"pk;column(key_id)"json:"keyId"form:"keyId"`
	//币种类型
	CoinType string `orm:"column(coin_type)"json:"coinType"form:"coinType"`
	//名称类型
	TranType string `orm:"column(tran_type)"json:"tranType"form:"tranType"`
	//交易时间
	PushTime time.Time `type(datetime);orm:"column(push_time)"json:"pushTime"form:"pushTime"`
	// 1收入 2支出
	Category int     `orm:"column(category)"json:"category"form:"category"`
	Money    float64 `orm:"column(money)"json:"money"form:"money"`
	//0确认中 1已完成 2已取消
	Status int `orm:"column(status)"json:"status"form:"status"`
	//关联类型 hashrateOrder tibi chongbi otcOrder hashrateOrderProfit chongElectricityOrder financialCurrent borrowOrder
	RelevanceCategory string `orm:"column(relevance_category)"json:"relevanceCategory"form:"relevanceCategory"`
	//关联ID
	RelevanceId string `orm:"column(relevance_id)"json:"relevanceId"form:"relevanceId"`
	InAddress   string `orm:"column(in_address)"json:"inAddress"form:"inAddress"`    //转入地址
	OutAddress  string `orm:"column(out_address)"json:"outAddress"form:"outAddress"` //转出地址
	//连表
	User     *TokenskyUser `orm:"rel(fk);column(user_id)"json:"-"form:"-"`
	UserId   int           `orm:"-"json:"userId"form:"userId"`
	Phone    string        `orm:"-"json:"phone"form:"phone"`
	NickName string        `orm:"-"json:"nickName"form:"-"` //昵称
}

func TokenskyTransactionRecordPageList(params *TokenskyTransactionRecordQueryParam) ([]*TokenskyTransactionRecord, int64) {
	o := orm.NewOrm()
	query := o.QueryTable(TokenskyTransactionRecordTBName())
	data := make([]*TokenskyTransactionRecord, 0)
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
	//状态
	if params.Status != "" && params.Status != "-1" {
		query = query.Filter("status__iexact", params.Status)
	}
	//时间段
	if params.StartTime > 0 {
		query = query.Filter("push_time__gte", time.Unix(params.StartTime, 0))
	}
	if params.EndTime > 0 {
		query = query.Filter("push_time__lte", time.Unix(params.EndTime, 0))
	}
	//手机号
	if params.Phone != "" {
		query = query.Filter("User__Phone__icontains", params.Phone)
	}
	//货币类型
	if params.CoinType != "" {
		query = query.Filter("coin_type__iexact", params.CoinType)
	}
	//名称类型
	if params.TranType != "" {
		query = query.Filter("tran_type__iexact", params.TranType)
	}
	//单号
	if params.TranNum != "" {
		query = query.Filter("relevance_id__icontains", params.TranNum)
	}
	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, (params.Offset-1)*params.Limit).RelatedSel().All(&data)
	for _, obj := range data {
		if obj.User != nil {
			//用户UID
			obj.UserId = obj.User.UserId
			//电话
			obj.Phone = obj.User.Phone
			obj.NickName = obj.User.NickName
		}
	}
	return data, total
}

//获取交易明细表
func TokenskyTransactionRecordOneByRelevance(o orm.Ormer, tranType, relevanceCategory, relevanceId string) *TokenskyTransactionRecord {
	if o == nil {
		o = orm.NewOrm()
	}
	obj := &TokenskyTransactionRecord{}
	query := o.QueryTable(TokenskyTransactionRecordTBName())
	query = query.Filter("relevance_category__exact", relevanceCategory)
	query = query.Filter("tran_type__iexact", tranType)
	query = query.Filter("relevance_id__exact", relevanceId)
	if err := query.One(obj); err != nil {
		return nil
	}
	return obj
}

//获取某段时间用户的收入
func TokenskyTransactionRecordIncomeByUidAndTm(uid int,StartTime,endTime *time.Time)float64{
	o := orm.NewOrm()
	query := o.QueryTable(TokenskyTransactionRecordTBName())
	query = query.Filter("User__user_id__exact",uid)
	query = query.Filter("category__exact",1)
	if StartTime != nil{
		query = query.Filter("push_time__gte",StartTime)
	}
	if endTime != nil{
		query = query.Filter("push_time__lt",endTime)
	}
	data := make([]*TokenskyTransactionRecord,0)
	query.All(&data)
	var num float64
	for _,obj := range data{
		num = utils.Float64Add(num,obj.Money)
	}
	return num
}

//获取某段时间用户的支出
func TokenskyTransactionRecordExpendByUidAndTm(uid int,StartTime,endTime *time.Time)float64{
	o := orm.NewOrm()
	query := o.QueryTable(TokenskyTransactionRecordTBName())
	query = query.Filter("User__user_id__exact",uid)
	query = query.Filter("category__exact",2)
	if StartTime != nil{
		query = query.Filter("push_time__gte",StartTime)
	}
	if endTime != nil{
		query = query.Filter("push_time__lt",endTime)
	}
	data := make([]*TokenskyTransactionRecord,0)
	query.All(&data)
	var num float64
	for _,obj := range data{
		num = utils.Float64Add(num,obj.Money)
	}
	return num
}
func TokenskyTransactionRecordExpendByUidAndTm2(o orm.Ormer,uid int,StartTime,endTime *time.Time)float64{
	query := o.QueryTable(TokenskyTransactionRecordTBName())
	query = query.Filter("User__user_id__exact",uid)
	query = query.Filter("category__exact",2)
	if StartTime != nil{
		query = query.Filter("push_time__gte",StartTime)
	}
	if endTime != nil{
		query = query.Filter("push_time__lt",endTime)
	}
	data := make([]*TokenskyTransactionRecord,0)
	query.All(&data)
	var num float64
	for _,obj := range data{
		num = utils.Float64Add(num,obj.Money)
	}
	return num
}