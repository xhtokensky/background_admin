package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io"
	"io/ioutil"
	"net/http"
	"tokensky_bg_admin/manage_balance/httpModel"
)

// fileLogs 生产环境下日志
var FileLogs *logs.BeeLogger

//请求实例

/*
{
	"source": 1,
	"change": {
		"uid": 1,
		"methodBalance": "add",
		"balance": "123",
		"methodFrozenBalance": "add",
		"frozenBalance": "0",
		"symbol":"货币"
	},
	"cont": "提币"
}
*/

func SetBalanceOne(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	var err error
	var ok bool
	var msg string
	resp := &httpModel.ResponseOne{Code: 0}
	reqOne := &httpModel.RequestOne{}
	//请求体
	reqBys, err := ioutil.ReadAll(req.Body) //把  body 内容读入字符串
	if err != nil {
		//返回错误信息
		resp.Code = 1
		resp.Msg = "请求异常:"+err.Error()
		log := fmt.Sprintf("请求异常 req:%+v err:%s",req.Body,err.Error())
		FileLogs.Error(log)
		goto End
	}
	err = json.Unmarshal(reqBys, reqOne)
	if err != nil {
		//返回错误信息
		resp.Code = 1
		log := fmt.Sprintf("解析数据异常 req:%s err:%s",string(reqBys),err.Error())
		resp.Msg = "解析数据异常: err:" + err.Error()
		FileLogs.Error(log)
		goto End
	}
	if reqOne.Change == nil {
		resp.Code = 1
		resp.Msg = "请求Change数据异常"
		log := fmt.Sprintf("请求Change数据异常 req:%+v",reqOne)
		FileLogs.Error(log)
		goto End
	}

	//业务处理
	ok, msg,resp.Balance = httpModel.BalanceChangeIsOne(reqOne)
	if !ok {
		resp.Code = 1
		resp.Msg = msg
		log := fmt.Sprintf("业务处理异常 req:%+v resp:%+v",reqOne,resp)
		FileLogs.Alert(log)
		goto End
	}
	resp.HashId = reqOne.HashId
End:
	w.WriteHeader(200)
	bys, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(400)
		bys = []byte("400")
		log := fmt.Sprintf("返回数据序列化异常 req:%+v resp:%+v",reqOne,resp)
		FileLogs.Alert(log)
	}
	_, err = io.WriteString(w, string(bys)) //返回内容
	if err != nil {
		//返回异常错误,打印错误日志
		log := fmt.Sprintf("返回中途异常 req:%+v resp:%+v",reqOne,resp)
		FileLogs.Alert(log)
	}
}

/*
{
	"source": 1,
	"changes": [{
		"uid": 1,
		"methodBalance": "add",
		"balance": "123",
		"methodFrozenBalance": "add",
		"frozenBalance": "0",
		"symbol": "货币"
	}, {
		"uid": 2,
		"methodBalance": "add",
		"balance": "123",
		"methodFrozenBalance": "add",
		"frozenBalance": "0",
		"symbol": "货币"
	}],
	"cont": "提币"
}
*/

func SetBalanceMulti(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	var err error
	var ok bool
	var msg string
	resp := &httpModel.ResponseMulti{Code: 0}
	reqMulti := &httpModel.RequestMulti{}
	//请求体
	reqBys, err := ioutil.ReadAll(req.Body) //把  body 内容读入字符串
	if err != nil {
		//返回错误信息
		resp.Code = 1
		resp.Msg = "请求异常:" + err.Error()
		log := fmt.Sprintf("请求异常 req:%+v err:%s",req.Body,err.Error())
		FileLogs.Error(log)
		goto End
	}
	err = json.Unmarshal(reqBys, reqMulti)
	if err != nil {
		//返回错误信息
		resp.Code = 1
		log := fmt.Sprintf("解析数据异常 req:%s err:%s",string(reqBys),err.Error())
		resp.Msg = "请求异常: err" + err.Error()
		FileLogs.Error(log)
		goto End
	}
	if len(reqMulti.Changes) <= 0 {
		resp.Code = 1
		resp.Msg = "请求Changes 数据不存在"
		log := fmt.Sprintf("请求Changes数据异常 req:%+v",reqMulti)
		FileLogs.Error(log)
		goto End
	}

	//业务处理
	ok, msg,resp.Balances = httpModel.BalanceChangeIsMulti(reqMulti)
	if !ok {
		resp.Code = 1
		resp.Msg = msg
		log := fmt.Sprintf("业务处理异常 req:%+v resp:%+v",reqMulti,resp)
		FileLogs.Alert(log)
		goto End
	}
	resp.HashId = reqMulti.HashId
End:
	w.WriteHeader(200)
	bys, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(400)
		bys = []byte("400")
		log := fmt.Sprintf("返回数据序列化异常 req:%+v resp:%+v",reqMulti,resp)
		FileLogs.Alert(log)
	}
	_, err = io.WriteString(w, string(bys)) //返回内容
	if err != nil {
		//返回异常错误,打印错误日志
		log := fmt.Sprintf("返回中途异常 req:%+v resp:%+v",reqMulti,resp)
		FileLogs.Alert(log)
	}
}
