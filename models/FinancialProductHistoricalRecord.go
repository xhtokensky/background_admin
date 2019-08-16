package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

//查询的类
type FinancialProductHistoricalRecordQueryParam struct {
	BaseQueryParam
	StartTime int64  `json:"startTime"` //开始时间
	EndTime   int64  `json:"endTime"`   //截止时间
	ConfId    string `json:"confId"`    //id
	Category  string `json:"category"`  //活期定期
}

func (a *FinancialProductHistoricalRecord) TableName() string {
	return FinancialProductHistoricalRecordTBName()
}

//数据历史记录
type FinancialProductHistoricalRecord struct {
	Id      int               `orm:"pk;column(id)"json:"id"form:"id"`
	Config  *FinancialProduct `orm:"rel(fk);column(config)"json:"-"form:"-"`
	ConfId  int               `orm:"-"json:"confId"form:"confId"`
	Admin   *AdminBackendUser `orm:"rel(fk);column(admin_id)"json:"-"form:"-"`
	AdminId int               `orm:"-"json:"adminId"form:"adminId"`
	//新利率
	NewRate float64 `orm:"column(new_rate)"json:"newRate"form:"newRate"`
	//原利率
	OldRate float64 `orm:"column(old_rate)"json:"oldRate"form:"oldRate"`
	//类型 1活期 2定期
	Category int `orm:"column(category)"json:"category"form:"category"`
	//描述
	Msg string `orm:"column(msg)"json:"msg"form:"msg"`
	//创建时间
	CreateTime time.Time `orm:"auto_now_add;type(datetime);column(create_time)"json:"createTime"form:"createTime"`
	//类型 集合[新增,编辑]
	RecordType string `orm:"column(record_type)"json:"recordType"form:"recordType"`

	Name string `orm:"-"json:"name"form:"name"`
}

func FinancialConfigHistoricalPageList(params *FinancialProductHistoricalRecordQueryParam) ([]*FinancialProductHistoricalRecord, int64) {
	o := orm.NewOrm()
	query := o.QueryTable(FinancialProductHistoricalRecordTBName())
	data := make([]*FinancialProductHistoricalRecord, 0)
	//默认排序
	sortorder := "id"
	switch params.Sort {
	case "id":
		sortorder = "id"
	}
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}
	if params.ConfId != "" {
		query = query.Filter("Config__id__exact", params.ConfId)
	}
	if params.Category != "" {
		query = query.Filter("category__exact", params.Category)
	}
	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, (params.Offset-1)*params.Limit).All(&data)
	for _,obj := range data{
		if obj.Config != nil{
			obj.ConfId = obj.Config.Id
		}
		if obj.Admin != nil{
			obj.Name = obj.Admin.UserName
			obj.AdminId = obj.Admin.Id
		}
	}
	return data, total
}
