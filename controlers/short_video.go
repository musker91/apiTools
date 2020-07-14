package controlers

import (
	"apiTools/libs/logger"
	"apiTools/modles"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func ShortVideoParse(c *gin.Context) {
	var requestForm modles.ShortVideoForm
	err := c.Bind(&requestForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "param has fail", "data": gin.H{}})
		logger.Echo.WithFields(logrus.Fields{
			"routers": c.Request.URL.Path,
			"err":     err,
			"query":   requestForm,
		}).Error("get param ShortVideoParse fields from short video parse")
		return
	}
	result, err := modles.ShortVideoParse(&requestForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "parse short video fail", "data": gin.H{}})
		logger.Echo.WithFields(logrus.Fields{
			"routers": c.Request.URL.Path,
			"err":     err,
			"query":   requestForm,
		}).Error("parse short video fail")
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "", "data": result})
}
