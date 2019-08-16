package tick

import (
	"fmt"
	"time"
	"tokensky_bg_admin/conf"
	"tokensky_bg_admin/models"
	"tokensky_bg_admin/utils"
)

//算力资产发放 定时器 每天12点30分
func TickHashrateOrderSendBalance() error {
	conf.HASHRATE_SEND_SIGN = true

	defer func() {
		conf.HASHRATE_SEND_SIGN = false
	}()

	//获取上一天的时间
	tm := time.Now().AddDate(0, 0, -1).Unix()
	//拉取收益
	HashrateSendBalanceRecords,err := models.HashrateOrderSendBalanceGetProfitRecord(tm, conf.Hashrate_Send_Balance_Allow_Coin_Type)
	if err !=nil{
		utils.EmailNotify("criticalToAddress","定时器:算力资产发放","创建收益表失败",err.Error())
		return nil
	}
	if len(HashrateSendBalanceRecords) > 0 {
		//创建收益
		models.HashrateOrderSendBalanceCreateProfitTb(tm, HashrateSendBalanceRecords)
		//奖励发放
		models.HashrateOrderSendBalanceProfit()
	} else {
		//无收益需要发放
		utils.LogCritical(fmt.Sprintf("TickHashrateOrderSendBalance 算力资产发放 无资产可发放"))
	}
	return nil
}
