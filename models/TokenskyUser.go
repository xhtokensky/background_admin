package models

import (
	"github.com/astaxie/beego/orm"
	"time"
	"tokensky_bg_admin/conf"
)

func (a *TokenskyUser) TableName() string {
	return TokenskyUserTBName()
}

//查询类
type TokenskyUserQueryParam struct {
	BaseQueryParam
	StartTime int64  `json:"startTime"` //开始时间
	EndTime   int64  `json:"endTime"`   //截止时间
	Phone     string `json:"phone"`     //电话
	UserId    int    `json:"userId"`    //Uid
	NickName  string `json:"nickName"`  //昵称
}

//用户
type TokenskyUser struct {
	UserId int `orm:"pk;column(user_id)"json:"userId"form:"userId"`
	SmsId  int `orm:"column(sms_id)"json:"smsId"form:"smsId"`
	//昵称
	NickName string `orm:"column(nick_name)"json:"nickName"form:"nickName"`
	//电话
	Phone string `orm:"column(phone)"json:"phone"form:"phone"`
	//用户真实姓名
	UserName string `orm:"column(user_name)"json:"userName"form:"userName"`
	//登陆账号
	Account string `orm:"unique;column(account)"json:"account"form:"account"`
	//密码
	Password            string `orm:"column(password)"json:"-"form:"password"`
	TransactionPassword string `orm:"column(transaction_password)"json:"transactionPassword"form:"transactionPassword"`
	//积分
	Points int `orm:"column(points)"json:"points"form:"points"`
	//账户状态 账号是否有效1为有效0为无效
	UserStatus int `orm:"column(user_status)"json:"userStatus"form:"userStatus"`
	//是否锁住
	IsLock int `orm:"column(is_lock)"json:"isLock"form:"isLock"`
	//是否登陆 0未登录，1登录
	IsLogin int `orm:"column(is_login)"json:"isLogin"form:"isLogin"`
	//密码错误次数
	PwdErrorNumber int `orm:"column(pwd_error_number)"json:"pwdErrorNumber"form:"pwdErrorNumber"`
	//电子邮件
	Email string `orm:"column(email)"json:"email"form:"email"`
	//性别 1男，2女，3保密
	Sex int `orm:"column(sex)"json:"sex"form:"sex"`
	//头像
	HeadImg string `orm:"column(head_img)"json:"headImg"form:"headImg"`
	//注册设备来源
	RegistDeviceType string `orm:"column(regist_device_type)"json:"registDeviceType"form:"registDeviceType"`
	//用户登记
	UserLevel int `orm:"column(user_level)"json:"userLevel"form:"userLevel"`
	//创建人 ID 1前台，2后台
	CreatorId int `orm:"column(creator_id)"json:"creatorId"form:"creatorId"`
	//创建时间
	CreateTime time.Time `orm:"auto_now_add;type(datetime);column(create_time)"json:"createTime"form:"createTime"`
	//IP地址
	CreateIp string `orm:"column(create_ip)"json:"createIp"form:"createIp"`
	//更新人ID
	UpdaterId int `orm:"column(updater_id)"json:"updaterId"form:"updaterId"`
	//
	Salt string `orm:"column(salt)"json:"salt"form:"salt"`
	//邀请码
	InviteCode string `orm:"column(invite_code)"json:"inviteCode"form:"inviteCode"`
	//最后登陆时间
	LastLoginTime time.Time `orm:"type(datetime);column(last_login_time)"json:"lastLoginTime"form:"lastLoginTime"`
	//等级
	Level int `orm:"column(level)"json:"level"form:"level"`
	//是否有权限推广 0没有权限 1拥有权限
	Invitation int `orm:"invitation"json:"invitation"form:"invitation"`
	//身份认证状态 0尚未认证 1认证成功 2认证失败
	RealAuth int `orm:"-"json:"realAuth"form:"realAuth"`
	//支付状态
	AccountBank []int `orm:"-"json:"accountBank"form:"accountBank"`
	//连表 黑名单 一对多反向关系
	Blacks []*RoleBlackList `orm:"reverse(many)"json:"blacks"form:"-"`
	//连表 用户审核 一对多反向关系
	TokenskyRealAuths []*TokenskyRealAuth `orm:"reverse(many)"json:"-"form:"-"`
	//连表 用户委托单 一对多反向关系
	OtcEntrustOrders []*OtcEntrustOrder `orm:"reverse(many)"json:"-"form:"-"`
	//连表 用户订单表 一对多反向关系
	OtcOrderVendor []*OtcOrder `orm:"reverse(many)"json:"-"form:"-"`
	OtcOrderVendee []*OtcOrder `orm:"reverse(many)"json:"-"form:"-"`
	//连表 算力订单表 一对多反向关系
	HashrateOrders []*HashrateOrder `orm:"reverse(many)"json:"-"form:"-"`
	//连表 提币审核表 一对多反向关系
	TokenskyUserTibis []*TokenskyUserTibi `orm:"reverse(many)"json:"-"form:"-"`
	//连表 交易明细表 一对多反向关系
	TokenskyTransactionRecords []*TokenskyTransactionRecord `orm:"reverse(many)"json:"-"form:"-"`
	//连表 用户充值表
	TokenskyUserDeposits []*TokenskyUserDeposit `orm:"reverse(many)"json:"-"form:"-"`
	//连表 用户收益表
	HashrateOrderProfits []*HashrateOrderProfit `orm:"reverse(many)"json:"-"form:"-"`
	//连表 用户消息
	TokenskyMessages []*TokenskyMessage `orm:"reverse(many)"json:"-"form:"-"`
	//理财相关
	FinancialProfits []*FinancialProfit `orm:"reverse(many)"json:"-"form:"-"`
	FinancialOrders []*FinancialOrder `orm:"reverse(many)"json:"-"form:"-"`
	FinancialLiveUserBalances []*FinancialLiveUserBalance `orm:"reverse(many)"json:"-"form:"-"`
	//借贷相关
	BorrowOrders []*BorrowOrder `orm:"reverse(many)"json:"-"form:"-"`
	BorrowOrdeLogs []*BorrowOrdeLog `orm:"reverse(many)"json:"-"form:"-"`
	BorrowLimitings []*BorrowLimiting `orm:"reverse(many)"json:"-"form:"-"`
	//其它
	TokenskyUserInvites []*TokenskyUserInvite  `orm:"reverse(many)"json:"-"form:"-"`
	//用户资产变化记录表
	TokenskyUserBalancesRecords []*TokenskyUserBalancesRecord`orm:"reverse(many)"json:"-"form:"-"`
	//电子资产
	Electricity float64 `orm:"-"json:"electricity"form:"electricity"`
}

//分页数据
func TokenskyUserPageList(params *TokenskyUserQueryParam) ([]*TokenskyUser, int64) {
	o := orm.NewOrm()
	query := o.QueryTable(TokenskyUserTBName())
	data := make([]*TokenskyUser, 0)
	//默认排序
	sortorder := "user_id"
	switch params.Sort {
	case "user_id":
		sortorder = "user_id"
	}
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}
	//电话筛选
	if params.Phone != "" {
		query = query.Filter("phone__icontains", params.Phone)
	}
	//昵称筛选
	if params.NickName != "" {
		query = query.Filter("nick_name__icontains", params.NickName)
	}
	//用户Id筛选
	if params.UserId > 0 {
		query = query.Filter("user_id__icontains", params.UserId)
	}
	//时间筛选 pass
	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, (params.Offset-1)*params.Limit).RelatedSel().All(&data)
	//交易密码是否存在
	for _, v := range data {
		if v.TransactionPassword != "" {
			v.TransactionPassword = "ok"
		} else {
			v.TransactionPassword = "no"
		}
	}
	ids := make([]int, 0)
	for _, v := range data {
		ids = append(ids, v.UserId)
	}
	//身份认证状态
	if mapp, err := TokenskyRealAuthByIds(ids); err == nil {
		for _, v := range data {
			if obj, found := mapp[v.UserId]; found {
				v.RealAuth = obj.Status
			}
		}
	}

	//用户支付方式
	if mapp := TokenskyAccountBanksByIds(ids); len(mapp) > 0 {
		for _, v := range data {
			if objs, found := mapp[v.UserId]; found {
				v.AccountBank = make([]int, conf.PAY_TYPE_MAX_NUM)
				for _, obj := range objs {
					if obj.Type > 0 {
						v.AccountBank[obj.Type-1] = 1
					}
				}
			} else {
				v.AccountBank = make([]int, conf.PAY_TYPE_MAX_NUM)
			}
		}
	} else {
		for _, v := range data {
			v.AccountBank = make([]int, conf.PAY_TYPE_MAX_NUM)
		}
	}
	//用户电力资产
	if mapp := TokenskyUserElectricityBalanceByUids(ids); len(mapp) > 0 {
		for _, v := range data {
			if ebj, found := mapp[v.UserId]; found {
				v.Electricity = ebj.Balance
			}

		}
	}
	return data, total
}

//根据电话查询用户uid
func TokenskyUserGetIdsByPhone(phone string) map[int]string {
	mapp := make(map[int]string)
	data := make([]*TokenskyUser, 0)
	query:= orm.NewOrm().QueryTable(TokenskyUserTBName())
	query.Filter("phone__icontains", phone).All(&data, "user_id", "phone")
	if len(data) > 0 {
		for _, v := range data {
			mapp[v.UserId] = v.Phone
		}
	}
	return mapp
}

//获取单条
func TokenskyUserOne(id int) (*TokenskyUser, error) {
	o := orm.NewOrm()
	m := TokenskyUser{UserId: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

//根据电话查询用户
func TokenskyUserOneByPhone(phone string) *TokenskyUser {
	o := orm.NewOrm()
	m := TokenskyUser{Phone: phone}
	err := o.Read(&m, "phone")
	if err != nil {
		return nil
	}
	return &m
}



