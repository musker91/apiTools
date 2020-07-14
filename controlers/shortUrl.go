package controlers

import (
	"apiTools/libs/logger"
	"apiTools/modles"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// 长链接转换为短链接 --> api
func ShortToShortUrl(c *gin.Context) {
	shortForm := &modles.ShortForm{}
	data := gin.H{
		"code":     1,       // 转换成功状态码(0 成功, 非零 失败)
		"domain":   gin.H{}, // 短地址配置的域名
		"shortUrl": "",      // 短链接地址
	}
	err := c.Bind(shortForm)
	// 接收数据失败
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "data": data, "msg": "Incorrect request data"})
		logger.Echo.WithFields(logrus.Fields{
			"routers": c.Request.URL.Path,
			"err":     err,
			"query":   shortForm,
		}).Error("get query to short url form param fail")
		return
	}
	// 设置默认域名
	if shortForm.Domain == "" {
		shortForm.Domain = c.Request.Host
	}
	// 设置默认过期时间
	if shortForm.ExpireTime == 0 {
		shortForm.ExpireTime = -1
	}

	shortInfo, err := modles.ToShortUrl(shortForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "data": data, "msg": "Short link generation failed, please try again later!!!"})
		logger.Echo.WithFields(logrus.Fields{
			"routers": c.Request.URL.Path,
			"err":     err,
			"query":   shortForm,
			"data":    shortInfo,
		}).Error("to short url fail")
		return
	}
	// 赋值
	data["code"] = 0
	data["domain"] = shortInfo.Domain
	data["shortUrl"] = fmt.Sprintf("http://%s/%s", shortInfo.Domain, shortInfo.ShortStr)

	// log info
	logger.Echo.WithFields(logrus.Fields{
		"routers": c.Request.URL.Path,
		"warn":    err,
		"query":   shortForm,
		"data":    shortInfo,
	}).Info("to short url success")

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": data, "msg": ""})
}

// 短链接解析回长链接 --> api
func ShortParseShortUrl(c *gin.Context) {
	shortUrlQuery := c.Query("shortUrl")
	data := gin.H{
		"code":    1,  // 转换成功状态码(0 成功, 非零 失败)
		"domain":  "", // 短地址配置的域名
		"longUrl": "", // 原长链接地址
	}
	// 获取请求参数失败
	if shortUrlQuery == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "data": data, "msg": "shortUrl parameter cannot be empty",})
		logger.Echo.WithFields(logrus.Fields{
			"routers": c.Request.URL.Path,
			"err":     data["msg"],
			"query":   "",
		}).Error("get query param short url form param fail")
		return
	}
	// 解析短链接
	shortInfo, err := modles.ParseShort(shortUrlQuery)
	// 解析失败
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "data": data, "msg": "parse short url fail"})
		logger.Echo.WithFields(logrus.Fields{
			"routers": c.Request.URL.Path,
			"err":     err,
			"query":   shortUrlQuery,
			"data":    shortInfo,
		}).Error("parse short url fail")
		return
	}
	// 解析成功, 赋值
	data["code"] = 0
	data["domain"] = shortInfo.Domain
	data["longUrl"] = shortInfo.LongUrl

	// log info
	logger.Echo.WithFields(logrus.Fields{
		"routers": c.Request.URL.Path,
		"warn":    err,
		"query":   shortUrlQuery,
		"data":    shortInfo,
	}).Info("parse short url success")

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": data, "msg": ""})
}

// 转换为三方短链接 --> api
func ShortToOfficial(c *gin.Context) {
	var shortForm modles.ShortForm
	if err := c.Bind(&shortForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "data": gin.H{}, "msg": "shortUrl parameter cannot be empty",})
		logger.Echo.WithFields(logrus.Fields{
			"routers": c.Request.URL.Path,
			"err":     err,
			"query":   "",
		}).Error("get query param short url form param fail")
		return
	}
	shortInfo, msg, err := modles.OfficialToShort(&shortForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "data": gin.H{}, "msg": msg,})
		logger.Echo.WithFields(logrus.Fields{
			"routers": c.Request.URL.Path,
			"err":     err,
			"query":   "",
		}).Error(msg)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{
		"longurl":  shortInfo.LongUrl,
		"shorturl": shortInfo.ShortStr,
	}, "msg": ""})
}

// 解析三方短链接 --> api
func ShortParseOfficial(c *gin.Context) {
	var shortForm modles.ShortForm
	if err := c.Bind(&shortForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "data": gin.H{}, "msg": "shortUrl parameter cannot be empty",})
		logger.Echo.WithFields(logrus.Fields{
			"routers": c.Request.URL.Path,
			"err":     err,
			"query":   "",
		}).Error("get query param short url form param fail")
		return
	}
	shortInfo, msg, err := modles.OfficialParseShort(&shortForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "data": gin.H{}, "msg": msg,})
		logger.Echo.WithFields(logrus.Fields{
			"routers": c.Request.URL.Path,
			"err":     err,
			"query":   "",
		}).Error(msg)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{
		"longurl":  shortInfo.LongUrl,
		"shorturl": shortInfo.ShortStr,
	}, "msg": ""})
}
