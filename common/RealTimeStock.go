package common

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strings"
	"sync"
	"time"
	"tokensky_bg_admin/conf"
)

type dbExchangeRateValue struct {
	Price float64
}
type dbExchangeRate struct {
	Id     bson.ObjectId `bson:"_id"`
	Symbol string
	Quote  map[string]dbExchangeRateValue
}

//更新ustd汇率
func UpdateUsdtExchangeRate() (float64,bool) {
	log.SetFlags(log.Flags() | log.Lshortfile)
	//连接
	session, err := mgo.Dial(conf.EXCHANGE_RATE_USDT_MGODB_URL)
	if err != nil {
		return 0,false
	}
	defer session.Close()
	//设置模式
	session.SetMode(mgo.Monotonic, true)
	//获取文档集
	collection := session.DB("tokenskyQuoteDB").C("quote")
	exchangeRate := &dbExchangeRate{}
	if err := collection.Find(bson.M{"symbol": "USDT"}).One(exchangeRate); err != nil {
		//索引异常
		return 0,false
	}
	if obj, found := exchangeRate.Quote["USD"]; found {
		return obj.Price,true
	}
	return 0,false
}

//汇率优化
var exchangeData sync.Map

//获取汇率
func GetSymbolExchangeRate(symbol string,currency string)(float64,bool){
	//symbol 虚拟币(如 BTC,USTD)  currency 货币(如美元USD)
	symbol = strings.ToUpper(symbol)
	currency = strings.ToUpper(currency)
	log.SetFlags(log.Flags() | log.Lshortfile)
	//连接
	session, err := mgo.Dial(conf.EXCHANGE_RATE_USDT_MGODB_URL)
	if err != nil {
		return 0,false
	}
	defer session.Close()
	//设置模式
	session.SetMode(mgo.Monotonic, true)
	//获取文档集
	collection := session.DB("tokenskyQuoteDB").C("quote")
	exchangeRate := &dbExchangeRate{}
	if err := collection.Find(bson.M{"symbol": symbol}).One(exchangeRate); err != nil {
		//索引异常
		return 0,false
	}
	if obj, found := exchangeRate.Quote["USD"]; found {
		return obj.Price,true
	}
	return 0,false
}

//获取汇率
func GetSymbolExchangeRate2(symbol string,currency string)(float64,bool){
	//symbol 虚拟币(如 BTC,USTD)  currency 货币(如美元USD)
	symbol = strings.ToUpper(symbol)
	currency = strings.ToUpper(currency)
	//时间
	now := time.Now()
	tm := time.Date(now.Year(),now.Month(),now.Day(),now.Hour(),now.Minute(),0,0,now.Location()).Unix()
	name := symbol + "&&" + currency
	if v,ok := exchangeData.Load(name);ok{
		if mapp,ok := v.(map[int64]float64);ok{
			if con,found := mapp[tm];found{
				return con,true
			}
		}
	}
	log.SetFlags(log.Flags() | log.Lshortfile)
	//连接
	session, err := mgo.Dial(conf.EXCHANGE_RATE_USDT_MGODB_URL)
	if err != nil {
		return 0,false
	}
	defer session.Close()
	//设置模式
	session.SetMode(mgo.Monotonic, true)
	//获取文档集
	collection := session.DB("tokenskyQuoteDB").C("quote")
	exchangeRate := &dbExchangeRate{}
	if err := collection.Find(bson.M{"symbol": symbol}).One(exchangeRate); err != nil {
		//索引异常
		return 0,false
	}
	obj, found := exchangeRate.Quote["USD"]
	if !found{
		return 0,false
	}
	exchangeData.Store(name,map[int64]float64{tm:obj.Price})
	return obj.Price,true
}