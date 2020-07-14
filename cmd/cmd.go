package cmd

import (
	"apiTools/libs/config"
	"apiTools/libs/logger"
	"apiTools/modles"
)

func InitApiTools() (err error) {
	// 解析命令行
	commandStatus := parseCommand()
	if !commandStatus {
		return
	}

	// 初始化配置文件
	err = config.InitConfig()
	if err != nil {
		return
	}

	// 初始化Logger
	err = logger.InitLogger()
	if err != nil {
		return
	}

	// 初始化redis连接池
	err = modles.InitRedis()
	if err != nil {
		return
	}

	// 初始化mysql连接
	err = modles.InitMysql()
	if err != nil {
		return
	}

	// 启动命令行处理
	err = command()
	if err != nil {
		println(err)
	}
	return

}
