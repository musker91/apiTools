package utils

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/robfig/cron/v3"
	"gopkg.in/gomail.v2"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

// 获取项目根目录
func GetRootPath() (rootPath string) {
	rootPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(fmt.Sprintf("get project root path faild: %s", err))
	}
	return
}

// 字符串计算md5值
func GetMD5(text string) string {
	h := md5.New()
	salt := "apiTools"
	io.WriteString(h, text+salt)
	urlmd5 := fmt.Sprintf("%x", h.Sum(nil))
	return urlmd5
}

// 获取随机的唯一短串
func GetShortStr() (tiny string) {
	// 时间戳随机加盐避免重复
	rand.Seed(time.Now().UnixNano() >> 3)
	num := rand.Int63n(time.Now().UnixNano() >> 3)
	alpha := merge(getRange(48, 57), getRange(65, 90))
	alpha = merge(alpha, getRange(97, 122))
	if num < 62 {
		tiny = string(alpha[num])
		return tiny
	} else {
		var runes []rune
		runes = append(runes, alpha[num%62])
		num = num / 62
		for num >= 1 {
			if num < 62 {
				runes = append(runes, alpha[num-1])
			} else {
				runes = append(runes, alpha[num%62])
			}
			num = num / 62
		}
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		tiny = string(runes)
		return
	}
}

func getRange(start, end rune) (ran []rune) {
	for i := start; i <= end; i++ {
		ran = append(ran, i)
	}
	return ran
}

func merge(a, b []rune) []rune {
	c := make([]rune, len(a)+len(b))
	copy(c, a)
	copy(c[len(a):], b)
	return c
}

// 重新定义cron定时任务初始化
func NewWithCron() *cron.Cron {
	secondParser := cron.NewParser(cron.Second | cron.Minute |
		cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
	return cron.New(cron.WithParser(secondParser), cron.WithChain())
}

// 定时器，用于堵塞进程
// second 定时时间 秒
func TimerUtil(second int64) {
	timer := time.NewTimer(time.Second * time.Duration(second))
	<-timer.C
}

// 判断值是否在一个切片中存在
func IsInSelic(data string, slice []string) bool {
	for _, s := range slice {
		if s == data {
			return true
		}
	}
	return false
}

// 判断字符串是否都是数字
func IsDigit(src string) bool {
	pattern := "\\d+" //反斜杠要转义
	result, err := regexp.MatchString(pattern, src)
	if err != nil {
		return false
	}

	return result
}

// 生成一个随机token
func CreateToken() string {
	// 获取当前时间的时间戳
	t := time.Now().Unix()

	// 生成一个MD5的哈希
	h := md5.New()

	// 将时间戳转换为byte，并写入哈希
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(t))
	h.Write([]byte(b))

	// 将字节流转化为16进制的字符串
	return hex.EncodeToString(h.Sum(nil))
}

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

// 正则匹配分组
func RegexMatchGroup(compile string, s string) (result map[string]string, err error) {
	videoMatch, err := regexp.Compile(compile)
	if err != nil {
		return
	}
	result = make(map[string]string)

	submatch := videoMatch.FindStringSubmatch(s)

	if len(submatch) == 0 {
		return
	}

	groupNames := videoMatch.SubexpNames()

	for i, name := range groupNames {
		if i != 0 && name != "" {
			result[name] = submatch[i]
		}
	}

	return
}

//判断文件文件夹是否存在
func IsFileExist(path string) (bool, error) {
	fileInfo, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false, nil
	}
	//我这里判断了如果是0也算不存在
	if fileInfo.Size() == 0 {
		return false, nil
	}
	if err == nil {
		return true, nil
	}
	return false, err
}

func InterfaceToBytes(v interface{}) (result []byte, err error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return
	}
	result = bytes
	return
}


func RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Int63n(max-min) + min
}