package models

import (
	"github.com/astaxie/beego/orm"
	"time"
	"tokensky_bg_admin/utils"
)

func (a *TokenskyJiguangRegistrationid) TableName() string {
	return TokenskyJiguangRegistrationidTBName()
}

//极光地址表
type TokenskyJiguangRegistrationid struct {
	Id int `orm:"pk;column(id)"json:"id"form:"id"`
	UserId int `orm:"column(user_id)"json:"userId"form:"userId"`
	RegistrationId string `orm:"column(registration_id)"json:"registrationId"form:"registrationId"`
	//创建时间
	CretaeTime time.Time `orm:"type(datetime);column(cretae_time)"json:"cretaeTime"form:"cretaeTime"`
	//更新时间
	UpdateTime time.Time `orm:"auto_now;type(datetime);column(update_time)"json:"updateTime"form:"updateTime"`
}

//获取极光地址
func TokenskyJiguangRegistrationidGetRegistrationId(uid int)(registrationId string){
	o := orm.NewOrm()
	obj := &TokenskyJiguangRegistrationid{}
	err := o.QueryTable(TokenskyJiguangRegistrationidTBName()).Filter("user_id__exact",uid).One(obj,"user_id","registration_id")
	if err ==nil{
		registrationId = obj.RegistrationId
	}
	return registrationId
}
func TokenskyJiguangRegistrationidGetRegistrationIds(ids []int)(mapp map[int]string){
	mapp = make(map[int]string,len(ids))
	if len(ids)>0{
		o := orm.NewOrm()
		data := make([]*TokenskyJiguangRegistrationid,0)
		_,err := o.QueryTable(TokenskyJiguangRegistrationidTBName()).Filter("user_id__in",ids).All(&data,"user_id","registration_id")
		if err != nil{
			return mapp
		}
		for _,obj := range data{
			mapp[obj.Id]= obj.RegistrationId
		}
	}
	return mapp
}

//对单人极光推送
func TokenskyJiguangRegistrationidSendByOne(uid int, alertTitle, alertContent, title, content string)  {
	if addr := TokenskyJiguangRegistrationidGetRegistrationId(uid);addr != ""{
		utils.JiGuangSendByAddr(addr,alertTitle,alertContent,title,content)
	}
}

//对多人极光推送
func TokenskyJiguangRegistrationidSendByIds(ids []int, alertTitle, alertContent, title, content string)  {
	addrs := TokenskyJiguangRegistrationidGetRegistrationIds(ids)
	list := make([]string,0)
	for _,addr := range addrs{
		list = append(list, addr)
	}
	if len(addrs) >0{
		utils.JiGuangSendByAddrs(list,alertTitle,alertContent,title,content)
	}
}