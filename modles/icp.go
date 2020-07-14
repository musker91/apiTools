package modles

import (
	"apiTools/utils"
	"bytes"
	"github.com/antchfx/htmlquery"
	"github.com/jonsen/gotld"
	"github.com/pkg/errors"
	"strings"
)

// ICP备案信息查询
const (
	icpApiUrl          = "http://icp.chinaz.com/"
	domainErr          = "malformed domain name"
	icpInfoNotFount    = "icp information does not exist"
	icpInfoQueryFailed = "icp information query failed"
)

// 请求表单
type ICPForm struct {
	Url string `form:"url" xml:"url" json:"xml" binding:"required"`
}

// 返回表单
type ICPResponse struct {
	OrganizerName          string `json:"organizer_name"`           // 主办单位名称
	OrganizerNature        string `json:"organizer_nature"`         // 主办单位性质
	RecordingLicenseNumber string `json:"recording_license_number"` // 网站备案/许可证号
	SiteName               string `json:"site_name"`                // 网站名称
	SiteIndexUrl           string `json:"site_index_url"`           // 网站首页地址
	ReviewTime             string `json:"review_time"`              // 审核时间
}

// 查询入口
func QueryICPInfo(icpForm *ICPForm) (resp *ICPResponse, msg string, err error) {
	_, domain, err := gotld.GetTld(icpForm.Url)
	if err != nil {
		msg = domainErr
		return
	}
	data, err := getIcpInfo(domain)
	if err != nil {
		msg = icpInfoQueryFailed
		return
	}
	icpResponse, status, err := parseIcpData(data)
	if err != nil {
		msg = icpInfoQueryFailed
		return
	}
	if !status {
		err = errors.New("domain icp info query not found")
		msg = icpInfoNotFount
		return
	}
	resp = icpResponse
	return
}

// 获取数据
func getIcpInfo(domain string) (data []byte, err error) {
	reqUrl := icpApiUrl + domain
	for i := 0; i < 3; i++ {
		proxyIp := GetOneProxyIp("chinazProxyPool")
		data, _, err = utils.HttpProxyGet(reqUrl, proxyIp, nil)
		if err != nil {
			_ = DelOneProxyFromRedis("chinazProxyPool", proxyIp)
			continue
		}
		err = nil
		break
	}
	if err != nil {
		data, _, err = utils.HttpProxyGet(reqUrl, "", nil)
	}
	if err != nil {
		return
	}
	return
}

// 解析数据
func parseIcpData(data []byte) (icpResponse *ICPResponse, status bool, err error) {
	htmlDom, err := htmlquery.Parse(bytes.NewReader(data))
	if err != nil {
		return
	}
	ulNode := htmlquery.Find(htmlDom, "//ul[@id='first']/li")
	if len(ulNode) == 0 {
		return
	}
	icpResponse = &ICPResponse{
		OrganizerName: func() string {
			textBody := htmlquery.InnerText(htmlquery.FindOne(ulNode[0], "/p"))
			return strings.Trim(textBody, "使用高级查询纠正信息")
		}(),
		OrganizerNature: func() string {
			textBody := htmlquery.InnerText(htmlquery.FindOne(ulNode[1], "/p"))
			return textBody
		}(),
		RecordingLicenseNumber: func() string {
			textBody := htmlquery.InnerText(htmlquery.FindOne(ulNode[2], "/p"))
			return strings.Trim(textBody, "查看截图")
		}(),
		SiteName: func() string {
			textBody := htmlquery.InnerText(htmlquery.FindOne(ulNode[3], "/p"))
			return textBody
		}(),
		SiteIndexUrl: func() string {
			textBody := htmlquery.InnerText(htmlquery.FindOne(ulNode[5], "/p"))
			return textBody
		}(),
		ReviewTime: func() string {
			textBody := htmlquery.InnerText(htmlquery.FindOne(ulNode[7], "/p"))
			return textBody
		}(),
	}
	status = true
	return
}
