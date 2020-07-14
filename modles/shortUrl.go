package modles

import (
	"apiTools/utils"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/jonsen/gotld"
	"github.com/pkg/errors"
	"io/ioutil"
	"math/rand"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	reqFailErr  = "request exception"
	tranFailErr = "transform fail"
)

var (
	toShortApiData    map[string]interface{}
	parseShortApiData map[string]interface{}
)

// 接收转换为短链接 api form表单
type ShortForm struct {
	Url        string `form:"url" json:"url" xml:"url" binding:"required"`   // 要进行转换的长链
	Domain     string `form:"domain" json:"domain" xml:"domain"`             // 短链接域名绑定自己的域名
	ExpireTime int    `form:"expireTime" json:"expireTime" xml:"expireTime"` // 设置过期时间, (以分钟为单位, -1代表用不过期)
	Type       int    `form:"type" json:"type" xml:"type"`                   // 要转换短链接的类型
}

// 返回的生成短链接信息数据结构
type ShortInfo struct {
	LongUrlMd5 string // 长链接转换后的md5值
	LongUrl    string // 原来长链接
	Domain     string // 短链接域名绑定的域名
	ShortStr   string // 短链接串
}

// 检测转换的域名是否存在，如果存在则返回数据
func checkShortUrl(shortForm *ShortForm) (*ShortInfo, bool, error) {
	shortInfo := &ShortInfo{
		LongUrlMd5: utils.GetMD5(shortForm.Url + shortForm.Domain),
	}
	// 获取redis连接
	redisClient := RedisPool.Get()
	defer redisClient.Close()
	exists, err := redis.Bool(redisClient.Do("EXISTS",
		fmt.Sprintf("short_long_%s", shortInfo.LongUrlMd5)))
	if err != nil || !exists {
		_, _ = redisClient.Do("DEL", fmt.Sprintf("short_long_%s", shortInfo.LongUrlMd5))
		return shortInfo, false, fmt.Errorf("check long exist, err: %v", err)
	}
	// 根据长串获取存储的短域名信息
	LongInfoBytes, err := redis.ByteSlices(redisClient.Do("HMGET",
		fmt.Sprintf("short_long_%s", shortInfo.LongUrlMd5),
		"longUrl", "domain", "shortStr"))
	if err != nil || len(LongInfoBytes) != 3 {
		return shortInfo, false, fmt.Errorf("get long info fail, err: %v", err)
	}
	// 赋值
	shortInfo.LongUrl = string(LongInfoBytes[0])
	shortInfo.Domain = string(LongInfoBytes[1])
	shortInfo.ShortStr = string(LongInfoBytes[2])

	return shortInfo, true, nil
}

// 转换成短链接
func ToShortUrl(shortForm *ShortForm) (*ShortInfo, error) {
	// 校验长连接是否存在
	tmpShortInfo, status, err := checkShortUrl(shortForm)
	if status && err == nil {
		return tmpShortInfo, nil
	}
	// 创建一个shortInfo数据对象
	shortInfo := &ShortInfo{
		LongUrl:    shortForm.Url,
		LongUrlMd5: utils.GetMD5(shortForm.Url + shortForm.Domain),
		ShortStr:   utils.GetShortStr(),
	}
	redisClient := RedisPool.Get()
	defer redisClient.Close()

	// domain
	_, _, err = gotld.GetTld(shortForm.Domain)
	if err != nil {
		return shortInfo, fmt.Errorf("short form domain [%s] not normal domain name", shortForm.Domain)
	}
	if ok := strings.HasPrefix(shortForm.Domain, "http"); !ok {
		shortForm.Domain = fmt.Sprintf("http://%s", shortForm.Domain)
	}
	parse, err := url.Parse(shortForm.Domain)
	if err != nil {
		return shortInfo, fmt.Errorf("short domain [%s] name parse fail", shortForm.Domain)
	}
	shortInfo.Domain = parse.Hostname()
	// 事务开始创建数据到redis
	_ = redisClient.Send("MULTI")
	// 设置hash   长地址md5值对应相关数据
	_ = redisClient.Send(
		"HMSET",
		fmt.Sprintf("short_long_%s", shortInfo.LongUrlMd5),
		"longUrl", shortInfo.LongUrl,
		"domain", shortInfo.Domain,
		"shortStr", shortInfo.ShortStr,
	)
	// 设置string 短串对应长地址md5值
	_ = redisClient.Send("SET", fmt.Sprintf("short_%s", shortInfo.ShortStr),
		shortInfo.LongUrlMd5)
	// 设置过期时间, 默认为用不过期
	if shortForm.ExpireTime != -1 {
		expireTime := 60 * shortForm.ExpireTime

		// 设置long md5过期
		_ = redisClient.Send("EXPIRE", fmt.Sprintf("short_long_%s", shortInfo.LongUrlMd5), expireTime)
		// 设置short str过期
		_ = redisClient.Send("EXPIRE", fmt.Sprintf("short_%s", shortInfo.ShortStr), expireTime)
	}
	_, err = redisClient.Do("EXEC")
	if err != nil {
		_, _ = redisClient.Do("DISCARD")
		return shortInfo, fmt.Errorf("set short data to redis fail, err: %v", err)
	}
	return shortInfo, nil
}

// 解析短链接
func ParseShort(shortUrl string) (*ShortInfo, error) {
	// 创建一个数据返回对象
	shortInfo := &ShortInfo{}
	// 解析短地址url
	urlParse, err := url.Parse(shortUrl)
	if err != nil {
		return shortInfo, fmt.Errorf("parse short url fail, err: %v", err)
	}
	// 短链接域名
	shortInfo.Domain = urlParse.Hostname()
	// 短链接短串
	if ok := strings.HasPrefix(urlParse.Path, "/"); !ok {
		return shortInfo, errors.New("short url not has param")
	}
	shortPathSlice := strings.Split(urlParse.Path, "/")
	if len(shortPathSlice) > 3 || len(shortPathSlice) < 2 {
		return shortInfo, errors.New("short url param malformed")
	}
	shortStr := shortPathSlice[1]
	if shortUrl == "" {
		return shortInfo, errors.New("short url param malformed")
	}
	shortInfo.ShortStr = shortStr
	// 获取redis连接
	redisClient := RedisPool.Get()
	defer redisClient.Close()
	// 根据短串获取长串md5
	longUrlMd5, err := redis.String(redisClient.Do("GET",
		fmt.Sprintf("short_%s", shortInfo.ShortStr)))
	if err != nil {
		return shortInfo, fmt.Errorf("get long url md5 fail from short str, err: %v", err)
	}
	if longUrlMd5 == "" {
		_, _ = redisClient.Do("DEL", fmt.Sprintf("short_%s", shortInfo.ShortStr))
		return shortInfo, errors.New("get long url md5 fail, because is empty")
	}
	shortInfo.LongUrlMd5 = longUrlMd5

	// 根据长串获取存储的短域名信息
	LongInfoBytes, err := redis.ByteSlices(redisClient.Do("HMGET", fmt.Sprintf("short_long_%s", shortInfo.LongUrlMd5),
		"longUrl", "domain", "shortStr"))
	if err != nil || len(LongInfoBytes) != 3 {
		return shortInfo, fmt.Errorf("get long info fail, err: %v", err)
	}
	// 赋值
	shortInfo.LongUrl = string(LongInfoBytes[0])
	shortInfo.Domain = string(LongInfoBytes[1])
	shortInfo.ShortStr = string(LongInfoBytes[2])

	return shortInfo, nil
}

func InitShortData() (err error) {
	data, err := ioutil.ReadFile(filepath.Join(utils.GetRootPath(), "data/shorturl/api.json"))
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &toShortApiData)
	if err != nil {
		err = fmt.Errorf("parse short url api data fail, err: %v", err)
		return
	}

	data, err = ioutil.ReadFile(filepath.Join(utils.GetRootPath(), "data/shorturl/parseapi.json"))
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &parseShortApiData)
	if err != nil {
		err = fmt.Errorf("parse short url api data fail, err: %v", err)
		return
	}

	return
}

// 第三方短链接转换
func OfficialToShort(form *ShortForm) (shortInfo *ShortInfo, msg string, err error) {
	api, err := getOneOffShortApi(form.Type, 0)
	if err != nil {
		msg = reqFailErr
		return
	}
	var shortUrl string
	if d, ok := api.(string); ok {
		shortUrl, msg = singleApiReq(form.Url, d)
	} else if d1, ok1 := api.(map[string]interface{}); ok1 {
		shortUrl, msg = jsonApiReq(form.Url, d1)
	} else {
		msg = reqFailErr
		return
	}
	if msg != "" {
		err = errors.New(msg)
		return
	}
	shortInfo = &ShortInfo{
		LongUrl:  form.Url,
		ShortStr: shortUrl,
	}
	return
}

// 第三方短链接解析
func OfficialParseShort(form *ShortForm) (shortInfo *ShortInfo, msg string, err error) {
	api, err := getOneOffShortApi(form.Type, 1)
	if err != nil {
		msg = reqFailErr
		return
	}
	var longUrl string
	if d, ok := api.(string); ok {
		longUrl, msg = singleApiReq(form.Url, d)
	} else if d1, ok1 := api.(map[string]interface{}); ok1 {
		longUrl, msg = jsonApiReq(form.Url, d1)
	} else {
		msg = reqFailErr
		return
	}
	if msg != "" {
		err = errors.New(msg)
		return
	}
	shortInfo = &ShortInfo{
		LongUrl:  longUrl,
		ShortStr: form.Url,
	}
	return
}

// 获取生成短链接的api接口
func getOneOffShortApi(t int, ty int) (d interface{}, err error) {
	var key string
	if ty == 1 { // 解析短链接
		switch t {
		case 1:
			key = "t.cn"
		case 2:
			key = "url.cn"
		default:
			key = "all"
		}
	} else { // 转换为短链接
		switch t {
		case 1:
			key = "t.cn"
		default:
			key = "url.cn"
		}
	}
	var di interface{}
	var ok bool
	if ty == 1 { // 解析短链接
		di, ok = parseShortApiData[key]
		if !ok {
			err = fmt.Errorf("get type is '%v' parse short api data fail", t)
			return
		}
	} else { // 转换为短链接
		di, ok = toShortApiData[key]
		if !ok {
			err = fmt.Errorf("get type is '%v' to short api data fail", t)
			return
		}
	}

	if ds, ok := di.([]interface{}); ok {
		rand.Seed(time.Now().UnixNano())
		intn := rand.Intn(len(ds))
		d = ds[intn]
	} else {
		err = fmt.Errorf("get type is '%v' short api data fail", t)
		return
	}
	return
}

// 当字符串类型的接口请求方法
func singleApiReq(formUrl string, apiUrl interface{}) (url string, msg string) {
	api_url := apiUrl.(string)
	api_url = strings.Replace(api_url, "<%url%>", formUrl, 1)
	var err error
	var data []byte
	for i := 0; i < 3; i++ {
		proxyIp := GetOneProxyIp("standardProxypool")
		data, _, err = utils.HttpProxyGet(api_url, proxyIp, nil)
		if err != nil {
			_ = DelOneProxyFromRedis("standardProxypool", proxyIp)
			continue
		}
		err = nil
		break
	}
	if err != nil {
		data, _, err = utils.HttpProxyGet(api_url, "", nil)
	}
	if err != nil {
		msg = tranFailErr
		return
	}
	bodyData := strings.TrimSpace(string(data))
	if !strings.HasPrefix(bodyData, "http") {
		msg = tranFailErr
		return
	}
	url = bodyData
	return
}

// json结构类型的接口请求方法
func jsonApiReq(longUrl string, apiUrl map[string]interface{}) (shortUrl string, msg string) {
	var api_url string
	if a, ok := apiUrl["url"]; ok {
		api_url = a.(string)
	} else {
		msg = reqFailErr
		return
	}
	api_url = strings.Replace(api_url, "<%url%>", longUrl, 1)
	var err error
	var data []byte
	for i := 0; i < 3; i++ {
		proxyIp := GetOneProxyIp("standardProxypool")
		data, _, err = utils.HttpProxyGet(api_url, proxyIp, nil)
		if err != nil {
			_ = DelOneProxyFromRedis("standardProxypool", proxyIp)
			continue
		}
		err = nil
		break
	}

	var err2 error
	if err != nil {
		data, _, err2 = utils.HttpProxyGet(api_url, "", nil)
	}
	if err2 != nil {
		msg = tranFailErr
		return
	}
	var bodyData map[string]interface{}
	err = json.Unmarshal(data, &bodyData)
	if err != nil {
		msg = tranFailErr
		return
	}
	var status string
	if s, ok := apiUrl["status"]; ok {
		status = s.(string)
	} else {
		msg = tranFailErr
		return
	}
	statusSlice := strings.Split(status, "<>")
	if ps, ok := bodyData[statusSlice[0]]; ok {
		if ips, ok2 := ps.(float64); ok2 {
			if (strconv.Itoa(int(ips)) != statusSlice[1]) {
				msg = tranFailErr
				return
			}
		} else {
			if ps.(string) != statusSlice[1] {
				msg = tranFailErr
				return
			}
		}
	} else {
		msg = tranFailErr
		return
	}

	var get_url_key string
	if key, ok := apiUrl["key"]; ok {
		get_url_key = key.(string)
	} else {
		msg = reqFailErr
		return
	}

	urlKeySlice := strings.Split(get_url_key, "<>")
	if len(urlKeySlice) == 2 {
		get_url_key = urlKeySlice[1]
		if newBodyData, ok := bodyData[urlKeySlice[0]]; ok {
			bodyData = newBodyData.(map[string]interface{})
		} else {
			msg = reqFailErr
			return
		}
	}

	if s, ok := bodyData[get_url_key]; ok {
		if s == "" {
			msg = reqFailErr
			return
		} else {
			shortUrl = s.(string)
		}
	} else {
		msg = reqFailErr
		return
	}
	return
}
