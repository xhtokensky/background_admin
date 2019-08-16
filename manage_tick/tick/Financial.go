package tick

import (
	"github.com/astaxie/beego/orm"
	"time"
	"tokensky_bg_admin/common"
	"tokensky_bg_admin/conf"
	"tokensky_bg_admin/models"
	"tokensky_bg_admin/utils"
)

//定时器相关
const (
	year_day float64 = 365 //一年天数
)

func FinancialTick()error{
	tm := time.Now()
	FinancialFixedDepositRate(tm)
	FinancialDemandDepositInterestRate(tm)
	return nil
}

//活期收益
var financialDemandDepositInterestRateSign bool = true
func FinancialDemandDepositInterestRate(tm time.Time) {
	if financialDemandDepositInterestRateSign {
		financialDemandDepositInterestRateSign = false
		defer func() { financialDemandDepositInterestRateSign = true}()
		//昨日
		yesterday := time.Date(tm.Year(), tm.Month(), tm.Day()-1, 0, 0, 0, 0, time.Local)
		//今日
		startTime := time.Date(tm.Year(), tm.Month(), tm.Day(), 0, 0, 0, 0, time.Local)
		//明日
		//endTime := time.Date(tm.Year(),tm.Month(),tm.Day()+1,0,0,0,0,time.Local)
		pTm := tm.Unix()
		rates := models.FinancialProductDemandDepositInterestRates()
		balanceIteration := models.TokenskyUserBalancesIteration(200)
		o := orm.NewOrm()
		logData := make(map[string]map[int]float64)
		logErrData := make(map[string][]int)
		isErr := false
		for str, _ := range rates {
			logData[str] = map[int]float64{}
			logErrData[str] = make([]int, 0)
		}
		var err error
		for {
			//用有资产的用户
			balances := balanceIteration()
			if len(balances) <= 0 {
				break
			}
			for _, tempBalance := range balances {
				//利率
				rate := rates[tempBalance.CoinType]
				//事务
				err = o.Begin()
				if err != nil {
					o.Rollback()
					logErrData[tempBalance.CoinType] = append(logErrData[tempBalance.CoinType], tempBalance.UserId)
					isErr = true
					continue
				}
				//1判断用户是否拥有昨天记录
				financialBalance := models.FinancialLiveUserBalanceOne(o, tempBalance.UserId, tempBalance.CoinType)
				balanceChange := common.NewTokenskyUserBalanceChange(3,"tick-financialLiveProfit","活期")
				balance := models.GetTokenskyUserBalanceByUidCoinType2(o, tempBalance.UserId, tempBalance.CoinType)
				if financialBalance != nil {
					/*利率发放*/
					expend := models.TokenskyTransactionRecordExpendByUidAndTm2(o, tempBalance.UserId, &yesterday, &startTime)
					num := utils.Float64Sub(financialBalance.Balance, expend)
					//存在利率收益
					if num > 0 {
						profit := utils.Float64Mul(rate, num)
						profit = utils.Float64Quo(profit, year_day)
						if profit > conf.FLOAT_PRECISE_8 {
							//存在活期收益
							balance.Balance = utils.Float64Add(balance.Balance, profit)
							//_, err = o.Update(balance)
							//if err != nil {
							//	o.Rollback()
							//	logErrData[tempBalance.CoinType] = append(logErrData[tempBalance.CoinType], tempBalance.UserId)
							//	isErr = true
							//	continue
							//}
							balanceChange.Add(tempBalance.UserId,tempBalance.CoinType,"",
								conf.CHANGE_ADD,profit,
								"",0,
							)
							//交易明细
							record := &models.TokenskyTransactionRecord{
								CoinType:          tempBalance.CoinType, //货币类型
								TranType:          "活期利息",
								PushTime:          tm, //时间
								Category:          1,  //1收入 2支出
								User:              &models.TokenskyUser{UserId: tempBalance.UserId},
								Money:             profit, //买方新增
								Status:            1,      //0确认中 1已完成
								RelevanceId:       "",
								RelevanceCategory: "financialLiveProfit",
							}
							_, err = o.Insert(record)
							if err != nil {
								o.Rollback()
								logErrData[tempBalance.CoinType] = append(logErrData[tempBalance.CoinType], tempBalance.UserId)
								isErr = true
								continue
							}
							//新增活期记录
							profitTb := &models.FinancialProfit{
								User:        &models.TokenskyUser{UserId: tempBalance.UserId},
								ProductId:   0,
								RelevanceId: 0,
								Product:     "live",
								Symbol:      tempBalance.CoinType,
								Balance:     balance.Balance,
								PayBalance:  expend,
								YearProfit:  rate,
								Profit:      num,
								IsDate:      tm,
								Status:      1,
								Category:1,
							}
							_, err = o.Insert(profitTb)
							if err != nil {
								o.Rollback()
								logErrData[tempBalance.CoinType] = append(logErrData[tempBalance.CoinType], tempBalance.UserId)
								isErr = true
								continue
							}
							//同步记录
							financialBalance.Balance = balance.Balance
							financialBalance.PushTime = pTm
							_, err = o.Update(financialBalance)
							if err != nil {
								o.Rollback()
								logErrData[tempBalance.CoinType] = append(logErrData[tempBalance.CoinType], tempBalance.UserId)
								isErr = true
								continue
							}
						}
					}
				} else {
					//新增同步记录
					financialBalance = &models.FinancialLiveUserBalance{
						User:     &models.TokenskyUser{UserId: tempBalance.UserId},
						Symbol:   tempBalance.CoinType,
						Balance:  balance.Balance,
						PushTime: pTm,
					}
					_, err = o.Insert(financialBalance)
					if err != nil {
						o.Rollback()
						logErrData[tempBalance.CoinType] = append(logErrData[tempBalance.CoinType], tempBalance.UserId)
						isErr = true
						continue
					}
				}

				if balanceChange.Count() >0{
					ok,_,tx := balanceChange.Send()
					if !ok{
						o.Rollback()
						isErr = true
						continue
					}
					ok = models.TokenskyUserBalanceHashSetStatus(o,tx)
					if !ok{
						o.Rollback()
						isErr = true
						continue
					}
				}
				//完成
				err = o.Commit()
				if err != nil {
					o.Rollback()
					logErrData[tempBalance.CoinType] = append(logErrData[tempBalance.CoinType], tempBalance.UserId)
					isErr = true
					continue
				}
			}
		}

		if isErr {
			data := make(map[string]interface{})
			data["title"] = "活期利息发放Err"
			data["发放成功数据"] = logData
			data["发放失败数据"] = logErrData
			utils.LogCritical(data)
		} else {
			data := make(map[string]interface{})
			data["title"] = "活期利息发放OK"
			data["httpModel"] = logData
			utils.LogNotice(data)
		}
	}
}

//定期收益
var financialFixedDepositRateSign bool = true
func FinancialFixedDepositRate(tm time.Time) {
	if financialFixedDepositRateSign {
		financialFixedDepositRateSign = false
		defer func() { financialFixedDepositRateSign = true}()
		startTime := time.Date(tm.Year(), tm.Month(), tm.Day(), 0, 0, 0, 0, time.Local)
		//当前时间需要发放的定期收益
		ids := models.FinancialOrderMaturityIdsByTm(startTime)
		o := orm.NewOrm()
		var err error
		for _, id := range ids {
			//事务
			err = o.Commit()
			if err != nil {
				o.Rollback()
				continue
			}
			obj := models.FinancialOrderOne(o, id)
			if obj == nil {
				o.Rollback()
				continue
			}
			ok,err := models.BorrowUseFinancialOrderIsLock(o,obj.Id)
			if err !=nil{
				//异常 需要提醒工作人员
				o.Rollback()
				continue
			}
			if ok{
				//订单已被锁定
				o.Rollback()
				continue
			}
			balanceChange := common.NewTokenskyUserBalanceChange(3,"tick-financialDeadProfit","定期")
			//利息发放
			if obj.QuantityLeft > 0 {
				yearProfit := utils.Float64Mul(obj.QuantityLeft, obj.YearProfit)
				num := utils.Float64Quo(float64(obj.Cycle), year_day)
				profit := utils.Float64Mul(yearProfit, num)
				//资产
				balance := models.GetTokenskyUserBalanceByUidCoinType2(o, obj.User.UserId, obj.Symbol)
				if balance == nil {
					//买定期肯定需要用户原先就有相应的货币，资产必然存在
					o.Rollback()
					continue
				}
				balance.Balance = utils.Float64Add(profit, balance.Balance)
				balance.Balance = utils.Float64Add(balance.Balance,obj.QuantityLeft)
				//_, err = o.Update(balance)
				//if err != nil {
				//	o.Rollback()
				//	continue
				//}
				balanceChange.Add(balance.UserId,balance.CoinType,obj.OrderId,
					conf.CHANGE_ADD, profit+obj.QuantityLeft,
					"",0)
				//明细记录
				record1 := &models.TokenskyTransactionRecord{
					CoinType:          balance.CoinType, //货币类型
					TranType:          "定期利息",
					PushTime:          tm, //时间
					Category:          1,  //1收入 2支出
					User:              &models.TokenskyUser{UserId: balance.UserId},
					Money:             profit, //买方新增
					Status:            1,      //0确认中 1已完成
					RelevanceId:       obj.OrderId,
					RelevanceCategory: "financialDeadProfit",
				}
				record2 := &models.TokenskyTransactionRecord{
					CoinType:          balance.CoinType, //货币类型
					TranType:          "定期成本",
					PushTime:          tm, //时间
					Category:          1,  //1收入 2支出
					User:              &models.TokenskyUser{UserId: balance.UserId},
					Money:             obj.QuantityLeft, //买方新增
					Status:            1,      //0确认中 1已完成
					RelevanceId:       obj.OrderId,
					RelevanceCategory: "financialDeadProfit",
				}
				_, err = o.InsertMulti(2,[]*models.TokenskyTransactionRecord{record1,record2})
				if err != nil {
					o.Rollback()
					continue
				}
				//收益记录
				financialProfit := &models.FinancialProfit{
					User:        &models.TokenskyUser{UserId: balance.UserId},
					ProductId:   obj.ProductId,
					RelevanceId: obj.Id,
					Category:2, //定期
					Product:     "dead",
					Symbol:      balance.CoinType,
					Balance:     obj.QuantityLeft,
					PayBalance:  0,
					YearProfit:  obj.YearProfit,
					Profit:      profit,
					IsDate:      tm,
					Status:      1,
				}
				_, err = o.Insert(financialProfit)
				if err != nil {
					o.Rollback()
					continue
				}
			}
			obj.Status = 2
			_, err = o.Update(obj)
			if err != nil {
				o.Rollback()
				continue
			}
			if balanceChange.Count()>0{
				ok,_,tx := balanceChange.Send()
				if !ok{
					o.Rollback()
					continue
				}
				ok = models.TokenskyUserBalanceHashSetStatus(o,tx)
				if !ok{
					o.Rollback()
					continue
				}
			}
			err = o.Commit()
			if err != nil {
				o.Rollback()
				continue
			}
		}
	}
}
