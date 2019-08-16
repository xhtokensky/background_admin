package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

// TableName 设置BackendUser表名
func (a *TokenskyMessage) TableName() string {
	return TokenskyMessageTBName()
}

// AdminBackendUserQueryParam 用于查询的类
type TokenskyMessageQueryParam struct {
	Status string `json:"status"`
	//发布时间
	StartTime int64  `json:"startTime"` //开始时间
	EndTime   int64  `json:"endTime"`   //截止时间
	Phone     string `json:"phone"`
	BaseQueryParam
}

// TokenskyMessage 实体类
type TokenskyMessage struct {
	//Id int "message_id"
	MessageId int    `orm:"pk;column(message_id)"json:"messageId"form:"messageId"`
	Title     string `orm:"column(title)"json:"title"form:"title"`
	Content   string `orm:"column(content)"json:"content"form:"content"`
	//0所有 1个人
	Type       int           `orm:"column(type)"json:"type"form:"type"`
	User       *TokenskyUser `orm:"rel(fk);column(user_id)"json:"-"form:"-"`
	UserId     int           `orm:"-"json:"user_id"form:"user_id"`
	Phone      string        `orm:"-"json:"phone"form:"phone"`
	EditorId   int           `orm:"column(editor_id)"json:"editorId"form:"editorId"`
	Status     int           `orm:"column(status)"json:"status"form:"status"`
	CreateTime time.Time     `orm:"auto_now_add;type(datetime);column(create_time)"json:"createTime"form:"createTime"`
	UpdateTime time.Time     `orm:"auto_now;type(datetime);column(update_time)"json:"updateTime"form:"updateTime"`
	AdminId    int           `orm:"column(admin_id)"json:"-"form:"-"`
	//0 admin 1 app
	System int `orm:"column(system)"json:"system"form:"system"`
	//关联ID
	RelevanceId string `orm:"column(relevance_id)"json:"relevanceId"form:"relevanceId"`
	//消息内容类型
	MsgCategory string `orm:"column(msg_category)"json:"msgCategory"form:"msgCategory"`
	//跳转路由
	MsgRoute string `orm:"column(msg_route)"json:"msgRoute"form:"msgRoute"`
	//是否已读
	IsRead int `orm:"column(is_read)"json:"isRead"form:"isRead"`
	//读取时间
	ReadTime time.Time`orm:"type(datetime);column(read_time)"json:"readTime"form:"readTime"`
}

// 获取分页数据
func TokenskyMessagePageList(params *TokenskyMessageQueryParam) ([]*TokenskyMessage, int64) {
	o := orm.NewOrm()
	query := o.QueryTable(TokenskyMessageTBName())
	data := make([]*TokenskyMessage, 0)
	//默认排序
	sortorder := "messageId"
	switch params.Sort {
	case "messageId":
		sortorder = "messageId"
	}
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}
	//状态
	if params.Status != "" && params.Status != "-1" {
		query = query.Filter("status__iexact", params.Status)
	}
	if params.Phone != "" {
		query = query.Filter("User__phone__icontains", params.Phone)
	}
	//后台过滤消息
	query = query.Filter("system__exact", 0)

	//时间段
	if params.StartTime > 0 {
		query = query.Filter("create_time__gte", time.Unix(params.StartTime, 0))
	}
	if params.EndTime > 0 {
		query = query.Filter("create_time__lte", time.Unix(params.EndTime, 0))
	}
	query = query.Filter("status__exact", 1)
	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, (params.Offset-1)*params.Limit).RelatedSel().All(&data)
	for _, v := range data {
		v.Phone = v.User.Phone
	}
	return data, total
}

// TokenskyMessageDelete 批量删除
func TokenskyMessageDelete(ids []int) (int64, error) {
	o := orm.NewOrm()
	query := o.QueryTable(TokenskyMessageTBName())
	num, err := query.Filter("messageId__in", ids).Delete()
	return num, err
}

// TokenskyMessageOne 根据id获取单条
func TokenskyMessageOne(id int) (*TokenskyMessage, error) {
	o := orm.NewOrm()
	m := TokenskyMessage{MessageId: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
