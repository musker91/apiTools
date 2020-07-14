package controlers

import (
	"apiTools/libs/logger"
	"apiTools/modles"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func QQInfoQueryInfo(c *gin.Context) {
	var requestForm modles.QQInfoForm
	err := c.Bind(&requestForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "param has fail", "data": gin.H{}})
		logger.Echo.WithFields(logrus.Fields{
			"routers": c.Request.URL.Path,
			"err":     err,
			"query":   requestForm,
		}).Error("get param fields fail from get qq info query")
		return
	}
	result, msg, err := modles.QueryQQInfo(&requestForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": msg, "data": gin.H{}})
		logger.Echo.WithFields(logrus.Fields{
			"routers": c.Request.URL.Path,
			"err":     err,
			"query":   requestForm,
		}).Error("query qq info fail")
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "", "data": result})
}
