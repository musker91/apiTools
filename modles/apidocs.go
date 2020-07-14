package modles

import (
	"apiTools/libs/config"
	"apiTools/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"gopkg.in/russross/blackfriday.v2"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// 全部api简要信息和其他该要信息
type AllApiInfo struct {
	ApiList      []*ApiInfo // 所有api简要信息
	ApiCount     int        // api个数
	RequestCount int        // 所有api总请求次数
	RunTime      int        // 系统总运行天数
}

// api简要信息
type ApiInfo struct {
	TitleName        string // Api标题名称
	Description      string // 描述信息
	ApiUrl           string // api url
	RequestCount     int    // api总请求次数
	RequestCountShow string // 页面显示次数
	Mainten          bool   // 是否处于维护状态
}

// api详细信息
type ApiDocInfo struct {
	ApiName       string     // api名称
	Description   string     // 描述信息
	ApiAddr       string     // api地址
	RsponseType   string     // 返回数据格式
	RequestMethod string     // 请求方式
	RequestDemo   string     // 请求示例
	ReponseDemo   string     // 返回示例
	ReqParam      [][]string // 请求参数列表
	RepParam      [][]string // 返回参数列表
	ErrParam      [][]string // 错误参照码列表
	DemoCode      string     // 示例代码
	Mainten       bool       // 接口是否维护中
	//RunTime       int        // 系统总运行天数
	//RequestCount  int        // Api调用请求次数
}

// 初始化api docs json数据
func InitApiDocsJsonData() (err error) {
	jsonFile, err := ioutil.ReadFile(filepath.Join(utils.GetRootPath(), "data/api", "apidocs.json"))
	if err != nil {
		return
	}
	err = json.Unmarshal(jsonFile, &JsonData)
	if err != nil {
		return
	}
	return
}

// 获取全部api简要信息和其他该要信息
func GetAllApiInfo() (allApiInfo *AllApiInfo, err error) {
	allApiInfo = &AllApiInfo{}
	var apiList []*ApiInfo
	// 循环json数据获取api list
	jsonData, ok := JsonData.(map[string]interface{})
	if !ok {
		return allApiInfo, errors.New("json data to map fail")
	}
	redisClient := RedisPool.Get()
	defer redisClient.Close()
	for key, value := range jsonData {
		apiInfo := &ApiInfo{}
		apiJson, ok := value.(map[string]interface{})
		if !ok {
			continue
		}
		// 判断api是否被禁用
		if apiJson["enable"] == false {
			continue
		}

		// 组装apiInfo数据
		apiInfo.ApiUrl = key
		apiInfo.Mainten = apiJson["mainten"].(bool)
		apiInfo.Description = apiJson["description"].(string)
		apiInfo.TitleName = fmt.Sprintf("%v", apiJson["titleName"])
		requestCount := getApiReqCount(fmt.Sprintf("%v", apiJson["countKey"]), redisClient)
		apiInfo.RequestCount = requestCount
		rwt := requestCount / 10000
		if (requestCount / 10000) >= 1 {
			dr := strconv.Itoa(requestCount % 10000)
			if len(dr) == 4 {
				apiInfo.RequestCountShow = fmt.Sprintf("%v.%v W", rwt, dr[0])
			} else {
				apiInfo.RequestCountShow = fmt.Sprintf("%v W", rwt)
			}
		} else {
			apiInfo.RequestCountShow = fmt.Sprintf("%v", requestCount)
		}
		apiList = append(apiList, apiInfo)
	}
	allApiInfo.ApiList = apiList
	// api总数
	allApiInfo.ApiCount = len(apiList)
	// 统计所有api总请求次数
	allApiInfo.RequestCount = countAllApiReq(apiList)
	// 计算系统当前运行时长
	allApiInfo.RunTime = calcSystemRunTime()
	return
}

// 获取单个api的请求统计
func getApiReqCount(countKey string, redisClient redis.Conn) (count int) {
	countKeyName := fmt.Sprintf("apiCount_%s", countKey)
	count, err := redis.Int(redisClient.Do("GET", countKeyName))
	if err != nil {
		return 0
	}
	return
}

// 统计所有api的请求次数
func countAllApiReq(apiList []*ApiInfo) (count int) {
	for _, apiInfo := range apiList {
		count += apiInfo.RequestCount
	}
	return
}

// 计算系统运行时长
func calcSystemRunTime() (day int) {
	startTime, ok := config.Get("web::startTime").(time.Time)
	if !ok {
		return
	}
	apartHours := int(time.Since(startTime).Hours())
	return (apartHours / 24) + 1
}

// 获取指定api文档的信息
func GetApiDocInfo(apiFileName string, urlPath string, countKey string, mainten bool) (apiDocInfo *ApiDocInfo, err error) {
	// 初始化返回数据结构
	apiDocInfo = &ApiDocInfo{}
	// 读取markdown文件数据
	docFile, err := ioutil.ReadFile(filepath.Join(utils.GetRootPath(), "data/api/docs", apiFileName))
	if err != nil {
		return
	}
	// 将markdown数据转换为html代码
	blockCode := bytes.NewReader(blackfriday.Run(docFile))
	document, err := goquery.NewDocumentFromReader(blockCode)
	// 查找api名称
	document.Find(".api-title").Each(func(i int, s *goquery.Selection) {
		apiDocInfo.ApiName = s.Text()
	})
	// 查找描述信息
	document.Find(".api-desc").Each(func(i int, s *goquery.Selection) {
		apiDocInfo.Description = s.Text()
	})
	// 查找api地址
	document.Find(".api-url").Each(func(i int, s *goquery.Selection) {
		apiDocInfo.ApiAddr = s.Text()
	})
	// 查找返回的数据格式
	document.Find(".api-reponse-format").Each(func(i int, s *goquery.Selection) {
		apiDocInfo.RsponseType = s.Text()
	})
	// 查找请求方法
	document.Find(".api-request-method").Each(func(i int, s *goquery.Selection) {
		apiDocInfo.RequestMethod = s.Text()
	})
	// 请求示例
	document.Find(".api-request-demo").Parent().NextFiltered("pre").Each(func(i int, s *goquery.Selection) {
		apiDocInfo.RequestDemo = strings.TrimSpace(s.Text())
	})
	// 返回示例
	document.Find(".api-reponse-demo").Parent().NextFiltered("pre").Each(func(i int, s *goquery.Selection) {
		apiDocInfo.ReponseDemo = strings.TrimSpace(s.Text())
	})
	// 请求参数说明列表
	document.Find(".request-param").Parent().NextFiltered("table").Each(func(i int, s *goquery.Selection) {
		trNodes := s.Find("tbody>tr")
		reqParam := make([][]string, 0, len(trNodes.Nodes))
		for _, tr := range trNodes.Nodes {
			tdNode := s.FindNodes(tr)
			for _, td := range tdNode.Nodes {
				tdText := strings.TrimSpace(s.FindNodes(td).Text())
				tdTextSlice := strings.Split(tdText, "\n")
				reqParam = append(reqParam, tdTextSlice)
			}
		}
		apiDocInfo.ReqParam = reqParam
	})

	// 返回参数说明列表
	document.Find(".reponse-param").Parent().NextFiltered("table").Each(func(i int, s *goquery.Selection) {
		trNodes := s.Find("tbody>tr")
		repParam := make([][]string, 0, len(trNodes.Nodes))
		for _, tr := range trNodes.Nodes {
			tdNode := s.FindNodes(tr)
			for _, td := range tdNode.Nodes {
				tdText := strings.TrimSpace(s.FindNodes(td).Text())
				tdTextSlice := strings.Split(tdText, "\n")
				repParam = append(repParam, tdTextSlice)
			}
		}
		apiDocInfo.RepParam = repParam
	})

	// 错误参照码列表
	document.Find(".error-param").Parent().NextFiltered("table").Each(func(i int, s *goquery.Selection) {
		trNodes := s.Find("tbody>tr")
		errParam := make([][]string, 0, len(trNodes.Nodes))
		for _, tr := range trNodes.Nodes {
			tdNode := s.FindNodes(tr)
			for _, td := range tdNode.Nodes {
				tdText := strings.TrimSpace(s.FindNodes(td).Text())
				tdTextSlice := strings.Split(tdText, "\n")
				errParam = append(errParam, tdTextSlice)
			}
		}
		apiDocInfo.ErrParam = errParam
	})

	// 示例代码
	document.Find(".code-demo").Parent().NextFiltered("pre").Each(func(i int, s *goquery.Selection) {
		apiDocInfo.DemoCode = strings.TrimSpace(s.Text())
	})

	// 是否维护中
	apiDocInfo.Mainten = mainten

	// 计算系统当前运行时长
	//apiDocInfo.RunTime = calcSystemRunTime()

	// 获取api调用次数
	//redisClient := RedisPool.Get()
	//defer redisClient.Close()
	//countKeyName := fmt.Sprintf("apiCount_%s", countKey)
	//count, err := redis.Int(redisClient.Do("GET", countKeyName))
	//if err != nil {
	//	count = 0
	//}
	//apiDocInfo.RequestCount = count
	return
}
