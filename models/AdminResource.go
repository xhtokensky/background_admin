package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"tokensky_bg_admin/utils"
)

// TableName 设置表名
func (a *AdminResource) TableName() string {
	return AdminResourceTBName()
}

// AdminResource 权限控制资源表
type AdminResource struct {
	Id int `orm:"pk"json:"id"form:"id"`
	//名称
	Name                 string                  `orm:"size(64)"json:"name"form:"name"`
	Parent               *AdminResource          `orm:"null;rel(fk)"json:"parent"form:"parent"` // RelForeignKey relation
	Rtype                int                     `json:"rtype"form:"rtype"`
	Seq                  int                     `json:"seq"form:"seq"`
	Sons                 []*AdminResource        `orm:"reverse(many)"json:"sons"form:"sons"` // fk 的反向关系
	SonNum               int                     `orm:"-"json:"sons"form:"sons"`
	Icon                 string                  `orm:"size(32)"json:"icon"form:"icon"`
	LinkUrl              string                  `orm:"-"json:"linkUrl"form:"linkUrl"`
	UrlFor               string                  `orm:"size(256)" Json:"-"form:"urlFor"`
	HtmlDisabled         int                     `orm:"-"json:"htmlDisabled"form:"htmlDisabled"`                   //在html里应用时是否可用
	Level                int                     `orm:"-"json:"level"form:"level"`                                 //第几级，从0开始
	AdminRoleResourceRel []*AdminRoleResourceRel `orm:"reverse(many)"json:"roleResourceRel"form:"roleResourceRel"` // 设置一对多的反向关系
	//给前端的特殊处理
	Childs []*AdminResource `orm:"-"json:"child"form:"-"`
	//是否选择
	Choice bool `orm:"-"json:"choice"form:"-"`
}

// AdminResourceOne 获取单条
func AdminResourceOne(id int) (*AdminResource, error) {
	o := orm.NewOrm()
	m := AdminResource{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// AdminResourceTreeGrid 获取treegrid顺序的列表
func AdminResourceTreeGrid() []*AdminResource {
	o := orm.NewOrm()
	query := o.QueryTable(AdminResourceTBName()).OrderBy("seq", "id")
	list := make([]*AdminResource, 0)
	query.All(&list)
	return AdminResourceList2TreeGrid(list)
}
func AdminResourceTreeGrid2() []*AdminResource {
	o := orm.NewOrm()
	query := o.QueryTable(AdminResourceTBName()).OrderBy("seq", "id")
	list := make([]*AdminResource, 0)
	sql := fmt.Sprintf(`SELECT id,name,parent_id,rtype,icon,seq,url_for FROM %s Where rtype <= ? Order By seq asc,Id asc`, AdminResourceTBName())
	o.Raw(sql, 1000).QueryRows(&list)
	query.All(&list)
	return list
}

// AdminResourceTreeGrid4Parent 获取可以成为某个节点父节点的列表
func AdminResourceTreeGrid4Parent(id int) []*AdminResource {
	tree := AdminResourceTreeGrid()
	if id == 0 {
		return tree
	}
	var index = -1
	//找出当前节点所在索引
	for i, _ := range tree {
		if tree[i].Id == id {
			index = i
			break
		}
	}
	if index == -1 {
		return tree
	} else {
		tree[index].HtmlDisabled = 1
		for _, item := range tree[index+1:] {
			if item.Level > tree[index].Level {
				item.HtmlDisabled = 1
			} else {
				break
			}
		}
	}
	return tree
}

// AdminResourceTreeGridByUserId 根据用户获取有权管理的资源列表，并整理成teegrid格式
func AdminResourceTreeGridByUserId(backuserid, maxrtype int) []*AdminResource {
	cachekey := fmt.Sprintf("rms_ResourceTreeGridByUserId_%v_%v", backuserid, maxrtype)
	var list []*AdminResource
	if err := utils.GetCache(cachekey, &list); err == nil {
		return list
	}
	o := orm.NewOrm()
	user, err := AdminBackendUserOne(backuserid)
	if err != nil || user == nil {
		return list
	}
	var sql string
	if user.IsSuper == 1 {
		//如果是管理员，则查出所有的
		sql = fmt.Sprintf(`SELECT id,name,parent_id,rtype,icon,seq,url_for FROM %s Where rtype <= ? Order By seq asc,Id asc`, AdminResourceTBName())
		o.Raw(sql, maxrtype).QueryRows(&list)
	} else {
		//联查多张表，找出某用户有权管理的
		sql = fmt.Sprintf(`SELECT DISTINCT T0.admin_resource_id,T2.id,T2.name,T2.parent_id,T2.rtype,T2.icon,T2.seq,T2.url_for
		FROM %s AS T0
		INNER JOIN %s AS T1 ON T0.admin_role_id = T1.admin_role_id
		INNER JOIN %s AS T2 ON T2.id = T0.admin_resource_id
		WHERE T1.admin_backend_user_id = ? and T2.rtype <= ?  Order By T2.seq asc,T2.id asc`, AdminRoleResourceRelTBName(), AdminRoleBackendUserRelTBName(), AdminResourceTBName())
		o.Raw(sql, backuserid, maxrtype).QueryRows(&list)
	}
	result := AdminResourceList2TreeGrid(list)
	utils.SetCache(cachekey, result, 30)
	return result
}

func AdminResourceTreeGridByUserId2(user *AdminBackendUser, maxrtype int) []*AdminResource {
	cachekey := fmt.Sprintf("rms_ResourceTreeGridByUserId_%v_%v", user.Id, maxrtype)
	var list []*AdminResource
	if err := utils.GetCache(cachekey, &list); err == nil {
		return list
	}
	o := orm.NewOrm()
	var sql string
	if user.IsSuper == 1 {
		//如果是管理员，则查出所有的
		sql = fmt.Sprintf(`SELECT id,name,parent_id,rtype,icon,seq,url_for FROM %s Where rtype <= ? Order By seq asc,Id asc`, AdminResourceTBName())
		o.Raw(sql, maxrtype).QueryRows(&list)
	} else {
		//联查多张表，找出某用户有权管理的
		sql = fmt.Sprintf(`SELECT DISTINCT T0.admin_resource_id,T2.id,T2.name,T2.parent_id,T2.rtype,T2.icon,T2.seq,T2.url_for
		FROM %s AS T0
		INNER JOIN %s AS T1 ON T0.admin_role_id = T1.admin_role_id
		INNER JOIN %s AS T2 ON T2.id = T0.admin_resource_id
		WHERE T1.admin_backend_user_id = ? and T2.rtype <= ?  Order By T2.seq asc,T2.id asc`, AdminRoleResourceRelTBName(), AdminRoleBackendUserRelTBName(), AdminResourceTBName())
		o.Raw(sql, user.Id, maxrtype).QueryRows(&list)
	}
	return list
}

//查询某角色拥有的所有权限
func AdminResourceTreeGridByUser(user *AdminBackendUser, maxrtype int) []*AdminResource {
	//查询
	if obj, err := AdminBackendUserOne(user.Id); err == nil {
		user = obj
	}
	//
	var list []*AdminResource
	o := orm.NewOrm()
	if user.IsSuper == 1 {
		//如果是管理员，则查出所有的
		sql := fmt.Sprintf(`SELECT id,name,parent_id,rtype,icon,seq,url_for FROM %s Where rtype <= ? Order By seq asc,Id asc`, AdminResourceTBName())
		o.Raw(sql, maxrtype).QueryRows(&list)
	} else {
		//用户拥有的角色
		roleList := make([]int, 0)
		for _, str := range strings.Split(user.RolesStr, ",") {
			if con, err := strconv.Atoi(str); err == nil {
				roleList = append(roleList, con)
			}
		}
		if len(roleList) > 0 {
			query := o.QueryTable(AdminRoleResourceRelTBName())
			mm2 := make([]*AdminRoleResourceRel, 0)
			query.Filter("AdminRole__in", roleList).All(&mm2)

			mm2Ids := make([]int, 0)
			for _, v := range mm2 {
				mm2Ids = append(mm2Ids, v.AdminResource.Id)
			}
			if len(mm2Ids) > 0 {
				query := o.QueryTable(AdminResourceTBName())
				query.Filter("id__in", mm2Ids).All(&list)
			}
		}
	}
	return list
}

func AdminResourceTreeGridByUserId3(uid int, maxrtype int) map[int]struct{} {
	cachekey := fmt.Sprintf("rms_ResourceTreeGridByUserId_%v_%v", uid, maxrtype)
	var list []*AdminResource
	mapp := make(map[int]struct{})
	if err := utils.GetCache(cachekey, &list); err == nil {
		return mapp
	}
	o := orm.NewOrm()
	var sql string
	//联查多张表，找出某用户有权管理的
	sql = fmt.Sprintf(`SELECT DISTINCT T0.admin_resource_id,T2.id,T2.name,T2.parent_id,T2.rtype,T2.icon,T2.seq,T2.url_for
		FROM %s AS T0
		INNER JOIN %s AS T1 ON T0.admin_role_id = T1.admin_role_id
		INNER JOIN %s AS T2 ON T2.id = T0.admin_resource_id
		WHERE T1.admin_backend_user_id = ? and T2.rtype <= ?  Order By T2.seq asc,T2.id asc`, AdminRoleResourceRelTBName(), AdminRoleBackendUserRelTBName(), AdminResourceTBName())
	o.Raw(sql, uid, maxrtype).QueryRows(&list)
	for _, obj := range list {
		mapp[obj.Id] = struct{}{}
	}
	return mapp
}

// AdminResourceList2TreeGrid 将资源列表转成treegrid格式
func AdminResourceList2TreeGrid(list []*AdminResource) []*AdminResource {
	result := make([]*AdminResource, 0)
	for _, item := range list {
		if item.Parent == nil || item.Parent.Id == 0 {
			item.Level = 0
			result = append(result, item)
			result = AdminResourceAddSons(item, list, result)
		}
	}
	return result
}

//AdminResourceAddSons 添加子菜单
func AdminResourceAddSons(cur *AdminResource, list, result []*AdminResource) []*AdminResource {
	for _, item := range list {
		if item.Parent != nil && item.Parent.Id == cur.Id {
			cur.SonNum++
			item.Level = cur.Level + 1
			result = append(result, item)
			result = AdminResourceAddSons(item, list, result)
		}
	}
	return result
}

//角色信息处理
func AdminRolesChildSort(objs []*AdminResource) []*AdminResource {
	list := make([]*AdminResource, 0)
	list2 := make([]*AdminResource, 0)
	mapp := make(map[int]*AdminResource)
	for _, obj := range objs {
		obj.Childs = make([]*AdminResource, 0)
		mapp[obj.Id] = obj
		if obj.Parent != nil {
			if obj.Parent.Id == 0 {
				list = append(list, obj)
			} else {
				list2 = append(list2, obj)
			}
		} else {
			list = append(list, obj)
		}
	}
	for _, obj := range list2 {
		if fobj, found := mapp[obj.Parent.Id]; found {
			fobj.Childs = append(fobj.Childs, obj)
		}
	}
	return list
}
