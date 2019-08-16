package models

import (
	"github.com/astaxie/beego/orm"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"tokensky_bg_admin/conf"
)

//唯一订单号id
func (a *TokenskyOrderIds) TableName() string {
	return TokenskyOrderIdsTBName()
}

func init() {
	conf.TokenskyOrderIdsIterationGetOid = TokenskyOrderIdsIteration()
}

//算力订单表
type TokenskyOrderIds struct {
	OrderId string `orm:"pk;column(order_id)"json:"orderId"form:"orderId"`
	//创建时间
	CreateTime time.Time `orm:"auto_now_add;type(datetime);column(create_time)"json:"createTime"form:"createTime"`
}

//生成订单号
func TokenskyOrderIdsRandOrderId(cn string) string {
	list := strings.Split(time.Now().Format("2006-01-02 15:04:05.000"), " ")
	s1s := strings.Split(list[0], "-")
	s2s := strings.Split(list[1], ".")
	code := strings.Join(s1s, "")[2:] + strings.Join(strings.Split(s2s[0], ":"), "") + s2s[1]
	rand.Seed(time.Now().UnixNano())
	rstr := strconv.FormatInt(int64(rand.Int31n(90000)+10000), 10)
	return code + cn + rstr
}

//创建多个订单号 并返回 最多100个
func TokenskyOrderIdsInsertMultiIds(num int, cn string) []string {
	if num > 50 {
		num = 50
	}
	if num < 1 {
		num = 1
	}
	i := 0
	for {
		mapp := make(map[string]struct{}, 0)
		objs := make([]*TokenskyOrderIds, 0, num)
		o := orm.NewOrm()
		o.Begin() //事务
		for i := 0; i < num; i++ {
			oid := TokenskyOrderIdsRandOrderId(cn)
			if _, found := mapp[oid]; !found {
				objs = append(objs, &TokenskyOrderIds{
					OrderId: TokenskyOrderIdsRandOrderId(cn),
				})
				mapp[oid] = struct{}{}
			} else {
				i--
			}
		}
		if _, err := o.InsertMulti(len(objs), &objs); err == nil {
			o.Commit()
			list := make([]string, 0, num)
			for _, obj := range objs {
				list = append(list, obj.OrderId)
			}
			return list
		}
		o.Rollback()
		//防止死循环
		i++
		if i >= 10 {
			list := make([]string, 0)
			for i := 0; i < num; i++ {
				list = append(list, TokenskyOrderIdsRandOrderId(cn))
			}
			i = 0
			return list
		}

	}
}

func TokenskyOrderIdsInsertOne(cn string) string {
	o := orm.NewOrm()
	msg := ""
	for i := 0; i < 10; i++ {
		msg = TokenskyOrderIdsRandOrderId(cn)
		if _, err := o.Insert(&TokenskyOrderIds{OrderId: msg}); err != nil {
			return msg
		}
	}
	return msg
}

//迭代版
func TokenskyOrderIdsIteration() func(sn string) string {
	size := 10
	indexs := make(map[string]int)
	mapp := make(map[string][]string)
	return func(sn string) string {
		list, found := mapp[sn]
		if !found {
			list = TokenskyOrderIdsInsertMultiIds(size, sn)
			mapp[sn] = list
			indexs[sn] = 0
		}
		if indexs[sn] >= size {
			list = TokenskyOrderIdsInsertMultiIds(size, sn)
			mapp[sn] = list
			indexs[sn] = 0
		}
		i := indexs[sn]
		indexs[sn] = i + 1
		return mapp[sn][i]
	}
}
