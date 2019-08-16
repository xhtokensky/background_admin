package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

//查询的类
type RoleBlackListQueryParam struct {
	BaseQueryParam
	Phone string `json:"phone"` //手机号 模糊查询
	//发布时间
	StartTime int64 `json:"startTime"` //开始时间
	EndTime   int64 `json:"endTime"`   //截止时间
}

func (a *RoleBlackList) TableName() string {
	return RoleBlackListTBName()
}

//角色黑名单
type RoleBlackList struct {
	Id int `orm:"pk;column(id)"json:"id"form:"id"`

	//手机号
	Phone string `orm:"column(phone)"json:"phone"form:"phone"`
	//黑名单类型
	BalckType int `orm:"column(balck_type)"json:"balckType"form:"balckType"`
	//开始时间
	StartTime time.Time `orm:"auto_now_add;type(datetime);column(start_time)"json:"startTime"form:"startTime"`
	//结束时间
	EndTime time.Time `orm:"type(datetime);column(end_time)"json:"endTime"form:"endTime"`
	//持续时间
	PeriodTime int64 `orm:"column(period_time)"json:"periodTime"form:"periodTime"`

	//用户Uid 连表用户 一对多关系
	User *TokenskyUser `orm:"rel(fk)"json:"-"form:"-"`
	//用户uid
	UserId int `orm:"-"json:"userId"form:"-"`
	//昵称
	NickName string `orm:"-"json:"nickName"form:"-"`
}

// OtcEntrustOrderPageList 获取分页数据
func RoleBlackListPageList(params *RoleBlackListQueryParam) ([]*RoleBlackList, int64) {
	o := orm.NewOrm()
	query := o.QueryTable(RoleBlackListTBName())
	data := make([]*RoleBlackList, 0)
	//默认排序
	sortorder := "id"
	switch params.Sort {
	case "id":
		sortorder = "id"
	}
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}
	//手机号模糊查询 起始位置开始
	if params.Phone != "" {
		query = query.Filter("phone__istartswith", params.Phone)
	}
	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, (params.Offset-1)*params.Limit).RelatedSel().All(&data)
	//数据处理
	for _, obj := range data {
		//if obj.User != nil{
		obj.UserId = obj.User.UserId
		obj.NickName = obj.User.NickName
		//}
	}
	return data, total
}

// RoleBlackListOneById 根据id获取单条
func RoleBlackListOneById(id int) (*RoleBlackList, error) {
	o := orm.NewOrm()
	m := RoleBlackList{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// RoleBlackListDelete 批量删除
func RoleBlackListDelete(ids []int) (int64, error) {
	query :=orm.NewOrm().QueryTable(RoleBlackListTBName())
	num, err := query.Filter("id__in", ids).Delete()
	return num, err
}
