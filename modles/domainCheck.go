package modles

import (
	"apiTools/utils"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"strings"
	"time"
)

// 检测域名在微信或者qq是否安全
/*
whitetype
1 未知
2 危险
3 安全
*/

const (
	// 第一个参数域名, 第二个参数是时间戳
	ymjcApiUrl        = "https://cgi.urlsec.qq.com/index.php?m=check&a=check&callback=jQuery1113049533065324068226_1585819829729&url=%s&_=%v"
	domainCheckFailed = "domain check failed"
)

var (
	domainCheckHeader = map[string]string{
		"Cookie": " tvfe_boss_uuid=b23cdf7df1cf0cc2; pgv_pvid=6034284609; uin=; skey=; pgv_pvi=3848481792; pgv_si=s9630232576; pgv_info=ssid=s1315500590",
	}
)

type DomainCheckForm struct {
	Url string `form:"url" json:"url" xml:"url" binding:"required"`
}

type DomainCheckResponse struct {
	Domain string `json:"domain"`
	Wx     string `json:"wx"` // ok,abnormal
	QQ     string `json:"qq"`
}

func QueryDomainStatus(form *DomainCheckForm) (resp *DomainCheckResponse, msg string, err error) {
	data, err := getDomainData(form.Url)
	if err != nil {
		msg = domainCheckFailed
		return
	}
	response, err := parseDomainData(data)
	if err != nil {
		msg = domainCheckFailed
		return
	}
	resp = response
	resp.Domain = form.Url
	return
}

func getDomainData(domain string) (data []byte, err error) {
	domainCheckHeader["Referer"] = fmt.Sprintf("https://urlsec.qq.com/check.html?url=%s", domain)
	apiUrl := fmt.Sprintf(ymjcApiUrl, domain, time.Now().Unix())
	data, _, err = utils.HttpProxyGet(apiUrl, "", domainCheckHeader)
	if err != nil {
		return
	}
	return
}

func parseDomainData(data []byte) (resp *DomainCheckResponse, err error) {
	newStr := strings.Replace(string(data), "jQuery1113049533065324068226_1585819829729(", "", -1)
	newStr = strings.Replace(newStr, ")", "", -1)
	var d interface{}
	err = json.Unmarshal([]byte(newStr), &d)
	if err != nil {
		return
	}
	var dm map[string]interface{}
	if tdm, ok := d.(map[string]interface{}); ok {
		dm = tdm
	} else {
		err = errors.New("data structure error")
		return
	}
	// 判断数据状态
	if s, ok := dm["reCode"]; !ok {
		err = errors.New("get data body has error")
		return
	} else {
		if s2, ok2 := s.(float64); ok2 {
			if int(s2) != 0 {
				err = errors.New("get data body has failed")
				return
			}
		} else {
			err = errors.New("get data body has error")
			return
		}
	}

	// 取数据
	var dm2 map[string]interface{}
	if tdm2, ok := dm["data"].(map[string]interface{}); ok {
		dm2 = tdm2
	} else {
		err = errors.New("data structure error")
		return
	}

	var results map[string]interface{}
	if rlt, ok := dm2["results"].(map[string]interface{}); ok {
		results = rlt
	} else {
		err = errors.New("data structure error")
		return
	}

	var statusInter interface{}
	if sc, ok := results["whitetype"]; ok {
		statusInter = sc
	} else {
		err = errors.New("get data body has error")
		return
	}

	var statusCode int
	if scode, ok := statusInter.(float64); ok {
		statusCode = int(scode)
	} else {
		err = errors.New("get data body has error")
		return
	}

	var statusStr string
	switch statusCode {
	case 2:
		statusStr = "danger"
	case 3:
		statusStr = "safe"
	default:
		statusStr = "unknown"
	}
	resp = &DomainCheckResponse{
		Wx: statusStr,
		QQ: statusStr,
	}
	return
}
