package cmd

import (
	"apiTools/apps/crontab"
	"apiTools/apps/proxyPool"
	"apiTools/apps/token"
	"apiTools/libs/config"
	"apiTools/modles"
	"apiTools/routers"
	"fmt"
)

// 运行所有服务
func runAllProgram() (err error) {
	go runCrontabApp()

	go runProxyPoolApp()

	if err = runHttpServer(); err != nil {
		return
	}
	return
}

// 运行http服务
func runHttpServer() (err error) {
	// 初始化api配置
	err = modles.InitApiConfig()
	if err != nil {
		return
	}

	// 初始化路由引擎
	routers.InitRouter()

	// 启动服务
	err = routers.Router.Run(fmt.Sprintf(":%s", config.GetString("web::port")))
	if err != nil {
		return
	}
	return
}

// 运行http代理池程序
func runProxyPoolApp() (err error) {
	err = proxyPool.RunProxyPoolApp()
	if err != nil {
		err = fmt.Errorf("run proxy pool app fail, Error Msg: %v", err)
		return
	}
	select {}
}

// 运行银行卡信息入库程序任务
func runBankInfoTask() (err error) {
	modles.ReadBankInfoToDB()
	return
}

// 运行计划任务程序
func runCrontabApp() (err error) {
	err = crontab.RunCrontabApp()
	if err != nil {
		err = fmt.Errorf("run crontab apps fail, Error Msg: %v", err)
		return
	}
	select {}
}

// 运行token管理app程序
func runTokenApp() (err error) {
	token.RunTokenApp()
	return
}