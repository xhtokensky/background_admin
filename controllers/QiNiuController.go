package controllers

import (
	"fmt"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"strconv"
	"time"
	"tokensky_bg_admin/conf"
	"tokensky_bg_admin/enums"
	"tokensky_bg_admin/utils"
)

//七牛
type QiNiuController struct {
	BaseController
}

//Prepare 参考beego官方文档说明
func (c *QiNiuController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	////如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	//c.checkAuthor("DataGrid", "DataList", "UpdateSeq")
	////如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	////权限控制里会进行登录验证，因此这里不用再作登录验证
	c.checkLogin()
}

//获取密钥
func (c *QiNiuController) GetQiNiuKey() {

	accessKey := conf.QINIU_ACCESS_KEY
	secretKey := conf.QINIU_SERVERT_KEY
	mac := qbox.NewMac(accessKey, secretKey)
	// 需要覆盖的文件名
	keyToOverwrite := "qiniu1.mp4"
	putPolicy := storage.PutPolicy{
		Scope: fmt.Sprintf("%s:%s", conf.QINIU_BUCHENT_NAME, keyToOverwrite),
	}
	upToken := putPolicy.UploadToken(mac)
	//定义返回的数据结构
	mapp := map[string]interface{}{
		"token": upToken,
		"url":   conf.QINIU_SERVER,
	}
	c.jsonResult(enums.JRCodeSucc, "", mapp)
}

//七牛上传插件
func (c QiNiuController) UploadFile() {
	f, h, err := c.GetFile("file")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "上传失败", 0)
		return
	}
	defer f.Close()
	// 这里获得的实际就是一个io,通过源码看到这个open方法最终返回的是一个结构体,其内部包含了 io.Reader的接口
	file, _ := h.Open()
	// type File interface {
	// io.Reader
	// io.ReaderAt
	// io.Seeker
	// io.Closer
	// }
	//unix := time.Now().Unix()
	//timeByte := []byte(strconv.Itoa(int(unix)))
	now := time.Now().Unix()
	filename := "admin" + "/" + strconv.Itoa(c.curUser.Id) + "/" + strconv.FormatInt(now, 10) + "/" + h.Filename
	if filename, err = utils.Upload(file, h.Size, filename); err != nil { // 通过h.size 即可获得文件大小
		c.jsonResult(enums.JRCodeFailed, "上传失败", 0)
	} else {
		data := map[string]string{
			"fileName": h.Filename,                       //文件名称
			"fileSize": strconv.FormatInt(h.Size, 10),    //文件大小
			"imgUrl":   utils.QiNiuDownload(filename, 0), //图片路径
			"imgKey":   filename,                         //路径
		}
		c.jsonResult(enums.JRCodeSucc, "上传成功", data)
	}
}
