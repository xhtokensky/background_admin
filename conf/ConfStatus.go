package conf

//运营-Banner
const (
	OPERATION_BANNER_STATUS_OFF int = 0 //关闭
	OPERATION_BANNER_STATUS_NO  int = 1 //开启
)

//OTC-申诉状态
const (
	OTC_APPEAL_STATUS_UOD      int = 0 //未处理
	OTC_APPEAL_STATUS_VALIDATE int = 1 //确认放币
	OTC_APPEAL_STATUS_CANCEL   int = 2 //取消订单
)

/*
等待对方支付 1:已完成 已完成 已完成 2:已支付 等待对方放币 对方已支付 3.已申诉 卖方已申诉 已申诉
4:客服审核申诉 客服审核中 客服审核中 5:已取消 已取消 对方已取消
6:超时取消 超时未支付自动取消  对方支付超时自动取消
*/
//OTC-订单状态
const (
	OTC_ORDER_STATUS_WAIT   int = 0 //待支付
	OTC_ORDER_STATUS_FOUND  int = 1 //已完成
	OTC_ORDER_STATUS_PAID   int = 2 //已支付 等待对方放币 对方已支付
	OTC_ORDER_STATUS_APPEAL int = 3 //已申诉 卖方已申诉 已申诉
	OTC_ORDER_STATUS_CANCEL int = 4 //已取消

)

//OTC-订单类型
const (
	OTC_ORDER_TYPE_VENDEE int = 1 //订单类型 买入
	OTC_ORDER_TYPE_VENDOR int = 2 //订单类型 卖出
)

//tokensky-角色审核
const (
	TOKENSKY_REAL_AUTH_STATUS_INITIAL int = 0 //未认证
	TOKENSKY_REAL_AUTH_STATUS_PASSED  int = 1 //认证通过
	TOKENSKY_REAL_AUTH_STATUS_FAILED  int = 2 //认证未通过
)

//tokensky-黑名单
const (
	ROLE_BLACK_LIST_STATUS_BAN_LANDING int = 1 //禁止登陆
	ROLE_BLACK_LIST_STATUS_BAN_TRADING int = 2 //禁止交易
)

//支付最大数量
const (
	PAY_TYPE_MAX_NUM int = 3
)

//订单生成规则 业务码
const (
	ORDER_BUSINESS_OTC_VENDEE_CODE           = "00" //OTC买入
	ORDER_BUSINESS_OTC_VENDOR_CODE           = "01" //OTC卖出
	ORDER_BUSINESS_HASHRATE_CATEGORY_CODE    = "02" //购买算力合约
	ORDER_BUSINESS_OTC_CHONEBI_CODE          = "03" //充币
	ORDER_BUSINESS_OTC_TIBI_CODE             = "04" //提币
	ORDER_BUSINESS_HASHRATE_SEND_PRICEP_CODE = "05" //发放合约收益
	ORDER_BUSINESS_ADD_ELECTRIC              = "06" //充电费
)

//提币状态码
const (
	TIBI_ERR_SUCCESS          = 0 //成功
	TIBI_ERR_BAD_PARAMETER    = 1 //参数错误
	TIBI_ERR_SERVER_ERROR     = 2 //服务器错误
	TIBI_ERR_WITHDRAW_EXISTED = 3 //已存在
	TIBI_ERR_INVALID_ADDRESS  = 4 //无效地址
	TIBI_ERR_INVALID_AMOUNT   = 5 //无效金额
)
