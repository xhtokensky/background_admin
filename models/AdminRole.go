package models

import (
	"github.com/astaxie/beego/orm"
)

// TableName 设置表名
func (a *AdminRole) TableName() string {
	return AdminRoleTBName()
}

// AdminRoleQueryParam 用于搜索的类
type AdminRoleQueryParam struct {
	BaseQueryParam
	NameLike string
}

// AdminRole 用户角色 实体类
type AdminRole struct {
	Id                 int                        `orm:"pk;"form:"id"json:"id"`
	Name               string                     `form:"name"json:"name"`
	Seq                int                        `form:"seq"json:"seq"`
	RoleResourceRel    []*AdminRoleResourceRel    `orm:"reverse(many)" json:"-"` // 设置一对多的反向关系
	RoleBackendUserRel []*AdminRoleBackendUserRel `orm:"reverse(many)" json:"-"` // 设置一对多的反向关系
}

// AdminRolePageList 获取分页数据
func AdminRolePageList(params *AdminRoleQueryParam) ([]*AdminRole, int64) {
	query := orm.NewOrm().QueryTable(AdminRoleTBName())
	data := make([]*AdminRole, 0)
	//默认排序
	sortorder := "Id"
	switch params.Sort {
	case "id":
		sortorder = "id"
	case "seq":
		sortorder = "seq"
	}
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}
	query = query.Filter("name__istartswith", params.NameLike)
	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, (params.Offset-1)*params.Limit).All(&data)
	return data, total
}

// AdminRoleDataList 获取角色列表
func AdminRoleDataList(params *AdminRoleQueryParam) []*AdminRole {
	params.Limit = -1
	params.Sort = "seq"
	params.Order = "asc"
	data, _ := AdminRolePageList(params)
	return data
}

// AdminRoleBatchDelete 批量删除
func AdminRoleBatchDelete(ids []int) (int64, error) {
	query := orm.NewOrm().QueryTable(AdminRoleTBName())
	num, err := query.Filter("id__in", ids).Delete()
	return num, err
}

// AdminRoleOne 获取单条
func AdminRoleOne(id int) (*AdminRole, error) {
	o := orm.NewOrm()
	m := AdminRole{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

//获取角色信息
func AdminRolesByIds(ids []int) map[int]*AdminRole {
	data := make(map[int]*AdminRole)
	if len(ids) > 0 {
		objs := make([]*AdminRole, 0)
		o := orm.NewOrm()
		query := o.QueryTable(AdminRoleTBName())
		query = query.Filter("id__in", ids)
		query.All(&objs)
		for _, obj := range objs {
			data[obj.Id] = obj
		}
	}
	return data
}
