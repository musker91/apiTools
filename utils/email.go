package utils

import "gopkg.in/gomail.v2"

// 发送邮件传递参数的结构体
type SendMailStruct struct {
	SendMail       string
	RecvMail       []string
	Subject        string
	Body           string
	SmtpHost       string
	SmtpPort       int
	SenderEmail    string
	SenderAuthCode string
}

// 发送邮件
func SendMail(sendMailStruct *SendMailStruct) (err error) {
	mail := gomail.NewMessage()
	mail.SetHeader("From", "<"+sendMailStruct.SendMail+">")        // 发送者
	mail.SetHeader("To", sendMailStruct.RecvMail...)               // 接收者
	mail.SetHeader("Subject", sendMailStruct.Subject)              //设置邮件主题
	mail.SetBody("text/html", sendMailStruct.Body)                 //设置邮件正文
	sendTo := gomail.NewDialer(sendMailStruct.SmtpHost, sendMailStruct.SmtpPort,
		sendMailStruct.SenderEmail, sendMailStruct.SenderAuthCode) //创建发送链接

	err = sendTo.DialAndSend(mail)
	if err != nil {
		return
	}
	return
}
