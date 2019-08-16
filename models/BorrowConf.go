package models

import (
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"time"
	"tokensky_bg_admin/utils"
)

//查询的类
type BorrowConfQueryParam struct {
	BaseQueryParam
	StartTime  int64  `json:"startTime"`  //开始时间
	EndTime    int64  `json:"endTime"`    //截止时间
	Status     string `json:"status"`     //状态
	Title       string `json:"title"` //标题
	CoinType   string `json:"coinType"`   //质押货币类型
	LoanSymbol string `json:"loanSymbol"` //借贷货币类型
}

// TableName 设置表名
func (a *BorrowConf) TableName() string {
	return BorrowConfTBName()
}

type BorrowConf struct {
	Id int `orm:"pk;column(id)"json:"id"form:"id"`
	//图标名
	Icon    string `orm:"column(icon)"json:"icon"form:"icon"`
	IconUrl string `orm:"-"json:"iconUrl"form:"iconUrl"`
	//标题内容
	Title string `orm:"column(title)"json:"title"form:"title"`
	//质押率最大值
	PledgeRateMax float64 `orm:"column(pledge_rate_max)"json:"pledgeRateMax"form:"pledgeRateMax"`
	//质押周期天数
	CycleDay int `orm:"column(cycle_day)"json:"cycleDay"form:"cycleDay"`
	//质押日利率
	DayRate float64 `orm:"column(day_rate)"json:"dayRate"form:"dayRate"`
	//逾期的日利率
	OverdueRate float64 `orm:"column(overdue_rate)"json:"overdueRate"form:"overdueRate"`
	//质押货币类型
	CoinType string `orm:"column(coin_type)"json:"coinType"form:"coinType"`
	//借贷货币类型
	LoanSymbol string `orm:"column(loan_symbol)"json:"loanSymbol"form:"loanSymbol"`
	//优先级
	Priority int `orm:"column(priority)"json:"priority"form:"priority"`
	//创建时间 单位毫秒
	CreateTime int64 `orm:"column(create_time)"json:"-"form:"-"`
	CreateTimeTwo time.Time `orm:"-"json:"createTime"form:"-"`
	//是否上架 0 待上架 1上架 2下架
	IsPutaway int `orm:"column(is_putaway)"json:"isPutaway"form:"isPutaway"`
	//介绍内容
	Introduce string `orm:"column(introduce)"json:"introduce"form:"introduce"`
	//创建人
	AdminId int `orm:"column(admin_id)"json:"adminId"form:"adminId"`
	//连表
	BorrowOrders []*BorrowOrder `orm:"reverse(many)"json:"-"form:"-"`
}

func (this *BorrowConf)Json() string{
	js,err := json.Marshal(this)
	if err != nil{
		return "err"
	}
	return string(js)
}

//获取分页数据
func BorrowConfPageList(params *BorrowConfQueryParam) ([]*BorrowConf, int64) {
	o := orm.NewOrm()
	query := o.QueryTable(BorrowConfTBName())
	data := make([]*BorrowConf, 0)
	//默认排序
	sortorder := "id"
	switch params.Sort {
	case "id":
		sortorder = "id"
	}
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}
	if params.CoinType != "" {
		query = query.Filter("coin_type__exact", params.CoinType)
	}
	if params.LoanSymbol != "" {
		query = query.Filter("loan_symbol__exact", params.LoanSymbol)
	}
	if params.Title != ""{
		query = query.Filter("title__icontains",params.Title)
	}
	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, (params.Offset-1)*params.Limit).All(&data)

	//图片下载凭证
	deadline := time.Now().Add(time.Second * 3600).Unix() //1小时有效期
	for _, obj := range data {
		obj.IconUrl = utils.QiNiuDownload(obj.Icon, deadline)
		obj.CreateTimeTwo = time.Unix(obj.CreateTime/1000,0)
	}
	return data, total
}

//
func BorrowConfOne(id int)*BorrowConf{
	o := orm.NewOrm()
	obj := &BorrowConf{Id:id}
	err := o.Read(obj)
	if err != nil{
		return nil
	}
	return obj
}

//上下架
func BorrowConfTheUppers(o orm.Ormer,status int, ids []int,aid int) (num int64, err error) {
	query := o.QueryTable(BorrowConfTBName())
	params := map[string]interface{}{
		"is_putaway": status,
		"admin_id":aid,
	}
	switch status {
	case 1:
		//上架
		num, err = query.Filter("id__in", ids).Filter("is_putaway__exact", 0).Update(params)
	case 2:
		//后台下架
		num, err = query.Filter("id__in", ids).Filter("is_putaway__exact", 1).Update(params)
	default:
		return 0, err
	}
	return num, err
}
