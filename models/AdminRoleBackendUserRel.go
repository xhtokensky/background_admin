package models

import (
	"time"
)

// RoleBackendUserRel 角色与用户关系
type AdminRoleBackendUserRel struct {
	Id               int
	AdminRole        *AdminRole        `orm:"rel(fk)"`  //外键
	AdminBackendUser *AdminBackendUser `orm:"rel(fk)" ` // 外键
	Created          time.Time         `orm:"auto_now_add;type(datetime)"`
}

// TableName 设置表名
func (a *AdminRoleBackendUserRel) TableName() string {
	return AdminRoleBackendUserRelTBName()
}
