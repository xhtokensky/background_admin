package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"time"
	"tokensky_bg_admin/common"
	"tokensky_bg_admin/conf"
	"tokensky_bg_admin/enums"
	"tokensky_bg_admin/models"
	"tokensky_bg_admin/utils"
)

//OtcAppealController OTC-申诉表
type OtcAppealController struct {
	BaseController
}

//Prepare 参考beego官方文档说明
func (c *OtcAppealController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	c.checkAuthor("DataGrid", "DataList", "UpdateSeq")
}

// DataGrid
func (c *OtcAppealController) DataGrid() {
	//直接反序化获取json格式的requestbody里的值
	var params models.OtcAppealQueryParam
	switch c.Ctx.Request.Method {
	case "POST":
		json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	case "GET":
		params.Sort = c.GetString("sort")
		params.Order = c.GetString("order")
		params.Limit, _ = c.GetInt64("limit")
		params.Offset, _ = c.GetInt64("offset")
		params.VendorPhone = c.GetString("vendorPhone")
		params.VendeePhone = c.GetString("vendeePhone")
		params.StartTime, _ = c.GetInt64("startTime")
		params.EndTime, _ = c.GetInt64("endTime")
		params.OrderId = c.GetString("orderId")
		params.Status = c.GetString("status")
	}
	//获取数据列表和总数
	data, total := models.OtcOrderAppealPageList(&params)
	//定义返回的数据结构
	mapp := map[string]interface{}{
		"rows":  data,
		"total": total,
	}
	c.jsonResult(enums.JRCodeSucc, "", mapp)
}

//审核结果
func (c *OtcAppealController) Examine() {
	if c.Ctx.Request.Method == "POST" {
		var err error
		var msg,tx string
		var ok bool
		var logData map[string]interface{}
		var logStr []byte
		balanceChange := common.NewTokenskyUserBalanceChange(2,"otc","申诉审核")
		m := models.OtcAppeal{}
		body := c.Ctx.Input.RequestBody //接收raw body内容
		if err := json.Unmarshal(body, &m); err != nil {
			c.jsonResult(enums.JRCodeFailed, "数据异常:", m.KeyId)
		}
		//1校验订单是否存在
		obj, err := models.OtcAppealOneById(m.KeyId)
		if err != nil {
			c.jsonResult(enums.JRCodeFailed, "无次订单号:", m.KeyId)
		}
		//是否审核校验
		switch obj.Status {
		case conf.OTC_APPEAL_STATUS_UOD:
			//pass
		case conf.OTC_APPEAL_STATUS_VALIDATE, conf.OTC_APPEAL_STATUS_CANCEL:
			c.jsonResult(enums.JRCodeFailed, "已经审核过", m.KeyId)
		default:
			c.jsonResult(enums.JRCodeFailed, "未知审核状态:", m.KeyId)
		}
		if err != nil {
			c.jsonResult(enums.JRCodeFailed, "申诉表是否存在", m.KeyId)
		}
		//订单状态修改
		obj.Status = m.Status
		//校验订单是否存在
		order := models.OtcOrderGetOrderById(obj.OrderId)
		if order == nil {
			c.jsonResult(enums.JRCodeFailed, "订单不存在", m.KeyId)
		}

		switch obj.Status {
		case conf.OTC_APPEAL_STATUS_CANCEL:
			//取消放币
			if order.Status == conf.OTC_ORDER_STATUS_FOUND || order.Status == conf.OTC_ORDER_STATUS_CANCEL {
				c.jsonResult(enums.JRCodeFailed, "订单状态不能是以完成或已取消", m.KeyId)
			}
		case conf.OTC_APPEAL_STATUS_VALIDATE:
			//确认放币
			if order.Status == conf.OTC_ORDER_STATUS_WAIT {
				c.jsonResult(enums.JRCodeFailed, "订单处于待支付状态", m.KeyId)
			}
			if order.Status == conf.OTC_ORDER_STATUS_FOUND {
				c.jsonResult(enums.JRCodeFailed, "订单处于已完成状态", m.KeyId)
			}
			if order.Status == conf.OTC_ORDER_STATUS_CANCEL {
				c.jsonResult(enums.JRCodeFailed, "订单处于取消状态", m.KeyId)
			}
		}
		//事务校验处理
		o := orm.NewOrm()
		err = o.Begin()
		if err != nil {
			c.jsonResult(enums.JRCodeFailed, "申诉订单事务处理失败", m.KeyId)
		}
		//日志
		logData = map[string]interface{}{
			"操作":    "申诉审核",
			"申诉状态":  obj.Status,
			"admin": c.curUser.Id,
			"申诉表":   obj.KeyId,
			"申诉订单表": obj.OrderId,
			"订单表状态": order.Status,
			"订单类型":  order.OrderType,
		}

		//订单状态操作
		switch obj.Status {
		case conf.OTC_APPEAL_STATUS_CANCEL:
			//取消放币
			order.Status = conf.OTC_ORDER_STATUS_CANCEL //取消订单状态
			order.CancelOrderTime = time.Now()          //订单取消时间
			_, err := o.Update(order)
			if err != nil {
				msg = "取消放币 订单状态修改失败"
				goto Err
			}
			//获取委托单管理
			otcEntrust := models.OtcEntrustOrderOneByKid2(o, order.EntrustOrderId)
			if otcEntrust == nil {
				msg = "获取委托管理订单失败"
				goto Err
			}

			//订单为卖单
			if order.OrderType == conf.OTC_ORDER_TYPE_VENDOR {
				//获取卖出委托冻结金表
				otcUserFrozenBalance := models.OtcUserFrozenBalanceOneByRelevanceIdAndType2(o, order.OrderId, 2)
				if otcUserFrozenBalance == nil {
					msg = "卖出委托冻结金表 信息获取失败"
					goto Err
				}
				//状态改0
				otcUserFrozenBalance.Status = 0
				//o.Using(models.OtcUserFreezeBalanceTBName())
				_, err = o.Update(otcUserFrozenBalance)
				if err != nil {
					msg = "更新 卖出委托冻结金表 失败"
					goto Err
				}

				//获取用户资产表
				tokenskyUserBalance := models.GetTokenskyUserBalanceByUidCoinType2(o, order.VendorUser.UserId, otcEntrust.CoinType)
				if tokenskyUserBalance == nil {
					msg = "用户资产表获取失败"
					goto Err
				}

				logData["卖方用户原冻结资产"] = tokenskyUserBalance.FrozenBalance
				logData["卖方用户原资产"] = tokenskyUserBalance.Balance
				tokenskyUserBalance.FrozenBalance = utils.Float64Sub(tokenskyUserBalance.FrozenBalance, otcUserFrozenBalance.FrozenBalance)
				if tokenskyUserBalance.FrozenBalance < 0 {
					//资产负数
					msg = "卖方用户冻结资产 为负数"
					goto Err
				}

				logData["卖方用户冻结资产"] = tokenskyUserBalance.FrozenBalance
				logData["卖方用户资产"] = tokenskyUserBalance.Balance
				//o.Using(models.TokenskyUserBalanceTBName())

				balanceChange.Add(tokenskyUserBalance.UserId,
					tokenskyUserBalance.CoinType,obj.OrderId,
					"",0,
					conf.CHANGE_SUB,otcUserFrozenBalance.FrozenBalance)

				//_, err = o.Update(tokenskyUserBalance)
				//if err != nil {
				//	msg = "用户资产表保存失败"
				//	goto Err
				//}
			}
			//委托单管理表
			otcEntrust.QuantityLeft = utils.Float64Add(otcEntrust.QuantityLeft, order.Quantity)
			_, err = o.Update(otcEntrust)
			if err != nil {
				msg = "委托订单管理表修改失败"
				goto Err
			}
			//日志
			logData["交易资产额度"] = order.Quantity
			logData["卖方用户uid"] = order.VendorUser.UserId
			logData["买方用户uid"] = order.VendeeUser.UserId
			logData["货币类型"] = otcEntrust.CoinType
			logData["结果"] = "取消放币"
			logData["委托单号"] = otcEntrust.KeyId

		case conf.OTC_APPEAL_STATUS_VALIDATE:
			//确认放币 不扣卖家手续费版

			//获取委托单表
			otcEntrust := models.OtcEntrustOrderOneByKid2(o, order.EntrustOrderId)
			if otcEntrust == nil {
				msg = "委托单表获取失败"
				goto Err
			}
			//订单状态只有3 已付款状态才可以修改
			if order.Status != conf.OTC_ORDER_STATUS_APPEAL {
				msg = "订单状态非3(已申诉)，不可修改"
				goto Err
			}
			//卖方用户资产表
			tokenskyUserBalance := models.GetTokenskyUserBalanceByUidCoinType2(o, order.VendorUser.UserId, otcEntrust.CoinType)
			if tokenskyUserBalance == nil {
				msg = "用户资产表获取失败"
				goto Err
			}
			////用户资产判断
			//if tokenskyUserBalance.FrozenBalance < order.Quantity+order.Quantity*otcEntrust.VendorServiceCharge {
			//	msg = "卖方用户资产不足"
			//	goto Err
			//}

			//扣除卖家币  包含手续费
			logData["卖方用户交易前资产"] = tokenskyUserBalance.Balance
			logData["卖方用户交易前冻结资产"] = tokenskyUserBalance.FrozenBalance
			logData["卖方用户所扣除费用"] = order.Quantity
			tokenskyUserBalance.FrozenBalance = utils.Float64Sub(tokenskyUserBalance.FrozenBalance, order.Quantity)
			if tokenskyUserBalance.FrozenBalance < 0 {
				//资产负数
				msg = "卖方用户冻结资产 为负数"
				goto Err
			}
			tokenskyUserBalance.Balance = utils.Float64Sub(tokenskyUserBalance.Balance, order.Quantity)
			if tokenskyUserBalance.Balance < 0 {
				//资产负数
				msg = "卖方用户资产 为负数"
				goto Err
			}
			logData["卖方用户交易后资产"] = tokenskyUserBalance.Balance
			logData["卖方用户交易后冻结资产"] = tokenskyUserBalance.FrozenBalance
			//o.Using(models.TokenskyUserBalanceTBName())
			//if _, err := o.Update(tokenskyUserBalance); err != nil {
			//	msg = "更新卖家用户资产表失败"
			//	goto Err
			//}
			balanceChange.Add(tokenskyUserBalance.UserId,
				tokenskyUserBalance.CoinType,obj.OrderId,
				conf.CHANGE_SUB,order.Quantity,
				conf.CHANGE_SUB,order.Quantity)

			//买家得到币 扣除手续费 有修改 没有便新增
			buyTokenskyUserBalance := models.GetTokenskyUserBalanceByUidCoinType2(o, order.VendeeUser.UserId, otcEntrust.CoinType)
			buyBalanceFee := utils.Float64Mul(order.Quantity, otcEntrust.VendeeServiceCharge)
			buyBalance := utils.Float64Sub(order.Quantity, buyBalanceFee)
			if buyBalance < 0 {
				//资产负数
				msg = "买方用户得到资产为负数"
				goto Err
			}

			logData["买方用户所扣手续费用"] = buyBalanceFee
			logData["买方用户得到费用"] = buyBalance
			if buyTokenskyUserBalance == nil {
				//新增
				buyTokenskyUserBalance = &models.TokenskyUserBalance{
					KeyId:         0,
					UserId:        order.VendeeUser.UserId,
					CoinType:      otcEntrust.CoinType,
					FrozenBalance: buyBalance,
					CreateTime:    time.Now(),
				}
				logData["买方用户交易前资产"] = 0
				logData["买方用户交易后资产"] = buyBalance
				//if _, err := o.Insert(buyTokenskyUserBalance); err != nil {
				//	msg = "新增买加家用户资产表失败"
				//	goto Err
				//}
				balanceChange.Add(buyTokenskyUserBalance.UserId,
					tokenskyUserBalance.CoinType,obj.OrderId,
					conf.CHANGE_ADD,order.Quantity,
					"",0)

			} else {
				//更新
				logData["买方用户交易前资产"] = buyTokenskyUserBalance.Balance
				buyTokenskyUserBalance.Balance = utils.Float64Add(buyTokenskyUserBalance.Balance, buyBalance)
				logData["买方用户交易后资产"] = buyTokenskyUserBalance.Balance
				//if _, err := o.Update(buyTokenskyUserBalance); err != nil {
				//	msg = "更新买加家用户资产表失败"
				//	goto Err
				//}
				balanceChange.Add(buyTokenskyUserBalance.UserId,
					tokenskyUserBalance.CoinType,obj.OrderId,
					conf.CHANGE_ADD,order.Quantity,
					"",0)
			}
			//订单状态修改
			order.Status = 1
			order.SendCoinTime = time.Now()
			if _, err := o.Update(order); err != nil {
				msg = "更新订单状态失败"
				goto Err
			}
			//买方记录
			vendeeLog := &models.TokenskyTransactionRecord{
				CoinType: otcEntrust.CoinType, //货币类型
				TranType: "OTC买入",
				PushTime: time.Now(), //时间
				Category: 1,          //1收入 2支出
				//UserId:      order.VendeeUser.UserId,
				User:              &models.TokenskyUser{UserId: order.VendeeUser.UserId},
				Money:             order.Quantity, //买方新增
				Status:            1,              //0确认中 1已完成
				RelevanceId:       order.OrderId,
				RelevanceCategory: "otcOrder",
			}
			//买方手续费
			vendeeFeeLog := &models.TokenskyTransactionRecord{
				CoinType: otcEntrust.CoinType, //货币类型
				TranType: "OTC买入手续费",
				PushTime: time.Now(), //时间
				Category: 2,          //1收入 2支出
				//UserId:      order.VendeeUser.UserId,
				User:              &models.TokenskyUser{UserId: order.VendeeUser.UserId},
				Money:             buyBalanceFee, //买方新增
				Status:            1,             //0确认中 1已完成
				RelevanceId:       order.OrderId,
				RelevanceCategory: "otcOrder",
			}
			//卖方记录
			vendorLog := &models.TokenskyTransactionRecord{
				CoinType:          otcEntrust.CoinType, //货币类型
				TranType:          "OTC卖出",
				PushTime:          time.Now(), //时间
				Category:          2,
				UserId:            order.VendorUser.UserId,
				User:              &models.TokenskyUser{UserId: order.VendorUser.UserId},
				Money:             order.Quantity, //卖方扣除费用
				Status:            1,
				RelevanceId:       order.OrderId,
				RelevanceCategory: "otcOrder",
			}
			_, err = o.InsertMulti(4, []*models.TokenskyTransactionRecord{vendeeLog, vendeeFeeLog, vendorLog})
			if err != nil {
				msg = "新增交易记录失败"
				goto Err
			}
			//是否自动取消订单
			if !c.OtcEntrustOrderCallOff(o, order, otcEntrust, tokenskyUserBalance) {
				msg = "取消委托订单失败"
				goto Err
			}
			//日志
			logData["交易资产额度"] = order.Quantity
			logData["卖方用户uid"] = order.VendorUser.UserId
			logData["买方用户uid"] = order.VendeeUser.UserId
			logData["货币类型"] = otcEntrust.CoinType
			logData["结果"] = "确认放币"
			logData["委托单号"] = otcEntrust.KeyId
		}

		_, err = o.Update(obj)
		if err != nil {
			msg = "更新申诉表失败"
			goto Err
		}
		//用户资产变动处理
		ok,msg,tx = balanceChange.Send()
		if !ok{
			goto Err
		}
		ok = models.TokenskyUserBalanceHashSetStatus(o,tx)
		if !ok{
			msg = "设置哈希表异常"
			goto Err
		}
		err = o.Commit()
		if err != nil {
			msg = "事务处理失败"
			goto Err
		}

		//日志
		logStr, _ = json.Marshal(logData)
		utils.LogNotice(logStr)
		//消息推送

		switch obj.Status {
		case conf.OTC_APPEAL_STATUS_CANCEL:
			//取消放币
			models.TokenskyJiguangRegistrationidSendByIds([]int{obj.Order.VendorUser.UserId,obj.Order.VendeeUser.UserId},
			"订单号"+order.OrderId+"审核结果","取消放币","订单号"+order.OrderId+"审核结果","取消放币")
		case conf.OTC_APPEAL_STATUS_VALIDATE:
			//确认放币
			models.TokenskyJiguangRegistrationidSendByIds([]int{obj.Order.VendorUser.UserId,obj.Order.VendeeUser.UserId},
			"订单号"+order.OrderId+"审核结果","确认放币","订单号"+order.OrderId+"审核结果","确认放币")
		}
		c.jsonResult(enums.JRCodeSucc, "编辑成功", obj.KeyId)
	Err:
		if err := o.Rollback(); err != nil {
			msg = "回滚失败"
		}
		c.jsonResultError(enums.JRCodeFailed, msg, obj.KeyId)
	}
}

//委托单自动取消
func (c *OtcAppealController) OtcEntrustOrderCallOff(o orm.Ormer, order *models.OtcOrder, orderEntrust *models.OtcEntrustOrder, userBalcnce *models.TokenskyUserBalance) bool {
	//1 otc_entrust_order表中的 quantity_left 要小于 min
	if orderEntrust.QuantityLeft >= orderEntrust.Min {
		//
		return true
	}

	//2 没有正在进行交易的订单
	// 查询otc_order表  并且order_id不等于当前order_id(否则永远都是在进行中)
	if obj := models.OtcOrderGetOrderByOidAndBidIsStatusNot1or4(o, order.OrderId, orderEntrust.KeyId); obj != nil {
		return true
	}

	switch order.OrderType {
	case 1:
		serviceCharge := utils.Float64Mul(orderEntrust.VendeeServiceCharge, orderEntrust.QuantityLeft)
		//扣除卖家手续费版本
		//sumMoney := utils.Float64Add(serviceCharge, orderEntrust.QuantityLeft, conf.FLOAT_PRECISE_NUM_6)
		//不扣卖家手续费版
		sumMoney := orderEntrust.QuantityLeft
		//买单
		OtcEntrustAutoRecord := &models.OtcEntrustAutoCancelRecord{
			EntrustType:    orderEntrust.EntrustType,
			EntrustOrderId: orderEntrust.KeyId,
			Money:          orderEntrust.QuantityLeft,
			ServiceCharge:  serviceCharge,
			SumMoney:       sumMoney,
			UserId:         orderEntrust.User.UserId,
		}
		if _, err := o.Insert(OtcEntrustAutoRecord); err != nil {
			return false
		}
		//取消时间
		orderEntrust.AutoCancelTime = time.Now()
		orderEntrust.Status = 3
		if _, err := o.Update(orderEntrust); err != nil {
			return false
		}
		//扣除用户冻结资产
		if userBalcnce != nil && orderEntrust.User.UserId == userBalcnce.UserId {
			userBalcnce.FrozenBalance = utils.Float64Sub(userBalcnce.FrozenBalance, sumMoney)
			if userBalcnce.FrozenBalance < 0 {
				return false
			}
			if _, err := o.Update(userBalcnce); err != nil {
				return false
			}
		}
	case 2:
		//卖单
		//serviceCharge := utils.Float64Mul(orderEntrust.QuantityLeft, orderEntrust.VendorServiceCharge, conf.FLOAT_PRECISE_NUM_6)
		//	sumMoney := utils.Float64Add(serviceCharge, orderEntrust.QuantityLeft, conf.FLOAT_PRECISE_NUM_6)
		OtcEntrustAutoRecord := &models.OtcEntrustAutoCancelRecord{
			EntrustType:    orderEntrust.EntrustType,
			EntrustOrderId: orderEntrust.KeyId,
			//Money:          orderEntrust.QuantityLeft,
			//ServiceCharge:  serviceCharge,
			//SumMoney:       sumMoney,
			UserId: orderEntrust.User.UserId,
		}
		if _, err := o.Insert(OtcEntrustAutoRecord); err != nil {
			return false
		}
		//取消时间
		orderEntrust.Status = 3
		orderEntrust.AutoCancelTime = time.Now()
		if _, err := o.Update(orderEntrust); err != nil {
			return false
		}
		////扣除用户冻结资产
		//if userBalcnce != nil && orderEntrust.User.UserId == userBalcnce.UserId{
		//	userBalcnce.FrozenBalance = utils.Float64Sub(userBalcnce.FrozenBalance, sumMoney, conf.FLOAT_PRECISE_NUM_6)
		//	if userBalcnce.FrozenBalance < 0 {
		//		return false
		//	}
		//	if _, err := o.Update(userBalcnce); err != nil {
		//		return false
		//	}
		//}
	}
	return true
}
