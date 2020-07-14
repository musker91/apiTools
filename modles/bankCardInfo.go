package modles

import (
	"apiTools/utils"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strconv"
)

const (
	BankErrorMsg01    = "bank card number malformed"
	BankErrorMsg02    = "not match card"
	alipayValidateUrl = "https://ccdcapi.alipay.com/validateAndCacheCardInfo.json?_input_charset=utf-8&cardBinCheck=true&cardNo="
)

var (
	cardType = map[string]string{
		"DC":  "储蓄卡",
		"CC":  "信用卡",
		"SCC": "准贷记卡",
		"PC":  "预付费卡",
	}
)

// 数据库表结构
type BankBinInfo struct {
	ID           int
	BankName     string `gorm:"size:64;not null"`
	BankNameEn   string `gorm:"size:64"`
	CardName     string `gorm:"size:64;not null"`
	CardType     string `gorm:"size:64;not null"`
	Bin          string `gorm:"size:16;not null"`
	NumberLength int    `gorm:"not null"`
	BinLength    int    `gorm:"not null"`
	Issueid      int    `gorm:"not null"`
}

// 客户端请求表单结构
type BankInfoForm struct {
	BankCard string `form:"bakCard" json:"bakCard" xml:"banCard" binding:"required"`
}

// 返回的结构
type BankInfoResult struct {
	CardType   string `json:"card_type"`
	BankName   string `json:"bank_name"`
	BankNameEn string `json:"bank_name_en"`
	CardNumber string `json:"card_number"`
}

// 根据银行卡英文名称获取银行卡中文名称
type bankNameStruct struct {
	BankName string `json:"bank_name"`
	Count    int    `json:"count"`
}

// 读取bank数据插入到数据库中
func ReadBankInfoToDB() {
	var bankBinInfo []map[string]interface{}
	jsonFile, err := ioutil.ReadFile(filepath.Join(utils.GetRootPath(), "data/bankInfo/bankbin.json"))
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(jsonFile, &bankBinInfo)
	if err != nil {
		panic(err)
	}

	for _, info := range bankBinInfo {
		bankBin := &BankBinInfo{
			ID:           int(info["id"].(float64)),
			BankName:     info["bankName"].(string),
			BankNameEn:   info["bankNameEn"].(string),
			CardName:     info["cardName"].(string),
			CardType:     info["cardType"].(string),
			Bin:          strconv.Itoa(int(info["bin"].(float64))),
			NumberLength: int(info["nLength"].(float64)),
			BinLength:    int(info["binLength"].(float64)),
			Issueid:      int(info["issueid"].(float64)),
		}
		err := SqlConn.Create(bankBin).Error
		if err != nil {
			fmt.Printf("insert data %#v fail, err: %v\n", bankBin, err)
		}
	}
	fmt.Println("read bank bin info insert to database success!!!")
}

//  查询bank信息
func QueryBankCardInfo(bankInfo *BankInfoForm) (result *BankInfoResult, msg string, err error) {
	isCard := isBinkCard(bankInfo.BankCard)
	if !isCard {
		err = errors.New(BankErrorMsg01)
		msg = BankErrorMsg01
		return
	}
	banLensSlice, err := getBinLens(bankInfo.BankCard)
	if err != nil {
		result, err = getAlipayBankInfo(bankInfo.BankCard)
		if err != nil {
			msg = BankErrorMsg02
		}
		return
	}
	qBankBinInfo, err := getBankInfo(bankInfo.BankCard, banLensSlice)
	if err != nil {
		result, err = getAlipayBankInfo(bankInfo.BankCard)
		if err != nil {
			msg = BankErrorMsg02
		}
		return
	}
	result = &BankInfoResult{
		CardType:   qBankBinInfo.CardType,
		BankName:   qBankBinInfo.BankName,
		BankNameEn: qBankBinInfo.BankNameEn,
		CardNumber: bankInfo.BankCard,
	}
	return
}

// 判断是否为银行卡号
func isBinkCard(card string) bool {
	rex, err := regexp.Compile("^[0-9][0-9]{14,18}$")
	if err != nil {
		return false
	}
	matchString := rex.MatchString(card)
	return matchString
}

// 根据卡长度获取可能的bin长度
func getBinLens(card string) (result []int, err error) {
	var queryData []*BankBinInfo
	err = SqlConn.Select([]string{"distinct bin_length"}).Where(&BankBinInfo{NumberLength: len(card)}).Find(&queryData).Error
	if err != nil {
		return
	}
	for _, d := range queryData {
		result = append(result, d.BinLength)
	}
	return
}

// 根据bin长度列表和卡号获取相应的银行
func getBankInfo(bankCard string, binSlice []int) (result *BankBinInfo, err error) {
	binInfoSlice := make([]string, 0, len(binSlice))
	bankCardRune := []rune(bankCard)
	for _, binLen := range binSlice {
		binInfoSlice = append(binInfoSlice, string(bankCardRune[:binLen]))
	}
	var d []BankBinInfo
	err = SqlConn.Where("bin in (?)", binInfoSlice).Find(&d).Error
	if err != nil {
		return
	}
	if len(d) > 1 {
		err = errors.New("get bank info to many")
		return
	} else if len(d) < 1 {
		err = errors.New("get bank is empty")
		return
	}
	result = &d[0]
	return
}

// 通过支付宝查询银行卡信息
func getAlipayBankInfo(card string) (result *BankInfoResult, err error) {
	reqUrl := alipayValidateUrl + card
	data, _, err := utils.HttpProxyGet(reqUrl, "", nil)
	if err != nil {
		return
	}
	var d map[string]interface{}
	err = json.Unmarshal(data, &d)
	if err != nil {
		return
	}
	if !d["validated"].(bool) {
		err = fmt.Errorf("get getAlipayBankInfo fail, err: %v", d["messages"])
		return
	}
	bankEnName := d["bank"].(string)
	bankName, err := getBankNameFenName(bankEnName)
	if err != nil {
		return
	}
	result = &BankInfoResult{
		CardType:   cardType[d["cardType"].(string)],
		BankName:   bankName,
		BankNameEn: bankEnName,
		CardNumber: card,
	}
	return
}

func getBankNameFenName(enName string) (bankName string, err error) {
	var d []bankNameStruct
	err = SqlConn.Model(&BankBinInfo{}).Select("bank_name, count(bank_name) as count").
		Where("bank_name_en = ?", enName).Group("bank_name").Order("count desc").
		Limit(1).Scan(&d).Error
	if err != nil {
		return
	}
	if len(d) == 0 {
		bankName = enName
	} else {
		bankName = d[0].BankName
	}
	return
}
