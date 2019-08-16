package models

import "time"

// AdminRoleResourceRel 角色与资源关系表
type AdminRoleResourceRel struct {
	Id            int
	AdminRole     *AdminRole     `orm:"rel(fk)"`  //外键
	AdminResource *AdminResource `orm:"rel(fk)" ` // 外键
	Created       time.Time      `orm:"auto_now_add;type(datetime)"`
}

// TableName 设置表名
func (a *AdminRoleResourceRel) TableName() string {
	return AdminRoleResourceRelTBName()
}
