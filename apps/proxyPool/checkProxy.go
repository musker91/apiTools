package proxyPool

import (
	"apiTools/libs/logger"
	"apiTools/modles"
	"apiTools/utils"
	"fmt"
	"time"
)

// 运行proxy ip校验
// processNum 启动校验协程的数量
func runCheckProxy(processNum int) {
	for i := 0; i < processNum; i++ {
		go checkProxyIP()
	}
	logger.Echo.Info("proxy app: check proxy coroutine is running...")
}

// proxyIP校验
func checkProxyIP() {
	for proxyInfo := range checkProxyJobChan {
		var status bool
		proxyAddr := fmt.Sprintf("%s:%s", proxyInfo.IP, proxyInfo.Port)
		startTime := time.Now()
		// 可用性校验
		switch proxyInfo.Protocol {
		case "http", "https":
			status = utils.CheckProtocolHttp(proxyAddr, checkToUrl)
		}
		useTime := int(time.Since(startTime).Milliseconds())
		if status == false {
			if proxyInfo.FailCount >= 3 {
				_ = modles.DelOneProxyFromDB(proxyInfo.IP)
				continue
			}
			proxyInfo.FailCount++
			proxyInfo.Speed = 0
		} else {
			proxyInfo.Speed = useTime
		}
		proxyInfo.VerifyTime = time.Now()

		// 入库
		checkProxyResultChan <- proxyInfo
	}
}

// 定时从数据库中提取proxy信息检测
func checkDBProxy() {
	logger.Echo.Debugf("proxy app: Extracting database data to check")
	proxyPools, err := modles.ExtractProxyInfo(100)
	if err != nil {
		logger.Echo.Errorf("proxy app: extract proxy info from database(cron name: checkDBProxy), err: %v", err)
		return
	}
	for _, info := range proxyPools {
		proxyInfo := &ProxyInfo{
			IP:         info.IP,
			Port:       info.Port,
			Anonymity:  info.Anonymity,
			Protocol:   info.Protocol,
			Address:    info.Address.String,
			Country:    info.Country.String,
			ISP:        info.ISP.String,
			FailCount:  info.FailCount,
			VerifyTime: info.VerifyTime,
			IsFormDB:   true,
		}
		checkProxyJobChan <- proxyInfo
	}
	logger.Echo.Debugf("proxy app: Extracting database data to check success")
}
