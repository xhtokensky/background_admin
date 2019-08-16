package models

import (
	"github.com/astaxie/beego/orm"
	"time"
	"tokensky_bg_admin/utils"
)

// OperationBannerQueryParam 用于查询的类
type OperationBannerQueryParam struct {
	BaseQueryParam
	//发布时间
	StartTime int64 `json:"startTime"` //开始时间
	EndTime   int64 `json:"endTime"`   //截止时间
	Status    int   `json:"status"`    //状态
}

func (a *OperationBanner) TableName() string {
	return OperationBannerTBName()
}

//运营Banner
type OperationBanner struct {
	Bid int `orm:"pk;column(bid)"json:"bid"form:"bid"`
	//链接
	Url string `orm:"column(url)"json:"url"form:"url"`
	//图片地址
	ImgKey string `orm:"column(img_key)"json:"imgKey"form:"imgKey"`
	//展示字段
	ImgUrl string `orm:"-"json:"imgUrl"`
	//顺序
	Seq int `orm:"column(seq)"json:"seq"form:"seq"`
	//创建时间
	CreateTime time.Time `orm:"auto_now_add;type(datetime);column(create_time)"json:"createTime"form:"createTime"`
	//最后更新时间
	UpdateTime time.Time `orm:"auto_now;type(datetime);column(update_time)"json:"updateTime"form:"updateTime"`
	//状态 0关闭 1开启
	Status int `orm:"column(status)"json:"status"form:"status"`
	//名称
	Name string `orm:"column(name)"json:"name"form:"name"`
	//用户
	AdminId int `orm:"column(admin_id)"json:"-"form:"-"`
}

//获取分页数据
func OperationBannerPageList(params *OperationBannerQueryParam) ([]*OperationBanner, int64) {
	o := orm.NewOrm()
	query := o.QueryTable(OperationBannerTBName())
	data := make([]*OperationBanner, 0)
	//默认排序
	sortorder := "bid"
	switch params.Sort {
	case "bid":
		sortorder = "bid"
	}
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}
	//状态
	if params.Status >= 0 {
		query = query.Filter("status__iexact", params.Status)
	}
	//时间段
	if params.StartTime > 0 {
		query = query.Filter("create_time__gte", time.Unix(params.StartTime, 0))
	}
	if params.EndTime > 0 {
		query = query.Filter("create_time__lte", time.Unix(params.EndTime, 0))
	}
	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, (params.Offset-1)*params.Limit).All(&data)
	//图片下载凭证
	deadline := time.Now().Add(time.Second * 3600).Unix() //1小时有效期
	for _, obj := range data {
		obj.ImgUrl = utils.QiNiuDownload(obj.ImgKey, deadline)
	}
	return data, total
}

//批量删除
func OperationBannerDelete(ids []int) (int64, error) {
	query := orm.NewOrm().QueryTable(OperationBannerTBName())
	num, err := query.Filter("bid__in", ids).Delete()
	return num, err
}

//获取单条
func OperationBannerOneById(id int) *OperationBanner {
	o := orm.NewOrm()
	m := OperationBanner{Bid: id}
	err := o.Read(&m)
	if err != nil {
		return nil
	}
	return &m
}
