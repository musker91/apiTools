package modles

import (
	"apiTools/libs/config"
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"
)

func init() {
	//initial()
}

func initial() {
	err := config.InitConfig()
	if err != nil {
		panic(err)
	}
	err = InitRedis()
	if err != nil {
		panic(err)
	}
	err = InitMysql()
	if err != nil {
		panic(err)
	}
	err = InitApiConfig()
	if err != nil {
		panic(err)
	}
}

func TestWhoisQuery(t *testing.T) {
	form := &WhoisForm{
		Domain:  "http://www.baidu.io",
		OutType: "json",
	}
	whoisInfo, err := QueryWhoisInfoToJson(form)
	if err != nil {
		t.Error("err", err)
	}
	t.Logf("%v\n", whoisInfo.TextInfo)
	t.Logf("%#v\n", whoisInfo.JsonInfo)
}

func TestToShortUrl(t *testing.T) {
	shortForm := &ShortForm{
		Url:        "https://www.runoob.com/python3/python-find-url-string.html",
		Domain:     "http://www.baidu.cn",
		ExpireTime: 7,
	}

	shortInfo, err := ToShortUrl(shortForm)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(shortInfo)
}

func TestParseShortUrl(t *testing.T) {
	shortUrl := "http://www.baidu.com/DF5m1YsVSf"
	shortInfo, err := ParseShort(shortUrl)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(shortInfo)
}

func TestInsertProxyInfo(t *testing.T) {
	proxyInfo := &ProxyPool{
		IP:         "1.1.1.2",
		Port:       "8089",
		Anonymity:  "高匿",
		Protocol:   "https",
		Speed:      sql.NullInt64{Int64: 1992, Valid: true},
		VerifyTime: time.Now(),
	}
	err := InsertProxyInfo(proxyInfo, false)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("insert success")
}

func TestCreateProxyInfo(t *testing.T) {
	for i := 1; i < 255; i++ {
		proxyInfo := &ProxyPool{
			IP:         fmt.Sprintf("%d:%d:%d:%d", i, i, i, i),
			Port:       fmt.Sprintf("8%d", i),
			Anonymity:  "透明",
			Protocol:   "http",
			VerifyTime: time.Now(),
		}
		err := InsertProxyInfo(proxyInfo, false)
		if err != nil {
			t.Error(err)
		}
		log.Printf("insert %s:%s success", proxyInfo.IP, proxyInfo.Port)
		time.Sleep(1 * time.Second)
	}
	t.Log("insert success")

}

func TestExtractProxyInfo(t *testing.T) {
	pools, err := ExtractProxyInfo(10)
	if err != nil {
		t.Error(err)
		return
	}
	for index, info := range pools {
		t.Logf("%d --> %#v\n", index, info)
	}
}

func TestGetLatestProxyInfo(t *testing.T) {
	proxyArray, err := GetLatestProxyInfo(10)
	if err != nil {
		t.Error(err)
		return
	}
	for index, info := range proxyArray {
		t.Logf("%d --> %s\n", index, info)
	}
}

func TestSetProxyInfoToRedis(t *testing.T) {
	proxyArray := []string{
		"254:254:254:254:8254",
		"253:253:253:253:8253",
		"252:252:252:252:8252",
		"251:251:251:251:8251",
		"250:250:250:250:8250",
		"249:249:249:249:8249",
		"248:248:248:248:8248",
		"247:247:247:247:8247",
		"246:246:246:246:8246",
		"245:245:245:245:8245",
	}
	err := SetProxyInfoToRedis("proxyPoolArray", proxyArray)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("set success")

}

func TestReadProxyInfoFromRedis(t *testing.T) {
	proxyArray, err := ReadProxyInfoFromRedis("proxyPoolArray")
	if err != nil {
		t.Error(err)
		return
	}
	for index, info := range proxyArray {
		t.Logf("%d --> %s\n", index, info)
	}
}

func TestDelOneProxyFromDB(t *testing.T) {
	ip := "183.166.133.146"
	err := DelOneProxyFromDB(ip)
	if err != nil {
		t.Error(err)
	}
	t.Log("del proxy ip success")
}

func TestQueryProxyPoolInfo(t *testing.T) {
	proxyPoolForm := &ProxyPoolForm{
		Page:      1,
		Country:   "中国",
		Protocol:  "https",
		Address:   "南京",
		OrderBy:   "speed",
		OrderRule: "desc",
	}
	proxyPools, err := QueryProxyPoolInfo(proxyPoolForm)
	if err != nil {
		t.Error(err)
		return
	}
	println("pages is: ", proxyPools.Pages)
	for index, info := range proxyPools.ProxyPools {
		fmt.Printf("%d --> %v\n", index, info)
	}
}

func TestGetBinLens(t *testing.T) {
	bankCard := "6222600260001072444"
	result, err := getBinLens(bankCard)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(result)
}

func TestGetBankInfo(t *testing.T) {
	bankCard := "6222600260001072444"
	binLenSlice := []int{6, 8, 4, 3, 5, 9, 7, 10}
	result, err := getBankInfo(bankCard, binLenSlice)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(result)
}

func TestIsBinkCard(t *testing.T) {
	bankCard := "88888888888888888888"
	card := isBinkCard(bankCard)
	t.Log(card)
}

func TestGetAlipayBankInfo(t *testing.T) {
	bankCard := "6222600260001072444"
	result, err := getAlipayBankInfo(bankCard)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(result)
}

func TestGetBankNameFenName(t *testing.T) {
	enName := "COMM"
	bankName, err := getBankNameFenName(enName)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(bankName)
}

func TestOfficialToShort(t *testing.T) {
	form := ShortForm{
		Domain: "http://www.mazhichao.com",
		Type:   0,
	}
	shortInfo, msg, err := OfficialToShort(&form)
	if err != nil {
		t.Error(err, msg)
		return
	}
	t.Log(msg)
	t.Logf("%#v\n", shortInfo)
}

func TestGetIcpInfo(t *testing.T) {
	domain := "baidu.com"
	data, err := getIcpInfo(domain)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(data))
}

func TestParseIcpData(t *testing.T) {
	domain := "hefupal.com"
	data, err := getIcpInfo(domain)
	if err != nil {
		t.Error(err)
		return
	}
	icpResponse, status, err := parseIcpData(data)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(icpResponse)
	t.Log(status)
}

func TestParseTelData(t *testing.T) {
	tel := "13520360558"
	data, err := getTelData(tel)
	if err != nil {
		t.Log(err)
		return
	}
	resp, status, err := parseTelData(data)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(resp)
	t.Log(status)
}

func TestGetDomainData(t *testing.T) {
	domain := "www.baidu.com"
	data, err := getDomainData(domain)
	if err != nil {
		t.Error(err)
		return
	}
	resp, err := parseDomainData(data)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(resp)
}

func TestShortVideoParse(t *testing.T) {
	url := "https://v.kuaishou.com/s/iMxeowhG"
	form := &ShortVideoForm{Url:url}

	result, err := ShortVideoParse(form)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf( "%#v\n",result)
}
