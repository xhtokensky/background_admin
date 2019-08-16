package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"strings"
	"tokensky_bg_admin/enums"
	"tokensky_bg_admin/models"
	"tokensky_bg_admin/utils"
)

type AdminHomeController struct {
	BaseController
}

func (c *AdminHomeController) Index() {
	//判断是否登录
	c.checkLogin()
	c.setTpl()

}

func (c *AdminHomeController) Page404() {
	c.setTpl()
}

func (c *AdminHomeController) Error() {
	c.Data["error"] = c.GetString(":error")
	c.setTpl("home/error.html", "shared/layout_pullbox.html")
}

func (c *AdminHomeController) Login() {
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "home/login_headcssjs.html"
	c.LayoutSections["footerjs"] = "home/login_footerjs.html"
	c.setTpl("home/login.html", "shared/layout_base.html")
}

//用户登陆
func (c *AdminHomeController) DoLogin() {
	//登陆安全校验

	var args models.AdminBackendUser
	body := c.Ctx.Input.RequestBody //接收raw body内容
	json.Unmarshal(body, &args)
	//POST请求
	username := args.UserName
	userpwd := args.UserPwd
	if len(username) == 0 || len(userpwd) == 0 {
		//校验
		//enums.UserDoLoginTimes.Load()

		c.jsonResult(enums.JRCodeFailed, "用户名和密码不正确", "")
	}
	userpwd = utils.String2md5(userpwd)
	user, err := models.AdminBackendUserOneByUserName(username, userpwd)
	if user != nil && err == nil {
		if user.Status == enums.Disabled {
			c.jsonResult(enums.JRCodeFailed, "用户被禁用，请联系管理员", "")
		}
		//保存用户信息到session
		c.setBackendUser2Session(user.Id)
		data := make(map[string]interface{})
		//返回用户信息
		obj := &models.AdminBackendUser{
			Id:       user.Id,
			RealName: user.RealName,
			UserName: user.UserName,
			Status:   user.Status,
			Email:    user.Email,
			Avatar:   user.Avatar,
			Mobile:   user.Mobile,
			IsSuper:  user.IsSuper,
		}
		data["user"] = obj
		c.jsonResult(enums.JRCodeSucc, "登录成功", data)
	} else {
		c.jsonResult(enums.JRCodeFailed, "用户名或者密码错误", "")
	}
	//
}

//退出登陆
func (c *AdminHomeController) Logout() {
	user := models.AdminBackendUser{}
	token := c.Ctx.Request.Header.Get("token")
	c.SetSession(token, user)
	//c.pageLogin()
	c.jsonResult(enums.JRCodeSucc, "退出成功", 0)
}

//用户注册
func (c *AdminHomeController) Register() {
	m := models.AdminBackendUser{}
	o := orm.NewOrm()
	var err error
	//获取form里的值
	if err = c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m.Id)
	}
	//数据初始化
	m.Id = 0
	m.IsSuper = 0
	m.Status = 1
	//创建用户(只能创建用户)
	if m.Id != 0 {
		m.Id = 0

	}
	//1 判断用户名称是否重复
	if obj, _ := models.AdminBackendUserOneByName(m.UserName); obj != nil {
		c.jsonResult(enums.JRCodeFailed, "用户名称重复", m.Id)
	}
	//2 判断电话是否已经注册
	if obj, _ := models.AdminBackendUserOneMobile(m.Mobile); obj != nil {
		c.jsonResult(enums.JRCodeFailed, "用户名称重复", m.Id)
	}
	//对密码进行加密
	m.UserPwd = utils.String2md5(m.UserPwd)
	if _, err := o.Insert(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "添加失败", m.Id)
	}
	//添加关系[暂未使用]
	var relations []models.AdminRoleBackendUserRel
	for _, roleId := range m.RoleIds {
		r := models.AdminRole{Id: roleId}
		relation := models.AdminRoleBackendUserRel{AdminBackendUser: &m, AdminRole: &r}
		relations = append(relations, relation)
	}
	if len(relations) > 0 {
		//批量添加
		if _, err := o.InsertMulti(len(relations), relations); err == nil {
			c.jsonResult(enums.JRCodeSucc, "保存成功", m.Id)
		} else {
			c.jsonResult(enums.JRCodeFailed, "保存失败", m.Id)
		}
	} else {
		c.jsonResult(enums.JRCodeSucc, "保存成功", m.Id)
	}
}

func (c *AdminHomeController) UrlFor2LinkOne(urlfor string) string {
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

//初始化数据库
func (c *AdminHomeController) DataReset() {
	//if ok, err := models.DataReset(); ok {
	//	c.jsonResult(enums.JRCodeSucc, "初始化成功", "")
	//} else {
	//	c.jsonResult(enums.JRCodeFailed, "初始化失败,可能原因:"+err.Error(), "")
	//}
	c.jsonResult(enums.JRCodeSucc, "接口已关闭", "")
}
