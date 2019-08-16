package conf

import (
	"github.com/astaxie/beego"
	"strings"
)

//数据库相关配置
const (
	DB_MYSQL_TYPE        string = "mysql"       //数据库类型
	DB_MYSQL_ADMIN_CONF  string = "mysql-admin" //admin数据库
	DB_MYSQL_ADMIN_ALIAS string = "default"     //admin数据库别名
)

/*表前缀*/
const (
	DB_ADMIN_DT_PREFIX     string = "rms_"       //admin
	DB_TOKENSKY_DT_PREFIX  string = "tokensky_"  //Tokensky
	DB_OTC_DT_PREFIX       string = "otc_"       //oct
	DB_ROLE_DT_PREFIX      string = "role_"      //AdminRole
	DB_OPERATION_DT_PREFIX string = "operation_" //Operation
	DB_HASHRATE_DT_PREFIX  string = "hashrate_"  //Hashrate
	DB_FINANCIAL_DT_PREFIX string = "financial_" //Financial
	DB_BORROW_DT_PREFIX    string = "borrow_"    //Borrow
	DB_CREDIT_DT_PREFIX    string = "credit"     //Credit
)

//组合查询 搜索默认参数
const (
	QUREY_PARAM_MIN_LIMIT  int64 = 10
	QUREY_PARAM_MAX_LIMIT  int64 = 1000
	QUREY_PARAM_OFFSET     int64 = 0
	QUREY_PARAM_ORDER_DESC       = "desc" //倒序
	QUREY_PARAM_ORDER_ASC        = "asc"  //正序
)

//七牛配置 默认配置空
var (
	QINIU_BUCHENT_NAME string = ""
	QINIU_ACCESS_KEY   string = ""
	QINIU_SERVERT_KEY  string = ""
	QINIU_SERVER       string = ""
)

//初始化配置信息
func init() {
	qinConfName := "qiniu"
	//七牛配置
	QINIU_BUCHENT_NAME = beego.AppConfig.String(qinConfName + "::buchent_name")
	QINIU_ACCESS_KEY = beego.AppConfig.String(qinConfName + "::access_key")
	QINIU_SERVERT_KEY = beego.AppConfig.String(qinConfName + "::servert_key")
	QINIU_SERVER = beego.AppConfig.String(qinConfName + "::server")
}

var (
	JIANG_SERVER_URL string
)

//地址
const (
	//JIANG_SERVER_URL  string = "http://192.168.3.91:8080"
	TBI_SERVER_ADDRESS_MAX  int64  = 100   //数量
	TBI_SERVER_ADDRESS_TICK string = "300" //单位秒
)

func init() {
	JIANG_SERVER_URL = beego.AppConfig.String("jiangurl" + "::url")
}

//浮点数精度计算
const (
	FLOAT_PRECISE_NUM_6  int     = 6
	FLOAT_PRECISE_NUM_8  int     = 8
	FLOAT_PRECISE_NUM_10 int     = 10
	FLOAT_PRECISE_8      float64 = 0.00000001
)

//汇率
var (
	//爬取汇率地址
	EXCHANGE_RATE_USDT_MGODB_URL string
)

func init() {
	EXCHANGE_RATE_USDT_MGODB_URL = beego.AppConfig.String("mongo" + "::usdt_url")
}

//支持提币的货币类型
var (
	TOKENSKY_ADDRESS_COIN_TYPES        []string
	TOKENSKY_ACCEPT_BALANCE_COIN_TYPES map[string]struct{}
)

//资产更变接口
var TOKENSKY_BALANCE_CHANGE_URL = "127.0.0.1:9000"

const (
	CHANGE_ADD = "add" //加
	CHANGE_SUB = "sub" //减
	CHANGE_MUL = "mul" //乘
	CHANGE_QUO = "quo" //除
)

const (
	BALANCE_CHANGE_SOURCE_APP   int = 1 //后端
	BALANCE_CHANGE_SOURCE_ADMIN int = 2 //管理后台
	BALANCE_CHANGE_SOURCE_TICK  int = 3 //定时器服务
)

func init() {
	if str := beego.AppConfig.String("balance" + "::coin_type_address"); str != "" {
		TOKENSKY_ADDRESS_COIN_TYPES = strings.Split(str, ",")
	}
	if str := beego.AppConfig.String("balance" + "::accept_balance_coin_type"); str != "" {
		TOKENSKY_ACCEPT_BALANCE_COIN_TYPES = make(map[string]struct{})
		for _, st := range strings.Split(str, ",") {
			TOKENSKY_ACCEPT_BALANCE_COIN_TYPES[st] = struct{}{}
		}
	}
	if str := beego.AppConfig.String("balance" + "::balance_change_url"); str != "" {
		TOKENSKY_BALANCE_CHANGE_URL = str
	}
}

//算力单位
const (
	HASHRATE_UNIT_H int64 = 1
	HASHRATE_UNIT_K int64 = 1000
	HASHRATE_UNIT_M int64 = 1000000
	HASHRATE_UNIT_G int64 = 1000000000
	HASHRATE_UNIT_T int64 = 1000000000000
	HASHRATE_UNIT_P int64 = 1000000000000000
	HASHRATE_UNIT_E int64 = 1000000000000000000
)

//算力发放状态
var (
	HASHRATE_SEND_SIGN = false
)

//获取订单号
var (
	TokenskyOrderIdsIterationGetOid func(sn string) string
)

//算力支持
var (
	Hashrate_Send_Balance_Allow_Coin_Type []string
)

func init() {
	Hashrate_Send_Balance_Allow_Coin_Type = []string{"BTC"}
}

//邮件推送相关
const (
	//错误等级
	EMAIL_ERROR_LEVEL_ERROR     = "error"
	EMAIL_ERROR_LEVEL_CRITICAL  = "critical"
	EMAIL_ERROR_LEVEL_ALERT     = "alert"
	EMAIL_ERROR_LEVEL_EMERGENCY = "emergency"
)

var (
	//用户
	EMAIL_SEND_USER_NAME = ""
	EMAIL_SEND_USER_PWD  = ""
	//邮件用户
	EMAIL_ADDRESS_ERROR     []string
	EMAIL_ADDRESS_CRITICAL  []string
	EMAIL_ADDRESS_ALERT     []string
	EMAIL_ADDRESS_EMERGENCY []string
	//邮件提醒是否开启
	EMAIL_NOTIFY_SIGN = false
)

func init() {
	if str := beego.AppConfig.String("email" + "::user"); str != "" {
		EMAIL_SEND_USER_NAME = str
	}
	if str := beego.AppConfig.String("email" + "::pwd"); str != "" {
		EMAIL_SEND_USER_PWD = str
	}
	if strs := beego.AppConfig.String("email" + "::errorToAddress"); strs != "" {
		EMAIL_ADDRESS_ERROR = strings.Split(strs, ",")
	}
	if strs := beego.AppConfig.String("email" + "::criticalToAddress"); strs != "" {
		EMAIL_ADDRESS_CRITICAL = strings.Split(strs, ",")
	}
	if strs := beego.AppConfig.String("email" + "::alertToAddress"); strs != "" {
		EMAIL_ADDRESS_ALERT = strings.Split(strs, ",")
	}
	if strs := beego.AppConfig.String("email" + "::emergencyToAddress"); strs != "" {
		EMAIL_ADDRESS_EMERGENCY = strings.Split(strs, ",")
	}
	if con, err := beego.AppConfig.Int("email" + "::sign"); err == nil {
		if con == 1 {
			EMAIL_NOTIFY_SIGN = true
		}
	}
}
