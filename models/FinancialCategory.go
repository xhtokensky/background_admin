package models

import (
	"github.com/astaxie/beego/orm"
	"time"
	"tokensky_bg_admin/utils"
)

//理财分类表

/*
CREATE TABLE `financial_category` (
  `id` int(11) NOT NULL,
  `avatar` varchar(255) NOT NULL,
  `symbol` varchar(255) NOT NULL,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='理财分类'
*/

//查询的类
type FinancialCategoryQueryParam struct {
	BaseQueryParam
	Name      string `json:"name"`
	StartTime int64  `json:"startTime"` //开始时间
	EndTime   int64  `json:"endTime"`   //截止时间
	Status    string `json:"status"`    //状态
	Id        int    `json:"id"`
}

// TableName 设置表名
func (a *FinancialCategory) TableName() string {
	return FinancialCategoryTBName()
}

type FinancialCategory struct {
	Id int `orm:"pk;column(id)"json:"id"form:"id"`
	//头像
	Avatar string `orm:"column(avatar)"json:"avatar"form:"avatar"`
	//路径
	AvatarUrl string `orm:"-"json:"avatarUrl"form:"avatarUrl"`
	//货币类型
	Symbol string `orm:"column(symbol)"json:"symbol"form:"symbol"`
	//创建时间
	CreateTime time.Time `orm:"auto_now_add;type(datetime);column(create_time)"json:"createTime"form:"createTime"`
	//更新时间
	//UpdateTime time.Time `orm:"auto_now;type(datetime);column(update_time)"json:"updateTime"form:"updateTime"`
	AdminId int `orm:"column(admin_id)"json:"adminId"form:"adminId"`
	//
	FinancialBaseConfigs []*FinancialProduct `orm:"reverse(many)"json:"-"form:"-"`
}

//获取分页数据
func FinancialCategoryPageList(params *FinancialCategoryQueryParam) ([]*FinancialCategory, int64) {
	o := orm.NewOrm()
	query := o.QueryTable(FinancialCategoryTBName())
	data := make([]*FinancialCategory, 0)
	//默认排序
	sortorder := "id"
	switch params.Sort {
	case "id":
		sortorder = "id"
	}
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}
	if params.Id != 0 {
		query = query.Filter("id__exact", params.Id)
	}
	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, (params.Offset-1)*params.Limit).All(&data)
	//图片下载凭证
	deadline := time.Now().Add(time.Second * 3600).Unix() //1小时有效期
	for _, obj := range data {
		obj.AvatarUrl = utils.QiNiuDownload(obj.Avatar, deadline)
	}
	return data, total
}

func FinancialCategoryDelete(id int) (bool, string) {
	o := orm.NewOrm()
	query := o.QueryTable(FinancialProductTBName())
	query.Filter("financial_category_id__exact", id)
	count, _ := query.Count()
	if count > 0 {
		return false, "关联中不可删除"
	}
	_,err := o.QueryTable(FinancialCategoryTBName()).Filter("id__exact",id).Delete()
	if err != nil {
		return false, err.Error()
	}
	return false, ""
}
