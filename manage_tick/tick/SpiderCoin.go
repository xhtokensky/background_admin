package tick

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"io/ioutil"
	"net/http"
	"time"
	"tokensky_bg_admin/models"
	"tokensky_bg_admin/utils"
)

const souderCoinUrl = "https://pool.viabtc.com/res/pool/state/new"

type spiderCoinResp struct {
	Code    int                        `json:"code"`
	Data    []*models.SpiderCoinMarket `json:"data"`
	Message string                     `json:"message"`
}

var coinGlobalSoiderSign bool = true
func CoinGlobalSoider()error {
	if coinGlobalSoiderSign{
		coinGlobalSoiderSign = false
		defer func() {coinGlobalSoiderSign=true}()
		now := time.Now()
		client := &http.Client{
			Timeout: 10 * time.Second, //请求超时时间
		}
		request, err := http.NewRequest("GET", souderCoinUrl, nil)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}
		request.Header.Set("Accept", "application/json, text/plain, */*")
		request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh" +
			"; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML" +
			", like Gecko) Chrome/75.0.3770.100 Safari/537.36")
		res, err := client.Do(request)
		if err != nil {
			return nil
		}
		o := orm.NewOrm()
		if body, err := ioutil.ReadAll(res.Body); err == nil {
			resp := &spiderCoinResp{}
			if err := json.Unmarshal(body, resp); err == nil {
				if len(resp.Data) >0{
					for _,obj := range resp.Data{
						obj.IsDate = now
						_,err := o.InsertOrUpdate(obj)
						if err != nil{
							fmt.Println("更新行情失败 err:",err.Error())
							utils.EmailNotify("criticalToAddress","更新行情数据失败","更新行情数据失败",err.Error())
						}
					}
				}
			}
		}else {
			//
			fmt.Println("爬取行情失败 err:",err.Error())
			utils.EmailNotify("criticalToAddress","爬取行情失败","爬取行情数据失败",err.Error())
		}
	}
	return nil
}
