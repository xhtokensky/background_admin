package models

import "github.com/astaxie/beego/orm"

//货币配置表

/*
CREATE TABLE `tokensky_user_balance_coin` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` varchar(255) DEFAULT NULL COMMENT '名称',
  `symbol` varchar(255) DEFAULT NULL COMMENT '标识',
  `sort` int(11) DEFAULT NULL COMMENT '排序',
  `status` int(1) DEFAULT '1' COMMENT '启动状态1为启动0为关闭',
  `avatar` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4;
*/

func (a *TokenskyUserBalanceCoin) TableName() string {
	return TokenskyUserBalanceCoinTBName()
}

type TokenskyUserBalanceCoin struct {
	Id int `orm:"pk;column(id)"json:"id"form:"id"`
	Name string `orm:"column(name)"json:"name"form:"name"`
	Symbol string `orm:"column(symbol)"json:"symbol"form:"symbol"`
	Sort int `orm:"column(sort)"json:"sort"form:"sort"`
	Status int `orm:"column(status)"json:"status"form:"status"`
	Avatar string `orm:"column(avatar)"json:"avatar"form:"avatar"`
}

//货币是否存在
func TokenskyUserBalanceCoinIsFound(name string)bool{
	o := orm.NewOrm()
	query := o.QueryTable(TokenskyUserBalanceCoinTBName())
	query = query.Filter("symbol__exact",name).Filter("status__exact",1)
	count,_ := query.Count()
	if count>0{
		return true
	}
	return false
}