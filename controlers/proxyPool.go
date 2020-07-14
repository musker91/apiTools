package controlers

import (
	"apiTools/libs/logger"
	"apiTools/modles"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// proxy pool信息查询 --> api
func ProxyPoolQuery(c *gin.Context) {
	var proxyPoolForm modles.ProxyPoolForm
	err := c.Bind(&proxyPoolForm)
	// 获取参数信息失败
	if err != nil {
		// 成功状态码(0 成功, 非零 失败)
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "request param error", "data": gin.H{}})
		logger.Echo.WithFields(logrus.Fields{
			"routers": c.Request.URL.Path,
			"err":     err,
			"query":   proxyPoolForm,
		}).Error("get proxy proxy info fail")
		return
	}
	proxyPoolResult, err := modles.QueryProxyPoolInfo(&proxyPoolForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "msg": "get fail", "data": gin.H{}})
		logger.Echo.WithFields(logrus.Fields{
			"routers": c.Request.URL.Path,
			"err":     err,
			"query":   proxyPoolForm,
		}).Error("get proxy proxy info fail")
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": proxyPoolResult, "msg": ""})
}
