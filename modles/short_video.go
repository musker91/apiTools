package modles

import (
	"apiTools/utils"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/tidwall/gjson"
	"regexp"
	"strings"
)

// 短视频解析

// 请求结构体
type ShortVideoForm struct {
	Url string `form:"url" json:"url" xml:"url" binding:"required"`
}

// 响应结构体
type ShortVideoResult struct {
	Desc  string `json:"desc"`
	Pic   string `json:"pic"`
	Video string `json:"video"`
	Music string `json:"music"`
}

type parseInterface interface {
	parse(url string) (result *baseShortVideo, err error)
}

type baseShortVideo struct {
	desc  string
	pic   string
	video string
	music string
}

const (
	shortUrlPattern = `^https?:\/\/(([a-zA-Z0-9_-])+(\.)?)*(:\d+)?(\/((\.)?(\?)?=?&?[a-zA-Z0-9_-](\?)?)*)*$`
	matchDomain     = `^((http://)|(https://))?([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,6}(/)`
)

var (
	shortVideoDomain = map[string]string{
		"v.douyin.com":      "douyin",
		"v.kuaishou.com":    "kuaishou",
		"h5.pipix.com":      "pipixia",
		"h5.weishi.qq.com":  "weishi",
		"share.izuiyou.com": "zuiyou",
		"share.huoshan.com": "huoshan",
	}
	iosHeader = map[string]string{"User-Agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1"}
	androidHeader = map[string]string{"User-Agent":"Mozilla/5.0 (Linux; Android 8.0.0; Pixel 2 XL Build/OPD1.170816.004) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.25 Mobile Safari/537.36"}
	)

// 解析入口
func ShortVideoParse(form *ShortVideoForm) (result *ShortVideoResult, err error) {
	// 判断是不是url
	matched, err := regexp.Match(shortUrlPattern, []byte(form.Url))
	if err != nil {
		return
	}
	if !matched {
		err = errors.New("from url param not url address")
		return
	}
	// 匹配域名
	matchShortUrl, err := regexp.Compile(matchDomain)
	if err != nil {
		return
	}

	allMatchShortUrl := matchShortUrl.FindAllString(form.Url, 1)
	if len(allMatchShortUrl) == 0 {
		err = errors.New("not match url")
		return
	}

	shortUrl := strings.Trim(allMatchShortUrl[0], "/")
	shortUrl = strings.TrimPrefix(shortUrl, "https://")
	shortUrl = strings.TrimPrefix(shortUrl, "http://")

	//  判断是那种类型短视频
	shortType, ok := shortVideoDomain[shortUrl]
	if !ok {
		err = errors.New("unsupported short video type")
		return
	}

	var parseObj parseInterface
	switch shortType {
	case "douyin":
		parseObj = &douyinVideo{}
	case "kuaishou":
		//parseObj = &kuaishouVideo{}
	case "weishi":
		parseObj = &weishiVideo{}
	case "pipixia":
		parseObj = &pipixiaVideo{}
	case "zuiyou":
		parseObj = &zuiyouVideo{}
	case "huoshan":
		parseObj = &huoshanVideo{}

	}
	if parseObj == nil {
		err = errors.New("parse short url interface is nil")
		return
	}

	resultData, err := parseObj.parse(form.Url)
	if err != nil {
		return
	}
	if resultData == nil  {
		err = fmt.Errorf("parse %s data return is nil", shortType)
		return
	}

	result = &ShortVideoResult{
		Desc:  resultData.desc,
		Pic:   resultData.pic,
		Video: resultData.video,
		Music: resultData.music,
	}
	return
}

// 抖音
type douyinVideo struct {
	baseShortVideo
	mid  string
	dytk string
}

func (this *douyinVideo) parse(url string) (result *baseShortVideo, err error) {
	data, err := this.getUrlData(url)
	if err != nil {
		return
	}
	err = this.parseUrlData(data)
	if err != nil {
		return
	}
	_ = this.getDescAndMusic()

	result = &this.baseShortVideo
	return
}

func (this *douyinVideo) getUrlData(url string) (data []byte, err error) {
	redirectUrl, _, err := utils.GetRedirectUrl(url, "", nil)
	if err != nil {
		return
	}

	// 匹配mid
	this.matchMid(redirectUrl)

	data, _, err = utils.HttpProxyGet(redirectUrl, "", nil)
	if err != nil {
		return
	}
	return
}

func (this *douyinVideo) parseUrlData(data []byte) (err error) {
	// video
	videoMatchGroup, err := utils.RegexMatchGroup(`playAddr: "(?P<playAddr>.*)"`, string(data))
	if err != nil {
		return
	}
	playAddr, playAddrOk := videoMatchGroup["playAddr"]
	if !playAddrOk {
		err = errors.New("not match douyin play address")
		return
	}
	newPlayAddr := strings.Replace(playAddr, "playwm", "play", 1)

	this.video, _, err = utils.GetRedirectUrl(newPlayAddr, "", iosHeader)
	if err != nil {
		return
	}
	// pic
	picMatchGroup, err := utils.RegexMatchGroup(`cover: "(?P<pic>.*)"`, string(data))
	if err != nil {
		return
	}
	pic, picMatchGroupOk := picMatchGroup["pic"]
	if !picMatchGroupOk {
		err = errors.New("not match douyin pic image")
		return
	}
	this.pic = pic
	// dytk
	dytkMatchGroup, err := utils.RegexMatchGroup(`dytk: "(?P<dytk>.*)"`, string(data))
	if err != nil {
		return
	}
	dytk, dytkMatchGroupOk := dytkMatchGroup["dytk"]
	if !dytkMatchGroupOk {
		err = errors.New("not match douyin dytk")
		return
	}
	this.dytk = dytk

	return
}

func (this *douyinVideo) getDescAndMusic() (err error) {
	if this.mid == "" || this.dytk == "" {
		err = errors.New("douyin request mid or dytk is null")
		return
	}
	url := fmt.Sprintf("https://www.iesdouyin.com/web/api/v2/aweme/iteminfo/?item_ids=%s&dytk=%s",
		this.mid, this.dytk)

	data, _, err := utils.HttpProxyGet(url, "", nil)
	if err != nil {
		return
	}
	var rspdData map[string]interface{}
	err = json.Unmarshal(data, &rspdData)
	if err != nil {
		return
	}

	status_code, status_code_ok := rspdData["status_code"].(float64)
	if status_code_ok && int(status_code) != 0 {
		err = errors.New("douyin get desc and music status code has error")
		return
	}

	item_list, item_list_ok := rspdData["item_list"].([]interface{})
	if !item_list_ok {
		err = errors.New("douyin get desc and music item_list has error")
		return
	}
	if len(item_list) == 0 {
		err = errors.New("douyin get desc and music item_list is null")
		return
	}
	item := item_list[0]
	item_data, item_data_ok := item.(map[string]interface{})
	if !item_data_ok {
		err = errors.New("douyin get desc and music item 0 has error")
		return
	}
	// desc
	desc, descOk := item_data["desc"].(string)
	if !descOk {
		err = errors.New("douyin get desc and music item desc has error")
		return
	}
	this.desc = desc

	// music
	musicList, musicListOK := item_data["music"].(map[string]interface{})
	if !musicListOK {
		err = errors.New("douyin get desc and music music list has error")
		return
	}
	play_url, play_url_ok := musicList["play_url"].(map[string]interface{})
	if !play_url_ok {
		err = errors.New("douyin get desc and music music list play_url has error")
		return
	}

	play_url_url, play_url_url_ok := play_url["uri"].(string)
	if !play_url_url_ok {
		err = errors.New("douyin get desc and music music list play_url url has error")
		return
	}

	this.music = play_url_url
	return
}

func (this *douyinVideo) matchMid(redirectUrl string) {
	result, err := utils.RegexMatchGroup("video/(?P<mid>.*)/", redirectUrl)
	if err != nil {
		return
	}
	mid, midOK := result["mid"]
	if !midOK {
		return
	}
	this.mid = mid
}

// 皮皮虾
type pipixiaVideo struct {
	baseShortVideo
	item string
}

func (this *pipixiaVideo) parse(url string) (result *baseShortVideo, err error) {
	err = this.getItem(url)
	if err != nil {
		return
	}
	data, err := this.getVideoData()
	if err != nil {
		return
	}
	err = this.parseResultData(data)
	if err != nil {
		return
	}

	result = &this.baseShortVideo

	return
}

func (this *pipixiaVideo) getItem(url string) (err error) {
	redirectUrl,_,  err := utils.GetRedirectUrl(url, "", nil)
	if err != nil {
		return
	}
	result, err := utils.RegexMatchGroup(`/item/(?P<item>.*)\?`, redirectUrl)
	if err != nil {
		return
	}
	item, itemOk := result["item"]
	if !itemOk {
		err = errors.New("get item is fail, he is null")
		return
	}
	this.item = item
	return
}

func (this *pipixiaVideo) getVideoData() (data []byte, err error) {
	url := fmt.Sprintf("https://is.snssdk.com/bds/item/detail/?app_name=super&aid=1319&item_id=%s", this.item)
	data, _, err = utils.HttpProxyGet(url, "", nil)
	if err != nil {
		return
	}
	return
}

func (this *pipixiaVideo) parseResultData(data []byte) (err error) {
	gjsonObj := gjson.ParseBytes(data)
	status_code := gjsonObj.Get("status_code").Int()
	if status_code != 0 {
		err = errors.New("get pipixia video url data fail, status_code error")
		return
	}
	this.desc = gjsonObj.Get("data.data.content").String()
	this.pic = gjsonObj.Get("data.data.cover.url_list").Array()[0].Get("url").String()
	this.video = gjsonObj.Get("data.data.video.video_fallback.url_list").Array()[0].Get("url").String()
	return
}

// 微视
type weishiVideo struct {
	baseShortVideo
	feedid string
}

func (this *weishiVideo) parse(url string) (result *baseShortVideo, err error) {
	err = this.getFeedid(url)
	if err != nil {
		return
	}
	err = this.getAndParseData()
	if err != nil {
		return
	}
	result = &this.baseShortVideo
	return
}

func (this *weishiVideo) getFeedid(url string) (err error) {
	result, err := utils.RegexMatchGroup(`/feed/(?P<feedid>.*)/`, url)
	if err != nil {
		return
	}
	feedid, ok := result["feedid"]
	if !ok {
		err = errors.New("weishi video get feedid is null")
		return
	}
	this.feedid = feedid
	return
}

func (this *weishiVideo) getAndParseData() (err error) {
	url := fmt.Sprintf("https://h5.qzone.qq.com/webapp/json/weishi/WSH5GetPlayPage?feedid=%s", this.feedid)
	data, _, err := utils.HttpProxyGet(url, "", nil)
	if err != nil {
		return
	}

	vData := gjson.ParseBytes(data)
	if vData.Get("ret").Int() != 0 {
		err = errors.New("get weisi video api data fail, status not is 0")
		return
	}

	feeds := vData.Get("data.feeds").Array()[0]
	this.pic = feeds.Get("images").Array()[0].Get("url").String()
	this.desc = feeds.Get("feed_desc").String()
	this.video = feeds.Get("video_url").String()
	return
}

// 最右
type zuiyouVideo struct {
	baseShortVideo
	pid string
}

func (this *zuiyouVideo) parse(url string) (result *baseShortVideo, err error) {
	err = this.getPid(url)
	if err != nil {
		return
	}
	err = this.getData()
	if err != nil {
		return
	}
	result = &this.baseShortVideo
	return
}

func (this *zuiyouVideo) getPid(url string) (err error) {
	result, err := utils.RegexMatchGroup(`/detail/(?P<pid>\d+)?`, url)
	if err != nil {
		return
	}
	pid, pidOk := result["pid"]
	if !pidOk {
		err = fmt.Errorf("not match zuiyou url pid, %s", url)
		return
	}

	this.pid = pid
	return
}

func (this *zuiyouVideo) getData() (err error) {
	url := fmt.Sprintf("https://share.izuiyou.com/hybrid/share/post?pid=%s", this.pid)
	data, _, err := utils.HttpProxyGet(url, "", nil)
	if err != nil {
		return
	}

	rootNode, err := htmlquery.Parse(bytes.NewReader(data))
	if err != nil {
		return
	}
	jsData := htmlquery.FindOne(rootNode, "//script[@id='appState']")
	jsText := htmlquery.InnerText(jsData)
	jsText = strings.TrimPrefix(jsText, "window.APP_INITIAL_STATE=")
	result := gjson.Parse(jsText)
	this.desc = result.Get("sharePost.postDetail.post.content").String()
	this.pic = result.Get("sharePost.postDetail.post.imgs").Array()[0].
		Get("urls.540Webp.urls").Array()[0].String()
	videos := result.Get("sharePost.postDetail.post.videos").Map()
	var keys []string
	for key, _ := range videos {
		keys = append(keys, key)
	}
	this.video = videos[keys[0]].Get("url").String()
	return
}


// 火山
type huoshanVideo struct {
	baseShortVideo
	itemId string
	videoId string
}

func (this *huoshanVideo) parse(url string) (result *baseShortVideo, err error) {
	err = this.getItemId(url)
	if err != nil {
		return
	}
	err = this.getPic()
	if err != nil {
		return
	}
	err = this.getVideoUrl()
	if err != nil {
		return
	}
	result = &this.baseShortVideo
	return
}

func (this *huoshanVideo) getItemId(url string) (err error) {
	redirectUrl,_,  err := utils.GetRedirectUrl(url, "", nil)
	if err != nil {
		return
	}
	result, err := utils.RegexMatchGroup(`item_id=(?P<itemId>\d+)&`, redirectUrl)
	if err != nil {
		return
	}
	itemId, itemIdOk := result["itemId"]
	if !itemIdOk {
		err = errors.New("get huoshan video item id fail")
		return
	}
	this.itemId = itemId
	return
}

func (this *huoshanVideo) getPic() (err error) {
	url := fmt.Sprintf("https://share.huoshan.com/api/item/info?item_id=%s", this.itemId)
	data, _, err := utils.HttpProxyGet(url, "", nil)
	if err != nil {
		return
	}
	result := gjson.ParseBytes(data)
	if result.Get("status_code").Int() != 0 {
		err = errors.New("get huoshan video pic data fail")
		return
	}
	this.pic = result.Get("data.item_info.cover").String()
	err = this.getVideoId(result.Get("data.item_info.url").String())
	if err != nil {
		return
	}
	return
}

func (this *huoshanVideo) getVideoId(videoUrl string) (err error) {
	if videoUrl == "" {
		err = errors.New("get video url is nil")
		if err != nil {
			return
		}
	}
	result, err := utils.RegexMatchGroup(`video_id=(?P<videoId>\w+)&`, videoUrl)
	if err != nil {
		return
	}
	videoId, videoIdOK := result["videoId"]
	if !videoIdOK {
		err = errors.New("get huoshan video video_id fail")
		return
	}
	this.videoId = videoId
	return
}

func (this *huoshanVideo) getVideoUrl() (err error) {
	fmt.Println(this.videoId)
	url := fmt.Sprintf("http://hotsoon.snssdk.com/hotsoon/item/video/_playback/?video_id=%s", this.videoId)
	redirectUrl,_,  err := utils.GetRedirectUrl(url, "", nil)
	if err != nil {
		return
	}
	this.video = redirectUrl
	return
}

// 快手，目前有问题，暂无法解析
type kuaishouVideo struct {
	baseShortVideo
}

func (this *kuaishouVideo) parse(url string) (result *baseShortVideo, err error) {
	err = this.getData(url)
	if err != nil {
		return
	}
	result = &this.baseShortVideo
	return
}

func (this *kuaishouVideo) getData(url string) (err error) {
	redirectUrl, _ ,  err := utils.GetRedirectUrl(url, "", nil)
	if err != nil {
		return
	}

	// soup = BeautifulSoup(share_response,'lxml')
	//noWaterMarkVideo = soup.find(attrs={'id': 'hide-pagedata'}).attrs['data-pagedata']
	//
	//print(noWaterMarkVideo)
	//
	//正则处理字符串获取真实地址
	//pattern = re.compile('\"srcNoMark\":"(.*?)"},',re.S)
	//
	//real_url = re.findall(pattern,noWaterMarkVideo)[0]
	// 这块有一个安全验证问题
	data, _, err := utils.HttpProxyGet(redirectUrl, "", iosHeader)
	if err != nil {
		return
	}
	fmt.Println(string(data))
	return
}
