package modles

import (
	"apiTools/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jonsen/gotld"
	"io"
	"io/ioutil"
	"net"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var (
	whoisServers map[string][]string
)

// whois信息接收的表单数据结构
type WhoisForm struct {
	Domain   string `form:"domain" json:"domain" xml:"domain" binding:"required"` // 域名
	OutType  string `form:"type" json:"type" xml:"type"`                          // 返回数据类型, json, text
	Standard bool   `form:"standard" json:"standard" xml:"standard"`              // 是否按照标准的固定字段输出成
}

// whois查询后返回信息数据结构
type WhoisInfo struct {
	WhoisForm     // whois查询数据表单数据
	Domain string // 解析后的域名
	Tld    string // 解析后域名tld信息
	// 0 获取到域名whois信息
	// 1 域名解析失败
	// 2 域名未注册
	// 3 暂不支持此域名后缀查询
	// 4 域名查询失败
	// 5 请求数据错误
	Status   uint                   // 域名 查询状态
	TextInfo string                 // 域名whois查询后文本信息
	JsonInfo map[string]interface{} // 域名whois查询后转换为map的信息
}

// 获取whois服务器列表
func InitWhoisServers() (err error) {
	whoisServers = make(map[string][]string)
	fileBytes, err := ioutil.ReadFile(filepath.Join(utils.GetRootPath(), "data/whois", "whois.servers.json"))
	if err != nil {
		err = fmt.Errorf("read whois servers json file fail, err: %v\n", err)
		return
	}
	err = json.Unmarshal(fileBytes, &whoisServers)
	if err != nil {
		err = fmt.Errorf("unmarshal whois server json file fail, err: %v\n", err)
		return
	}
	return
}

// 返回域名whois信息,文本数据信息
func QueryWhoisInfo(whoisForm *WhoisForm) (whoisInfo *WhoisInfo, err error) {
	whoisInfo = &WhoisInfo{
		WhoisForm: *whoisForm,
	}
	whoisInfo, err = whoisInfo.getWhoisInfo()
	if err != nil {
		return
	}
	whoisInfo, err = whoisInfo.matchWhois()
	if err != nil {
		return
	}
	return
}

// 返回域名whois信息,map数据信息
func QueryWhoisInfoToJson(whoisForm *WhoisForm) (whoisInfo *WhoisInfo, err error) {
	whoisInfo, err = QueryWhoisInfo(whoisForm)
	if err != nil {
		return
	}
	whoisInfo, err = whoisInfo.textToJson()
	if whoisInfo.Status != 0 {
		return
	}
	standardJsonData := make(map[string]interface{})
	// 组装标准固定字段输出
	if whoisForm.Standard {
		standardJsonData["domainName"] = whoisInfo.JsonInfo["domainName"]
		standardJsonData["domainStatus"] = whoisInfo.JsonInfo["domainStatus"]
		standardJsonData["dnsNameServer"] = whoisInfo.JsonInfo["nameServer"]
		// 注册时间
		if registrationTime, ok := whoisInfo.JsonInfo["registrationTime"]; ok {
			standardJsonData["registrationTime"] = registrationTime
		} else if registrationTime2, ok := whoisInfo.JsonInfo["creationDate"]; ok {
			standardJsonData["registrationTime"] = registrationTime2
		} else {
			standardJsonData["registrationTime"] = ""
		}

		// 过期时间
		if registryExpiryDate, ok := whoisInfo.JsonInfo["registryExpiryDate"]; ok {
			standardJsonData["expirationTime"] = registryExpiryDate
		} else if registryExpiryDate2, ok := whoisInfo.JsonInfo["expirationTime"]; ok {
			standardJsonData["expirationTime"] = registryExpiryDate2
		} else {
			standardJsonData["expirationTime"] = ""
		}

		// 更新时间
		if updatedDate, ok := whoisInfo.JsonInfo["updatedDate"]; ok {
			standardJsonData["updatedDate"] = updatedDate
		} else {
			standardJsonData["updatedDate"] = ""
		}

		// whois server
		if registrarWHOISServer, ok := whoisInfo.JsonInfo["registrarWHOISServer"]; ok {
			standardJsonData["registrarWHOISServer"] = registrarWHOISServer
		} else {
			standardJsonData["registrarWHOISServer"] = ""
		}

		// 注册商和联系人
		if registrar, ok := whoisInfo.JsonInfo["registrar"]; ok {
			standardJsonData["registrar"] = registrar
		} else if registrar2, ok := whoisInfo.JsonInfo["sponsoringRegistrar"]; ok {
			standardJsonData["registrar"] = registrar2
		} else {
			standardJsonData["registrar"] = ""
		}

		// 联系人
		if registrant, ok := whoisInfo.JsonInfo["registrant"]; ok {
			standardJsonData["registrant"] = registrant
		} else if registrant2, ok := whoisInfo.JsonInfo["registrantOrganization"]; ok {
			standardJsonData["registrant"] = registrant2
		} else {
			standardJsonData["registrant"] = ""
		}

		// 联系邮箱
		if contactEmail, ok := whoisInfo.JsonInfo["registrarAbuseContactEmail"]; ok {
			standardJsonData["contactEmail"] = contactEmail
		} else if contactEmail2, ok := whoisInfo.JsonInfo["registrantContactEmail"]; ok {
			standardJsonData["contactEmail"] = contactEmail2
		} else {
			standardJsonData["contactEmail"] = ""
		}

		// 联系电话
		if contactPhone, ok := whoisInfo.JsonInfo["registrarAbuseContactPhone"]; ok {
			standardJsonData["contactPhone"] = contactPhone
		} else if contactPhone2, ok := whoisInfo.JsonInfo["registrantContactPhone"]; ok {
			standardJsonData["contactPhone"] = contactPhone2
		} else {
			standardJsonData["contactPhone"] = ""
		}
		whoisInfo.JsonInfo = standardJsonData
	}
	return
}

// 获取whois信息
func (whoisInfo *WhoisInfo) getWhoisInfo() (*WhoisInfo, error) {
	// 赋值
	// 获取域名信息
	tld, domain, err := gotld.GetTld(whoisInfo.WhoisForm.Domain)
	if err != nil {
		whoisInfo.Status = 1
		return whoisInfo, err
	}
	whoisInfo.Domain = domain
	whoisInfo.Tld = tld.Tld
	// 获取域名服务器列表
	domainSuffix := tld.Tld
	servers, ok := whoisServers[domainSuffix]
	if !ok || len(servers) == 0 {
		whoisInfo.Status = 3
		err = fmt.Errorf("%s get whois server faild, because is empty\n", domain)
		return whoisInfo, err
	}
	// 获取域名whois信息
	ctx, cancel := context.WithCancel(context.Background())
	dataChan := make(chan string)
	wgChan := make(chan int, len(servers))
	defer close(dataChan)
	defer close(wgChan)
	for _, server := range servers {
		go connWhoisServer(domain, server, dataChan, ctx, wgChan)
	}

	// 进程个数统计
	var wgCount int
DoneLabel:
	for {
		select {
		case whoisInfo.TextInfo = <-dataChan:
			if len(whoisInfo.TextInfo) < 3 {
				whoisInfo.Status = 4
			} else {
				whoisInfo.Status = 0
			}
			cancel()
			break DoneLabel
		case <-ctx.Done():
			break DoneLabel
		case <-wgChan:
			wgCount++
			if wgCount == len(servers) {
				whoisInfo.Status = 4
				break DoneLabel
			}
		}
	}
	return whoisInfo, err
}

// 连接whois服务器获取whois信息
func connWhoisServer(domain string, server string, dataChan chan string, ctx context.Context, wgChan chan int) {
	// 捕获结协程中的异常
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	// 协程结束
	defer func() {
		wgChan <- 1
	}()
	// 连接whois服务器获取数据
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:43", server), time.Second*10)
	if err != nil {
		return
	}
	defer conn.Close()
	_, err = conn.Write([]byte(domain))
	if err != nil {
		return
	}
	_, err = conn.Write([]byte("\r\n"))
	if err != nil {
		return
	}
	var content []byte
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err == io.EOF || err != nil {
			content = append(content, buf[:n]...)
			break
		} else {
			content = append(content, buf[:n]...)
		}
	}

	select {
	// 是否已经关闭
	case <-ctx.Done():
		break
	// 写入数据到chan
	case dataChan <- string(content):
		break
	}
}

// 匹配查询后的域名信息
func (whoisInfo *WhoisInfo) matchWhois() (*WhoisInfo, error) {
	if whoisInfo.Status != 0 {
		return whoisInfo, fmt.Errorf("match domain [%s] fail", whoisInfo.Domain)
	}
	// 处理文本格式whois信息换行符
	textInfoSlice := strings.Split(whoisInfo.TextInfo, "\n")
	if len(textInfoSlice) <= 1 {
		textInfoSlice = strings.Split(whoisInfo.TextInfo, "\r\n")
	}
	newTextInfo := ""
	for _, line := range textInfoSlice {
		newTextInfo += line + "\n"
	}
	whoisInfo.TextInfo = newTextInfo

	// 匹配域名是否查询成功
	matchedUpper, _ := regexp.Match(fmt.Sprintf("Domain Name: %s",
		strings.ToUpper(whoisInfo.Domain)), []byte(whoisInfo.TextInfo))
	matchedLower, _ := regexp.Match(fmt.Sprintf("Domain Name: %s",
		strings.ToLower(whoisInfo.Domain)), []byte(whoisInfo.TextInfo))
	if matchedUpper || matchedLower {
		whoisInfo.Status = 0
		return whoisInfo, nil
	}

	whoisInfo.Status = 2
	return whoisInfo, nil
}

// 查询到的whois信息转换为文本信息
func (whoisInfo *WhoisInfo) textToJson() (*WhoisInfo, error) {
	// 域名状态不正常直接返回
	if whoisInfo.Status != 0 {
		return whoisInfo, errors.New("domain query status is fail\n")
	}
	// 统计key值
	keyCount := make(map[string]int)
	textInfoSlice := strings.Split(whoisInfo.TextInfo, "\n")
	for _, line := range textInfoSlice {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "DNSSEC") {
			break
		}
		lineSlice := strings.Split(line, ":")
		key := lineSlice[0]
		c, ok := keyCount[key]
		if ok {
			keyCount[key] = c + 1
		} else {
			keyCount[key] = 1
		}
	}
	// 创建存储whios json数据的map
	whoisJsonInfo := make(map[string]interface{})
	// 正则匹配查找
	for key, count := range keyCount {
		// 生成map中的key
		keySlice := strings.Split(key, " ")
		newKey := strings.ToLower(keySlice[0])
		for _, ks := range keySlice[1:] {
			newKey += ks
		}
		// 匹配数据
		if count > 1 {
			//`Domain Name: (.+)`
			whoisJsonInfo[newKey] = matchText(fmt.Sprintf("%s: (.+)", key), whoisInfo.TextInfo, "slice")
		} else {
			whoisJsonInfo[newKey] = matchText(fmt.Sprintf("%s: (.+)", key), whoisInfo.TextInfo, "str")
		}
	}
	// 赋值
	whoisInfo.JsonInfo = whoisJsonInfo

	return whoisInfo, nil
}

// ty返回数据类型, str, slice
func matchText(pattern string, text string, ty string) (data interface{}) {
	re := regexp.MustCompile(pattern)
	submatch := re.FindAllStringSubmatch(text, -1)
	if len(submatch) == 0 {
		if ty == "str" {
			data = ""
		} else {
			data = make([]string, 0)
		}
	} else {
		if ty == "str" {
			data = strings.TrimSpace(submatch[0][1])
		} else {
			rslice := make([]string, 0, len(submatch))
			for _, match := range submatch {
				rslice = append(rslice, strings.TrimSpace(match[1]))
			}
			data = rslice
		}
	}
	return
}
