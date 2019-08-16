package models

import (
	"github.com/astaxie/beego/orm"
	"strconv"
	"time"
	"tokensky_bg_admin/common"
	"tokensky_bg_admin/conf"
	"tokensky_bg_admin/utils"
)

//查询的类
type HashrateSendBalanceRecordParam struct {
	BaseQueryParam
	StartTime int64  `json:"startTime"` //开始时间
	EndTime   int64  `json:"endTime"`   //截止时间
	Treaty    string `json:"treaty"`    //订单号
	Status    string `json:"status"`    //状态
	CoinType  string `json:"coinType"`  //币种
}

func (a *HashrateSendBalanceRecord) TableName() string {
	return HashrateSendBalanceRecordTBName()
}

//算力资源资产发放记录表 防止重复发放 Unknown column 'T0.-' in 'field list'
type HashrateSendBalanceRecord struct {
	KeyId int `orm:"pk;column(id)"json:"id"form:"id"`
	//货币类型
	CoinType string `orm:"column(coin_type)"json:"coinType"form:"coinType"`
	//总数
	TotalQuantity float64 `orm:"column(total_quantity)"json:"totalQuantity"form:"totalQuantity"`
	//总算力
	TotalHashrate int64 `orm:"column(total_hashrate)"json:"totalHashrate"form:"totalHashrate"`
	//发放数量
	SendQuantity float64 `orm:"column(send_quantity)"json:"sendQuantity"form:"sendQuantity"`
	//电费
	Electric float64 `orm:"column(electric)"json:"electric"form:"electric"`
	//更新时间
	UpdateTime time.Time `orm:"auto_now;type(datetime);column(update_time)"json:"updateTime"form:"updateTime"`
	//创建时间
	CreateTime time.Time `orm:"auto_now_add;type(datetime);column(create_time)"json:"createTime"form:"createTime"`
	//算力资源日期
	Isdate time.Time `orm:"type(date);column(isdate)"json:"isdate"form:"isdate"`
	//状态 是否完成 0 未完成，1完成
	Status int
	//每T收益
	UnitOutput float64 `orm:"column(unit_output)"json:"unitOutput"form:"unitOutput"`

	/*其它字段*/

	//单分收益
	Profit float64 `orm:"-"json:"-"form:"-"`
	CTName string  `orm:"-"json:"-"form:"-"`

	//每T收益
	ProfitTOne     float64 `orm:"-"json:"profitTOne"form:"-"`
	TotalHashrateP float64 `orm:"-"json:"totalHashrateP"form:"-"`
}

//分页数据
func HashrateSendBalanceRecordPageList(params *HashrateSendBalanceRecordParam) ([]*HashrateSendBalanceRecord, int64) {
	o := orm.NewOrm()
	query := o.QueryTable(HashrateSendBalanceRecordTBName())
	data := make([]*HashrateSendBalanceRecord, 0)
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
		query = query.Filter("isdate__gte", time.Unix(params.StartTime, 0))
	}
	if params.EndTime > 0 {
		query = query.Filter("isdate__lte", time.Unix(params.EndTime, 0))
	}
	//订单号
	if params.Treaty != "" {
		query = query.Filter("treaty__iexact", params.Treaty)
	}
	if params.CoinType != "" {
		query = query.Filter("coin_type__iexact", params.CoinType)
	}

	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, (params.Offset-1)*params.Limit).RelatedSel().All(&data)
	for _, obj := range data {
		//每T收益
		obj.ProfitTOne = utils.Float64Quo(obj.TotalQuantity, float64(obj.TotalHashrate/conf.HASHRATE_UNIT_T))
		//总算力
		obj.TotalHashrateP = utils.Float64Quo(float64(obj.TotalHashrate), float64(conf.HASHRATE_UNIT_T))
	}

	return data, total
}

//获取资产记录表，根据订单号和日期
func HashrateSendBalanceRecordOneByTreatyAndDay(treaty int, now int64) *HashrateSendBalanceRecord {
	var t1, t2 time.Time
	if now <= 0 {
		//取今天
		t1 = time.Now()
	} else {
		t1 = time.Unix(now, 0)
	}
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, t1.Location())
	t2 = t1.AddDate(0, 0, 1)
	obj := &HashrateSendBalanceRecord{}
	o := orm.NewOrm()
	query := o.QueryTable(HashrateSendBalanceRecordTBName())
	err := query.Filter("treaty__exact", treaty).Filter("update_time__gte", t1).Filter("update_time__lt", t2).One(obj)
	if err != nil {
		return nil
	}
	return obj
}

func HashrateSendBalanceRecordOneByCoinTypeAndDay(coinType string, now int64) *HashrateSendBalanceRecord {
	var t1, t2 time.Time
	if now <= 0 {
		//取今天
		t1 = time.Now()
	} else {
		t1 = time.Unix(now, 0)
	}
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, t1.Location())
	t2 = t1.AddDate(0, 0, 1)
	obj := &HashrateSendBalanceRecord{}
	o := orm.NewOrm()
	query := o.QueryTable(HashrateSendBalanceRecordTBName())
	err := query.Filter("coin_type__exact", coinType).Filter("update_time__gte", t1).Filter("update_time__lt", t2).One(obj)
	if err != nil {
		return nil
	}
	return obj
}

//获取资产记录表
func HashrateSendBalanceRecordByDateAndCions(tm time.Time, coin []string) ([]*HashrateSendBalanceRecord, error) {
	o := orm.NewOrm()
	query := o.QueryTable(HashrateSendBalanceRecordTBName())
	//时间转凌晨
	now := time.Date(tm.Year(), tm.Month(), tm.Day(), 0, 0, 0, 0, tm.Location())
	query = query.Filter("coin_type__in", coin)
	query = query.Filter("isdate__exact", now)
	data := make([]*HashrateSendBalanceRecord, 0)
	if _, err := query.All(&data); err != nil {
		return data, err
	}
	return data, nil
}

//拉去记录数据

/*资产发放*/

//拉取资产收益
func HashrateOrderSendBalanceGetProfitRecord(tm int64, cions []string) (map[string]*HashrateSendBalanceRecord,error) {
	if len(cions) == 0 {
		return make(map[string]*HashrateSendBalanceRecord),nil
	}
	now := time.Unix(tm, 0)

	data, err := HashrateSendBalanceRecordByDateAndCions(now, cions)
	mapp := make(map[string]*HashrateSendBalanceRecord)
	if err != nil {
		return mapp,err
	}
	for _, obj := range data {
		mapp[obj.CoinType] = obj
	}
	//新增表
	newList := make([]*HashrateSendBalanceRecord, 0)
	for _, coin := range cions {
		if _, found := mapp[coin]; !found {
			//获取历史收益
			profit, found := utils.GetViabtcProfitHistoryData(tm, coin)
			if !found {
				continue
			}
			//获取历史算力
			hashrate, found := utils.GetViabtcHashrateHistoryData(tm, coin)
			if !found {
				continue
			}
			obj := &HashrateSendBalanceRecord{
				CoinType:      coin,
				TotalQuantity: profit,
				TotalHashrate: hashrate,
				Isdate:        now,
			}
			market := SpiderCoinMarketOne(coin)
			//每T收益
			num, err := strconv.ParseFloat(market.UnitOutput, 64)
			if err != nil{
				return mapp,err
			}
			obj.UnitOutput = num
			newList = append(newList, obj)
		}
	}
	if len(newList) == 0 {
		mapp2 := make(map[string]*HashrateSendBalanceRecord)
		for _, obj := range mapp {
			if obj.Status == 0 {
				mapp2[obj.CoinType] = obj
			}
		}
		return mapp2,nil
	}
	//新增
	o := orm.NewOrm()
	if _, err := o.InsertMulti(len(newList), newList); err != nil {
		//保存错误
		return make(map[string]*HashrateSendBalanceRecord),err
	}

	data, err = HashrateSendBalanceRecordByDateAndCions(now, cions)
	mapp = make(map[string]*HashrateSendBalanceRecord)
	if err != nil {
		return mapp,err
	}
	for _, obj := range data {
		if obj.Status == 0 {
			mapp[obj.CoinType] = obj
		}
	}
	return mapp,nil
}

//创建收益表
func HashrateOrderSendBalanceCreateProfitTb(t1 int64, sendRecordMap map[string]*HashrateSendBalanceRecord) {
	// t1 发放收益的的时间

	//拉取汇率 USD汇率
	exchangeRateUsdt,found := common.UpdateUsdtExchangeRate()
	if !found{
		utils.EmailNotify("criticalToAddress","算力收益创建","获取汇率失败","")
		return
	}

	//凌晨时间戳
	now := time.Now()
	smallHours := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).Unix()
	//时间处理
	var opTime time.Time
	//指定时间
	if t1 > 0 {
		opTime = time.Unix(t1, 0)
	} else {
		return
	}
	//遍历所有的订单
	OrderIteration := HashrateOrderIteration(20, opTime)
	//算力收益
	hsByTid := HashrateOrderSendBalanceByTreaty(sendRecordMap)
	//异常的资产发放表
	errHs := make(map[int]float64)
	isIds := make(map[string]struct{})
	for {
		orders := OrderIteration()
		if len(orders) <= 0 {
			break
		}
		//算力订单收益
		hashrateOrderProfits := make([]*HashrateOrderProfit, 0)
		//防止重复收益
		ids := make([]string, 0, len(orders))
		for _, order := range orders {
			ids = append(ids, order.OrderId)
		}
		//防止重复发放奖励
		mapp := HashrateOrderProfitMapByIds(ids, opTime)
		o := orm.NewOrm()
		//开启事务
		err := o.Begin()
		if err != nil {
			for _, v := range sendRecordMap {
				errHs[v.KeyId] += 0
			}
			break
		}
		for _, order := range orders {
			if _,found := isIds[order.OrderId];found{
				continue
			}else {
				isIds[order.OrderId] = struct{}{}
			}
			//查看是否已经建立过该收益表
			if _, found := mapp[order.OrderId]; found {
				//该时间段已经有过该订单段收益
				continue
			}
			//收益
			hobj := hsByTid(order.HashrateTreaty.KeyId)
			if hobj == nil {
				continue
			}
			//if hobj.Profit <= 0 {
			//	continue
			//}
			if order.BuyQuantity < 1 {
				continue
			}
			//购买份数
			buyQuantity := float64(order.BuyQuantity)
			//实际收益
			profitNum := utils.Float64Mul(hobj.UnitOutput, buyQuantity)
			//管理费
			managementProfit := utils.Float64Mul(profitNum,order.HashrateTreaty.Management)
			//扣除管理费
			profitNum = utils.Float64Sub(profitNum,managementProfit)
			//电费
			electricity := utils.Float64Mul(buyQuantity, order.HashrateTreaty.ElectricBill/exchangeRateUsdt)
			//记录发放收益
			hobj.SendQuantity = utils.Float64Add(hobj.SendQuantity, profitNum)
			hobj.Electric = utils.Float64Add(hobj.Electric, electricity)
			hashrateOrderProfits = append(hashrateOrderProfits, &HashrateOrderProfit{
				Order:          order,
				OrderId:        TokenskyOrderIdsInsertOne(conf.ORDER_BUSINESS_HASHRATE_SEND_PRICEP_CODE),
				CategoryName:   hobj.CTName,
				User:           order.User,
				Profit:         profitNum,
				CoinType:       hobj.CoinType,
				Electricity:    electricity,
				Status:         0,
				Isdate:         opTime,
				RecordId:       hobj.KeyId,
			})
		}
		size := len(hashrateOrderProfits)
		if size == 0 {
			o.Rollback()
			continue
		}
		_, err = o.InsertMulti(len(hashrateOrderProfits), &hashrateOrderProfits)
		if err != nil {
			for _, v := range hashrateOrderProfits {
				errHs[v.RecordId] = utils.Float64Add(errHs[v.RecordId], v.Profit)
			}
			//回滚
			o.Rollback()
			break
		}
		//执行
		if err := o.Commit(); err != nil {
			for _, v := range hashrateOrderProfits {
				errHs[v.RecordId] = utils.Float64Add(errHs[v.RecordId], v.Profit)
			}
		}
		//是否最后一天发放收益
		for _, order := range orders {
			if order.EndTime.Unix() == smallHours {
				TokenskyJiguangRegistrationidSendByOne(order.User.UserId, "云算力已到期", order.OrderId, "云算力已到期", order.OrderId)
			}
		}
	}
	//完成 更新资产表
	o := orm.NewOrm()
	for _, obj := range sendRecordMap {
		if profit, found := errHs[obj.KeyId]; !found {
			obj.Status = 1
		} else {
			obj.Profit -= profit
		}
		o.Update(obj)
	}
}

//收益发放
func HashrateOrderSendBalanceProfit() {
	/*收益发放*/
	hashrateOrderProfitIteration := HashrateOrderProfitIteration(200)
	now := time.Now()
	//电费不足用户
	electricityIds := make([]int, 0)
	isIds := make(map[int]struct{})
	for {
		profitObjs := hashrateOrderProfitIteration()
		if len(profitObjs) <= 0 {
			break
		}
		for _, obj := range profitObjs {
			//不支持类型
			if _, found := conf.TOKENSKY_ACCEPT_BALANCE_COIN_TYPES[obj.CoinType]; !found {
				continue
			}
			if _,found := isIds[obj.Id];found{
				continue
			}else {
				isIds[obj.Id] = struct{}{}
			}
			o := orm.NewOrm()
			//开启事务
			err := o.Begin()
			if err != nil {
				break
			}
			electricityBalance := TokenskyUserElectricityBalanceByUid(o, obj.User.UserId)
			obj.Status = 0
			balanceChange := common.NewTokenskyUserBalanceChange(3,"hashrateOrderProfit","算力奖励发放")
			if electricityBalance == nil {
				//没有电力表信息
				obj.Status = 2
				//o.Rollback()
				//continuex
			} else {
				//扣除电费
				electricityBalance.Balance = utils.Float64Sub(electricityBalance.Balance, obj.Electricity)
				if electricityBalance.Balance < 0 {
					//电费不足
					obj.Status = 2
					//o.Rollback()
					//continue
				}
			}
			//非电费不足
			if obj.Status != 2 {
				balanceChange.Add(obj.User.UserId,obj.CoinType,obj.Order.OrderId,
					conf.CHANGE_ADD,obj.Profit,"",0)

				//用户资产表
				//balance := GetTokenskyUserBalanceByUidCoinType2(o, obj.User.UserId, obj.CoinType)
				//if balance != nil {
				//	//存在 更新操作
				//	balance.Balance = utils.Float64Add(balance.Balance, obj.Profit)
				//	if _, err := o.Update(balance); err != nil {
				//		o.Rollback()
				//		continue
				//	}
				//} else {
				//	//不存在 新增操作
				//	balance = &TokenskyUserBalance{
				//		UserId:   obj.User.UserId,
				//		CoinType: obj.CoinType,
				//		Balance:  obj.Profit,
				//	}
				//	if _, err := o.Insert(balance); err != nil {
				//		o.Rollback()
				//		continue
				//	}
				//}
				//保存电费
				if _, err := o.Update(electricityBalance); err != nil {
					o.Rollback()
					continue
				}
				//保存收益表
				obj.Status = 1
				//新增交易明细
				record := &TokenskyTransactionRecord{
					CoinType: obj.CoinType,
					TranType: "算力合约收益",
					PushTime: now,
					Category: 1, //收入
					Money:    obj.Profit,
					Status:   1,
					//TranNum:strconv.Itoa(obj.Id), 交易编号
					User:              &TokenskyUser{UserId: obj.User.UserId},
					RelevanceId:       obj.OrderId,
					RelevanceCategory: "hashrateOrderProfit",
				}
				if _, err := o.Insert(record); err != nil {
					o.Rollback()
					continue
				}
			} else {
				//电费不足
				electricityIds = append(electricityIds, obj.User.UserId)
			}
			if _, err := o.Update(obj); err != nil {
				o.Rollback()
				continue
			}
			//资产变动
			if balanceChange.Count() >0{
				ok,_,tx := balanceChange.Send()
				if !ok{
					o.Rollback()
					continue
				}
				ok = TokenskyUserBalanceHashSetStatus(o,tx)
				if !ok{
					o.Rollback()
					continue
				}
			}
			if err := o.Commit(); err != nil {
				o.Rollback()
				continue
			}
		}
	}
	//电费不足用户
	TokenskyJiguangRegistrationidSendByIds(electricityIds, "算力电费不足", "无法获取收益", "算力电费不足", "无法获取收益")

}

//获取收益 返回值为负数，代表收益不存在
func HashrateOrderSendBalanceGetProfitByCoinType(objs map[string]*HashrateSendBalanceRecord) func(tid int) float64 {
	mapp := make(map[string]map[string]float64)
	for _, obj := range objs {
		if obj.TotalHashrate <= 0 {
			continue
		}
		mapp[obj.CoinType] = make(map[string]float64)
		for _, st := range []string{"H", "K","M", "G", "T", "P", "E"} {
			switch st {
			case "H":
				mapp[obj.CoinType]["H"] = utils.Float64Quo(obj.TotalQuantity, float64(obj.TotalHashrate/conf.HASHRATE_UNIT_H))
			case "K":
				mapp[obj.CoinType]["K"] = utils.Float64Quo(obj.TotalQuantity, float64(obj.TotalHashrate/conf.HASHRATE_UNIT_K))
			case "M":
				mapp[obj.CoinType]["M"] = utils.Float64Quo(obj.TotalQuantity, float64(obj.TotalHashrate/conf.HASHRATE_UNIT_M))
			case "G":
				mapp[obj.CoinType]["G"] = utils.Float64Quo(obj.TotalQuantity, float64(obj.TotalHashrate/conf.HASHRATE_UNIT_G))
			case "T":
				mapp[obj.CoinType]["T"] = utils.Float64Quo(obj.TotalQuantity, float64(obj.TotalHashrate/conf.HASHRATE_UNIT_T))
			case "P":
				mapp[obj.CoinType]["P"] = utils.Float64Quo(obj.TotalQuantity, float64(obj.TotalHashrate/conf.HASHRATE_UNIT_P))
			case "E":
				mapp[obj.CoinType]["E"] = utils.Float64Quo(obj.TotalQuantity, float64(obj.TotalHashrate/conf.HASHRATE_UNIT_E))
			}
		}
	}
	//订单收益
	orders := make(map[int]float64)
	return func(tid int) float64 {
		if v, found := orders[tid]; found {
			return v
		} else {
			if obj := HashrateTreatyOneById(tid); obj != nil {
				if vs, found := mapp[obj.HashrateCategoryObj.Name]; found {
					if v, found := vs[obj.HashrateCategoryObj.Unit]; found {
						orders[obj.KeyId] = v
						return v
					}
				}
			}
		}
		orders[tid] = -1
		return -1
	}
}

//根据收益表获取对应的资产发放表
func HashrateOrderSendBalanceByTreaty(objs map[string]*HashrateSendBalanceRecord) func(tid int) *HashrateSendBalanceRecord {
	mapp := make(map[string]map[string]float64)
	hs := make(map[int]*HashrateSendBalanceRecord)
	for _, obj := range objs {
		if obj.TotalHashrate <= 0 {
			continue
		}
		mapp[obj.CoinType] = make(map[string]float64)
		for _, st := range []string{"H", "K","M", "G", "T", "P", "E"} {
			var num float64
			switch st {
			case "H":
				num = obj.TotalQuantity * float64(conf.HASHRATE_UNIT_H) / float64(obj.TotalHashrate)
			case "K":
				num = obj.TotalQuantity * float64(conf.HASHRATE_UNIT_K) / float64(obj.TotalHashrate)
			case "M":
				num = obj.TotalQuantity * float64(conf.HASHRATE_UNIT_M) / float64(obj.TotalHashrate)
			case "G":
				num = obj.TotalQuantity * float64(conf.HASHRATE_UNIT_G) / float64(obj.TotalHashrate)
			case "T":
				num = obj.TotalQuantity * float64(conf.HASHRATE_UNIT_T) / float64(obj.TotalHashrate)
			case "P":
				num = obj.TotalQuantity * float64(conf.HASHRATE_UNIT_P) / float64(obj.TotalHashrate)
			case "E":
				num = obj.TotalQuantity * float64(conf.HASHRATE_UNIT_E) / float64(obj.TotalHashrate)
			}
			if num >= utils.FloatGetPrec(conf.FLOAT_PRECISE_NUM_8) {
				mapp[obj.CoinType][st] = num
			}
		}
	}
	//订单收益
	return func(tid int) *HashrateSendBalanceRecord {
		if h, found := hs[tid]; found {
			return h
		} else {
			if obj := HashrateTreatyOneById(tid); obj != nil {
				if h, found := objs[obj.HashrateCategoryObj.Name]; found {
					if vs, found := mapp[obj.HashrateCategoryObj.Name]; found {
						if v, found := vs[obj.HashrateCategoryObj.Unit]; found {
							h.Profit = v
							h.CTName = obj.HashrateCategoryObj.Name
							hs[tid] = h
							return h
						}
					}
				}
			}
		}
		hs[tid] = nil
		return nil
	}
}
