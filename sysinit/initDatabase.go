package sysinit

import (
	"fmt"
	"tokensky_bg_admin/conf"
	_ "tokensky_bg_admin/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	//_ "github.com/mattn/go-sqlite3"
	_ "github.com/go-sql-driver/mysql"
)

//初始化数据连接
func InitDatabase() {

	/*admin数据库*/
	dbAdminHost := beego.AppConfig.String(conf.DB_MYSQL_ADMIN_CONF + "::db_host")
	dbAdminPort := beego.AppConfig.String(conf.DB_MYSQL_ADMIN_CONF + "::db_port")
	dbAdminUser := beego.AppConfig.String(conf.DB_MYSQL_ADMIN_CONF + "::db_user")
	dbAdminPwd := beego.AppConfig.String(conf.DB_MYSQL_ADMIN_CONF + "::db_pwd")
	dbAdminName := beego.AppConfig.String(conf.DB_MYSQL_ADMIN_CONF + "::db_name")
	dbAdminCharset := beego.AppConfig.String(conf.DB_MYSQL_ADMIN_CONF + "::db_charset")
	dbAdminMaxActive, _ := beego.AppConfig.Int(conf.DB_MYSQL_ADMIN_CONF + "::db_max_active")

	fmt.Println(dbAdminHost, dbAdminPort, dbAdminUser, dbAdminPwd, dbAdminName, dbAdminCharset, dbAdminMaxActive)
	//注册ADMIN数据库
	if err := dbAddMysql(conf.DB_MYSQL_TYPE, dbAdminHost, dbAdminPort, dbAdminUser, dbAdminPwd, dbAdminName, conf.DB_MYSQL_ADMIN_ALIAS, dbAdminCharset, dbAdminMaxActive); err != nil {
		//数据库未正常连接
		panic("数据库连接失败Err:"+err.Error())
	}

	//如果是开发模式，则显示命令信息
	isDev := (beego.AppConfig.String("runmode") == "dev")
	//自动建表
	//orm.RunSyncdb("default", false, isDev)
	if isDev {
		orm.Debug = isDev
	}
}

//注册Mysql
func dbAddMysql(dbType, host, port, user, pwd, name, alias, dbCharset string, maxActive int) error {
	err := orm.RegisterDataBase(alias, dbType, user+":"+pwd+"@tcp("+host+":"+
		port+")/"+name+"?charset="+dbCharset+"&loc=Local", maxActive)
	return err
}
