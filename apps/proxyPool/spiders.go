package proxyPool

import (
	"apiTools/libs/config"
	"apiTools/libs/logger"
	"apiTools/modles"
	"apiTools/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/pkg/errors"
	"strings"
	"sync"
)

/*
 代理ip网站:
	1. 高可用全球免费代理IP库 https://ip.jiangxianli.com/ (已实现)
	2. 西刺免费代理IP https://www.xicidaili.com	(已实现-暂时无法访问)
	3. 西拉免费代理IP http://www.xiladaili.com
	4. 89免费代理IP	http://www.89ip.cn
	5. 快代理		https://www.kuaidaili.com
	6. ip3366云代理	http://www.ip3366.net
	7. 89免费代理	http://www.89ip.cn
	8. Emailtry		http://emailtry.com
	9. 青花代理ip	http://www.qinghuadaili.com
	10. 开心代理		http://www.kxdaili.com
	11. 泥马IP代理	http://www.nimadaili.com
	12. 极速数据		http://www.superfastip.com (已实现-页面dom改版，暂时无法获取)
	13. github proxy_pool   http://127.0.0.1:15010 (实现中)
*/

// 所有定义爬虫

/*
自定义爬虫:
	1. 实现spiderInterface接口中的方法
	2. 创建一个返回spiderOption变量类型的结构体
	3. 在runSpiders函数中，注册爬虫
*/

const (
	spiderCatchTaskSpec     = "0 0 */1 * * * "
	extractProxyToRedisSpec = "*/10 * * * * ?"
)

// 定义爬虫的接口
type spiderInterface interface {
	startSpider() error
	redisProxyPoolInfo() (keyName, checkUrl string)
	spiderTaskSpec() string
	spiderName() string
}

// 定义启动爬虫结构
type spiderOption func() (spider spiderInterface)

// 定义爬虫基本的结构体
type spiderBaseInfo struct {
	name     string // 爬虫的名称
	keyName  string // redis中存储proxy pool的key
	checkUrl string // proxy pool 校验的url地址
	spec     string // 爬虫定时任务时间
}

// 运行爬虫
func RunSpiders() {
	go recordSpiders()
	spiderCron := utils.NewWithCron()
	_, err := spiderCron.AddFunc(spiderCatchTaskSpec, func() {
		recordSpiders()
	})
	if err != nil {
		logger.Echo.Errorf("run spiders cron task fail, err: %v", err)
		return
	}
	spiderCron.Start()
}

// 爬虫程序注册声明
func recordSpiders() {
	// withFreeipSpider()
	// withXiciSpider()
	// withSuperipSpider()
	// withLocalProxySpider()
	registerSpiders(withFreeipSpider())
}

func startSingleSpider(spiderObj spiderInterface) {
	// redis proxy pool cron
	keyName, checkUrl := spiderObj.redisProxyPoolInfo()
	scron := utils.NewWithCron()
	defer scron.Stop()
	_, err := scron.AddFunc(extractProxyToRedisSpec, func() {
		logger.Echo.Debugf("run extracting proxy to redis task(keyName: %s, checkUrl: %s)", keyName, checkUrl)
		modles.ExtractProxyToRedis(keyName, checkUrl)
	})
	if err != nil {
		logger.Echo.Errorf("start %s redis proxy pool cron fail, err: %v", keyName, err)
		return
	}
	scron.Start()
	logger.Echo.Infof("start running %s spider", spiderObj.spiderName())
	err = spiderObj.startSpider()
	if err != nil {
		logger.Echo.Error(err)
		return
	}
}

// 注册爬虫
func registerSpiders(spiders ...spiderOption) {
	for _, spider := range spiders {
		s := spider()
		logger.Echo.Infof("start running %s spiders", s.spiderName())
		startSingleSpider(s)
	}
}

// ----------------------定义爬虫---------------------- //

/*
  高可用全球免费代理IP库
  https://www.freeip.top
*/
type freeipSpider struct {
	spiderBaseInfo
	url       string
	pages     int
	failCount int
}

func (s *freeipSpider) startSpider() (err error) {
	err = s.getProxyData()
	if err != nil {
		return
	}
	return
}

func (s *freeipSpider) spiderName() string {
	return s.spiderBaseInfo.name
}

func (s *freeipSpider) spiderTaskSpec() string {
	return s.spec
}

func (s *freeipSpider) redisProxyPoolInfo() (keyName, checkUrl string) {
	return s.spiderBaseInfo.keyName, s.spiderBaseInfo.checkUrl
}

func (s *freeipSpider) getProxyData() (err error) {
	var failCount = 0
	for i := 1; i <= s.pages; i++ {
		if failCount > s.failCount {
			logger.Echo.Errorf("%s site may have been hung up", s.name)
			break
		}
		logger.Echo.Debugf("start running %s spider,count page is: %d, page is: %d", s.name, s.pages, i)
		proxyIp := modles.GetOneProxyIp(s.spiderBaseInfo.keyName)
		url := fmt.Sprintf("%s%d", s.url, i)
		data, _, err := utils.HttpProxyGet(url, proxyIp, nil)
		if err != nil {
			logger.Echo.Debugf("run %s spider, page is: %d, err: %v", s.name, i, err)
			failCount++
			i--
			_ = modles.DelOneProxyFromRedis(s.spiderBaseInfo.keyName, proxyIp)
			utils.TimerUtil(2)
			continue
		}
		lastPage, err := s.parseData(data)
		if err != nil {
			logger.Echo.Debugf("run %s spider, page is: %d, err: %v", s.name, i, err)
			failCount++
			i--
			utils.TimerUtil(3)
			continue
		}
		if lastPage != 0 {
			if i >= lastPage {
				break
			} else {
				s.pages = lastPage
				logger.Echo.Debugf("modify page, page count is: %d", lastPage)
			}
		}
		failCount = 0
		logger.Echo.Debugf("run %s spider,count page is: %d, get page %d data success", s.name, s.pages, i)
		utils.TimerUtil(3)
	}
	logger.Echo.Debugf("run %s spider end", s.name)
	return
}

func (s *freeipSpider) parseData(data []byte) (lastPage int, err error) {
	var response map[string]interface{}
	err = json.Unmarshal(data, &response)
	if err != nil {
		return
	}
	respData, ok := response["data"].(map[string]interface{})
	if !ok {
		err = errors.New("parse response json data fail")
		return
	}
	respLastPage, ok := respData["last_page"].(float64)
	if ok {
		lastPage = int(respLastPage)
	}

	respProxyList, ok := respData["data"].([]interface{})
	if !ok {
		err = errors.New("parse response proxy json data fail")
		return
	}
	for _, info := range respProxyList {
		proxyInfo, ok := info.(map[string]interface{})
		if !ok {
			continue
		}
		newProxyInfo := &ProxyInfo{
			IP:   proxyInfo["ip"].(string),
			Port: proxyInfo["port"].(string),
			Anonymity: func(j float64) string {
				i := int(j)
				var name string
				switch i {
				case 1:
					name = "透明"
				case 2:
					name = "高匿"
				default:
					name = "未知"
				}
				return name
			}(proxyInfo["anonymity"].(float64)),
			Protocol: proxyInfo["protocol"].(string),
			Address: func(addr string) string {
				return strings.TrimSpace(strings.Replace(addr, "X", "", -1))
			}(proxyInfo["ip_address"].(string)),
			Country: func(country string) string {
				return strings.TrimSpace(strings.Replace(country, "X", "", -1))
			}(proxyInfo["country"].(string)),
			ISP: func(isp string) string {
				return strings.TrimSpace(strings.Replace(isp, "X", "", -1))
			}(proxyInfo["isp"].(string)),
		}

		completionProxyInfo(newProxyInfo)

		checkProxyJobChan <- newProxyInfo
	}
	return
}

// 声明配置爬虫
func withFreeipSpider() spiderOption {
	return func() (spider spiderInterface) {
		freeip := &freeipSpider{
			spiderBaseInfo: spiderBaseInfo{
				name:     "高可用全球免费代理IP库",
				keyName:  "freeipSpiderProxyPool",
				checkUrl: "www.freeip.top",
			},
			url:       "https://ip.jiangxianli.com/api/proxy_ips?page=",
			pages:     20,
			failCount: 10,
		}
		spider = freeip
		return
	}
}

/*
  西刺免费代理IP
  https://www.xicidaili.com
*/
type xiciSpider struct {
	spiderBaseInfo
	urls      []string
	failCount int
	pages     int
	wg        sync.WaitGroup
}

func (s *xiciSpider) startSpider() error {
	for _, url := range s.urls {
		logger.Echo.Debugf("start run %s, url: %s", s.name, url)
		s.wg.Add(1)
		go s.getProxyData(url)
	}
	s.wg.Wait()
	return nil
}

func (s *xiciSpider) getProxyData(url string) {
	defer s.wg.Done()
	failCount := 0
	for i := 1; i <= s.pages; i++ {
		if failCount > s.failCount {
			logger.Echo.Errorf("%s site may have been hung up", s.name)
			break
		}
		dataUrl := fmt.Sprintf("%s%d", url, i)
		proxyIp := modles.GetOneProxyIp(s.spiderBaseInfo.keyName)
		data, response, err := utils.HttpProxyGet(dataUrl, proxyIp, nil)
		if err != nil {
			failCount++
			i--
			logger.Echo.Debugf("run %s spider,url: %s, err: %v", s.name, dataUrl, err)

			_ = modles.DelOneProxyFromRedis(s.spiderBaseInfo.keyName, proxyIp)
			utils.TimerUtil(1)
			continue
		}
		if response.StatusCode == 404 {
			break
		}
		err = s.parseProxyData(data)
		if err != nil {
			failCount++
			i--
			logger.Echo.Debug("run %s spider,url: %s, err: %v", s.name, dataUrl, err)
			utils.TimerUtil(1)
			continue
		}
		failCount = 0
		utils.TimerUtil(2)
		logger.Echo.Debugf("run %s spider,get %s data success", s.name, dataUrl)
	}
	logger.Echo.Debugf("run %s spider end", s.name)
	return
}

func (s *xiciSpider) parseProxyData(data []byte) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("parse html dom fail, err: %v", e)
			return
		}
	}()
	rootDom, err := htmlquery.Parse(bytes.NewReader(data))
	if err != nil {
		err = fmt.Errorf("parse html dom fail, err: %v", err)
		return
	}

	trs := htmlquery.Find(rootDom, "//table[@id='ip_list']//tr[@class]")
	for _, tr := range trs {
		tds := htmlquery.Find(tr, ".//tr/*")
		proxyInfo := &ProxyInfo{
			IP:        htmlquery.InnerText(tds[1]),
			Port:      htmlquery.InnerText(tds[2]),
			Anonymity: htmlquery.InnerText(tds[4]),
			Protocol:  strings.ToLower(htmlquery.InnerText(tds[5])),
			Country:   "中国",
			Address: func(src string) string {
				var address string
				if src == "" {
					return ""
				} else {
					srcSlice := config.Seg.Cut(src, true)
					address += "中国"
					for _, item := range srcSlice {
						address += strings.TrimSpace(item) + " "
					}
				}
				return strings.TrimSpace(address)
			}(htmlquery.InnerText(tds[3])),
		}

		completionProxyInfo(proxyInfo)

		checkProxyJobChan <- proxyInfo
	}
	return
}

func (s *xiciSpider) spiderName() string {
	return s.spiderBaseInfo.name
}

func (s *xiciSpider) redisProxyPoolInfo() (keyName, checkUrl string) {
	return s.spiderBaseInfo.keyName, s.spiderBaseInfo.checkUrl
}

func (s *xiciSpider) spiderTaskSpec() string {
	return s.spec
}

// 声明配置爬虫
func withXiciSpider() spiderOption {
	return func() (spider spiderInterface) {
		xici := &xiciSpider{
			spiderBaseInfo: spiderBaseInfo{
				name:     "西刺免费代理ip",
				keyName:  "xiciSpiderProxyPool",
				checkUrl: "www.xicidaili.com",
			},
			urls: []string{
				"https://www.xicidaili.com/nn/",
				"https://www.xicidaili.com/nt/",
				"https://www.xicidaili.com/wn/",
				"https://www.xicidaili.com/wt/",
			},
			failCount: 6,
			pages:     10,
		}
		spider = xici
		return
	}
}

/*
  极速数据
  http://www.superfastip.com/
*/
type superfipSpider struct {
	spiderBaseInfo
	url       string
	failCount int
	pages     int
}

func (s *superfipSpider) startSpider() error {
	err := s.getProxyData()
	if err != nil {
		return err
	}
	return nil
}

func (s *superfipSpider) getProxyData() (err error) {
	failCount := 0
	for i := 1; i <= s.pages; i++ {
		if failCount > s.failCount {
			logger.Echo.Errorf("%s site may have been hung up", s.name)
			break
		}
		dataUrl := fmt.Sprintf("%s%d", s.url, i)
		proxyIp := modles.GetOneProxyIp(s.spiderBaseInfo.keyName)
		data, response, err := utils.HttpProxyGet(dataUrl, proxyIp, nil)
		if err != nil {
			failCount++
			i--
			logger.Echo.Debugf("run %s spider,url: %s, err: %v", s.name, dataUrl, err)

			_ = modles.DelOneProxyFromRedis(s.spiderBaseInfo.keyName, proxyIp)
			utils.TimerUtil(1)
			continue
		}
		if response.StatusCode == 404 {
			break
		}
		err = s.parseProxyData(data)
		if err != nil {
			failCount++
			i--
			logger.Echo.Debugf("run %s spider,url: %s, err: %v", s.name, dataUrl, err)
			utils.TimerUtil(1)
			continue
		}
		failCount = 0
		utils.TimerUtil(2)
		logger.Echo.Debugf("run %s spider,get %s data success", s.name, dataUrl)
	}
	logger.Echo.Debugf("run %s spider end", s.name)
	return
}

func (s *superfipSpider) parseProxyData(data []byte) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("parse html dom fail, err: %v", e)
			return
		}
	}()

	rootDom, err := htmlquery.Parse(bytes.NewReader(data))
	if err != nil {
		err = fmt.Errorf("parse html dom fail, err: %v", err)
		return
	}
	table := htmlquery.Find(rootDom, "//table/tbody[1]")[1]
	trs := htmlquery.Find(table, ".//tr")
	for _, tr := range trs {
		tds := htmlquery.Find(tr, ".//tr/*")
		proxyInfo := &ProxyInfo{
			IP:   htmlquery.InnerText(tds[0]),
			Port: htmlquery.InnerText(tds[1]),
			Anonymity: func(src string) string {
				switch src {
				case "高级隐私":
					return "高匿"
				default:
					return "透明"
				}
			}(htmlquery.InnerText(tds[2])),
			Protocol: strings.ToLower(htmlquery.InnerText(tds[3])),
			Country:  htmlquery.InnerText(tds[4]),
		}

		completionProxyInfo(proxyInfo)

		checkProxyJobChan <- proxyInfo
	}
	return
}

func (s *superfipSpider) spiderName() string {
	return s.spiderBaseInfo.name
}

func (s *superfipSpider) redisProxyPoolInfo() (keyName, checkUrl string) {
	return s.spiderBaseInfo.keyName, s.spiderBaseInfo.checkUrl
}

func (s *superfipSpider) spiderTaskSpec() string {
	return s.spec
}

// 声明配置爬虫
func withSuperipSpider() spiderOption {
	return func() (spider spiderInterface) {
		superfip := &superfipSpider{
			spiderBaseInfo: spiderBaseInfo{
				name:     "极速数据",
				keyName:  "SuperipSpiderProxyPool",
				checkUrl: "www.superfastip.com",
			},
			url:       "http://www.superfastip.com/welcome/freeip/",
			failCount: 20,
			pages:     10,
		}
		spider = superfip
		return
	}
}


/*
startSpider() error
	redisProxyPoolInfo() (keyName, checkUrl string)
	spiderTaskSpec() string
	spiderName() string
 */
/*
  local github proxy pool
  http://www.superfastip.com/
*/
type localProxySpider struct {
	spiderBaseInfo
	url       string
	failCount int
	pages     int
}

func (s *localProxySpider) startSpider() (err error) {
	return
}

func (s *localProxySpider) redisProxyPoolInfo()(keyName, checkUrl string) {
	return s.keyName, s.checkUrl
}

func (s *localProxySpider) spiderTaskSpec() string {
	return s.spec
}

func (s *localProxySpider) spiderName() string {
	return s.name
}

// 声明配置爬虫
func withLocalProxySpider() spiderOption {
	return func() (spider spiderInterface) {
		localProxy := &localProxySpider{
			spiderBaseInfo: spiderBaseInfo{
				name:     "本地IP池(Github)",
				keyName:  "localProxyPool",
				checkUrl: "127.0.0.1:15010",
			},
			url:       "http://127.0.0.1:15010/get_all/",
			failCount: 20,
			pages:     10,
		}
		spider = localProxy
		return
	}
}