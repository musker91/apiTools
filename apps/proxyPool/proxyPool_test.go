package proxyPool

import (
	"apiTools/libs/config"
	"apiTools/libs/logger"
	"apiTools/modles"
	"apiTools/utils"
	"testing"
)

func init() {
	err := config.InitConfig()
	if err != nil {
		panic(err)
	}
	err = logger.InitLogger()
	if err != nil {
		panic(err)
	}
	err = modles.InitRedis()
	if err != nil {
		panic(err)
	}
	err = modles.InitMysql()
	if err != nil {
		panic(err)
	}
}

func TestCheckProtocolHttp(t *testing.T) {
	proxyAddr := "123.212.131.133:8008"
	utils.CheckProtocolHttp(proxyAddr, "www.baidu.com")
	t.Logf("%s\n", proxyAddr)
}

func TestFreeipSpider(t *testing.T) {
	spider := withFreeipSpider()
	err := spider().startSpider()
	if err != nil {
		t.Error(err)
	}
	t.Log("sider run success")
}

func TestXiciSpider(t *testing.T) {
	spider := withXiciSpider()
	spider().startSpider()
}