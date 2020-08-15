package config

import (
	"time"
)

/*
config struct tag说明:
ini: 配置文件中名称，默认空则同名
conf: 调用配置默认名称，否则以struct字段名第一个字母小写为准
default: 配置默认值, 单字符串表示直接设置值，func:xxx表示调用函数设置值
env: 优先读取环境变量中的值，默认env设置名称是struct field 以大写字母分割, 每个分割的单词全大写 以 _ 相连接
func: 调用获取值时返回的类型
none: 标记字段为none时的值，后续扫描会判断
panic: 当这个值为空时 panic错误消息提示
pass: 忽略初始化扫描
*/

// app configure
type AppConfigInfo struct {
	Service *serviceInfo `ini:"web" conf:"web"`
	Redis   *redisInfo   `ini:"redis"`
	Mysql   *mysqlInfo   `ini:"mysql"`
	Email   *emailInfo   `ini:"email"`
	proxyPoolApp proxyPoolApp `pass:"-"`
}

// config file //

// 服务配置信息
type serviceInfo struct {
	Port                  int    `default:"8091" env:"HTTP_PORT"`
	AppMode               string `default:"development"`
	EnablePprof           bool
	LogLevel              string `default:"info"`
	LogSaveDay            int    `default:"7"`
	LogSplitTime          int    `default:"24"`
	LogOutType            string `default:"json"`
	LogOutPath            string `default:"file"`
	StartTime             string `default:"func:StartTime"`
	EnableIpLimiting      bool
	IpLimitingTimeSeconds int `default:"10"`
	IpLimitingCount       int `default:"8"`
	LiftIpLimiting        int `default:"5"`
}

// redis 配置信息
type redisInfo struct {
	Host     string `panic:"redis host not is empty" env:"REDIS_HOST"`
	Port     int64  `default:"3306" env:"REDIS_PORT"`
	Password string `env:"REDIS_PASSWORD"`
}

// mysql 配置信息
type mysqlInfo struct {
	Host        string `panic:"mysql host not is empty" env:"MYSQL_HOST"`
	Port        int    `default:"3306"  env:"MYSQL_PORT"`
	User        string `panic:"mysql user not is empty"  env:"MYSQL_USER"`
	Password    string `panic:"mysql password not is empty" env:"MYSQL_PASSWORD"`
	DB          string `panic:"mysql db name not is empty" env:"MYSQL_DB" conf:"db"`
	EnableDebug bool
}

// email 配置
type emailInfo struct {
	RecverMail     []string
	SmtpHost       string
	SmtpPort       int
	SenderMail     string
	SenderAuthCode string
}

// other config //

// 代理池在redis中存储结构
type RedisProxyPool struct {
	KeyName  string `json:"keyName"`  // redis 键名称
	CheckUrl string `json:"checkUrl"` // proxy 检测url名称
}

// proxy pool app配置信息
type proxyPoolApp struct {
	RedisProxyPools []*RedisProxyPool
}

// 默认值设置，回调
type defaultConfCallBack struct{}

func (*defaultConfCallBack) StartTime(value string) string {
	_, err := time.Parse("2006/01/02", value)
	if err != nil {
		return time.Now().Format("2006/01/02")
	}
	return value
}

// 调用配置默认回调函数绑定
type getConfigCallBack struct{}

func (*getConfigCallBack) StartTime(value string) time.Time {
	return time.Now()
}
