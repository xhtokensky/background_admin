package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"strconv"
	"time"
	"tokensky_bg_admin/common"
	"tokensky_bg_admin/conf"
	"tokensky_bg_admin/enums"
	"tokensky_bg_admin/models"
	"tokensky_bg_admin/utils"
)

//充值接口
type PersonalDepositController struct {
	BaseController
}

//Prepare 参考beego官方文档说明
func (c *PersonalDepositController) Prepare() {
	////先执行
	//c.BaseController.Prepare()
	////如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	//c.checkAuthor("DataGrid", "DataList", "UpdateSeq")
	////如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	////权限控制里会进行登录验证，因此这里不用再作登录验证
	//c.checkLogin()
}

/*
提现回调:
withdraw_id 请求的提现id
asset 币种名称
status状态永远为0
回调返回200视为回调成功

充值回调:
deposit_id 钱包服务生成的充值id
asset 币种名称
to_address 充值目标地址
txid 充值订单号（最多66位（一般64）字符串）
height 充值订单所在区块高度
amount 充值数额
chain_height 保留
*/

/*  mfmQKPcWUVLwZsYkmzvK3QuiNt2XVziaoT

"to_address": "n1NryMjpmnYASfd9bzqHcu2DUvYkGwnYXS",
	"height": 103,
	"amount": "1.23400000",
	"asset": "BTC",
	"chain_height": 103,
	"deposit_id": 1,
	"txid": "d5a7e23a28c0c2adae2e6f0798a62e95a868afd4611ab6ea8face8e9982f4a11"
*/

type personalDepositData struct {
	ToAddress   string `json:"to_address"` //充值目标地址
	Height      int    `json:"height"`     //充值订单所在区块高度
	Amount      string `json:"amount"`     //充值数额
	amount      float64
	Asset       string `json:"asset"`        //币种名称
	ChainHeight int    `json:"chain_height"` //保留字段
	DepositId   int    `json:"deposit_id"`   //钱包服务生成的充值id
	Txid        string `json:"txid"`         //充值哈希字段
}

/*
{"to_address": "msVM1FjdAM297EPiwKM2pLDRnH1EroPAya",
"height": 138,
"amount": "5",
"asset": "BTC",
"chain_height": 139,
"deposit_id": 1,
"txid": "b41a1770cc8e32c1b3637a06cbf2ad6aae2c72ec6b7d608d62c34fb9710b9e51"}
*/

//充值回调
func (c *PersonalDepositController) AddBalance() {
	resp := personalDepositData{}
	body := c.Ctx.Input.RequestBody //接收raw body内容
	if err := json.Unmarshal(body, &resp); err != nil {
		c.jsonResultError2(enums.JRCodeError, "数据异常")
	}
	//判断是否重续充值
	if obj := models.TokenskyUserDepositByDid(resp.DepositId); obj != nil {
		c.jsonResultError2(enums.JRCodeFailed, "depositId 已重复")
	}
	//浮点转化
	if con, err := strconv.ParseFloat(resp.Amount, 64); err == nil {
		resp.amount = con
	}
	//1获取用户地址表
	userAddress := models.TokenskyUserAddressByCoinTypeAndAddress(resp.Asset, resp.ToAddress)
	if userAddress == nil {
		c.jsonResultError2(enums.JRCodeError, "用户地址薄不存在")
	}
	//货币类型校验
	if _, found := conf.TOKENSKY_ACCEPT_BALANCE_COIN_TYPES[resp.Asset]; !found {
		c.jsonResultError2(enums.JRCodeError, "不支持该货币类型")
	}

	/* 充币配置相关 */
	//confObj := models.TokenskyChongbiConfigGetLastOne(resp.Asset)
	//var serviceCharge,baseServiceCharge float64
	//if confObj !=nil{
	//	if resp.amount < confObj.Min{
	//		//小于最低提币数量
	//		if _, found := conf.TOKENSKY_ACCEPT_BALANCE_COIN_TYPES[resp.Asset]; !found {
	//			c.jsonResultError2(enums.JRCodeError, "充币数量小于最低限制")
	//		}
	//		serviceCharge = utils.Float64Mul(resp.amount,serviceCharge,conf.FLOAT_PRECISE_NUM_8)
	//		baseServiceCharge = confObj.BaseServiceCharge
	//	}
	//}
	////总共扣除的手续费
	//sumServiceCharge := serviceCharge + baseServiceCharge

	//change := models.NewTokenskyUserBalanceChange("充币",conf.BALANCE_CHANGE_SOURCE_ADMIN)

	//开启事务
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		c.jsonResultError2(enums.JRCodeError, "开启事务失败")
	}
	//2 获取用户资产表 没有就新增
	var logBalance float64
	//实际数量
	quantity := resp.amount
	/* 充币配置相关 扣除手续费 */
	//if sumServiceCharge > conf.FLOAT_PRECISE_8{
	//	quantity = utils.Float64Sub(quantity,serviceCharge,conf.FLOAT_PRECISE_NUM_8)
	//	if quantity <=0{
	//		o.Rollback()
	//		c.jsonResultError2(enums.JRCodeError, "手续费大于充值费用")
	//	}
	//}

	oid := models.TokenskyOrderIdsInsertOne(conf.ORDER_BUSINESS_OTC_CHONEBI_CODE)

	balanceChange := common.NewTokenskyUserBalanceChange(2,"chongbi","充币")
	balanceChange.Add(
		userAddress.UserId,userAddress.CoinType,
		oid,conf.CHANGE_ADD,quantity,
		"",0)

	tokenskyUserBalance := models.GetTokenskyUserBalanceByUidCoinType2(o, userAddress.UserId, resp.Asset)
	if tokenskyUserBalance == nil {
		tokenskyUserBalance = &models.TokenskyUserBalance{
			UserId:     userAddress.UserId,
			CoinType:   resp.Asset,
			Balance:    quantity,
			CreateTime: time.Now(),
		}
		if err != nil {
			//充值失败
			o.Rollback()
			c.jsonResultError2(enums.JRCodeError, "充值失败")
		}
	} else {
		//提币前资产
		logBalance = tokenskyUserBalance.Balance
		//提币后资产
		tokenskyUserBalance.Balance = utils.Float64Add(tokenskyUserBalance.Balance, quantity)
	}
	//_, err = o.InsertOrUpdate(tokenskyUserBalance)
	//if err != nil {
	//	o.Rollback()
	//	c.jsonResultError2(enums.JRCodeError, "更新资产失败")
	//}


	//3新建充值纪录表
	tokenskyUserDeposit := &models.TokenskyUserDeposit{
		DepositId:   resp.DepositId,
		OrderId:     oid,
		CoinType:    resp.Asset,
		Txid:        resp.Txid,
		Height:      resp.Height,
		Amount:      resp.amount,
		ChainHeight: resp.ChainHeight,
		ToAddress:   resp.ToAddress,
		/* 充币配置相关 手续费*/
		//ServiceCharge:serviceCharge,
		FinishTime: time.Now(),
		Status:     1, //成功
		User:       &models.TokenskyUser{UserId: userAddress.UserId},
	}
	_, err = o.Insert(tokenskyUserDeposit)
	if err != nil {
		o.Rollback()
		c.jsonResultError2(enums.JRCodeError, "新建记录表失败")
	}
	//新增记录
	feeLog := &models.TokenskyTransactionRecord{
		CoinType:          resp.Asset, //货币类型
		TranType:          "充币",
		PushTime:          time.Now(), //时间
		Category:          1,          //1收入
		UserId:            userAddress.UserId,
		Money:             resp.amount, //充值数额
		Status:            1,
		RelevanceId:       oid, //"Deposit:${" +depositId+ "}",
		RelevanceCategory: "chongbi",
		InAddress:         resp.ToAddress, //转入地址
		User:              &models.TokenskyUser{UserId: userAddress.UserId},
	}
	if _, err := o.Insert(feeLog); err != nil {
		//新建日志记录失败
		o.Rollback()
		c.jsonResultError2(enums.JRCodeError, "新建日志记录失败")
	}

	/* 充币配置相关 */
	//if serviceCharge > conf.FLOAT_PRECISE_8{
	//	//新增手续费记录表
	//	feeLog2 := &models.TokenskyTransactionRecord{
	//		CoinType:          resp.Asset, //货币类型
	//		TranType:          "充币手续费",
	//		PushTime:          time.Now(), //时间
	//		Category:          2,          //2支出
	//		UserId:            userAddress.UserId,
	//		Money:             serviceCharge, //数额
	//		Status:            1,
	//		RelevanceId:       oid, //"Deposit:${" +depositId+ "}",
	//		RelevanceCategory: "chongbi",
	//		InAddress:         resp.ToAddress, //转入地址
	//		User:              &models.TokenskyUser{UserId: userAddress.UserId},
	//	}
	//	if _, err := o.Insert(feeLog2); err != nil {
	//		//新建日志记录失败
	//		o.Rollback()
	//		c.jsonResultError2(enums.JRCodeError, "新建日志记录失败")
	//	}
	//}

	//资金变化
	ok,msg,tx := balanceChange.Send()
	if !ok{
		o.Rollback()
		c.jsonResultError2(enums.JRCodeError, msg)
	}
	ok = models.TokenskyUserBalanceHashSetStatus(o,tx)
	if !ok{
		o.Rollback()
		c.jsonResultError2(enums.JRCodeError, "哈希表异常")
	}
	//事务
	if err = o.Commit(); err != nil {
		o.Rollback()
		c.jsonResultError2(enums.JRCodeError, "事务失败")
	}

	/* 提币详情日志 */
	logData := map[string]interface{}{
		"操作":      "充币(回调)",
		"admin":   c.curUser.Id,
		"用户":      userAddress.UserId,
		"提币单号":    oid,
		"提币类型":    resp.Asset,
		"提币前资产":   logBalance,
		"提币后资产":   tokenskyUserBalance.Balance,
		"提币后冻结资产": tokenskyUserBalance.FrozenBalance,
	}
	logStr, _ := json.Marshal(logData)
	utils.LogNotice(logStr)
	models.TokenskyJiguangRegistrationidSendByOne(userAddress.UserId,"充币已到账",oid,"充币已到账",oid)
	c.jsonResult2(enums.JRCodeSucc, "ok")
}

/*
{"status": 0, "withdraw_id": 3, "asset": "BTC"}
*/

type callbackObj struct {
	Status     int    `json:"status"`
	WithdrawId int    `json:"withdraw_id"`
	Asset      string `json:"asset"`
	Txid       string `json:"txid"`
}

//提币回调
func (c *PersonalDepositController) Callback() {
	value := callbackObj{}
	body := c.Ctx.Input.RequestBody //接收raw body内容
	if err := json.Unmarshal(body, &value); err != nil {
		//数据解析异常
		c.jsonResultError2(enums.JRCodeError, "数据解析异常")
	}
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		//开启事务失败
		c.jsonResultError2(enums.JRCodeError, "开启事务失败")
	}
	obj := models.TokenskyUserTibiById(value.WithdrawId)
	if obj == nil {
		//不存在
		o.Rollback()
		c.jsonResultError2(enums.JRCodeError, "提币信息表不存在")
	}

	switch obj.Status {
	case 0:
		//未处理
	case 1:
		//已成功处理
		o.Rollback()
		c.jsonResultError2(enums.JRCodeFailed, "重复")
	case 3:
		//处理中
	case 4:
		//异常专状态通过
	default:
		//状态异常
		o.Rollback()
		c.jsonResultError2(enums.JRCodeFailed, "状态异常")
	}

	//状态修改
	obj.Status = 1
	obj.Txid = value.Txid
	obj.FinishTime = time.Now()
	if _, err := o.Update(obj); err != nil {
		o.Rollback()
		c.jsonResultError2(enums.JRCodeError, "更新异常")
	}
	//用户资产修改
	balance := models.GetTokenskyUserBalanceByUidCoinType2(o, obj.User.UserId, obj.CoinType)
	if balance == nil {
		//用户资产不存在
		o.Rollback()
		c.jsonResultError2(enums.JRCodeError, "用户资产异常")
	}
	//提币前资产
	logBalance := balance.Balance
	logFrozenBalance := balance.FrozenBalance
	//扣除用户冻结资金
	balance.FrozenBalance = utils.Float64Sub(balance.FrozenBalance, obj.SumQuantity)
	if balance.FrozenBalance < 0 {
		//冻结资产负数
		o.Rollback()
		c.jsonResultError2(enums.JRCodeError, "扣除用户冻结资金为负数")
	}
	balance.Balance = utils.Float64Sub(balance.Balance, obj.SumQuantity)
	if balance.Balance < 0 {
		//资产负数
		o.Rollback()
		c.jsonResultError2(enums.JRCodeError, "扣除用户资产，为负数")
	}
	//if _, err := o.Update(balance); err != nil {
	//	//资产表更新失败
	//	o.Rollback()
	//	c.jsonResultError2(enums.JRCodeError, "扣除用户冻结资金为负数")
	//}

	balanceChange := common.NewTokenskyUserBalanceChange(2,"tibi","提币回调")
	balanceChange.Add(balance.UserId,balance.CoinType,obj.OrderId,
		conf.CHANGE_SUB, obj.SumQuantity,
		conf.CHANGE_SUB,obj.SumQuantity)

	//资金变化
	ok,msg,tx := balanceChange.Send()
	if !ok{
		o.Rollback()
		c.jsonResultError2(enums.JRCodeError, msg)
	}
	ok = models.TokenskyUserBalanceHashSetStatus(o,tx)
	if !ok{
		o.Rollback()
		c.jsonResultError2(enums.JRCodeError, "哈希表异常")
	}

	//修改记录状态
	if feeLog := models.TokenskyTransactionRecordOneByRelevance(o, "提币", "tibi", obj.OrderId); feeLog != nil {
		feeLog.Status = 1
		if _, err := o.Update(feeLog); err != nil {
			o.Rollback()
			c.jsonResultError2(enums.JRCodeError, "记录表保存失败")
		}
	}
	if feeLog2 := models.TokenskyTransactionRecordOneByRelevance(o, "提币手续费", "tibi", obj.OrderId); feeLog2 != nil {
		feeLog2.Status = 1
		if _, err := o.Update(feeLog2); err != nil {
			//记录表保存失败
			o.Rollback()
			c.jsonResultError2(enums.JRCodeError, "记录表保存失败")
		}
	}
	if err = o.Commit(); err != nil {
		//事务滚动失败
		o.Rollback()
		c.jsonResultError2(enums.JRCodeError, "事务滚动失败")
	}

	/* 提币详情日志 */
	logData := map[string]interface{}{
		"操作":      "提币审核(回调)",
		"审核状态":    obj.Status,
		"admin":   c.curUser.Id,
		"用户":      balance.UserId,
		"提币单号":    obj.OrderId,
		"提币类型":    obj.CoinType,
		"提币金额":    obj.Quantity,
		"提币手续费":   obj.ServiceChargeQuantity,
		"提币前资产":   logBalance,
		"提币前冻结资产": logFrozenBalance,
		"提币后资产":   balance.Balance,
		"提币后冻结资产": balance.FrozenBalance,
	}
	logStr, _ := json.Marshal(logData)
	utils.LogNotice(logStr)
	models.TokenskyJiguangRegistrationidSendByOne(obj.User.UserId,"提币已到账",obj.OrderId,"提币已到账",obj.OrderId)
	c.jsonResult2(enums.JRCodeSucc, "ok")
}
