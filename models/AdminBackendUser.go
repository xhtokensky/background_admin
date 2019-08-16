package models

import (
	"crypto/md5"
	"fmt"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
)

// TableName 设置BackendUser表名
func (a *AdminBackendUser) TableName() string {
	return AdminBackendUserTBName()
}

// AdminBackendUserQueryParam 用于查询的类
type AdminBackendUserQueryParam struct {
	BaseQueryParam
	UserName string `json:"userName"` //模糊查询
	RealName string `json:"realName"` //模糊查询
	Mobile   string `json:"mobile"`   //精确查询
	Status   string `json:"status"`   //为空不查询，有值精确查询
	Id       int    `json:"id"`       //Id查询单条
}

// AdminBackendUser 实体类
type AdminBackendUser struct {
	Id       int    `orm:"pk;column(id)"json:"id"form:"id"`
	RealName string `orm:"column(real_name)"json:"realName"form:"realName"`
	UserName string `orm:"column(user_name)"json:"userName"form:"userName"`
	//密码
	UserPwd string `orm:"column(user_pwd)"json:"userPwd"form:"userPwd"`
	//是否超级用户 0关闭 1启用
	IsSuper int `orm:"column(is_super)"json:"isSuper"form:"isSuper"`
	//状态 0关闭 1启用
	Status int `orm:"column(status)"json:"status"form:"status"`
	//电话
	Mobile string `orm:"column(mobile);size(16)"json:"mobile"form:"mobile"`
	//邮件
	Email string `orm:"column(email);size(256)"json:"email"form:"email"`
	//头像
	Avatar                  string                     `orm:"column(avatar);size(256)"json:"avatar"form:"avatar"`
	RoleIds                 []int                      `orm:"-" form:"roleIds"json:"roleIds"`
	AdminRoleBackendUserRel []*AdminRoleBackendUserRel `orm:"reverse(many)"json:"_"` // 设置一对多的反向关系
	ResourceUrlForList      []string                   `orm:"-"json:"_"`
	//CreateCourses      []*AdminCourse             `rom:"reverse(many)"` // 设置一对多的反向关系
	//Creator   *AdminRoleBackendUserRel `orm:"rel(fk)"` //设置一对多关系
	RolesStr string `orm:"column(roles_str)"json:"-"form:"-"`
	//连表
	FinancialConfigHistoricalRecords []*FinancialProductHistoricalRecord `orm:"reverse(many)"json:"-"form:"-"`
	FinancialProducts  []*FinancialProduct `orm:"reverse(many)"json:"-"form:"-"`
	//特殊字段
	Ids       string   `orm:"-"json:"ids"form:"ids"`
	RolesName []string `orm:"-"json:"rolesName"form:"-"`
}

func (this *AdminBackendUser) GetToken() string {
	data := []byte(this.UserPwd + this.UserName)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

// AdminBackendUserPageList 获取分页数据
func AdminBackendUserPageList(params *AdminBackendUserQueryParam) ([]*AdminBackendUser, int64) {
	query := orm.NewOrm().QueryTable(AdminBackendUserTBName())
	data := make([]*AdminBackendUser, 0)
	//默认排序
	sortorder := "Id"
	switch params.Sort {
	case "Id":
		sortorder = "Id"
	}
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}
	if params.UserName != "" {
		query = query.Filter("username__icontains", params.UserName)
	}
	if params.RealName != "" {
		query = query.Filter("realname__icontains", params.RealName)
	}
	if len(params.Mobile) > 0 {
		query = query.Filter("mobile__icontains", params.Mobile)
	}
	if len(params.Status) > 0 {
		query = query.Filter("status__iexact", params.Status)
	}
	total, _ := query.Count()
	query = query.OrderBy(sortorder).Limit(params.Limit, (params.Offset-1)*params.Limit)
	//query = query.RelatedSel("AdminRoleBackendUserRel")
	query.All(&data)
	//AdminRolesByIds
	roleIds := make([]int, 0)
	roleIdsMap := make(map[int]struct{})
	for _, obj := range data {
		list := strings.Split(obj.RolesStr, ",")
		for _, s := range list {
			if i, err := strconv.Atoi(s); err == nil {
				obj.RoleIds = append(obj.RoleIds, i)
				if _, found := roleIdsMap[i]; !found {
					roleIdsMap[i] = struct{}{}
					roleIds = append(roleIds, i)
				}
			}
		}
	}
	//角色
	roleObjs := AdminRolesByIds(roleIds)
	for _, obj := range data {
		//拥有角色查询
		for _, rid := range obj.RoleIds {
			if role, found := roleObjs[rid]; found {
				obj.RolesName = append(obj.RolesName, role.Name)
			}
		}
		//密码不返回
		obj.UserPwd = ""
	}
	return data, total
}

// AdminBackendUserOne 根据id获取单条
func AdminBackendUserOne(id int) (*AdminBackendUser, error) {
	o := orm.NewOrm()
	m := AdminBackendUser{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// AdminBackendUserOneByName 根据 用户名获取单条
func AdminBackendUserOneByName(username string) (*AdminBackendUser, error) {
	m := AdminBackendUser{}
	if err := orm.NewOrm().QueryTable(AdminBackendUserTBName()).Filter("userName", username).One(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

//AdminBackendUserOneMobile 根据 用户电话获取单条
func AdminBackendUserOneMobile(mobile string) (*AdminBackendUser, error) {
	m := AdminBackendUser{}
	if err := orm.NewOrm().QueryTable(AdminBackendUserTBName()).Filter("mobile", mobile).One(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

// AdminBackendUserOneByUserName 根据用户名密码获取单条
func AdminBackendUserOneByUserName(username, userpwd string) (*AdminBackendUser, error) {
	m := AdminBackendUser{}
	err := orm.NewOrm().QueryTable(AdminBackendUserTBName()).Filter("userName", username).Filter("userPwd", userpwd).One(&m)
	//err := orm.NewOrm().QueryTable(AdminBackendUserTBName()).Filter("userName", username).One(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
