package modles

import (
	"apiTools/libs/config"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

var (
	RedisPool *redis.Pool
	JsonData  interface{}
	SqlConn   *gorm.DB
)

func InitRedis() (err error) {
	RedisPool = &redis.Pool{
		Dial: func() (conn redis.Conn, err error) {
			conn, err = redis.Dial("tcp", fmt.Sprintf("%s:%d",
				config.GetString("redis::host"), config.GetInt("redis::port")))
			if err != nil {
				return nil, err
			}
			if config.GetString("redis::password") != "" {
				if _, err := conn.Do("AUTH", config.GetString("redis::password")); err != nil {
					conn.Close()
					return nil, err
				}
			}
			return
		},
		MaxIdle:     16,              // 最大空闲连接数
		MaxActive:   32,              // 最大活跃连接数
		IdleTimeout: time.Second * 3, // 最大空闲超时时间
	}
	pool := RedisPool.Get()
	defer pool.Close()
	if err := pool.Err(); err != nil {
		return fmt.Errorf("init redis failed, err: %v", err)
	}
	return
}

func InitMysql() (err error) {
	// CREATE DATABASE `apitools` CHARACTER SET utf8 COLLATE utf8_general_ci;
	// ALTER user 'apitools'@'localhost' IDENTIFIED BY 'apitools#.*'

	db, err := gorm.Open("mysql",
		fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			config.GetString("mysql::user"),
			config.GetString("mysql::password"),
			config.GetString("mysql::host"),
			config.GetInt("mysql::port"),
			config.GetString("mysql::db"),
		))
	if err != nil {
		return
	}
	SqlConn = db

	if config.GetBool("mysql::enableDebug") {
		db.LogMode(true)
	}
	// 单数表名
	SqlConn.SingularTable(true)

	// 自动映表
	SqlConn.AutoMigrate(&ProxyPool{}, &BankBinInfo{}, &UserToken{}, &UserVisitAppli{})

	return
}

// 关闭IO流
func CloseIO() {
	defer func() {
		recover()
	}()

	// 关闭redis io
	RedisPool.Close()

	// 关闭mysql io
	_ = SqlConn.Close()
}

// 初始化api配置
func InitApiConfig() (err error) {
	// 初始化api docs json数据
	err = InitApiDocsJsonData()
	if err != nil {
		return
	}

	// 初始化whois服务器列表
	err = InitWhoisServers()
	if err != nil {
		return
	}

	// 初始化ipv4db数据库信息
	err = InitIp4DB()
	if err != nil {
		return
	}

	// 初始化短链接转换api
	err = InitShortData()
	if err != nil {
		return
	}
	return
}
