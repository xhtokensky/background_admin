package models

import "time"

//用户行为记录

func (a *AdminModelRecord) TableName() string {
	return AdminModelRecordTBName()
}

type AdminModelRecord struct {
	Id int `orm:"pk;column(id)"json:"id"form:"id"`
	//uid
	Uid int `orm:"column(uid)"json:"uid"form:"uid"`
	//操作
	Handle string `orm:"column(handle)"json:"handle"form:"handle"`
	//表
	Model string `orm:"column(model)"json:"model"form:"model"`
	//id
	Tbid string `orm:"column(tbid)"json:"tbid"form:"tbid"`
	//旧数据
	OldData string `orm:"column(old_data)"json:"oldData"form:"oldData"`
	//新数据
	NewData string `orm:"column(new_data)"json:"newData"form:"newData"`
	//创建时间
	CreateTime time.Time `orm:"auto_now_add;type(datetime);column(create_time)"json:"createTime"form:"createTime"`

}