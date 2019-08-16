package models

import (
	"github.com/astaxie/beego/orm"
	//"tokensky_bg_admin/utils"
	"time"
	"tokensky_bg_admin/utils"
)

//查询的类
type TokenskyRealAuthQueryParam struct {
	BaseQueryParam
	//发布时间
	StartTime    int64  `json:"startTime"`    //开始时间
	EndTime      int64  `json:"endTime"`      //截止时间
	Status       string `json:"status"`       //状态
	Name         string `json:"name"`         //姓名 模糊搜索
	Phone        string `json:"phone"`        //手机号 模糊搜索
	IdentityCard string `json:"identityCard"` //身份证号 模糊搜索
}

func (a *TokenskyRealAuth) TableName() string {
	return TokenskyRealAuthTBName()
}

//身份审核
type TokenskyRealAuth struct {
	KeyId int `orm:"pk;column(key_id)"json:"keyId"form:"keyId"`
	//用户uid
	//UserId int `orm:"column(user_id)"json:"userId"form:"userId"`
	//UserId int `orm:"-"json:"userId"form:"userId"`
	//姓名
	Name string `orm:"column(name)"json:"name"form:"name"`
	//身份证
	IdentityCard string `orm:"column(identity_card)"json:"identityCard"form:"identityCard"`
	//证件照 正反照片
	IdentityCardPicturev  string `orm:"column(identity_card_picture)"json:"-"form:"-"`
	IdentityCardPicture2v string `orm:"column(identity_card_picture2)"json:"-"form:"-"`

	IdentityCardPicture  string `orm:"-"json:"identityCardPicture"form:"identityCardPicture"`
	IdentityCardPicture2 string `orm:"-"json:"identityCardPicture2"form:"identityCardPicture2"`

	//创建(提交)时间
	CreateTime time.Time `orm:"auto_now_add;type(datetime);column(create_time)"json:"createTime"form:"createTime"`
	//审核时间
	UpdateTime time.Time `orm:"auto_now;type(datetime);column(update_time)"json:"updateTime"form:"updateTime"`
	//状态 默认1起来  0尚未认证 1认证成功 2认证失败
	Status int    `orm:"column(status)"json:"status"form:"status"`
	Phone  string `orm:"-"json:"phone"form:"-"`
	//人脸识别 巴拉巴拉

	//关联用户表
	User *TokenskyUser `orm:"rel(fk);column(user_id)"json:"-"form:"-"`
}

//获取分页数据
func TokenskyRealAuthPageList(params *TokenskyRealAuthQueryParam) ([]*TokenskyRealAuth, int64) {
	o := orm.NewOrm()
	query := o.QueryTable(TokenskyRealAuthTBName())
	data := make([]*TokenskyRealAuth, 0)
	//默认排序
	sortorder := "keyId"
	switch params.Sort {
	case "keyId":
		sortorder = "key_id"
	}
	switch params.Order {
	case "":
		sortorder = "-" + sortorder
	case "desc":
		sortorder = "-" + sortorder
	default:
		sortorder = sortorder
	}
	//姓名模糊查询
	if params.Name != "" {
		query = query.Filter("name__icontains", params.Name)
	}
	//身份证号模糊查询
	if params.IdentityCard != "" {
		query = query.Filter("identity_card__icontains", params.IdentityCard)
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
	//手机号
	//phones := make(map[int]string)
	//手机号模糊查询 起始位置开始
	if params.Phone != "" {
		query = query.Filter("User__Phone__icontains", params.Phone)
	}
	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, (params.Offset-1)*params.Limit).RelatedSel("user").All(&data)
	//图片下载凭证
	deadline := time.Now().Add(time.Second * 3600).Unix() //1小时有效期
	for _, obj := range data {
		//用户电话
		if obj.User != nil {
			obj.Phone = obj.User.Phone
		}
		//证书
		obj.IdentityCardPicture = utils.QiNiuDownload(obj.IdentityCardPicturev, deadline)
		obj.IdentityCardPicture2 = utils.QiNiuDownload(obj.IdentityCardPicture2v, deadline)
	}

	return data, total
}

//根据id获取单条记录
func TokenskyRealAuthOneById(id int) (*TokenskyRealAuth, error) {
	o := orm.NewOrm()
	m := TokenskyRealAuth{KeyId: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	m.IdentityCardPicture = utils.QiNiuDownload(m.IdentityCardPicturev, 0)
	m.IdentityCardPicture2 = utils.QiNiuDownload(m.IdentityCardPicture2v, 0)
	return &m, nil
}

//获取多条记录
func TokenskyRealAuthByIds(ids []int) (map[int]*TokenskyRealAuth, error) {
	var err error
	o := orm.NewOrm()
	query := o.QueryTable(TokenskyRealAuthTBName())
	data := make([]*TokenskyRealAuth, 0)
	if len(ids) > 0 {
		_, err = query.Filter("user_id__in", ids).All(&data, "User", "key_id", "status")
	}
	mapp := make(map[int]*TokenskyRealAuth)
	for _, v := range data {
		if v.User != nil {
			mapp[v.User.UserId] = v
		}
	}
	return mapp, err
}
