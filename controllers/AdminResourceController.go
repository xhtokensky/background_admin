package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"tokensky_bg_admin/enums"
	"tokensky_bg_admin/models"

	"github.com/astaxie/beego/orm"
)

type AdminResourceController struct {
	BaseController
}

func (c *AdminResourceController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的少数Action需要权限控制，则将验证放到需要控制的Action里
	c.checkAuthor("TreeGrid", "UserMenuTree", "ParentTreeGrid", "Select")
}

func (c *AdminResourceController) Index() {
	//需要权限控制
	c.checkAuthor()
	//将页面左边菜单的某项激活
	c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "resource/index_headcssjs.html"
	c.LayoutSections["footerjs"] = "resource/index_footerjs.html"
	//页面里按钮权限控制
	c.Data["canEdit"] = c.checkActionAuthor("AdminResourceController", "Edit")
	c.Data["canDelete"] = c.checkActionAuthor("AdminResourceController", "Delete")
}

//TreeGrid 获取所有资源的列表
func (c *AdminResourceController) TreeGrid() {
	//tree := models.AdminResourceTreeGrid()
	tree := models.AdminResourceTreeGrid2()
	//转换UrlFor 2 LinkUrl
	//c.UrlFor2Link(tree)
	for _, obj := range tree {
		if obj.UrlFor != "" {
			obj.LinkUrl = c.UrlFor2LinkOne(obj.UrlFor)
		}
	}
	tree = models.AdminRolesChildSort(tree)
	c.jsonResult(enums.JRCodeSucc, "", tree)
}

//UserMenuTree 获取用户有权管理的菜单、区域列表
func (c *AdminResourceController) UserMenuTree() {
	//userid := c.curUser.Id
	////获取用户有权管理的菜单列表（包括区域）
	//tree := models.AdminResourceTreeGridByUserId(userid, 1)
	////转换UrlFor 2 LinkUrl
	//c.UrlFor2Link(tree)
	//c.jsonResult(enums.JRCodeSucc, "", tree)

	//父子关系版
	data := make(map[string]interface{})
	//测试 权限相关
	c.curUser.IsSuper = 1

	tree := models.AdminResourceTreeGridByUser(&c.curUser, 10000)
	for _, obj := range tree {
		if obj.UrlFor != "" {
			obj.LinkUrl = c.UrlFor2LinkOne(obj.UrlFor)
		}
	}
	tree = models.AdminRolesChildSort(tree)
	data["tree"] = tree
	c.jsonResult(enums.JRCodeSucc, "成功", data)
}

//ParentTreeGrid 获取可以成为某节点的父节点列表
func (c *AdminResourceController) ParentTreeGrid() {
	Id, _ := c.GetInt("id", 0)
	tree := models.AdminResourceTreeGrid4Parent(Id)
	//转换UrlFor 2 LinkUrl
	c.UrlFor2Link(tree)
	c.jsonResult(enums.JRCodeSucc, "", tree)
}

// UrlFor2LinkOne 使用URLFor方法，将资源表里的UrlFor值转成LinkUrl
func (c *AdminResourceController) UrlFor2LinkOne(urlfor string) string {
	if len(urlfor) == 0 {
		return ""
	}
	// AdminResourceController.Edit,:id,1
	strs := strings.Split(urlfor, ",")
	if len(strs) == 1 {
		return c.URLFor(strs[0])
	} else if len(strs) > 1 {
		var values []interface{}
		for _, val := range strs[1:] {
			values = append(values, val)
		}
		return c.URLFor(strs[0], values...)
	}
	return ""
}

//UrlFor2Link 使用URLFor方法，批量将资源表里的UrlFor值转成LinkUrl
func (c *AdminResourceController) UrlFor2Link(src []*models.AdminResource) {
	for _, item := range src {
		item.LinkUrl = c.UrlFor2LinkOne(item.UrlFor)
	}
}

//Edit 资源编辑页面
func (c *AdminResourceController) Edit() {
	//需要权限控制
	c.checkAuthor()
	//如果是POST请求，则由Save处理
	if c.Ctx.Request.Method == "POST" {
		c.Save()
	}
}

//Save 资源添加编辑 保存
func (c *AdminResourceController) Save() {
	var err error
	o := orm.NewOrm()
	parent := &models.AdminResource{}
	m := models.AdminResource{}
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m.Id)
	}
	mapp := make(map[string]int)
	json.Unmarshal(c.Ctx.Input.RequestBody, &mapp)
	parentId := mapp["parentId"]
	//获取父节点
	if parentId > 0 {
		parent, err = models.AdminResourceOne(parentId)
		if err == nil && parent != nil {
			m.Parent = parent
		} else {
			c.jsonResult(enums.JRCodeFailed, "父节点无效", "")
		}
	}
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

// Delete 删除
func (c *AdminResourceController) Delete() {
	//需要权限控制
	c.checkAuthor()
	body := c.Ctx.Input.RequestBody
	mapp := make(map[string]int)
	json.Unmarshal(body, &mapp)
	Id := mapp["id"]
	if Id == 0 {
		c.jsonResult(enums.JRCodeFailed, "选择的数据无效", 0)
	}
	query := orm.NewOrm().QueryTable(models.AdminResourceTBName())
	if _, err := query.Filter("id", Id).Delete(); err == nil {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("删除成功"), 0)
	} else {
		c.jsonResult(enums.JRCodeFailed, "删除失败", 0)
	}
}

// Select 通用选择面板
func (c *AdminResourceController) Select() {
	//获取调用者的类别 1表示 角色
	desttype, _ := c.GetInt("desttype", 0)
	//获取调用者的值
	destval, _ := c.GetInt("destval", 0)
	//返回的资源列表
	var selectedIds []string
	o := orm.NewOrm()
	if desttype > 0 && destval > 0 {
		//如果都大于0,则获取已选择的值，例如：角色，就是获取某个角色已关联的资源列表
		switch desttype {
		case 1:
			{
				role := models.AdminRole{Id: destval}
				o.LoadRelated(&role, "AdminRoleResourceRel")
				for _, item := range role.RoleResourceRel {
					selectedIds = append(selectedIds, strconv.Itoa(item.AdminResource.Id))
				}
			}
		}
	}
	c.Data["selectedIds"] = strings.Join(selectedIds, ",")
	c.setTpl("resource/select.html", "shared/layout_pullbox.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "resource/select_headcssjs.html"
	c.LayoutSections["footerjs"] = "resource/select_footerjs.html"
}

//CheckUrlFor 填写UrlFor时进行验证
func (c *AdminResourceController) CheckUrlFor() {
	body := c.Ctx.Input.RequestBody
	mapp := make(map[string]string)
	json.Unmarshal(body, &mapp)

	link := c.UrlFor2LinkOne(mapp["urlfor"])
	if len(link) > 0 {
		c.jsonResult(enums.JRCodeSucc, "解析成功", link)
	} else {
		c.jsonResult(enums.JRCodeFailed, "解析失败", link)
	}
}
func (c *AdminResourceController) UpdateSeq() {
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
	oM, err := models.AdminResourceOne(Id)
	if err != nil || oM == nil {
		c.jsonResult(enums.JRCodeFailed, "选择的数据无效", 0)
	}
	oM.Seq = value
	if _, err := orm.NewOrm().Update(oM); err == nil {
		c.jsonResult(enums.JRCodeSucc, "修改成功", oM.Id)
	} else {
		c.jsonResult(enums.JRCodeFailed, "修改失败", oM.Id)
	}
}

//获取用户权限列表
func (c *AdminResourceController) GetUserResources() {
	//用户角色权限菜单
	tree := models.AdminResourceTreeGridByUserId(c.curUser.Id, 1)
	//url处理
	c.UrlFor2Link(tree)
	//角色url处理
	data := make(map[string]interface{})
	data["tree"] = models.AdminRolesChildSort(tree)
	c.jsonResult(enums.JRCodeSucc, "获取数据成功", data)
}

//获取角色已拥有的权限
func (c *AdminResourceController) TreeGridByRole() {
	body := c.Ctx.Input.RequestBody
	mapp := make(map[string]int)
	json.Unmarshal(body, &mapp)
	Id := mapp["id"]
	o := orm.NewOrm()
	list := make([]*models.AdminRoleResourceRel, 0)
	o.QueryTable(models.AdminRoleResourceRelTBName()).Filter("AdminRole__id__exact", Id).All(&list)
	ids := make([]int, 0)
	for _, obj := range list {
		ids = append(ids, obj.AdminResource.Id)
	}
	data := make([]*models.AdminResource, 0)
	if len(ids) > 0 {
		o.QueryTable(models.AdminResourceTBName()).Filter("id__in", ids).All(&data)
	}
	c.jsonResult(enums.JRCodeSucc, "", data)

	//ids := make(map[int]struct{})
	//for _,obj := range list{
	//	ids[obj.AdminResource.Id] = struct{}{}
	//}
	//tree := models.AdminResourceTreeGrid2()
	//for _,obj := range tree{
	//	if obj.UrlFor != ""{
	//		obj.LinkUrl = c.UrlFor2LinkOne(obj.UrlFor)
	//	}
	//	//
	//	if _,found := ids[obj.Id];found{
	//		obj.Choice = true
	//	}
	//}
	//tree = models.AdminRolesChildSort(tree)
	//c.jsonResult(enums.JRCodeSucc, "", tree)
}
