package models

import (
	"github.com/astaxie/beego/orm"
	"time"
	"tokensky_bg_admin/utils"
)

//查询的类
type HashrateCategoryQueryParam struct {
	BaseQueryParam
	Name      string `json:"name"`
	StartTime int64  `json:"startTime"` //开始时间
	EndTime   int64  `json:"endTime"`   //截止时间
	Status    string `json:"status"`    //状态
}

func (a *HashrateCategory) TableName() string {
	return HashrateCategoryTBName()
}

//算力合约分类
type HashrateCategory struct {
	KeyId int `orm:"pk;column(key_id)"json:"keyId"form:"keyId"`
	//图片key值
	ImgKey string `orm:"column(img_key)"json:"imgKey"form:"imgKey"`
	//展示字段
	ImgUrl string `orm:"-"json:"imgUrl"`
	//名称
	Name string `orm:"column(name)"json:"name"form:"name"`
	//类别
	Symbol string `orm:"column(symbol)"json:"symbol"form:"symbol"`
	//单位
	Unit string `orm:"column(unit)"json:"unit"form:"unit"`
	//状态 1为正常
	Status     int       `orm:"column(status)"json:"status"form:"status"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime);column(create_time)"json:"createTime"form:"createTime"`
	UpdateTime time.Time `orm:"auto_now;type(datetime);column(update_time)"json:"updateTime"form:"updateTime"`
	//用户id
	AdminUserId int `orm:"column(admin_user_id)"json:"adminUserId"form:"adminUserId"`
	//连表 算力合约表 一对多反向关系
	HashrateTreatys []*HashrateTreaty `orm:"reverse(many)"json:"-"form:"-"`
}

//获取分页数据
func HashrateCategoryPageList(params *HashrateCategoryQueryParam) ([]*HashrateCategory, int64) {
	o := orm.NewOrm()
	query := o.QueryTable(HashrateCategoryTBName())
	data := make([]*HashrateCategory, 0)
	//默认排序
	sortorder := "key_id"
	switch params.Sort {
	case "keyId":
		sortorder = "key_id"
	}
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}
	//姓名模糊查询
	if params.Name != "" {
		query = query.Filter("name__icontains", params.Name)
	}
	//状态
	if params.Status != "" && params.Status != "-1" {
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
func HashrateCategoryDelete(ids []int) (int64, error) {
	o := orm.NewOrm()
	query := o.QueryTable(HashrateCategoryTBName())
	num, err := query.Filter("key_id__in", ids).Delete()
	return num, err
}

//
func HashrateCategoryByIds(ids []int) map[int]*HashrateCategory {
	mapp := make(map[int]*HashrateCategory)
	data := make([]*HashrateCategory, 0)
	if len(ids) > 0 {
		query := orm.NewOrm().QueryTable(HashrateCategoryTBName())
		query = query.Filter("key_id__in", ids)
		query.All(&data)
	}
	for _, obj := range data {
		mapp[obj.KeyId] = obj
	}
	return mapp
}

func HashrateCategoryById(id int) *HashrateCategory {
	o :=orm.NewOrm()
	m := HashrateCategory{KeyId: id}
	err := o.Read(&m)
	if err != nil {
		return nil
	}
	return &m
}
