package config

import (
	"apiTools/utils"
	"fmt"
	"gopkg.in/ini.v1"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

/*
config struct tag说明:
ini: 配置文件中名称，默认空则同名
conf: 调用配置默认名称，否则以struct字段名第一个字母小写为准
default: 配置默认值, 单字符串表示直接设置值，func:xxx表示调用函数设置值
env: 优先读取环境变量中的值，默认读取struct字段名全大写值
func: 调用获取值时返回的类型
none: 标记字段为none时的值，后续扫描会判断
panic: 当这个值为空时 panic错误消息提示
*/

// app configure
type AppConfigInfo struct {
	Service *serviceInfo `ini:"web" conf:"web"`
	Redis   *redisInfo   `ini:"redis"`
	Mysql   *mysqlInfo   `ini:"mysql"`
	Email   *emailInfo   `ini:"email"`
}

// config file //

// 服务配置信息
type serviceInfo struct {
	Port                  int    `default:"8091"`
	AppMode               string `default:"development"`
	EnablePprof           bool
	LogLevel              string `default:"info"`
	LogSaveDay            int    `default:"7"`
	LogSplitTime          int    `default:"24"`
	LogOutType            string `default:"json"`
	LogOutPath            string `default:""`
	StartTime             string `default:"func:startTime"`
	EnableIpLimiting      bool
	IpLimitingTimeSeconds int `default:"10"`
	IpLimitingCount       int `default:"8"`
	LiftIpLimiting        int `default:"5"`
}

// redis 配置信息
type redisInfo struct {
	Host     string
	Port     int64 `default:"3306"`
	Password string
}

// mysql 配置信息
type mysqlInfo struct {
	Host        string
	Port        int `default:"3306"`
	User        string
	Password    string
	DB          string
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

var (
	AppConfig *AppConfigInfo
	//Seg       gse.Segmenter // 分词

	defaultCallBack confCallBack
)

func InitConfig() (err error) {
	AppConfig = &AppConfigInfo{}
	// 获取配置文件
	configPath := filepath.Join(utils.GetRootPath(), "config", "apitools.ini")
	iniFile, err := ini.Load(configPath)
	if err != nil {
		return
	}
	err = iniFile.MapTo(AppConfig)
	if err != nil {
		return
	}
	// 配置文件扫描从新配置
	scanConfig(AppConfig)

	// 初始化分词功能, 加载默认字典
	//_ = Seg.LoadDict()
	return
}

type confCallBack struct{}

func (*confCallBack) startTime(value string) time.Time {
	if value == "" {

	}
	parse, _ := time.Parse("2006/01/02", value)
	return parse
}

func setConfigValue(fieldType *reflect.StructField, fieldVal *reflect.Value, value interface{}) {
	confValue := value.(string)
	switch fieldType.Type.Kind() {
	case reflect.String:
		fieldVal.SetString(confValue)
	case reflect.Int:
		v, _ := strconv.Atoi(confValue)
		fieldVal.SetInt(int64(v))
	}
}

func updateConfigValue(fieldType *reflect.StructField, fieldVal *reflect.Value) {
	name := fieldType.Name
	value := fieldVal.Interface()
	// env name
	envName := fieldType.Tag.Get("env")
	if envName == "" {
		envName = strings.ToUpper(name)
	}
	// default value
	defaultType := "" // "","func"
	defaultValue := fieldType.Tag.Get("default")
	if strings.HasPrefix(defaultValue, "func:") {
		defaultType = "func"
		defaultValue = strings.Split(defaultValue, "func:")[1]
	}
	// none value
	noneValue := fieldType.Tag.Get("none")
	if noneValue == "" {
		switch fieldType.Type.Kind() {
		case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint64:
			noneValue = "0"
		case reflect.Bool:
			noneValue = "false"
		}
	}

	// 优先级  env, source, default, none
	envValue := os.Getenv(envName)
	if envValue != "-" {
		if envValue != "" {
			setConfigValue(fieldType, fieldVal, envValue)
			return
		}
	}
	newVal := fmt.Sprintf("%v", value)
	if newVal != strings.TrimSpace(noneValue) {
		return
	}

	if defaultValue != "" {
		if defaultType == "func" {
			callBack := reflect.ValueOf(defaultCallBack)
			callFunc := callBack.MethodByName(defaultValue)
			callValue := callFunc.Call([]reflect.Value{})[0].Interface()
			setConfigValue(fieldType, fieldVal, callValue)
		} else {
			setConfigValue(fieldType, fieldVal, defaultValue)
		}
	}
	// error panic
	panicValue := fieldType.Tag.Get("panic")
	if panicValue != "" {
		if fieldVal.String() == "" && defaultValue == "" && envValue == "" {

			panic(panicValue)
		}
	}
}

func scanConfig(config interface{}) {
	tp := reflect.TypeOf(config)
	val := reflect.ValueOf(config)
	if tp.Kind() == reflect.Ptr {
		tp = tp.Elem()
		val = val.Elem()
	}
	for i := 0; i < tp.NumField(); i++ {
		fieldTp := tp.Field(i)
		fieldVal := val.Field(i)
		var tyKind reflect.Kind
		if fieldTp.Type.Kind() == reflect.Ptr {
			tyKind = fieldTp.Type.Elem().Kind()
		}
		if tyKind == reflect.Struct {
			scanConfig(fieldVal.Interface())
		} else {
			updateConfigValue(&fieldTp, &fieldVal)
		}
	}
}

/*
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
*/

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
	appType := reflect.TypeOf(AppConfig)
	appVal := reflect.ValueOf(AppConfig)
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

func GetInt(ck string) int {
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

func GetBool(ck string) bool {
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

//func GetRedisProxyPools() []*RedisProxyPool {
//	return AppConfig.proxyPoolApp.RedisProxyPools
//}
