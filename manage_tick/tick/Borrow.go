package tick

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
	"tokensky_bg_admin/conf"
	"tokensky_bg_admin/models"
	"tokensky_bg_admin/utils"
)

var (
	//警告时间间隔
	borrow_warn_time int64 = 86400
	//警告系数
	borrow_warn_ratio float64 = 0.1
)

func init()  {
	if con,err := beego.AppConfig.Int64("borrow" + "::warn_time");err ==nil{
		borrow_warn_time = con
	}
	if con,err := beego.AppConfig.Float("borrow" + "::warn_ratio");err == nil{
		borrow_warn_ratio = con
	}
}

//借贷风险控制系统
var borrowWindControlSystemSign bool = true
func BorrowWindControlSystem()error{
	if borrowWindControlSystemSign{
		borrowWindControlSystemSign = false
		defer func() {borrowWindControlSystemSign=true}()
		errObj := make(map[string][]int,0)
		iteration := models.BorrowOrderIteration()
		now := time.Now()
		nowInt := now.Unix()
		o := orm.NewOrm()
		var err error
		for {
			objs := iteration()
			if len(objs)<=0{
				break
			}
			for _,obj := range objs{
				//逾期校验
				if nowInt >obj.ExpireTime.Unix(){
					err = models.BorrowOrderExpireLimiting(o,obj)
					if err != nil{
						if _,found := errObj["逾期处理异常"];!found{
							errObj["逾期处理异常"] = make([]int,0)
						}
						errObj["逾期处理异常"] = append(errObj["逾期处理异常"], obj.Id)
					}
				}
				//质押校验
				if obj.GetRealTimePledge(){
					if obj.RealTimePledge >= 0.9{
						err = models.BorrowOrderPledgeLimiting(o,obj)
						if err != nil{
							if _,found := errObj["质押异常"];!found{
								errObj["质押异常"] = make([]int,0)
							}
							errObj["质押异常"] = append(errObj["质押异常"], obj.Id)
						}
						continue
					}
					//质押警告
					if obj.RealTimePledge > 0.9-borrow_warn_ratio{
						if obj.ExpireTime.Unix() + borrow_warn_time < nowInt{
							models.BorrowOrderWarn(o,obj)
						}
					}
				}else {
					if _,found := errObj["风控计算异常"];!found{
						errObj["风控计算异常"] = make([]int,0)
					}
					errObj["风控计算异常"] = append(errObj["风控计算异常"], obj.Id)
				}
			}
		}
		if len(errObj)>0{
			//异常警告
			utils.EmailNotify(conf.EMAIL_ERROR_LEVEL_ERROR,"借贷风控","风控异常",errObj)
		}
	}
	return nil
}
