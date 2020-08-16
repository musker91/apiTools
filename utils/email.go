package utils

import "gopkg.in/gomail.v2"

// 发送邮件传递参数的结构体
type SendMailParams struct {
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
func SendMail(sendMailParams *SendMailParams) (err error) {
	mail := gomail.NewMessage()
	mail.SetHeader("From", "<"+sendMailParams.SendMail+">")        // 发送者
	mail.SetHeader("To", sendMailParams.RecvMail...)               // 接收者
	mail.SetHeader("Subject", sendMailParams.Subject)              //设置邮件主题
	mail.SetBody("text/html", sendMailParams.Body)                 //设置邮件正文
	sendTo := gomail.NewDialer(sendMailParams.SmtpHost, sendMailParams.SmtpPort,
		sendMailParams.SenderEmail, sendMailParams.SenderAuthCode) //创建发送链接

	err = sendTo.DialAndSend(mail)
	if err != nil {
		return
	}
	return
}
