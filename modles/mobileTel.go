package modles

import (
	"apiTools/utils"
	"errors"
	"github.com/axgle/mahonia"
	"strings"
)

// 接口
// https://tcc.taobao.com/cc/json/mobile_tel_segment.htm?tel=手机号
// https://www.baifubao.com/callback?cmd=1059&callback=phone&phone=手机号

const (
	telApiUrl         = "https://tcc.taobao.com/cc/json/mobile_tel_segment.htm?tel="
	notFoundTelErr    = "telephone not exists"
	queryTelFailedErr = "get tel info failed"
)

// 请求表单
type MobileTelForm struct {
	Tel string `form:"tel" json:"tel" xml:"tel" binding:"required"` // 手机号码
}

type MobileTelResponse struct {
	Tel     string `json:"tel"`     // 手机号
	Mst     string `json:"mst"`     // 号码区段
	Carrier string `json:"carrier"` // 运营商
	Area    string `json:"area"`    // 号码归属地
}

func QueryTelInfo(form *MobileTelForm) (resp *MobileTelResponse, msg string, err error) {
	data, err := getTelData(form.Tel)
	if err != nil {
		msg = queryTelFailedErr
		return
	}
	response, status, err := parseTelData(data)
	if err != nil {
		msg = queryTelFailedErr
		return
	}
	if !status {
		msg = notFoundTelErr
		err = errors.New("mobile tel not found")
		return
	}
	resp = response
	return
}

func getTelData(tel string) (data []byte, err error) {
	reqUrl := telApiUrl + tel
	data, _, err = utils.HttpProxyGet(reqUrl, "", nil)
	if err != nil {
		return
	}
	return
}

func parseTelData(data []byte) (resp *MobileTelResponse, status bool, err error) {
	strData := mahonia.NewDecoder("gbk").ConvertString(string(data))
	dataSlice := strings.Split(strData, ",")
	if len(dataSlice) < 5 {
		return
	}
	resp = &MobileTelResponse{
		Tel: func() string {
			split := strings.Split(dataSlice[3], ":")
			s := strings.Trim(split[1], "'")
			return s
		}(),
		Mst: func() string {
			split := strings.Split(dataSlice[0], ":")
			s := strings.Trim(split[1], "'")
			return s
		}(),
		Carrier: func() string {
			split := strings.Split(dataSlice[2], ":")
			s := strings.Trim(split[1], "'")
			return s
		}(),
		Area: func() string {
			split := strings.Split(dataSlice[1], ":")
			s := strings.Trim(split[1], "'")
			return s
		}(),
	}
	status = true
	return
}
