package utils

import (
	"github.com/astaxie/beego"
	"github.com/zwczou/jpush"
)

var (
	appKey = "9a11d6ce355150887087d0ca"
	secret = "af4025100bbfc437e3df1726"
)

func init()  {

	if str := beego.AppConfig.String("jiguang"+ "::appKey");str != ""{
		appKey = str
	}
	if str := beego.AppConfig.String("jiguang"+ "::secret");str != ""{
		secret = str
	}
}

/*极光推送 所有用户*/
func JiGuangSendAll(alertTitle, alertContent, title, content string) {

	//1初始化客户端
	client := jpush.NewJpushClient(appKey, secret)

	//2获取推送唯一标识符cid
	//cidList, err = client.PushCid(1, "push")

	//推送消息
	payload := &jpush.Payload{
		Platform: jpush.NewPlatform().All(),
		Audience: jpush.NewAudience().All(),
		Notification: &jpush.Notification{
			Alert: "后台推送",
			//提醒
			Android: &jpush.AndroidNotification{
				Alert: alertContent, //提醒内容
				Title: alertTitle,   //提醒标题
			},
			Ios: &jpush.IosNotification{
				Alert: alertContent,
				Sound: title,
			},
			WinPhone: &jpush.WinPhoneNotification{
				Alert: alertContent,
				Title: alertTitle,
			},
		},
		Options: &jpush.Options{
			TimeLive:       60,
			ApnsProduction: false,
		},
		//内容
		Message: &jpush.Message{
			Title:   title,
			Content: content,
		},
	}
	msgId, err := client.Push(payload)
	// msgId, err = client.PushValidate(payload)
	if err != nil {
		//异常 pass
	} else {
		msgId = msgId
	}
	//4创建计划任务
	//client.ScheduleCreate
}

//极光推送单用户
func JiGuangSendByAddr(addr, alertTitle, alertContent, title, content string) {
	//1初始化客户端
	client := jpush.NewJpushClient(appKey, secret)

	//2获取推送唯一标识符cid
	//cidList, err = client.PushCid(1, "push")

	audience := jpush.NewAudience()
	audience.SetRegistrationId(addr)
	//推送消息
	payload := &jpush.Payload{
		Platform: jpush.NewPlatform().All(),
		Audience: audience,
		Notification: &jpush.Notification{
			Alert: "后台推送",
			//提醒
			Android: &jpush.AndroidNotification{
				Alert: alertContent, //提醒内容
				Title: alertTitle,   //提醒标题
			},
			Ios: &jpush.IosNotification{
				Alert: alertContent,
				Sound: title,
			},
			WinPhone: &jpush.WinPhoneNotification{
				Alert: alertContent,
				Title: alertTitle,
			},
		},
		Options: &jpush.Options{
			TimeLive:       60,
			ApnsProduction: false,
		},
		//内容
		Message: &jpush.Message{
			Title:   title,
			Content: content,
		},
	}
	msgId, err := client.Push(payload)
	// msgId, err = client.PushValidate(payload)
	if err != nil {
		//异常 pass
	} else {
		msgId = msgId
	}
	//4创建计划任务
	//client.ScheduleCreate
}

//极光推送多用户
func JiGuangSendByAddrs(addrs []string, alertTitle, alertContent, title, content string)  {
	//1初始化客户端
	client := jpush.NewJpushClient(appKey, secret)

	//2获取推送唯一标识符cid
	//cidList, err = client.PushCid(1, "push")

	audience := jpush.NewAudience()
	audience.SetRegistrationId(addrs...)
	//推送消息
	payload := &jpush.Payload{
		Platform: jpush.NewPlatform().All(),
		Audience: audience,
		Notification: &jpush.Notification{
			Alert: "后台推送",
			//提醒
			Android: &jpush.AndroidNotification{
				Alert: alertContent, //提醒内容
				Title: alertTitle,   //提醒标题
			},
			Ios: &jpush.IosNotification{
				Alert: alertContent,
				Sound: title,
			},
			WinPhone: &jpush.WinPhoneNotification{
				Alert: alertContent,
				Title: alertTitle,
			},
		},
		Options: &jpush.Options{
			TimeLive:       60,
			ApnsProduction: false,
		},
		//内容
		Message: &jpush.Message{
			Title:   title,
			Content: content,
		},
	}
	msgId, err := client.Push(payload)
	// msgId, err = client.PushValidate(payload)
	if err != nil {
		//异常 pass
	} else {
		msgId = msgId
	}
	//4创建计划任务
	//client.ScheduleCreate
}