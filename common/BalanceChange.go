package common

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
	"tokensky_bg_admin/conf"
)

//用户资产变动
type balanceChange struct {
	Uid int `json:"uid"`
	//货币类型
	Symbol string `json:"symbol"`
	//操作  add加 sub减 mul乘法 qup除法
	MethodBalance       string  `json:"methodBalance"`
	Balance             string  `json:"balance"`
	balance             float64 `json:"-"`
	MethodFrozenBalance string  `json:"methodFrozenBalance"`
	FrozenBalance       string  `json:"frozenBalance"`
	frozenBalance       float64 `json:"-"`
	SignId              string  `json:"signId"`
}

type requestOne struct {
	//来源 1后端 2后台
	Source int `json:"source"`
	//单个
	Change *balanceChange `json:"change"`
	//说明
	Cont string `json:"cont"`
	//操作模版
	Mold string `json:"mold"`
	//时间戳 单位毫秒
	PushTime int64 `json:"pushTime"`
	//唯一哈希
	HashId string `json:"hashId"`
}

type requestMulti struct {
	//来源 1后端 2后台 3定时任务
	Source int `json:"source"`
	//资产变动多个
	Changes []*balanceChange `json:"changes"`
	//说明
	Cont string `json:"cont"`
	//操作模版
	Mold string `json:"mold"`
	//时间戳 单位毫秒
	PushTime int64 `json:"pushTime"`
	//唯一哈希
	HashId string `json:"hashId"`
}

//用户资产
type userBalance struct {
	Uid           int    `json:"uid"`
	Symbol        string `json:"symbol"`
	Balance       string `json:"balance"`
	FrozenBalance string `json:"frozenBalance"`
}

//响应
type responseOne struct {
	//0正常
	Code int `json:"code"`
	//返回最新数据
	Balance *userBalance `json:"balance"`
	//说明
	Msg string `json:"msg"`
	//哈希
	HashId string `json:"hashId"`
}

type responseMulti struct {
	//0正常
	Code int `json:"code"`
	//返回最新数据
	Balances []*userBalance `json:"balances"`
	//说明
	Msg string `json:"msg"`
	//哈希
	HashId string `json:"hashId"`
}

//用户资产变化
type tokenskyUserBalanceChange struct {
	cont     string
	source   int
	hashId   string
	mold     string
	pushTime int64
	data     []*balanceChange
}

func (this *tokenskyUserBalanceChange) Add(uid int, symbol string, signId string,
	methodBalance string, balance float64, methodFrozenBalance string, frozenBalance float64) {
	var Balance, FrozenBalance string
	if balance != 0 {
		Balance = strconv.FormatFloat(balance, 'f', 8, 64)
	}
	if frozenBalance != 0 {
		FrozenBalance = strconv.FormatFloat(frozenBalance, 'f', 8, 64)
	}
	this.data = append(this.data,&balanceChange{
		Uid:                 uid,
		Symbol:              symbol,
		SignId:              signId,
		MethodBalance:       methodBalance,
		MethodFrozenBalance: methodFrozenBalance,
		Balance:             Balance,
		FrozenBalance:       FrozenBalance,
	})
}

func (this *tokenskyUserBalanceChange)Count()int{
	return len(this.data)
}

func (this *tokenskyUserBalanceChange) Send() (bool, string,string) {
	var bys []byte
	var err error
	num := len(this.data)
	var url string
	//单位毫秒
	this.pushTime = time.Now().UnixNano() / int64(time.Millisecond)

	switch num {
	case 0:
		return false, "没有待处理数据",""
	case 1:
		obj := this.data[0]
		res := requestOne{
			Source:   this.source,
			Cont:     this.cont,
			Mold:     this.mold,
			PushTime: this.pushTime,
			Change:   obj,
		}
		bys, err = json.Marshal(res)
		if err != nil {
			return false, "序列化异常",""
		}
		url = conf.TOKENSKY_BALANCE_CHANGE_URL + "/balance/one"
	default:
		res := requestMulti{
			Source:  this.source,
			Cont:    this.cont,
			Changes: this.data,
		}
		bys, err = json.Marshal(res)
		if err != nil {
			return false, "序列化异常",""
		}
		url = conf.TOKENSKY_BALANCE_CHANGE_URL + "/balance/multi"
	}
	tx := ",\"hashId\":\"" + fmt.Sprintf("%x", md5.Sum(bys)) + "\"}"
	bys = append(bys[:len(bys)-13], []byte(tx)...)
	client := &http.Client{}
	reader := bytes.NewReader(bys)
	request, err := http.NewRequest("POST",url, reader)
	if err != nil {
		return false, "创建请求异常",""
	}
	response, err := client.Do(request)
	if err != nil {
		return false, "请求异常",""
	}
	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	//数据返回异常
	if err != nil {
		return false, "获取数据异常",""
	}
	hashId := ""
	switch num {
	case 1:
		obj := &responseOne{}
		err = json.Unmarshal(body, obj)
		if err != nil {
			//解析异常
			return false, "解析异常",""
		}
		if obj.Code != 0 {
			return false, obj.Msg,""
		}
		hashId = obj.HashId
	default:
		obj := &responseMulti{}
		err = json.Unmarshal(body, obj)
		if err != nil {
			//解析异常
			return false, "解析异常",""
		}
		if obj.Code != 0 {
			return false, obj.Msg,""
		}
		hashId = obj.HashId
	}
	return true, "",hashId
}

func NewTokenskyUserBalanceChange(source int, mold string, cont string) *tokenskyUserBalanceChange {
	return &tokenskyUserBalanceChange{
		cont:   cont,
		source: source,
		mold:   mold,
	}
}


