package utils

import (
	"encoding/json"
	"net/smtp"
	"strings"
	"time"
	"tokensky_bg_admin/conf"
)

/* 邮件通知系统 */

const (
	email_host  = "smtp.qq.com:587"
	email_max   = 200
	email_sleep = 10
)

type emailData struct {
	errLevel string
	title    string
	content  string
	data     interface{}
}

var emailChans chan emailData

func init() {
	emailChans = make(chan emailData, email_max)
}

func EmailGo() {
	if conf.EMAIL_NOTIFY_SIGN {
	Go:
		for len(emailChans) > 0 {
			obj := <-emailChans
			EmailNotify(obj.errLevel, obj.title, obj.content, obj.data)
		}
		time.Sleep(time.Second * email_sleep)
		goto Go
	}
}

func EmailAdd(errLevel string, title string, content string, data interface{}) {
	if conf.EMAIL_NOTIFY_SIGN {
		emailChans <- emailData{errLevel, title, content, data}
	}
}

func EmailNotify(errLevel string, title string, content string, data interface{}) {
	addrs := make([]string, 0)
	user := conf.EMAIL_SEND_USER_NAME
	pwd := conf.EMAIL_SEND_USER_PWD
	switch errLevel {
	case conf.EMAIL_ERROR_LEVEL_ERROR:
		addrs = conf.EMAIL_ADDRESS_ERROR
	case conf.EMAIL_ERROR_LEVEL_CRITICAL:
		addrs = conf.EMAIL_ADDRESS_CRITICAL
	case conf.EMAIL_ERROR_LEVEL_ALERT:
		addrs = conf.EMAIL_ADDRESS_ALERT
	case conf.EMAIL_ERROR_LEVEL_EMERGENCY:
		addrs = conf.EMAIL_ADDRESS_EMERGENCY
	}
	if user != "" && pwd != "" && len(addrs) > 0 {
		now := time.Now().Format("2006-01-02 15:04:05")
		subject := "后台异常:" + title + " 级别:" + errLevel + " 时间:" + now
		msg := content
		if bus, err := json.Marshal(data); err == nil {
			msg += " \n 值:" + string(bus)
		}
		for _, to := range addrs {
			if err := sendToMail(user, pwd, email_host, to, subject, msg, "html"); err != nil {
				//异常
			}
		}
	}
}

func sendToMail(user, password, host, to, subject, body, mailtype string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}
	msg := []byte("To: " + to + "\r\nFrom: " + user + "\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	err := smtp.SendMail(host, auth, user, []string{to}, msg)
	return err
}
