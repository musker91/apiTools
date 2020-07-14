package controlers

import (
	"apiTools/libs/logger"
	"apiTools/modles"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func TextSegCut(c *gin.Context) {
	var requestForm modles.SegForm
	err := c.Bind(&requestForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "param has fail", "data": gin.H{}})
		logger.Echo.WithFields(logrus.Fields{
			"routers": c.Request.URL.Path,
			"err":     err,
			"query":   requestForm,
		}).Error("get param segmentation fields from text segmentation cut")
		return
	}
	result, msg, err := modles.SegTextCut(&requestForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": msg, "data": gin.H{}})
		logger.Echo.WithFields(logrus.Fields{
			"routers": c.Request.URL.Path,
			"err":     err,
			"query":   requestForm,
		}).Error("text segmentation cut fail")
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "", "data": result})
}
