package proxyPool

import (
	"apiTools/modles"
	"strings"
	"time"
)

// 定义ip代理存储的结构体
type ProxyInfo struct {
	IP         string    `json:"ip" xml:"ip"`                 // proxy ip
	Port       string    `json:"port" xml:"ip"`               // proxy 端口
	Anonymity  string    `json:"anonymity" xml:"anonymity"`   // 匿名类型(透明/高匿)
	Protocol   string    `json:"protocol" xml:"protocol"`     // 协议类型(http/https)
	Country    string    `json:"country" xml:"country"`       // 所在国家
	Address    string    `json:"address" xml:"address"`       // 所在地区
	ISP        string    `json:"isp" xml:"isp"`               // 运营商
	Speed      int       `json:"speed" xml:"speed"`           // 响应速度
	FailCount  uint      `json:"-" xml:"-"`                   // 校验失败的次数
	VerifyTime time.Time `json:"verifyTime" xml:"verifyTime"` // 最后验证时间
	IsFormDB   bool                                           // 标记是否从此数据库中提取出来的数据
}

const (
	checkToUrl = "www.baidu.com" // proxy校验的地址
)

var (
	// 校验ip管道
	checkProxyJobChan chan *ProxyInfo
	// 校验完成后数据存储管道
	checkProxyResultChan chan *ProxyInfo
)

// 补全proxy ip的信息
func completionProxyInfo(proxyInfo *ProxyInfo) {
	if proxyInfo.Address == "" || proxyInfo.Country == "" || proxyInfo.ISP == "" {
		ipv4Form := modles.Ipv4Form{Ip: proxyInfo.IP}
		ipv4Info, err := modles.Ipv4Query(ipv4Form)
		if err != nil {
			return
		}
		if proxyInfo.Country == "" && ipv4Info.Country != "" {
			proxyInfo.Country = ipv4Info.Country
		}

		if proxyInfo.ISP == "" && ipv4Info.ISP != "" {
			proxyInfo.ISP = ipv4Info.ISP
		}

		if proxyInfo.Address == "" && (ipv4Info.City != "" || ipv4Info.Province != "") {
			proxyInfo.Address = strings.TrimSpace(ipv4Info.Country + " " + ipv4Info.Province + " " + ipv4Info.City)
		}
	}
}
