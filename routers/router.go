// @APIVersion 1.0.0
// @Title mobile API
// @Description mobile has every tool to get any job done, so codename for the new mobile APIs.
// @Contact astaxie@gmail.com
package routers

import (
	"github.com/astaxie/beego"
	"tokensky_bg_admin/controllers"
)

func init() {

	//后台-角色 增删改查
	//beego.Router("/role/index", &controllers.AdminRoleController{}, "*:Index")
	beego.Router("/admin/role/datagrid", &controllers.AdminRoleController{}, "Get,Post:DataGrid")
	beego.Router("/admin/role/edit/?:id", &controllers.AdminRoleController{}, "Get,Post:Edit")
	beego.Router("/admin/role/delete", &controllers.AdminRoleController{}, "Post:Delete")
	beego.Router("/admin/role/datalist", &controllers.AdminRoleController{}, "Post:DataList")
	beego.Router("/admin/role/allocate", &controllers.AdminRoleController{}, "Post:Allocate")
	beego.Router("/admin/role/updateseq", &controllers.AdminRoleController{}, "Post:UpdateSeq")

	//后台-资源 增删改查
	//beego.Router("/resource/index", &controllers.AdminResourceController{}, "*:Index")
	beego.Router("/admin/resource/treegrid", &controllers.AdminResourceController{}, "Get,POST:TreeGrid")
	beego.Router("/admin/resource/treeGridByRole", &controllers.AdminResourceController{}, "Get,POST:TreeGridByRole")
	beego.Router("/admin/resource/edit/?:id", &controllers.AdminResourceController{}, "Get,Post:Edit")
	beego.Router("/admin/resource/parent", &controllers.AdminResourceController{}, "Post:ParentTreeGrid")
	beego.Router("/admin/resource/delete", &controllers.AdminResourceController{}, "Post:Delete")

	//快速修改顺序
	beego.Router("/admin/resource/updateseq", &controllers.AdminResourceController{}, "Post:UpdateSeq")

	//通用选择面板
	beego.Router("/admin/resource/select", &controllers.AdminResourceController{}, "Get:Select")
	//用户有权管理的菜单列表（包括区域）
	beego.Router("/admin/resource/usermenutree", &controllers.AdminResourceController{}, "POST:UserMenuTree")
	beego.Router("/admin/resource/checkurlfor", &controllers.AdminResourceController{}, "POST:CheckUrlFor")

	//后台-用户 增删改查
	//beego.Router("/backenduser/index", &controllers.AdminBackendUserController{}, "*:Index")
	beego.Router("/admin/backenduser/datagrid", &controllers.AdminBackendUserController{}, "GET,POST:DataGrid")
	beego.Router("/admin/backenduser/edit/?:id", &controllers.AdminBackendUserController{}, "Post:Edit")
	beego.Router("/admin/backenduser/delete", &controllers.AdminBackendUserController{}, "Post:Delete")

	//后台用户中心
	beego.Router("/admin/usercenter/profile", &controllers.AdminUserCenterController{}, "Get:Profile")
	beego.Router("/admin/usercenter/basicinfosave", &controllers.AdminUserCenterController{}, "Post:BasicInfoSave")
	beego.Router("/admin/usercenter/uploadimage", &controllers.AdminUserCenterController{}, "Post:UploadImage")
	beego.Router("/admin/usercenter/passwordsave", &controllers.AdminUserCenterController{}, "Post:PasswordSave")

	//beego.Router("/home/index", &controllers.AdminHomeController{}, "*:Index")
	//beego.Router("/home/login", &controllers.AdminHomeController{}, "*:Login")
	beego.Router("/home/dologin", &controllers.AdminHomeController{}, "Get,Post:DoLogin")
	beego.Router("/home/logout", &controllers.AdminHomeController{}, "*:Logout")
	//beego.Router("/home/datareset", &controllers.AdminHomeController{}, "Post:DataReset")

	//用户注册
	//beego.Router("/home/register", &controllers.AdminHomeController{}, "Post:Register")
	//beego.Router("/home/404", &controllers.AdminHomeController{}, "*:Page404")
	//beego.Router("/hofme/error/?:error", &controllers.AdminHomeController{}, "*:Error")

	//钱包-消息 增删改查
	beego.Router("/tokensky/message/datagrid", &controllers.TokenskyMessageController{}, "Get,POST:DataGrid")
	beego.Router("/tokensky/message/edit/?:id", &controllers.TokenskyMessageController{}, "Post:Edit")
	beego.Router("/tokensky/message/delete", &controllers.TokenskyMessageController{}, "Post:Delete")
	//钱包-角色 -查
	beego.Router("/tokensky/user/datagrid", &controllers.TokenskyUserController{}, "Get,POST:DataGrid")
	beego.Router("/tokensky/user/setLevel", &controllers.TokenskyUserController{}, "POST:SetLevel")
	beego.Router("/tokensky/user/invitation", &controllers.TokenskyUserController{}, "POST:SetInvitation")
	beego.Router("/tokensky/user/getAddr", &controllers.TokenskyUserController{}, "POST:GetAddr")
	//钱包-资产-查
	beego.Router("/tokensky/userBalance/Balance", &controllers.TokenskyUserBalanceController{}, "Get,POST:GetBalances")

	//钱包-OTC 配置
	beego.Router("/otc/conf/datagrid", &controllers.OtcConfController{}, "Get,Post:DataGrid")
	beego.Router("/otc/conf/getconf", &controllers.OtcConfController{}, "*:GetConf")
	beego.Router("/otc/conf/edit", &controllers.OtcConfController{}, "Post:Edit")
	//钱包-OTC 委托单
	beego.Router("/otc/entrustOrder/datagrid", &controllers.OtcEntrustOrderController{}, "Get,Post:DataGrid")
	//钱包-OTC 订单
	beego.Router("/otc/order/datagrid", &controllers.OtcOrderController{}, "Get,Post:DataGrid")
	//钱包-OTC 申诉
	beego.Router("/otc/appeal/datagrid", &controllers.OtcAppealController{}, "Get,Post:DataGrid")
	beego.Router("/otc/appeal/examine", &controllers.OtcAppealController{}, "Post:Examine")

	//钱包-角色 黑名单
	beego.Router("/role/blackList/datagrid", &controllers.RoleBlackListController{}, "Get,Post:DataGrid")
	beego.Router("/role/blackList/edit", &controllers.RoleBlackListController{}, "Post:Edit")
	beego.Router("/role/blackList/delete", &controllers.RoleBlackListController{}, "Post:Delete")

	//钱包-角色 身份审核
	beego.Router("/tokensky/realAuth/datagrid", &controllers.TokenskyRealAuthController{}, "Get,Post:DataGrid")
	beego.Router("/tokensky/realAuth/auditing", &controllers.TokenskyRealAuthController{}, "Post:Auditing")

	//钱包-提币配置
	beego.Router("/tokensky/tokenskyTibiConfig/datagrid", &controllers.TokenskyTibiConfigController{}, "Get,Post:DataGrid")
	beego.Router("/tokensky/tokenskyTibiConfig/edit", &controllers.TokenskyTibiConfigController{}, "Post:Edit")

	//钱包-冲币配置
	beego.Router("/tokensky/tokenskyChongbiConfig/datagrid", &controllers.TokenskyChongbiConfigController{}, "Get,Post:DataGrid")
	beego.Router("/tokensky/tokenskyChongbiConfig/edit", &controllers.TokenskyChongbiConfigController{}, "Post:Edit")

	//钱包-提币审核
	beego.Router("/tokensky/tokenskyUserTibi/datagrid", &controllers.TokenskyUserTibiController{}, "Get,Post:DataGrid")
	beego.Router("/tokensky/tokenskyUserTibi/examine", &controllers.TokenskyUserTibiController{}, "Post:Examine")

	//钱包-运营-banner
	beego.Router("/operation/banner/datagrid", &controllers.OperationBannerController{}, "Get,Post:DataGrid")
	beego.Router("/operation/banner/edit", &controllers.OperationBannerController{}, "Post:Edit")
	beego.Router("/operation/banner/delete", &controllers.OperationBannerController{}, "Post:Delete")

	//钱包-算力-算力合约分类
	beego.Router("/hashrate/category/datagrid", &controllers.HashrateCategoryController{}, "Get,Post:DataGrid")
	beego.Router("/hashrate/category/edit", &controllers.HashrateCategoryController{}, "Post:Edit")
	beego.Router("/hashrate/category/delete", &controllers.HashrateCategoryController{}, "Post:Delete")

	//钱包-算力-算力合约
	beego.Router("/hashrate/transaction/datagrid", &controllers.HashrateTreatyController{}, "Get,Post:DataGrid")
	beego.Router("/hashrate/transaction/edit", &controllers.HashrateTreatyController{}, "Post:Edit")
	beego.Router("/hashrate/transaction/delete", &controllers.HashrateTreatyController{}, "Post:Delete")
	beego.Router("/hashrate/transaction/shelves", &controllers.HashrateTreatyController{}, "Post:Shelves")
	beego.Router("/hashrate/transaction/isNotFutures", &controllers.HashrateTreatyController{}, "Post:IsNotFutures")
	//钱包-算力-订单
	beego.Router("/hashrate/order/datagrid", &controllers.HashrateOrderController{}, "Get,Post:DataGrid")

	//钱包-算力-算力合约收益表
	beego.Router("/hashrate/hashrateOrderProfit/datagrid", &controllers.HashrateOrderProfitController{}, "Get,Post:DataGrid")

	//钱包-算力-算力合约资产发放表
	beego.Router("/hashrate/hashrateSendBalanceRecord/datagrid", &controllers.HashrateSendBalanceRecordController{}, "Get,Post:DataGrid")
	//钱包-算力-算力合约补发收益
	beego.Router("/hashrate/hashrateSendBalanceRecord/sendBalcnce", &controllers.HashrateSendBalanceRecordController{}, "Get,Post:SendBalcnce")

	//钱包-财务-交易明细
	beego.Router("/tokensky/tokenskyTransactionRecord/datagrid", &controllers.TokenskyTransactionRecordController{}, "Get,Post:DataGrid")
	//钱包-财务-充值记录表
	beego.Router("/tokensky/tokenskyUserDeposit/datagrid", &controllers.TokenskyUserDepositController{}, "Get,Post:DataGrid")

	/*理财*/

	//理财分类
	beego.Router("/financial/category/datagrid", &controllers.FinancialCategoryController{}, "Get,Post:DataGrid")
	beego.Router("/financial/category/edit", &controllers.FinancialCategoryController{}, "Post:Edit")
	beego.Router("/financial/category/delete", &controllers.FinancialCategoryController{}, "Post:Delete")
	//理财配置
	beego.Router("/financial/product/datagrid", &controllers.FinancialProductController{}, "Get,Post:DataGrid")
	beego.Router("/financial/product/edit", &controllers.FinancialProductController{}, "Post:Edit")
	beego.Router("/financial/product/delete", &controllers.FinancialProductController{}, "Post:Delete")
	beego.Router("/financial/product/theUpper", &controllers.FinancialProductController{}, "Post:TheUpper") //上下架
	//配置修改记录表
	beego.Router("/financial/productRecord/datagrid", &controllers.FinancialProductHistoricalRecordController{}, "Get,Post:DataGrid")
	//理财用户资产表
	beego.Router("/financial/liveUserBalance/datagrid", &controllers.FinancialLiveUserBalanceController{}, "Get,Post:DataGrid")
	//理财用户订单表
	beego.Router("/financial/order/datagrid", &controllers.FinancialOrderController{}, "Get,Post:DataGrid")
	//理财用户订单提币表
	beego.Router("/financial/orderWithdrawal/datagrid", &controllers.FinancialOrderWithdrawalController{}, "Get,Post:DataGrid")
	//理财用户收益表
	beego.Router("/financial/profit/datagrid", &controllers.FinancialProfitController{}, "Get,Post:DataGrid")

	/*借贷*/

	//配置列表
	beego.Router("/borrow/conf/datagrid", &controllers.BorrowConfController{}, "Get,Post:DataGrid")
	beego.Router("/borrow/conf/edit", &controllers.BorrowConfController{}, "Get,Post:Edit")
	beego.Router("/borrow/conf/theUpper", &controllers.BorrowConfController{}, "Get,Post:TheUpper")
	//订单表
	beego.Router("/borrow/order/datagrid", &controllers.BorrowOrderContoller{}, "Get,Post:DataGrid")
	beego.Router("/borrow/order/pledge", &controllers.BorrowOrderContoller{}, "Get,Post:GetMaxPledge")
	//日志表
	beego.Router("/borrow/order/addPledgeDataGrid", &controllers.BorrowOrdeLogController{}, "Get,Post:AddPledgeDataGrid")
	//强平表
	beego.Router("/borrow/limiting/datagrid", &controllers.BorrowLimitingController{}, "Get,Post:DataGrid")
	beego.Router("/borrow/limiting/sell", &controllers.BorrowLimitingController{}, "Get,Post:Sell")


	//用户邀请表
	beego.Router("/invite/formToAmount", &controllers.TokenskyUserInviteContoller{}, "Get,Post:FormToAmount")

	//七牛
	beego.Router("/qiniu/getkey", &controllers.QiNiuController{}, "Post:GetQiNiuKey")
	beego.Router("/qiniu/uploadFile", &controllers.QiNiuController{}, "Post:UploadFile")

	//私有接口 充值接口
	beego.Router("/personal/deposit/addBalance", &controllers.PersonalDepositController{}, "Post:AddBalance")
	beego.Router("/personal/withdraw/callback", &controllers.PersonalDepositController{}, "Post:Callback")

	//Options用于跨域复杂请求预检
	beego.NSRouter("/*", &controllers.BaseController{}, "Options:Options")

}
