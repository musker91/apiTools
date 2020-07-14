package config

import (
	"apiTools/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-ego/gse"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

// 项目配置信息
type appConfig struct {
	serverConfig `conf:"web"`
	redisConfig  `conf:"redis"`
	mysqlConfig  `conf:"mysql"`
	proxyPoolApp `conf:"proxyPoolApp"`
	emailConfig  `conf:"email"`
}

// 服务配置信息
type serverConfig struct {
	Port                  string    `conf:"port"`                  // http监听端口
	AppMode               string    `conf:"appMode"`               // app运行模式: production, development
	EnablePprof           bool      `conf:"enablePprof"`           // 是否开启pprof
	LogLevel              string    `conf:"logLevel"`              // 日志级别: debug, info, error, warn, panic
	LogSaveDay            uint      `conf:"logSaveDay"`            // 日志文件保留天数
	LogSplitTime          uint      `conf:"logSplitTime"`          // 日志切割时间间隔
	LogOutType            string    `conf:"logOutType"`            // 日志输入类型, json, text
	LogOutPath            string    `conf:"logOutPath"`            // 文件输出位置, file console
	StartTime             time.Time `conf:"startTime"`             // 系统开始运行时间
	EnableIpLimiting      bool      `conf:"enableIpLimiting"`      // 是否开启ip限流
	IpLimitingTimeSeconds uint      `conf:"ipLimitingTimeSeconds"` // IP限流时间段(单位: 秒)
	IpLimitingCount       uint      `conf:"ipLimitingCount"`       // IP限流时间段内请求不能超过的次数
	LiftIpLimiting        uint      `conf:"liftIpLimiting"`        // 解除ip限流的时间(单位: 秒)
}

// redis 配置信息
type redisConfig struct {
	Host     string `conf:"host"`     // redis连接地址
	Port     string `conf:"port"`     // redis连接端口
	Password string `conf:"password"` // redis连接密码
}

// mysql 配置信息
type mysqlConfig struct {
	Host        string `conf:"host"`        // mysql连接地址
	Port        string `conf:"port"`        // mysql连接端口
	User        string `conf:"user"`        // mysql用户名
	Password    string `conf:"password"`    // mysql连接密码
	DB          string `conf:"db"`          // mysql 数据库名称
	EnableDebug bool   `conf:"enableDebug"` // 是否开启sql调试模式
}

// email 配置
type emailConfig struct {
	RecvMail       []string `conf:"recvMail"`
	SmtpHost       string   `conf:"smtpHost"`
	SmtpPort       int      `conf:"smtpPort"`
	SenderMail     string   `conf:"senderMail"`
	SenderAuthCode string   `conf:"senderAuthCode"`
}

// 代理池在redis中存储结构
type RedisProxyPool struct {
	KeyName  string `json:"keyName"`  // redis 键名称
	CheckUrl string `json:"checkUrl"` // proxy 检测url名称
}

// proxy pool app配置信息
type proxyPoolApp struct {
	RedisProxyPools []*RedisProxyPool
}

var (
	appConf appConfig
	Seg     gse.Segmenter // 分词
)

// ----------- 初始化配置 ----------- //
// 初始化配置
func InitConfig() (err error) {
	// 获取配置文件
	configPath := filepath.Join(utils.GetRootPath(), "config", "apitools.ini")
	iniFile, err := ini.Load(configPath)
	if err != nil {
		return
	}
	// 读取web server配置
	err = readServerConfig(iniFile)
	if err != nil {
		return
	}
	// 读取redis配置
	err = readRedisConfig(iniFile)
	if err != nil {
		return
	}

	// 读取mysql配置
	err = readMysqlConfig(iniFile)
	if err != nil {
		return
	}

	// 读取email配置
	err = readEmailConfig(iniFile)
	if err != nil {
		return
	}

	// 读取proxy app 配置
	err = readProxyAppConfig(iniFile)
	if err != nil {
		return
	}
	// 初始化分词功能, 加载默认字典
	_ = Seg.LoadDict()
	return
}

// 读取web server配置
func readServerConfig(iniFile *ini.File) (err error) {
	serverConf := iniFile.Section("web")

	httpPort := serverConf.Key("http_port").String()
	if httpPort == "" {
		httpPort = "8091"
	}
	appConf.serverConfig.Port = httpPort

	runMode := serverConf.Key("app_mode").String()
	if runMode == "" {
		runMode = "development"
	}
	appConf.serverConfig.AppMode = runMode

	enablePprof, err := serverConf.Key("enable_pprof").Bool()
	if err != nil {
		enablePprof = false
	}
	appConf.serverConfig.EnablePprof = enablePprof

	logLevel := serverConf.Key("logLevel").String()
	if logLevel == "" && runMode == "development" {
		logLevel = "debug"
	} else {
		logLevel = "info"
	}
	appConf.serverConfig.LogLevel = logLevel

	logSaveDay, err := serverConf.Key("logSaveDay").Uint()
	if err != nil {
		logSaveDay = 7
	}
	appConf.serverConfig.LogSaveDay = logSaveDay

	logSplitTime, err := serverConf.Key("logSplitTime").Uint()
	if err != nil {
		logSplitTime = 24
	}
	appConf.serverConfig.LogSplitTime = logSplitTime

	logOutType := serverConf.Key("LogOutType").String()
	if logOutType == "" {
		logOutType = "json"
	}
	appConf.serverConfig.LogOutType = logOutType

	logOutPath := serverConf.Key("logOutPath").String()
	if logOutPath == "" {
		logOutPath = "file"
	}
	appConf.serverConfig.LogOutPath = logOutPath

	startTime := serverConf.Key("startTime").String()
	if startTime == "" {
		appConf.serverConfig.StartTime = time.Now()
	} else {
		runTime, err := time.Parse("2006/01/02", startTime)
		if err != nil {
			appConf.serverConfig.StartTime = time.Now()
		} else {
			appConf.serverConfig.StartTime = runTime
		}
	}

	// ip限流
	enableIpLimiting, err := serverConf.Key("enableIpLimiting").Bool()
	if err != nil {
		enableIpLimiting = false
	}
	appConf.serverConfig.EnableIpLimiting = enableIpLimiting
	if enableIpLimiting {
		ipLimitingTimeSeconds, err := serverConf.Key("ipLimitingTimeSeconds").Uint()
		if err != nil || ipLimitingTimeSeconds == 0 {
			ipLimitingTimeSeconds = 10
		}
		appConf.serverConfig.IpLimitingTimeSeconds = ipLimitingTimeSeconds

		ipLimitingCount, err := serverConf.Key("ipLimitingCount").Uint()
		if err != nil || ipLimitingCount == 0 {
			ipLimitingCount = 8
		}
		appConf.serverConfig.IpLimitingCount = ipLimitingCount

		liftIpLimiting, err := serverConf.Key("liftIpLimiting").Uint()
		if err != nil || liftIpLimiting == 0 {
			liftIpLimiting = 5
		}
		appConf.serverConfig.LiftIpLimiting = liftIpLimiting
	}

	return
}

// 读取redis配置
func readRedisConfig(iniFile *ini.File) (err error) {
	redisConf := iniFile.Section("redis")

	host := redisConf.Key("host").String()
	if host == "" {
		return errors.New("config file redis host can not be empty")
	}
	appConf.redisConfig.Host = host

	port := redisConf.Key("port").String()
	if port == "" {
		port = "6379"
	}
	appConf.redisConfig.Port = port

	password := redisConf.Key("password").String()
	appConf.redisConfig.Password = password

	return
}

// 读取mysql配置
func readMysqlConfig(iniFile *ini.File) (err error) {
	mysqlConf := iniFile.Section("mysql")

	host := mysqlConf.Key("host").String()
	if host == "" {
		return errors.New("config file mysql host can not be empty")
	}
	appConf.mysqlConfig.Host = host

	port := mysqlConf.Key("port").String()
	if port == "" {
		port = "3306"
	}
	appConf.mysqlConfig.Port = port

	user := mysqlConf.Key("user").String()
	if user == "" {
		return errors.New("config file mysql connection user can not be empty")
	}
	appConf.mysqlConfig.User = user

	password := mysqlConf.Key("password").String()
	if password == "" {
		return errors.New("config file mysql connection password can not be empty")
	}
	appConf.mysqlConfig.Password = password

	db := mysqlConf.Key("db").String()
	if db == "" {
		return errors.New("config file mysql db name can not be empty")
	}
	appConf.mysqlConfig.DB = db

	enableDebug, err := mysqlConf.Key("enableDebug").Bool()
	if err != nil {
		enableDebug = false
	}
	appConf.mysqlConfig.EnableDebug = enableDebug

	return
}

// 读取email配置
func readEmailConfig(iniFile *ini.File) (err error) {
	emailConf := iniFile.Section("email")

	recvMail := emailConf.Key("recverMail").String()
	if recvMail != "" {
		recvMailSlice := strings.Split(recvMail, ",")
		appConf.emailConfig.RecvMail = recvMailSlice
	}
	smtpHost := emailConf.Key("smtpHost").String()
	appConf.emailConfig.SmtpHost = smtpHost

	smtpPort, _ := emailConf.Key("smtpPort").Int()
	appConf.emailConfig.SmtpPort = smtpPort

	senderMail := emailConf.Key("senderMail").String()
	appConf.emailConfig.SenderMail = senderMail

	senderAuthCode := emailConf.Key("senderAuthCode").String()
	appConf.emailConfig.SenderAuthCode = senderAuthCode

	return
}

// 读取proxy app配置
func readProxyAppConfig(iniFile *ini.File) (err error) {
	err = initRedisProxyPools()
	if err != nil {
		return
	}
	return
}

// 初始化redis proxy ip 池配置
func initRedisProxyPools() (err error) {
	redisProxyPoolsFile, err := ioutil.ReadFile(filepath.Join(utils.GetRootPath(), "config", "redisProxyPools.json"))
	if err != nil {
		return
	}
	redisProxyPools := make([]*RedisProxyPool, 0, 10)
	err = json.Unmarshal(redisProxyPoolsFile, &redisProxyPools)
	if err != nil {
		err = fmt.Errorf("config parse redisProxyPools.json fail, err: %v", err)
		return
	}
	appConf.proxyPoolApp.RedisProxyPools = redisProxyPools
	return
}

// ----------- 获取数据 ----------- //

// 获取指定字段值
// Get("web::port")
// Get("appname")
func Get(ck string) interface{} {
	if ck == "" {
		return nil
	}
	keys := strings.Split(ck, "::")
	if len(keys) == 0 {
		return nil
	}
	appType := reflect.TypeOf(appConf)
	appVal := reflect.ValueOf(appConf)
	for i := 0; i < appType.NumField(); i++ {
		appTypeFiled := appType.Field(i)
		regionTag := appTypeFiled.Tag.Get("conf")
		if regionTag == keys[0] {
			if len(keys) > 1 {
				regionType := appVal.Field(i).Type()
				regionVal := appVal.Field(i)
				for j := 0; j < regionType.NumField(); j++ {
					regionTypeFiled := regionType.Field(j)
					confTag := regionTypeFiled.Tag.Get("conf")
					if confTag == keys[1] {
						val := regionVal.FieldByName(regionTypeFiled.Name).Interface()
						return val
					}
				}
			} else {
				val := appVal.FieldByName(appTypeFiled.Name).Interface()
				return val
			}
		}
	}
	return nil
}

func GetString(ck string) (val string) {
	value := Get(ck)
	if value == nil {
		return
	}
	val, ok := value.(string)
	if !ok {
		return ""
	}
	return
}

func GetInt(ck string) (int) {
	value := Get(ck)
	if value == nil {
		return 0
	}
	if v01, ok01 := value.(uint); ok01 {
		return int(v01)
	}
	if v02, ok02 := value.(int); ok02 {
		return v02
	}
	return 0
}

func GetBool(ck string) (bool) {
	value := Get(ck)
	if value == nil {
		return false
	}
	if v, ok := value.(bool); ok {
		return v
	}

	return false
}

func GetStrings(ck string) (reslut []string) {
	value := Get(ck)
	if value == nil {
		return
	}
	if v, ok := value.([]string); ok {
		return v
	}
	return
}

func GetRedisProxyPools() []*RedisProxyPool {
	return appConf.proxyPoolApp.RedisProxyPools
}
