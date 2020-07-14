package modles

import (
	"apiTools/utils"
	"github.com/jonsen/gotld"
	"github.com/lionsoul2014/ip2region/binding/golang/ip2region"
	"github.com/pkg/errors"
	"net"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	ipReg = "((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})(\\.((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})){3}"
)

// ipv4表单结构
type Ipv4Form struct {
	Ip string `form:"ip" json:"ip" xml:"ip" binding:"required"`
}

// ipv4数据返回结构
type Ipv4Info struct {
	IP       string `json:"ip"`       // ip地址
	CityId   int64  `json:"cityId"`   // 城市ID号
	Country  string `json:"country"`  // 国家名称
	Region   string `json:"region"`   // 区域号
	Province string `json:"province"` // 省
	City     string `json:"city"`     // 市
	ISP      string `json:"isp"`      // isp厂商
}

var (
	ipv4db *ip2region.Ip2Region
)

// 初始化ipv4db数据库信息
func InitIp4DB() (err error) {
	region, err := ip2region.New(filepath.Join(utils.GetRootPath(), "data/ip/ip2region.db"))
	if err != nil {
		return
	}
	ipv4db = region
	return
}

// ipv4信息查询
func Ipv4Query(ipv4Form Ipv4Form) (Ipv4Info, error) {
	ipInfo := Ipv4Info{}
	ip, err := judType(ipv4Form.Ip)
	if err != nil {
		return ipInfo, err
	}
	queryIpInfo, err := ipv4db.MemorySearch(ip)
	if err != nil {
		return ipInfo, err
	}
	ipInfo.IP = ip
	ipInfo.CityId = queryIpInfo.CityId

	if queryIpInfo.Country == "0" {
		queryIpInfo.Country = ""
	}
	ipInfo.Country = queryIpInfo.Country

	if queryIpInfo.Region == "0" {
		queryIpInfo.Region = ""
	}
	ipInfo.Region = queryIpInfo.Region

	if queryIpInfo.Province == "0" {
		queryIpInfo.Province = ""
	}
	ipInfo.Province = queryIpInfo.Province

	if queryIpInfo.City == "0" {
		queryIpInfo.City = ""
	}
	ipInfo.City = queryIpInfo.City

	if queryIpInfo.ISP == "0" {
		queryIpInfo.ISP = ""
	}
	ipInfo.ISP = queryIpInfo.ISP

	return ipInfo, nil
}

// 判断输入ip的类型
func judType(text string) (ip string, err error) {
	// ip
	matched, err := regexp.MatchString(ipReg, text)
	if err != nil {
		return
	}
	if matched {
		return text, err
	}
	// 域名
	if !strings.HasPrefix(text, "http", ) {
		text = "http://" + text
	}
	parse, err := url.Parse(text)
	if err != nil {
		return
	}
	_, _, err = gotld.GetTld(parse.Hostname())
	if err != nil {
		return
	}

	addrs, err := net.LookupHost(parse.Hostname())
	if err != nil {
		return
	}
	if len(addrs) <= 0 {
		return "", errors.New("dns parse address slice len is zero")
	}
	ip = addrs[0]

	return
}
