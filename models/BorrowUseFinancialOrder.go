package models

import "github.com/astaxie/beego/orm"

/*
DROP TABLE IF EXISTS `borrow_use_financial_order`;
CREATE TABLE `borrow_use_financial_order` (
  `id` bigint(16) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL COMMENT '用户id',
  `loan_order_id` varchar(50) DEFAULT NULL COMMENT '借贷记录id',
  `create_way` tinyint(1) DEFAULT '1' COMMENT '创建的方式 1 借贷  2 借贷增加质押数',
  `financial_id` int(11) DEFAULT NULL COMMENT '理财包id',
  PRIMARY KEY (`id`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

SET FOREIGN_KEY_CHECKS = 1;
*/

func (a *BorrowUseFinancialOrder) TableName() string {
	return BorrowUseFinancialOrderTBName()
}

type BorrowUseFinancialOrder struct {
	Id            int64 `orm:"pk;column(id)"json:"id"form:"id"`
	UserId        int `orm:"column(user_id)"json:"userId"form:"userId"`
	BorrowOrderId string `orm:"column(loan_order_id)"json:"-"form:"-"`
	//借贷方式 1 借贷 2 借贷增加质押数
	CreateWay int `orm:"column(create_way)"json:"createWay"form:"createWay"`
	//理财包id
	FinancialOrderId int `orm:"column(financial_id)"json:"-"form:"-"`
}

//获取质押都理财保ids
func BorrowUseFinancialOrderGetFinancials()(map[int]struct{},error){
	mapp := make(map[int]struct{})
	o := orm.NewOrm()
	query := o.QueryTable(BorrowUseFinancialOrderTBName())
	count,_:= query.Count()
	num := 1
	for count>0{
		data := make([]*BorrowUseFinancialOrder,0)
		query.Limit(500,(num-1)*500).All(&data)
		count -= int64(len(data))
		num++
		for _,v := range data{
			mapp[v.FinancialOrderId] = struct{}{}
		}
		if len(data)==0{
			break
		}
	}
	return mapp,nil
}

//判断定期是否被锁定
func BorrowUseFinancialOrderIsLock(o orm.Ormer,financialOrderId int)(bool,error){
	query := o.QueryTable(BorrowUseFinancialOrderTBName())
	count,err := query.Filter("financial_id__exact",financialOrderId).Count()
	if err != nil{
		return false,err
	}
	if count>0{
		return true,nil
	}
	return false,nil
}

//删除指定订单
func BorrowUseFinancialOrderDelete(o orm.Ormer,orderId string)error{
	query := o.QueryTable(BorrowUseFinancialOrderTBName())
	_,err := query.Filter("loan_order_id__exact",orderId).Delete()
	return err
}