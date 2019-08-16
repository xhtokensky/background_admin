package models

import (
	"github.com/astaxie/beego/orm"
	"time"
	"tokensky_bg_admin/utils"
)

//查询的类
type HashrateOrderQueryParam struct {
	BaseQueryParam
	Phone     string `json:"phone"`     //手机号
	StartTime int64  `json:"startTime"` //开始时间
	EndTime   int64  `json:"endTime"`   //截止时间
	Status    string `json:"status"`    //状态
	OrderId   string `json:"orderId"`   //订单号
}

func (a *HashrateOrder) TableName() string {
	return HashrateOrderTBName()
}

//算力订单表
type HashrateOrder struct {
	OrderId string `orm:"pk;column(order_id)"json:"orderId"form:"orderId"`
	//购买数量
	BuyQuantity int `orm:"column(buy_quantity)"json:"buyQuantity"form:"buyQuantity"`
	//支付时间 付款时间
	PayTime time.Time `orm:"type(datetime);column(pay_time)"json:"payTime"form:"payTime"`
	//创建时间
	CreateTime time.Time `orm:"auto_now_add;type(datetime);column(create_time)"json:"createTime"form:"createTime"`
	//开挖时间
	ExcavateTime time.Time `orm:"type(datetime);column(excavate_time)"json:"excavateTime"` //开挖时间
	//结束时间
	EndTime time.Time `orm:"type(date);column(end_time)"json:"endTime"` //挖矿截止时间
	//状态 1交易中 2结算中 3已完成
	Status int `orm:"column(status)"json:"status"`
	//用户表
	User     *TokenskyUser `orm:"rel(fk);column(user_id)"json:"-"form:"-"`
	UserId   int           `orm:"-"json:"userId"form:"userId"` //用户ID
	Phone    string        `orm:"-"form:"-" json:"phone"`      //手机号
	NickName string        `orm:"-"json:"nickName"form:"-"`    //昵称
	//算力合约表
	HashrateTreaty   *HashrateTreaty `orm:"rel(fk);column(hashrate_treaty_id)"json:"-"form:"-"`
	HashrateTreatyId int             `orm:"-"json:"hashrateTreatyId"form:"hashrateTreatyId"` //算力合约ID
	Titie            string          `orm:"-"form:"-" json:"titie"`
	Price            float64         `orm:"-"json:"price"`                    //单价
	MiningCycle      int             `orm:"-"form:"-"json:"miningCycle"`      //挖矿周期(运行周期)
	DayElectricBill  float64         `orm:"-"form:"-"json:"dayElectricBills"` //电费

	//算力合约分类表
	ImgUrl string `orm:"-"form:"-"json:"imgUrl"` //图片路径
	//算力合约关联表
	PayType          string  `orm:"-"json:"payType"`          //支付方式
	PayTypeList      []int   `orm:"-"json:"payTypeList"`      //支付方式
	TransactionMoney float64 `orm:"-"json:"transactionMoney"` //交易金额
	/*未关联表*/
	Income float64 `orm:"-"json:"income"` //累计收益(计算的)
	// 1交易中 2结算中 3已完成 后台计算
	TransactionStatus int `orm:"-"json:"transactionStatus"` //交易状态(计算的)

}

//获取分页数据
func HashrateOrderPageList(params *HashrateOrderQueryParam) ([]*HashrateOrder, int64) {
	o := orm.NewOrm()
	query := o.QueryTable(HashrateOrderTBName())
	data := make([]*HashrateOrder, 0)
	//默认排序
	sortorder := "create_time"
	switch params.Sort {
	case "createTime":
		sortorder = "create_time"
	}
	switch params.Order {
	case "":
		sortorder = "-" + sortorder
	case "desc":
		sortorder = "-" + sortorder
	default:
		sortorder = sortorder
	}

	//时间段
	if params.StartTime > 0 {
		query = query.Filter("create_time__gte", time.Unix(params.StartTime, 0))
	}

	if params.EndTime > 0 {
		query = query.Filter("create_time__lte", time.Unix(params.EndTime, 0))
	}

	//订单号查询
	if params.OrderId != "" {
		query = query.Filter("order_id__contains", params.OrderId)

	}
	//手机查询
	if params.Phone != "" {
		query = query.Filter("User__Phone__icontains", params.Phone)
	}

	total, _ := query.Count()
	query = query.OrderBy(sortorder).Limit(params.Limit, (params.Offset-1)*params.Limit)
	query.RelatedSel().All(&data)

	//算力合约关联表
	hashrateOrderTransactionIds := make([]string, 0, len(data))
	for _, obj := range data {
		hashrateOrderTransactionIds = append(hashrateOrderTransactionIds, obj.OrderId)
	}
	hashrateOrderTransactions := HashrateOrderTransactionsByIds(hashrateOrderTransactionIds)
	//资源
	deadline := time.Now().Add(time.Second * 3600).Unix()
	//当前时间
	now := time.Now().Unix()
	for _, obj := range data {
		//用户表
		if obj.User != nil {
			obj.Phone = obj.User.Phone
			obj.UserId = obj.User.UserId
			obj.NickName = obj.User.NickName
		}
		//算力合约表
		if obj.HashrateTreaty != nil {
			obj.HashrateTreatyId = obj.HashrateTreaty.KeyId
			obj.Titie = obj.HashrateTreaty.Title
			obj.Price = obj.HashrateTreaty.Price
			obj.MiningCycle = obj.HashrateTreaty.RunCycle
			//电费
			obj.DayElectricBill = obj.HashrateTreaty.ElectricBill * float64(obj.BuyQuantity)
			//算力合约分类表
			if obj.HashrateTreaty.HashrateCategoryObj != nil {
				obj.ImgUrl = utils.QiNiuDownload(obj.HashrateTreaty.HashrateCategoryObj.ImgKey, deadline)
			}
			//算力合约关联表[暂时单条处理]
			//obj.PayTypeList = make([]int, conf.PAY_TYPE_MAX_NUM)
			if hobjs, found := hashrateOrderTransactions[obj.OrderId]; found {
				for _, hobj := range hobjs {
					obj.TransactionMoney = hobj.TransactionMoney
					//支付方式
					obj.PayType = hobj.PayType
					////支付方式处理
					//for _, str := range strings.Split(hobj.PayType, ",") {
					//
					//	if con, err := strconv.Atoi(str); err == nil {
					//		if con > 0 && con <= conf.PAY_TYPE_MAX_NUM {
					//			obj.PayTypeList[con-1] = 1
					//		}
					//	}
					//}
					break
				}

			}
			//交易状态计算
			if now >= obj.EndTime.Unix() {
				obj.Status = 2
			} else if obj.ExcavateTime.Unix() > now {
				obj.Status = 1
			}
		}
	}

	return data, total
}

func HashrateOrderIteration(num int, st time.Time) func() []*HashrateOrder {
	query := orm.NewOrm().QueryTable(HashrateOrderTBName())
	tn := time.Date(st.Year(), st.Month(), st.Day(), 0, 0, 0, 0, st.Location())
	query = query.Filter("excavate_time__lte", tn)
	query = query.Filter("end_time__gte", tn)
	page := 1
	count,_ := query.Count()
	return func() []*HashrateOrder {
		data := make([]*HashrateOrder, 0)
		if count>0{
			query.Limit(num, (page-1)*num).RelatedSel("hashrate_treaty_id").All(&data)
			page++
			count-=int64(num)
		}
		return data
	}
}

//算力订单表改非期货
func HashrateOrderIsNotFutures(o orm.Ormer, obj *HashrateTreaty) (int64, error) {
	query := o.QueryTable(HashrateOrderTBName())
	query.Filter("hashrate_treaty_id__exact", obj.KeyId)
	endTime2 := time.Unix(obj.EffectiveDate.Unix(), 0)
	endTime := time.Date(endTime2.Year(), endTime2.Month(), endTime2.Day()+obj.RunCycle, 0, 0, 0, 0, time.Local)
	params := make(map[string]interface{})
	params["excavate_time"] = obj.EffectiveDate
	params["end_time"] = endTime
	return query.Update(params)
}
