package modles

import (
	"apiTools/libs/logger"
	"apiTools/utils"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"math/rand"
	"strings"
	"sync"
	"time"
)

// 定义ip代理存储的结构体
type ProxyPool struct {
	ID         int
	IP         string `gorm:"size:32;not null;unique_index"`
	Port       string `gorm:"size:8;not null"`             // proxy 端口
	Anonymity  string `gorm:"size:16;not null"`            // 匿名类型(透明/高匿)
	Protocol   string `gorm:"size:8;not null"`             // 协议类型(http/https)
	Country    sql.NullString                              // 所在国家
	Address    sql.NullString                              // 所在地区
	ISP        sql.NullString `gorm:"size:50" `            // 运营商
	Speed      sql.NullInt64  `gorm:"not null;default:0" ` // 响应速度(毫秒)
	FailCount  uint           `gorm:"not null;default:0"`  // 测试失败次数, 失败次数超过三次删除
	VerifyTime time.Time                                   // 最后验证时间
}

type ProxyPoolForm struct {
	// 每页默认30个
	Page      uint   `form:"page" json:"page" xml:"page"`                   // 当前页码
	Protocol  string `form:"protocol" json:"protocol" xml:"protocol"`       // 代理的协议类型(http/https)
	Anonymity string `form:"anonymity" json:"anonymity" xml:"anonymity"`    // 匿名类型(透明/匿名/高匿)
	Country   string `form:"country" json:"country" xml:"country"`          // 所在国家
	Address   string `form:"address" json:"address" xml:"address"`          // 所在地区
	ISP       string `form:"isp" json:"isp" xml:"isp"`                      // 运营商
	OrderBy   string `form:"order_by" json:"order_by" xml:"order_by"`       // 排序字段 (speed:响应速度,verify_time:校验时间)
	OrderRule string `form:"order_rule" json:"order_rule" xml:"order_rule"` // 排序规则(desc:降序 asc:升序)
}

// 单个代理数据结构
type SingleProxyInfo struct {
	IP         string `json:"ip"`
	Port       string `json:"port"`
	Anonymity  string `json:"anonymity"`
	Protocol   string `json:"protocol"`
	Country    string `json:"country"`
	Address    string `json:"address"`
	ISP        string `json:"isp"`
	Speed      int    `json:"speed"`
	VerifyTime string `json:"verify_time"`
}

// api查询返回的数据结构
type ProxyPoolResult struct {
	ProxyPools []*SingleProxyInfo `json:"data"`
	Pages      uint               `json:"pages"`
}

// 查询proxy数据
func QueryProxyPoolInfo(proxyPoolForm *ProxyPoolForm) (proxyPoolResult *ProxyPoolResult, err error) {
	var pageSize uint = 15
	proxyPoolResult = &ProxyPoolResult{
		ProxyPools: make([]*SingleProxyInfo, 0, 30),
	}

	if proxyPoolForm.Page == 0 {
		proxyPoolForm.Page = 1
	}
	offset := (proxyPoolForm.Page - 1) * pageSize
	selectObj := SqlConn.Model(&ProxyPool{}).Where("speed != 0")
	if proxyPoolForm.Protocol != "" {
		proxyPoolForm.Protocol = strings.ToLower(proxyPoolForm.Protocol)
		protocolSlice := []string{"http", "https"}
		isIn := utils.IsInSelic(proxyPoolForm.Protocol, protocolSlice)
		if !isIn {
			err = errors.New("protocol parameter passed incorrectly")
			return
		}
		selectObj = selectObj.Where("protocol = ?", proxyPoolForm.Protocol)

	}
	if proxyPoolForm.Anonymity != "" {
		anonymitySlice := []string{"透明", "高匿"}
		isIn := utils.IsInSelic(proxyPoolForm.Anonymity, anonymitySlice)
		if !isIn {
			err = errors.New("anonymity parameter passed incorrectly")
			return
		}
		selectObj = selectObj.Where("anonymity = ?", proxyPoolForm.Anonymity)

	}
	if proxyPoolForm.Country != "" {
		selectObj = selectObj.Where("country = ?", proxyPoolForm.Country)
	}
	if proxyPoolForm.Address != "" {
		selectObj = selectObj.Where("address LIKE ?", fmt.Sprintf("%%%s%%", proxyPoolForm.Address))
	}
	if proxyPoolForm.ISP != "" {
		selectObj = selectObj.Where("isp LIKE ?", fmt.Sprintf("%%%s%%", proxyPoolForm.ISP))
	}
	if proxyPoolForm.OrderBy != "" {
		orderBySlice := []string{"speed", "verify_time"}
		isIn := utils.IsInSelic(proxyPoolForm.OrderBy, orderBySlice)
		if !isIn {
			err = errors.New("order_by parameter passed incorrectly")
			return
		}
		if proxyPoolForm.OrderRule == "" {
			proxyPoolForm.OrderRule = "asc"
		} else {
			orderRuleSlice := []string{"desc", "asc"}
			isIn := utils.IsInSelic(proxyPoolForm.OrderRule, orderRuleSlice)
			if !isIn {
				err = errors.New("order_rule parameter passed incorrectly")
				return
			}
		}
		selectObj = selectObj.Order(fmt.Sprintf("%s %s", proxyPoolForm.OrderBy, proxyPoolForm.OrderRule))
	}
	countObj := *selectObj
	var proxyCount uint
	err = countObj.Count(&proxyCount).Error
	// 获取总页数
	if err == nil {
		switch {
		case proxyCount > 30:
			t := proxyCount / pageSize
			if proxyCount%pageSize > 0 {
				t++
			}
			proxyPoolResult.Pages = t
		case proxyCount > 0 && proxyCount <= 30:
			proxyPoolResult.Pages = 1
		default:
			proxyPoolResult.Pages = 0
		}
	}
	// 查询详细数据
	var ProxyPools []*ProxyPool
	err = selectObj.Offset(offset).Limit(pageSize).Find(&ProxyPools).Error
	for _, info := range ProxyPools {
		proxyInfo := &SingleProxyInfo{
			IP:         info.IP,
			Port:       info.Port,
			Anonymity:  info.Anonymity,
			Protocol:   info.Protocol,
			Country:    info.Country.String,
			Address:    info.Address.String,
			ISP:        info.ISP.String,
			Speed:      int(info.Speed.Int64),
			VerifyTime: info.VerifyTime.Format("2006/01/02 15:04:05"),
		}
		proxyPoolResult.ProxyPools = append(proxyPoolResult.ProxyPools, proxyInfo)
	}
	return
}

// 插入ip代理数据，如果ip不存在
// ip 查询ip
// proxyInfo 代理数据信息
// isFromDB 标记此数据是否从数据库提取出来的
func InsertProxyInfo(proxyInfo *ProxyPool, isFromDB bool) (err error) {
	oldProxyInfo := &ProxyPool{}
	SqlConn.Where(&ProxyPool{IP: proxyInfo.IP}).First(oldProxyInfo)
	if oldProxyInfo.IP == "" {
		if err = SqlConn.Create(proxyInfo).Error; err != nil {
			return
		}
	} else {
		proxyInfo.ID = oldProxyInfo.ID
		if isFromDB == false {
			proxyInfo.FailCount += oldProxyInfo.FailCount
		}
		if err = SqlConn.Save(proxyInfo).Error; err != nil {
			return
		}
	}
	return
}

// 从数据库中提取代理数据(定时任务数据库检测)
func ExtractProxyInfo(count int) ([]*ProxyPool, error) {
	var proxyPoolList []*ProxyPool
	err := SqlConn.Limit(count).Order("verify_time").Find(&proxyPoolList).Error
	if err != nil {
		return proxyPoolList, err
	}
	return proxyPoolList, nil
}

// 从数据库中读取最近校验成功的代理信息(存入redis用于其他api使用)
func GetLatestProxyInfo(count int) ([]string, error) {
	var proxyPoolList []*ProxyPool
	err := SqlConn.Limit(count).Select([]string{"ip", "port", "fail_count", "country", "speed", "verify_time"}).
		Where("speed BETWEEN ? AND ? AND fail_count <= ? AND country LIKE ?", 1, 5000, 1, "中国%").
		Order("verify_time desc, speed").Find(&proxyPoolList).Error
	if err != nil {
		return nil, err
	}
	proxyInfoArray := make([]string, 0, count)
	for _, info := range proxyPoolList {
		proxyInfoArray = append(proxyInfoArray, fmt.Sprintf("%s:%s", info.IP, info.Port))
	}

	return proxyInfoArray, nil
}

// 删除数据库中的一条代理信息
func DelOneProxyFromDB(ip string) error {
	err := SqlConn.Where(ProxyPool{IP: ip}).Delete(&ProxyPool{}).Error
	return err
}

// 提取数据库代理校验并存入redis
func ExtractProxyToRedis(keyName string, checkUrl string) bool {
	rand.Seed(time.Now().UnixNano())
	extractCount := rand.Intn(15)
	proxyArray, err := GetLatestProxyInfo(extractCount)
	if err != nil {
		logger.Echo.Errorf("proxy app: extract proxy info from database(cron name: extractProxyToRedis), err: %v", err)
		return false
	}
	// 协程去校验ip
	waitGroup := sync.WaitGroup{}
	resultChan := make(chan string, extractCount)
	for _, proxyAddr := range proxyArray {
		waitGroup.Add(1)
		go func(proxyAddr string, resultChan chan string) {
			defer waitGroup.Done()
			status := utils.CheckProtocolHttp(proxyAddr, checkUrl)
			if status {
				resultChan <- proxyAddr
			}
		}(proxyAddr, resultChan)
	}
	waitGroup.Wait()
	close(resultChan)
	// 读取校验完成的ip
	newProxyArray := make([]string, 0, extractCount)
	for info := range resultChan {
		newProxyArray = append(newProxyArray, info)
	}
	if len(newProxyArray) == 0 {
		return false
	}
	err = SetProxyInfoToRedis(keyName, newProxyArray)
	if err != nil {
		logger.Echo.Errorf("proxy app: save proxy info to database fail(cron name: extractProxyToRedis), err: %v", err)
		return false
	}
	return true
}

// 设置代理信息到redis数据库
func SetProxyInfoToRedis(keyName string, proxyList []string) error {
	redisClient := RedisPool.Get()
	defer redisClient.Close()
	redisClient.Send("MULTI")
	redisClient.Do("DEL", keyName)
	for _, ip := range proxyList {
		redisClient.Do("LPUSH", keyName, ip)
	}
	redisClient.Send("EXEC")
	return nil
}

// 从redis数据库中读取代理信息
func ReadProxyInfoFromRedis(keyName string) ([]string, error) {
	redisClient := RedisPool.Get()
	defer redisClient.Close()
	proxyArray, err := redis.Strings(redisClient.Do("LRANGE", keyName, 0, -1))
	if err != nil {
		return nil, err
	}
	return proxyArray, nil
}

// 删除redis 中指定的一个代理
func DelOneProxyFromRedis(keyName string, proxyIp string) error {
	redisClient := RedisPool.Get()
	defer redisClient.Close()
	_, err := redisClient.Do("LREM", keyName, 1, proxyIp)
	return err
}

// 从redis中随机获取一个proxy ip
func GetOneProxyIp(keyName string) string {
	// 从redis中读取可用代理信息
	proxyArray, err := ReadProxyInfoFromRedis(keyName)
	if err != nil || len(proxyArray) == 0 {
		return ""
	}
	// 随机获取一个代理ip
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(proxyArray))
	return proxyArray[index]
}
