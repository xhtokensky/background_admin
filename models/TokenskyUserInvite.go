package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

//邀请表
func (a *TokenskyUserInvite) TableName() string {
	return TokenskyUserInviteTBName()
}

//用户邀请表
type TokenskyUserInvite struct {
	Id int `orm:"column(id)"json:"id"form:"id"`
	//邀请人
	From *TokenskyUser `orm:"rel(fk);column(from)"json:"-"form:"-"`
	//被邀请人
	To int `orm:"column(to)"json:"-"form:"-"`
	//创建时间
	CreateTime time.Time `orm:"type(datetime);column(create_time)"json:"-"form:"-"`
	//
	UserId int `orm:"-"json:"userId"form:"-"`
	Count int `orm:"-"json:"count"form:"-"`
}

//
func TokenskyUserInviteFormToAmount()map[int]*TokenskyUserInvite{
	o := orm.NewOrm()
	query := o.QueryTable(TokenskyUserInviteTBName())
	num := 1
	count,_ := query.Count()
	mapp2 := make(map[int]*TokenskyUserInvite)
	for count>0{
		data := make([]*TokenskyUserInvite,0)
		query.Limit(500,(num-1)*500).All(&data)
		count -= 500
		num++
		for _,obj := range data{
			if obj2,found := mapp2[obj.From.UserId];!found{
				mapp2[obj.From.UserId] = obj
				obj.Count++
				obj.UserId = obj.From.UserId
			}else {
				obj2.Count++
			}
		}
	}
	return mapp2
}