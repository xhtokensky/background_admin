package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

/*
DROP TABLE IF EXISTS `tokensky_account_bank`;
CREATE TABLE `tokensky_account_bank` (
  `key_id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `bank_user_name` varchar(50) DEFAULT NULL,
  `bank_card_no` varchar(50) DEFAULT NULL,
  `bank_name` varchar(50) DEFAULT NULL,
  `bank_branch_name` varchar(200) DEFAULT NULL,
  `alipay_user_name` varchar(100) DEFAULT NULL,

  `alipay_account` varchar(100) DEFAULT NULL,
  `alipay_qr_code` varchar(255) DEFAULT NULL,

  `wechat_user_name` varchar(100) DEFAULT NULL,
  `wechat_account` varchar(100) DEFAULT NULL,

  `wechat_qr_code` varchar(100) DEFAULT NULL,
  `status` int(11) NOT NULL DEFAULT '1',
  `type` int(11) DEFAULT '1' COMMENT '1银行卡 2支付宝 3微信',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`key_id`),
  KEY `user_id_index` (`user_id`),
  KEY `key_id_index` (`key_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='收款管理地址表';
*/

//查询的类
type TokenskyAccountBankQueryParam struct {
	BaseQueryParam
}

func (a *TokenskyAccountBank) TableName() string {
	return TokenskyAccountBankTBName()
}

//用户支付方式
type TokenskyAccountBank struct {
	KeyId          int `orm:"pk;column(key_id)"json:"keyId"form:"keyId"`
	UserId         int `orm:"column(user_id)"json:"userId"form:"userId"`
	BankUserName   string
	BankCardNo     string
	BankName       string
	BankBranchName string
	AlipayUserName string
	AlipayAccount  string
	AlipayQrCode   string
	WechatUserName string
	WechatAccount  string
	WechatQrCode   string
	Status         int
	//1银行卡 2支付宝 3微信
	Type       int `orm:"column(type)"json:"type"form:"type"`
	CreateTime time.Time
	UpdateTime time.Time
}

//获取用户支付方式
func TokenskyAccountBanksByIds(ids []int) map[int][]*TokenskyAccountBank {
	mapp := make(map[int][]*TokenskyAccountBank)
	o := orm.NewOrm()
	query := o.QueryTable(TokenskyAccountBankTBName())
	data := make([]*TokenskyAccountBank, 0)
	if len(ids) > 0 {
		query.Filter("user_id__in", ids).All(&data, "key_id", "user_id", "type")
	}
	for _, v := range data {
		if _, found := mapp[v.UserId]; !found {
			mapp[v.UserId] = make([]*TokenskyAccountBank, 0)
		}
		mapp[v.UserId] = append(mapp[v.UserId], v)
	}
	return mapp
}
