package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"strings"
	"time"
	"tokensky_bg_admin/common"
	"tokensky_bg_admin/conf"
	"tokensky_bg_admin/utils"
)

//查询的类
type BorrowLimitingQueryParam struct {
	BaseQueryParam
	StartTime  int64  `json:"startTime"`  //开始时间
	EndTime    int64  `json:"endTime"`    //截止时间
	OrderId     string `json:"orderId"`   //订单号
	Name       string `json:"name"`
}

func (a *BorrowLimiting) TableName() string {
	return BorrowLimitingTBName()
}

//强屏表
type BorrowLimiting struct {
	Order *BorrowOrder`orm:"pk;rel(one);column(order_id)"json:"-"form:"-"`
	//用户uid
	User  *TokenskyUser `orm:"rel(fk);column(user_id)"json:"-"form:"-"`
	//货币类型
	Symbol string `orm:"column(symbol)"json:"symbol"form:"symbol"`
	//强平时质押率
	Pledge float64 `orm:"column(pledge)"json:"pledge"form:"pledge"`
	//强平时货币价格
	SymbolPrice float64  `orm:"column(symbol_price)"json:"symbolPrice"form:"symbolPrice"`
	//售卖价格
	TotalPrice float64 `orm:"column(total_price)"json:"totalPrice"form:"totalPrice"`
	//还款额度
	PayBackPrice float64 `orm:"column(pay_back_price)"json:"payBackPrice"form:"payBackPrice"`
	//实际还款额度
	PracticalPrice float64 `orm:"column(practical_price)"json:"practicalPrice"form:"practicalPrice"`
	//执行人
	AdminId int `orm:"column(admin_id)"json:"-"form:"-"`
	//逾期时间
	ExceedTime time.Time`orm:"type(datetime);column(exceed_time)"json:"exceedTime"form:"exceedTime"`
	//售卖时间
	SellTime time.Time `orm:"type(datetime);column(sell_time)"json:"sellTime"form:"sellTime"`
	//创建时间
	CreateTime time.Time `orm:"auto_now_add;type(datetime);column(create_time)"json:"createTime"form:"createTime"`

	//关联键
	Name string `orm:"-"json:"name"form:"-"`
	//质押方式 1 活期钱包 2 理财包
	PledgeWay int `orm:"-"json:"pledgeWay"form:"-"`
	//质押货币数量
	PledgeAmount float64 `orm:"-"json:"pledgeAmount"form:"-"`
	//借贷金额
	Amount float64 `orm:"-"json:"amount"form:"-"`
	//借贷的货币类型，如:USDT
	LoanSymbol string `orm:"-"json:"loanSymbol"form:"-"`
}


//获取分页数据
func BorrowLimitingPageList(params *BorrowLimitingQueryParam) ([]*BorrowLimiting, int64) {
	o := orm.NewOrm()
	query := o.QueryTable(BorrowLimitingTBName())
	data := make([]*BorrowLimiting, 0)
	//默认排序
	sortorder := "createTime"
	switch params.Sort {
	case "create_time":
		sortorder = "id"
	}
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}
	if params.Name != ""{
		query = query.Filter("User__nick_name__exact",params.Name)
	}
	if params.OrderId != ""{
		query = query.Filter("Order__order_id__exact",params.OrderId)
	}
	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, (params.Offset-1)*params.Limit).RelatedSel().All(&data)
	for _,v := range data{
		if v.User != nil{
			v.Name = v.User.NickName
		}
		if v.Order != nil{
			v.PledgeWay = v.Order.PledgeWay
			v.PledgeAmount = v.Order.PledgeAmount
			v.Amount = v.Order.Amount
			v.LoanSymbol = v.Order.LoanSymbol
		}
	}
	return data, total
}

//售卖
func BorrowLimitingsSell(orderId string,totalPrice float64,adminId int)error{
	var err error
	obj := &BorrowLimiting{}
	o := orm.NewOrm()
	query := o.QueryTable(BorrowLimitingTBName())
	err = query.Filter("order_id__exact",orderId).RelatedSel().One(obj)
	if err != nil{
		return err
	}
	if obj.Order.Status !=5 && obj.Order.Status != 6{
		return fmt.Errorf("状态异常")
	}
	//售卖价格
	err = o.Begin()
	if err != nil{
		return fmt.Errorf("开启事务失败")
	}
	//质押理财包
	if obj.Order.RelevStatus == 1{
		order := obj.Order
		order.RepayTime = time.Now()
		//order.ForcedPrice = obj.SymbolPrice
		order.SellTotalPrice = totalPrice
		_,err :=o.Update(order)
		if err != nil{
			o.Rollback()
			return fmt.Errorf("保存订单表失败")
		}
		ids := strings.Split(obj.Order.RelevanceId,",")
		if len(ids)>0{
			params := map[string]interface{}{
				"status": 3,
			}
			financialQuery := o.QueryTable(FinancialOrderTBName())
			num, err := financialQuery.Filter("id__in", ids).Update(params)
			if err !=nil{
				o.Rollback()
				return fmt.Errorf("质押物状态修改异常 err:"+err.Error())
			}
			if num != int64(len(ids)){
				o.Rollback()
				return fmt.Errorf("质押物数量不一致")
			}
		}
		//删除记录
		err = BorrowUseFinancialOrderDelete(o,obj.Order.OrderId)
		if err != nil{
			return fmt.Errorf("删除抵押记录异常 err:"+err.Error())
		}
	}

	//售卖价格
	obj.TotalPrice = totalPrice
	//
	now := time.Now()
	newRecord1 := &TokenskyTransactionRecord{
		CoinType:obj.Symbol,
		TranType:"强平卖出",
		PushTime:now,
		Category:1,
		Money:totalPrice,
		Status:1,
		RelevanceCategory:"borrowOrder",
		RelevanceId:obj.Order.OrderId,
		User:&TokenskyUser{UserId:obj.User.UserId},
	}
	_,err = o.Insert(newRecord1)
	if err != nil{
		o.Rollback()
		return fmt.Errorf("新增强平售卖记录表失败")
	}
	//如果售卖价格小于还款额度 那么扣除金额为售卖价格
	obj.PracticalPrice = obj.PayBackPrice
	if obj.PracticalPrice>totalPrice{
		obj.PracticalPrice = totalPrice
	}
	newRecord2 := &TokenskyTransactionRecord{
		CoinType:obj.Symbol,
		TranType:"强平还款",
		PushTime:now,
		Category:2,
		Money:obj.PayBackPrice,
		Status:1,
		RelevanceCategory:"borrowOrder",
		RelevanceId:obj.Order.OrderId,
		User:&TokenskyUser{UserId:obj.User.UserId},
	}
	_,err = o.Insert(newRecord2)
	if err != nil{
		o.Rollback()
		return fmt.Errorf("新增强平还款记录表失败")
	}
	//剩余资产
	num := utils.Float64Sub(totalPrice,obj.PracticalPrice)
	obj.AdminId = adminId
	_,err = o.Update(obj)
	if err != nil{
		o.Rollback()
		return fmt.Errorf("更新强平表异常")
	}
	//资产变动
	if num > conf.FLOAT_PRECISE_8{
		balanceChange := common.NewTokenskyUserBalanceChange(3,"borrowLimitingsSell","强平售卖")
		balanceChange.Add(obj.User.UserId,obj.Symbol,obj.Order.OrderId,conf.CHANGE_ADD,num,"",0)
		ok,_,tx := balanceChange.Send()
		if !ok{
			o.Rollback()
			return fmt.Errorf("用户资产更变失败")
		}
		ok = TokenskyUserBalanceHashSetStatus(o,tx)
		if !ok{
			return fmt.Errorf("设置哈希表异常")
		}
	}
	err = o.Commit()
	if err != nil{
		o.Rollback()
		return fmt.Errorf("事务执行失败")
	}
	return nil
}
