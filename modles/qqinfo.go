package modles

import (
	"apiTools/utils"
	"encoding/json"
	"github.com/axgle/mahonia"
	"github.com/pkg/errors"
	"strings"
)

// 获取qq信息
// 头像和昵称

// https://users.qzone.qq.com/fcg-bin/cgi_get_portrait.fcg?uins=652600996

const (
	apiQQUrl        = "https://users.qzone.qq.com/fcg-bin/cgi_get_portrait.fcg?uins="
	getQQInfoFailed = "get qq info failed"
	qqInfoNotFound  = "qq info not found"
)

type QQInfoForm struct {
	QQNum string `form:"qq" json:"qq" xml:"qq" binding:"required"`
}

type QQInfoResponse struct {
	HeadPortrait string `json:"head_portrait"` // 头像
	NickName     string `json:"nickname"`      // 昵称
}

func QueryQQInfo(form *QQInfoForm) (resp *QQInfoResponse, msg string, err error) {
	data, err := getQQInfoData(form.QQNum)
	if err != nil {
		msg = getQQInfoFailed
		return
	}

	response, status, err := parseQQInfoData(data, form.QQNum)

	if err != nil {
		msg = getQQInfoFailed
		return
	}
	if !status {
		err = errors.New(qqInfoNotFound)
		msg = qqInfoNotFound
	}
	resp = response
	return
}

func getQQInfoData(qqNum string) (data []byte, err error) {
	apiUrl := apiQQUrl + qqNum
	data, _, err = utils.HttpProxyGet(apiUrl, "", nil)
	if err != nil {
		return
	}
	return
}

func parseQQInfoData(data []byte, qqNum string) (resp *QQInfoResponse, status bool, err error) {
	strData := mahonia.NewDecoder("gbk").ConvertString(string(data))
	if strings.HasPrefix(strData, "portraitCallBack") { // 成功
		s := strings.Replace(strData, "portraitCallBack(", "", 1)
		s = strings.Replace(s, ")", "", 1)
		var d interface{}
		err = json.Unmarshal([]byte(s), &d)
		if err != nil {
			return
		}
		resp = &QQInfoResponse{}
		var dmap map[string]interface{}
		if dm, ok := d.(map[string]interface{}); !ok {
			err = errors.New("data parse failed")
			return
		} else {
			dmap = dm
		}

		if qqInfo, ok2 := dmap[qqNum]; ok2 {
			if qqInfoSlice, ok3 := qqInfo.([]interface{}); ok3 {
				status = true
				resp.HeadPortrait = qqInfoSlice[0].(string)
				resp.NickName = qqInfoSlice[6].(string)
			} else {
				err = errors.New("data parse failed")
				return
			}
		} else {
			return
		}
	}
	return
}
