// 代理池主入口文件
package proxyPool

import "apiTools/modles"

func RunProxyPoolApp() (err error) {
	// 初始化数据管道
	checkProxyJobChan = make(chan *ProxyInfo, 1024)
	checkProxyResultChan = make(chan *ProxyInfo, 1024)

	// 初始化ipv4db数据库信息
	err = modles.InitIp4DB()
	if err != nil {
		return
	}

	// 启动运行爬虫抓取代理
	RunSpiders()

	// 启动校验
	runCheckProxy(16)

	// 启动定时提库校验
	runCrontab()

	// 启动入库
	runSaveToDB(8)

	return
}
