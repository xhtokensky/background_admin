package controllers

import (
	"encoding/json"
	"fmt"

	"tokensky_bg_admin/enums"
	"tokensky_bg_admin/models"

	"strconv"
	"strings"

	"github.com/astaxie/beego/orm"
)

//AdminRoleController 角色管理
type AdminRoleController struct {
	BaseController
}

//Prepare 参考beego官方文档说明
func (c *AdminRoleController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	c.checkAuthor()
}

//Index 角色管理首页
func (c *AdminRoleController) Index() {
	//是否显示更多查询条件的按钮
	c.Data["showMoreQuery"] = false
	//将页面左边菜单的某项激活
	c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "role/index_headcssjs.html"
	c.LayoutSections["footerjs"] = "role/index_footerjs.html"
	//页面里按钮权限控制
	c.Data["canEdit"] = c.checkActionAuthor("AdminRoleController", "Edit")
	c.Data["canDelete"] = c.checkActionAuthor("AdminRoleController", "Delete")
	c.Data["canAllocate"] = c.checkActionAuthor("AdminRoleController", "Allocate")
}

// DataGrid 角色管理首页 表格获取数据
func (c *AdminRoleController) DataGrid() {
	//直接反序化获取json格式的requestbody里的值
	var params models.AdminRoleQueryParam
	switch c.Ctx.Request.Method {
	case "POST":
		json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	case "GET":
		params.Sort = c.GetString("sort")
		params.Order = c.GetString("order")
		params.Limit, _ = c.GetInt64("limit")
		params.Offset, _ = c.GetInt64("offset")
	}
	if params.Sort == "" {
		params.Sort = "id"
	}
	//获取数据列表和总数
	data, total := models.AdminRolePageList(&params)
	mapp := make(map[string]interface{})
	mapp["rows"] = data
	mapp["total"] = total
	c.jsonResult(enums.JRCodeSucc, "", mapp)
}

//DataList 角色列表
func (c *AdminRoleController) DataList() {
	var params = models.AdminRoleQueryParam{}
	//获取数据列表和总数
	data := models.AdminRoleDataList(&params)
	//定义返回的数据结构
	c.jsonResult(enums.JRCodeSucc, "", data)
}

//Edit 添加、编辑角色界面
func (c *AdminRoleController) Edit() {
	if c.Ctx.Request.Method == "POST" {
		c.Save()
	}
	//Id, _ := c.GetInt(":id", 0)
	//m := models.AdminRole{Id: Id, Seq: 100}
	//if Id > 0 {
	//	o := orm.NewOrm()
	//	err := o.Read(&m)
	//	if err != nil {
	//		c.pageError("数据无效，请刷新后重试")
	//	}
	//}
	//c.Data["m"] = m
}

//Save 添加、编辑页面 保存
func (c *AdminRoleController) Save() {
	m := models.AdminRole{}
	var err error
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m.Id)
	}
	o := orm.NewOrm()
	if m.Id == 0 {
		if _, err = o.Insert(&m); err == nil {
			c.jsonResult(enums.JRCodeSucc, "添加成功", m.Id)
		} else {
			c.jsonResult(enums.JRCodeFailed, "添加失败", m.Id)
		}
	} else {
		if _, err = o.Update(&m); err == nil {
			c.jsonResult(enums.JRCodeSucc, "编辑成功", m.Id)
		} else {
			c.jsonResult(enums.JRCodeFailed, "编辑失败", m.Id)
		}
	}

}

//Delete 批量删除
func (c *AdminRoleController) Delete() {
	body := c.Ctx.Input.RequestBody
	mapp := make(map[string]string)
	json.Unmarshal(body, &mapp)
	ids := make([]int, 0)
	for _, str := range strings.Split(mapp["ids"], ",") {
		if id, err := strconv.Atoi(str); err == nil {
			ids = append(ids, id)
		}
	}
	if len(ids) > 0 {
		if num, err := models.AdminRoleBatchDelete(ids); err == nil {
			c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), 0)
		} else {
			c.jsonResult(enums.JRCodeFailed, "删除失败", 0)
		}
	} else {
		c.jsonResult(enums.JRCodeFailed, "没有等删除数据", 0)
	}
}

//Allocate 给角色分配资源界面
func (c *AdminRoleController) Allocate() {
	var roleId int
	body := c.Ctx.Input.RequestBody
	mapp := make(map[string]string)
	json.Unmarshal(body, &mapp)
	if con, err := strconv.Atoi(mapp["id"]); err == nil {
		roleId = con
	}
	o := orm.NewOrm()
	m := models.AdminRole{Id: roleId}
	if err := o.Read(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "数据无效，请刷新后重试", "")
	}
	//删除已关联的历史数据
	if _, err := o.QueryTable(models.AdminRoleResourceRelTBName()).Filter("adminRole__id", m.Id).Delete(); err != nil {
		c.jsonResult(enums.JRCodeFailed, "删除历史关系失败", "")
	}
	var relations []models.AdminRoleResourceRel
	for _, str := range strings.Split(mapp["ids"], ",") {
		if id, err := strconv.Atoi(str); err == nil {
			r := models.AdminResource{Id: id}
			relation := models.AdminRoleResourceRel{AdminRole: &m, AdminResource: &r}
			relations = append(relations, relation)
		}
	}
	if len(relations) > 0 {
		//批量添加
		if _, err := o.InsertMulti(len(relations), relations); err == nil {
			c.jsonResult(enums.JRCodeSucc, "保存成功", "")
		}
	}
	c.jsonResult(0, "保存失败", "")
}

//更新seq
func (c *AdminRoleController) UpdateSeq() {
	var Id, value int
	body := c.Ctx.Input.RequestBody
	mapp := make(map[string]string)
	json.Unmarshal(body, &mapp)
	if con, err := strconv.Atoi(mapp["pk"]); err == nil {
		Id = con
	}
	if con, err := strconv.Atoi(mapp["value"]); err == nil {
		value = con
	}
	oM, err := models.AdminRoleOne(Id)
	if err != nil || oM == nil {
		c.jsonResult(enums.JRCodeFailed, "选择的数据无效", 0)
	}
	oM.Seq = value
	o := orm.NewOrm()
	if _, err := o.Update(oM); err == nil {
		c.jsonResult(enums.JRCodeSucc, "修改成功", oM.Id)
	} else {
		c.jsonResult(enums.JRCodeFailed, "修改失败", oM.Id)
	}
}
