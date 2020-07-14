package controlers

import (
	"apiTools/libs/logger"
	"apiTools/modles"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// ipv4信息查询 --> api
func Ipv4Query(c *gin.Context) {
	var ipv4Form modles.Ipv4Form
	err := c.Bind(&ipv4Form)
	// 获取参数信息失败
	if err != nil {
		// 成功状态码(0 成功, 非零 失败)
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "param has fail", "data": gin.H{}})
		logger.Echo.WithFields(logrus.Fields{
			"routers": c.Request.URL.Path,
			"err":    err,
			"query":  ipv4Form,
		}).Error("get param fields fail from ipv4 info query")
		return
	}
	ipv4Info, err := modles.Ipv4Query(ipv4Form)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "msg": "query fail", "data": gin.H{}})
		logger.Echo.WithFields(logrus.Fields{
			"routers": c.Request.URL.Path,
			"err":    err,
			"query":  ipv4Form,
		}).Error("query ipv4 info fail")
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": ipv4Info, "msg": ""})
}
