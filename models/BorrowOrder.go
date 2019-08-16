package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"time"
	"tokensky_bg_admin/common"
)

//查询的类
type BorrowOrderQueryParam struct {
	BaseQueryParam
	Name       string `json:"name"`
	StartTime  int64  `json:"startTime"`  //开始时间
	EndTime    int64  `json:"endTime"`    //截止时间
	Symbol     string `json:"symbol"`     //质押货币类型
	LoanSymbol string `json:"loanSymbol"` //借贷货币类型
	Status     string `json:"status"`     //状态
}

func (a *BorrowOrder) TableName() string {
	return BorrowOrderTBName()
}

type BorrowOrder struct {
	Id int `orm:"pk;column(id)"json:"id"form:"id"`
	//借贷记录id
	OrderId string `orm:"column(order_id)"json:"orderId"form:"orderId"`
	//配置id
	Conf *BorrowConf `orm:"rel(fk);column(conf_id)"json:"-"form:"-"`
	//用户id
	User   *TokenskyUser `orm:"rel(fk);column(user_id)"json:"userId"form:"userId"`
	UserId int           `orm:"-"json:"userId"form:"-"`
	//质押货币类型
	Symbol string `orm:"column(symbol)"json:"symbol"form:"symbol"`
	//质押货币数量
	PledgeAmount float64 `orm:"column(pledge_amount)"json:"pledgeAmount"form:"pledgeAmount"`
	//质押率
	ForcedingPledgeRate float64 `orm:"column(forceding_pledge_rate)"json:"forcedingPledgeRate"form:"forcedingPledgeRate"`
	//质押的周期月数
	CycleMonth int `orm:"column(cycle_month)"json:"cycleMonth"form:"cycleMonth"`
	//周期月数的天数
	CycleMonthDay int `orm:"column(cycle_month_day)"json:"cycleMonthDay"form:"cycleMonthDay"`
	//质押的日利率
	PledgeDayRate float64 `orm:"column(pledge_day_rate)"json:"pledgeDayRate"form:"pledgeDayRate"`
	//质押方式 1 活期钱包 2 理财包
	PledgeWay int `orm:"column(pledge_way)"json:"pledgeWay"form:"pledgeWay"`
	//被质押的理财包id集合 理财订单表id
	RelevanceId string `orm:"column(relevance_id)"json:"relevanceId"form:"relevanceId"`
	//被质押的理财包的状态 1 是 2 否
	RelevStatus int `orm:"column(relev_status)"json:"relevStatus"form:"relevStatus"`
	//借贷金额
	Amount float64 `orm:"column(amount)"json:"amount"form:"amount"`
	//借贷的货币类型，如:  USDT
	LoanSymbol string `orm:"column(loan_symbol)"json:"loanSymbol"form:"loanSymbol"`
	//还款时的利息
	RepayInterest float64 `orm:"column(repay_interest)"json:"repayInterest"form:"repayInterest"`
	//借贷时间
	BorrowTime time.Time `orm:"type(datetime);column(borrow_time)"json:"borrowTime"form:"borrowTime"`
	//到期时间
	ExpireTime time.Time `orm:"type(datetime);column(expire_time)"json:"expireTime"form:"expireTime"`
	//手动还贷时间或被强平的时间
	RepayTime time.Time `orm:"type(datetime);column(repay_time)"json:"repayTime"form:"repayTime"`
	//进入强平的时间(逾期或超过最大质押率)
	ForcedingTime time.Time `orm:"type(datetime);column(forceding_time)"json:"forcedingTime"form:"forcedingTime"`
	//最近一次增加质押时间
	AddPledgeTime time.Time `orm:"type(datetime);column(add_pledge_time)"json:"addPledgeTime"form:"addPledgeTime"`
	//上次提醒时间
	WarnTime time.Time `orm:"type(datetime);column(warn_time)"json:"warnTime"form:"warnTime"`
	//关联表
	BorrowLimiting *BorrowLimiting `orm:"reverse(one)"json:"coinData"form:"-"`
	//强平时结算价格
	ForcedPrice float64 `orm:"column(forced_price)"json:"forcedPrice"form:"forcedPrice"`
	//售卖总价格
	SellTotalPrice float64 `orm:"column(sell_total_price)"json:"sellTotalPrice"form:"sellTotalPrice"`
	//状态 记录状态 1 使用中; 2 还贷日;  4 已还贷; 5 逾期被强平中; 6 最大质押率被强平中;  7 逾期已强平;  8 最大质押率已强平;
	Status int `orm:"column(status)"json:"status"form:"status"`



	//关联字段
	Name string `orm:"-"json:"name"form:"-"`

	//实时质押率
	RealTimePledge float64 `orm:"-"json:"realTimePledge"form:"-"`
	//实时利息
	RealTimeInterest float64 `orm:"-"json:"realTimeInterest"form:"-"`
	//货币价格
	SymbolPrice float64 `orm:"-"json:"-"form:"-"`
	LoanSymbolPrice float64 `orm:"-"json:"-"form:"-"`
	//最高质押数额
	MaxPledge float64  `orm:"-"json:"-"form:"-"`
}

//实时质押率
func (this *BorrowOrder) GetRealTimePledge()bool{
	rate, ok := common.GetSymbolExchangeRate2(this.Symbol, "USD")
	if !ok {
		this.RealTimePledge = 0
		return false
	}
	this.SymbolPrice = rate
	loanRate, ok := common.GetSymbolExchangeRate2(this.LoanSymbol, "USD")
	if !ok {
		this.RealTimePledge = 0
		return false
	}
	this.LoanSymbolPrice = loanRate
	this.RealTimePledge = this.Amount / (this.PledgeAmount * (rate / loanRate))
	return true
}

//实时利息
func (this *BorrowOrder) GetRealTimeInterest() {
	now := time.Now()
	tm := now.Sub(this.BorrowTime)
	//实时天数
	day := int(tm.Hours())/24 + 1
	max := this.CycleMonth * this.CycleMonthDay
	//逾期
	if day > max {
		day = max
	}
	this.RealTimeInterest = this.Amount * float64(day) * this.RepayInterest
}

//最高质押数额
func (this *BorrowOrder)GetMaxPledge()bool{
	rate, ok := common.GetSymbolExchangeRate2(this.Symbol, "USD")
	if !ok{
		return false
	}
	this.MaxPledge = rate*this.PledgeAmount
	return true
}

func BorrowOrderGetMaxPledge(id int)float64{
	o := orm.NewOrm()
	query := o.QueryTable(BorrowOrderTBName())
	obj := &BorrowOrder{}
	err:= query.Filter("id__exact",id).One(obj)
	if err != nil{
		return 0
	}
	obj.GetMaxPledge()
	return obj.MaxPledge
}

//获取分页数据
func BorrowOrderPageList(params *BorrowOrderQueryParam) ([]*BorrowOrder, int64) {
	o := orm.NewOrm()
	query := o.QueryTable(BorrowOrderTBName())
	data := make([]*BorrowOrder, 0)
	//默认排序
	sortorder := "id"
	switch params.Sort {
	case "id":
		sortorder = "id"
	}
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}
	if params.Symbol != "" {
		query = query.Filter("symbol__exact", params.Symbol)
	}
	if params.LoanSymbol != "" {
		query = query.Filter("loan_symbol__exact", params.LoanSymbol)
	}
	if params.Name != "" {
		query = query.Filter("User__nick_name__exact", params.Name)
	}
	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, (params.Offset-1)*params.Limit).RelatedSel().All(&data)
	for _, obj := range data {
		if obj.User != nil {
			obj.Name = obj.User.NickName
		}
	}

	return data, total
}

func BorrowOrderIteration()func()[]*BorrowOrder{
	o := orm.NewOrm()
	num := 1
	query := o.QueryTable(BorrowOrderTBName())
	query.Filter("status__lt",4)
	count,_ := query.Count()
	return func() []*BorrowOrder {
		data := make([]*BorrowOrder,0)
		if count>0{
			query.Limit(500,(num-1)*500).All(&data)
			count -= int64(len(data))
			num++
		}
	return data
	}
}

//逾期强平
func BorrowOrderExpireLimiting(o orm.Ormer,obj *BorrowOrder)error{
	var err error
	if !obj.GetRealTimePledge(){
		return fmt.Errorf("GetRealTimePledge 获取实时质押失败")
	}
	obj.GetRealTimeInterest()
	err = o.Begin()
	if err != nil{
		return err
	}
	obj.Status = 6
	obj.RepayInterest = obj.RealTimeInterest
	obj.ForcedingTime = time.Now()
	obj.ForcedPrice = obj.SymbolPrice
	_,err = o.Update(obj)
	if err != nil{
		return err
	}
	newObj := &BorrowLimiting{
		Order:&BorrowOrder{Id:obj.Id},
		User:&TokenskyUser{UserId:obj.User.UserId},
		Symbol:obj.Symbol,
		Pledge:obj.RealTimePledge,
		SymbolPrice:obj.SymbolPrice,
		ExceedTime:obj.ExpireTime,
		PayBackPrice:obj.RepayInterest+obj.Amount,
	}
	_,err = o.Insert(newObj)
	if err != nil{
		return err
	}
	err = o.Commit()
	if err != nil{
		return err
	}
	return nil
}

//质押强平
func BorrowOrderPledgeLimiting(o orm.Ormer,obj *BorrowOrder) error {
	var err error
	if !obj.GetRealTimePledge(){
		return fmt.Errorf("GetRealTimePledge 获取实时质押失败")
	}
	obj.GetRealTimeInterest()
	err = o.Begin()
	if err != nil{
		return err
	}
	obj.Status = 5
	obj.RepayInterest = obj.RealTimeInterest
	obj.ForcedingPledgeRate = obj.RealTimePledge
	//obj.RepayTime = time.Now()
	_,err = o.Update(obj)
	if err != nil{
		return err
	}
	newObj := &BorrowLimiting{
		Order:&BorrowOrder{Id:obj.Id},
		User:&TokenskyUser{UserId:obj.User.UserId},
		Symbol:obj.Symbol,
		Pledge:obj.RealTimePledge,
		SymbolPrice:obj.SymbolPrice,
		ExceedTime:obj.ExpireTime,
		PayBackPrice:obj.RepayInterest+obj.Amount,
	}
	_,err = o.Insert(newObj)
	if err != nil{
		return err
	}
	err = o.Commit()
	if err != nil{
		return err
	}
	return nil
}

//警告提醒
func BorrowOrderWarn(o orm.Ormer,obj *BorrowOrder)  {
	//警告代码...


	obj.ExpireTime = time.Now()
	o.Update(obj)
}