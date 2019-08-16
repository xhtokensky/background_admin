package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

//理财定期配置表

/*
CREATE TABLE `financial_base_config` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `financial_category_id` int(11) NOT NULL,
  `category` int(11) NOT NULL DEFAULT '1' COMMENT '1活期 2定期',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '0待上架 1上架 2下架 ',
  `cycle` int(11) NOT NULL COMMENT '周期  以天为单位',
  `year_profit` double NOT NULL COMMENT '年化收益率',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `admin_id` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='理财定期配置表'
*/

//查询的类
type FinancialProductQueryParam struct {
	BaseQueryParam
	StartTime           int64  `json:"startTime"`           //开始时间
	EndTime             int64  `json:"endTime"`             //截止时间
	Status              string `json:"status"`              //状态
	FinancialCategoryId string `json:"financialCategoryId"` //关联id
	Category            string `json:"category"`            //活期定期
}

func (a *FinancialProduct) TableName() string {
	return FinancialProductTBName()
}

type FinancialProduct struct {
	Id int `orm:"pk;column(id)"json:"id"form:"id"`
	//关联id
	FinancialCategoryObj *FinancialCategory `orm:"rel(fk);column(financial_category_id)"json:"-"form:"-"`
	FinancialCategoryId  int                `orm:"-"json:"financialCategoryId"form:"financialCategoryId"`
	//1活期 2定期
	Category int `orm:"column(category)"json:"category"form:"category"`
	//状态 0待上架 1上架 2下架
	Status int `orm:"column(status)"json:"status"form:"status"`
	//周期 天为单位
	Cycle int `orm:"column(cycle)"json:"cycle"form:"cycle"`
	//最小起投额
	MinQuantity float64 `orm:"column(min_quantity)"json:"minQuantity"form:"minQuantity"`

	//年化收益
	YearProfit float64 `orm:"column(year_profit)"json:"yearProfit"form:"yearProfit"`
	//创建时间
	CreateTime time.Time `orm:"auto_now_add;type(datetime);column(create_time)"json:"createTime"form:"createTime"`
	//更新时间
	UpdateTime time.Time `orm:"auto_now;type(datetime);column(update_time)"json:"updateTime"form:"updateTime"`
	//排序
	Sort int `orm:"column(sort)"json:"sort"form:"sort"`
	//标题
	Titile string `orm:"column(title)"json:"title"form:"title"`
	//admin
	Admin   *AdminBackendUser `orm:"rel(fk);column(admin_id)"json:"-"form:"-"`
	AdminId int               `orm:"-"json:"adminId"form:"adminId"`
	//连表
	FinancialConfigHistoricalRecords []*FinancialProductHistoricalRecord `orm:"reverse(many)"json:"-"form:"-"`
	//描述[只接受数据]
	Msg string `orm:"-"json:"-"form:"msg"`

	//创建人
	Name string `orm:"-"json:"name"form:"name"`
	//支持币种
	Symbol string  `orm:"-"json:"symbol"form:"symbol"`

}

//获取分页数据
func FinancialProductPageList(params *FinancialProductQueryParam) ([]*FinancialProduct, int64) {
	o := orm.NewOrm()
	query := o.QueryTable(FinancialProductTBName())
	data := make([]*FinancialProduct, 0)
	//默认排序
	sortorder := "id"
	switch params.Sort {
	case "id":
		sortorder = "id"
	}
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}
	if params.FinancialCategoryId != "" {
		query = query.Filter("financial_category_id__exact", params.FinancialCategoryId)
	}
	if params.Category != "" {
		query = query.Filter("category__exact", params.Category)
	}
	if params.Status != "" {
		query = query.Filter("status__exact", params.Status)
	}
	//时间段
	if params.StartTime > 0 {
		query = query.Filter("create_time__gte", time.Unix(params.StartTime, 0))
	}
	if params.EndTime > 0 {
		query = query.Filter("create_time__lte", time.Unix(params.EndTime, 0))
	}
	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, (params.Offset-1)*params.Limit).RelatedSel().All(&data)
	for _,obj := range data{
		obj.FinancialCategoryId = obj.FinancialCategoryObj.Id
		if obj.Admin != nil{
			obj.Name = obj.Admin.UserName
			obj.AdminId = obj.Admin.Id
		}
		if obj.FinancialCategoryObj != nil{
			obj.Symbol = obj.FinancialCategoryObj.Symbol
		}
	}
	return data, total
}

//上下架
func FinancialProductTheUppers(status int, ids []int,aid int) (num int64, err error) {
	o := orm.NewOrm()
	query := o.QueryTable(FinancialProductTBName())
	params := map[string]interface{}{
		"status": status,
		"admin_id":aid,
	}
	switch status {
	case 1:
		//上架
		num, err = query.Filter("category__exact", 2).Filter("id__in", ids).Filter("status__exact", 0).Update(params)
	case 2:
		//后台下架
		num, err = query.Filter("category__exact", 2).Filter("id__in", ids).Filter("status__exact", 1).Update(params)
	default:
		return 0, err
	}
	return num, err
}

//获取单条
func FinancialProductOne(id int) *FinancialProduct {
	o := orm.NewOrm()
	obj := &FinancialProduct{Id: id}
	if err := o.Read(obj, "id"); err != nil {
		return nil
	}
	return obj
}

func FinancialProductIsAddObj(fid int) bool {
	o := orm.NewOrm()
	query := o.QueryTable(FinancialProductTBName())
	num, err := query.Filter("financial_category_id__exact", fid).Filter("category__exact", 1).Filter("status__in", []int{1, 0}).Count()
	if err != nil {
		return false
	}
	if num > 0 {
		return false
	}
	return true
}

func FinancialProductDelete(id int) (bool, string) {
	o := orm.NewOrm()
	query := o.QueryTable(FinancialProductTBName())
	_, err := query.Filter("id__exact", id).Filter("status__exact", 0).Delete()
	if err != nil {
		return false, err.Error()
	}
	return true, ""
}

//获取所有的定期利率
func FinancialProductDemandDepositInterestRates()map[string]float64{
	o := orm.NewOrm()
	mapp := make(map[string]float64)
	data := make([]*FinancialProduct, 0)
	query := o.QueryTable(FinancialProductTBName())
	query.Filter("category__exact",1).RelatedSel().All(&data)
	for _,obj := range data{
		if obj.FinancialCategoryObj != nil{
			mapp[obj.FinancialCategoryObj.Symbol] = obj.YearProfit
		}
	}
	return mapp
}