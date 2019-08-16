package models

import (
	//"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"tokensky_bg_admin/conf"
)

// init 初始化
func init() {
	//admin
	orm.RegisterModel(new(AdminBackendUser), new(AdminResource), new(AdminRole), new(AdminRoleResourceRel), new(AdminRoleBackendUserRel))
	orm.RegisterModel(new(AdminModelRecord))
	orm.RegisterModel(new(TokenskyMessage), new(TokenskyUser), new(TokenskyUserBalance), new(TokenskyAccountBank))
	orm.RegisterModel(new(TokenskyRealAuth), new(TokenskyTransactionRecord), new(TokenskyUserElectricityBalance))
	orm.RegisterModel(new(TokenskyTibiConfig), new(TokenskyTibiConfigBak), new(TokenskyChongbiConfig), new(TokenskyChongbiConfigBak))
	orm.RegisterModel(new(TokenskyUserTibi), new(TokenskyUserAddress), new(TokenskyUserDeposit), new(TokenskyOrderIds), new(TokenskyJiguangRegistrationid))
	orm.RegisterModel(new(TokenskyUserBalancesRecord),new(TokenskyUserBalanceHash),new(TokenskyUserBalanceCoin))
	//OCT
	orm.RegisterModel(new(OtcConf), new(OtcConfBak), new(OtcEntrustOrder), new(OtcOrder), new(OtcAppeal))
	orm.RegisterModel(new(OtcUserFrozenBalance), new(OtcEntrustAutoCancelRecord))
	//Operation
	orm.RegisterModel(new(OperationBanner))
	//角色管理
	orm.RegisterModel(new(RoleBlackList))
	//算力 Hashrate
	orm.RegisterModel(new(HashrateCategory), new(HashrateOrder), new(HashrateOrderTransaction), new(HashrateTreaty), new(HashrateOrderProfit))
	orm.RegisterModel(new(HashrateSendBalanceRecord))
	//理财
	orm.RegisterModel(new(FinancialProduct),new(FinancialCategory),new(FinancialProductHistoricalRecord))
	orm.RegisterModel(new(FinancialLiveUserBalance),new(FinancialOrder),new(FinancialProfit))
	orm.RegisterModel(new(FinancialOrderWithdrawal))
	//借贷
	orm.RegisterModel(new(BorrowConf),new(BorrowOrder),new(BorrowUseFinancialOrder),new(BorrowLimiting),new(BorrowOrdeLog))
	//其它
	orm.RegisterModel(new(TokenskyUserInvite))
	//爬虫
	orm.RegisterModel(new(SpiderCoinMarket))
}

/*Admin相关*/

// AdminBackendUserTBName 获取 AdminBackendUser 对应的表名称
func AdminBackendUserTBName() string {
	return conf.DB_ADMIN_DT_PREFIX + "backend_user"
}

// AdminResourceTBName 获取 AdminResource 对应的表名称
func AdminResourceTBName() string {
	return conf.DB_ADMIN_DT_PREFIX + "resource"
}

// AdminRoleTBName 获取 AdminRole 对应的表名称
func AdminRoleTBName() string {
	return conf.DB_ADMIN_DT_PREFIX + "role"
}

// AdminRoleResourceRelTBName 角色与资源多对多关系表
func AdminRoleResourceRelTBName() string {
	return conf.DB_ADMIN_DT_PREFIX + "role_resource_rel"
}

// AdminRoleBackendUserRelTBName 角色与用户多对多关系表
func AdminRoleBackendUserRelTBName() string {
	return conf.DB_ADMIN_DT_PREFIX + "role_backenduser_rel"
}

//用户修改记录表
func AdminModelRecordTBName()string{
	return conf.DB_ADMIN_DT_PREFIX + "model_record"
}

/*Tokensky相关*/

// TokenskyMessageTBName 获取TokenskyMessage对应多名称
func TokenskyMessageTBName() string {
	return conf.DB_TOKENSKY_DT_PREFIX + "message"
}

// TokenskyUserTBName 获取TokenskyUser对应多名称
func TokenskyUserTBName() string {
	return conf.DB_TOKENSKY_DT_PREFIX + "user"
}
//TokenskyUserBalanceCoinTBName
func TokenskyUserBalanceCoinTBName()string {
	return conf.DB_TOKENSKY_DT_PREFIX + "user_balance_coin"
}
//TokenskyUserBalance 用户资产表
func TokenskyUserBalanceTBName() string {
	return conf.DB_TOKENSKY_DT_PREFIX + "user_balance"
}

//用户资产记录表
func TokenskyUserBalancesRecordTBName()string{
	return conf.DB_TOKENSKY_DT_PREFIX + "user_balance_record"
}

//交易记录明细表
func TokenskyTransactionRecordTBName() string {
	return conf.DB_TOKENSKY_DT_PREFIX + "transaction_record"
}

//用户付款设置表
func TokenskyAccountBankTBName() string {
	return conf.DB_TOKENSKY_DT_PREFIX + "account_bank"
}

//提币配置表
func TokenskyTibiConfigTBName() string {
	return conf.DB_TOKENSKY_DT_PREFIX + "tibi_config"
}

//提币配置bak表
func TokenskyTibiConfigBakTBName() string {
	return conf.DB_TOKENSKY_DT_PREFIX + "tibi_config_bak"
}

//充币配置表
func TokenskyChongbiConfigTBName() string {
	return conf.DB_TOKENSKY_DT_PREFIX + "chongbi_config"
}

//充币配置bak表
func TokenskyChongbiConfigBakTBName() string {
	return conf.DB_TOKENSKY_DT_PREFIX + "chongbi_config_bak"
}

//提币审核表
func TokenskyUserTibiTBName() string {
	return conf.DB_TOKENSKY_DT_PREFIX + "user_tibi"
}

//用户充值地址表
func TokenskyUserAddressTBName() string {
	return conf.DB_TOKENSKY_DT_PREFIX + "user_address"
}

//用户充值记录表
func TokenskyUserDepositTBName() string {
	return conf.DB_TOKENSKY_DT_PREFIX + "user_deposit"
}

//用户资产哈希表
func TokenskyUserBalanceHashTBName()string{
	return conf.DB_TOKENSKY_DT_PREFIX + "user_balance_hash"
}

// TokenskyUserElectricityBalance用户电力资产表
func TokenskyUserElectricityBalanceTBName() string {
	return conf.DB_TOKENSKY_DT_PREFIX + "user_electricity_balance"
}

//名称表
func TokenskyOrderIdsTBName() string {
	return conf.DB_TOKENSKY_DT_PREFIX + "order_ids"
}

//极光地址表
func TokenskyJiguangRegistrationidTBName() string {
	return conf.DB_TOKENSKY_DT_PREFIX + "jiguang_registrationid"
}

//用户邀请表
func TokenskyUserInviteTBName()string{
	return conf.DB_TOKENSKY_DT_PREFIX + "user_invite"
}

/*Oct相关*/

//OtcConfTBName 获取OctConf对应多名称
func OtcConfTBName() string {
	return conf.DB_OTC_DT_PREFIX + "conf"
}

//OtcConf 记录副表
func OtcConfBakTBName() string {
	return conf.DB_OTC_DT_PREFIX + "conf_bak"
}

//OtcEntrustOrder 委托订单表
func OtcEntrustOrderTBName() string {
	return conf.DB_OTC_DT_PREFIX + "entrust_order"
}

// OtcOrder订单管理表
func OtcOrderTBName() string {
	return conf.DB_OTC_DT_PREFIX + "order"
}

// OtcAppeal 订单申诉表
func OtcAppealTBName() string {
	return conf.DB_OTC_DT_PREFIX + "appeal"
}

//卖出委托订单表
func OtcUserFrozenBalanceTBName() string {
	return conf.DB_OTC_DT_PREFIX + "user_frozen_balance"
}

//委托单记录取消表
func OtcEntrustAutoCancelRecordTBName() string {
	return conf.DB_OTC_DT_PREFIX + "entrust_auto_cancel_record"
}

/*算力相关*/

//算力合约表
func HashrateCategoryTBName() string {
	return conf.DB_HASHRATE_DT_PREFIX + "category"
}

//算力订单交易关联表
func HashrateOrderTransactionTBName() string {
	return conf.DB_HASHRATE_DT_PREFIX + "order_transaction"
}

//算力订单表
func HashrateOrderTBName() string {
	return conf.DB_HASHRATE_DT_PREFIX + "order"
}

//算力合约表
func HashrateTreatyTBName() string {
	return conf.DB_HASHRATE_DT_PREFIX + "treaty"
}

//算力订单收益表
func HashrateOrderProfitTBName() string {
	return conf.DB_HASHRATE_DT_PREFIX + "order_profit"
}

// 算力奖励发放记录表
func HashrateSendBalanceRecordTBName() string {
	return conf.DB_HASHRATE_DT_PREFIX + "send_balance_record"
}

/*运营相关*/

//OperationBanner Banner表
func OperationBannerTBName() string {
	return conf.DB_OPERATION_DT_PREFIX + "banner"
}

/*用户管理相关*/

//黑名单表
func RoleBlackListTBName() string {
	return conf.DB_ROLE_DT_PREFIX + "black_list"
}

//身份审核
func TokenskyRealAuthTBName() string {
	return conf.DB_TOKENSKY_DT_PREFIX + "real_auth"
}

/* 财务管理 */

//财务类型表
func FinancialCategoryTBName() string {
	return conf.DB_FINANCIAL_DT_PREFIX + "category"
}

//财务配置表
func FinancialProductTBName() string {
	return conf.DB_FINANCIAL_DT_PREFIX + "product"
}

//财务收益表
func FinancialProfitTBName() string {
	return conf.DB_FINANCIAL_DT_PREFIX + "profit"
}

//财务配置历史记录表
func FinancialProductHistoricalRecordTBName() string {
	return conf.DB_FINANCIAL_DT_PREFIX + "product_historical_record"
}

//财务订单表
func FinancialOrderTBName() string {
	return conf.DB_FINANCIAL_DT_PREFIX + "order"
}

//财务订单提币表
func FinancialOrderWithdrawalTBName()string{
	return conf.DB_FINANCIAL_DT_PREFIX + "order_withdrawal"
}

//财务资产表
func FinancialLiveUserBalanceTBName() string {
	return conf.DB_FINANCIAL_DT_PREFIX + "live_user_balance"
}

/*借贷*/

//借贷配置表
func BorrowConfTBName()string{
	return conf.DB_BORROW_DT_PREFIX + "conf"
}

//借贷订单表
func BorrowOrderTBName()string{
	return conf.DB_BORROW_DT_PREFIX + "order"
}

//借贷日志表
func BorrowOrdeLogTBName()string{
	return conf.DB_BORROW_DT_PREFIX + "order_log"
}

//借贷关联表
func BorrowUseFinancialOrderTBName()string{
	return conf.DB_BORROW_DT_PREFIX + "use_financial_order"
}

//借贷强屏表
func BorrowLimitingTBName()string{
	return conf.DB_BORROW_DT_PREFIX + "limiting"
}

//爬虫表
func SpiderCoinMarketTBName()string{
	return "spider_coin_market"
}