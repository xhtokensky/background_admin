package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"strings"
	"tokensky_bg_admin/enums"
	"tokensky_bg_admin/models"
	"tokensky_bg_admin/utils"
)

type BaseController struct {
	beego.Controller
	controllerName string                  //当前控制名称
	actionName     string                  //当前action名称
	curUser        models.AdminBackendUser //当前用户信息
}

func (c *BaseController) Prepare() {
	//附值
	c.controllerName, c.actionName = c.GetControllerAndAction()
	//从Session里获取数据 设置用户信息
	c.adapterUserInfo()
	//测试 权限
	c.curUser.Id = 22
	//日志
	switch c.Ctx.Request.Method {
	//现在只使用了 GET POST请求
	case "POST":
		if c.controllerName == "AdminHomeController" {
			if c.actionName == "DoLogin" {
				return
			}
		}

		//登陆用户
		logData := map[string]interface{}{
			"controllerName": c.controllerName,
			"actionName":     c.actionName,
			"uid":            c.curUser.Id,
			"ip":             c.Ctx.Request.Host,
			"ok":             0,
		}
		logData["body"] = string(c.Ctx.Input.RequestBody)
		//日志记录
		logStr, _ := json.Marshal(logData)
		utils.LogInfo(logStr)
	}
}

// checkLogin判断用户是否登录，未登录则跳转至登录页面
// 一定要在BaseController.Prepare()后执行
func (c *BaseController) checkLogin() {
	//if c.curUser.Id == 0 {
	//	c.jsonResult(enums.JRCodeNotLogin, "请登陆", 0)
	//}
}

// 判断某 Controller.Action 当前用户是否有权访问
func (c *BaseController) checkActionAuthor(ctrlName, ActName string) bool {
	//测试 权限
	return true


	//if c.curUser.Id == 0 {
	//	return false
	//}
	////从session获取用户信息
	//user := c.GetSession("backenduser")
	////类型断言
	//v, ok := user.(models.AdminBackendUser)
	//if ok {
	//	//如果是超级管理员，则直接通过
	//	if v.IsSuper == 1 {
	//		return true
	//	}
	//	//遍历用户所负责的资源列表
	//	for _, str := range v.ResourceUrlForList {
	//		urlfor := strings.TrimSpace(str)
	//		if len(urlfor) == 0 {
	//			continue
	//		}
	//		// TestController.Get,:last,xie,:first,asta
	//		strs := strings.Split(urlfor, ",")
	//		if len(strs) > 0 && strs[0] == (ctrlName+"."+ActName) {
	//			return true
	//		}
	//	}
	//}
	//return false
}

// checkLogin判断用户是否有权访问某地址，无权则会跳转到错误页面
//一定要在BaseController.Prepare()后执行
// 会调用checkLogin
// 传入的参数为忽略权限控制的Action
func (c *BaseController) checkAuthor(ignores ...string) {
	//先判断是否登录
	c.checkLogin()
	//如果Action在忽略列表里，则直接通用
	for _, ignore := range ignores {
		if ignore == c.actionName {
			return
		}
	}
	hasAuthor := c.checkActionAuthor(c.controllerName, c.actionName)
	if !hasAuthor {
		//如果没有权限
		c.jsonResult(enums.JRCode401, "无权访问", 0)
	}
}

//从session里取用户信息
func (c *BaseController) adapterUserInfo() {
	a := c.GetSession("backenduser")
	if a != nil {
		c.curUser = a.(models.AdminBackendUser)
		c.Data["backenduser"] = a
	}
}

//SetBackendUser2Session 获取用户信息（包括资源UrlFor）保存至Session
func (c *BaseController) setBackendUser2Session(userId int) error {
	m, err := models.AdminBackendUserOne(userId)
	if err != nil {
		return err
	}
	//获取这个用户能获取到的所有资源列表
	resourceList := models.AdminResourceTreeGridByUserId(userId, 1000)
	for _, item := range resourceList {
		m.ResourceUrlForList = append(m.ResourceUrlForList, strings.TrimSpace(item.UrlFor))
	}
	c.SetSession("backenduser", *m)
	return nil
}

// 设置模板
// 第一个参数模板，第二个参数为layout
func (c *BaseController) setTpl(template ...string) {
	var tplName string
	layout := "shared/layout_page.html"
	switch {
	case len(template) == 1:
		tplName = template[0]
	case len(template) == 2:
		tplName = template[0]
		layout = template[1]
	default:
		//不要Controller这个10个字母
		ctrlName := strings.ToLower(c.controllerName[0 : len(c.controllerName)-10])
		actionName := strings.ToLower(c.actionName)
		tplName = ctrlName + "/" + actionName + ".html"
	}
	c.Layout = layout
	c.TplName = tplName
}

func (c *BaseController) jsonResult(code enums.JsonResultCode, msg string, content interface{}) {
	r := &models.JsonResult{code, msg, content}
	c.Data["json"] = r
	//记录日志
	logData := map[string]interface{}{
		"controllerName": c.controllerName,
		"actionName":     c.actionName,
		"uid":            c.curUser.Id,
		"ip":             c.Ctx.Request.Host,
		"ok":             1,
	}
	switch code {
	case enums.JRCodeSucc:
		//正常
		logStr, _ := json.Marshal(logData)
		utils.LogDebug(logStr)
	case enums.JRCodeFailed:
		//错误操作
		if c.Ctx.Request.Method == "GET" {
			logData["form"] = c.Ctx.Request.Form
		} else {
			logData["postForm"] = c.Ctx.Request.PostForm
			logData["body"] = string(c.Ctx.Input.RequestBody)
		}
		logStr, _ := json.Marshal(logData)
		utils.LogDebug(logStr)
	case enums.JRCode401:
		//未经授权
		logStr, _ := json.Marshal(logData)
		utils.LogDebug(logStr)
	default:
		//其它状态
		logStr, _ := json.Marshal(logData)
		utils.LogDebug(logStr)
	}

	c.ServeJSON()
	c.StopRun()
}

func (c *BaseController) jsonResult2(code enums.JsonResultCode, msg string) {
	r := &models.JsonResult2{code, msg}
	c.Data["json"] = r
	//记录日志
	logData := map[string]interface{}{
		"controllerName": c.controllerName,
		"actionName":     c.actionName,
		"uid":            c.curUser.Id,
		"ip":             c.Ctx.Request.Host,
		"ok":             1,
	}
	switch code {
	case enums.JRCodeSucc:
		//正常
		logStr, _ := json.Marshal(logData)
		utils.LogDebug(logStr)
	case enums.JRCodeFailed:
		//错误操作
		if c.Ctx.Request.Method == "GET" {
			logData["form"] = c.Ctx.Request.Form
		} else {
			logData["postForm"] = c.Ctx.Request.PostForm
			logData["body"] = string(c.Ctx.Input.RequestBody)
		}
		logStr, _ := json.Marshal(logData)
		utils.LogDebug(logStr)
	case enums.JRCode401:
		//未经授权
		logStr, _ := json.Marshal(logData)
		utils.LogDebug(logStr)
	default:
		//其它状态
		logStr, _ := json.Marshal(logData)
		utils.LogDebug(logStr)
	}
	c.ServeJSON()
	c.StopRun()
}

func (c *BaseController) jsonResultError2(code enums.JsonResultCode, msg string) {
	r := &models.JsonResult2{code, msg}
	c.Data["json"] = r
	//记录日志记录所有信息
	logData := map[string]interface{}{
		"controllerName": c.controllerName,
		"actionName":     c.actionName,
		"uid":            c.curUser.Id,
		"ok":             1,
	}
	//请求类型
	if c.Ctx.Request.Method == "GET" {
		logData["form"] = c.Ctx.Request.Form
	} else {
		logData["postForm"] = c.Ctx.Request.PostForm
		logData["body"] = string(c.Ctx.Input.RequestBody)
	}
	//相应数据
	logData["data"] = c.Data
	//日志记录
	utils.LogError(logData)

	c.ServeJSON()
	c.StopRun()
}

//记录所有的异常信息
func (c *BaseController) jsonResultError(code enums.JsonResultCode, msg string, content interface{}) {
	r := &models.JsonResult{code, msg, content}
	c.Data["json"] = r
	//记录日志记录所有信息
	logData := map[string]interface{}{
		"controllerName": c.controllerName,
		"actionName":     c.actionName,
		"uid":            c.curUser.Id,
		"ok":             1,
	}
	//请求类型
	if c.Ctx.Request.Method == "GET" {
		logData["form"] = c.Ctx.Request.Form
	} else {
		logData["postForm"] = c.Ctx.Request.PostForm
		logData["body"] = string(c.Ctx.Input.RequestBody)
	}
	//相应数据
	logData["data"] = c.Data
	//日志记录
	utils.LogError(logData)

	c.ServeJSON()
	c.StopRun()
}

// 重定向
func (c *BaseController) redirect(url string) {
	c.Redirect(url, 302)
	c.StopRun()
}

// 重定向 去错误页
func (c *BaseController) pageError(msg string) {
	errorurl := c.URLFor("AdminHomeController.Error") + "/" + msg
	c.Redirect(errorurl, 302)
	c.StopRun()
}

//// 重定向 去登录页
//func (c *BaseController) pageLogin() {
//	url := c.URLFor("AdminHomeController.Login")
//	c.Redirect(url, 302)
//	c.StopRun()
//}

func (c *BaseController) Options() {
	c.Data["json"] = map[string]interface{}{"status": 200, "message": "ok", "moreinfo": ""}
	c.ServeJSON()
}
