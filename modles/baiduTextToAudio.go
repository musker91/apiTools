package modles

import (
	"apiTools/utils"
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

/*
百度文字转语音合成接口

http://tts.baidu.com/text2audio?lan=zh&ie=UTF-8&spd=2&text=要转换的文字

lan=zh：语言是中文，如果改为lan=en，则语言是英文。

ie=UTF-8：文字格式。

spd=2：语速，可以是1-9的数字，数字越大，语速越快。

text=**：这个就是你要转换的文字。
*/

const (
	bdApiAudioUrl    = "http://tts.baidu.com/text2audio"
	conversionFailed = "text conversion audio failed"
)

type AudioForm struct {
	Lang    string `form:"lang" json:"lang" xml:"lang" binding:"required"`          // 语言
	Charset string `form:"charset" json:"charset" xml:"charset" binding:"required"` // 文件编码格式
	Speed   int    `form:"speed" json:"speed" xml:"speed" binding:"required"`       // 语速
	Text    string `form:"text" json:"text" xml:"text" binding:"required"`          // 要转换的文字
}

type AudioResponse struct {
	//Body          io.Reader
	Data          []byte
	FileName      string
	ConTextType   string
	ContentLength int64
	Msg           string // 消息
}

func BdTextToAudio(form *AudioForm) (resp *AudioResponse, err error) {
	resp = &AudioResponse{}
	apiUrl := fmt.Sprintf("%s?lan=%s&ie=%s&spd=%d&text=%s",
		bdApiAudioUrl, form.Lang, form.Charset, form.Speed, form.Text)
	data, response, err := utils.HttpProxyGet(apiUrl, "", nil)
	if err != nil {
		resp.Msg = conversionFailed
		return
	}
	conTextType := response.Header.Get("Content-Type")
	resp.ConTextType = conTextType
	if conTextType == "audio/x-bd-bv" { // 转换成功
		resp.FileName = strings.ToLower(utils.GetShortStr())
		//resp.Body = response.Body
		resp.Data = data
		resp.ContentLength = response.ContentLength
		return
	}
	err = errors.New(conversionFailed)
	resp.Msg = conversionFailed
	return
}
