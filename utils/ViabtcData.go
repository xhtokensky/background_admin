package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
	"tokensky_bg_admin/conf"
)

//拉取ViaBTC 数据

const (
	Viabtc_History_Api_Key    = "601bbe87e5e8e1446e701f5cfc47f5be"
	Viabtc_History_Secret_Key = "f58672bef89f93a60d829e97044721c0375bcd728fbb8ed1bb29f0d174ce2403"

	//获取历史收益
	viabtc_Profit_History_Url = "https://pool.viabtc.com/res/openapi/v1/profit/history"
	//获取历史算力
	viabtc_Hashrate_History_Url = "https://pool.viabtc.com/res/openapi/v1/hashrate/history"
)

//历史收益
type HistoryProfit struct {
	Coin        string  `json:"coin"`
	Date        string  `json:"date"`
	PplnsProfit string   `json:"pplns_profit"` // pplns收益
	pplnsProfit float64 `json:"-"`
	PpsProfit   string `json:"pps_profit"`   //pps收益
	ppsProfit float64 `json:"-"`
	SoloProfit  string   `json:"solo_profit"`  //solo收益
	soloProfit int64 `json:"-"`
	TotalProfit string `json:"total_profit"` //总收益
	totalProfit float64 `json:"-"`
}

type HistoryProfitJson struct {
	Data      []*HistoryProfit `json:"data"`
	Count     int64           `json:"count"`
	CurrPage  int64           `json:"curr_page"`
	HasNext   bool            `json:"has_next"`
	Total     int64           `json:"total"`
	TotalPage int64           `json:"total_page"`
}




type HistoryProfitJsonData struct {
	Code    int               `json:"code"`
	Data    *HistoryProfitJson `json:"data"`
	Message string            `json:"message"`
}

//历史算力
type HistoryHashrate struct {
	Coin       string `json:"coin"`
	Date       string `json:"date"`
	date       time.Time
	Hashrate   string `json:"hashrate"`
	hashrate int64 `json:"-"`
	RejectRate string `json:"reject_rate"`
	rejectRate float64 `json:"-"`
}

type HistoryHashrateJson struct {
	Data      []*HistoryHashrate `json:"data"`
	Count     int64             `json:"count"`
	CurrPage  int64             `json:"curr_page"`
	HasNext   bool              `json:"has_next"`
	Total     int64             `json:"total"`
	TotalPage int64             `json:"total_page"`
}

//
type HistoryHashrateData struct {
	Code    int                 `json:"code"`
	Data    *HistoryHashrateJson `json:"data"`
	Message string              `json:"message"`
}

/*
| coin | string | yes | 币种 |
| start_date | string | no | 起始日期 |
| end_date | string | no | 截止日期 |
| utc | string | no | true/false, 默认false |
| page | int | no | 页码 |
| limit | int | no | 每页条数 |
*/

//获取历史收益
func GetViabtcProfitHistoryData(tm int64, coin string) (float64, bool) {
	//coin 货币类型 tm时间
	tm1 := time.Unix(tm, 0)
	//开始时间年月日
	tm2 := tm1.AddDate(0, 0, 1)
	params := []string{
		"coin=" + coin,                           //货币类型
		"start_date=" + tm1.Format("2006-01-02"), //开始时间
		"end_date=" + tm2.Format("2006-01-02"),   //结束时间
		"page=1",                                 //页码
		"limit=1000",                             //条数
	}
	//url传参
	paramUrl := strings.Join(params, "&")
	url := viabtc_Profit_History_Url + "?" + paramUrl
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil && request != nil {
		LogCritical(fmt.Sprintf("GetViabtcProfitHistoryData 获取历史收益 创建请求对象失败 err :%+v", err))
		return 0, false
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-API-KEY", Viabtc_History_Api_Key)
	request.Header.Set("X-SIGNATURE", Viabtc_History_Secret_Key)

	response, err := client.Do(request)
	if err != nil {
		LogCritical(fmt.Sprintf("GetViabtcProfitHistoryData 获取历史收益 发起请求失败 err :%+v", err))
		return 0, false
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		LogCritical(fmt.Sprintf("GetViabtcProfitHistoryData 获取历史收益 加载数据失败 err :%+v", err))
		return 0, false
	}
	//数据解析
	req := HistoryProfitJsonData{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		LogCritical(fmt.Sprintf("GetViabtcProfitHistoryData 获取历史收益 数据解析异常 err :%+v", err))
		return 0, false
	}
	if req.Data != nil{
		for _,obj := range req.Data.Data{
			if obj.TotalProfit != ""{
				obj.totalProfit,err = strconv.ParseFloat(obj.TotalProfit,64)
				if err != nil{
					LogCritical(fmt.Sprintf("GetViabtcProfitHistoryData 获取历史收益 解析TotalProfit err :%+v", err))
					return 0, false
				}
				obj.ppsProfit,err = strconv.ParseFloat(obj.PpsProfit,64)
				if err != nil{
					LogCritical(fmt.Sprintf("GetViabtcProfitHistoryData 获取历史收益 解析PplnsProfit err :%+v", err))
					return 0, false
				}
				obj.pplnsProfit,err = strconv.ParseFloat(obj.PplnsProfit,64)
				if err != nil{
					LogCritical(fmt.Sprintf("GetViabtcProfitHistoryData 获取历史收益 解析PplnsProfit err :%+v", err))
					return 0, false
				}
				obj.soloProfit,err = strconv.ParseInt(obj.SoloProfit,10,64)
				if err != nil{
					LogCritical(fmt.Sprintf("GetViabtcProfitHistoryData 获取历史收益 解析SoloProfit err :%+v", err))
					return 0, false
				}
			}
		}
	}


	if len(req.Data.Data) == 0 {
		LogCritical(fmt.Sprintf("GetViabtcProfitHistoryData 获取历史收益 拉取数据不存在"))
		return 0, false
	}

	//总收益
	if req.Data.Data[0].Coin != coin {
		LogCritical(fmt.Sprintf("GetViabtcProfitHistoryData 获取历史收益 货币类型错误"))
		return 0, false
	}

	//无收益
	if req.Data.Data[0].totalProfit < conf.FLOAT_PRECISE_8 {
		LogCritical(fmt.Sprintf("GetViabtcProfitHistoryData 获取历史收益 无收益"))
		return 0, false
	}
	for _,obj := range req.Data.Data{
		objTm,err := time.Parse("2006-01-02", obj.Date)
		if err == nil{
			objTm = time.Date(objTm.Year(),objTm.Month(),objTm.Day(),0,0,0,0,time.Local)
			if objTm.Unix() == tm{
				return obj.totalProfit,true
			}
		}
	}
	LogCritical(fmt.Sprintf("GetViabtcProfitHistoryData 历史收益不存在 无收益"))
	//收益正常
	return req.Data.Data[0].totalProfit, true
}

/*
|参数名|类型|必须|备注|
| -----|----|----|----|
| coin | string | yes | 币种 |
| start_date | string | no | 起始日期，格式 2019-01-24 |
| end_date | string | no | 截止日期 |
| utc | string | no | true/false, 默认false |
| page | int | no | 页码 |
| limit | int | no | 每页条数 |
*/

//获取历史算力
func GetViabtcHashrateHistoryData(tm int64, coin string) (int64, bool) {
	//coin 货币类型 tm时间
	tm1 := time.Unix(tm, 0)
	//开始时间年月日
	tm2 := tm1.AddDate(0, 0, 1)
	params := []string{
		"coin=" + coin,                           //货币类型
		"start_date=" + tm1.Format("2006-01-02"), //开始时间
		"end_date=" + tm2.Format("2006-01-02"),   //结束时间
		"page:=",                                 //页码
		"limit=100",                              //条数
	}
	//url传参
	paramUrl := strings.Join(params, "&")
	url := viabtc_Hashrate_History_Url + "?" + paramUrl
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil && request != nil {
		LogCritical(fmt.Sprintf("GetViabtcHashrateHistoryData 获取历史算力 创建请求对象失败 err:%+v", err))
		return 0, false
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-API-KEY", Viabtc_History_Api_Key)
	request.Header.Set("X-SIGNATURE", Viabtc_History_Secret_Key)

	response, err := client.Do(request)
	if err != nil {
		LogCritical(fmt.Sprintf("GetViabtcHashrateHistoryData 获取历史算力 发起请求失败 err:%+v", err))
		return 0, false
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		LogCritical(fmt.Sprintf("GetViabtcHashrateHistoryData 获取历史算力 加载数据失败 err:%+v", err))
		return 0, false
	}
	//数据解析
	req := HistoryHashrateData{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		LogCritical(fmt.Sprintf("GetViabtcHashrateHistoryData 获取历史算力 数据解析异常2 err:%+v", err))
		return 0, false
	}
	if req.Data!=nil{
		for _,obj := range req.Data.Data{
			obj.hashrate,err = strconv.ParseInt(obj.Hashrate,10,64)
			if err != nil{
				LogCritical(fmt.Sprintf("GetViabtcHashrateHistoryData 获取历史算力 Hashrate解析异常 err:%+v", err))
				return 0, false
			}
			obj.rejectRate,err = strconv.ParseFloat(obj.RejectRate,64)
			if err != nil{
				LogCritical(fmt.Sprintf("GetViabtcHashrateHistoryData 获取历史算力 RejectRate解析异常 err:%+v", err))
				return 0, false
			}
		}
	}

	if len(req.Data.Data) == 0 {
		LogCritical(fmt.Sprintf("GetViabtcHashrateHistoryData 获取历史算力 不存在数据"))
		return 0, false
	}

	//总收益
	if req.Data.Data[0].Coin != coin {
		LogCritical(fmt.Sprintf("GetViabtcHashrateHistoryData 获取历史算力 货币类型错误"))
		return 0, false
	}
	//无收益
	if req.Data.Data[0].hashrate <= 0 {
		LogCritical(fmt.Sprintf("GetViabtcHashrateHistoryData 获取历史算力 无收益"))
		return 0, false
	}
	for _,obj := range req.Data.Data{
		objTm,err := time.Parse("2006-01-02", obj.Date)
		if err == nil{
			objTm = time.Date(objTm.Year(),objTm.Month(),objTm.Day(),0,0,0,0,time.Local)
			if objTm.Unix() == tm{
				return obj.hashrate,true
			}
		}
	}
	LogCritical(fmt.Sprintf("GetViabtcHashrateHistoryData 算力数据不存在 无收益"))
	//收益正常
	return req.Data.Data[0].hashrate, true
}
