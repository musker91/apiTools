package proxyPool

import (
	"apiTools/libs/logger"
	"apiTools/modles"
	"database/sql"
	"github.com/sirupsen/logrus"
)

// 运行入库进程
// processNum运行多少个协程
func runSaveToDB(processNum int) {
	for i := 0; i < processNum; i++ {
		go saveToDB()
	}
	logger.Echo.Info("proxy app: proxy saveToDB coroutine is running...")
}

func saveToDB() {
	for proxyInfo := range checkProxyResultChan {
		proxyInfoObj := &modles.ProxyPool{
			IP:         proxyInfo.IP,
			Port:       proxyInfo.Port,
			Anonymity:  proxyInfo.Anonymity,
			Protocol:   proxyInfo.Protocol,
			Address:    sql.NullString{String: proxyInfo.Address, Valid: true},
			Country:    sql.NullString{String: proxyInfo.Country, Valid: true},
			ISP:        sql.NullString{String: proxyInfo.ISP, Valid: true},
			Speed:      sql.NullInt64{Int64: int64(proxyInfo.Speed), Valid: true},
			FailCount:  proxyInfo.FailCount,
			VerifyTime: proxyInfo.VerifyTime,
		}
		err := modles.InsertProxyInfo(proxyInfoObj, proxyInfo.IsFormDB)
		if err != nil {
			logger.Echo.WithFields(
				logrus.Fields{"proxyInfo:": proxyInfoObj},
			).Error("save proxy info to db fail")
		}
	}
}
