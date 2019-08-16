package models

import (
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"io/ioutil"
	"net/http"
	"time"
	"tokensky_bg_admin/conf"
)

//TokenskyUserBalance tokensky_user_address 用户充提地址表

/*
DROP TABLE IF EXISTS `tokensky_user_address`;
CREATE TABLE `tokensky_user_address` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `address` varchar(255) NOT NULL,
  `coin_type` varchar(255) NOT NULL,
  `user_id` int(11) NOT NULL DEFAULT '0',
  `status` int(11) NOT NULL DEFAULT '1',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `user_id_index` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户充提币地址表';

SET FOREIGN_KEY_CHECKS = 1;

*/

type TokenskyUserAddressQueryParam struct {
	BaseQueryParam
	StartTime string `json:"startTime"` //开始时间
	EndTime   string `json:"endTime"`   //截止时间

}

func (a *TokenskyUserAddress) TableName() string {
	return TokenskyUserAddressTBName()
}

//用户充值地址簿
type TokenskyUserAddress struct {
	Id         int       `orm:"pk;column(id)"json:"id"form:"id"`
	Address    string    `orm:"column(address)"json:"address"form:"address"`
	CoinType   string    `orm:"column(coin_type)"json:"coinType"form:"coinType"`
	UserId     int       `orm:"column(user_id)"json:"userId"form:"userId"`
	Status     int       `orm:"column(status)"json:"status"form:"status"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime);column(create_time)"json:"createTime"form:"createTime"`
}

//获取空的地址簿数量
func TokenskyUserAddressGetNotUsedCont(coinType string) int64 {
	var total int64
	o := orm.NewOrm()
	query := o.QueryTable(TokenskyUserAddressTBName())
	query = query.Filter("user_id__exact", 0)
	query = query.Filter("coin_type__iexact", coinType)
	total, _ = query.Count()
	return total
}

type tokenskyUserAddressAddNumBody struct {
	Code    int    `json:"code"`
	Rresult string `json:"result"`
}

//新增地址薄
func TokenskyUserAddressAddNum(coinType string, num int64) {
	data := make([]*TokenskyUserAddress, 0, num)
	client := &http.Client{}

	var obj tokenskyUserAddressAddNumBody
	url := conf.JIANG_SERVER_URL + "/service/" + coinType + "/address"
	for i := int64(0); i < num; i++ {
		resp, err := client.Get(url)
		if err != nil {
			//异常
			continue
		}
		if resp.StatusCode >= 300 && resp.StatusCode < 200 {
			//异常
			continue
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			//异常
			continue
		}
		if err := json.Unmarshal(body, &obj); err != nil {
			//异常
			continue
		}
		if obj.Code != 0{
			//异常
			continue
		}
		data = append(data, &TokenskyUserAddress{
			Address:  obj.Rresult,
			CoinType: coinType,
			UserId:   0,
			Status:   1, //状态默认1
		})
	}
	//新增
	if con := len(data); con > 0 {
		o := orm.NewOrm()
		_, err := o.InsertMulti(con, &data)
		if err != nil {
			//异常
		}
	}
}

//根据地址获取
func TokenskyUserAddressByCoinTypeAndAddress(coinType string, address string) *TokenskyUserAddress {
	obj := TokenskyUserAddress{}
	o := orm.NewOrm()
	query := o.QueryTable(TokenskyUserAddressTBName())
	query = query.Filter("coin_type__iexact", coinType).Filter("address__exact", address)
	if err := query.One(&obj); err != nil {
		return nil
	}
	return &obj

}
