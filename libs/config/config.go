package config

import (
	"apiTools/utils"
	"encoding/json"
	"fmt"
	"github.com/go-ego/gse"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	AppConfig *AppConfigInfo
	iniFile   *ini.File
	Seg       gse.Segmenter // 分词
)

/* init config */
func InitConfig() (err error) {
	AppConfig = &AppConfigInfo{}
	// 获取配置文件
	configPath := filepath.Join(utils.GetRootPath(), "config", "apitools.ini")
	iniFile, err = ini.Load(configPath)
	if err != nil {
		return
	}
	err = iniFile.MapTo(AppConfig)
	if err != nil {
		return
	}
	err = initRedisProxyPools()
	if err != nil {
		return
	}
	// 配置文件扫描从新配置
	scanConfig(AppConfig)

	// 初始化分词功能, 加载默认字典
	_ = Seg.LoadDict()
	return
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
		reCompile := regexp.MustCompile("[A-Z]+[a-z0-9]+")
		keySlice := reCompile.FindAllString(name, -1)
		for i := 0; i < len(keySlice); i++ {
			keySlice[i] = strings.ToUpper(keySlice[i])
		}
		envName = strings.Join(keySlice, "_")
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
			callBack := reflect.ValueOf(&defaultConfCallBack{})
			callFunc := callBack.MethodByName(defaultValue)
			args := []reflect.Value{reflect.ValueOf("xxx")}
			callValue := callFunc.Call(args)
			if len(callValue) > 0 {
				callResultValue := callValue[0].String()
				setConfigValue(fieldType, fieldVal, callResultValue)
			}
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
		passTag := fieldTp.Tag.Get("pass")
		if passTag != "" {
			continue
		}
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
	AppConfig.proxyPoolApp.RedisProxyPools = redisProxyPools
	return
}


/* get config data */
func readConfigVal(fileConfig *ini.File, keys []string) interface{} {
	sectionStrings := fileConfig.SectionStrings()
	isHas := utils.IsInSelic(keys[0], sectionStrings)
	if !isHas {
		return nil
	}
	section := fileConfig.Section(keys[0])
	keyStrings := section.KeyStrings()
	isHas = utils.IsInSelic(keys[1], keyStrings)
	if !isHas {
		return nil
	}
	value := section.Key(keys[1]).String()
	return value
}

func parseConfig(config interface{}, key string) interface{} {
	tp := reflect.TypeOf(config)
	val := reflect.ValueOf(config)
	if tp.Kind() == reflect.Ptr {
		tp = tp.Elem()
		val = val.Elem()
	}
	for i := 0; i < tp.NumField(); i++ {
		fieldTp := tp.Field(i)
		fieldVal := val.Field(i)
		fieldConfName := fieldTp.Tag.Get("conf")
		if fieldConfName == "" {
			fieldConfName = utils.LowerCase(fieldTp.Name)
		}
		if key == fieldConfName {
			funcTag := fieldTp.Tag.Get("func")
			if funcTag != "" {

			}
			return fieldVal.Interface()
		}
	}
	return nil
}

// get config field value
// Get("web::port")
// Get("config")
func Get(ck string) interface{} {
	if ck == "" {
		return nil
	}
	keys := strings.Split(ck, "::")
	if len(keys) == 0 {
		return nil
	}
	var key string
	var config interface{}
	config = AppConfig
	for i := 0; i < len(keys); i++ {
		key = keys[i]
		config = parseConfig(config, key)
		if config == nil {
			break
		}
	}
	if config == nil && len(keys) == 2 {
		config = readConfigVal(iniFile, keys)
	}
	return config
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
	valueTp := reflect.ValueOf(value).Kind().String()
	isInt  := strings.Index(valueTp, "int")
	if isInt == -1 {
		return 0
	}
	return int(reflect.ValueOf(value).Int())
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

func GetRedisProxyPools() []*RedisProxyPool {
	return AppConfig.proxyPoolApp.RedisProxyPools
}
