package controlers

import (
	"apiTools/libs/logger"
	"apiTools/modles"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// 域名whois信息查询 --> api
func WhoisQuery(c *gin.Context) {
	var whoisForm modles.WhoisForm
	err := c.Bind(&whoisForm)
	data := gin.H{
		"data":   gin.H{}, // whois数据
		"status": 5,  // 域名查询状态
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "data": data, "msg": "request param fail"})
		logger.Echo.WithFields(logrus.Fields{
			"routers": c.Request.URL.Path,
			"err":    err,
			"query":  whoisForm,
		}).Error("get query whois form param fail")
		return
	}
	whoisInfo := &modles.WhoisInfo{}
	if whoisForm.OutType == "text" {
		whoisInfo.WhoisForm.OutType = "text"
		whoisInfo, err = modles.QueryWhoisInfo(&whoisForm)
		data["status"] = whoisInfo.Status
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 0, "data": data, "msg": "query fail"})
			logger.Echo.WithFields(logrus.Fields{
				"routers": c.Request.URL.Path,
				"err":    err,
				"query":  whoisForm,
				//"data":   whoisInfo,
			}).Error("query whois info fail")
			return
		}
		data["data"] = whoisInfo.TextInfo
	} else {
		whoisInfo.WhoisForm.OutType = "json"
		whoisInfo, err = modles.QueryWhoisInfoToJson(&whoisForm)
		data["status"] = whoisInfo.Status
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 0, "data": data, "msg": "query fail"})
			logger.Echo.WithFields(logrus.Fields{
				"routers": c.Request.URL.Path,
				"err":    err,
				"query":  whoisForm,
				//"data":   whoisInfo,
			}).Error("query whois info fail")
			return
		}
		data["data"] = whoisInfo.JsonInfo
	}
	// log debug
	logger.Echo.WithFields(logrus.Fields{
		"routers": c.Request.URL.Path,
		"warn":   err,
		"query":  whoisForm,
		"data":   whoisInfo,
	}).Debug("query whois info success")
	// log info
	logger.Echo.WithFields(logrus.Fields{
		"routers": c.Request.URL.Path,
		"warn":   err,
		"query":  whoisForm,
		"data":   logrus.Fields{"domain": whoisInfo.Domain},
	}).Info("query whois info success")

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": data, "msg": ""})
}
