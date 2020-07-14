package modles

import (
	"apiTools/libs/config"
	"apiTools/utils"
	"bytes"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"html/template"
)

const (
	tokenRedisKey = "user_token"
)

var sendEmailTmpl = `
<p>qq号码: {{.QQ}}</p>
<p>消息:  {{.Msg}}</p>
<p>--------------------------------</p>
<p>这封邮件来ApiTools白名单申请</p>
`

type UserVisitAppli struct {
	gorm.Model
	QQ  string `gorm:"size:16;not null" form:"qq" json:"qq" xml:"qq" binding:"required"`
	Msg string `form:"msg" json:"msg" xml:"msg"`
}

type UserToken struct {
	gorm.Model
	QQ    string `gorm:"size:16;not null" json:"qq"`
	Token string `gorm:"not null;unique" json:"token"`
}

// 从数据库中获取user token
func GetTokensFromDB() (tokens []*UserToken, err error) {
	err = SqlConn.Find(&tokens).Error
	if err != nil {
		return
	}
	return
}

// 写入token信息到缓存中
func WriteTokensToCache(tokens []string) (err error) {
	client := RedisPool.Get()
	defer client.Close()
	client.Do("MULTI")
	client.Do("DEL", tokenRedisKey)
	for _, token := range tokens {
		client.Do("LPUSH", tokenRedisKey, token)
	}
	client.Do("EXEC")

	return
}

func CreateToken(userToken *UserToken) (err error) {
	err = SqlConn.Create(userToken).Error
	if err != nil {
		return
	}
	return
}

func DeleteToken(userToken *UserToken) (err error) {
	err = SqlConn.Delete(userToken).Error
	if err != nil {
		return
	}
	return
}

func GetTokensFromCache() (tokens []string, err error) {
	client := RedisPool.Get()
	tokens, err = redis.Strings(client.Do("LRANGE", tokenRedisKey, "0", "-1"))
	if err != nil {
		return
	}
	return
}

func CreateUserVisit(userVisitAppli *UserVisitAppli) (err error) {
	err = SqlConn.Create(userVisitAppli).Error
	if err != nil {
		return
	}

	// 格式化模版
	t := template.Must(template.New("").Parse(sendEmailTmpl))
	buf := new(bytes.Buffer)
	t.Execute(buf, userVisitAppli)

	sendMailStruct := &utils.SendMailStruct{
		SendMail:       config.GetString("email::senderMail"),
		RecvMail:       config.GetStrings("email::recvMail"),
		Subject:        "ApiTools白名单申请",
		Body:           buf.String(),
		SmtpHost:       config.GetString("email::smtpHost"),
		SmtpPort:       config.GetInt("email::smtpPort"),
		SenderEmail:    config.GetString("email::senderMail"),
		SenderAuthCode: config.GetString("email::senderAuthCode"),
	}

	err = utils.SendMail(sendMailStruct)
	if err != nil {
		return
	}
	return
}
