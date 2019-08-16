package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
	"tokensky_bg_admin/common"
	"tokensky_bg_admin/conf"
	"tokensky_bg_admin/enums"
	"tokensky_bg_admin/models"
	"tokensky_bg_admin/utils"
)

//提币审核
type TokenskyUserTibiController struct {
	BaseController
}

//Prepare 参考beego官方文档说明
func (c *TokenskyUserTibiController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	c.checkAuthor()
}

// DataGrid
func (c *TokenskyUserTibiController) DataGrid() {
	//直接反序化获取json格式的requestbody里的值
	var params models.TokenskyUserTibiQueryParam
	switch c.Ctx.Request.Method {
	case "POST":
		json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	case "GET":
		params.Sort = c.GetString("sort")
		params.Order = c.GetString("order")
		params.Limit, _ = c.GetInt64("limit")
		params.Offset, _ = c.GetInt64("offset")
		params.Status = c.GetString("status")
		params.StartTime, _ = c.GetInt64("startTime")
		params.EndTime, _ = c.GetInt64("endTime")
		params.Phone = c.GetString("phone")
		params.CoinType = c.GetString("coinType")
		params.OutAddress = c.GetString("outAddress")
		params.InAddress = c.GetString("inAddress")
	}
	//获取数据列表和总数
	data, total := models.TokenskyUserTibiPageList(&params)
	//定义返回的数据结构
	mapp := map[string]interface{}{
		"rows":  data,
		"total": total,
	}
	c.jsonResult(enums.JRCodeSucc, "", mapp)
}

//审核
func (c *TokenskyUserTibiController) Examine() {
	mapp := make(map[string]string)
	body := c.Ctx.Input.RequestBody //接收raw body内容
	if err := json.Unmarshal(body, &mapp); err != nil {
		c.jsonResult(enums.JRCodeFailed, "数据异常", 0)
	}
	ids := strings.Split(mapp["ids"], ",")

	var status int
	if con, err := strconv.Atoi(mapp["status"]); err == nil {
		status = con
	} else {
		c.jsonResult(enums.JRCodeFailed, "数据异常", 0)
	}
	data, total := models.TokenskyUserTibiByIds(ids)
	if total == 0 {
		c.jsonResult(enums.JRCodeFailed, "数据不存在", 0)
	}

	//1状态校验
	var con int
	for i, obj := range data {
		switch i {
		case 0:
			con = obj.Status
		default:
			if con != obj.Status {
				c.jsonResult(enums.JRCodeFailed, "数据状态不统一", obj.KeyId)
			}
		}
	}
	//提币
	switch status {
	case 1, 3:
		//审核通过
		client := &http.Client{}

		for _, obj := range data {
			url := conf.JIANG_SERVER_URL + "/service/" + obj.CoinType + "/withdraw"
			if obj.Quantity < obj.ServiceCharge {
				//提币数量小于手续费
				continue
			}
			//开启事务
			o := orm.NewOrm()
			err := o.Begin()
			if err != nil {
				//开启事务失败
				c.jsonResultError2(enums.JRCodeError, "开启事务失败")
			}
			//状态只有0 和 3 可以继续
			obj = models.TokenskyUserTibiById2(o, obj.KeyId)
			if obj.Status != 0 && obj.Status != 3 {
				o.Rollback()
				continue
			}

			//withdraw:
			//	amount 提币数量
			//	withdraw_id 请求发起方填的提现id，不能重复
			//	to 提现目标地址

			mapp := map[string]interface{}{
				"withdraw_id": obj.KeyId,                                      //请求的提现id
				"amount":      strconv.FormatFloat(obj.Quantity, 'E', -1, 64), //数额
				"to":          obj.InAddress,                                  //地址
			}
			bytesData, err := json.Marshal(mapp)
			reader := bytes.NewReader(bytesData)
			request, err := http.NewRequest("POST", url, reader)
			if err != nil && request != nil {
				o.Rollback()
				c.jsonResult(enums.JRCodeFailed, "创建请求对象失败", obj.KeyId)
			}
			request.Header.Set("Content-type", "application/json")
			response, err := client.Do(request)
			if err != nil {
				o.Rollback()
				c.jsonResult(enums.JRCodeFailed, "连接API失败", obj.KeyId)
			}
			body, err := ioutil.ReadAll(response.Body)
			//数据返回异常
			if err != nil {
				o.Rollback()
				c.jsonResult(enums.JRCodeFailed, "加载数据失败", obj.KeyId)
			}
			req := make(map[string]interface{})
			err = json.Unmarshal(body, &req)
			//状态
			code2, ok := req["code"].(float64)
			code := int(code2)
			//数据解析失败
			if !ok {
				o.Rollback()
				c.jsonResult(enums.JRCodeFailed, "审核表返回状态解析异常", obj.KeyId)
			}
			//通过人Id记录
			if obj.AdminId == 0 {
				obj.AdminId = c.curUser.Id
				obj.VerifyTime = time.Now()
			}

			/*提币审核日志*/
			logData := map[string]interface{}{
				"操作":    "提币审核(通过)",
				"审核状态":  obj.Status,
				"admin": c.curUser.Id,
				"用户":    obj.User.UserId,
				"提币单号":  obj.OrderId,
				"提币类型":  obj.CoinType,
				"提币金额":  obj.Quantity,
				"提币手续费": obj.ServiceChargeQuantity,
				"响应状态":  code,
			}
			logStr, _ := json.Marshal(logData)
			utils.LogNotice(logStr)
			//状态码
			switch code {
			case conf.TIBI_ERR_SUCCESS:
				//成功
				obj.Status = 3
				if _, err := o.Update(obj); err != nil {
					o.Rollback()
					c.jsonResult(enums.JRCodeFailed, "状态更新失败", obj.KeyId)
				}
				o.Commit()
				//极光推送状态
				models.TokenskyJiguangRegistrationidSendByOne(obj.User.UserId,"提币审核通过",obj.OrderId,"提币审核通过",obj.OrderId)
			case conf.TIBI_ERR_BAD_PARAMETER:
				//参数错误
				obj.Status = 4
				if _, err := o.Update(obj); err != nil {
					o.Rollback()
					c.jsonResult(enums.JRCodeFailed, "状态更新失败", obj.KeyId)
				}
				o.Commit()
				c.jsonResult(enums.JRCodeFailed, "参数异常", obj.KeyId)
			case conf.TIBI_ERR_SERVER_ERROR:
				//服务器错误
				o.Rollback()
				c.jsonResult(enums.JRCodeFailed, "服务器错误，请稍等尝试", obj.KeyId)
			case conf.TIBI_ERR_WITHDRAW_EXISTED:
				//重复申请
				o.Rollback()
				c.jsonResult(enums.JRCodeFailed, "提币订单已在申请中", obj.KeyId)
			case conf.TIBI_ERR_INVALID_ADDRESS:
				//无效地址
				o.Rollback()
				c.jsonResult(enums.JRCodeFailed, "无效地址", obj.KeyId)
			case conf.TIBI_ERR_INVALID_AMOUNT:
				//无效金额
				o.Rollback()
				c.jsonResult(enums.JRCodeFailed, "无效金额", obj.KeyId)
			default:
				o.Rollback()
				if _, err := o.Update(obj); err != nil {
					o.Rollback()
					c.jsonResult(enums.JRCodeFailed, "状态更新失败", obj.KeyId)
				}
				o.Commit()
				c.jsonResult(enums.JRCodeFailed, "未知异常", obj.KeyId)
			}
		}
	case 2:
		//审核未通过
		for _, obj := range data {
			//审核中的无法选择审核未通过
			switch obj.Status {
			case 3:
				continue
			}
			//通过人Id记录
			if obj.AdminId == 0 {
				obj.AdminId = c.curUser.Id
				obj.VerifyTime = time.Now()
			}
			//用户资产表 资产一定存在
			o := orm.NewOrm()
			err := o.Begin()
			if err != nil {
				c.jsonResult(enums.JRCodeFailed, "开启事务失败", 0)
			}
			balance := models.GetTokenskyUserBalanceByUidCoinType2(o, obj.User.UserId, obj.CoinType)
			if balance == nil {
				//用户资产表一定存在
				o.Rollback()
				continue
			}
			//
			obj = models.TokenskyUserTibiById2(o, obj.KeyId)
			if obj == nil {
				o.Rollback()
				continue
			}
			//状态校验
			if obj.Status != 0 {
				//已经通过审核过
				continue
			}
			//提币前资产
			logBalance := balance.Balance
			logFrozenBalance := balance.FrozenBalance
			//资产修改
			balanceChange := common.NewTokenskyUserBalanceChange(2,"tibi","提币")

			balance.FrozenBalance = utils.Float64Sub(balance.FrozenBalance, obj.SumQuantity)
			if balance.FrozenBalance < 0 {
				//资产负数
				o.Rollback()
				c.jsonResultError2(enums.JRCodeError, "扣除用户冻结资产，为负数")
			}

			obj.Status = 2
			//if _, err := o.Update(balance); err != nil {
			//	//资产表更新失败
			//	o.Rollback()
			//	c.jsonResult(enums.JRCodeFailed, "资产更新失败", obj.KeyId)
			//}
			if _, err := o.Update(obj); err != nil {
				//审核表保存失败
				if err := o.Rollback(); err != nil {
					c.jsonResult(enums.JRCodeFailed, "审核回滚失败", obj.KeyId)
				}
				c.jsonResult(enums.JRCodeFailed, "审核表保存失败", obj.KeyId)
			}
			balanceChange.Add(balance.UserId,balance.CoinType,
				obj.OrderId,"",0,
				conf.CHANGE_SUB,obj.SumQuantity,
				)

			//修改记录状态
			feeLog := models.TokenskyTransactionRecordOneByRelevance(o, "提币", "tibi", obj.OrderId)
			if feeLog == nil {
				//审核表保存失败
				if err := o.Rollback(); err != nil {
					c.jsonResult(enums.JRCodeFailed, "获取记录表回滚失败", obj.KeyId)
				}
				c.jsonResult(enums.JRCodeFailed, "获取记录表失败", obj.KeyId)
			}
			feeLog.Status = 2
			if _, err := o.Update(feeLog); err != nil {
				//记录表保存失败
				if err := o.Rollback(); err != nil {
					c.jsonResult(enums.JRCodeFailed, "记录表回滚失败", obj.KeyId)
				}
			}
			feeLog2 := models.TokenskyTransactionRecordOneByRelevance(o, "提币手续费", "tibi", obj.OrderId)
			if feeLog2 == nil {
				//审核表保存失败
				if err := o.Rollback(); err != nil {
					c.jsonResult(enums.JRCodeFailed, "获取记录表回滚失败2", obj.KeyId)
				}
				c.jsonResult(enums.JRCodeFailed, "获取记录表失败2", obj.KeyId)
			}
			feeLog2.Status = 2
			if _, err := o.Update(feeLog2); err != nil {
				//记录表保存失败
				if err := o.Rollback(); err != nil {
					c.jsonResult(enums.JRCodeFailed, "记录表回滚失败", obj.KeyId)
				}
			}
			//
			ok,msg,tx := balanceChange.Send()
			if !ok{
				o.Rollback()
				c.jsonResult(enums.JRCodeFailed, msg, obj.KeyId)
			}
			ok = models.TokenskyUserBalanceHashSetStatus(o,tx)
			if !ok{
				o.Rollback()
				c.jsonResult(enums.JRCodeFailed, "哈希表异常", obj.KeyId)
			}
			if err := o.Commit(); err != nil {
				c.jsonResult(enums.JRCodeFailed, "事务提交失败", obj.KeyId)
			}

			/* 提币详情日志 */
			logData := map[string]interface{}{
				"操作":      "提币审核(未通过)",
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
		}
	default:
		c.jsonResult(enums.JRCodeFailed, "状态异常", 0)
	}
	c.jsonResult(enums.JRCodeSucc, "审核完成", 0)
}
