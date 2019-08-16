package httpModel

import (
	"github.com/astaxie/beego/orm"
	"strconv"
	"tokensky_bg_admin/models"
	"tokensky_bg_admin/utils"
)

const (
	//减少数据，允许精度
	FLOAT_PRECISE_8 float64 = 0.00000001
	FLOAT_NUM_8     int     = 8
	//允许货币类型校验[暂缺]
)

var (
	//允许的货币类型
	acceptSymbol map[string]struct{}
)
func init()  {
	acceptSymbol = make(map[string]struct{})
}

//用户资产变动
type BalanceChange struct {
	Uid int `json:"uid"`
	//货币类型
	Symbol string `json:"symbol"`
	//操作  add加 sub减 mul乘法 qup除法
	MethodBalance       string  `json:"methodBalance"`
	Balance             string  `json:"balance"`
	balance             float64 `json:"-"`
	MethodFrozenBalance string  `json:"methodFrozenBalance"`
	FrozenBalance       string  `json:"frozenBalance"`
	frozenBalance       float64 `json:"-"`
	SignId              string  `json:"signId"`
}

func (this *BalanceChange) Check()(bool,string)  {
	if this.MethodBalance == "" && this.MethodFrozenBalance == "" {
		return false,"操作空"
	}
	if this.Uid <= 0 {
		return false,"uid空"
	}
	//允许的资产类型校验
	var err error
	if _,found := acceptSymbol[this.Symbol];!found{
		if models.TokenskyUserBalanceCoinIsFound(this.Symbol){
			acceptSymbol[this.Symbol] = struct{}{}
		}else {
			return false,"资产类型不存在"
		}
	}
	found := false
	if this.Balance != "" {
		this.balance, err = strconv.ParseFloat(this.Balance, 64)
		switch this.MethodBalance {
		case "add", "sub", "mul", "qup":
		default:
			return false,"操作错误"
		}
		if this.balance < 0 {
			return false,"资产为负数"
		}
		found = true
	}
	if err != nil {
		return false,err.Error()
	}
	if this.FrozenBalance != "" {
		this.frozenBalance, err = strconv.ParseFloat(this.FrozenBalance, 64)
		switch this.MethodFrozenBalance {
		case "add", "sub", "mul", "qup":
		default:
			return false,"冻结资产操作错误"
		}
		if this.frozenBalance < 0 {
			return false,"冻结资产负数"
		}
		found = true
	}
	if err != nil {
		return false,err.Error()
	}
	return found,"ok"
}

//用户资产变化
func BalanceChangeIsOne(res *RequestOne) (bool, string, *UserBalance) {
	//校验是否重复
	o := orm.NewOrm()
	hashObj := models.TokenskyUserBalanceHashOne(o,res.HashId)
	if hashObj != nil{
		if hashObj.BalanceStatus == 1{
			return false, "重复请求", nil
		}
	}else {
		hashObj = &models.TokenskyUserBalanceHash{
			BalanceStatus:1,
			Source:res.Source,
			HashId:res.HashId,
		}
	}

	if res.Change == nil {
		return false, "数据不存在", nil
	}
	obj := res.Change
	//判断用户是否存在
	if found,msg := obj.Check();!found {
		return false, "数据校验未通过:" + msg, nil
	}

	err := o.Begin()
	if err != nil {
		return false, "开启事务失败", nil
	}
	new := false
	balance := &models.TokenskyUserBalance{}
	balanceRecord := &models.TokenskyUserBalancesRecord{
		User:   &models.TokenskyUser{UserId:obj.Uid},
		Symbol:   obj.Symbol,
		Cont:     res.Cont,
		Source:   res.Source,
		Mold:     res.Mold,
		PushTime: res.PushTime,
		SignId:   obj.SignId,
		MethodBalance:obj.MethodBalance,
		MethodFrozenBalance:obj.MethodFrozenBalance,
		Balance:obj.Balance,
		FrozenBalance:obj.FrozenBalance,
		HashId:res.HashId,
	}
	query := o.QueryTable(models.TokenskyUserBalanceTBName())
	err = query.Filter("user_id__exact", obj.Uid).Filter("coin_type__exact", obj.Symbol).One(balance)
	if err != nil {
		if err.Error() == "<QuerySeter> no row found" {
			//新增
			balance.UserId = obj.Uid
			balance.CoinType = obj.Symbol
			new = true
		} else {
			o.Rollback()
			return false, "获取数据异常", nil
		}
	}
	balanceRecord.OldBalance = balance.Balance
	balanceRecord.OldFrozenBalance = balance.FrozenBalance
	//资产计算
	ok, msg := BalanceAmount(balance, obj)
	if !ok {
		o.Rollback()
		return false, msg, nil
	}
	if new {
		_, err = o.Insert(balance)
		if err != nil {
			o.Rollback()
			return false, "新增数据失败", nil
		}
	} else {
		_, err = o.Update(balance)
		if err != nil {
			o.Rollback()
			return false, "更新数据失败", nil
		}
	}

	balanceRecord.NewBalance = balance.Balance
	balanceRecord.NewFrozenBalance = balance.FrozenBalance
	_, err = o.Insert(balanceRecord)
	if err != nil {
		o.Rollback()
		return false, "新增记录失败", nil
	}
	//新增哈希记录表
	_,err = o.InsertOrUpdate(hashObj)
	if err !=nil{
		o.Rollback()
		return false, "创建哈希记录失败", nil
	}
	err = o.Commit()
	if err != nil {
		return false, "事务执行失败", nil
	}
	resp := &UserBalance{
		Uid:           balance.UserId,
		Symbol:        balance.CoinType,
		Balance:       strconv.FormatFloat(balance.Balance, 'f', FLOAT_NUM_8, 64),
		FrozenBalance: strconv.FormatFloat(balance.FrozenBalance, 'f', FLOAT_NUM_8, 64),
	}

	return true, "", resp
}

func BalanceChangeIsMulti(res *RequestMulti) (bool, string, []*UserBalance) {
	//校验是否重复
	o := orm.NewOrm()
	hashObj := models.TokenskyUserBalanceHashOne(o,res.HashId)
	if hashObj != nil{
		if hashObj.BalanceStatus == 1{
			return false, "重复请求", nil
		}
	}else {
		hashObj = &models.TokenskyUserBalanceHash{
			BalanceStatus:1,
			Source:res.Source,
			HashId:res.HashId,
		}
	}
	objs := res.Changes
	if len(objs) < 1 {
		return false, "数据不存在", nil
	}
	changes := make(map[string]map[int]*BalanceChange)
	records := make([]*models.TokenskyUserBalancesRecord, 0)
	balances := make(map[string]map[int]*models.TokenskyUserBalance)
	resp := make([]*UserBalance, 0)
	for _, obj := range objs {
		if found,msg := obj.Check();found {
			if _, found := changes[obj.Symbol]; !found {
				changes[obj.Symbol] = make(map[int]*BalanceChange)
				//records[obj.Symbol] = make(map[int]*models.TokenskyUserBalancesRecord)
				balances[obj.Symbol] = make(map[int]*models.TokenskyUserBalance)
			}
			changes[obj.Symbol][obj.Uid] = obj
		} else {
			return false, "用户:" + strconv.Itoa(obj.Uid) + ";货币:" + obj.Symbol + " 校验失败 err:" + msg, nil
		}
	}
	err := o.Begin()
	if err != nil {
		return false, "开启事务失败", nil
	}
	for symbol, mapp := range changes {
		query := o.QueryTable(models.TokenskyUserBalanceTBName())
		query = query.Filter("coin_type__exact", symbol)
		ids := make([]int, 0,len(mapp))
		for id, _ := range mapp {
			ids = append(ids, id)
		}
		data := make([]*models.TokenskyUserBalance, 0)
		query = query.Filter("user_id__in", ids)
		_, err = query.All(&data)
		if err != nil {
			o.Rollback()
			return false, "获取数据失败", nil
		}
		ids2 := make(map[int]*models.TokenskyUserBalance)
		for _, obj := range data {
			ids2[obj.UserId] = obj
		}
		data2 := make([]*models.TokenskyUserBalance, 0)
		for _, obj := range mapp {
			if balance, found := ids2[obj.Uid]; found {
				//修改
				record := &models.TokenskyUserBalancesRecord{
					User:   &models.TokenskyUser{UserId:obj.Uid},
					Source:           res.Source,
					Cont:             res.Cont,
					Symbol:           obj.Symbol,
					OldBalance:       balance.Balance,
					OldFrozenBalance: balance.FrozenBalance,
					Mold:     res.Mold,
					PushTime: res.PushTime,
					SignId:   obj.SignId,
					MethodBalance:obj.MethodBalance,
					MethodFrozenBalance:obj.MethodFrozenBalance,
					Balance:obj.Balance,
					FrozenBalance:obj.FrozenBalance,
					HashId:res.HashId,
				}
				ok, msg := BalanceAmount(balance, obj)
				if !ok {
					o.Rollback()
					return false, msg, nil
				}
				_, err = o.Update(balance)
				if err != nil {
					o.Rollback()
					return false, "更新数据失败", nil
				}
				record.NewBalance = balance.Balance
				record.NewFrozenBalance = balance.FrozenBalance
				records = append(records, record)
				resp = append(resp, &UserBalance{
					Uid:           balance.UserId,
					Symbol:        balance.CoinType,
					Balance:       strconv.FormatFloat(balance.Balance, 'f', FLOAT_NUM_8, 64),
					FrozenBalance: strconv.FormatFloat(balance.FrozenBalance, 'f', FLOAT_NUM_8, 64),
				})
			} else {
				//新增
				balance = &models.TokenskyUserBalance{UserId: obj.Uid, CoinType: symbol}
				ok, msg := BalanceAmount(balance, obj)
				if !ok {
					o.Rollback()
					return false, msg, nil
				}
				data2 = append(data2, balance)
				record := &models.TokenskyUserBalancesRecord{
					User:   &models.TokenskyUser{UserId:obj.Uid},
					Source:           res.Source,
					Cont:             res.Cont,
					Symbol:           obj.Symbol,
					OldBalance:       0,
					OldFrozenBalance: 0,
					NewBalance:       balance.Balance,
					NewFrozenBalance: balance.FrozenBalance,
					Mold:     res.Mold,
					PushTime: res.PushTime,
					SignId:   obj.SignId,
					MethodBalance:obj.MethodBalance,
					MethodFrozenBalance:obj.MethodFrozenBalance,
					Balance:obj.Balance,
					FrozenBalance:obj.FrozenBalance,
					HashId:res.HashId,
				}
				records = append(records, record)
				resp = append(resp, &UserBalance{
					Uid:           balance.UserId,
					Symbol:        balance.CoinType,
					Balance:       strconv.FormatFloat(balance.Balance, 'f', FLOAT_NUM_8, 64),
					FrozenBalance: strconv.FormatFloat(balance.FrozenBalance, 'f', FLOAT_NUM_8, 64),
				})
			}
		}
		if len(data2) > 0 {
			_, err = o.InsertMulti(len(data2), data2)
			if err != nil {
				o.Rollback()
				return false, "新增数据失败", nil
			}
		}
	}
	if len(records) > 0 {
		_, err = o.InsertMulti(len(records), records)
		if err != nil {
			o.Rollback()
			return false, "新增记录失败", nil
		}
	}
	//新增哈希记录表
	_,err = o.InsertOrUpdate(hashObj)
	if err !=nil{
		o.Rollback()
		return false, "创建哈希记录失败", nil
	}
	err = o.Commit()
	if err != nil {
		o.Rollback()
		return false, "事务更新数据失败", nil
	}

	return true, "", resp
}

func BalanceAmount(balance *models.TokenskyUserBalance, obj *BalanceChange) (bool, string) {
	switch obj.MethodBalance {
	case "add":
		//加
		balance.Balance = utils.Float64Add(balance.Balance, obj.balance)
	case "sub":
		//减
		balance.Balance = utils.Float64Sub(balance.Balance, obj.balance)
		if balance.Balance < FLOAT_PRECISE_8 && balance.Balance > -FLOAT_PRECISE_8 {
			balance.Balance = 0
		}
		if balance.Balance < 0 {
			return false, "资产为负数"
		}
	case "mul":
		//乘
		balance.Balance = utils.Float64Mul(balance.Balance, obj.balance)
	case "quo":
		//除法
		balance.Balance = utils.Float64Quo(balance.Balance, obj.balance)
	}
	switch obj.MethodFrozenBalance {
	case "add":
		//加
		balance.FrozenBalance = utils.Float64Add(balance.FrozenBalance, obj.frozenBalance)
	case "sub":
		//减
		balance.FrozenBalance = utils.Float64Sub(balance.FrozenBalance, obj.frozenBalance)
		if balance.FrozenBalance < FLOAT_PRECISE_8 && balance.FrozenBalance > -FLOAT_PRECISE_8 {
			balance.FrozenBalance = 0
		}
		if balance.FrozenBalance < 0 {
			return false, "冻结资产为负数"
		}
	case "mul":
		//乘
		balance.FrozenBalance = utils.Float64Mul(balance.FrozenBalance, obj.frozenBalance)
	case "quo":
		//除法
		balance.FrozenBalance = utils.Float64Quo(balance.FrozenBalance, obj.frozenBalance)
	}
	//冻结资产大于实际资产
	if balance.FrozenBalance > balance.Balance {
		return false, "冻结资产大于实际资产"
	}
	return true, ""
}

//用户资产
type UserBalance struct {
	Uid           int    `json:"uid"`
	Symbol        string `json:"symbol"`
	Balance       string `json:"balance"`
	FrozenBalance string `json:"frozenBalance"`
}

//响应
type ResponseMulti struct {
	//0正常
	Code int `json:"code"`
	//返回最新数据
	Balances []*UserBalance `json:"balances"`
	//说明
	Msg string `json:"msg"`
	//哈希
	HashId string `json:"hashId"`
}

//响应
type ResponseOne struct {
	//0正常
	Code int `json:"code"`
	//返回最新数据
	Balance *UserBalance `json:"balance"`
	//说明
	Msg string `json:"msg"`
	//哈希
	HashId string `json:"hashId"`
}

//请求
type RequestMulti struct {
	//来源 1后端 2后台 3定时任务
	Source int `json:"source"`
	//资产变动多个
	Changes []*BalanceChange `json:"changes"`
	//说明
	Cont string `json:"cont"`
	//操作模版
	Mold string `json:"mold"`
	//时间戳 单位毫秒
	PushTime int64 `json:"pushTime"`
	//唯一哈希
	HashId string `json:"hashId"`
}

type RequestOne struct {
	//来源 1后端 2后台
	Source int `json:"source"`
	//单个
	Change *BalanceChange `json:"change"`
	//说明
	Cont string `json:"cont"`
	//操作模版
	Mold string `json:"mold"`
	//时间戳 单位毫秒
	PushTime int64 `json:"pushTime"`
	//唯一哈希
	HashId string `json:"hashId"`
}
