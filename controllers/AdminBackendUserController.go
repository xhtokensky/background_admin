package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"tokensky_bg_admin/enums"
	"tokensky_bg_admin/models"
	"tokensky_bg_admin/utils"
)

type AdminBackendUserController struct {
	BaseController
}

func (c *AdminBackendUserController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	c.checkAuthor()
	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//权限控制里会进行登录验证，因此这里不用再作登录验证
	//c.checkLogin()
}

func (c *AdminBackendUserController) Index() {
	//是否显示更多查询条件的按钮
	//c.Data["showMoreQuery"] = true
	//将页面左边菜单的某项激活
	//c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
	//页面模板设置
	//c.setTpl()
	//c.LayoutSections = make(map[string]string)
	//c.LayoutSections["headcssjs"] = "backenduser/index_headcssjs.html"
	//c.LayoutSections["footerjs"] = "backenduser/index_footerjs.html"
	//页面里按钮权限控制
	//c.Data["canEdit"] = c.checkActionAuthor("AdminBackendUserController", "Edit")
	//c.Data["canDelete"] = c.checkActionAuthor("AdminBackendUserController", "Delete")
}

//获取信息
func (c *AdminBackendUserController) DataGrid() {
	//直接反序化获取json格式的requestbody里的值（要求配置文件里 copyrequestbody=true)
	var params models.AdminBackendUserQueryParam
	switch c.Ctx.Request.Method {
	case "POST":
		json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	case "GET":
		params.Sort = c.GetString("sort")
		params.Order = c.GetString("order")
		params.Limit, _ = c.GetInt64("limit")
		params.Offset, _ = c.GetInt64("offset")
		params.Status = c.GetString("status")
		params.RealName = c.GetString("realName")
		params.UserName = c.GetString("userName")
		params.Mobile = c.GetString("mobile")
		params.Id, _ = c.GetInt("id")
	}
	mapp := make(map[string]interface{})
	//查单条
	if params.Id > 0 {
		obj, _ := models.AdminBackendUserOne(params.Id)
		if obj != nil {
			var roleIds []int
			o := orm.NewOrm()
			o.LoadRelated(obj, "AdminRoleBackendUserRel")
			for _, item := range obj.AdminRoleBackendUserRel {
				roleIds = append(roleIds, item.AdminRole.Id)
			}
			obj.RoleIds = roleIds
			mapp["rows"] = []*models.AdminBackendUser{obj}
			mapp["total"] = 1
			c.jsonResult(enums.JRCodeSucc, "", mapp)
		} else {
			c.jsonResult(enums.JRCodeFailed, "获取失败", mapp)
		}
	}
	//获取数据列表和总数
	data, total := models.AdminBackendUserPageList(&params)
	mapp["rows"] = data
	mapp["total"] = total
	c.jsonResult(enums.JRCodeSucc, "", mapp)
}

// Edit 添加 编辑 页面
func (c *AdminBackendUserController) Edit() {
	//如果是Post请求，则由Save处理
	if c.Ctx.Request.Method == "POST" {
		c.Save()
	}
}
func (c *AdminBackendUserController) Save() {
	m := models.AdminBackendUser{}
	o := orm.NewOrm()
	var err error
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m.Id)
	}
	//删除已关联的历史数据
	if _, err := o.QueryTable(models.AdminRoleBackendUserRelTBName()).Filter("adminBackendUser__id", m.Id).Delete(); err != nil {
		c.jsonResult(enums.JRCodeFailed, "删除历史关系失败", "")
	}
	//关系
	strs := make([]string, 0)
	for _, i := range strings.Split(m.Ids, ",") {
		if rid, err := strconv.Atoi(i); err == nil {
			strs = append(strs, i)
			m.RoleIds = append(m.RoleIds, rid)
		}
	}
	m.RolesStr = strings.Join(strs, ",")
	if m.Id == 0 {
		//对密码进行加密
		m.UserPwd = utils.String2md5(m.UserPwd)
		if uid, err := o.Insert(&m); err != nil {
			c.jsonResult(enums.JRCodeFailed, "添加失败", m.Id)
		} else {
			m.Id = int(uid)
		}
	} else {
		if oM, err := models.AdminBackendUserOne(m.Id); err != nil {
			c.jsonResult(enums.JRCodeFailed, "数据无效，请刷新后重试", m.Id)
		} else {
			m.UserPwd = strings.TrimSpace(m.UserPwd)
			if len(m.UserPwd) == 0 {
				//如果密码为空则不修改
				m.UserPwd = oM.UserPwd
			} else {
				m.UserPwd = utils.String2md5(m.UserPwd)
			}
			//本页面不修改头像和密码，直接将值附给新m
			m.Avatar = oM.Avatar
		}
		if _, err := o.Update(&m); err != nil {
			c.jsonResult(enums.JRCodeFailed, "编辑失败", m.Id)
		}
		//删除编辑用户的缓存session

	}
	//添加关系
	var relations []models.AdminRoleBackendUserRel
	for _, roleId := range m.RoleIds {
		r := models.AdminRole{Id: roleId}
		relation := models.AdminRoleBackendUserRel{AdminBackendUser: &m, AdminRole: &r}
		relations = append(relations, relation)
	}
	if len(relations) > 0 {
		//批量添加
		if _, err := o.InsertMulti(len(relations), relations); err == nil {
			//日志
			//models.AdminOperationLogsAdd(c.curUser.Id,)
			c.jsonResult(enums.JRCodeSucc, "保存成功", m.Id)
		} else {
			c.jsonResult(enums.JRCodeFailed, "保存失败", m.Id)
		}
	} else {
		c.jsonResult(enums.JRCodeSucc, "保存成功", m.Id)
	}
}

//删除
func (c *AdminBackendUserController) Delete() {
	c.jsonResult(enums.JRCodeFailed, "不允许删除", 0)

	body := c.Ctx.Input.RequestBody
	mapp := make(map[string]string)
	json.Unmarshal(body, &mapp)
	ids := make([]int, 0)
	for _, str := range strings.Split(mapp["ids"], ",") {
		if id, err := strconv.Atoi(str); err == nil {
			ids = append(ids, id)
		}
	}
	query := orm.NewOrm().QueryTable(models.AdminBackendUserTBName())
	if num, err := query.Filter("id__in", ids).Delete(); err == nil {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), 0)
	} else {
		c.jsonResult(enums.JRCodeFailed, "删除失败", 0)
	}
}
